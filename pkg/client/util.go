package client

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pdok/pdok-metadata-tool/v2/internal/common"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeXML  = "application/xml"
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
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := string(bodyBytes)

		fmt.Println(bodyStr)

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
	contentType string,
) ([]byte, error) {
	var req *http.Request
	if requestBody != nil {
		req, _ = http.NewRequest(method, url, bytes.NewBufferString(*requestBody))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	xsrfToken, err := obtainXSRFToken(ngrConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain XSRF token: %w", err)
	}

	setXsrfToken(req, xsrfToken)

	username := ngrConfig.NgrUserName
	password := ngrConfig.NgrPassword
	req.SetBasicAuth(*username, *password)

	req.Header.Set("User-Agent", "pdok.nl (pdok-metadata-tool)")
	req.Header.Set("Accept", "*/*;q=0.8,application/signed-exchange")
	req.Header.Set("Content-Type", contentType)

	//nolint:bodyclose // We use common.SafeClose to handle closing the response body
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while calling NGR using url %s: %w", url, err)
	}
	defer common.SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf(
			"error while calling NGR using url %s\nhttp status is %d",
			url,
			resp.StatusCode,
		)
	}

	return io.ReadAll(resp.Body)
}

func setXsrfToken(req *http.Request, xsrfToken string) {
	req.Header.Set("X-Xsrf-Token", xsrfToken)
	req.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: xsrfToken})
}

func getCookieValueByName(cookies []string, name string) (string, error) {
	for _, c := range cookies {
		if strings.Contains(c, name+"=") {
			parts := strings.Split(c, ";")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, name+"=") {
					two := 2

					split := strings.SplitN(part, "=", two)
					if len(split) == two {
						return split[1], nil
					}
				}
			}
		}
	}

	return "", errors.New("cannot obtain " + name + " from cookie")
}
