package loadtest

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// KindNewConnection is the request kind of a connection test.
const KindNewConnection = "New connection (GET)"

// ConnSuite is one target of the connection test. Control suites are neutral targets that are
// not part of the system under investigation; they show whether the network path itself is sound.
type ConnSuite struct {
	Name    string
	URL     string
	Control bool
}

// DefaultConnSuites are the NGR (with and without www) and a number of neutral control targets.
func DefaultConnSuites() []ConnSuite {
	return []ConnSuite{
		{Name: "ngr-www", URL: "https://www.nationaalgeoregister.nl" + CswInspirePath +
			"?service=CSW&request=GetCapabilities"},
		{Name: "ngr-plain", URL: "https://nationaalgeoregister.nl" + CswInspirePath +
			"?service=CSW&request=GetCapabilities"},
		{Name: "control-kadaster", URL: "https://www.kadaster.nl/", Control: true},
		{Name: "control-rijksoverheid", URL: "https://www.rijksoverheid.nl/", Control: true},
		{Name: "control-example", URL: "https://example.com/", Control: true},
	}
}

// ParseConnSuite parses "name=url" or just "url". Suites named "control-*" are control suites.
func ParseConnSuite(spec string) (ConnSuite, error) {
	name, rawURL, found := strings.Cut(spec, "=")
	if !found {
		rawURL = name
		name = ""
	}

	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ConnSuite{}, fmt.Errorf("invalid suite %q, expected name=https://host/path", spec)
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = parsed.Host
	}

	return ConnSuite{
		Name:    name,
		URL:     strings.TrimSpace(rawURL),
		Control: strings.HasPrefix(name, "control-"),
	}, nil
}

// ConnOptions are the inputs of a connection test.
type ConnOptions struct {
	Suites []ConnSuite
	// Requests is the number of new connections (one GET each) per suite.
	Requests          int
	Workers           int
	RequestsPerMinute int
	Timeout           time.Duration
	// Label describes where the test runs from, e.g. "laptop" or "test-cluster".
	Label     string
	OutputDir string
	// WANIP overrides the lookup; WANIPLookupURL empty disables the lookup.
	WANIP          string
	WANIPLookupURL string
	Progress       io.Writer
}

// ExecuteConnectionTests runs every suite one after the other. Each suite opens a brand new
// connection (DNS, TCP, TLS) for every request and gets its own run directory and report.
func ExecuteConnectionTests(ctx context.Context, opts ConnOptions) ([]Outcome, error) {
	if len(opts.Suites) == 0 {
		return nil, errors.New("at least one suite is required")
	}

	if opts.Requests < 1 {
		return nil, errors.New("requests must be at least 1")
	}

	wan := opts.WANIP
	if wan == "" && opts.WANIPLookupURL != "" {
		wan = LookupWANIP(ctx, opts.WANIPLookupURL)
	}

	if wan == "" {
		wan = "not looked up"
	}

	outcomes := make([]Outcome, 0, len(opts.Suites))

	for _, suite := range opts.Suites {
		if ctx.Err() != nil {
			break
		}

		if opts.Progress != nil {
			_, _ = fmt.Fprintf(opts.Progress, "Suite %s: %d new connections to %s\n",
				suite.Name, opts.Requests, suite.URL)
		}

		outcome, err := runConnSuite(ctx, opts, suite, wan)
		if err != nil {
			return outcomes, fmt.Errorf("suite %s: %w", suite.Name, err)
		}

		outcomes = append(outcomes, outcome)
	}

	return outcomes, ctx.Err()
}

func runConnSuite(
	ctx context.Context,
	opts ConnOptions,
	suite ConnSuite,
	wan string,
) (Outcome, error) {
	target, err := url.Parse(suite.URL)
	if err != nil {
		return Outcome{}, err
	}

	origin := target.Scheme + "://" + target.Host
	cfg := Config{
		Hosts:             []string{origin},
		Workers:           max(opts.Workers, 1),
		RequestsPerMinute: opts.RequestsPerMinute,
		Timeout:           opts.Timeout,
		Requests:          opts.Requests,
		KeepAlive:         false,
	}

	var (
		sinkMu  sync.Mutex
		results []Result
	)

	run := &runner{
		cfg:     cfg,
		client:  newHTTPClient(cfg),
		limiter: newRateLimiter(cfg.RequestsPerMinute),
		sink: func(r Result) {
			sinkMu.Lock()
			defer sinkMu.Unlock()

			results = append(results, r)
		},
	}

	info := RunInfo{Started: time.Now()}

	var (
		next atomic.Int64
		wg   sync.WaitGroup
	)

	for range cfg.Workers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for ctx.Err() == nil {
				n := int(next.Add(1))
				if n > opts.Requests {
					return
				}

				run.do(ctx, job{iteration: n, host: origin}, KindNewConnection, suite.URL, "", nil)
			}
		}()
	}

	wg.Wait()

	info.Finished = time.Now()
	info.Aborted = ctx.Err() != nil

	label := opts.Label + " / " + suite.Name
	meta := RunMeta{
		Label:    label,
		Hostname: Hostname(),
		WANIP:    wan,
		Config:   cfg,
		Started:  info.Started,
		Finished: info.Finished,
		Aborted:  info.Aborted,
		Type:     TypeConnectionTest,
		Origin:   opts.Label,
		Suite:    suite.Name,
		Target:   suite.URL,
		Control:  suite.Control,
	}
	summary := Summarize(results, info, cfg.Hosts)

	dir := filepath.Join(opts.OutputDir, fmt.Sprintf(
		"%s-%s", info.Started.Format("20060102-150405"), sanitize(label),
	))
	if err := os.MkdirAll(dir, outDirPerm); err != nil {
		return Outcome{}, err
	}

	if err := writeResultsJSONL(filepath.Join(dir, "results.jsonl"), results); err != nil {
		return Outcome{}, err
	}

	if err := writeOutputs(dir, meta, summary, results); err != nil {
		return Outcome{}, err
	}

	return Outcome{Dir: dir, Summary: summary, Meta: meta}, nil
}

func writeResultsJSONL(path string, results []Result) error {
	file, err := os.OpenFile( //nolint:gosec // path is built from the user supplied output dir
		path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, outFilePerm,
	)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)

	for _, res := range results {
		if err := encoder.Encode(res); err != nil {
			return err
		}
	}

	return writer.Flush()
}
