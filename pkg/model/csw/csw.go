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

type GetRecordsConstraint struct {
	MetadataType     *MetadataType
	OrganisationName *string
}

func (grc *GetRecordsConstraint) ToQueryParameter() (constraint string) {
	var constraints []string

	if grc.MetadataType != nil {
		constraints = append(constraints, fmt.Sprintf("type='%s'", *grc.MetadataType))
	}
	if grc.OrganisationName != nil {
		constraints = append(constraints, fmt.Sprintf("OrganisationName='%s'", *grc.OrganisationName))
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
			ConsistencyResult struct {
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

	harmonized := inspire.Harmonized
	asIs := inspire.AsIs

	for _, report := range m.DQDataQuality.Report {
		if strings.Contains(report.ConsistencyResult.SpecificationTitle, inspireRegulation) {
			isInspire = true
		}
		if strings.HasPrefix(report.ConsistencyResult.SpecificationTitle, "INSPIRE") && report.ConsistencyResult.Pass != "true" {
			isConformant = false
		}
	}

	switch {
	case isInspire && isConformant:
		return &harmonized
	case isInspire && !isConformant:
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
	Identifier string `xml:"identifier"`
	Title      string `xml:"title"`
	Type       string `xml:"type"`
}
