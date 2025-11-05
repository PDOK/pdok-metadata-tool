// Package csw holds models for handling CSW requests.
package csw

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
)

// MetadataType holds the possible types of metadata.
type MetadataType string

// Possible values for MetadataType.
const (
	Service MetadataType = "service"
	Dataset MetadataType = "dataset"
)

func (m MetadataType) String() string {
	return string(m)
}

// GetRecordByIDResponse struct for unmarshalling a CSW GetRecordByID response.
type GetRecordByIDResponse struct {
	XMLName    xml.Name   `xml:"GetRecordByIdResponse"`
	Text       string     `xml:",chardata"`
	Csw        string     `xml:"csw,attr"`
	MDMetadata MDMetadata `xml:"MD_Metadata"`
}

// GetRecordsResponse struct for unmarshalling a CSW GetRecords response.
type GetRecordsResponse struct {
	XMLName       xml.Name `xml:"GetRecordsResponse"`
	SearchResults struct {
		NumberOfRecordsMatched string          `xml:"numberOfRecordsMatched,attr"`
		NextRecord             string          `xml:"nextRecord,attr"`
		SummaryRecords         []SummaryRecord `xml:"SummaryRecord"`
	} `xml:"SearchResults"`
}

// GetRecordsCQLConstraint struct for creating a CQL constraint.
type GetRecordsCQLConstraint struct {
	MetadataType     *MetadataType
	OrganisationName *string
}

// ToQueryParameter returns a query parameter string based on the constraints.
func (c *GetRecordsCQLConstraint) ToQueryParameter() (constraint string) {
	var constraints []string

	if c.MetadataType != nil {
		constraints = append(constraints, fmt.Sprintf("type='%s'", *c.MetadataType))
	}

	if c.OrganisationName != nil {
		constraints = append(constraints, fmt.Sprintf("OrganisationName='%s'", *c.OrganisationName))
	}

	if len(constraints) == 0 {
		return
	}

	constraint += "&constraintLanguage=CQL_TEXT"
	constraint += "&constraint_language_version=1.1.0"
	constraint += "&constraint="

	for i, c := range constraints {
		constraint += c
		if i < len(constraints)-1 {
			constraint += "+AND+"
		}
	}

	return
}

// GetRecordsOgcFilter struct for creating an OgcFilter for a CSW GetRecords request.
type GetRecordsOgcFilter struct {
	MetadataType MetadataType
	Title        *string
	Identifier   *string
}

// ToRequestBody Returns a request body string for a CSW GetRecords request.
func (f *GetRecordsOgcFilter) ToRequestBody() (ogcFilter string, err error) {
	template := `<csw:GetRecords xmlns:csw="http://www.opengis.net/cat/csw/2.0.2" xmlns:ogc="http://www.opengis.net/ogc" service="CSW" version="2.0.2" resultType="results" startPosition="1" maxRecords="5" outputFormat="application/xml" outputSchema="http://www.opengis.net/cat/csw/2.0.2" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.opengis.net/cat/csw/2.0.2 http://schemas.opengis.net/csw/2.0.2/CSW-discovery.xsd">
	<csw:Query typeNames="csw:Record">
		<csw:ElementSetName>summary</csw:ElementSetName>
                <csw:Constraint version="1.1.0">
                <ogc:Filter>
                    %s
                </ogc:Filter>
            </csw:Constraint>
	</csw:Query>
</csw:GetRecords>`

	metadataTypeFilter := f.getPropertyIsEqualToClause("dc:type", f.MetadataType.String())

	var filters []string
	if f.Title != nil && len(*f.Title) > 0 {
		filters = append(filters, f.getPropertyIsLikeClause("dc:title", *f.Title))
	}

	if f.Identifier != nil && len(*f.Identifier) > 0 {
		filters = append(filters, f.getPropertyIsEqualToClause("dc:identifier", *f.Identifier))
	}

	var filter string

	switch len(filters) {
	case 0:
		filter = metadataTypeFilter
	case 1:
		filter = wrapFiltersInAndOperator([]string{metadataTypeFilter, filters[0]})
	default:
		filter = wrapFiltersInAndOperator(
			[]string{metadataTypeFilter, wrapFiltersInOrOperator(filters)},
		)
	}

	requestBody := fmt.Sprintf(template, filter)

	return requestBody, nil
}

