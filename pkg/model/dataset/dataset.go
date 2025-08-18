package dataset

import (
	"github.com/google/uuid"
	"pdok-metadata-tool/pkg/model/csw"
	"pdok-metadata-tool/pkg/model/hvd"
	"pdok-metadata-tool/pkg/model/inspire"
)

type NLDatasetMetadata struct {
	MetadataId     uuid.UUID
	SourceId       string
	Title          string
	Abstract       string
	ContactName    string
	ContactEmail   string
	Keywords       []string
	LicenceUrl     string
	ThumbnailUrl   *string
	InspireVariant *inspire.InspireVariant
	InspireThemes  []string
	HVDCategories  []hvd.HVDCategory
}

func NewNLDatasetMetadataFromMDMetadata(m *csw.MDMetadata) *NLDatasetMetadata {
	// TODO What to do with invalid uuids?
	metadataId, err := uuid.Parse(m.UUID)
	if err != nil {
		metadataId = uuid.Nil
	}

	return &NLDatasetMetadata{
		MetadataId:     metadataId,
		SourceId:       m.UUID,
		Title:          m.IdentificationInfo.MDDataIdentification.Title,
		Abstract:       m.IdentificationInfo.MDDataIdentification.Abstract,
		ContactName:    m.IdentificationInfo.MDDataIdentification.ContactName,
		ContactEmail:   m.IdentificationInfo.MDDataIdentification.ContactEmail,
		Keywords:       m.GetKeywords(),
		LicenceUrl:     m.GetLicenseUrl(),
		ThumbnailUrl:   m.GetThumbnailUrl(),
		InspireVariant: m.GetInspireVariant(),
		InspireThemes:  m.GetInspireThemes(),
		HVDCategories:  m.GetHVDCategories(),
	}
}
