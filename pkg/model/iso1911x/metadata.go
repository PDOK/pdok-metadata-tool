// Package iso1911x holds models for ISO 19115/19119 metadata (read/write).
// This file defines a unified raw model for unmarshalling CSW MD_Metadata records
// for both datasets (ISO 19115) and services (ISO 19119). Helper methods to
// extract commonly used values are attached to this type.
package iso1911x

import (
	"html"
	"net/url"
	"strings"

	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
)

// MetadataType holds the possible types of metadata for MD_Metadata records.
type MetadataType string

// Possible values for MetadataType.
const (
	Service MetadataType = "service"
	Dataset MetadataType = "dataset"
)

// String returns the string representation of the MetadataType.
func (m MetadataType) String() string { return string(m) }

// MDMetadata struct for unmarshalling the CSW response (MD_Metadata).
// Supports both dataset (MD_DataIdentification) and service (SV_ServiceIdentification).
type MDMetadata struct {
	SelfURL string
	MdType  *struct {
		CodeListValue string `xml:",chardata"`
		TextValue     string `xml:"codeListValue,attr"`
	} `xml:"hierarchyLevel>MD_ScopeCode"`
	MetadataStandardVersion string               `xml:"metadataStandardVersion>CharacterString"`
	UUID                    string               `xml:"fileIdentifier>CharacterString"`
	ResponsibleParty        *CSWResponsibleParty `xml:"contact>CI_ResponsibleParty>organisationName"`
	IdentificationInfo      struct {
		SVServiceIdentification *struct {
			Title               string                  `xml:"citation>CI_Citation>title>CharacterString"`
			Abstract            string                  `xml:"abstract>CharacterString"`
			ResponsibleParty    *CSWResponsibleParty    `xml:"pointOfContact>CI_ResponsibleParty>organisationName"`
			GraphicOverview     *CSWGraphicOverview     `xml:"graphicOverview"`
			DescriptiveKeywords []CSWDescriptiveKeyword `xml:"descriptiveKeywords"`
			ServiceType         string                  `xml:"serviceType>LocalName"`
			LicenseURL          []CSWAnchor             `xml:"resourceConstraints>MD_LegalConstraints>otherConstraints>Anchor"`
			OperatesOn          []struct {
				Uuidref string `xml:"uuidref,attr"`
				Href    string `xml:"href,attr"`
			} `xml:"operatesOn"`
		} `xml:"SV_ServiceIdentification"`
		MDDataIdentification *struct {
			Title               string                  `xml:"citation>CI_Citation>title>CharacterString"`
			Source              Source                  `xml:"citation>CI_Citation>identifier>MD_Identifier>code"`
			Abstract            string                  `xml:"abstract>CharacterString"`
			GraphicOverview     *CSWGraphicOverview     `xml:"graphicOverview"`
			DescriptiveKeywords []CSWDescriptiveKeyword `xml:"descriptiveKeywords"`
			ContactName         string                  `xml:"pointOfContact>CI_ResponsibleParty>individualName>CharacterString"`
			ContactEmail        string                  `xml:"pointOfContact>CI_ResponsibleParty>contactInfo>CI_Contact>address>CI_Address>electronicMailAddress>CharacterString"`
			ContactURL          string                  `xml:"pointOfContact>CI_ResponsibleParty>contactInfo>CI_Contact>onlineResource>CI_OnlineResource>linkage>URL"`
			LicenseURL          []CSWAnchor             `xml:"resourceConstraints>MD_LegalConstraints>otherConstraints>Anchor"`
			UseLimitation       string                  `xml:"resourceConstraints>MD_Constraints>useLimitation>CharacterString"`
			ResponsibleParty    *CSWResponsibleParty    `xml:"pointOfContact>CI_ResponsibleParty>OrganisationName"`
			Extent              struct {
				WestBoundLongitude string `xml:"westBoundLongitude>Decimal"`
				EastBoundLongitude string `xml:"eastBoundLongitude>Decimal"`
				SouthBoundLatitude string `xml:"southBoundLatitude>Decimal"`
				NorthBoundLatitude string `xml:"northBoundLatitude>Decimal"`
			} `xml:"extent>EX_Extent>geographicElement>EX_GeographicBoundingBox"`
		} `xml:"MD_DataIdentification"`
	} `xml:"identificationInfo"`
	// todo: also implement transferOptions>MD_DigitalTransferOptions>onLine for datasets
	OnLine []struct {
		URL      string `xml:"CI_OnlineResource>linkage>URL"`
		Protocol struct {
			Anchor CSWAnchor `xml:"Anchor"`
		} `xml:"CI_OnlineResource>protocol"`
	} `xml:"distributionInfo>MD_Distribution>transferOptions>MD_DigitalTransferOptions>onLine"`
	DQDataQuality struct {
		Report []struct {
			ConsistencyResult []struct {
				Specification struct {
					CharacterString string    `xml:"CharacterString"`
					Anchor          CSWAnchor `xml:"Anchor"`
				} `xml:"DQ_ConformanceResult>specification>CI_Citation>title"`
				Explanation string `xml:"DQ_ConformanceResult>explanation>CharacterString"`
				Pass        string `xml:"DQ_ConformanceResult>pass>Boolean"`
			} `xml:"DQ_DomainConsistency>result"`
		} `xml:"report"`
	} `xml:"dataQualityInfo>DQ_DataQuality"`
}

