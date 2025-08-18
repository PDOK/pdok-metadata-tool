package repository

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"pdok-metadata-tool/pkg/client"
	"pdok-metadata-tool/pkg/model/csw"
	"pdok-metadata-tool/pkg/model/dataset"
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

func (mr *MetadataRepository) GetDatasetMetadata(logPrefix string, limit int) (metadataRecords []dataset.NLDatasetMetadata, err error) {
	metadataRecords, err = mr.HarvestMDMetadata(csw.Dataset, limit, logPrefix)
	return
}

func (mr *MetadataRepository) GetDatasetMetadataById(id string, logPrefix string) (datasetMetadata *dataset.NLDatasetMetadata, err error) {
	mdMetadata, err := mr.CswClient.GetRecordById(id, logPrefix)
	if err != nil {
		return
	}

	datasetMetadata = dataset.NewNLDatasetMetadataFromMDMetadata(&mdMetadata)
	return
}

func (mr *MetadataRepository) HarvestMDMetadata(metadataType csw.MetadataType, limit int, logPrefix string) (metadataRecords []dataset.NLDatasetMetadata, err error) {
	summaryRecords, err := mr.HarvestSummaryRecords(metadataType, limit, logPrefix)
	if err != nil {
		return nil, err
	}
	for _, summaryRecord := range summaryRecords {
		if metadataRecord, err := mr.GetDatasetMetadataById(summaryRecord.Identifier, logPrefix); err != nil {
			log.Warningf("Failed to get dataset metadata for id %s: %v", summaryRecord.Identifier, err)
		} else {
			metadataRecords = append(metadataRecords, *metadataRecord)
		}
	}
	return
}

func (mr *MetadataRepository) HarvestSummaryRecords(metadataType csw.MetadataType, limit int, logPrefix string) (summaryRecords []csw.SummaryRecord, err error) {
	var processed = 0
	var offset = 1
harvestSummaryRecordsLoop:
	for {
		constraint := csw.GetRecordsConstraint{MetadataType: &metadataType}
		records, nextRecord, err := mr.CswClient.GetRecords(&constraint, offset, logPrefix)
		if err != nil {
			log.Fatalf("CSW Harvest could not determine if more records are available, %v", err)
		}
		for _, summaryRecord := range records {
			summaryRecords = append(summaryRecords, summaryRecord)
			processed++

			if limit > 0 && processed >= limit {
				break harvestSummaryRecordsLoop
			}

		}
		if nextRecord == 0 {
			break harvestSummaryRecordsLoop
		}
		offset = nextRecord
	}
	return
}
