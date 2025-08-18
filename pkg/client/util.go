package client

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func getUnmarshalledXMLResponse(resultStruct interface{}, url string, client http.Client, logPrefix string) error {
	resp, err := getResponse(url, client, logPrefix)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Read body: %v", err)
		return fmt.Errorf("Error while reading NGR response from url %s ; %v.", url, err)
	}

	err = xml.Unmarshal(data, resultStruct)
	if err != nil {
		log.Errorf("Unmarshalling: %v", err)
		return fmt.Errorf("Error unmarshalling NGR response from url %s ; %v.", url, err)
	}
	return nil
}

func getUnmarshalledJSONResponse(resultStruct interface{}, url string, client http.Client, logPrefix string) error {
	resp, err := getResponse(url, client, logPrefix)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Read body: %v", err)
		return fmt.Errorf("Error while reading NGR response from url %s ; %v.", url, err)
	}

	err = json.Unmarshal(data, resultStruct)
	if err != nil {
		log.Errorf("Unmarshalling: %v", err)
		return fmt.Errorf("Error unmarshalling NGR response from url %s ; %v.", url, err)
	}
	return nil
}

func getResponse(url string, client http.Client, logPrefix string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "pdok.nl (inspire-etf-validator)")
	req.Header.Set("Accept", "*/*;q=0.8,application/signed-exchange")
	log.Infof("%s: get metadata using url %s", logPrefix, req.URL)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Retrieving service fron NGR: %v", err)
		return nil, fmt.Errorf("Error while calling NGR using url %s ; %v.", url, err)
	}
	if resp.StatusCode != 200 {
		log.Errorf("Retrieving services from NGR: HTTP Statuscode %d", resp.StatusCode)
		return nil, fmt.Errorf("Error while calling NGR using url %s ; http status is %d.", url, resp.StatusCode)
	}
	return resp, nil
}