// CSWResponsibleParty represents a text or Anchor value for organisation names in CSW records.
type CSWResponsibleParty struct {
	Char   string `xml:"CharacterString"`
	Anchor string `xml:"Anchor"`
}

// CSWAnchor struct for unmarshalling text + href anchors in CSW responses.
type CSWAnchor struct {
	Text string `xml:",chardata"`
	Href string `xml:"href,attr"`
}

// CSWKeywordEntry represents either a CharacterString or Anchor keyword value.
type CSWKeywordEntry struct {
	CharacterString string    `xml:"CharacterString"`
	Anchor          CSWAnchor `xml:"Anchor"`
}

// CSWKeywordType wraps MD_KeywordTypeCode attributes.
type CSWKeywordType struct {
	MDKeywordTypeCode struct {
		CodeList      string `xml:"codeList,attr"`
		CodeListValue string `xml:"codeListValue,attr"`
	} `xml:"MD_KeywordTypeCode"`
}

// CSWThesaurus describes the thesaurus name used for keywords.
type CSWThesaurus struct {
	CharacterString string    `xml:"CharacterString"`
	Anchor          CSWAnchor `xml:"Anchor"`
}

// CSWGraphicOverview is a shared type for dataset/service graphic overviews in CSW records.
type CSWGraphicOverview struct {
	MDBrowseGraphic CSWBrowseGraphic `xml:"MD_BrowseGraphic"`
}

// CSWBrowseGraphic represents gmd:MD_BrowseGraphic for CSW reader model.
type CSWBrowseGraphic struct {
	FileName        string `xml:"fileName>CharacterString"`
	FileDescription string `xml:"fileDescription>CharacterString"`
	FileType        string `xml:"fileType>CharacterString"`
}

// CSWMDKeywords models MD_Keywords for both dataset and service records.
type CSWMDKeywords struct {
	Keyword   []CSWKeywordEntry `xml:"keyword"`
	Type      CSWKeywordType    `xml:"type"`
	Thesaurus CSWThesaurus      `xml:"thesaurusName>CI_Citation>title"`
}

// CSWDescriptiveKeyword wraps MD_Keywords blocks for reuse across dataset and service structures.
type CSWDescriptiveKeyword struct {
	MDKeywords CSWMDKeywords `xml:"MD_Keywords"`
}

// Source models the source identification, which can be stored in an anchor or character string
type Source struct {
	CharacterString string    `xml:"CharacterString"`
	Anchor          CSWAnchor `xml:"Anchor"`
}

func (s *Source) GetID() string {
	if s.Anchor.Text != "" {
		return s.Anchor.Text
	}

	return s.CharacterString
}

// ServiceEndpoint represents an access endpoint for the service, including protocol information.
type ServiceEndpoint struct {
	URL      string
	Protocol string
}

