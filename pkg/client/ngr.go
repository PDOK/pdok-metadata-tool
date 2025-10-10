package client

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pdok/pdok-metadata-tool/pkg/model/ngr"
)

type NgrClient struct {
	NgrClient *http.Client
	NgrConfig NgrConfig
}

type NgrConfig struct {
	NgrUrl      string
	NgrUserName string
	NgrPassword string
}

const API_RECORDS_TEMPlATE = "/geonetwork/srv/api/records"
const API_LOGIN_PART = "/geonetwork/srv/dut/info?type=me"
const INSPIRE_TAG = 224342

func NewNgrClient(config NgrConfig) NgrClient {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	return NgrClient{
		NgrConfig: config,
		NgrClient: client,
	}
}

func (c *NgrClient) getHeaders(ngrConfig *NgrConfig) (http.Header, error) {
	xsrfToken, err := obtainXSRFToken(ngrConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain XSRF token: %w", err)
	}

	// In the old selfservice was a check to block update for published records.
	// In the new selfservice we need to update the record regardless of publish state

	username := c.NgrConfig.NgrUserName
	password := c.NgrConfig.NgrPassword
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	headers := http.Header{}
	headers.Set("Authorization", "Basic "+auth)
	headers.Set("X-XSRF-TOKEN", xsrfToken)
	headers.Set("Cookie", "XSRF-TOKEN="+xsrfToken)
	headers.Set("Content-Type", "application/xml")
	headers.Set("Accept", "application/json, application/xml")
	return headers, nil
}

// TODO Use this for harvesting only INSPIRE service metadata in ETF-validator-go
func (c NgrClient) GetRecordTags(uuid string) (ngr.RecordTagsResponse, error) {
	mdTagUrl := fmt.Sprintf("%s/geonetwork/srv/api/records/%s/tags", c.NgrConfig.NgrUrl, uuid)

	recordTagsResponse := ngr.RecordTagsResponse{}

	err := getUnmarshalledJSONResponse(&recordTagsResponse, mdTagUrl, *c.NgrClient)
	if err != nil {
		return nil, err
	}
	return recordTagsResponse, nil
}

func obtainXSRFToken(ngrConfig *NgrConfig) (string, error) {
	// Build URL
	url := ngrConfig.NgrUrl + API_LOGIN_PART

	// Create HTTP client and request
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", fmt.Errorf("cannot create request: %w", err)
	}

	// Set Basic Auth
	username := ngrConfig.NgrUserName
	password := ngrConfig.NgrPassword
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// In Kotlin, success (200) is considered wrong and should throw an error
	if resp.StatusCode < 400 {
		return "", errors.New("cannot obtain XSRF token, server shouldn't allow POST action")
	}

	// Look for 403 Forbidden
	if resp.StatusCode == http.StatusForbidden {
		// Parse cookies from headers
		cookies := resp.Header.Values("Set-Cookie")
		for _, c := range cookies {
			if strings.Contains(c, "XSRF-TOKEN=") {
				parts := strings.Split(c, ";")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if strings.HasPrefix(part, "XSRF-TOKEN=") {
						split := strings.SplitN(part, "=", 2)
						if len(split) == 2 {
							return split[1], nil
						}
					}
				}
			}
		}
		return "", errors.New("cannot obtain XSRF token from cookie")
	}

	return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func (c *NgrClient) createOrUpdateServiceMetadataRecord(
	record string,
	categoryId *string,
	groupId *string,
	ngrConfig *NgrConfig,
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
		c.NgrConfig.NgrUrl,
		API_RECORDS_TEMPlATE,
		params,
	)

	// PUT the record
	putReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(record))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}
	headers, err := c.getHeaders(ngrConfig)
	if err != nil {
		return fmt.Errorf("failed to create headers request: %w", err)
	}
	putReq.Header = headers
	putResp, err := c.NgrClient.Do(putReq)
	if err != nil {
		return fmt.Errorf("PUT request error: %w", err)
	}
	defer putResp.Body.Close()
	body, _ := io.ReadAll(putResp.Body)
	if putResp.StatusCode < 200 || putResp.StatusCode >= 300 {
		return fmt.Errorf(
			"PUT request failed with status %d: %s",
			putResp.StatusCode,
			string(body),
		)
	}
	return nil
}

