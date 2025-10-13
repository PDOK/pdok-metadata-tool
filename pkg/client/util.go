package client

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pdok/pdok-metadata-tool/internal/common"
)

func getUnmarshalledXMLResponse(
	resultStruct any,
	url string,
	method string,
	requestBody *string,
	client http.Client,
) error {
	responseBody, err := getResponseBody(url, method, requestBody, client)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(responseBody, resultStruct)
	if err != nil {
		return fmt.Errorf("error unmarshalling NGR response from url %s: %w", url, err)
	}

	return nil
}

func getUnmarshalledJSONResponse(resultStruct any, url string, client http.Client) error {
	responseBody, err := getResponseBody(url, "GET", nil, client)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, resultStruct)
	if err != nil {
		return fmt.Errorf("error unmarshalling NGR response from url %s: %w", url, err)
	}

	return nil
}

func getResponseBody(
	url string,
	method string,
	requestBody *string,
	client http.Client,
) ([]byte, error) {
	var body io.Reader
	if method == "POST" && requestBody != nil {
		body = strings.NewReader(*requestBody)
	}

	req, _ := http.NewRequest(method, url, body)

	req.Header.Set("User-Agent", "pdok.nl (pdok-metadata-tool)")
	req.Header.Set("Accept", "*/*;q=0.8,application/signed-exchange")
	req.Header.Set("Content-Type", "application/xml")

	//nolint:bodyclose // We use common.SafeClose to handle closing the response body
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while calling NGR using url %s: %w", url, err)
	}
	defer common.SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"error while calling NGR using url %s\nhttp status is %d",
			url,
			resp.StatusCode,
		)
	}

	return io.ReadAll(resp.Body)
}

func getNgrResponseBody(
	ngrConfig *NgrConfig,
	url string,
	method string,
	requestBody *string,
	client http.Client,
) ([]byte, error) {
	var body io.Reader
	if method == "POST" && requestBody != nil {
		body = strings.NewReader(*requestBody)
	}

	req, _ := http.NewRequest(method, url, body)
	xsrfToken, err := obtainXSRFToken(ngrConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain XSRF token: %w", err)
	}
	req.Header.Set("Cookie", "XSRF-TOKEN="+xsrfToken)
	username := ngrConfig.NgrUserName
	password := ngrConfig.NgrPassword
	//auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.SetBasicAuth(username, password) //Header.Set("Authorization", "Basic "+auth)

	req.Header.Set("User-Agent", "pdok.nl (pdok-metadata-tool)")
	req.Header.Set("Accept", "*/*;q=0.8,application/signed-exchange")
	req.Header.Set("Content-Type", "application/xml")

	//nolint:bodyclose // We use common.SafeClose to handle closing the response body
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while calling NGR using url %s: %w", url, err)
	}
	defer common.SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"error while calling NGR using url %s\nhttp status is %d",
			url,
			resp.StatusCode,
		)
	}

	return io.ReadAll(resp.Body)
}
