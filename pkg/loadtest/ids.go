// Package loadtest contains a load/stress test that mimics the requests of the INSPIRE
// (ETF) validator towards the NGR CSW endpoints and reports on the results.
package loadtest

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pdok/pdok-metadata-tool/v2/internal/common"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/csw"
)

const (
	// CswInspirePath is the CSW path the INSPIRE validator uses for all its requests.
	CswInspirePath = "/geonetwork/srv/dut/csw-inspire"
	// TagsPathFormat is the NGR API path to retrieve the tags of a record.
	TagsPathFormat = "/geonetwork/srv/api/records/%s/tags"
	// UserAgent is the User-Agent of the INSPIRE validator, so NGR sees identical traffic.
	UserAgent = "pdok.nl (inspire-etf-validator)"
	// AcceptHeader is the Accept header of the INSPIRE validator.
	AcceptHeader = "*/*;q=0.8,application/signed-exchange"

	idPageSize       = 100
	idScrapeTimeout  = 60 * time.Second
	idCacheDirPerm   = 0o750
	idCacheFilePerm  = 0o600
	maxIDScrapePages = 10000
	serviceCQL       = "type='service' AND organisationName='Beheer PDOK'"
)

// ServiceRecord is a PDOK service metadata record as returned by the CSW.
type ServiceRecord struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type idCache struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Host      string          `json:"host"`
	Records   []ServiceRecord `json:"records"`
}

// IDCacheOptions configures how the service ID cache is used.
type IDCacheOptions struct {
	Path    string
	TTL     time.Duration
	Refresh bool
}

// LoadServiceRecords returns all PDOK service records. A fresh cache is used when available,
// otherwise the records are scraped from the CSW and the cache is rewritten.
// The second return value indicates whether the cache was used.
func LoadServiceRecords(
	ctx context.Context,
	host string,
	opts IDCacheOptions,
) ([]ServiceRecord, bool, error) {
	if !opts.Refresh && opts.Path != "" {
		if cached, ok := readIDCache(opts, host); ok {
			return cached, true, nil
		}
	}

	records, err := ScrapeServiceRecords(ctx, host)
	if err != nil {
		return nil, false, err
	}

	if opts.Path != "" {
		if err := writeIDCache(opts.Path, idCache{
			FetchedAt: time.Now().UTC(),
			Host:      host,
			Records:   records,
		}); err != nil {
			slog.Warn("could not write ID cache", "path", opts.Path, "err", err)
		}
	}

	return records, false, nil
}

func readIDCache(opts IDCacheOptions, host string) ([]ServiceRecord, bool) {
	raw, err := os.ReadFile(opts.Path)
	if err != nil {
		return nil, false
	}

	var cache idCache
	if err := json.Unmarshal(raw, &cache); err != nil {
		return nil, false
	}

	if cache.Host != host || len(cache.Records) == 0 || time.Since(cache.FetchedAt) > opts.TTL {
		return nil, false
	}

	return cache.Records, true
}

func writeIDCache(path string, cache idCache) error {
	if err := os.MkdirAll(filepath.Dir(path), idCacheDirPerm); err != nil {
		return err
	}

	raw, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, raw, idCacheFilePerm)
}

// ScrapeServiceRecords pages through the CSW GetRecords results for all PDOK services,
// following the same query as the INSPIRE validator.
func ScrapeServiceRecords(ctx context.Context, host string) ([]ServiceRecord, error) {
	httpClient := &http.Client{Timeout: idScrapeTimeout}

	var (
		records []ServiceRecord
		seen    = map[string]bool{}
		start   = 1
	)

	for range maxIDScrapePages {
		page, err := getRecordsPage(ctx, httpClient, host, start)
		if err != nil {
			return nil, fmt.Errorf("scraping service records at position %d: %w", start, err)
		}

		for _, rec := range page.SearchResults.SummaryRecords {
			if rec.Identifier == "" || seen[rec.Identifier] {
				continue
			}

			seen[rec.Identifier] = true

			records = append(records, ServiceRecord{ID: rec.Identifier, Title: rec.Title})
		}

		next, _ := strconv.Atoi(page.SearchResults.NextRecord)
		matched, _ := strconv.Atoi(page.SearchResults.NumberOfRecordsMatched)

		// Geonetwork sometimes returns a nextRecord beyond the matched total on the last page.
		if next <= start || next > matched {
			return records, nil
		}

		start = next
	}

	return nil, errors.New("too many pages while scraping service records")
}

func getRecordsPage(
	ctx context.Context,
	httpClient *http.Client,
	host string,
	start int,
) (csw.GetRecordsResponse, error) {
	query := url.Values{}
	query.Set("request", "GetRecords")
	query.Set("service", "CSW")
	query.Set("version", "2.0.2")
	query.Set("typeNames", "gmd:MD_Metadata")
	query.Set("resultType", "results")
	query.Set("constraintLanguage", "CQL_TEXT")
	query.Set("constraint_language_version", "1.1.0")
	query.Set("constraint", serviceCQL)
	query.Set("startPosition", strconv.Itoa(start))
	query.Set("maxRecords", strconv.Itoa(idPageSize))

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		host+CswInspirePath+"?"+query.Encode(),
		nil,
	)
	if err != nil {
		return csw.GetRecordsResponse{}, err
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", AcceptHeader)

	//nolint:bodyclose // closed with common.SafeClose below
	resp, err := httpClient.Do(req)
	if err != nil {
		return csw.GetRecordsResponse{}, err
	}

	defer common.SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return csw.GetRecordsResponse{}, fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return csw.GetRecordsResponse{}, err
	}

	var page csw.GetRecordsResponse
	if err := xml.Unmarshal(body, &page); err != nil { //nolint:musttag // model types contain tags
		return csw.GetRecordsResponse{}, err
	}

	return page, nil
}
