package metadata

import "github.com/pdok/pdok-metadata-tool/v2/pkg/model/iso1911x"

type BoundingBox struct {
	WestBoundLongitude string
	EastBoundLongitude string
	SouthBoundLatitude string
	NorthBoundLatitude string
}

type Thumbnail struct {
	URL         string
	Description string
	Type        string
}

func getBoundingBox(m *iso1911x.MDMetadata) *BoundingBox {
	if m == nil {
		return nil
	}

	return &BoundingBox{
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
	}
}

// getThumbnail returns the thumbnail from either dataset or service metadata.
func getThumbnail(m *iso1911x.MDMetadata) *Thumbnail {
	if m == nil {
		return nil
	}

	switch m.GetMetaDataType() {
	case iso1911x.Service:
		if m.IdentificationInfo.SVServiceIdentification == nil ||
			m.IdentificationInfo.SVServiceIdentification.GraphicOverview == nil {
			return nil
		}

		browseGraphic := m.IdentificationInfo.SVServiceIdentification.GraphicOverview.MDBrowseGraphic

		return &Thumbnail{
			URL:         iso1911x.NormalizeXMLText(browseGraphic.FileName),
			Description: browseGraphic.FileDescription,
			Type:        browseGraphic.FileType,
		}
	case iso1911x.Dataset:
		if m.IdentificationInfo.MDDataIdentification == nil ||
			m.IdentificationInfo.MDDataIdentification.GraphicOverview == nil {
			return nil
		}

		browseGraphic := m.IdentificationInfo.MDDataIdentification.GraphicOverview.MDBrowseGraphic

		return &Thumbnail{
			URL:         iso1911x.NormalizeXMLText(browseGraphic.FileName),
			Description: browseGraphic.FileDescription,
			Type:        browseGraphic.FileType,
		}
	}

	return nil
}
