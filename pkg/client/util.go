package client

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getUnmarshalledXMLResponse(resultStruct interface{}, url string, method string, requestBody *string, client http.Client) error {
	resp, err := getResponse(url, method, requestBody, client)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error while reading NGR response from url: %s\n\n%v", url, err)
	}

	err = xml.Unmarshal(data, resultStruct)
	if err != nil {
		return fmt.Errorf("Error unmarshalling NGR response from url: %s\n\n%v", url, err)
	}
	return nil
}

func getUnmarshalledJSONResponse(resultStruct interface{}, url string, client http.Client) error {
	resp, err := getResponse(url, "GET", nil, client)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error while reading NGR response from url: %s\n\n%v", url, err)
	}

	err = json.Unmarshal(data, resultStruct)
	if err != nil {
		return fmt.Errorf("Error unmarshalling NGR response from url: %s\n\n%v", url, err)
	}
	return nil
}

func getResponse(url string, method string, requestBody *string, client http.Client) (*http.Response, error) {
	var body io.Reader = nil
	if method == "POST" && requestBody != nil {
		body = strings.NewReader(*requestBody)
	}
	req, _ := http.NewRequest(method, url, body)

	req.Header.Set("User-Agent", "pdok.nl (pdok-metadata-tool)")
	req.Header.Set("Accept", "*/*;q=0.8,application/signed-exchange")
	// log.Infof("get metadata using url %s", req.URL) // Used for debugging
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while calling NGR using url: %s\n\n%v", url, err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error while calling NGR using url %s\nhttp status is %d", url, resp.StatusCode)
	}
	return resp, nil
}
