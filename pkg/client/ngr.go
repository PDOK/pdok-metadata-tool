package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pdok/pdok-metadata-tool/pkg/model/ngr"
)

type NgrClient struct {
	host   *url.URL
	client *http.Client
}

func NewNgrClient(host *url.URL) NgrClient {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	return NgrClient{
		host:   host,
		client: client,
	}
}

// TODO Use this for harvesting only INSPIRE service metadata in ETF-validator-go
func (c NgrClient) GetRecordTags(uuid string) (ngr.RecordTagsResponse, error) {
	mdTagUrl := fmt.Sprintf("%s/geonetwork/srv/api/records/%s/tags", c.host.String(), uuid)

	recordTagsResponse := ngr.RecordTagsResponse{}

	err := getUnmarshalledJSONResponse(&recordTagsResponse, mdTagUrl, *c.client)
	if err != nil {
		return nil, err
	}
	return recordTagsResponse, nil
}
