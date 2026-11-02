package loadtest

import (
	"fmt"
	"strings"
)

const (
	levelGood = "good"
	levelWarn = "warn"
	levelBad  = "bad"

	warnRate     = 0.05
	minorRate    = 0.01
	slowTailFold = 10.0
	percent      = 100.0
)

// Verdict is a plain language conclusion of a run.
type Verdict struct {
	Level  string
	Status string
	Title  string
	Points []string
}

func levelForRate(rate float64) string {
	switch {
	case rate == 0:
		return levelGood
	case rate < warnRate:
		return levelWarn
	default:
		return levelBad
	}
}

func statusForLevel(level string) string {
	switch level {
	case levelGood:
		return "OK"
	case levelWarn:
		return "Some failures"
	default:
		return "Failing"
	}
}

// failureExplanations explain in plain language what a failure phase or category means.
var failureExplanations = map[string]string{
	ErrTLS: "The TLS handshake did not complete: the server accepted the TCP connection but did not " +
		"finish setting up the encrypted connection in time. That happens in the network edge " +
		"(firewall, load balancer, WAF), before the CSW application is involved.",
	ErrConnect: "A TCP connection could not be opened. Something between client and server " +
		"(firewall, load balancer, NAT) dropped or refused new connections.",
	"dns":           "The host name could not be resolved.",
	ErrTimeout:      "The server did not answer within the request timeout.",
	ErrConnReset:    "The connection was closed by the other side in the middle of the request.",
	ErrHTTPStatus:   "The server answered with an HTTP error status; the application itself returned an error.",
	ErrCSWException: "The CSW answered with an ExceptionReport.",
	ErrBodyRead:     "The connection broke while the response was being read.",
	ErrParse:        "The response was received but could not be parsed.",
	ErrOther:        "An unclassified error occurred.",
	ErrCanceled:     "The request was canceled.",
}

func explanationFor(group ErrorGroup) string {
	switch group.Phase {
	case ErrTLS, ErrConnect, ErrDNS:
		return failureExplanations[group.Phase]
	default:
		return failureExplanations[group.Category]
	}
}

// Describe explains in plain language what the run did.
func Describe(meta RunMeta) string {
	if meta.IsConnTest() {
		role := "This is a control target (neutral, not part of NGR) to check that the network path itself is healthy."
		if !meta.Control {
			role = "This target is part of the system under investigation (NGR)."
		}

		return fmt.Sprintf(
			"Opens %d new connections to %s, %d at a time. Every request uses a brand new "+
				"connection (DNS lookup, TCP connect and TLS handshake) and the connection is closed afterwards, "+
				"so it measures how reliably and how fast a server accepts new connections. %s",
			meta.Config.Requests, meta.Target, max(meta.Config.Workers, 1), role,
		)
	}

	connections := "Connections are reused (keep-alive), like a normal client: only a few new connections are needed."
	if !meta.Config.KeepAlive {
		connections = "Every request uses a new connection (no keep-alive): a worst case for connection setup."
	}

	return fmt.Sprintf(
		"Mimics the INSPIRE validator towards the NGR CSW. For each of %d PDOK services it requests the "+
			"service record and, for the datasets it operates on, the tags and the dataset record. "+
			"%d services are processed in parallel, %d iteration(s). %s",
		meta.ServiceCount,
		max(meta.Config.Workers, 1),
		meta.Config.Iterations,
		connections,
	)
}

