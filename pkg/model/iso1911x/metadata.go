// Package iso1911x holds models for ISO 19115/19119 metadata (read/write).
// This file defines a unified raw model for unmarshalling CSW MD_Metadata records
// for both datasets (ISO 19115) and services (ISO 19119). Helper methods to
// extract commonly used values are attached to this type.
package iso1911x

import (
	"strings"

	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
)

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
	ResponsibleParty        *CSWResponsibleParty `xml:"contact>CI_ResponsibleParty>OrganisationName"`
	IdentificationInfo      struct {
		SVServiceIdentification *struct {
			Title               string               `xml:"citation>CI_Citation>title>CharacterString"`
			ResponsibleParty    *CSWResponsibleParty `xml:"pointOfContact>CI_ResponsibleParty>OrganisationName"`
			GraphicOverview     *CSWGraphicOverview  `xml:"graphicOverview"`
			DescriptiveKeywords []struct {
				MDKeywords CSWMDKeywords `xml:"MD_Keywords"`
			} `xml:"descriptiveKeywords"`
			ServiceType string `xml:"serviceType>LocalName"`
			OperatesOn  []struct {
				Uuidref string `xml:"uuidref,attr"`
				Href    string `xml:"href,attr"`
			} `xml:"operatesOn"`
		} `xml:"SV_ServiceIdentification"`
		MDDataIdentification *struct {
			Title               string              `xml:"citation>CI_Citation>title>CharacterString"`
			Abstract            string              `xml:"abstract>CharacterString"`
			GraphicOverview     *CSWGraphicOverview `xml:"graphicOverview"`
			DescriptiveKeywords []struct {
				MDKeywords CSWMDKeywords `xml:"MD_Keywords"`
			} `xml:"descriptiveKeywords"`
			ContactName      string               `xml:"pointOfContact>CI_ResponsibleParty>individualName>CharacterString"`
			ContactEmail     string               `xml:"pointOfContact>CI_ResponsibleParty>contactInfo>CI_Contact>address>CI_Address>electronicMailAddress>CharacterString"`
			ContactURL       string               `xml:"pointOfContact>CI_ResponsibleParty>contactInfo>CI_Contact>onlineResource>CI_OnlineResource>linkage>URL"`
			LicenseURL       []CSWAnchor          `xml:"resourceConstraints>MD_LegalConstraints>otherConstraints>Anchor"`
			UseLimitation    string               `xml:"resourceConstraints>MD_Constraints>useLimitation>CharacterString"`
			ResponsibleParty *CSWResponsibleParty `xml:"pointOfContact>CI_ResponsibleParty>OrganisationName"`
			Extent           struct {
				WestBoundLongitude string `xml:"westBoundLongitude>Decimal"`
				EastBoundLongitude string `xml:"eastBoundLongitude>Decimal"`
				SouthBoundLatitude string `xml:"southBoundLatitude>Decimal"`
				NorthBoundLatitude string `xml:"northBoundLatitude>Decimal"`
			} `xml:"extent>EX_Extent>geographicElement>EX_GeographicBoundingBox"`
		} `xml:"MD_DataIdentification"`
	} `xml:"identificationInfo"`
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

// ResponsibleParty struct for unmarshalling the CSW response.
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

// GetDatasetKeywords returns a slice of keywords (datasets).
func (m *MDMetadata) GetDatasetKeywords() (keywords []string) {
	if m.IdentificationInfo.MDDataIdentification == nil {
		return
	}
	for _, descriptiveKeyword := range m.IdentificationInfo.MDDataIdentification.DescriptiveKeywords {
		for _, keyword := range descriptiveKeyword.MDKeywords.Keyword {
			if keyword.CharacterString != "" {
				keywords = append(keywords, keyword.CharacterString)
			} else if keyword.Anchor.Text != "" {
				keywords = append(keywords, keyword.Anchor.Text)
			}
		}
	}
	return
}

