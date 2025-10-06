// Package dataset holds a model containing the relevant fields for dataset metadata.
package dataset

import (
	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
)

// NLDatasetMetadata is used for retrieving the relevant fields from dataset metadata.
type NLDatasetMetadata struct {
	MetadataID     string
	SourceID       string
	Title          string
	Abstract       string
	ContactName    string
	ContactEmail   string
	Keywords       []string
	LicenceURL     string
	ThumbnailURL   *string
	InspireVariant *inspire.InspireVariant
	InspireThemes  []string
	HVDCategories  []hvd.HVDCategory
}

// NewNLDatasetMetadataFromMDMetadata creates a new instance based on dataset metadata from a CSW response.
func NewNLDatasetMetadataFromMDMetadata(m *csw.MDMetadata) *NLDatasetMetadata {
	return &NLDatasetMetadata{
		MetadataID:     m.UUID,
		SourceID:       m.UUID,
		Title:          m.IdentificationInfo.MDDataIdentification.Title,
		Abstract:       m.IdentificationInfo.MDDataIdentification.Abstract,
		ContactName:    m.IdentificationInfo.MDDataIdentification.ContactName,
		ContactEmail:   m.IdentificationInfo.MDDataIdentification.ContactEmail,
		Keywords:       m.GetKeywords(),
		LicenceURL:     m.GetLicenseURL(),
		ThumbnailURL:   m.GetThumbnailURL(),
		InspireVariant: m.GetInspireVariant(),
		InspireThemes:  m.GetInspireThemes(),
		HVDCategories:  m.GetHVDCategories(),
	}
}
