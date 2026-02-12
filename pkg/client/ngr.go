package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pdok/pdok-metadata-tool/v2/internal/common"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/ngr"
)

type NgrClient struct { //nolint:recvcheck
	NgrClient *http.Client
	NgrConfig *NgrConfig
}

type NgrConfig struct {
	NgrUrl      *string
	NgrUserName *string
	NgrPassword *string
}

const API_RECORDS_TEMPLATE = "/geonetwork/srv/api/records"
const API_LOGIN_PART = "/geonetwork/srv/dut/info?type=me"
const INSPIRE_TAG = 224342
const NGR_CLIENT_TIMEOUT = 20 * time.Second

func NewNgrClient(config NgrConfig) NgrClient {
	client := &http.Client{
		Timeout: NGR_CLIENT_TIMEOUT,
	}

	return NgrClient{
		NgrConfig: &config,
		NgrClient: client,
	}
}

// TODO Use this for harvesting only INSPIRE service metadata in ETF-validator-go
// GetRecordTags does a GET request on NGR to get the Tags for a record.
func (c *NgrClient) GetRecordTags(uuid string) (ngr.RecordTagsResponse, error) {
	mdTagUrl := fmt.Sprintf("%s/geonetwork/srv/api/records/%s/tags", *c.NgrConfig.NgrUrl, uuid)

	recordTagsResponse := ngr.RecordTagsResponse{}

	err := getUnmarshalledJSONResponse(&recordTagsResponse, mdTagUrl, *c.NgrClient)
	if err != nil {
		return nil, err
	}

	return recordTagsResponse, nil
}

// CreateOrUpdateServiceMetadataRecord does a PUT request on NGR to create or update a record.
func (c *NgrClient) CreateOrUpdateServiceMetadataRecord(
	record string,
	categoryId *string,
	groupId *string,
	toBePublished bool,
) error {
	// Build URL with query params
	params := "?metadataType=METADATA&uuidProcessing=REMOVE_AND_REPLACE"
	if toBePublished {
		params += "&publishToAll=true"
	} else {
		params += "&publishToAll=false"
	}

	if groupId != nil {
		params += "&group=" + *groupId
	}

	if categoryId != nil {
		params += "&category=" + *categoryId
	}

	url := fmt.Sprintf(
		"%s%s/%s",
		*c.NgrConfig.NgrUrl,
		API_RECORDS_TEMPLATE,
		params,
	)

	_, err := getNgrResponseBody(
		c.NgrConfig,
		url,
		http.MethodPut,
		&record,
		*c.NgrClient,
		ContentTypeXML,
	)

	return err
}

// GetRecord does a GET request on NGR to get a record.
func (c *NgrClient) GetRecord(uuid string) (string, error) {
	ngrUrl := fmt.Sprintf("%s%s/%s",
		*c.NgrConfig.NgrUrl,
		API_RECORDS_TEMPLATE,
		uuid,
	)

	var responseBodyString = ""

	responseBodyByteArr, err := getNgrResponseBody(
		c.NgrConfig,
		ngrUrl,
		http.MethodGet,
		nil,
		*c.NgrClient,
		ContentTypeXML,
	)

	if responseBodyByteArr != nil {
		responseBodyString = string(responseBodyByteArr)
	}

	return responseBodyString, err
}

// DeleteRecord does  DELETE request on NGR to delete a record.
func (c *NgrClient) DeleteRecord(uuid string) error {
	ngrUrl := fmt.Sprintf("%s%s/%s",
		*c.NgrConfig.NgrUrl,
		API_RECORDS_TEMPLATE,
		uuid,
	)
	_, err := getNgrResponseBody(
		c.NgrConfig,
		ngrUrl,
		http.MethodDelete,
		nil,
		*c.NgrClient,
		ContentTypeXML,
	)

	return err
}

// AddTagToRecord does a PUT request on NGR to add a Tag to a record.
func (c *NgrClient) AddTagToRecord(uuid string, tagId int) error {
	ngrUrl := fmt.Sprintf("%s%s/%s/tags?id=%d",
		*c.NgrConfig.NgrUrl,
		API_RECORDS_TEMPLATE,
		uuid,
		tagId,
	)
	_, err := getNgrResponseBody(
		c.NgrConfig,
		ngrUrl,
		http.MethodPut,
		nil,
		*c.NgrClient,
		ContentTypeXML,
	)

	return err
}

// ValidateRecord does a PUT request on NGR to validate a record.
func (c *NgrClient) ValidateRecord(uuid string) (ValidationResult, error) {
	ngrUrl := fmt.Sprintf("%s%s/validate?uuids=%s",
		*c.NgrConfig.NgrUrl,
		API_RECORDS_TEMPLATE,
		uuid,
	)

	response, err := getNgrResponseBody(
		c.NgrConfig,
		ngrUrl,
		http.MethodPut,
		nil,
		*c.NgrClient,
		ContentTypeJSON,
	)
	if err != nil {
		return ValidationResult{}, err
	}

	ngrResponse := ValidationResult{}
	if err = json.Unmarshal(response, &ngrResponse); err != nil {
		return ValidationResult{}, err
	}

	return ngrResponse, nil
}

func obtainXSRFToken(ngrConfig *NgrConfig) (string, error) {
	// Build URL
	url := *ngrConfig.NgrUrl + API_LOGIN_PART

	// Create HTTP client and request
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", fmt.Errorf("cannot create request: %w", err)
	}

	username := ngrConfig.NgrUserName
	password := ngrConfig.NgrPassword
	req.SetBasicAuth(*username, *password)

	client := &http.Client{}
	//nolint:bodyclose // We use common.SafeClose to handle closing the response body
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}

	defer common.SafeClose(resp.Body)

	// Look for 403 Forbidden
	if resp.StatusCode == http.StatusForbidden { //nolint:nestif
		// Parse cookies from headers
		cookies := resp.Header.Values("Set-Cookie")

		return getCookieValueByName(cookies, "XSRF-TOKEN")
	}

	return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
