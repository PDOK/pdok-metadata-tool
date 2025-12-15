package repository

import (
	"fmt"
	"net/url"

	"github.com/pdok/pdok-metadata-tool/pkg/client"
	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/metadata"
)

// MetadataRepository is used for looking up metadata using the given CSW endpoint.
type MetadataRepository struct {
	CswClient *client.CswClient
	HVDRepo   hvd.CategoryProvider
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
		datasetMetadata = metadata.NewNLDatasetMetadataFromMDMetadataWithHVDRepo(&mdMetadata, mr.HVDRepo)
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

// SetCache enables caching on the underlying CSW client.
func (mr *MetadataRepository) SetCache(cacheDir string, ttlHours int) {
	if mr != nil && mr.CswClient != nil {
		mr.CswClient.SetCache(cacheDir, ttlHours)
	}
}

// UnsetCache disables caching on the underlying CSW client.
func (mr *MetadataRepository) UnsetCache() {
	if mr != nil && mr.CswClient != nil {
		mr.CswClient.UnsetCache()
	}
}

// HarvestByCQLConstraint is a generic harvester that returns flat models based on the MetadataType in the constraint.
// Usage:
//   - For services: repo.HarvestByCQLConstraint[metadata.NLServiceMetadata](constraintWithTypeService)
//   - For datasets: repo.HarvestByCQLConstraint[metadata.NLDatasetMetadata](constraintWithTypeDataset)
func HarvestByCQLConstraint[T any](
	mr *MetadataRepository,
	constraint *csw.GetRecordsCQLConstraint,
) (result []T, err error) {
	if constraint == nil {
		return nil, fmt.Errorf("constraint must not be nil")
	}
	if constraint.MetadataType == nil {
		return nil, fmt.Errorf("constraint.MetadataType must be set to 'service' or 'dataset'")
	}

	// Determine expected T vs MetadataType
	var zero T
	switch any(zero).(type) {
	case metadata.NLServiceMetadata:
		if *constraint.MetadataType != csw.Service {
			return nil, fmt.Errorf("type parameter mismatch: T=NLServiceMetadata but MetadataType=%s", constraint.MetadataType.String())
		}
		mds, err := mr.CswClient.HarvestByCQLConstraint(constraint)
		if err != nil {
			return nil, err
		}
		for i := range mds {
			if mds[i].IdentificationInfo.SVServiceIdentification != nil {
				sm := metadata.NewNLServiceMetadataFromMDMetadataWithHVDRepo(&mds[i], mr.HVDRepo)
				if sm != nil {
					result = append(result, any(*sm).(T))
				}
			}
		}
		return result, nil
	case metadata.NLDatasetMetadata:
		if *constraint.MetadataType != csw.Dataset {
			return nil, fmt.Errorf("type parameter mismatch: T=NLDatasetMetadata but MetadataType=%s", constraint.MetadataType.String())
		}
		mds, err := mr.CswClient.HarvestByCQLConstraint(constraint)
		if err != nil {
			return nil, err
		}
		for i := range mds {
			if mds[i].IdentificationInfo.MDDataIdentification.Title != "" {
				dm := metadata.NewNLDatasetMetadataFromMDMetadataWithHVDRepo(&mds[i], mr.HVDRepo)
				if dm != nil {
					result = append(result, any(*dm).(T))
				}
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported type parameter T; must be NLServiceMetadata or NLDatasetMetadata")
	}
}

// SetHVDRepo sets the HVD category provider used to enrich HVD categories in flat models.
func (mr *MetadataRepository) SetHVDRepo(repo hvd.CategoryProvider) {
	mr.HVDRepo = repo
}
