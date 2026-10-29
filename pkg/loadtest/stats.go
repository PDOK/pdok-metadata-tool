package loadtest

import (
	"math"
	"sort"
	"strings"
	"time"
)

const (
	percentile50       = 0.50
	percentile90       = 0.90
	percentile95       = 0.95
	percentile99       = 0.99
	targetBuckets      = 60
	minBucket          = 5 * time.Second
	slowestCount       = 15
	maxDistinctErrors  = 25
	maxReportedFailure = 500
)

// LatencyStats summarises the response times (in ms) of a group of requests.
type LatencyStats struct {
	Count      int     `json:"count"`
	Errors     int     `json:"errors"`
	ErrorRate  float64 `json:"errorRate"`
	AvgMs      float64 `json:"avgMs"`
	MinMs      float64 `json:"minMs"`
	P50Ms      float64 `json:"p50Ms"`
	P90Ms      float64 `json:"p90Ms"`
	P95Ms      float64 `json:"p95Ms"`
	P99Ms      float64 `json:"p99Ms"`
	MaxMs      float64 `json:"maxMs"`
	AvgDNSMs   float64 `json:"avgDnsMs"`
	AvgConnMs  float64 `json:"avgConnectMs"`
	AvgTLSMs   float64 `json:"avgTlsMs"`
	AvgTTFBMs  float64 `json:"avgTtfbMs"`
	AvgBytes   float64 `json:"avgBytes"`
	RequestsPS float64 `json:"requestsPerSecond"`
}

// GroupStats is LatencyStats for a named group (host, request kind, ...).
type GroupStats struct {
	LatencyStats

	Name string `json:"name"`
}

// ErrorGroup counts failures with the same host, category and message.
type ErrorGroup struct {
	Host     string `json:"host"`
	Category string `json:"category"`
	Phase    string `json:"phase,omitempty"`
	Message  string `json:"message"`
	Count    int    `json:"count"`
}

// Bucket is one time slice of the run, per host.
type Bucket struct {
	Offset time.Duration `json:"offsetNs"`
	Host   string        `json:"host"`
	Count  int           `json:"count"`
	Errors int           `json:"errors"`
	AvgMs  float64       `json:"avgMs"`
	MaxMs  float64       `json:"maxMs"`
}

// HistogramRow counts requests per latency range, per host.
type HistogramRow struct {
	Label  string         `json:"label"`
	Counts map[string]int `json:"counts"`
}

// Summary is the aggregated result of a run.
type Summary struct {
	Overall      LatencyStats   `json:"overall"`
	ByHost       []GroupStats   `json:"byHost"`
	ByKind       []GroupStats   `json:"byKind"`
	ByHostKind   []GroupStats   `json:"byHostKind"`
	ByRemoteIP   []GroupStats   `json:"byRemoteIp"`
	Errors       []ErrorGroup   `json:"errors"`
	ErrorsByCat  []GroupStats   `json:"errorsByCategory"`
	Timeline     []Bucket       `json:"timeline"`
	BucketSize   time.Duration  `json:"bucketSizeNs"`
	Histogram    []HistogramRow `json:"histogram"`
	Failed       []Result       `json:"failed"`
	FailedTotal  int            `json:"failedTotal"`
	Slowest      []Result       `json:"slowest"`
	Hosts        []string       `json:"hosts"`
	ElapsedSecs  float64        `json:"elapsedSeconds"`
	ServicesSeen int            `json:"servicesSeen"`
}

var histogramEdges = []struct {
	label string
	upper float64
}{
	{"< 100 ms", 100},
	{"100-250 ms", 250},
	{"250-500 ms", 500},
	{"0.5-1 s", 1000},
	{"1-2.5 s", 2500},
	{"2.5-5 s", 5000},
	{"5-10 s", 10000},
	{"10-20 s", 20000},
	{"> 20 s", math.Inf(1)},
}

