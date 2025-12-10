// Package client holds client logic, i.e. for doing specific CSW or NGR requests.
package client

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
	"github.com/pdok/pdok-metadata-tool/pkg/model/iso1911x"
)

// CswClient is used as a client for doing CSW requests.
type CswClient struct {
	endpoint *url.URL
	client   *http.Client
	cacheDir *string
	cacheTTL time.Duration
}

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

func (c *CswClient) getRecordByIDUrl(uuid string) string {
	return c.endpoint.String() +
		"?service=CSW" +
		"&request=GetRecordById" +
		"&version=2.0.2" +
		"&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full" +
		"&id=" + uuid + "#MD_DataIdentification"
}

// GetRecordByID returns a metadata record for a given id.
func (c *CswClient) GetRecordByID(uuid string) (iso1911x.MDMetadata, error) {
	cswURL := c.getRecordByIDUrl(uuid)

	// todo: use cache trough GetRawRecordByID
	cswResponse := csw.GetRecordByIDResponse{}

	err := getUnmarshalledXMLResponse(&cswResponse, cswURL, "GET", nil, *c.client)
	if err != nil {
		return iso1911x.MDMetadata{}, err
	}

	cswResponse.MDMetadata.SelfURL = cswURL

	return cswResponse.MDMetadata, nil
}

func (c *CswClient) GetRawRecordByID(uuid string) (rawRecord []byte, err error) {
	// Try cache first when enabled and fresh
	if c.cacheDir != nil {
		if cached, ok, cacheErr := c.getCachedRecordIfFresh(uuid); cacheErr == nil && ok {
			fmt.Println("Using cached record")
			return cached, nil
		}
		// if cacheErr != nil we ignore and proceed to fetch
	}

	// Fetch from remote
	cswURL := c.getRecordByIDUrl(uuid)
	fmt.Println("Using NGR: " + cswURL)
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

// getCachedRecordIfFresh returns cached bytes and true when cache exists and is within TTL.
// When cache does not exist or is stale, returns (nil, false, nil).
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

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}
	return data, true, nil
}

// storeRecordInCache ensures dir exists and writes bytes to cache file.
func (c *CswClient) storeRecordInCache(uuid string, data []byte) error {
	if c.cacheDir == nil {
		return nil
	}
	if err := c.ensureCacheDir(); err != nil {
		return err
	}
	path := c.getCachePath(uuid)
	return os.WriteFile(path, data, 0o644)
}

func (c *CswClient) ensureCacheDir() error {
	if c.cacheDir == nil {
		return nil
	}
	return os.MkdirAll(*c.cacheDir, 0o755)
}

func (c *CswClient) getCachePath(uuid string) string {
	// simple filename: <uuid>.xml inside cacheDir
	return filepath.Join(*c.cacheDir, uuid+".xml")
}

// GetRecordPage returns one page of summary metadata, possibly using a constraint.
// TODO Use this for harvesting service metadata in ETF-validator-go.
// todo: when the NumberOfRecordsMatched changed while paging, this could upset the paging. We should check during paging if the number of records has changed. If so we should restart paging from 1.
func (c *CswClient) GetRecordPage(constraint *csw.GetRecordsCQLConstraint, offset int) ([]csw.SummaryRecord, int, error) {
	if offset == 0 {
		offset = 1
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
		return nil, -1, err
	}

	nextRecord, err := strconv.Atoi(cswResponse.SearchResults.NextRecord)
	if err != nil {
		return nil, -1, err
	}

	return cswResponse.SearchResults.SummaryRecords, nextRecord, nil
}

func (c *CswClient) getRecordsRecursive(constraint *csw.GetRecordsCQLConstraint, offset int, result *[]csw.SummaryRecord) (err error) {

	fmt.Printf("recursing with offset %d\n", offset)
	records, nextOffset, err := c.GetRecordPage(constraint, offset)
	if err != nil {
		return err
	}

	*result = append(*result, records...)

	if nextOffset == 0 {
		return
	}

	return c.getRecordsRecursive(constraint, nextOffset, result)
}

// GetAllRecords returns all metadata records based on recursive paging, possibly using a constraint.
func (c *CswClient) GetAllRecords(constraint *csw.GetRecordsCQLConstraint) ([]csw.SummaryRecord, error) {

	var result []csw.SummaryRecord
	err := c.getRecordsRecursive(constraint, 1, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *CswClient) HarvestByCQLConstraint(constraint *csw.GetRecordsCQLConstraint) (result []iso1911x.MDMetadata, err error) {
	records, err := c.GetAllRecords(constraint)

	fmt.Printf("Found %d records for pdok services\n", len(records))
	if err != nil {
		return
	}

	for _, record := range records {
		fmt.Printf("Harvesting record: %s\n", record.Identifier)
		raw, err := c.GetRawRecordByID(record.Identifier)
		if err != nil {
			fmt.Printf("Error retrieving %s\t%s\t%s\n", record.Identifier, record.Type, record.Title)
		}

		cswResponse := csw.GetRecordByIDResponse{}
		err = xml.Unmarshal(raw, &cswResponse)
		if err != nil {
			fmt.Printf("error unmarshalling NGR response from record %s: %s", record.Identifier, err)
		} else {
			result = append(result, cswResponse.MDMetadata)
		}

	}

	return
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

	err = getUnmarshalledXMLResponse(&cswResponse, c.endpoint.String(), "POST", &requestBody, *c.client)
	if err != nil {
		return nil, err
	}

	return cswResponse.SearchResults.SummaryRecords, nil
}