func wrapFiltersInOrOperator(filters []string) string {
	template := `<ogc:Or>
    %s
</ogc:Or>`

	return fmt.Sprintf(template, strings.Join(filters, "\n"))
}

func wrapFiltersInAndOperator(filters []string) string {
	template := `<ogc:And>
    %s
</ogc:And>`

	return fmt.Sprintf(template, strings.Join(filters, "\n"))
}

func (f *GetRecordsOgcFilter) getPropertyIsLikeClause(
	property string,
	value string,
) (clause string) {
	template := `<ogc:PropertyIsLike escapeChar="" wildCard="%%" singleChar="_">
    <ogc:PropertyName>%s</ogc:PropertyName>
    <ogc:Literal>%%%s%%</ogc:Literal>
</ogc:PropertyIsLike>`

	return fmt.Sprintf(template, property, value)
}

func (f *GetRecordsOgcFilter) getPropertyIsEqualToClause(
	property string,
	value string,
) (clause string) {
	template := `<ogc:PropertyIsEqualTo>
    <ogc:PropertyName>%s</ogc:PropertyName>
    <ogc:Literal>%s</ogc:Literal>
</ogc:PropertyIsEqualTo>`

	return fmt.Sprintf(template, property, value)
}

// MDMetadata struct for unmarshalling the CSW response.
type MDMetadata struct {
	SelfURL string
	MdType  *struct {
		CodeListValue string `xml:",chardata"`
		TextValue     string `xml:"codeListValue,attr"`
	} `xml:"hierarchyLevel>MD_ScopeCode"`
	MetadataStandardVersion string            `xml:"metadataStandardVersion>CharacterString"`
	UUID                    string            `xml:"fileIdentifier>CharacterString"`
	ResponsibleParty        *ResponsibleParty `xml:"contact>CI_ResponsibleParty>OrganisationName"`
	IdentificationInfo      struct {
		SVServiceIdentification *struct {
			Title            string            `xml:"citation>CI_Citation>title>CharacterString"`
			ResponsibleParty *ResponsibleParty `xml:"pointOfContact>CI_ResponsibleParty>OrganisationName"`
			GraphicOverview  *struct {
				MDBrowseGraphic struct {
					FileName        string `xml:"fileName>>CharacterString"`
					FileDescription string `xml:"fileDescription>CharacterString"`
				} `xml:"MD_BrowseGraphic"`
			} `xml:"graphicOverview"`
			DescriptiveKeywords []struct {
				MDKeywords struct {
					Keyword []struct {
						CharacterString string `xml:"CharacterString"`
						Anchor          Anchor `xml:"Anchor"`
					} `xml:"keyword"`
					Type struct {
						MDKeywordTypeCode struct {
							CodeList      string `xml:"codeList,attr"`
							CodeListValue string `xml:"codeListValue,attr"`
						} `xml:"MD_KeywordTypeCode"`
					} `xml:"type"`
				} `xml:"MD_Keywords"`
			} `xml:"descriptiveKeywords"`
			ServiceType string `xml:"serviceType>LocalName"`
			OperatesOn  []struct {
				Uuidref string `xml:"uuidref,attr"`
				Href    string `xml:"href,attr"`
			} `xml:"operatesOn"`
		} `xml:"SV_ServiceIdentification"`
		MDDataIdentification *struct {
			Title           string `xml:"citation>CI_Citation>title>CharacterString"`
			Abstract        string `xml:"abstract>CharacterString"`
			GraphicOverview *struct {
				MDBrowseGraphic struct {
					FileName        string `xml:"fileName>CharacterString"`
					FileDescription string `xml:"fileDescription>CharacterString"`
				} `xml:"MD_BrowseGraphic"`
			} `xml:"graphicOverview"`
			DescriptiveKeywords []struct {
				MDKeywords struct {
					Keyword []struct {
						CharacterString string `xml:"CharacterString"`
						Anchor          Anchor `xml:"Anchor"`
					} `xml:"keyword"`
					Type struct {
						MDKeywordTypeCode struct {
							CodeList      string `xml:"codeList,attr"`
							CodeListValue string `xml:"codeListValue,attr"`
						} `xml:"MD_KeywordTypeCode"`
					} `xml:"type"`
					Thesaurus struct {
						CharacterString string `xml:"CharacterString"`
						Anchor          Anchor `xml:"Anchor"`
					} `xml:"thesaurusName>CI_Citation>title"`
				} `xml:"MD_Keywords"`
			} `xml:"descriptiveKeywords"`
			ContactName      string            `xml:"pointOfContact>CI_ResponsibleParty>individualName>CharacterString"`
			ContactEmail     string            `xml:"pointOfContact>CI_ResponsibleParty>contactInfo>CI_Contact>address>CI_Address>electronicMailAddress>CharacterString"`
			LicenseURL       []Anchor          `xml:"resourceConstraints>MD_LegalConstraints>otherConstraints>Anchor"`
			UseLimitation    string            `xml:"resourceConstraints>MD_Constraints>useLimitation>CharacterString"`
			ResponsibleParty *ResponsibleParty `xml:"pointOfContact>CI_ResponsibleParty>OrganisationName"`
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
			Anchor Anchor `xml:"Anchor"`
		} `xml:"CI_OnlineResource>protocol"`
	} `xml:"distributionInfo>MD_Distribution>transferOptions>MD_DigitalTransferOptions>onLine"`
	DQDataQuality struct {
		Report []struct {
			ConsistencyResult []struct {
				Specification struct {
					CharacterString string `xml:"CharacterString"`
					Anchor          Anchor `xml:"Anchor"`
				} `xml:"DQ_ConformanceResult>specification>CI_Citation>title"`
				Explanation string `xml:"DQ_ConformanceResult>explanation>CharacterString"`
				Pass        string `xml:"DQ_ConformanceResult>pass>Boolean"`
			} `xml:"DQ_DomainConsistency>result"`
		} `xml:"report"`
	} `xml:"dataQualityInfo>DQ_DataQuality"`
}