// BuildVerdict summarises a run in a title and a few plain language points.
func BuildVerdict(sum Summary) Verdict {
	all := sum.Overall
	if all.Count == 0 {
		return Verdict{Level: levelBad, Status: "No data", Title: "No requests were completed."}
	}

	verdict := Verdict{Level: levelForRate(all.ErrorRate)}
	verdict.Status = statusForLevel(verdict.Level)

	if all.Errors == 0 {
		verdict.Title = fmt.Sprintf("All %d requests succeeded.", all.Count)
	} else {
		verdict.Title = fmt.Sprintf("%d of %d requests failed (%.1f%%).",
			all.Errors, all.Count, all.ErrorRate*percent)
		verdict.Points = append(verdict.Points, failurePoints(sum)...)
	}

	setup := all.AvgDNSMs + all.AvgConnMs + all.AvgTLSMs
	if setup > 0 {
		verdict.Points = append(verdict.Points, fmt.Sprintf(
			"Setting up a connection took %s on average (DNS %s, TCP %s, TLS %s) of %s per successful request.",
			formatMs(setup),
			formatMs(all.AvgDNSMs),
			formatMs(all.AvgConnMs),
			formatMs(all.AvgTLSMs),
			formatMs(all.AvgMs),
		))
	}

	if all.AvgMs > 0 && all.P99Ms > slowTailFold*all.AvgMs {
		verdict.Points = append(verdict.Points, fmt.Sprintf(
			"Most requests are fast (average %s) but the slowest 1%% take %s or more (p99): "+
				"occasional very slow requests.", formatMs(all.AvgMs), formatMs(all.P99Ms)))
	}

	return verdict
}

func failurePoints(sum Summary) []string {
	var points []string

	top := sum.Errors[0]
	for _, group := range sum.Errors {
		if group.Count > top.Count {
			top = group
		}
	}

	points = append(points, fmt.Sprintf("Most common failure (%d times) on %s: %s%s.",
		top.Count, top.Host, top.Category, phaseSuffix(top.Phase)))

	if text := explanationFor(top); text != "" {
		points = append(points, text)
	}

	if len(sum.ByHost) > 1 {
		parts := make([]string, 0, len(sum.ByHost))
		for _, host := range sum.ByHost {
			parts = append(parts, fmt.Sprintf("%s: %d of %d (%.1f%%)",
				host.Name, host.Errors, host.Count, host.ErrorRate*percent))
		}

		points = append(points, "Failures per host - "+strings.Join(parts, "; ")+".")
	}

	return points
}

func phaseSuffix(phase string) string {
	if phase == "" {
		return ""
	}

	return " (during " + phase + ")"
}

// Finding is one conclusion on the overview page, drawn from several runs.
type Finding struct {
	Level string
	Title string
	Text  string
}

type tally struct{ errors, count, runs int }

func (t *tally) String() string {
	rate := 0.0
	if t.count > 0 {
		rate = float64(t.errors) / float64(t.count) * percent
	}

	return fmt.Sprintf(
		"%d of %d requests failed (%.2f%%) in %d run(s)",
		t.errors,
		t.count,
		rate,
		t.runs,
	)
}

func (t *tally) add(stats LatencyStats) {
	t.errors += stats.Errors
	t.count += stats.Count
	t.runs++
}

func buildFindings(runs []RunFile) []Finding {
	return append(keepAliveFinding(runs), connectionFindings(runs)...)
}

func keepAliveFinding(runs []RunFile) []Finding {
	var withKA, withoutKA tally

	for _, run := range runs {
		switch {
		case run.Meta.IsConnTest():
		case run.Meta.Config.KeepAlive:
			withKA.add(run.Summary.Overall)
		default:
			withoutKA.add(run.Summary.Overall)
		}
	}

	if withKA.runs == 0 || withoutKA.runs == 0 {
		return nil
	}

	text := "With keep-alive: " + withKA.String() + ". Without keep-alive: " + withoutKA.String() + "."

	switch {
	case withKA.errors == 0 && withoutKA.errors > 0:
		return []Finding{{
			Level: levelWarn,
			Title: "Failures only appear when every request needs a new connection",
			Text: text + " The CSW answers reliably as long as connections are reused. " +
				"Failures show up when every request has to set up a new connection, " +
				"which points at connection setup in the network edge rather than at the CSW application.",
		}}
	case withKA.errors > 0:
		return []Finding{{
			Level: levelBad,
			Title: "Failures also occur when connections are reused",
			Text:  text + " So not only connection setup is a problem.",
		}}
	default:
		return []Finding{{Level: levelGood, Title: "No failures in any load test", Text: text}}
	}
}