func (c *NgrClient) getRecord(uuid string, ngrConfig *NgrConfig) (string, int, error) {
	ngrUrl := fmt.Sprintf("%s%s/%s",
		c.NgrConfig.NgrUrl,
		API_RECORDS_TEMPlATE,
		uuid,
	)
	bodyString := ""

	getReq, err := http.NewRequest(http.MethodGet, ngrUrl, nil)

	headers, err := c.getHeaders(ngrConfig)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("failed to create headers request: %w", err)
	}
	getReq.Header = headers
	getResp, err := c.NgrClient.Do(getReq)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("failed to send http request: %w", err)
	}
	if err == nil {
		defer getResp.Body.Close()
		switch getResp.StatusCode {
		case http.StatusOK:
			bodyBytes, err := io.ReadAll(getResp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString = string(bodyBytes)
			return bodyString, http.StatusOK, nil
		case http.StatusNotFound:
		case http.StatusForbidden:
		default:
			return bodyString, getResp.StatusCode, fmt.Errorf(
				"unexpected http status %d retrieving sharing info for record %s",
				getResp.StatusCode,
				uuid,
			)
		}
	}
	return bodyString, getResp.StatusCode, fmt.Errorf(
		"unexpected http status %d retrieving sharing info for record %s",
		getResp.StatusCode,
		uuid,
	)
}

func (c *NgrClient) deleteRecord(uuid string, ngrConfig *NgrConfig) (int, error) {
	ngrUrl := fmt.Sprintf("%s%s/%s",
		c.NgrConfig.NgrUrl,
		API_RECORDS_TEMPlATE,
		uuid,
	)
	deleteReq, err := http.NewRequest(http.MethodDelete, ngrUrl, nil)

	headers, err := c.getHeaders(ngrConfig)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create headers request: %w", err)
	}
	deleteReq.Header = headers
	deleteResp, err := c.NgrClient.Do(deleteReq)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send http request: %w", err)
	}
	if err == nil {
		defer deleteResp.Body.Close()
		switch deleteResp.StatusCode {
		case http.StatusNoContent:
			return http.StatusNoContent, nil
		case http.StatusNotFound:
		case http.StatusForbidden:
		default:
			return deleteResp.StatusCode, fmt.Errorf(
				"unexpected http status %d retrieving sharing info for record %s",
				deleteResp.StatusCode,
				uuid,
			)
		}
	}
	return deleteResp.StatusCode, fmt.Errorf(
		"unexpected http status %d retrieving sharing info for record %s",
		deleteResp.StatusCode,
		uuid,
	)
}

func (c *NgrClient) addTagToRecord(uuid string, ngrConfig *NgrConfig, tagId int) (int, error) {
	ngrUrl := fmt.Sprintf("%s%s/%s/tags?id=%d",
		c.NgrConfig.NgrUrl,
		API_RECORDS_TEMPlATE,
		uuid,
		tagId,
	)
	putTagReq, err := http.NewRequest(http.MethodPut, ngrUrl, nil)

	headers, err := c.getHeaders(ngrConfig)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create headers request: %w", err)
	}
	putTagReq.Header = headers
	putTagResp, err := c.NgrClient.Do(putTagReq)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send http request: %w", err)
	}
	if err == nil {
		defer putTagResp.Body.Close()
		switch putTagResp.StatusCode {
		case http.StatusCreated:
			return http.StatusCreated, nil
		case http.StatusNotFound:
		case http.StatusForbidden:
		default:
			return putTagResp.StatusCode, fmt.Errorf(
				"unexpected http status %d retrieving sharing info for record %s",
				putTagResp.StatusCode,
				uuid,
			)
		}
	}
	return putTagResp.StatusCode, fmt.Errorf(
		"unexpected http status %d retrieving sharing info for record %s",
		putTagResp.StatusCode,
		uuid,
	)
}

func (c *NgrClient) getTagsByRecord(uuid string, ngrConfig *NgrConfig) (string, int, error) {
	ngrUrl := fmt.Sprintf("%s%s/%s/tags",
		c.NgrConfig.NgrUrl,
		API_RECORDS_TEMPlATE,
		uuid,
	)
	bodyString := ""
	getTagReq, err := http.NewRequest(http.MethodGet, ngrUrl, nil)

	headers, err := c.getHeaders(ngrConfig)
	if err != nil {
		return bodyString, http.StatusInternalServerError, fmt.Errorf("failed to create headers request: %w", err)
	}
	getTagReq.Header = headers
	getTagResp, err := c.NgrClient.Do(getTagReq)
	if err != nil {
		return bodyString, http.StatusInternalServerError, fmt.Errorf("failed to send http request: %w", err)
	}
	if err == nil {
		defer getTagResp.Body.Close()
		switch getTagResp.StatusCode {
		case http.StatusOK:
			bodyBytes, err := io.ReadAll(getTagResp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString = string(bodyBytes)
			return bodyString, http.StatusOK, nil
		case http.StatusNotFound:
		case http.StatusForbidden:
		default:
			return bodyString, getTagResp.StatusCode, fmt.Errorf(
				"unexpected http status %d retrieving sharing info for record %s",
				getTagResp.StatusCode,
				uuid,
			)
		}
	}
	return bodyString, getTagResp.StatusCode, fmt.Errorf(
		"unexpected http status %d retrieving sharing info for record %s",
		getTagResp.StatusCode,
		uuid,
	)
}
