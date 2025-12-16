// Package dataset holds a model containing the relevant fields for dataset metadata.
package metadata

import (
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
	"github.com/pdok/pdok-metadata-tool/pkg/model/iso1911x"
)

// NLDatasetMetadata is used for retrieving the relevant fields from dataset metadata.
type NLDatasetMetadata struct {
	MetadataID     string
	SourceID       string
	Title          string
	Abstract       string
	ContactName    string
	ContactEmail   string
	ContactURL     string
	Keywords       []string
	LicenceURL     string
	UseLimitation  string
	ThumbnailURL   *string
	InspireVariant *inspire.InspireVariant
	InspireThemes  []string
	HVDCategories  []hvd.HVDCategory
	BoundingBox    *BoundingBox
}

// NewNLDatasetMetadataFromMDMetadata creates a new instance based on dataset metadata from a CSW response.
func NewNLDatasetMetadataFromMDMetadata(m *iso1911x.MDMetadata) *NLDatasetMetadata {
	return NewNLDatasetMetadataFromMDMetadataWithHVDRepo(m, nil)
}

// NewNLDatasetMetadataFromMDMetadataWithHVDRepo creates a new instance and enriches HVD categories
// using the provided HVD category provider.
func NewNLDatasetMetadataFromMDMetadataWithHVDRepo(
	m *iso1911x.MDMetadata,
	hvdRepo hvd.CategoryProvider,
) *NLDatasetMetadata {
	return &NLDatasetMetadata{
		MetadataID:     m.UUID,
		SourceID:       m.UUID,
		Title:          m.IdentificationInfo.MDDataIdentification.Title,
		Abstract:       m.IdentificationInfo.MDDataIdentification.Abstract,
		ContactName:    m.IdentificationInfo.MDDataIdentification.ContactName,
		ContactEmail:   m.IdentificationInfo.MDDataIdentification.ContactEmail,
		ContactURL:     m.IdentificationInfo.MDDataIdentification.ContactURL,
		Keywords:       m.GetKeywords(),
		LicenceURL:     m.GetLicenseURL(),
		UseLimitation:  m.IdentificationInfo.MDDataIdentification.UseLimitation,
		ThumbnailURL:   m.GetThumbnailURL(),
		InspireVariant: m.GetInspireVariant(),
		InspireThemes:  m.GetInspireThemes(),
		HVDCategories:  m.GetHVDCategories(hvdRepo),
		BoundingBox: &BoundingBox{
			WestBoundLongitude: m.IdentificationInfo.MDDataIdentification.Extent.WestBoundLongitude,
			EastBoundLongitude: m.IdentificationInfo.MDDataIdentification.Extent.EastBoundLongitude,
			SouthBoundLatitude: m.IdentificationInfo.MDDataIdentification.Extent.SouthBoundLatitude,
			NorthBoundLatitude: m.IdentificationInfo.MDDataIdentification.Extent.NorthBoundLatitude,
		},
	}
}
