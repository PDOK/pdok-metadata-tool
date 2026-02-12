// Package client holds client logic, i.e. for doing specific CSW or NGR requests.
package client

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"log/slog"

	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/csw"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/iso1911x"
)

// CswClient is used as a client for doing CSW requests.
type CswClient struct {
	endpoint *url.URL
	client   *http.Client
	cacheDir *string
	cacheTTL time.Duration
}

const (
	permDir0750  = 0o750
	permFile0600 = 0o600
)

// NewCswClient creates a new instance of NgrClient.
func NewCswClient(endpoint *url.URL) CswClient {
	const defaultTimeoutSeconds = 20

	client := &http.Client{
		Timeout: defaultTimeoutSeconds * time.Second,
	}

	return CswClient{
		endpoint: endpoint,
		client:   client,
		cacheDir: nil,
	}
}

// SetCache enables on-disk caching of raw CSW records.
// ttlHours is the time-to-live expressed in hours.
func (c *CswClient) SetCache(cacheDir string, ttlHours int) {
	c.cacheDir = &cacheDir
	c.cacheTTL = time.Duration(ttlHours) * time.Hour
}

func (c *CswClient) UnsetCache() {
	c.cacheDir = nil
}

// GetRecordByID returns a metadata record for a given id.
func (c *CswClient) GetRecordByID(uuid string) (iso1911x.MDMetadata, error) {
	raw, err := c.GetRawRecordByID(uuid)
	if err != nil {
		return iso1911x.MDMetadata{}, err
	}

	cswResponse := csw.GetRecordByIDResponse{}
	if err := xml.Unmarshal(
		raw,
		&cswResponse, //nolint:musttag // model types contain tags
	); err != nil {
		return iso1911x.MDMetadata{}, err
	}

	return cswResponse.MDMetadata, nil
}

func (c *CswClient) GetRawRecordByID(uuid string) (rawRecord []byte, err error) {
	// Try cache first when enabled and fresh
	if c.cacheDir != nil {
		if cached, ok, cacheErr := c.getCachedRecordIfFresh(uuid); cacheErr == nil && ok {
			slog.Debug("Harvesting record from cache", "uuid", uuid)

			return cached, nil
		}
		// if cacheErr != nil we ignore and proceed to fetch
	}

	// Fetch from remote
	cswURL := c.getRecordByIDUrl(uuid)
	slog.Debug("Harvesting record from", "url", cswURL)

	rawRecord, err = getResponseBody(cswURL, "GET", nil, *c.client)
	if err != nil {
		return nil, err
	}

	// Store in cache when enabled
	if c.cacheDir != nil {
		_ = c.storeRecordInCache(uuid, rawRecord) // best-effort caching
	}

	return rawRecord, nil
}

// GetRecordPage returns the full CSW GetRecords response for a page, possibly using a constraint.
// TODO Use this for harvesting service metadata in ETF-validator-go.
// Note: Interface changed to return the entire GetRecordsResponse (not only records and next offset).
func (c *CswClient) GetRecordPage(
	constraint *csw.GetRecordsCQLConstraint,
	offset int,
) (csw.GetRecordsResponse, error) {
	if offset == 0 {
		return csw.GetRecordsResponse{}, errors.New("offset must be greater than 0")
	}

	cswURL := c.endpoint.String() +
		"?service=CSW" +
		"&request=GetRecords" +
		"&version=2.0.2" +
		"&typeNames=gmd:MD_Metadata" +
		"&resultType=results" +
		"&startPosition=" + strconv.Itoa(offset) +
		"&maxRecords=50"

	if constraint != nil {
		cswURL += constraint.ToQueryParameter()
	}

	var cswResponse = csw.GetRecordsResponse{}

	err := getUnmarshalledXMLResponse(&cswResponse, cswURL, "GET", nil, *c.client)
	if err != nil {
		return csw.GetRecordsResponse{}, err
	}

	return cswResponse, nil
}

