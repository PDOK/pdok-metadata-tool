package loadtest

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/csw"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/ngr"
)

// Request kinds, in the order the INSPIRE validator performs them.
const (
	KindServiceRecord = "GetRecordById (service)"
	KindDatasetTags   = "Tags (dataset)"
	KindDatasetRecord = "GetRecordById (dataset)"
)

// Error categories.
const (
	ErrTimeout      = "timeout"
	ErrDNS          = "dns"
	ErrConnect      = "connect"
	ErrTLS          = "tls"
	ErrConnReset    = "connection reset"
	ErrBodyRead     = "body read"
	ErrHTTPStatus   = "http status"
	ErrCSWException = "csw exception"
	ErrParse        = "parse"
	ErrCanceled     = "canceled"
	ErrOther        = "other"
)

const (
	maxRedirects      = 10
	tlsHandshakeLimit = 10 * time.Second
	dialTimeout       = 30 * time.Second
	keepAlive         = 30 * time.Second
	idleConnTimeout   = 90 * time.Second
	msPerSecond       = 1000.0
	secondsPerMinute  = 60
)

var operatesOnIDRegex = regexp.MustCompile(`(?i)id=([0-9a-zA-Z-]+)#MD_DataIdentification`)

// Config configures a load test run.
type Config struct {
	// Hosts are base URLs (scheme + host) that are tested, e.g. with and without "www.".
	Hosts []string `json:"hosts"`
	// Workers is the number of services that are processed concurrently.
	Workers int `json:"workers"`
	// RequestsPerMinute caps the total request rate over all hosts; 0 means unlimited.
	RequestsPerMinute int `json:"requestsPerMinute"`
	// Timeout is the per request timeout, including reading the body.
	Timeout time.Duration `json:"timeout"`
	// Iterations is how many times all services are processed; 0 means until Duration is reached.
	Iterations int `json:"iterations"`
	// Duration stops scheduling new work after this time; 0 means no limit.
	Duration time.Duration `json:"duration"`
	// IncludeDatasets also fetches tags and records of the datasets a service operates on.
	IncludeDatasets bool `json:"includeDatasets"`
	// MaxDatasetsPerService limits the datasets per service; 0 means all.
	MaxDatasetsPerService int `json:"maxDatasetsPerService"`
	// KeepAlive reuses connections like the validator does; when false every request is a
	// fresh connection (which exposes DNS/connect/TLS timings on every request).
	KeepAlive bool `json:"keepAlive"`
	// Shuffle randomises the service order for every iteration.
	Shuffle bool `json:"shuffle"`
}

// Result is the outcome of a single request.
type Result struct {
	Time          time.Time `json:"time"`
	Iteration     int       `json:"iteration"`
	Host          string    `json:"host"`
	Kind          string    `json:"kind"`
	ServiceID     string    `json:"serviceId"`
	DatasetID     string    `json:"datasetId,omitempty"`
	URL           string    `json:"url"`
	DurationMs    float64   `json:"durationMs"`
	Status        int       `json:"status"`
	Bytes         int64     `json:"bytes"`
	OK            bool      `json:"ok"`
	ErrorCategory string    `json:"errorCategory,omitempty"`
	ErrorPhase    string    `json:"errorPhase,omitempty"`
	Error         string    `json:"error,omitempty"`
	DNSMs         float64   `json:"dnsMs"`
	ConnectMs     float64   `json:"connectMs"`
	TLSMs         float64   `json:"tlsMs"`
	TTFBMs        float64   `json:"ttfbMs"`
	RemoteIP      string    `json:"remoteIp,omitempty"`
	ConnReused    bool      `json:"connReused"`
	Proto         string    `json:"proto,omitempty"`
	Redirects     []string  `json:"redirects,omitempty"`
}

// RunInfo holds the timing of a finished run.
type RunInfo struct {
	Started  time.Time
	Finished time.Time
	// Aborted is true when the run was interrupted before it completed.
	Aborted bool
}

type redirectKey struct{}

type runner struct {
	cfg     Config
	client  *http.Client
	sink    func(Result)
	limiter *rateLimiter
	hostSeq atomic.Int64
	stopAt  time.Time
}

