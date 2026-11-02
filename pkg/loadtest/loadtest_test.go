package loadtest_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pdok/pdok-metadata-tool/v2/pkg/loadtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const serviceRecordXML = `<?xml version="1.0"?>
<csw:GetRecordByIdResponse xmlns:csw="http://www.opengis.net/cat/csw/2.0.2" xmlns:gmd="http://www.isotc211.org/2005/gmd" xmlns:srv="http://www.isotc211.org/2005/srv" xmlns:xlink="http://www.w3.org/1999/xlink">
<gmd:MD_Metadata><gmd:identificationInfo><srv:SV_ServiceIdentification>
<srv:serviceType><gco:LocalName xmlns:gco="http://www.isotc211.org/2005/gco">view</gco:LocalName></srv:serviceType>
<srv:operatesOn xlink:href="http://x/csw?id=aaaa-1#MD_DataIdentification"/>
<srv:operatesOn xlink:href="http://x/csw?id=bbbb-2#MD_DataIdentification"/>
</srv:SV_ServiceIdentification></gmd:identificationInfo></gmd:MD_Metadata>
</csw:GetRecordByIdResponse>`

func fakeNGR(t *testing.T, scrapes *atomic.Int32) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		switch {
		case q.Get("request") == "GetRecords":
			scrapes.Add(1)

			if q.Get("startPosition") == "1" {
				fmt.Fprint(
					w,
					`<csw:GetRecordsResponse xmlns:csw="http://www.opengis.net/cat/csw/2.0.2" xmlns:dc="http://purl.org/dc/elements/1.1/"><csw:SearchResults numberOfRecordsMatched="2" nextRecord="2"><csw:SummaryRecord><dc:identifier>svc-1</dc:identifier><dc:title>One</dc:title></csw:SummaryRecord></csw:SearchResults></csw:GetRecordsResponse>`,
				)

				return
			}

			fmt.Fprint(
				w,
				`<csw:GetRecordsResponse xmlns:csw="http://www.opengis.net/cat/csw/2.0.2" xmlns:dc="http://purl.org/dc/elements/1.1/"><csw:SearchResults numberOfRecordsMatched="2" nextRecord="0"><csw:SummaryRecord><dc:identifier>svc-2</dc:identifier><dc:title>Two</dc:title></csw:SummaryRecord></csw:SearchResults></csw:GetRecordsResponse>`,
			)
		case q.Get("request") == "GetRecordById" && q.Get("id") == "svc-2":
			w.WriteHeader(http.StatusServiceUnavailable)
		case q.Get("request") == "GetRecordById" && strings.HasPrefix(q.Get("id"), "svc-"):
			fmt.Fprint(w, serviceRecordXML)
		case q.Get("request") == "GetRecordById":
			fmt.Fprint(
				w,
				"<csw:GetRecordByIdResponse xmlns:csw=\"http://www.opengis.net/cat/csw/2.0.2\"/>",
			)
		case strings.HasSuffix(r.URL.Path, "/tags") && strings.Contains(r.URL.Path, "aaaa-1"):
			fmt.Fprint(w, `[{"id":1,"name":"inspire"}]`)
		case strings.HasSuffix(r.URL.Path, "/tags"):
			fmt.Fprint(w, `[{"id":2,"name":"other"}]`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestLoadServiceRecordsScrapesPagesAndUsesCache(t *testing.T) {
	t.Parallel()

	var scrapes atomic.Int32

	srv := fakeNGR(t, &scrapes)
	defer srv.Close()

	opts := loadtest.IDCacheOptions{
		Path: filepath.Join(t.TempDir(), "ids.json"),
		TTL:  time.Hour,
	}

	records, cached, err := loadtest.LoadServiceRecords(context.Background(), srv.URL, opts)
	require.NoError(t, err)
	assert.False(t, cached)
	assert.Equal(
		t,
		[]loadtest.ServiceRecord{{ID: "svc-1", Title: "One"}, {ID: "svc-2", Title: "Two"}},
		records,
	)
	assert.EqualValues(t, 2, scrapes.Load())

	_, cached, err = loadtest.LoadServiceRecords(context.Background(), srv.URL, opts)
	require.NoError(t, err)
	assert.True(t, cached)
	assert.EqualValues(t, 2, scrapes.Load(), "cache hit must not scrape again")

	opts.Refresh = true
	_, cached, err = loadtest.LoadServiceRecords(context.Background(), srv.URL, opts)
	require.NoError(t, err)
	assert.False(t, cached)
	assert.EqualValues(t, 4, scrapes.Load())
}

func TestRunMimicsValidatorAndClassifiesErrors(t *testing.T) {
	t.Parallel()

	var scrapes atomic.Int32

	srv := fakeNGR(t, &scrapes)
	defer srv.Close()

	var results []loadtest.Result

	info, err := loadtest.Run(context.Background(), loadtest.Config{
		Hosts:           []string{srv.URL},
		Workers:         2,
		Timeout:         5 * time.Second,
		Iterations:      2,
		IncludeDatasets: true,
		KeepAlive:       true,
	}, []loadtest.ServiceRecord{{ID: "svc-1"}, {ID: "svc-2"}}, func(r loadtest.Result) {
		results = append(results, r)
	})
	require.NoError(t, err)

	kinds := map[string]int{}
	failed := 0

	for _, r := range results {
		kinds[r.Kind]++

		if !r.OK {
			failed++

			assert.Equal(t, loadtest.ErrHTTPStatus, r.ErrorCategory)
			assert.Equal(t, http.StatusServiceUnavailable, r.Status)
		}
	}

	// svc-1: service record + 2 tag requests + 1 record (only aaaa-1 is INSPIRE tagged), twice.
	assert.Equal(t, 4, kinds[loadtest.KindServiceRecord], "two iterations of both services")
	assert.Equal(t, 4, kinds[loadtest.KindDatasetTags])
	assert.Equal(t, 2, kinds[loadtest.KindDatasetRecord])
	assert.Equal(t, 2, failed, "svc-2 fails in both iterations")

	sum := loadtest.Summarize(results, info, []string{srv.URL})
	assert.Equal(t, 2, sum.FailedTotal)
	assert.Equal(t, 2, sum.ServicesSeen)
	require.NotEmpty(t, sum.Errors)
	assert.Equal(t, 2, sum.Errors[0].Count)
}

func TestRunTimeoutIsClassified(t *testing.T) {
	t.Parallel()

	slow := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(2 * time.Second):
		}
	}))
	defer slow.Close()

	var results []loadtest.Result

	_, err := loadtest.Run(context.Background(), loadtest.Config{
		Hosts:      []string{slow.URL},
		Workers:    1,
		Timeout:    100 * time.Millisecond,
		Iterations: 1,
	}, []loadtest.ServiceRecord{{ID: "svc-1"}}, func(r loadtest.Result) { results = append(results, r) })
	require.NoError(t, err)

	require.Len(t, results, 1)
	assert.False(t, results[0].OK)
	assert.Equal(t, loadtest.ErrTimeout, results[0].ErrorCategory)
	assert.Equal(t, "wait response", results[0].ErrorPhase)
}

func TestRunRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	svcs := []loadtest.ServiceRecord{{ID: "x"}}
	sink := func(loadtest.Result) {}

	_, err := loadtest.Run(context.Background(), loadtest.Config{Iterations: 1}, svcs, sink)
	require.Error(t, err)

	_, err = loadtest.Run(
		context.Background(),
		loadtest.Config{Hosts: []string{"http://x"}},
		svcs,
		sink,
	)
	require.Error(t, err)
}

func TestExecuteWritesReport(t *testing.T) {
	t.Parallel()

	var scrapes atomic.Int32

	srv := fakeNGR(t, &scrapes)
	defer srv.Close()

	out, err := loadtest.Execute(context.Background(), loadtest.Options{
		Config: loadtest.Config{
			Hosts:           []string{srv.URL},
			Workers:         2,
			Timeout:         5 * time.Second,
			Iterations:      1,
			IncludeDatasets: true,
			KeepAlive:       true,
		},
		IDCache: loadtest.IDCacheOptions{
			Path: filepath.Join(t.TempDir(), "ids.json"),
			TTL:  time.Hour,
		},
		Label:     "unit test/lab",
		OutputDir: t.TempDir(),
		WANIP:     "203.0.113.7",
	})
	require.NoError(t, err)

	for _, name := range []string{"results.jsonl", "summary.json", "report.html", "failures.csv"} {
		assert.FileExists(t, filepath.Join(out.Dir, name))
	}

	html, err := os.ReadFile(filepath.Join(out.Dir, "report.html"))
	require.NoError(t, err)
	assert.Contains(t, string(html), "203.0.113.7")
	assert.Contains(t, string(html), "HTTP 503 Service Unavailable")
	assert.Contains(t, string(html), "<svg")
	assert.Contains(t, loadtest.TextSummary(out), "TOTAL")
}

