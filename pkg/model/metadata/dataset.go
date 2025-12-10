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
	return &NLDatasetMetadata{
		MetadataID:     m.UUID,
		SourceID:       m.UUID,
		Title:          m.IdentificationInfo.MDDataIdentification.Title,
		Abstract:       m.IdentificationInfo.MDDataIdentification.Abstract,
		ContactName:    m.IdentificationInfo.MDDataIdentification.ContactName,
		ContactEmail:   m.IdentificationInfo.MDDataIdentification.ContactEmail,
		ContactURL:     m.IdentificationInfo.MDDataIdentification.ContactURL,
		Keywords:       m.GetDatasetKeywords(),
		LicenceURL:     m.GetDatasetLicenseURL(),
		UseLimitation:  m.IdentificationInfo.MDDataIdentification.UseLimitation,
		ThumbnailURL:   m.GetDatasetThumbnailURL(),
		InspireVariant: m.GetInspireVariant(),
		InspireThemes:  m.GetDatasetInspireThemes(),
		HVDCategories:  m.GetDatasetHVDCategories(),
		BoundingBox: &BoundingBox{
			WestBoundLongitude: m.IdentificationInfo.MDDataIdentification.Extent.WestBoundLongitude,
			EastBoundLongitude: m.IdentificationInfo.MDDataIdentification.Extent.EastBoundLongitude,
			SouthBoundLatitude: m.IdentificationInfo.MDDataIdentification.Extent.SouthBoundLatitude,
			NorthBoundLatitude: m.IdentificationInfo.MDDataIdentification.Extent.NorthBoundLatitude,
		},
	}
}