// ResponsibleParty struct for unmarshalling the CSW response.
type ResponsibleParty struct {
	Char   string `xml:"CharacterString"`
	Anchor string `xml:"Anchor"`
}

// Anchor struct for unmarshalling the CSW response.
type Anchor struct {
	Text string `xml:",chardata"`
	Href string `xml:"href,attr"`
}

// GetKeywords returns a slice of keywords.
func (m *MDMetadata) GetKeywords() (keywords []string) {
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

// GetLicenseURL returns the license href.
func (m *MDMetadata) GetLicenseURL() string {
	for _, val := range m.IdentificationInfo.MDDataIdentification.LicenseURL {
		if strings.Contains(val.Href, "creativecommons.org") {
			return val.Href
		}
	}

	return ""
}

// GetThumbnailURL returns the thumbnail URL.
func (m *MDMetadata) GetThumbnailURL() *string {
	if m.IdentificationInfo.MDDataIdentification != nil &&
		m.IdentificationInfo.MDDataIdentification.GraphicOverview != nil {
		thumbnailURL := m.IdentificationInfo.MDDataIdentification.GraphicOverview.MDBrowseGraphic.FileName

		return &thumbnailURL
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

// GetInspireThemes retrieves the INSPIRE themes from dataset metadata.
func (m *MDMetadata) GetInspireThemes() (themes []string) {
	const (
		thesaurusName            = "GEMET - INSPIRE themes, version 1.0"
		thesaurusVocabularyDutch = "http://www.eionet.europa.eu/gemet/nl/inspire-theme/"
	)

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

// GetHVDCategories retrieves the HVD categories from dataset metadata.
func (m *MDMetadata) GetHVDCategories() (categories []hvd.HVDCategory) {
	const thesaurusVocabulary = "http://data.europa.eu/bna/"

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

// SummaryRecord struct for unmarshalling the CSW response.
type SummaryRecord struct {
	Identifier string `json:"identifier" xml:"identifier"`
	Title      string `json:"title"      xml:"title"`
	Type       string `json:"type"       xml:"type"`
}
