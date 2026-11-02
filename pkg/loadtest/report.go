package loadtest

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	unknownValue   = "unknown"
	wanIPTimeout   = 5 * time.Second
	wanIPMaxBytes  = 128
	chartWidth     = 960.0
	chartHeight    = 240.0
	chartPadLeft   = 56.0
	chartPadRight  = 12.0
	chartPadTop    = 12.0
	chartPadBottom = 28.0
)

//go:embed report.html.tmpl
var reportTemplate string

var hostColors = []string{"#4f9dff", "#ff8a4c", "#43d19e", "#c792ea", "#ffd166"}

// RunMeta describes the environment and configuration of a run, shown in the report.
type RunMeta struct {
	Label        string        `json:"label"`
	Hostname     string        `json:"hostname"`
	WANIP        string        `json:"wanIp"`
	Config       Config        `json:"config"`
	Started      time.Time     `json:"started"`
	Finished     time.Time     `json:"finished"`
	Aborted      bool          `json:"aborted"`
	ServiceCount int           `json:"serviceCount"`
	IDsFromCache bool          `json:"idsFromCache"`
	IDsFetchTime time.Duration `json:"idsFetchTimeNs"`
}

// LookupWANIP asks a public "what is my IP" service for the WAN address of this machine.
func LookupWANIP(ctx context.Context, lookupURL string) string {
	ctx, cancel := context.WithTimeout(ctx, wanIPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, lookupURL, nil)
	if err != nil {
		return unknownValue
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return unknownValue
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, wanIPMaxBytes))
	if err != nil || resp.StatusCode != http.StatusOK {
		return unknownValue
	}

	return strings.TrimSpace(string(body))
}

// Hostname returns the local hostname, or "unknown".
func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return unknownValue
	}

	return name
}

type chartSeries struct {
	Color string
	Name  string
}

type reportData struct {
	Meta       RunMeta
	Summary    Summary
	Duration   string
	Colors     map[string]string
	Legend     []chartSeries
	LatencySVG template.HTML
	CountSVG   template.HTML
	BucketSize string
	Hist       []histView
	FailedNote string
}

type histView struct {
	Label string
	Bars  []histBar
}

type histBar struct {
	Host  string
	Count int
	Pct   float64
	Color string
}

