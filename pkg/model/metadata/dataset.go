// Package dataset holds a model containing the relevant fields for dataset metadata.
package metadata

import (
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/inspire"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/iso1911x"
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
	LicenseText    string
	UseLimitation  string
	ThumbnailURL   string
	InspireVariant inspire.InspireVariant
	InspireThemes  []string
	HVDCategories  []hvd.HVDCategory
	BoundingBox    *BoundingBox
	CreationDate   string

	// Temp fields for validation script
	ContactOrganisationName string
	DigitalTransferOptions  []OnLine
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
		MetadataID: iso1911x.NormalizeXMLText(m.UUID),
		SourceID: iso1911x.NormalizeXMLText(
			m.IdentificationInfo.MDDataIdentification.Source.GetID(),
		),
		Title: iso1911x.NormalizeXMLText(m.IdentificationInfo.MDDataIdentification.Title),
		Abstract: iso1911x.NormalizeXMLText(
			m.IdentificationInfo.MDDataIdentification.Abstract,
		),
		ContactName: iso1911x.NormalizeXMLText(
			m.IdentificationInfo.MDDataIdentification.ContactName,
		),
		ContactEmail: iso1911x.NormalizeXMLText(
			m.IdentificationInfo.MDDataIdentification.ContactEmail,
		),
		ContactURL: iso1911x.NormalizeXMLText(
			m.IdentificationInfo.MDDataIdentification.ContactURL,
		),
		Keywords:    m.GetKeywords(),
		LicenceURL:  m.GetLicenseURL(),
		LicenseText: m.GetLicenseText(),
		UseLimitation: iso1911x.NormalizeXMLText(
			m.GetUseLimitation(),
		),
		ThumbnailURL:   m.GetThumbnailURL(),
		InspireVariant: m.GetInspireVariantForDataset(),
		InspireThemes:  m.GetInspireThemes(),
		HVDCategories:  m.GetHVDCategories(hvdRepo),
		CreationDate:   m.GetCreationDate(),
		BoundingBox: &BoundingBox{
			WestBoundLongitude: iso1911x.NormalizeXMLText(
				m.IdentificationInfo.MDDataIdentification.Extent.WestBoundLongitude,
			),
			EastBoundLongitude: iso1911x.NormalizeXMLText(
				m.IdentificationInfo.MDDataIdentification.Extent.EastBoundLongitude,
			),
			SouthBoundLatitude: iso1911x.NormalizeXMLText(
				m.IdentificationInfo.MDDataIdentification.Extent.SouthBoundLatitude,
			),
			NorthBoundLatitude: iso1911x.NormalizeXMLText(
				m.IdentificationInfo.MDDataIdentification.Extent.NorthBoundLatitude,
			),
		},
		// Fill temp fields for validation script
		ContactOrganisationName: iso1911x.NormalizeXMLText(
			m.GetContactOrganisationName(),
		),
		DigitalTransferOptions: getDigitalTransferOptions(m),
	}
}

func getDigitalTransferOptions(m *iso1911x.MDMetadata) []OnLine {
	var ret []OnLine

	for _, do := range m.OnLine {
		online := OnLine{
			URL: iso1911x.NormalizeXMLText(do.URL),
		}
		if do.Protocol.Anchor.Text != "" {
			online.Protocol = iso1911x.NormalizeXMLText(do.Protocol.Anchor.Text)
		} else {
			online.Protocol = iso1911x.NormalizeXMLText(do.Protocol.CharacterString)
		}

		ret = append(ret, online)
	}

	return ret
}