func connectionFindings(runs []RunFile) []Finding {
	type split struct{ subject, control tally }

	byOrigin := map[string]*split{}

	var origins []string

	for _, run := range runs {
		if !run.Meta.IsConnTest() {
			continue
		}

		entry, ok := byOrigin[run.Meta.Origin]
		if !ok {
			entry = &split{}
			byOrigin[run.Meta.Origin] = entry

			origins = append(origins, run.Meta.Origin)
		}

		if run.Meta.Control {
			entry.control.add(run.Summary.Overall)
		} else {
			entry.subject.add(run.Summary.Overall)
		}
	}

	findings := make([]Finding, 0, len(origins))

	for _, origin := range origins {
		entry := byOrigin[origin]
		text := fmt.Sprintf(
			"NGR: %s. Control targets: %s.",
			entry.subject.String(),
			entry.control.String(),
		)

		switch {
		case entry.control.runs == 0 || entry.subject.runs == 0:
			findings = append(findings, Finding{
				Level: levelForRate(rateOf(entry.subject)),
				Title: "New connections from " + origin,
				Text:  text + " Run both NGR and control suites to be able to compare.",
			})
		case entry.control.errors == 0 && entry.subject.errors > 0:
			level, title := levelBad, fmt.Sprintf("From %s the problem is specific to NGR", origin)
			if rateOf(entry.subject) < minorRate {
				level = levelWarn
				title = fmt.Sprintf("From %s only NGR shows a few failed new connections", origin)
			}

			findings = append(findings, Finding{
				Level: level,
				Title: title,
				Text: text + " Neutral targets accept every new connection, so the network path " +
					"from here looks fine and the connection failures are caused by NGR or its network edge.",
			})
		case entry.control.errors > 0:
			findings = append(findings, Finding{
				Level: levelWarn,
				Title: fmt.Sprintf(
					"From %s new connections also fail towards neutral targets",
					origin,
				),
				Text: text + " The problem is at least partly in the path from here " +
					"(egress, NAT, proxy or firewall), not only in NGR.",
			})
		default:
			findings = append(findings, Finding{
				Level: levelGood,
				Title: fmt.Sprintf("From %s new connections are healthy", origin),
				Text:  text + " At this load NGR accepts new connections as reliably as the control targets.",
			})
		}
	}

	return findings
}

func rateOf(t tally) float64 {
	if t.count == 0 {
		return 0
	}

	return float64(t.errors) / float64(t.count)
}

// matrixCell is the latest connection test result of one suite from one origin.
type matrixCell struct {
	Has       bool
	Index     int
	Errors    int
	Count     int
	ErrorRate float64
	SetupMs   float64
	P99Ms     float64
}

type matrixRow struct {
	Suite   string
	Target  string
	Control bool
	Cells   []matrixCell
}

type matrix struct {
	Origins []string
	Rows    []matrixRow
}

func buildMatrix(runs []RunFile) matrix {
	var result matrix

	rowIndex := map[string]int{}
	colIndex := map[string]int{}

	for _, run := range runs {
		if run.Meta.IsConnTest() {
			if _, ok := colIndex[run.Meta.Origin]; !ok {
				colIndex[run.Meta.Origin] = len(result.Origins)
				result.Origins = append(result.Origins, run.Meta.Origin)
			}

			if _, ok := rowIndex[run.Meta.Suite]; !ok {
				rowIndex[run.Meta.Suite] = len(result.Rows)
				result.Rows = append(result.Rows, matrixRow{
					Suite: run.Meta.Suite, Target: run.Meta.Target, Control: run.Meta.Control,
				})
			}
		}
	}

	for i := range result.Rows {
		result.Rows[i].Cells = make([]matrixCell, len(result.Origins))
	}

	for i, run := range runs {
		if !run.Meta.IsConnTest() {
			continue
		}

		all := run.Summary.Overall
		result.Rows[rowIndex[run.Meta.Suite]].Cells[colIndex[run.Meta.Origin]] = matrixCell{
			Has: true, Index: i + 1, Errors: all.Errors, Count: all.Count, ErrorRate: all.ErrorRate,
			SetupMs: all.AvgDNSMs + all.AvgConnMs + all.AvgTLSMs, P99Ms: all.P99Ms,
		}
	}

	return result
}