// Run executes the load test and calls sink for every result. sink is called from multiple
// goroutines, but never concurrently.
func Run(
	ctx context.Context,
	cfg Config,
	services []ServiceRecord,
	sink func(Result),
) (RunInfo, error) {
	if len(cfg.Hosts) == 0 {
		return RunInfo{}, errors.New("at least one host is required")
	}

	if len(services) == 0 {
		return RunInfo{}, errors.New("no services to test")
	}

	if cfg.Iterations == 0 && cfg.Duration == 0 {
		return RunInfo{}, errors.New("set at least one of iterations or duration")
	}

	cfg.Workers = max(cfg.Workers, 1)

	var sinkMu sync.Mutex

	run := &runner{
		cfg:     cfg,
		client:  newHTTPClient(cfg),
		limiter: newRateLimiter(cfg.RequestsPerMinute),
		sink: func(r Result) {
			sinkMu.Lock()
			defer sinkMu.Unlock()

			sink(r)
		},
	}

	info := RunInfo{Started: time.Now()}
	if cfg.Duration > 0 {
		run.stopAt = info.Started.Add(cfg.Duration)
	}

	jobs := make(chan job)

	var wg sync.WaitGroup

	for range cfg.Workers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := range jobs {
				run.processService(ctx, j)
			}
		}()
	}

	run.produce(ctx, services, jobs)
	close(jobs)
	wg.Wait()

	info.Finished = time.Now()
	info.Aborted = ctx.Err() != nil

	return info, nil
}

type job struct {
	iteration int
	service   ServiceRecord
	host      string
}