// GetServiceKeywords returns a slice of keywords from service metadata.
func (m *MDMetadata) GetServiceKeywords() (keywords []string) {
	if m.IdentificationInfo.SVServiceIdentification == nil {
		return
	}
	for _, descriptiveKeyword := range m.IdentificationInfo.SVServiceIdentification.DescriptiveKeywords {
		for _, keyword := range descriptiveKeyword.MDKeywords.Keyword {
			if keyword.CharacterString != "" {
				keywords = append(keywords, keyword.CharacterString)
			} else if keyword.Anchor.Text != "" {
				keywords = append(keywords, keyword.Anchor.Text)
			}
		}
	}
	return
}

// GetDatasetLicenseURL returns the license href (dataset).
func (m *MDMetadata) GetDatasetLicenseURL() string {
	if m.IdentificationInfo.MDDataIdentification == nil {
		return ""
	}
	for _, val := range m.IdentificationInfo.MDDataIdentification.LicenseURL {
		if strings.Contains(val.Href, "creativecommons.org") {
			return val.Href
		}
	}
	return ""
}

// GetDatasetThumbnailURL returns the thumbnail URL (dataset).
func (m *MDMetadata) GetDatasetThumbnailURL() *string {
	if m.IdentificationInfo.MDDataIdentification != nil &&
		m.IdentificationInfo.MDDataIdentification.GraphicOverview != nil {
		thumbnailURL := m.IdentificationInfo.MDDataIdentification.GraphicOverview.MDBrowseGraphic.FileName
		return &thumbnailURL
	}
	return nil
}

// GetServiceThumbnailURL returns the thumbnail URL from a service metadata record.
func (m *MDMetadata) GetServiceThumbnailURL() *string {
	if m.IdentificationInfo.SVServiceIdentification != nil &&
		m.IdentificationInfo.SVServiceIdentification.GraphicOverview != nil {
		thumbnailURL := m.IdentificationInfo.SVServiceIdentification.GraphicOverview.MDBrowseGraphic.FileName
		if thumbnailURL != "" {
			return &thumbnailURL
		}
	}
	return nil
}

// GetInspireVariant retrieves the INSPIRE variant from dataset metadata.
func (m *MDMetadata) GetInspireVariant() *inspire.InspireVariant {
	isInspire := false
	isConformant := true
	inspireRegulation := "VERORDENING (EU) Nr. 1089/2010"

	harmonised := inspire.Harmonised
	asIs := inspire.AsIs

	for _, report := range m.DQDataQuality.Report {
		for _, result := range report.ConsistencyResult {
			specificationTitle := ""
			if result.Specification.CharacterString != "" {
				specificationTitle = result.Specification.CharacterString
			} else if result.Specification.Anchor.Text != "" {
				specificationTitle = result.Specification.Anchor.Text
			}

			hasInspireRegulation := strings.Contains(specificationTitle, inspireRegulation)
			hasInspireSpecTitle := strings.HasPrefix(specificationTitle, "INSPIRE")

			if hasInspireRegulation {
				isInspire = true
			}

			if hasInspireSpecTitle && result.Pass != "true" {
				isConformant = false
			}
		}
	}

	switch {
	case isInspire && isConformant:
		return &harmonised
	case isInspire:
		return &asIs
	default:
		return nil
	}
}

// GetDatasetInspireThemes retrieves the INSPIRE themes from dataset metadata.
func (m *MDMetadata) GetDatasetInspireThemes() (themes []string) {
	const (
		thesaurusName            = "GEMET - INSPIRE themes, version 1.0"
		thesaurusVocabularyDutch = "http://www.eionet.europa.eu/gemet/nl/inspire-theme/"
	)

	if m.IdentificationInfo.MDDataIdentification == nil {
		return
	}

	for _, descriptiveKeyword := range m.IdentificationInfo.MDDataIdentification.DescriptiveKeywords {
		thesaurus := descriptiveKeyword.MDKeywords.Thesaurus

		if thesaurus.CharacterString != thesaurusName && thesaurus.Anchor.Text != thesaurusName {
			// Skip, this is not the right thesaurus
			continue
		}

		for _, keyword := range descriptiveKeyword.MDKeywords.Keyword {
			if keyword.Anchor.Href != "" {
				// Try to get the INSPIRE theme from the anchor
				// according to TG Recommendation 1.5: metadata/2.0/rec/datasets-and-series/use-anchors-for-gemet
				if strings.Contains(keyword.Anchor.Href, thesaurusVocabularyDutch) {
					theme := strings.ReplaceAll(
						keyword.Anchor.Href,
						thesaurusVocabularyDutch,
						"",
					)
					themes = append(themes, theme)
				}
			} else if keyword.CharacterString != "" {
				// Otherwise match the keyword with values of the GEMET vocabulary
				// according to TG Requirement 1.4: metadata/2.0/req/datasets-and-series/inspire-theme-keyword
				theme := inspire.GetInspireThemeIDForDutchLabel(keyword.CharacterString)
				if theme != "" {
					themes = append(themes, theme)
				}
			}
		}
	}

	return themes
}

