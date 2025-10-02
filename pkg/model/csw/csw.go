package csw

import (
	"encoding/xml"
	"fmt"
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
	"strings"
)

type MetadataType string

const (
	Service MetadataType = "service"
	Dataset MetadataType = "dataset"
)

func (m MetadataType) String() string {
	return string(m)
}

type GetRecordByIdResponse struct {
	XMLName    xml.Name   `xml:"GetRecordByIdResponse"`
	Text       string     `xml:",chardata"`
	Csw        string     `xml:"csw,attr"`
	MDMetadata MDMetadata `xml:"MD_Metadata"`
}

type GetRecordsResponse struct {
	XMLName       xml.Name `xml:"GetRecordsResponse"`
	nextRecord    int      ``
	SearchResults struct {
		NumberOfRecordsMatched string          `xml:"numberOfRecordsMatched,attr"`
		NextRecord             string          `xml:"nextRecord,attr"`
		SummaryRecords         []SummaryRecord `xml:"SummaryRecord"`
	} `xml:"SearchResults"`
}

type GetRecordsCQLConstraint struct {
	MetadataType     *MetadataType
	OrganisationName *string
}

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

type GetRecordsOgcFilter struct {
	MetadataType MetadataType
	Title        *string
	Identifier   *string
}

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
		filter = wrapFiltersInAndOperator([]string{metadataTypeFilter, wrapFiltersInOrOperator(filters)})
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

func (f *GetRecordsOgcFilter) getPropertyIsLikeClause(property string, value string) (clause string) {
	template := `<ogc:PropertyIsLike escapeChar="" wildCard="%%" singleChar="_">
    <ogc:PropertyName>%s</ogc:PropertyName>
    <ogc:Literal>%%%s%%</ogc:Literal>
</ogc:PropertyIsLike>`

	return fmt.Sprintf(template, property, value)
}

func (f *GetRecordsOgcFilter) getPropertyIsEqualToClause(property string, value string) (clause string) {
	template := `<ogc:PropertyIsEqualTo>
    <ogc:PropertyName>%s</ogc:PropertyName>
    <ogc:Literal>%s</ogc:Literal>
</ogc:PropertyIsEqualTo>`

	return fmt.Sprintf(template, property, value)
}

type MDMetadata struct {
	SelfUrl string
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
					ThesaurusName Anchor `xml:"thesaurusName>CI_Citation>title>Anchor"`
				} `xml:"MD_Keywords"`
			} `xml:"descriptiveKeywords"`
			ContactName  string   `xml:"pointOfContact>CI_ResponsibleParty>individualName>CharacterString"`
			ContactEmail string   `xml:"pointOfContact>CI_ResponsibleParty>contactInfo>CI_Contact>address>CI_Address>electronicMailAddress>CharacterString"`
			LicenseUrl   []Anchor `xml:"resourceConstraints>MD_LegalConstraints>otherConstraints>Anchor"`

			ResponsibleParty *ResponsibleParty `xml:"pointOfContact>CI_ResponsibleParty>OrganisationName"`
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
				SpecificationAnchor Anchor `xml:"DQ_ConformanceResult>specification>CI_Citation>title>Anchor"`
				SpecificationTitle  string `xml:"DQ_ConformanceResult>specification>CI_Citation>title>CharacterString"`
				Explanation         string `xml:"DQ_ConformanceResult>explanation>CharacterString"`
				Pass                string `xml:"DQ_ConformanceResult>pass>Boolean"`
			} `xml:"DQ_DomainConsistency>result"`
		} `xml:"report"`
	} `xml:"dataQualityInfo>DQ_DataQuality"`
}

type ResponsibleParty struct {
	Char   string `xml:"CharacterString"`
	Anchor string `xml:"Anchor"`
}

type Anchor struct {
	Text string `xml:",chardata"`
	Href string `xml:"href,attr"`
}

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

func (m *MDMetadata) GetLicenseUrl() string {
	for _, val := range m.IdentificationInfo.MDDataIdentification.LicenseUrl {
		if strings.Contains(val.Href, "creativecommons.org") {
			return val.Href
		}
	}
	return ""
}

func (m *MDMetadata) GetThumbnailUrl() *string {
	if m.IdentificationInfo.MDDataIdentification != nil && m.IdentificationInfo.MDDataIdentification.GraphicOverview != nil {
		thumbnailUrl := m.IdentificationInfo.MDDataIdentification.GraphicOverview.MDBrowseGraphic.FileName
		return &thumbnailUrl
	}
	return nil
}

func (m *MDMetadata) GetInspireVariant() *inspire.InspireVariant {
	isInspire := false
	isConformant := true
	inspireRegulation := "VERORDENING (EU) Nr. 1089/2010"

	harmonised := inspire.Harmonised
	asIs := inspire.AsIs

	for _, report := range m.DQDataQuality.Report {
		for _, result := range report.ConsistencyResult {
			hasInspireRegulation := strings.Contains(result.SpecificationTitle, inspireRegulation)
			hasInspireSpecTitle := strings.HasPrefix(result.SpecificationTitle, "INSPIRE")

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

func (m *MDMetadata) GetInspireThemes() (themes []string) {
	for _, descriptiveKeyword := range m.IdentificationInfo.MDDataIdentification.DescriptiveKeywords {
		if descriptiveKeyword.MDKeywords.ThesaurusName.Href == "http://www.eionet.europa.eu/gemet/inspire_themes" {
			for _, keyword := range descriptiveKeyword.MDKeywords.Keyword {
				inspireThemeHrefPrefix := "http://www.eionet.europa.eu/gemet/nl/inspire-theme/"
				if strings.Contains(keyword.Anchor.Href, inspireThemeHrefPrefix) {
					themeValue := strings.ReplaceAll(keyword.Anchor.Href, inspireThemeHrefPrefix, "")
					themes = append(themes, themeValue)
				}
			}
		}
	}

	return
}

func (m *MDMetadata) GetHVDCategories() (categories []hvd.HVDCategory) {
	for _, descriptiveKeyword := range m.IdentificationInfo.MDDataIdentification.DescriptiveKeywords {
		if descriptiveKeyword.MDKeywords.ThesaurusName.Href == "http://publications.europa.eu/resource/dataset/high-value-dataset-category" {
			for _, keyword := range descriptiveKeyword.MDKeywords.Keyword {
				hvdCategoryHrefPrefix := "http://data.europa.eu/bna/"
				if strings.Contains(keyword.Anchor.Href, hvdCategoryHrefPrefix) {
					categoryValue := strings.ReplaceAll(keyword.Anchor.Href, hvdCategoryHrefPrefix, "")
					categories = append(categories, hvd.HVDCategory{Id: categoryValue})
				}
			}
		}
	}
	return
}

type SummaryRecord struct {
	Identifier string `xml:"identifier" json:"identifier"`
	Title      string `xml:"title" json:"title"`
	Type       string `xml:"type" json:"type"`
}
