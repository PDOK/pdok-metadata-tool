package metadata

import (
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
	"github.com/pdok/pdok-metadata-tool/pkg/model/iso1911x"
)

// NLServiceMetadata is used for retrieving the relevant fields from service metadata.
// It flattens ISO 19119 Service Identification to an easier-to-use internal model,
// similar in spirit to NLDatasetMetadata for datasets.
type NLServiceMetadata struct {
	MetadataID       string
	SourceID         string
	Title            string
	OrganisationName string
	Keywords         []string
	ServiceType      string
	OperatesOn       []OperatesOnRef
	Endpoints        []ServiceEndpoint
	ThumbnailURL     *string
	InspireVariant   *inspire.InspireVariant
	InspireThemes    []string
	HVDCategories    []hvd.HVDCategory
}

// OperatesOnRef represents a coupled dataset reference from a service metadata record.
type OperatesOnRef struct {
	UUIDRef string
	Href    string
}

// ServiceEndpoint represents an access endpoint for the service, including protocol information.
type ServiceEndpoint struct {
	URL          string
	Protocol     string
	ProtocolHref string
}

// NewNLServiceMetadataFromMDMetadata creates a new instance based on service metadata from a CSW response.
func NewNLServiceMetadataFromMDMetadata(m *iso1911x.MDMetadata) *NLServiceMetadata {
	sm := &NLServiceMetadata{
		MetadataID:     m.UUID,
		SourceID:       m.UUID,
		Title:          "",
		Keywords:       nil,
		ServiceType:    "",
		OperatesOn:     nil,
		Endpoints:      nil,
		ThumbnailURL:   nil,
		InspireVariant: nil,
		InspireThemes:  nil,
		HVDCategories:  nil,
	}

	if m.IdentificationInfo.SVServiceIdentification != nil {
		svc := m.IdentificationInfo.SVServiceIdentification
		sm.Title = svc.Title
		sm.ServiceType = svc.ServiceType

		// Organisation (point of contact organisation name)
		if svc.ResponsibleParty != nil {
			if svc.ResponsibleParty.Char != "" {
				sm.OrganisationName = svc.ResponsibleParty.Char
			} else if svc.ResponsibleParty.Anchor != "" {
				sm.OrganisationName = svc.ResponsibleParty.Anchor
			}
		}

		// Keywords
		sm.Keywords = m.GetServiceKeywords()

		// Coupled datasets
		for _, op := range svc.OperatesOn {
			sm.OperatesOn = append(sm.OperatesOn, OperatesOnRef{
				UUIDRef: op.Uuidref,
				Href:    op.Href,
			})
		}

		// Thumbnail
		sm.ThumbnailURL = m.GetServiceThumbnailURL()

		// INSPIRE & HVD (service)
		sm.InspireThemes = m.GetServiceInspireThemes()
		sm.HVDCategories = m.GetServiceHVDCategories()
		// Variant is determined from dataQualityInfo (applies to the record overall)
		sm.InspireVariant = m.GetInspireVariant()
	}

	// Distribution endpoints
	for _, ol := range m.OnLine {
		ep := ServiceEndpoint{URL: ol.URL}
		if ol.Protocol.Anchor.Text != "" {
			ep.Protocol = ol.Protocol.Anchor.Text
		}
		if ol.Protocol.Anchor.Href != "" {
			ep.ProtocolHref = ol.Protocol.Anchor.Href
		}
		sm.Endpoints = append(sm.Endpoints, ep)
	}

	return sm
}

func (m *NLServiceMetadata) GetInspireVariant() string {

	if m.InspireVariant != nil {
		return string(*m.InspireVariant)
	}
	return ""

}
