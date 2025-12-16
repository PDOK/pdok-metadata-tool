package metadata

import (
	"html"
	"net/url"
	"strings"

	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
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
	LicenceURL       string
	//	InspireVariant *inspire.InspireVariant // A service does not have an inspire variant. Datasets have. The service can be conform inspire or not. But when a dataset is as-is the service is still 100% conform inspire. Thus, the conformity of a service is separate from the dataset it serves. This is our current interpretation. In the future we might call this field conformInspire. But this only reflects if the service is conform inspire and not if the dataset in the service is conform.
	InspireThemes []string
	HVDCategories []hvd.HVDCategory
}

// OperatesOnRef represents a coupled dataset reference from a service metadata record.
type OperatesOnRef struct {
	UUIDRef string
	Href    string
}

func (o *OperatesOnRef) GetID() string {
	// The UUIDRef field has been deprecated, see https://geonovum.github.io/Metadata-ISO19119/#gekoppelde-bron
	// It can still be present in NGR metadata, but it does not always match the id in the CSW href.
	// Therefor we first try to parse the id from the CSW href.
	unescapedHref := html.UnescapeString(o.Href)

	hrefUrl, err := url.Parse(unescapedHref)
	if err == nil {
		for _, key := range []string{"id", "ID"} {
			id := hrefUrl.Query().Get(key)
			if id != "" {
				// remove whitespace
				id = strings.ReplaceAll(id, " ", "")

				return id
			}
		}
	}

	return strings.ReplaceAll(o.UUIDRef, " ", "")
}

// ServiceEndpoint represents an access endpoint for the service, including protocol information.
type ServiceEndpoint struct {
	URL          string
	Protocol     string
	ProtocolHref string
}

// NewNLServiceMetadataFromMDMetadata creates a new instance based on service metadata from a CSW response.
func NewNLServiceMetadataFromMDMetadata(m *iso1911x.MDMetadata) *NLServiceMetadata {
	return NewNLServiceMetadataFromMDMetadataWithHVDRepo(m, nil)
}

// NewNLServiceMetadataFromMDMetadataWithHVDRepo creates a new instance and enriches HVD categories
// using the provided HVD category provider.
func NewNLServiceMetadataFromMDMetadataWithHVDRepo(
	m *iso1911x.MDMetadata,
	hvdRepo hvd.CategoryProvider,
) *NLServiceMetadata {
	sm := &NLServiceMetadata{
		MetadataID:   m.UUID,
		SourceID:     m.UUID,
		Title:        "",
		Keywords:     nil,
		ServiceType:  "",
		OperatesOn:   nil,
		Endpoints:    nil,
		ThumbnailURL: nil,
		LicenceURL:   "",
		//InspireVariant: nil,
		InspireThemes: nil,
		HVDCategories: nil,
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

		// Keywords (generic, excludes INSPIRE and HVD groups)
		sm.Keywords = m.GetKeywords()

		// Coupled datasets
		for _, op := range svc.OperatesOn {
			sm.OperatesOn = append(sm.OperatesOn, OperatesOnRef{
				UUIDRef: op.Uuidref,
				Href:    op.Href,
			})
		}

		// Thumbnail (generic)
		sm.ThumbnailURL = m.GetThumbnailURL()

		// License URL (from MD_LegalConstraints)
		sm.LicenceURL = m.GetLicenseURL()

		// INSPIRE & HVD (service)
		sm.InspireThemes = m.GetInspireThemes()
		sm.HVDCategories = m.GetHVDCategories(hvdRepo)
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

// func (m *NLServiceMetadata) GetInspireVariant() string {
//
//	if m.InspireVariant != nil {
//		return string(*m.InspireVariant)
//	}
//	return ""
//
//}