// GetAllRecords returns all metadata records based on recursive paging, possibly using a constraint.
func (c *CswClient) GetAllRecords(
	constraint *csw.GetRecordsCQLConstraint,
) ([]csw.SummaryRecord, error) {
	var result []csw.SummaryRecord

	err := c.getRecordsRecursive(constraint, 1, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *CswClient) HarvestByCQLConstraint(
	constraint *csw.GetRecordsCQLConstraint,
) (result []iso1911x.MDMetadata, err error) {
	records, err := c.GetAllRecords(constraint)
	if err != nil {
		return result, err
	}

	for _, record := range records {
		raw, err := c.GetRawRecordByID(record.Identifier)
		if err != nil {
			slog.Debug(
				"Error retrieving record",
				"identifier",
				record.Identifier,
				"type",
				record.Type,
				"title",
				record.Title,
				"err",
				err,
			)
		}

		cswResponse := csw.GetRecordByIDResponse{}

		err = xml.Unmarshal(raw, &cswResponse) //nolint:musttag // model types contain tags
		if err != nil {
			slog.Debug(
				"Error unmarshalling NGR response",
				"identifier",
				record.Identifier,
				"err",
				err,
			)
		} else {
			result = append(result, cswResponse.MDMetadata)
		}
	}

	return result, err
}

// GetRecordsWithOGCFilter returns summary metadata records, using an OGC filter.
func (c *CswClient) GetRecordsWithOGCFilter(
	filter *csw.GetRecordsOgcFilter,
) ([]csw.SummaryRecord, error) {
	requestBody, err := filter.ToRequestBody()
	if err != nil {
		return nil, err
	}

	var cswResponse = csw.GetRecordsResponse{}

	err = getUnmarshalledXMLResponse(
		&cswResponse,
		c.endpoint.String(),
		"POST",
		&requestBody,
		*c.client,
	)
	if err != nil {
		return nil, err
	}

	return cswResponse.SearchResults.SummaryRecords, nil
}

// --- Helper and unexported methods must be placed after exported methods (funcorder) ---

// --- Helper and unexported methods must be placed after exported methods (funcorder) ---

func (c *CswClient) getRecordByIDUrl(
	uuid string,
) string {
	return c.endpoint.String() +
		"?service=CSW" +
		"&request=GetRecordById" +
		"&version=2.0.2" +
		"&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full" +
		"&id=" + uuid + "#MD_DataIdentification"
}

func (c *CswClient) getCachedRecordIfFresh(uuid string) ([]byte, bool, error) {
	path := c.getCachePath(uuid)

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}

		return nil, false, err
	}

	// Check freshness
	if time.Since(fi.ModTime()) > c.cacheTTL {
		return nil, false, nil
	}

	// #nosec G304 -- reading from a constructed path under controlled cache directory
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}

	return data, true, nil
}

func (c *CswClient) storeRecordInCache(uuid string, data []byte) error {
	if c.cacheDir == nil {
		return nil
	}

	if err := c.ensureCacheDir(); err != nil {
		return err
	}

	path := c.getCachePath(uuid)

	return os.WriteFile(path, data, permFile0600)
}

func (c *CswClient) ensureCacheDir() error {
	if c.cacheDir == nil {
		return nil
	}

	return os.MkdirAll(*c.cacheDir, permDir0750)
}

func (c *CswClient) getCachePath(uuid string) string {
	// simple filename: <uuid>.xml inside cacheDir
	return filepath.Join(*c.cacheDir, uuid+".xml")
}

// getRecordsRecursive recursively pages through all records.
// It guards against changes in NumberOfRecordsMatched by restarting from offset 1 if detected.
func (c *CswClient) getRecordsRecursive(
	constraint *csw.GetRecordsCQLConstraint,
	offset int,
	result *[]csw.SummaryRecord,
) (err error) {
	return c.getRecordsRecursiveWithState(constraint, offset, result, -1, 0)
}

// getRecordsRecursiveWithState is the internal implementation tracking baseline matched and restart attempts.
func (c *CswClient) getRecordsRecursiveWithState(
	constraint *csw.GetRecordsCQLConstraint,
	offset int,
	result *[]csw.SummaryRecord,
	baselineMatched int,
	restarts int,
) (err error) {
	const maxRestarts = 3

	resp, err := c.GetRecordPage(constraint, offset)
	if err != nil {
		return err
	}

	nextOffset, err := strconv.Atoi(resp.SearchResults.NextRecord)
	if err != nil {
		return err
	}

	matched, err := strconv.Atoi(resp.SearchResults.NumberOfRecordsMatched)
	if err != nil {
		return err
	}

	if offset == 1 {
		// First page logging and establish baseline
		baselineMatched = matched
		slog.Info("Found records; paging recursively until end", "total", matched)
	} else {
		// Detect changes in total matches and restart if necessary
		slog.Debug(
			"Paging recursively until end",
			"total",
			matched,
			"baseline",
			baselineMatched,
			"offset",
			offset,
		)

		if baselineMatched >= 0 && matched != baselineMatched {
			if restarts >= maxRestarts {
				slog.Warn(
					"NumberOfRecordsMatched keeps changing; giving up restarts",
					"baseline",
					baselineMatched,
					"current",
					matched,
					"restarts",
					restarts,
				)
			} else {
				slog.Warn(
					"NumberOfRecordsMatched changed during paging; restarting from offset 1",
					"baseline",
					baselineMatched,
					"current",
					matched,
					"restarts",
					restarts+1,
				)
				// Reset results and restart from beginning with new baseline
				*result = (*result)[:0]

				return c.getRecordsRecursiveWithState(constraint, 1, result, matched, restarts+1)
			}
		}
	}

	*result = append(*result, resp.SearchResults.SummaryRecords...)

	if nextOffset == 0 {
		return nil
	}

	return c.getRecordsRecursiveWithState(constraint, nextOffset, result, baselineMatched, restarts)
}