// RenderReport renders the self-contained HTML report.
func RenderReport(meta RunMeta, sum Summary) (string, error) {
	colors := hostColorMap(sum.Hosts)

	tmpl, err := template.New("report").Funcs(reportFuncs(colors)).Parse(reportTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, buildReportData(meta, sum, colors)); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func hostColorMap(hosts []string) map[string]string {
	colors := make(map[string]string, len(hosts))
	for i, host := range hosts {
		colors[host] = hostColors[i%len(hostColors)]
	}

	return colors
}

func buildReportData(meta RunMeta, sum Summary, colors map[string]string) reportData {
	legend := make([]chartSeries, 0, len(sum.Hosts))
	for _, host := range sum.Hosts {
		legend = append(legend, chartSeries{Color: colors[host], Name: host})
	}

	data := reportData{
		Meta:       meta,
		Summary:    sum,
		Duration:   meta.Finished.Sub(meta.Started).Round(time.Second).String(),
		Colors:     colors,
		Legend:     legend,
		BucketSize: sum.BucketSize.String(),
		LatencySVG: lineChart(sum, colors, func(b Bucket) float64 { return b.AvgMs }, "ms"),
		CountSVG:   barChart(sum, colors),
		Hist:       histViews(sum, colors),
	}

	if sum.FailedTotal > len(sum.Failed) {
		data.FailedNote = fmt.Sprintf(
			"Showing the first %d of %d failed requests; all of them are in results.jsonl.",
			len(sum.Failed), sum.FailedTotal,
		)
	}

	return data
}

func reportFuncs(colors map[string]string) template.FuncMap {
	return template.FuncMap{
		"ms":       formatMs,
		"pct":      func(f float64) string { return fmt.Sprintf("%.2f%%", f*100) },
		"f1":       func(f float64) string { return fmt.Sprintf("%.1f", f) },
		"f2":       func(f float64) string { return fmt.Sprintf("%.2f", f) },
		"kb":       func(f float64) string { return fmt.Sprintf("%.1f KB", f/1024) }, //nolint:mnd // bytes per KB
		"ts":       func(t time.Time) string { return t.Format("15:04:05.000") },
		"datetime": func(t time.Time) string { return t.Format("2006-01-02 15:04:05 MST") },
		"color": func(host string) string {
			if c, ok := colors[host]; ok {
				return c
			}

			return "#888"
		},
		"rateClass": func(f float64) string {
			switch {
			case f == 0:
				return "good"
			case f < 0.05: //nolint:mnd // 5% threshold
				return "warn"
			default:
				return "bad"
			}
		},
	}
}

func formatMs(v float64) string {
	switch {
	case v >= 1000: //nolint:mnd // switch to seconds
		return fmt.Sprintf("%.2f s", v/1000) //nolint:mnd // ms to s
	case v >= 100: //nolint:mnd // no decimals above 100 ms
		return fmt.Sprintf("%.0f ms", v)
	default:
		return fmt.Sprintf("%.1f ms", v)
	}
}

func histViews(sum Summary, colors map[string]string) []histView {
	maxCount := 1

	for _, row := range sum.Histogram {
		for _, c := range row.Counts {
			maxCount = max(maxCount, c)
		}
	}

	views := make([]histView, 0, len(sum.Histogram))

	for _, row := range sum.Histogram {
		view := histView{Label: row.Label}
		for _, host := range sum.Hosts {
			view.Bars = append(view.Bars, histBar{
				Host:  host,
				Count: row.Counts[host],
				Pct: float64(
					row.Counts[host],
				) / float64(
					maxCount,
				) * 100, //nolint:mnd // percentage
				Color: colors[host],
			})
		}

		views = append(views, view)
	}

	return views
}

type plot struct {
	sb           strings.Builder
	maxY         float64
	buckets      int
	innerW       float64
	innerH       float64
	yLabelSuffix string
	centered     bool
}

func newPlot(maxY float64, buckets int, suffix string) *plot {
	p := &plot{
		maxY:         niceMax(maxY),
		buckets:      max(buckets, 1),
		innerW:       chartWidth - chartPadLeft - chartPadRight,
		innerH:       chartHeight - chartPadTop - chartPadBottom,
		yLabelSuffix: suffix,
	}

	fmt.Fprintf(&p.sb, `<svg viewBox="0 0 %.0f %.0f" role="img">`, chartWidth, chartHeight)

	const gridLines = 4
	for i := 0; i <= gridLines; i++ {
		y := chartPadTop + p.innerH - p.innerH*float64(i)/gridLines
		fmt.Fprintf(
			&p.sb,
			`<line x1="%.1f" x2="%.1f" y1="%.1f" y2="%.1f" class="grid"/><text x="%.1f" y="%.1f" class="axis" text-anchor="end">%s</text>`,
			chartPadLeft,
			chartWidth-chartPadRight,
			y,
			y,
			chartPadLeft-6,
			y+4,
			axisLabel(p.maxY*float64(i)/gridLines, suffix),
		)
	}

	return p
}

func (p *plot) x(i int) float64 {
	if p.centered {
		return chartPadLeft + p.innerW*(float64(i)+0.5)/float64(
			p.buckets,
		) //nolint:mnd // bucket centre
	}

	if p.buckets <= 1 {
		return chartPadLeft + p.innerW/2 //nolint:mnd // centre
	}

	return chartPadLeft + p.innerW*float64(i)/float64(p.buckets-1)
}

func (p *plot) y(v float64) float64 {
	return chartPadTop + p.innerH - p.innerH*min(v/p.maxY, 1)
}

func (p *plot) xLabels(sum Summary) {
	step := max(p.buckets/8, 1) //nolint:mnd // about 8 labels

	for i := 0; i < p.buckets; i += step {
		offset := time.Duration(i) * sum.BucketSize

		fmt.Fprintf(&p.sb, `<text x="%.1f" y="%.1f" class="axis" text-anchor="middle">+%s</text>`,
			p.x(i), chartHeight-8, offset.Round(time.Second)) //nolint:mnd // label offset
	}
}

func (p *plot) close() template.HTML {
	p.sb.WriteString(`</svg>`)

	return template.HTML(p.sb.String()) //nolint:gosec // built from numbers and escaped host names
}

func niceMax(v float64) float64 {
	if v <= 0 {
		return 1
	}

	magnitude := 1.0
	for magnitude*10 <= v { //nolint:mnd // decimal magnitude
		magnitude *= 10
	}

	for _, m := range []float64{1, 2, 5, 10} { //nolint:mnd // nice steps
		if m*magnitude >= v {
			return m * magnitude
		}
	}

	return v
}

func axisLabel(v float64, suffix string) string {
	if suffix == "ms" && v >= 1000 { //nolint:mnd // switch to seconds
		return fmt.Sprintf("%.1f s", v/1000) //nolint:mnd // ms to s
	}

	if v == float64(int(v)) {
		return fmt.Sprintf("%d %s", int(v), suffix)
	}

	return fmt.Sprintf("%.1f %s", v, suffix)
}

func bucketsByHost(sum Summary) (map[string][]Bucket, int) {
	byHost := map[string][]Bucket{}
	count := 0

	for _, b := range sum.Timeline {
		byHost[b.Host] = append(byHost[b.Host], b)
		count = max(count, len(byHost[b.Host]))
	}

	return byHost, count
}

func lineChart(
	sum Summary,
	colors map[string]string,
	value func(Bucket) float64,
	suffix string,
) template.HTML {
	byHost, count := bucketsByHost(sum)

	maxY := 0.0

	for _, bs := range byHost {
		for _, b := range bs {
			maxY = max(maxY, value(b))
		}
	}

	p := newPlot(maxY, count, suffix)

	for _, host := range sum.Hosts {
		var points []string

		for i, b := range byHost[host] {
			if b.Count-b.Errors > 0 {
				points = append(points, fmt.Sprintf("%.1f,%.1f", p.x(i), p.y(value(b))))
			}
		}

		if len(points) == 0 {
			continue
		}

		fmt.Fprintf(&p.sb, `<polyline fill="none" stroke="%s" stroke-width="2" points="%s"/>`,
			colors[host], strings.Join(points, " "))

		for _, pt := range points {
			xy := strings.Split(pt, ",")
			fmt.Fprintf(
				&p.sb,
				`<circle cx="%s" cy="%s" r="2.5" fill="%s"/>`,
				xy[0],
				xy[1],
				colors[host],
			)
		}
	}

	p.xLabels(sum)

	return p.close()
}

// barChart shows the number of requests per bucket per host, with failures in red on top.
func barChart(sum Summary, colors map[string]string) template.HTML {
	byHost, count := bucketsByHost(sum)

	maxY := 0.0

	for _, bs := range byHost {
		for _, b := range bs {
			maxY = max(maxY, float64(b.Count))
		}
	}

	p := newPlot(maxY, count, "")
	p.centered = true
	groupW := p.innerW / float64(p.buckets)
	barW := groupW * 0.8 / float64(max(len(sum.Hosts), 1)) //nolint:mnd // spacing between groups

	for hi, host := range sum.Hosts {
		for i, b := range byHost[host] {
			x := chartPadLeft + groupW*float64(
				i,
			) + groupW*0.1 + barW*float64(
				hi,
			) //nolint:mnd // spacing
			okH := p.innerH * float64(b.Count-b.Errors) / p.maxY
			errH := p.innerH * float64(b.Errors) / p.maxY
			base := chartPadTop + p.innerH

			fmt.Fprintf(&p.sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`,
				x, base-okH, barW, okH, colors[host])

			if b.Errors > 0 {
				fmt.Fprintf(
					&p.sb,
					`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" class="errbar"/>`,
					x,
					base-okH-errH,
					barW,
					errH,
				)
			}
		}
	}

	p.xLabels(sum)

	return p.close()
}
