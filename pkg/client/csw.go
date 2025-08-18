package client

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"pdok-metadata-tool/pkg/model/csw"
)

type CswClient struct {
	host     *url.URL
	endpoint string
	client   *http.Client
}

func NewCswClient(host *url.URL) CswClient {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	return CswClient{
		host:   host,
		client: client,
	}
}

func (c CswClient) GetRecordById(uuid string, logPrefix string) (csw.MDMetadata, error) {
	cswUrl := c.host.String() +
		"?service=CSW" +
		"&request=GetRecordById" +
		"&version=2.0.2" +
		"&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full" +
		"&id=" + uuid + "#MD_DataIdentification"

	cswResponse := csw.GetRecordByIdResponse{}
	err := getUnmarshalledXMLResponse(&cswResponse, cswUrl, *c.client, logPrefix)
	if err != nil {
		return csw.MDMetadata{}, err
	}
	cswResponse.MDMetadata.SelfUrl = cswUrl
	return cswResponse.MDMetadata, nil
}

func (c CswClient) GetRecords(constraint *csw.GetRecordsConstraint, offset int, logPrefix string) ([]csw.SummaryRecord, int, error) {
	cswUrl := c.host.String() +
		"?service=CSW" +
		"&request=GetRecords" +
		"&version=2.0.2" +
		"&typeNames=gmd:MD_Metadata" +
		"&resultType=results" +
		"&startPosition=" + strconv.Itoa(offset)

	if constraint != nil {
		cswUrl += constraint.ToQueryParameter()
	}

	var cswResponse = csw.GetRecordsResponse{}
	err := getUnmarshalledXMLResponse(&cswResponse, cswUrl, *c.client, logPrefix)
	if err != nil {
		return nil, -1, err
	}

	nextRecord, err := strconv.Atoi(cswResponse.SearchResults.NextRecord)
	if err != nil {
		return nil, -1, err
	}

	return cswResponse.SearchResults.SummaryRecords, nextRecord, nil
}
