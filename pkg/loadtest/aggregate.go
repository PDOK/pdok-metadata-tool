package loadtest

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"time"
)

//go:embed aggregate.html.tmpl
var aggregateTemplate string

// RunFile is a finished run as stored in a run directory (summary.json).
type RunFile struct {
	Dir     string  `json:"-"`
	Meta    RunMeta `json:"meta"`
	Summary Summary `json:"summary"`
}

// LoadRuns reads the summary.json of every run directory in dir, oldest first.
// Directories without a summary.json are skipped.
func LoadRuns(dir string) ([]RunFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var runs []RunFile

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		summaryPath := filepath.Join(dir, entry.Name(), "summary.json")

		raw, err := os.ReadFile(summaryPath) //nolint:gosec // user supplied dir
		if err != nil {
			continue
		}

		run := RunFile{Dir: entry.Name()}
		if err := json.Unmarshal(raw, &run); err != nil {
			return nil, fmt.Errorf("reading %s: %w", entry.Name(), err)
		}

		runs = append(runs, run)
	}

	sort.SliceStable(
		runs,
		func(i, j int) bool { return runs[i].Meta.Started.Before(runs[j].Meta.Started) },
	)

	return runs, nil
}

type overviewRow struct {
	Index   int
	Dir     string
	Meta    RunMeta
	Overall LatencyStats
	Hosts   []GroupStats
	Verdict Verdict
	Detail  template.HTML
}

type overviewBar struct {
	Host   string
	Color  string
	Value  float64
	Pct    float64
	Format string
}

type overviewChart struct {
	Label string
	Bars  []overviewBar
}

type aggregateData struct {
	Rows        []overviewRow
	LoadRows    []overviewRow
	ConnRows    []overviewRow
	Findings    []Finding
	Matrix      matrix
	Legend      []chartSeries
	ConnLegend  []chartSeries
	LatencyBars []overviewChart
	ErrorBars   []overviewChart
	ConnLatency []overviewChart
	ConnErrors  []overviewChart
	Generated   time.Time
	TotalReqs   int
	TotalFailed int
}

// RenderAggregate renders one self-contained HTML with an overview of all runs and a
// detail page per run.
func RenderAggregate(runs []RunFile) (string, error) {
	loadHosts, connHosts := hostsByType(runs)
	colors := hostColorMap(uniqueHosts(loadHosts, connHosts))

	tmpl, err := template.New("report").Funcs(reportFuncs(colors)).Parse(reportTemplate)
	if err != nil {
		return "", err
	}

	if _, err := tmpl.New("aggregate").Parse(aggregateTemplate); err != nil {
		return "", err
	}

	data := aggregateData{
		Generated:  time.Now(),
		Findings:   buildFindings(runs),
		Matrix:     buildMatrix(runs),
		Legend:     legendFor(loadHosts, colors),
		ConnLegend: legendFor(connHosts, colors),
	}

	for i, run := range runs {
		var detail bytes.Buffer
		if err := tmpl.ExecuteTemplate(
			&detail,
			"runbody",
			buildReportData(run.Meta, run.Summary, colors),
		); err != nil {
			return "", err
		}

		row := overviewRow{
			Index:   i + 1,
			Dir:     run.Dir,
			Meta:    run.Meta,
			Overall: run.Summary.Overall,
			Hosts:   run.Summary.ByHost,
			Verdict: BuildVerdict(run.Summary),
			Detail:  template.HTML(detail.String()), //nolint:gosec // rendered by our own template
		}

		data.Rows = append(data.Rows, row)

		if run.Meta.IsConnTest() {
			data.ConnRows = append(data.ConnRows, row)
		} else {
			data.LoadRows = append(data.LoadRows, row)
		}

		data.TotalReqs += run.Summary.Overall.Count
		data.TotalFailed += run.Summary.Overall.Errors
	}

	avg := func(s LatencyStats) float64 { return s.AvgMs }
	errRate := func(s LatencyStats) float64 { return s.ErrorRate }
	pct := func(v float64) string { return fmt.Sprintf("%.2f%%", v*percent) }

	data.LatencyBars = overviewCharts(data.LoadRows, loadHosts, colors, avg, formatMs)
	data.ErrorBars = overviewCharts(data.LoadRows, loadHosts, colors, errRate, pct)
	data.ConnLatency = overviewCharts(data.ConnRows, connHosts, colors, avg, formatMs)
	data.ConnErrors = overviewCharts(data.ConnRows, connHosts, colors, errRate, pct)

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "aggregate", data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func hostsByType(runs []RunFile) ([]string, []string) {
	var load, conn []string

	seenLoad, seenConn := map[string]bool{}, map[string]bool{}

	for _, run := range runs {
		target, seen := &load, seenLoad
		if run.Meta.IsConnTest() {
			target, seen = &conn, seenConn
		}

		for _, host := range run.Summary.Hosts {
			if !seen[host] {
				seen[host] = true

				*target = append(*target, host)
			}
		}
	}

	sort.Strings(load)
	sort.Strings(conn)

	return load, conn
}

func legendFor(hosts []string, colors map[string]string) []chartSeries {
	legend := make([]chartSeries, 0, len(hosts))
	for _, host := range hosts {
		legend = append(legend, chartSeries{Color: colors[host], Name: host})
	}

	return legend
}

func uniqueHosts(lists ...[]string) []string {
	var out []string

	for _, list := range lists {
		for _, host := range list {
			if !contains(out, host) {
				out = append(out, host)
			}
		}
	}

	return out
}
func overviewCharts(
	rows []overviewRow,
	hosts []string,
	colors map[string]string,
	value func(LatencyStats) float64,
	format func(float64) string,
) []overviewChart {
	maxVal := 0.0

	for _, row := range rows {
		for _, h := range row.Hosts {
			maxVal = max(maxVal, value(h.LatencyStats))
		}
	}

	charts := make([]overviewChart, 0, len(rows))

	for _, row := range rows {
		chart := overviewChart{Label: fmt.Sprintf("%d. %s", row.Index, row.Meta.Label)}

		for _, host := range hosts {
			for _, h := range row.Hosts {
				if h.Name != host {
					continue
				}

				v := value(h.LatencyStats)
				pct := 0.0

				if maxVal > 0 {
					pct = v / maxVal * 100 //nolint:mnd // percentage
				}

				chart.Bars = append(chart.Bars, overviewBar{
					Host: host, Color: colors[host], Value: v, Pct: pct, Format: format(v),
				})
			}
		}

		charts = append(charts, chart)
	}

	return charts
}
