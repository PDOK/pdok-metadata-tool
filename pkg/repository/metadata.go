package repository

import (
	"net/url"

	"github.com/pdok/pdok-metadata-tool/pkg/client"
	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
	"github.com/pdok/pdok-metadata-tool/pkg/model/metadata"
)

// MetadataRepository is used for looking up metadata using the given CSW endpoint.
type MetadataRepository struct {
	CswClient *client.CswClient
}

// NewMetadataRepository creates a new instance of MetadataRepository.
func NewMetadataRepository(cswEndpoint string) (*MetadataRepository, error) {
	h, err := url.Parse(cswEndpoint)
	if err != nil {
		return nil, err
	}

	cswClient := client.NewCswClient(h)

	return &MetadataRepository{
		CswClient: &cswClient,
	}, nil
}

// GetDatasetMetadataByID retrieves dataset metadata by id.
func (mr *MetadataRepository) GetDatasetMetadataByID(
	id string,
) (datasetMetadata *metadata.NLDatasetMetadata, err error) {
	mdMetadata, err := mr.CswClient.GetRecordByID(id)
	if err != nil {
		return
	}

	if mdMetadata.IdentificationInfo.MDDataIdentification != nil {
		datasetMetadata = metadata.NewNLDatasetMetadataFromMDMetadata(&mdMetadata)
	}

	return
}

// SearchDatasetMetadata searches for dataset metadata by title or id.
func (mr *MetadataRepository) SearchDatasetMetadata(
	title *string,
	id *string,
) (records []csw.SummaryRecord, err error) {
	filter := csw.GetRecordsOgcFilter{
		MetadataType: csw.Dataset,
		Title:        title,
		Identifier:   id,
	}
	records, err = mr.CswClient.GetRecordsWithOGCFilter(&filter)

	return
}
