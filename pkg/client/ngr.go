package client

import (
	"fmt"
	"github.com/pdok/pdok-metadata-tool/pkg/model/ngr"
	"net/http"
	"net/url"
	"time"
)

type NgrClient struct {
	host      *url.URL
	client    *http.Client
	cswClient CswClient
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
func (c NgrClient) GetRecordTags(uuid string, logPrefix string) (ngr.RecordTagsResponse, error) {
	mdTagUrl := fmt.Sprintf("%s/geonetwork/srv/api/records/%s/tags", c.host.String(), uuid)

	recordTagsResponse := ngr.RecordTagsResponse{}

	err := getUnmarshalledJSONResponse(&recordTagsResponse, mdTagUrl, *c.client, logPrefix)
	if err != nil {
		return nil, err
	}
	return recordTagsResponse, nil
}
