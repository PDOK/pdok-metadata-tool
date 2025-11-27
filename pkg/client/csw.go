// Package client holds client logic, i.e. for doing specific CSW or NGR requests.
package client

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
)

// CswClient is used as a client for doing CSW requests.
type CswClient struct {
	host   *url.URL
	client *http.Client
}

// NewCswClient creates a new instance of NgrClient.
func NewCswClient(host *url.URL) CswClient {
	const defaultTimeoutSeconds = 20

	client := &http.Client{
		Timeout: defaultTimeoutSeconds * time.Second,
	}

	return CswClient{
		host:   host,
		client: client,
	}
}

// GetRecordByID returns a metadata record for a given id.
func (c CswClient) GetRecordByID(uuid string) (csw.MDMetadata, error) {
	cswURL := c.host.String() +
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

// GetRecords returns summary metadata records, possibly using a constraint.
// TODO Use this for harvesting service metadata in ETF-validator-go.
func (c CswClient) GetRecords(
	constraint *csw.GetRecordsCQLConstraint,
	offset int,
) ([]csw.SummaryRecord, int, error) {
	cswURL := c.host.String() +
		"?service=CSW" +
		"&request=GetRecords" +
		"&version=2.0.2" +
		"&typeNames=gmd:MD_Metadata" +
		"&resultType=results" +
		"&startPosition=" + strconv.Itoa(offset)

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

// GetRecordsWithOGCFilter returns summary metadata records, using an OGC filter.
func (c CswClient) GetRecordsWithOGCFilter(
	filter *csw.GetRecordsOgcFilter,
) ([]csw.SummaryRecord, error) {
	requestBody, err := filter.ToRequestBody()
	if err != nil {
		return nil, err
	}

	var cswResponse = csw.GetRecordsResponse{}

	err = getUnmarshalledXMLResponse(&cswResponse, c.host.String(), "POST", &requestBody, *c.client)
	if err != nil {
		return nil, err
	}

	return cswResponse.SearchResults.SummaryRecords, nil
}
