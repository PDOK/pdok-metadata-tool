package repository

import (
	"github.com/pdok/pdok-metadata-tool/pkg/client"
	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
	"github.com/pdok/pdok-metadata-tool/pkg/model/dataset"
	"net/url"
)

type MetadataRepository struct {
	CswClient *client.CswClient
}

func NewMetadataRepository(cswHost string, cswPath string) (*MetadataRepository, error) {
	h, err := url.Parse(cswHost)
	if err != nil {
		return nil, err
	}
	h.Path = cswPath

	cswClient := client.NewCswClient(h)

	return &MetadataRepository{
		CswClient: &cswClient,
	}, nil
}

func (mr *MetadataRepository) GetDatasetMetadataById(id string, logPrefix string) (datasetMetadata *dataset.NLDatasetMetadata, err error) {
	mdMetadata, err := mr.CswClient.GetRecordById(id, logPrefix)
	if err != nil {
		return
	}

	datasetMetadata = dataset.NewNLDatasetMetadataFromMDMetadata(&mdMetadata)
	return
}

func (mr *MetadataRepository) SearchDatasetMetadata(title *string, id *string, logPrefix string) (records []csw.SummaryRecord, err error) {
	filter := csw.GetRecordsOgcFilter{
		MetadataType: csw.Dataset,
		Title:        title,
		Identifier:   id,
	}
	records, err = mr.CswClient.GetRecordsWithOGCFilter(&filter, logPrefix)
	return
}