// Summarize aggregates results. hosts determines the order of the hosts in the output.
func Summarize(results []Result, info RunInfo, hosts []string) Summary {
	elapsed := info.Finished.Sub(info.Started)
	sum := Summary{Hosts: hosts, ElapsedSecs: elapsed.Seconds()}

	sum.Overall = computeStats(results, elapsed)
	sum.ByHost = groupBy(results, elapsed, func(r Result) string { return r.Host }, hosts)
	sum.ByKind = groupBy(
		results, elapsed, func(r Result) string { return r.Kind },
		[]string{KindServiceRecord, KindDatasetTags, KindDatasetRecord},
	)
	sum.ByHostKind = groupBy(
		results,
		elapsed,
		func(r Result) string { return r.Host + " | " + r.Kind },
		nil,
	)
	sum.ByRemoteIP = groupBy(results, elapsed, func(r Result) string {
		if r.RemoteIP == "" {
			return ""
		}

		return r.Host + " -> " + r.RemoteIP
	}, nil)
	sum.ErrorsByCat = errorCategories(results)
	sum.Errors = errorGroups(results)
	sum.Timeline, sum.BucketSize = timeline(results, info, hosts)
	sum.Histogram = histogram(results)

	services := map[string]bool{}

	for _, r := range results {
		services[r.ServiceID] = true

		if !r.OK {
			sum.FailedTotal++

			if len(sum.Failed) < maxReportedFailure {
				sum.Failed = append(sum.Failed, r)
			}
		}
	}

	sum.ServicesSeen = len(services)
	sum.Slowest = slowest(results)

	return sum
}

func computeStats(results []Result, elapsed time.Duration) LatencyStats {
	var (
		stats     LatencyStats
		latencies []float64
		okCount   float64
	)

	stats.Count = len(results)

	for _, r := range results {
		if !r.OK {
			stats.Errors++

			continue
		}

		okCount++

		latencies = append(latencies, r.DurationMs)
		stats.AvgDNSMs += r.DNSMs
		stats.AvgConnMs += r.ConnectMs
		stats.AvgTLSMs += r.TLSMs
		stats.AvgTTFBMs += r.TTFBMs
		stats.AvgBytes += float64(r.Bytes)
	}

	if stats.Count > 0 {
		stats.ErrorRate = float64(stats.Errors) / float64(stats.Count)
	}

	if elapsed > 0 {
		stats.RequestsPS = float64(stats.Count) / elapsed.Seconds()
	}

	if okCount == 0 {
		return stats
	}

	stats.AvgDNSMs /= okCount
	stats.AvgConnMs /= okCount
	stats.AvgTLSMs /= okCount
	stats.AvgTTFBMs /= okCount
	stats.AvgBytes /= okCount

	sort.Float64s(latencies)

	total := 0.0
	for _, l := range latencies {
		total += l
	}

	stats.AvgMs = total / okCount
	stats.MinMs = latencies[0]
	stats.MaxMs = latencies[len(latencies)-1]
	stats.P50Ms = percentile(latencies, percentile50)
	stats.P90Ms = percentile(latencies, percentile90)
	stats.P95Ms = percentile(latencies, percentile95)
	stats.P99Ms = percentile(latencies, percentile99)

	return stats
}

// percentile uses the nearest-rank method on a sorted slice.
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}

	rank := int(math.Ceil(p*float64(len(sorted)))) - 1

	return sorted[min(max(rank, 0), len(sorted)-1)]
}

// groupBy groups results by key; an empty key skips the result. Names in order come first.
func groupBy(
	results []Result,
	elapsed time.Duration,
	key func(Result) string,
	order []string,
) []GroupStats {
	groups := map[string][]Result{}

	for _, r := range results {
		if k := key(r); k != "" {
			groups[k] = append(groups[k], r)
		}
	}

	names := make([]string, 0, len(groups))

	for _, name := range order {
		if _, ok := groups[name]; ok {
			names = append(names, name)
		}
	}

	var rest []string

	for name := range groups {
		if !contains(order, name) {
			rest = append(rest, name)
		}
	}

	sort.Strings(rest)

	names = append(names, rest...)

	out := make([]GroupStats, 0, len(names))
	for _, name := range names {
		out = append(out, GroupStats{Name: name, LatencyStats: computeStats(groups[name], elapsed)})
	}

	return out
}