// GetDatasetHVDCategories retrieves the HVD categories from dataset metadata.
func (m *MDMetadata) GetDatasetHVDCategories() (categories []hvd.HVDCategory) {
	const thesaurusVocabulary = "http://data.europa.eu/bna/"

	if m.IdentificationInfo.MDDataIdentification == nil {
		return
	}

	for _, descriptiveKeyword := range m.IdentificationInfo.MDDataIdentification.DescriptiveKeywords {
		for _, keyword := range descriptiveKeyword.MDKeywords.Keyword {
			if keyword.Anchor.Href != "" {
				if strings.Contains(keyword.Anchor.Href, thesaurusVocabulary) {
					parts := strings.Split(keyword.Anchor.Href, "/")
					category := parts[len(parts)-1]
					categories = append(categories, hvd.HVDCategory{ID: category})
				}
			}
		}
	}

	return
}

// GetServiceInspireThemes retrieves the INSPIRE themes from service metadata.
func (m *MDMetadata) GetServiceInspireThemes() (themes []string) {
	const (
		thesaurusName            = "GEMET - INSPIRE themes, version 1.0"
		thesaurusVocabularyDutch = "http://www.eionet.europa.eu/gemet/nl/inspire-theme/"
	)

	if m.IdentificationInfo.SVServiceIdentification == nil {
		return
	}

	for _, descriptiveKeyword := range m.IdentificationInfo.SVServiceIdentification.DescriptiveKeywords {
		thesaurus := descriptiveKeyword.MDKeywords.Thesaurus

		if thesaurus.CharacterString != thesaurusName && thesaurus.Anchor.Text != thesaurusName {
			// Skip, this is not the right thesaurus
			continue
		}

		for _, keyword := range descriptiveKeyword.MDKeywords.Keyword {
			if keyword.Anchor.Href != "" {
				if strings.Contains(keyword.Anchor.Href, thesaurusVocabularyDutch) {
					theme := strings.ReplaceAll(keyword.Anchor.Href, thesaurusVocabularyDutch, "")
					themes = append(themes, theme)
				}
			} else if keyword.CharacterString != "" {
				theme := inspire.GetInspireThemeIDForDutchLabel(keyword.CharacterString)
				if theme != "" {
					themes = append(themes, theme)
				}
			}
		}
	}

	return themes
}

// GetServiceHVDCategories retrieves the HVD categories from service metadata.
func (m *MDMetadata) GetServiceHVDCategories() (categories []hvd.HVDCategory) {
	const thesaurusVocabulary = "http://data.europa.eu/bna/"

	if m.IdentificationInfo.SVServiceIdentification == nil {
		return
	}

	for _, descriptiveKeyword := range m.IdentificationInfo.SVServiceIdentification.DescriptiveKeywords {
		for _, keyword := range descriptiveKeyword.MDKeywords.Keyword {
			if keyword.Anchor.Href != "" {
				if strings.Contains(keyword.Anchor.Href, thesaurusVocabulary) {
					parts := strings.Split(keyword.Anchor.Href, "/")
					category := parts[len(parts)-1]
					categories = append(categories, hvd.HVDCategory{ID: category})
				}
			}
		}
	}

	return
}
