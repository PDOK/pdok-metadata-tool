// Package client holds client logic, i.e. for doing specific CSW or NGR requests.
package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
)

// CswClient is used as a client for doing CSW requests.
type CswClient struct {
	endpoint *url.URL
	client   *http.Client
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
	}
}

// GetRecordByID returns a metadata record for a given id.
func (c CswClient) GetRecordByID(uuid string) (csw.MDMetadata, error) {
	cswURL := c.endpoint.String() +
		"?service=CSW" +
		"&request=GetRecordById" +
		"&version=2.0.2" +
		"&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full" +
		"&id=" + uuid + "#MD_DataIdentification"

	cswResponse := csw.GetRecordByIDResponse{}

	err := getUnmarshalledXMLResponse(&cswResponse, cswURL, "GET", nil, *c.client)
	if err != nil {
		return csw.MDMetadata{}, err
	}

	cswResponse.MDMetadata.SelfURL = cswURL

	return cswResponse.MDMetadata, nil
}

// GetRecordPage returns one page of summary metadata, possibly using a constraint.
// TODO Use this for harvesting service metadata in ETF-validator-go.
// todo: when the NumberOfRecordsMatched changed while paging, this could upset the paging. We should check during paging if the number of records has changed. If so we should restart paging from 1.
func (c CswClient) GetRecordPage(constraint *csw.GetRecordsCQLConstraint, offset int) ([]csw.SummaryRecord, int, error) {
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

func (c CswClient) getRecordsRecursive(constraint *csw.GetRecordsCQLConstraint, offset int, result *[]csw.SummaryRecord) (err error) {

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
func (c CswClient) GetAllRecords(constraint *csw.GetRecordsCQLConstraint) ([]csw.SummaryRecord, error) {

	var result []csw.SummaryRecord
	err := c.getRecordsRecursive(constraint, 1, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetRecordsWithOGCFilter returns summary metadata records, using an OGC filter.
func (c CswClient) GetRecordsWithOGCFilter(
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