func errorCategories(results []Result) []GroupStats {
	counts := map[string]int{}
	total := 0

	for _, r := range results {
		if !r.OK {
			counts[r.ErrorCategory]++
			total++
		}
	}

	out := make([]GroupStats, 0, len(counts))
	for name, count := range counts {
		out = append(out, GroupStats{Name: name, LatencyStats: LatencyStats{
			Count:     count,
			Errors:    count,
			ErrorRate: float64(count) / float64(total),
		}})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}

		return out[i].Name < out[j].Name
	})

	return out
}

// errorGroups groups failures by host, category and message, with the volatile parts
// (URLs, addresses with ports) of the message stripped so equal failures group together.
func errorGroups(results []Result) []ErrorGroup {
	type key struct{ host, category, phase, message string }

	counts := map[key]int{}

	for _, r := range results {
		if !r.OK {
			counts[key{r.Host, r.ErrorCategory, r.ErrorPhase, normalizeError(r)}]++
		}
	}

	out := make([]ErrorGroup, 0, len(counts))
	for k, count := range counts {
		out = append(out, ErrorGroup{
			Host: k.host, Category: k.category, Phase: k.phase, Message: k.message, Count: count,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}

		return out[i].Message < out[j].Message
	})

	if len(out) > maxDistinctErrors {
		out = out[:maxDistinctErrors]
	}

	return out
}

func normalizeError(r Result) string {
	msg := strings.ReplaceAll(r.Error, r.URL, "<url>")
	// Go wraps client errors as `Get "<url>": reason`; keep only the reason.
	if i := strings.LastIndex(msg, `": `); i >= 0 {
		msg = msg[i+len(`": `):]
	}

	return msg
}

func timeline(results []Result, info RunInfo, hosts []string) ([]Bucket, time.Duration) {
	if len(results) == 0 {
		return nil, 0
	}

	size := max((info.Finished.Sub(info.Started) / targetBuckets).Round(time.Second), minBucket)

	type acc struct {
		count, errors int
		sum, max      float64
		ok            int
	}

	grid := map[int]map[string]*acc{}
	last := 0

	for _, r := range results {
		idx := max(int(r.Time.Sub(info.Started)/size), 0)
		last = max(last, idx)

		if grid[idx] == nil {
			grid[idx] = map[string]*acc{}
		}

		a := grid[idx][r.Host]
		if a == nil {
			a = &acc{}
			grid[idx][r.Host] = a
		}

		a.count++

		if !r.OK {
			a.errors++

			continue
		}

		a.ok++
		a.sum += r.DurationMs
		a.max = math.Max(a.max, r.DurationMs)
	}

	var out []Bucket

	for idx := 0; idx <= last; idx++ {
		for _, host := range hosts {
			b := Bucket{Offset: time.Duration(idx) * size, Host: host}
			if a := grid[idx][host]; a != nil {
				b.Count, b.Errors, b.MaxMs = a.count, a.errors, a.max
				if a.ok > 0 {
					b.AvgMs = a.sum / float64(a.ok)
				}
			}

			out = append(out, b)
		}
	}

	return out, size
}

func histogram(results []Result) []HistogramRow {
	rows := make([]HistogramRow, len(histogramEdges))
	for i, e := range histogramEdges {
		rows[i] = HistogramRow{Label: e.label, Counts: map[string]int{}}
	}

	for _, r := range results {
		if !r.OK {
			continue
		}

		for i, e := range histogramEdges {
			if r.DurationMs < e.upper {
				rows[i].Counts[r.Host]++

				break
			}
		}
	}

	return rows
}

func slowest(results []Result) []Result {
	ok := make([]Result, 0, len(results))

	for _, r := range results {
		if r.OK {
			ok = append(ok, r)
		}
	}

	sort.Slice(ok, func(i, j int) bool { return ok[i].DurationMs > ok[j].DurationMs })

	return ok[:min(len(ok), slowestCount)]
}

func contains(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}

	return false
}
