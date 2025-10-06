package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pdok/pdok-metadata-tool/pkg/model/ngr"
)

// NgrClient is used as a client for doing NGR requests.
type NgrClient struct {
	host   *url.URL
	client *http.Client
}

// NewNgrClient creates a new instance of NgrClient.
func NewNgrClient(host *url.URL) NgrClient {
	const defaultTimeoutSeconds = 20

	client := &http.Client{
		Timeout: defaultTimeoutSeconds * time.Second,
	}

	return NgrClient{
		host:   host,
		client: client,
	}
}

// GetRecordTags returns the NGR tags for a give record id.
// TODO Use this for harvesting only INSPIRE service metadata in ETF-validator-go.
func (c NgrClient) GetRecordTags(uuid string) (ngr.RecordTagsResponse, error) {
	mdTagURL := fmt.Sprintf("%s/geonetwork/srv/api/records/%s/tags", c.host.String(), uuid)

	recordTagsResponse := ngr.RecordTagsResponse{}

	err := getUnmarshalledJSONResponse(&recordTagsResponse, mdTagURL, *c.client)
	if err != nil {
		return nil, err
	}

	return recordTagsResponse, nil
}
