package loadtest

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	outDirPerm   = 0o750
	outFilePerm  = 0o600
	progressStep = 25
)

var unsafeFilenameChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

// Options are all inputs of a complete load test: ID scrape, run and report.
type Options struct {
	Config Config
	// IDHost is the host used to scrape the service IDs; defaults to the first host.
	IDHost  string
	IDCache IDCacheOptions
	// LimitServices restricts the test to the first N services; 0 means all.
	LimitServices int
	Label         string
	OutputDir     string
	// WANIP overrides the lookup; WANIPLookupURL empty disables the lookup.
	WANIP          string
	WANIPLookupURL string
	// Progress receives human readable progress lines; may be nil.
	Progress io.Writer
}

// Outcome describes where the results were written.
type Outcome struct {
	Dir     string
	Summary Summary
	Meta    RunMeta
}

// Execute scrapes the service IDs, runs the load test and writes results.jsonl, summary.json
// and report.html to a new directory in the output directory.
func Execute(ctx context.Context, opts Options) (Outcome, error) {
	progress := func(format string, args ...any) {
		if opts.Progress != nil {
			_, _ = fmt.Fprintf(opts.Progress, format+"\n", args...)
		}
	}

	idHost := opts.IDHost
	if idHost == "" {
		idHost = opts.Config.Hosts[0]
	}

	idsStart := time.Now()

	services, fromCache, err := LoadServiceRecords(ctx, idHost, opts.IDCache)
	if err != nil {
		return Outcome{}, fmt.Errorf("fetching service IDs: %w", err)
	}

	idsTook := time.Since(idsStart)

	if opts.LimitServices > 0 && opts.LimitServices < len(services) {
		services = services[:opts.LimitServices]
	}

	progress(
		"Service records: %d (%s, %s)",
		len(services),
		cacheLabel(fromCache),
		idsTook.Round(time.Millisecond),
	)

	wan := opts.WANIP
	if wan == "" && opts.WANIPLookupURL != "" {
		wan = LookupWANIP(ctx, opts.WANIPLookupURL)
	}

	if wan == "" {
		wan = "not looked up"
	}

	dir := filepath.Join(opts.OutputDir, fmt.Sprintf(
		"%s-%s", time.Now().Format("20060102-150405"), sanitize(opts.Label),
	))
	if err := os.MkdirAll(dir, outDirPerm); err != nil {
		return Outcome{}, err
	}

	logFile, err := os.OpenFile( //nolint:gosec // path is built from the user supplied output dir
		filepath.Join(dir, "results.jsonl"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, outFilePerm,
	)
	if err != nil {
		return Outcome{}, err
	}

	defer func() { _ = logFile.Close() }()

	writer := bufio.NewWriter(logFile)
	encoder := json.NewEncoder(writer)

	var (
		results []Result
		failed  int
	)

	progress(
		"Running against %s with %d workers...",
		strings.Join(opts.Config.Hosts, ", "),
		opts.Config.Workers,
	)

	info, err := Run(ctx, opts.Config, services, func(res Result) {
		results = append(results, res)
		if !res.OK {
			failed++
		}

		if err := encoder.Encode(res); err != nil {
			slog.Warn("could not write result", "err", err)
		}

		if len(results)%progressStep == 0 {
			progress("  %d requests done, %d failed", len(results), failed)
		}
	})
	if err != nil {
		return Outcome{}, err
	}

	if err := writer.Flush(); err != nil {
		return Outcome{}, err
	}

	meta := RunMeta{
		Label:        opts.Label,
		Hostname:     Hostname(),
		WANIP:        wan,
		Config:       opts.Config,
		Started:      info.Started,
		Finished:     info.Finished,
		Aborted:      info.Aborted,
		ServiceCount: len(services),
		IDsFromCache: fromCache,
		IDsFetchTime: idsTook,
	}
	summary := Summarize(results, info, opts.Config.Hosts)

	if err := writeOutputs(dir, meta, summary); err != nil {
		return Outcome{}, err
	}

	return Outcome{Dir: dir, Summary: summary, Meta: meta}, nil
}

type summaryFile struct {
	Meta    RunMeta `json:"meta"`
	Summary Summary `json:"summary"`
}

func writeOutputs(dir string, meta RunMeta, summary Summary) error {
	raw, err := json.MarshalIndent(summaryFile{meta, summary}, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "summary.json"), raw, outFilePerm); err != nil {
		return err
	}

	html, err := RenderReport(meta, summary)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "report.html"), []byte(html), outFilePerm)
}

// TextSummary returns a short plain text summary per host.
func TextSummary(out Outcome) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Run '%s' from %s (WAN IP %s), %s\n",
		out.Meta.Label, out.Meta.Hostname, out.Meta.WANIP,
		out.Meta.Finished.Sub(out.Meta.Started).Round(time.Second))
	fmt.Fprintf(&sb, "%-40s %8s %8s %8s %9s %9s %9s %7s\n",
		"host", "requests", "failed", "err%", "avg", "p95", "p99", "req/s")

	rows := append(
		[]GroupStats{{Name: "TOTAL", LatencyStats: out.Summary.Overall}},
		out.Summary.ByHost...)
	for _, row := range rows {
		fmt.Fprintf(&sb, "%-40s %8d %8d %7.2f%% %9s %9s %9s %7.2f\n",
			row.Name, row.Count, row.Errors, row.ErrorRate*100, //nolint:mnd // percentage
			formatMs(row.AvgMs), formatMs(row.P95Ms), formatMs(row.P99Ms), row.RequestsPS)
	}

	return sb.String()
}

func cacheLabel(fromCache bool) string {
	if fromCache {
		return "from cache"
	}

	return "scraped from CSW"
}

func sanitize(label string) string {
	label = unsafeFilenameChars.ReplaceAllString(strings.TrimSpace(label), "-")
	if label == "" {
		return "run"
	}

	return label
}