func TestAggregateRendersOverviewAndDetails(t *testing.T) {
	t.Parallel()

	var scrapes atomic.Int32

	srv := fakeNGR(t, &scrapes)
	defer srv.Close()

	outDir := t.TempDir()

	for _, label := range []string{"first", "second"} {
		_, err := loadtest.Execute(context.Background(), loadtest.Options{
			Config: loadtest.Config{
				Hosts: []string{srv.URL}, Workers: 1, Timeout: 5 * time.Second, Iterations: 1,
				IncludeDatasets: true, KeepAlive: true,
			},
			IDCache: loadtest.IDCacheOptions{
				Path: filepath.Join(t.TempDir(), "ids.json"),
				TTL:  time.Hour,
			},
			Label:     label,
			OutputDir: outDir,
			WANIP:     "203.0.113.7",
		})
		require.NoError(t, err)
	}

	runs, err := loadtest.LoadRuns(outDir)
	require.NoError(t, err)
	require.Len(t, runs, 2)

	html, err := loadtest.RenderAggregate(runs)
	require.NoError(t, err)

	for _, want := range []string{`id="overview"`, `id="run-1"`, `id="run-2"`, "first", "second", "HTTP 503"} {
		assert.Contains(t, html, want)
	}
}

func TestConnectionTestsAndAggregate(t *testing.T) {
	t.Parallel()

	good := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	}))
	defer good.Close()

	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer bad.Close()

	outDir := t.TempDir()

	outcomes, err := loadtest.ExecuteConnectionTests(context.Background(), loadtest.ConnOptions{
		Suites: []loadtest.ConnSuite{
			{Name: "ngr-fake", URL: bad.URL + "/csw"},
			{Name: "control-fake", URL: good.URL + "/", Control: true},
		},
		Requests:  6,
		Workers:   2,
		Timeout:   5 * time.Second,
		Label:     "unit",
		OutputDir: outDir,
		WANIP:     "203.0.113.7",
	})
	require.NoError(t, err)
	require.Len(t, outcomes, 2)
	assert.Equal(t, 6, outcomes[0].Summary.Overall.Errors)
	assert.Equal(t, 0, outcomes[1].Summary.Overall.Errors)

	runs, err := loadtest.LoadRuns(outDir)
	require.NoError(t, err)
	require.Len(t, runs, 2)

	html, err := loadtest.RenderAggregate(runs)
	require.NoError(t, err)

	for _, want := range []string{
		"Connection tests: new connections per target", "From unit", "ngr-fake", "control-fake",
		"the problem is specific to NGR", "6 of 6 requests failed",
	} {
		assert.Contains(t, html, want)
	}
}

func TestParseConnSuite(t *testing.T) {
	t.Parallel()

	suite, err := loadtest.ParseConnSuite("control-x=https://example.org/a")
	require.NoError(t, err)
	assert.True(t, suite.Control)
	assert.Equal(t, "https://example.org/a", suite.URL)

	_, err = loadtest.ParseConnSuite("not a url")
	require.Error(t, err)
}