func (r *runner) produce(ctx context.Context, services []ServiceRecord, jobs chan<- job) {
	order := make([]ServiceRecord, len(services))
	copy(order, services)

	for iteration := 1; r.cfg.Iterations == 0 || iteration <= r.cfg.Iterations; iteration++ {
		if r.cfg.Shuffle {
			//nolint:gosec // shuffling the test order needs no secure randomness
			rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
		}

		for _, svc := range order {
			if r.finished(ctx) {
				return
			}

			host := r.cfg.Hosts[(r.hostSeq.Add(1)-1)%int64(len(r.cfg.Hosts))]

			select {
			case jobs <- job{iteration: iteration, service: svc, host: host}:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (r *runner) finished(ctx context.Context) bool {
	if ctx.Err() != nil {
		return true
	}

	return !r.stopAt.IsZero() && time.Now().After(r.stopAt)
}

// processService performs the sequence of requests the validator does for a single service.
func (r *runner) processService(ctx context.Context, j job) {
	svcURL := recordByIDURL(j.host, j.service.ID)

	res, body := r.do(ctx, j, KindServiceRecord, svcURL, "", validateCSWRecord)
	if !res.OK || !r.cfg.IncludeDatasets {
		return
	}

	for _, datasetID := range r.linkedDatasets(body) {
		if r.finished(ctx) {
			return
		}

		tagsRes, tagsBody := r.do(
			ctx, j, KindDatasetTags, j.host+fmt.Sprintf(TagsPathFormat, datasetID), datasetID,
			validateTags,
		)
		if !tagsRes.OK || !hasInspireTag(tagsBody) {
			continue
		}

		r.do(
			ctx,
			j,
			KindDatasetRecord,
			recordByIDURL(j.host, datasetID),
			datasetID,
			validateCSWRecord,
		)
	}
}

func (r *runner) linkedDatasets(serviceBody []byte) []string {
	var parsed csw.GetRecordByIDResponse

	err := xml.Unmarshal(serviceBody, &parsed) //nolint:musttag // model types contain tags
	if err != nil {
		return nil
	}

	svc := parsed.MDMetadata.IdentificationInfo.SVServiceIdentification
	if svc == nil || strings.EqualFold(svc.ServiceType, "discovery") {
		return nil
	}

	var (
		ids  []string
		seen = map[string]bool{}
	)

	for _, op := range svc.OperatesOn {
		match := operatesOnIDRegex.FindStringSubmatch(op.Href)
		if len(match) < 2 || seen[match[1]] { //nolint:mnd // match[0] is the full match
			continue
		}

		seen[match[1]] = true

		ids = append(ids, match[1])
		if r.cfg.MaxDatasetsPerService > 0 && len(ids) >= r.cfg.MaxDatasetsPerService {
			break
		}
	}

	return ids
}

func hasInspireTag(body []byte) bool {
	var tags ngr.RecordTagsResponse
	if err := json.Unmarshal(body, &tags); err != nil {
		return false
	}

	for _, tag := range tags {
		if strings.EqualFold(tag.Name, "inspire") {
			return true
		}
	}

	return false
}

func recordByIDURL(host, id string) string {
	query := url.Values{}
	query.Set("service", "CSW")
	query.Set("request", "GetRecordById")
	query.Set("version", "2.0.2")
	query.Set("outputSchema", "http://www.isotc211.org/2005/gmd")
	query.Set("elementSetName", "full")
	query.Set("id", id)

	return host + CswInspirePath + "?" + query.Encode() + "#MD_DataIdentification"
}

func validateCSWRecord(body []byte) (string, string) {
	if strings.Contains(string(body), "ExceptionReport") {
		return ErrCSWException, "CSW returned an ExceptionReport"
	}

	return "", ""
}

func validateTags(body []byte) (string, string) {
	var tags ngr.RecordTagsResponse
	if err := json.Unmarshal(body, &tags); err != nil {
		return ErrParse, "tags response is not valid JSON: " + err.Error()
	}

	return "", ""
}

// do performs one request, measures it and reports the result to the sink.
func (r *runner) do(
	ctx context.Context,
	j job,
	kind, rawURL, datasetID string,
	validate func([]byte) (string, string),
) (Result, []byte) {
	res := Result{
		Iteration: j.iteration,
		Host:      j.host,
		Kind:      kind,
		ServiceID: j.service.ID,
		DatasetID: datasetID,
		URL:       rawURL,
	}

	if err := r.limiter.wait(ctx); err != nil {
		res.Time = time.Now()
		res.ErrorCategory, res.Error = ErrCanceled, err.Error()

		return res, nil
	}

	body := r.execute(ctx, rawURL, &res)
	if res.OK && validate != nil {
		if category, msg := validate(body); category != "" {
			res.OK = false
			res.ErrorCategory, res.Error = category, msg
		}
	}

	r.sink(res)

	return res, body
}

type traceState struct {
	start         time.Time
	dnsStart      time.Time
	connectStart  time.Time
	tlsStart      time.Time
	wroteRequest  time.Time
	phase         atomic.Value
	res           *Result
	mu            sync.Mutex
	redirectChain []string
}

func (t *traceState) setPhase(p string) { t.phase.Store(p) }

func (t *traceState) trace() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		DNSStart: func(httptrace.DNSStartInfo) {
			t.setPhase("dns")
			t.mu.Lock()
			t.dnsStart = time.Now()
			t.mu.Unlock()
		},
		DNSDone: func(httptrace.DNSDoneInfo) {
			t.mu.Lock()
			t.res.DNSMs += ms(time.Since(t.dnsStart))
			t.mu.Unlock()
		},
		ConnectStart: func(string, string) {
			t.setPhase("connect")
			t.mu.Lock()
			if t.connectStart.IsZero() {
				t.connectStart = time.Now()
			}
			t.mu.Unlock()
		},
		ConnectDone: func(_, addr string, err error) {
			if err != nil {
				return
			}

			t.mu.Lock()
			t.res.ConnectMs = ms(time.Since(t.connectStart))
			t.res.RemoteIP = hostOnly(addr)
			t.mu.Unlock()
		},
		TLSHandshakeStart: func() {
			t.setPhase("tls")
			t.mu.Lock()
			t.tlsStart = time.Now()
			t.mu.Unlock()
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			t.mu.Lock()
			t.res.TLSMs = ms(time.Since(t.tlsStart))
			t.mu.Unlock()
		},
		GotConn: func(info httptrace.GotConnInfo) {
			t.setPhase("request")
			t.mu.Lock()

			t.res.ConnReused = info.Reused
			if addr := info.Conn.RemoteAddr(); addr != nil {
				t.res.RemoteIP = hostOnly(addr.String())
			}
			t.mu.Unlock()
		},
		WroteRequest: func(httptrace.WroteRequestInfo) {
			t.setPhase("wait response")
			t.mu.Lock()
			t.wroteRequest = time.Now()
			t.mu.Unlock()
		},
		GotFirstResponseByte: func() {
			t.setPhase("read body")
			t.mu.Lock()
			t.res.TTFBMs = ms(time.Since(t.wroteRequest))
			t.mu.Unlock()
		},
	}
}

func (r *runner) execute(ctx context.Context, rawURL string, res *Result) []byte {
	reqCtx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()

	state := &traceState{res: res, start: time.Now()}
	state.setPhase("start")

	res.Time = state.start

	reqCtx = httptrace.WithClientTrace(reqCtx, state.trace())
	reqCtx = context.WithValue(reqCtx, redirectKey{}, state)

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, rawURL, nil)
	if err != nil {
		res.ErrorCategory, res.Error = ErrOther, err.Error()

		return nil
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", AcceptHeader)

	resp, err := r.client.Do(req)
	if err != nil {
		r.fail(res, state, err, errors.Is(reqCtx.Err(), context.DeadlineExceeded))

		return nil
	}

	defer func() { _ = resp.Body.Close() }()

	res.Status = resp.StatusCode
	res.Proto = resp.Proto
	res.Redirects = state.redirectChain

	body, err := io.ReadAll(resp.Body)
	res.DurationMs = ms(time.Since(state.start))
	res.Bytes = int64(len(body))

	if err != nil {
		r.fail(res, state, err, errors.Is(reqCtx.Err(), context.DeadlineExceeded))

		return nil
	}

	if resp.StatusCode != http.StatusOK {
		res.ErrorCategory = ErrHTTPStatus
		res.Error = fmt.Sprintf("HTTP %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))

		return body
	}

	res.OK = true

	return body
}

func (r *runner) fail(res *Result, state *traceState, err error, timedOut bool) {
	res.DurationMs = ms(time.Since(state.start))
	res.Error = err.Error()

	if phase, ok := state.phase.Load().(string); ok {
		res.ErrorPhase = phase
	}

	res.ErrorCategory = classify(err, timedOut, res.ErrorPhase)
}

func classify(err error, timedOut bool, phase string) string {
	var (
		dnsErr  *net.DNSError
		netErr  net.Error
		opErr   *net.OpError
		certErr *tls.CertificateVerificationError
	)

	switch {
	case errors.Is(err, context.Canceled) && !timedOut:
		return ErrCanceled
	case errors.Is(err, context.DeadlineExceeded),
		errors.As(err, &netErr) && netErr.Timeout():
		return ErrTimeout
	case errors.As(err, &dnsErr):
		return ErrDNS
	case errors.As(err, &certErr), phase == "tls":
		return ErrTLS
	case errors.Is(err, io.ErrUnexpectedEOF) || strings.Contains(err.Error(), "connection reset"):
		return ErrConnReset
	case errors.As(err, &opErr) && opErr.Op == "dial", phase == "connect":
		return ErrConnect
	case phase == "read body":
		return ErrBodyRead
	case errors.Is(err, io.EOF):
		return ErrConnReset
	default:
		return ErrOther
	}
}

func newHTTPClient(cfg Config) *http.Client {
	dialer := &net.Dialer{Timeout: dialTimeout, KeepAlive: keepAlive}

	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         dialer.DialContext,
		ForceAttemptHTTP2:   true,
		TLSHandshakeTimeout: tlsHandshakeLimit,
		IdleConnTimeout:     idleConnTimeout,
		MaxIdleConnsPerHost: max(cfg.Workers, 2), //nolint:mnd // default idle conns per host
		DisableKeepAlives:   !cfg.KeepAlive,
	}

	return &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return errors.New("stopped after too many redirects")
			}

			if state, ok := req.Context().Value(redirectKey{}).(*traceState); ok {
				status := 0
				if req.Response != nil {
					status = req.Response.StatusCode
				}

				state.mu.Lock()
				state.redirectChain = append(
					state.redirectChain,
					fmt.Sprintf("%d -> %s", status, req.URL.String()),
				)
				state.mu.Unlock()
			}

			return nil
		},
	}
}

func ms(d time.Duration) float64 { return float64(d.Microseconds()) / msPerSecond }

func hostOnly(addr string) string {
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}

	return addr
}

// rateLimiter hands out evenly spaced request slots shared by all workers.
type rateLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	next     time.Time
}

func newRateLimiter(requestsPerMinute int) *rateLimiter {
	if requestsPerMinute <= 0 {
		return &rateLimiter{}
	}

	return &rateLimiter{
		interval: time.Duration(secondsPerMinute) * time.Second / time.Duration(requestsPerMinute),
	}
}

func (l *rateLimiter) wait(ctx context.Context) error {
	if l.interval == 0 {
		return ctx.Err()
	}

	l.mu.Lock()

	now := time.Now()
	if l.next.Before(now) {
		l.next = now
	}

	slot := l.next
	l.next = l.next.Add(l.interval)
	l.mu.Unlock()

	timer := time.NewTimer(time.Until(slot))
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
