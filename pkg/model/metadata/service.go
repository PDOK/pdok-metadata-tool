package metadata

import (
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/iso1911x"
)

// NLServiceMetadata is used for retrieving the relevant fields from service metadata.
// It flattens ISO 19119 Service Identification to an easier-to-use internal model,
// similar in spirit to NLDatasetMetadata for datasets.
type NLServiceMetadata struct {
	MetadataID       string
	Title            string
	Abstract         string
	OrganisationName string
	Keywords         []string
	ServiceType      string
	OperatesOn       []string
	Endpoints        []iso1911x.ServiceEndpoint
	ThumbnailURL     string
	LicenceURL       string
	UseLimitation    string
	//	InspireVariant *inspire.InspireVariant // A service does not have an inspire variant. Datasets have. The service can be conform inspire or not. But when a dataset is as-is the service is still 100% conform inspire. Thus, the conformity of a service is separate from the dataset it serves. This is our current interpretation. In the future we might call this field conformInspire. But this only reflects if the service is conform inspire and not if the dataset in the service is conform.
	InspireThemes []string
	HVDCategories []hvd.HVDCategory
	CreationDate  string
	RevisionDate  string
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
		MetadataID: iso1911x.NormalizeXMLText(m.UUID),
		Title: iso1911x.NormalizeXMLText(
			m.IdentificationInfo.SVServiceIdentification.Title,
		),
		Abstract: iso1911x.NormalizeXMLText(
			m.IdentificationInfo.SVServiceIdentification.Abstract,
		),
		OrganisationName: m.GetServiceContactForService(),
		Keywords:         m.GetKeywords(),
		ServiceType: iso1911x.NormalizeXMLText(
			m.IdentificationInfo.SVServiceIdentification.ServiceType,
		),
		OperatesOn:    m.GetOperatesOnForService(),
		Endpoints:     m.GetServiceEndpointsForService(),
		ThumbnailURL:  m.GetThumbnailURL(),
		LicenceURL:    m.GetLicenseURL(),
		UseLimitation: m.GetUseLimitation(),
		InspireThemes: m.GetInspireThemes(),
		HVDCategories: m.GetHVDCategories(hvdRepo),
		CreationDate:  m.GetCreationDate(),
		RevisionDate:  m.GetRevisionDate(),
	}

	// Organisation (point of contact organisation name)

	return sm
}