// NormalizeXMLText removes leading and trailing whitespace from XML text nodes.
func NormalizeXMLText(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// GetMetaDataType returns whether this MDMetadata represents a dataset or a service.
// It uses the hierarchyLevel>MD_ScopeCode value (codeListValue or text) when present.
// Defaults to Dataset when unknown.
func (m *MDMetadata) GetMetaDataType() MetadataType {
	if m.MdType != nil {
		// Prefer CodeListValue (element content), fallback to attribute if needed
		val := m.MdType.CodeListValue
		if val == "" {
			val = m.MdType.TextValue
		}

		switch strings.ToLower(strings.TrimSpace(val)) {
		case string(Service):
			return Service
		case string(Dataset):
			return Dataset
		}
	}
	// Default to dataset when unknown
	return Dataset
}

// GetKeywords returns non-INSPIRE, non-HVD keywords for both dataset and service metadata.
func (m *MDMetadata) GetKeywords() (keywords []string) {
	var dks []CSWDescriptiveKeyword

	switch m.GetMetaDataType() {
	case Service:
		if m.IdentificationInfo.SVServiceIdentification != nil {
			dks = m.IdentificationInfo.SVServiceIdentification.DescriptiveKeywords
		}
	case Dataset:
		if m.IdentificationInfo.MDDataIdentification != nil {
			dks = m.IdentificationInfo.MDDataIdentification.DescriptiveKeywords
		}
	}

	const (
		inspireThesaurusName    = "GEMET - INSPIRE themes, version 1.0"
		inspireDutchVocabulary  = "http://www.eionet.europa.eu/gemet/nl/inspire-theme/"
		hvdConceptVocabulary    = "http://data.europa.eu/bna/"
		hvdThesaurusTitleAnchor = "http://publications.europa.eu/resource/dataset/high-value-dataset-category"
	)

	for _, dk := range dks {
		th := dk.MDKeywords.Thesaurus
		// Skip INSPIRE keywords groups
		if NormalizeXMLText(th.CharacterString) == inspireThesaurusName ||
			NormalizeXMLText(th.Anchor.Text) == inspireThesaurusName ||
			(th.Anchor.Href != "" && strings.Contains(th.Anchor.Href, inspireDutchVocabulary)) {
			continue
		}
		// Skip HVD keywords groups
		isHVDGroup := false
		if th.Anchor.Href != "" && strings.Contains(th.Anchor.Href, hvdThesaurusTitleAnchor) {
			isHVDGroup = true
		} else {
			for _, kw := range dk.MDKeywords.Keyword {
				if kw.Anchor.Href != "" && strings.Contains(kw.Anchor.Href, hvdConceptVocabulary) {
					isHVDGroup = true

					break
				}
			}
		}

		if isHVDGroup {
			continue
		}
		// Collect generic keywords
		for _, kw := range dk.MDKeywords.Keyword {
			if kw.CharacterString != "" {
				keywords = append(keywords, NormalizeXMLText(kw.CharacterString))
			} else if kw.Anchor.Text != "" {
				keywords = append(keywords, NormalizeXMLText(kw.Anchor.Text))
			}
		}
	}

	return keywords
}

// GetLicenseURL returns a license URL for either dataset or service (if present).
func (m *MDMetadata) GetLicenseURL() string {
	switch m.GetMetaDataType() {
	case Dataset:
		if m.IdentificationInfo.MDDataIdentification == nil {
			return ""
		}

		for _, val := range m.IdentificationInfo.MDDataIdentification.LicenseURL {
			if strings.Contains(val.Href, "creativecommons.org") {
				return NormalizeXMLText(val.Href)
			}
		}
	case Service:
		if m.IdentificationInfo.SVServiceIdentification == nil {
			return ""
		}

		for _, val := range m.IdentificationInfo.SVServiceIdentification.LicenseURL {
			if strings.Contains(val.Href, "creativecommons.org") {
				return NormalizeXMLText(val.Href)
			}
		}
	}

	return ""
}

// GetThumbnailURL returns the thumbnail URL from either dataset or service metadata.
func (m *MDMetadata) GetThumbnailURL() string {
	switch m.GetMetaDataType() {
	case Service:
		if m.IdentificationInfo.SVServiceIdentification == nil ||
			m.IdentificationInfo.SVServiceIdentification.GraphicOverview == nil {
			return ""
		}

		thumbnailURL := m.IdentificationInfo.SVServiceIdentification.GraphicOverview.MDBrowseGraphic.FileName
		if thumbnailURL != "" {
			return NormalizeXMLText(thumbnailURL)
		}
	case Dataset:
		if m.IdentificationInfo.MDDataIdentification == nil ||
			m.IdentificationInfo.MDDataIdentification.GraphicOverview == nil {
			return ""
		}

		thumbnailURL := m.IdentificationInfo.MDDataIdentification.GraphicOverview.MDBrowseGraphic.FileName
		if thumbnailURL != "" {
			return NormalizeXMLText(thumbnailURL)
		}
	}

	return ""
}

// GetInspireVariantForDataset retrieves the INSPIRE variant from dataset metadata.
func (m *MDMetadata) GetInspireVariantForDataset() inspire.InspireVariant {
	isInspire := false
	isConformant := true
	inspireRegulation := "VERORDENING (EU) Nr. 1089/2010"

	harmonised := inspire.Harmonised
	asIs := inspire.AsIs

	foundInspireRegulation := false

	for _, report := range m.DQDataQuality.Report {
		for _, result := range report.ConsistencyResult {
			specificationTitle := ""
			if result.Specification.CharacterString != "" {
				specificationTitle = result.Specification.CharacterString
			} else if result.Specification.Anchor.Text != "" {
				specificationTitle = result.Specification.Anchor.Text
			}

			foundInspireRegulation = strings.Contains(
				NormalizeXMLText(specificationTitle),
				inspireRegulation,
			)

			if foundInspireRegulation {
				isInspire = true
			}

			if result.Pass != "true" {
				isConformant = false
			}

			if foundInspireRegulation {
				break
			}
		}

		if foundInspireRegulation {
			break
		}
	}

	switch {
	case isInspire && isConformant:
		return harmonised
	case isInspire:
		return asIs
	default:
		return ""
	}
}

// GetInspireThemes is a generic version that returns INSPIRE themes for dataset or service.
func (m *MDMetadata) GetInspireThemes() (themes []string) {
	const (
		thesaurusName              = "GEMET - INSPIRE themes, version 1.0"
		thesaurusVocabularyDutch   = "http://www.eionet.europa.eu/gemet/nl/inspire-theme/"
		thesaurusVocabularyEnglish = "http://www.eionet.europa.eu/gemet/en/inspire-theme/"
		inspireThemeRegistry       = "http://inspire.ec.europa.eu/theme/"
	)

	var dks []CSWDescriptiveKeyword

	switch m.GetMetaDataType() {
	case Service:
		if m.IdentificationInfo.SVServiceIdentification != nil {
			dks = m.IdentificationInfo.SVServiceIdentification.DescriptiveKeywords
		}
	case Dataset:
		if m.IdentificationInfo.MDDataIdentification != nil {
			dks = m.IdentificationInfo.MDDataIdentification.DescriptiveKeywords
		}
	}

	for _, descriptiveKeyword := range dks {
		thesaurus := descriptiveKeyword.MDKeywords.Thesaurus

		if NormalizeXMLText(thesaurus.CharacterString) != thesaurusName &&
			NormalizeXMLText(thesaurus.Anchor.Text) != thesaurusName {
			// Skip, this is not the right thesaurus
			continue
		}

		for _, keyword := range descriptiveKeyword.MDKeywords.Keyword {
			if keyword.Anchor.Href != "" {
				// Try to get the INSPIRE theme from the anchor according to TG Recommendation 1.5
				expectedPrefixes := []string{
					thesaurusVocabularyDutch,
					thesaurusVocabularyEnglish,
					inspireThemeRegistry,
				}

				for _, prefix := range expectedPrefixes {
					if strings.Contains(keyword.Anchor.Href, prefix) {
						theme := strings.ReplaceAll(keyword.Anchor.Href, prefix, "")
						themes = append(themes, theme)

						break
					}
				}
			} else if keyword.CharacterString != "" {
				// Otherwise match the keyword with values of the GEMET vocabulary (TG Requirement 1.4)
				theme := inspire.GetInspireThemeIDForDutchLabel(keyword.CharacterString)
				if theme != "" {
					themes = append(themes, theme)
				}
			}
		}
	}

	return themes
}

// GetHVDCategories is the generic method to retrieve HVD categories from metadata.
// It detects whether the MDMetadata is a dataset or service, extracts HVD category codes
// from the appropriate descriptive keywords, and (optionally) enriches them via the
// provided HVDRepository to include full details.
func (m *MDMetadata) GetHVDCategories(
	hvdRepo hvd.CategoryProvider,
) (categories []hvd.HVDCategory) { //nolint:lll
	const thesaurusVocabulary = "http://data.europa.eu/bna/"

	// Use a map to avoid duplicates while preserving order via slice append checks
	seen := map[string]bool{}

	var dks []CSWDescriptiveKeyword

	switch m.GetMetaDataType() {
	case Service:
		if m.IdentificationInfo.SVServiceIdentification != nil {
			dks = m.IdentificationInfo.SVServiceIdentification.DescriptiveKeywords
		}
	case Dataset:
		if m.IdentificationInfo.MDDataIdentification != nil {
			dks = m.IdentificationInfo.MDDataIdentification.DescriptiveKeywords
		}
	}

	for _, descriptiveKeyword := range dks {
		for _, keyword := range descriptiveKeyword.MDKeywords.Keyword {
			if keyword.Anchor.Href == "" {
				continue
			}

			if !strings.Contains(keyword.Anchor.Href, thesaurusVocabulary) {
				continue
			}

			parts := strings.Split(keyword.Anchor.Href, "/")
			code := strings.TrimSpace(parts[len(parts)-1])
			label := strings.TrimSpace(keyword.Anchor.Text)

			if seen[code] {
				continue
			}

			seen[code] = true
			if hvdRepo != nil {
				if cat, err := hvdRepo.GetHVDCategoryByCode(code); err == nil && cat != nil {
					categories = append(categories, *cat)

					continue
				}
			}

			categories = append(categories, hvd.HVDCategory{ID: code, LabelDutch: label})
		}
	}

	return categories
}

func (m *MDMetadata) GetOperatesOnForService() (result []string) {
	for _, val := range m.IdentificationInfo.SVServiceIdentification.OperatesOn {
		unescapedHref := html.UnescapeString(val.Href)

		hrefUrl, err := url.Parse(unescapedHref)
		if err == nil {
			for _, key := range []string{"id", "ID"} {
				id := hrefUrl.Query().Get(key)
				if id != "" {
					// remove whitespace
					id = strings.ReplaceAll(id, " ", "")

					result = append(result, id)
				}
			}
		} else {
			result = append(result, strings.ReplaceAll(val.Uuidref, " ", ""))
		}
	}

	return
}

func (m *MDMetadata) GetServiceEndpointsForService() (result []ServiceEndpoint) {
	for _, ol := range m.OnLine {
		ep := ServiceEndpoint{URL: NormalizeXMLText(ol.URL)}
		if ol.Protocol.Anchor.Text != "" {
			ep.Protocol = NormalizeXMLText(ol.Protocol.Anchor.Text)
		}

		result = append(result, ep)
	}

	return
}

func (m *MDMetadata) GetServiceContactForService() string {
	if m.IdentificationInfo.SVServiceIdentification.ResponsibleParty != nil {
		if m.IdentificationInfo.SVServiceIdentification.ResponsibleParty.Char != "" {
			return strings.TrimSpace(
				m.IdentificationInfo.SVServiceIdentification.ResponsibleParty.Char,
			)
		} else if m.IdentificationInfo.SVServiceIdentification.ResponsibleParty.Anchor != "" {
			return strings.TrimSpace(m.IdentificationInfo.SVServiceIdentification.ResponsibleParty.Anchor)
		}
	}

	return ""
}
