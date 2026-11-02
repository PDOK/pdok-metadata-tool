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

		raw, err := os.ReadFile(
			filepath.Join(dir, entry.Name(), "summary.json"),
		) //nolint:gosec // user supplied dir
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
	Legend      []chartSeries
	LatencyBars []overviewChart
	ErrorBars   []overviewChart
	Generated   time.Time
	TotalReqs   int
	TotalFailed int
}

// RenderAggregate renders one self-contained HTML with an overview of all runs and a
// detail page per run.
func RenderAggregate(runs []RunFile) (string, error) {
	var hosts []string

	seen := map[string]bool{}

	for _, run := range runs {
		for _, host := range run.Summary.Hosts {
			if !seen[host] {
				seen[host] = true

				hosts = append(hosts, host)
			}
		}
	}

	sort.Strings(hosts)

	colors := hostColorMap(hosts)

	tmpl, err := template.New("report").Funcs(reportFuncs(colors)).Parse(reportTemplate)
	if err != nil {
		return "", err
	}

	if _, err := tmpl.New("aggregate").Parse(aggregateTemplate); err != nil {
		return "", err
	}

	data := aggregateData{Generated: time.Now()}
	for _, host := range hosts {
		data.Legend = append(data.Legend, chartSeries{Color: colors[host], Name: host})
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

		data.Rows = append(data.Rows, overviewRow{
			Index:   i + 1,
			Dir:     run.Dir,
			Meta:    run.Meta,
			Overall: run.Summary.Overall,
			Hosts:   run.Summary.ByHost,
			Detail:  template.HTML(detail.String()), //nolint:gosec // rendered by our own template
		})
		data.TotalReqs += run.Summary.Overall.Count
		data.TotalFailed += run.Summary.Overall.Errors
	}

	data.LatencyBars = overviewCharts(
		data.Rows,
		hosts,
		colors,
		func(s LatencyStats) float64 { return s.AvgMs },
		formatMs,
	)
	data.ErrorBars = overviewCharts(
		data.Rows,
		hosts,
		colors,
		func(s LatencyStats) float64 { return s.ErrorRate },
		func(v float64) string { return fmt.Sprintf("%.2f%%", v*100) },
	) //nolint:mnd // percentage

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "aggregate", data); err != nil {
		return "", err
	}

	return buf.String(), nil
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
