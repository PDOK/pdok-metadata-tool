// Package generator holds the logic for generating metadata.
package generator

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pdok/pdok-metadata-tool/internal/common"
	"github.com/pdok/pdok-metadata-tool/pkg/model/codelist"
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/iso1911x"
	"github.com/pdok/pdok-metadata-tool/pkg/repository"
)

// ISO19119Generator is used for generating service metadata according to ISO19119 format.
type ISO19119Generator struct {
	MetadataHolder map[string]*MetadataEntry
	currentID      *string
	Codelist       *codelist.Codelist
	HVDRepository  *repository.HVDRepository
	outputDir      string
}

// NewISO19119Generator creates a new instance of the ISO19119 generator and sets up the metadata holder.
func NewISO19119Generator(
	serviceSpecifics ServiceSpecifics,
	outputDir string,
) (*ISO19119Generator, error) {
	// Setup holder for serviceSpecifics-config and output
	holder := make(map[string]*MetadataEntry)
	for _, serviceConfig := range serviceSpecifics.Services {
		holder[serviceConfig.ID] = &MetadataEntry{
			Config: serviceConfig,
		}
	}

	codelists, err := codelist.NewCodelist()
	if err != nil {
		return nil, err
	}

	hvdRepo := repository.NewHVDRepository(hvd.HvdEndpoint, common.HvdLocalRDFPath)

	return &ISO19119Generator{
		MetadataHolder: holder,
		Codelist:       codelists,
		HVDRepository:  hvdRepo,
		outputDir:      outputDir,
	}, nil
}

// CurrentEntry returns the current MetadataEntry for processing.
func (g *ISO19119Generator) CurrentEntry() (*MetadataEntry, error) {
	entry, ok := g.MetadataHolder[*g.currentID]
	if !ok {
		return nil, fmt.Errorf("no entry found for id: %s", *g.currentID)
	}

	return entry, nil
}

// MetadataEntry is used for processing the metadata generation.
type MetadataEntry struct {
	Config   ServiceConfig
	Metadata iso1911x.ISO19119
	filename string
}

func (e *MetadataEntry) setFilename() {
	e.filename = e.Config.ID + ".xml"
}

// Generate generates the metadata for each entry in the metadata holder.
func (g *ISO19119Generator) Generate() error {
	for id := range g.MetadataHolder {
		g.currentID = &id
		if err := g.SetMetadata(); err != nil {
			return err
		}

		if err := g.WriteToXML(); err != nil {
			return err
		}

		g.currentID = nil
	}

	return nil
}

// SetMetadata sets all the values for the metadata.
func (g *ISO19119Generator) SetMetadata() error {
	if err := g.setGeneralInfo(); err != nil {
		return err
	}

	if err := g.setIdentificationInfo(); err != nil {
		return err
	}

	if err := g.setDistributionInfo(); err != nil {
		return err
	}

	if err := g.setDataQualityInfo(); err != nil {
		return err
	}

	return nil
}

// WriteToXML writes the available metadata to XML.
func (g *ISO19119Generator) WriteToXML() error {
	entry, err := g.CurrentEntry()
	if err != nil {
		return err
	}

	output, err := xml.MarshalIndent(entry.Metadata, "", "  ")
	if err != nil {
		return err
	}

	perm := 750
	if err1 := os.MkdirAll(g.outputDir, os.FileMode(perm)); err1 != nil {
		return err1
	}

	entry.setFilename()

	path := filepath.Join(g.outputDir, entry.filename)

	perm = 0600
	if err = os.WriteFile(path, output, os.FileMode(perm)); err != nil {
		return err
	}

	return nil
}

// PrintSummary prints a summary of the generated metadata files.
func (g *ISO19119Generator) PrintSummary() {
	fmt.Printf("The following metadata has been created in %s: \n", g.outputDir)

	for _, entry := range g.MetadataHolder {
		fmt.Printf("  - %s\n", entry.filename)
	}
}

func (g *ISO19119Generator) setGeneralInfo() error {
	entry, err := g.CurrentEntry()
	if err != nil {
		return err
	}

	config := entry.Config

	entry.Metadata = iso1911x.ISO19119{
		XmlnsSrv:          "http://www.isotc211.org/2005/srv",
		XmlnsGmd:          "http://www.isotc211.org/2005/gmd",
		XmlnsGco:          "http://www.isotc211.org/2005/gco",
		XmlnsGml:          "http://www.opengis.net/gml",
		XmlnsXsi:          "http://www.w3.org/2001/XMLSchema-instance",
		XmlnsXs:           "http://www.w3.org/2001/XMLSchema",
		XmlnsCsw:          "http://www.opengis.net/cat/csw/2.0.2",
		XmlnsGmx:          "http://www.isotc211.org/2005/gmx",
		XmlnsGts:          "http://www.isotc211.org/2005/gts",
		XmlnsXlink:        "http://www.w3.org/1999/xlink",
		XsiSchemaLocation: "http://www.isotc211.org/2005/gmd  http://schemas.opengis.net/csw/2.0.2/profiles/apiso/1.0.0/apiso.xsd",
		FileIdentifier: iso1911x.CharacterStringTag{
			// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#metadata-unieke-identifier
			CharacterString: config.ID,
		},
		Language: iso1911x.LanguageTag{
			// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#taal-van-de-metadata
			LanguageCode: iso1911x.CodeListValueTag{
				CodeList:      "http://www.loc.gov/standards/iso639-2/",
				CodeListValue: "dut",
				Value:         "Nederlands; Vlaams",
			},
		},
		CharacterSet: iso1911x.CharacterSetTag{
			// https://docs.geostandaarden.nl/md/mdprofiel-iso19115/#x5-2-8-karakterset-van-de-bron
			MDCharacterSetCode: iso1911x.CodeListValueTag{
				CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#MD_CharacterSetCode",
				CodeListValue: "utf8",
				Value:         "utf8",
			},
		},
		HierarchyLevel: iso1911x.HierarchyLevelTag{
			// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#hiërarchieniveau
			// The fixed value for INSPIRE services is 'service'
			MDScopeCode: iso1911x.CodeListValueTag{
				CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#MD_ScopeCode",
				CodeListValue: "service",
				Value:         "service",
			},
		},
		HierarchyLevelName: iso1911x.CharacterStringTag{
			// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#hiërarchieniveaunaam
			CharacterString: "service",
		},
		Contact: iso1911x.ContactTag{
			ResponsibleParty: iso1911x.ResponsibleParty{
				// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#verantwoordelijke-organisatie-metadata
				OrganisationName: iso1911x.OrganisationNameTag{
					Anchor: iso1911x.AnchorTag{
						Href:  config.GetContactOrganisationURI(),
						Value: config.GetContactOrganisationName(),
					},
				},
				ContactInfo: iso1911x.ContactInfoTag{
					Contact: iso1911x.ContactDetails{
						Address: iso1911x.AddressTag{
							CIAddress: iso1911x.CIAddressTag{
								Email: iso1911x.CharacterStringTag{
									CharacterString: config.GetContactEmail(),
								},
							},
						},
						OnlineResource: iso1911x.OnlineResourceTag{
							CIOnlineResource: iso1911x.CIOnlineResourceTag{
								Linkage: iso1911x.URLTag{
									URL: config.GetContactURL(),
								},
							},
						},
					},
				},
				Role: iso1911x.RoleTag{
					// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#verantwoordelijke-organisatie-metadata:-rol
					CIRoleCode: iso1911x.CodeListValueTag{
						CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#CI_RoleCode",
						CodeListValue: "pointOfContact",
						Value:         "contactpunt",
					},
				},
			},
		},
		DateStamp: iso1911x.DateTag{
			// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#metadatadatum -->
			// Date on which the metadata was created or modified (format YYYY-MM-DD)
			Date: config.GetRevisionDate(),
		},
		MetadataStandardName: iso1911x.CharacterStringTag{
			// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#metadata-standaard-naam
			CharacterString: "ISO 19119",
		},
		MetadataStandardVersion: iso1911x.CharacterStringTag{
			// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#metadatastandaard-versie
			CharacterString: "Nederlands metadata profiel op ISO 19119 voor services 2.1.0",
		},
	}

	return nil
}

//nolint:funlen,maintidx
func (g *ISO19119Generator) setIdentificationInfo() error {
	entry, err := g.CurrentEntry()
	if err != nil {
		return err
	}

	config := entry.Config

	entry.Metadata.IdentificationInfo = iso1911x.IdentificationInfo{
		ServiceIdentification: iso1911x.ServiceIdentification{
			Citation: iso1911x.Citation{

				CICitation: iso1911x.CICitation{
					// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#titel-van-de-bron
					// This element must match with the element WMS_Capabilities/Service/Title in the Capabilities document
					Title: iso1911x.TitleTag{
						CharacterString: common.Ptr(config.GetTitle()),
					},
					Dates: []iso1911x.CIDateTag{
						{
							CIDate: iso1911x.CIDate{
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#x5-2-2-datum-van-de-bron
								// Date on which the service was created, format YYYY-MM-DD
								Date: iso1911x.DateTag{
									Date: config.GetCreationDate(),
								},
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#datum-type-van-de-bron
								DateType: iso1911x.DateTypeTag{
									CIDateTypeCode: iso1911x.CodeListValueTag{
										CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#CI_DateTypeCode",
										CodeListValue: "creation",
										Value:         "creatie", //nolint:misspell
									},
								},
							},
						},
						{
							CIDate: iso1911x.CIDate{
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#x5-2-2-datum-van-de-bron
								// Date on which the service was last revised, format YYYY-MM-DD
								Date: iso1911x.DateTag{
									Date: config.GetRevisionDate(),
								},
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#datum-type-van-de-bron
								DateType: iso1911x.DateTypeTag{
									CIDateTypeCode: iso1911x.CodeListValueTag{
										CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#CI_DateTypeCode",
										CodeListValue: "revision",
										Value:         "revisie",
									},
								},
							},
						},
					},
				},
			},
			Abstract: iso1911x.CharacterStringTag{
				// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#samenvatting
				// This element must match with the element WMS_Capabilities/Service/Abstract in the Capabilities document
				CharacterString: config.GetAbstract(),
			},
			PointOfContact: iso1911x.ContactTag{
				// The organisation which is responsible for the service
				ResponsibleParty: iso1911x.ResponsibleParty{
					// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#verantwoordelijke-organisatie-bron
					OrganisationName: iso1911x.OrganisationNameTag{
						Anchor: iso1911x.AnchorTag{
							Href:  config.GetContactOrganisationURI(),
							Value: config.GetContactOrganisationName(),
						},
					},
					ContactInfo: iso1911x.ContactInfoTag{
						Contact: iso1911x.ContactDetails{
							Address: iso1911x.AddressTag{
								CIAddress: iso1911x.CIAddressTag{
									// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#verantwoordelijke-organisatie-bron-email
									Email: iso1911x.CharacterStringTag{
										CharacterString: config.GetContactEmail(),
									},
								},
							},
							OnlineResource: iso1911x.OnlineResourceTag{
								CIOnlineResource: iso1911x.CIOnlineResourceTag{
									Linkage: iso1911x.URLTag{
										URL: config.GetContactURL(),
									},
								},
							},
						},
					},
					Role: iso1911x.RoleTag{
						// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#verantwoordelijke-organisatie-bron:-rol
						CIRoleCode: iso1911x.CodeListValueTag{
							CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#CI_RoleCode",
							CodeListValue: "custodian",
							Value:         "beheerder",
						},
					},
				},
			},
		},
	}

	// setThumbnails
	thumbnails := config.GetThumbnails()

	if thumbnails != nil {
		entry.Metadata.IdentificationInfo.ServiceIdentification.GraphicOverview = []iso1911x.GraphicOverviewTag{}
		for _, thumbnail := range thumbnails {
			graphicOverview := iso1911x.GraphicOverviewTag{
				BrowseGraphic: iso1911x.BrowseGraphic{
					FileName: iso1911x.CharacterStringTag{CharacterString: thumbnail.File},
					FileDescription: iso1911x.CharacterStringTag{
						CharacterString: thumbnail.Description,
					},
					FileType: nil,
				},
			}
			if thumbnail.Filetype != "" {
				graphicOverview.BrowseGraphic.FileType = common.Ptr(
					iso1911x.CharacterStringTag{CharacterString: thumbnail.Filetype},
				)
			}

			entry.Metadata.IdentificationInfo.ServiceIdentification.GraphicOverview = append(
				entry.Metadata.IdentificationInfo.ServiceIdentification.GraphicOverview,
				graphicOverview,
			)
		}
	}

	// The descriptive keywords match with the elements under:
	// - WMS_Capabilities/Services/KeywordList
	// - WMS_Capabilities/Capability/inspire_vs:ExtendedCapabilities/inspire_common:MandatoryKeyword
	// - WMS_Capabilities/Capability/inspire_vs:ExtendedCapabilities/inspire_common:Keyword
	// in the Capabilities document
	keywords := config.GetKeywords()

	protocol, ok := g.Codelist.GetProtocolDetailsByProtocol(config.Type)
	if !ok {
		return fmt.Errorf("no protcol found for service type: %s", config.Type)
	}

	entry.Metadata.IdentificationInfo.ServiceIdentification.DescriptiveKeywords = []iso1911x.DescriptiveKeywordsTag{}

	descriptiveKeyword := iso1911x.DescriptiveKeywordsTag{
		// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#trefwoord
		Keywords: &iso1911x.MDKeywords{
			Keyword: []iso1911x.KeywordTag{},
		},
	}

	// For INSPIRE, at least one of the mandatory keywords from Metadata IR Part 4 is required (i.e. infoMapAccessService, infoFeatureAccessService, etc.)
	if len(config.GetInspireThemes()) > 0 {
		// There has been discussion about an inconsistency in the spec regarding the NL profile vs INSPIRE.
		// See https://github.com/INSPIRE-MIF/helpdesk-validator/issues/84
		// NL Profile requires the anchor to have the value SpatialDataserviceCategoryURI, but INSPIRE requires the value to be SpatialDataserviceCategory
		keywordTag := iso1911x.KeywordTag{
			Anchor: &iso1911x.AnchorTag{
				Href:  protocol.SpatialDataserviceCategoryURI,
				Value: protocol.SpatialDataserviceCategory,
			},
		}
		descriptiveKeyword.Keywords.Keyword = append(
			descriptiveKeyword.Keywords.Keyword,
			keywordTag,
		)
	}

	for _, keyword := range keywords {
		keywordTag := iso1911x.KeywordTag{
			CharacterString: &keyword,
		}
		descriptiveKeyword.Keywords.Keyword = append(
			descriptiveKeyword.Keywords.Keyword,
			keywordTag,
		)
	}

	entry.Metadata.IdentificationInfo.ServiceIdentification.DescriptiveKeywords = append(
		entry.Metadata.IdentificationInfo.ServiceIdentification.DescriptiveKeywords,
		descriptiveKeyword,
	)

	// INSPIRE theme as keyword
	if len(config.GetInspireThemes()) > 0 {
		inspireDescriptiveKeyword := iso1911x.DescriptiveKeywordsTag{
			Keywords: &iso1911x.MDKeywords{
				Keyword: []iso1911x.KeywordTag{},
				Type: &iso1911x.KeywordTypeTag{
					Code: iso1911x.CodeListValueTag{
						CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#MD_KeywordTypeCode",
						CodeListValue: "theme",
						Value:         "theme",
					},
				},
				ThesaurusName: &iso1911x.Citation{
					// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#thesaurus
					// The GEMET Thesaurus in which the INSPIRE theme is defined
					CICitation: iso1911x.CICitation{
						Title: iso1911x.TitleTag{
							Anchor: &iso1911x.AnchorTag{
								Href:  "https://www.eionet.europa.eu/gemet/nl/inspire-themes/",
								Value: "GEMET - INSPIRE themes, version 1.0",
							},
						},
						Dates: []iso1911x.CIDateTag{
							{
								CIDate: iso1911x.CIDate{
									Date: iso1911x.DateTag{
										// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#thesaurusdatum
										Date: "2008-06-01",
									},
									DateType: iso1911x.DateTypeTag{
										// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#thesaurusdatum-type
										CIDateTypeCode: iso1911x.CodeListValueTag{
											CodeList:      "https://standards.iso.org/ittf/PubliclyAvailableStandards/ISO_19139_Schemas/resources/Codelist/gmxCodelists.xml#CI_DateTypeCode",
											CodeListValue: "publication",
											Value:         "publicatie",
										},
									},
								},
							},
						},
						Identifier: &iso1911x.IdentifierTag{
							MDIdentifier: iso1911x.MDIdentifier{
								Code: iso1911x.CodeTag{
									Anchor: iso1911x.AnchorTag{
										Href:  "https://www.nationaalgeoregister.nl/geonetwork/srv/api/registries/vocabularies/external.theme.httpinspireeceuropaeutheme-theme",
										Value: "geonetwork.thesaurus.external.theme.httpinspireeceuropaeutheme-theme",
									},
								},
							},
						},
					},
				},
			},
		}

		for _, inspireTheme := range config.GetInspireThemes() {
			inspireThemeLabel, ok := g.Codelist.GetINSPIREThemeLabelByURI(inspireTheme)
			if !ok {
				return fmt.Errorf("no INSPIRE theme found for code: %s", inspireTheme)
			}

			inspireKeyword := iso1911x.KeywordTag{
				// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#trefwoord
				// Must match the element WMS_Capabilities/Capability/inspire_vs:ExtendedCapability/inspire_common:Keyword in the Capabilities file
				// in which the INSPIRE theme as defined in the GEMET Thesaurus is included
				// Name of the INSPIRE theme as defined in the GEMET Thesaurus and written in the language of this metadata document
				Anchor: &iso1911x.AnchorTag{
					Href:  inspireTheme,
					Value: *inspireThemeLabel,
				},
			}
			inspireDescriptiveKeyword.Keywords.Keyword = append(
				inspireDescriptiveKeyword.Keywords.Keyword,
				inspireKeyword,
			)
		}

		entry.Metadata.IdentificationInfo.ServiceIdentification.DescriptiveKeywords = append(
			entry.Metadata.IdentificationInfo.ServiceIdentification.DescriptiveKeywords,
			inspireDescriptiveKeyword,
		)
	}

	// If an HVD category is linked, this must be made clear by means of a keyword, see https://docs.geostandaarden.nl/eu/handreiking-hvd/#409368F9
	hvdCategories := config.GetHvdCategories()

	if len(hvdCategories) > 0 {
		hvdKeywordTag := iso1911x.KeywordTag{
			Anchor: &iso1911x.AnchorTag{
				Href:  "http://data.europa.eu/eli/reg_impl/2023/138/oj",
				Value: "HVD",
			},
		}
		descKeyword := entry.Metadata.IdentificationInfo.ServiceIdentification.DescriptiveKeywords[0]
		descKeyword.Keywords.Keyword = append(descKeyword.Keywords.Keyword, hvdKeywordTag)

		var keywordTags []iso1911x.KeywordTag

		filteredHvdCategories, err := g.HVDRepository.GetFilteredHvdCategories(hvdCategories)
		if err != nil {
			return err
		}

		for _, hvdCategory := range filteredHvdCategories {
			fullHvdCategoryCode := "http://data.europa.eu/bna/" + hvdCategory.ID

			hvdCategoryKeyword := iso1911x.KeywordTag{
				Anchor: &iso1911x.AnchorTag{
					Href:  fullHvdCategoryCode,
					Value: hvdCategory.LabelDutch,
				},
			}
			keywordTags = append(keywordTags, hvdCategoryKeyword)
		}

		descriptiveKeyword := iso1911x.DescriptiveKeywordsTag{
			Keywords: &iso1911x.MDKeywords{
				Keyword: keywordTags,
				// reference to the value list for HVD themes and subthemes
				ThesaurusName: &iso1911x.Citation{
					CICitation: iso1911x.CICitation{
						Title: iso1911x.TitleTag{
							Anchor: &iso1911x.AnchorTag{
								Href:  "http://publications.europa.eu/resource/dataset/high-value-dataset-category",
								Value: "High-value dataset categories",
							},
						},
						Dates: []iso1911x.CIDateTag{
							{
								CIDate: iso1911x.CIDate{
									Date: iso1911x.DateTag{
										Date: "2023-09-27",
									},
									DateType: iso1911x.DateTypeTag{
										CIDateTypeCode: iso1911x.CodeListValueTag{
											CodeList:      "http://standards.iso.org/iso/19139/resources/gmxCodelists.xml#CI_DateTypeCode",
											CodeListValue: "publication",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		entry.Metadata.IdentificationInfo.ServiceIdentification.DescriptiveKeywords = append(
			entry.Metadata.IdentificationInfo.ServiceIdentification.DescriptiveKeywords,
			descriptiveKeyword,
		)
	}

	// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#x5-2-12-juridische-toegangsrestricties
	entry.Metadata.IdentificationInfo.ServiceIdentification.ResourceConstraints = []iso1911x.ResourceConstraint{}

	constraint1 := iso1911x.ResourceConstraint{
		MDConstraints: &iso1911x.MDConstraints{
			// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#gebruiksbeperkingen
			// Applications for which the service is not suitable.
			UseLimitation: iso1911x.CharacterStringTag{
				CharacterString: config.GetUseLimitation(),
			},
		},
	}

	licenseURI := config.GetServiceLicense()

	dataLicense, ok := g.Codelist.GetDataLicenseByURI(licenseURI)
	if !ok {
		return fmt.Errorf("no data license found for license URI: %s", licenseURI)
	}

	entry.Metadata.IdentificationInfo.ServiceIdentification.ResourceConstraints = append(
		entry.Metadata.IdentificationInfo.ServiceIdentification.ResourceConstraints,
		constraint1,
	)
	resourceConstraint := iso1911x.ResourceConstraint{
		MDLegalConstraints: &iso1911x.MDLegalConstraints{
			// Must match the element WMS_Capabilities/Service/AccessConstraints in the Capabilities file
			// If there are no usage restrictions: use "otherRestrictions" in the MD_RestrictionCode element and include a reference to a Public Domain declaration or CC0 in the otherConstraints
			// Otherwise, use another Creative Commons license; if that’s not sufficient, create a geo-shared license and include a reference to that license in otherConstraints
			// For INSPIRE, also include a code from the ConditionsApplyingToAccessAndUse code list in a second otherConstraints element within the same MD_LegalConstraints
			AccessConstraints: []iso1911x.AccessConstraintTag{
				{
					// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#x5-2-12-juridische-toegangsrestricties
					MDRestrictionCode: iso1911x.CodeListValueTag{
						CodeListValue: "otherRestrictions",
						CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#MD_RestrictionCode",
						Value:         "anders",
					},
				},
			},
			OtherConstraints: []iso1911x.OtherConstraintTag{
				// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#overige-beperkingen
				{
					Anchor: iso1911x.AnchorTag{
						Href:  licenseURI,
						Value: dataLicense.Description,
					},
				},
			},
		},
	}

	if config.ServiceInspireType != nil {
		inspireOtherConstraint := iso1911x.OtherConstraintTag{
			Anchor: iso1911x.AnchorTag{
				Href:  "http://inspire.ec.europa.eu/metadata-codelist/ConditionsApplyingToAccessAndUse/noConditionsApply",
				Value: "Geen condities voor toegang en gebruik",
			},
		}
		resourceConstraint.MDLegalConstraints.OtherConstraints = append(
			resourceConstraint.MDLegalConstraints.OtherConstraints,
			inspireOtherConstraint,
		)
	}

	entry.Metadata.IdentificationInfo.ServiceIdentification.ResourceConstraints = append(
		entry.Metadata.IdentificationInfo.ServiceIdentification.ResourceConstraints,
		resourceConstraint,
	)

	if config.ServiceInspireType != nil {
		// For INSPIRE, also include a code from the LimitationsOnPublicAccess code list in an additional MD_LegalConstraints element
		inspireResourceConstraint := iso1911x.ResourceConstraint{
			MDLegalConstraints: &iso1911x.MDLegalConstraints{
				AccessConstraints: []iso1911x.AccessConstraintTag{
					{
						MDRestrictionCode: iso1911x.CodeListValueTag{
							CodeListValue: "otherRestrictions",
							CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#MD_RestrictionCode",
							Value:         "anders",
						},
					},
				},
				OtherConstraints: []iso1911x.OtherConstraintTag{
					{
						Anchor: iso1911x.AnchorTag{
							Href:  "http://inspire.ec.europa.eu/metadata-codelist/LimitationsOnPublicAccess/noLimitations",
							Value: "Geen beperkingen",
						},
					},
				},
			},
		}

		entry.Metadata.IdentificationInfo.ServiceIdentification.ResourceConstraints = append(
			entry.Metadata.IdentificationInfo.ServiceIdentification.ResourceConstraints,
			inspireResourceConstraint,
		)
	}

	// serviceType
	inspireServiceType, ok := g.Codelist.GetInspireServiceTypeByServiceType(config.Type)
	if !ok {
		return fmt.Errorf("no INSPIRE service type found for type: %s", config.Type)
	}

	entry.Metadata.IdentificationInfo.ServiceIdentification.ServiceType = iso1911x.ServiceTypeTag{
		// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#service-type
		// Must match the element WMS_Capabilities/Capability/inspire_vs:ExtendedCapabilities/inspire_common:SpatialDataServiceType
		LocalName: iso1911x.LocalNameTag{
			CodeSpace: "http://inspire.ec.europa.eu/metadata-codelist/SpatialDataServiceType",
			Value:     inspireServiceType.InspireServiceType,
		},
	}

	// Extent
	boundingBox := config.GetBoundingBox()
	entry.Metadata.IdentificationInfo.ServiceIdentification.Extent = iso1911x.ExtentTag{
		// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#Omgrenzende%20rechthoek
		// Must match with the element WMS_Capabilities/Capability/Layer/Ex_GeographicBoundingBox in the Capabilities document
		EXExtent: iso1911x.EXExtentTag{
			GeographicElement: iso1911x.GeographicElementTag{
				GeographicBoundingBox: iso1911x.GeographicBoundingBoxTag{
					WestBoundLongitude: iso1911x.DecimalTag{Value: boundingBox.MinX},
					EastBoundLongitude: iso1911x.DecimalTag{Value: boundingBox.MaxX},
					SouthBoundLatitude: iso1911x.DecimalTag{Value: boundingBox.MinY},
					NorthBoundLatitude: iso1911x.DecimalTag{Value: boundingBox.MaxY},
				},
			},
		},
	}

	// Coupling type
	entry.Metadata.IdentificationInfo.ServiceIdentification.CouplingType = iso1911x.CouplingTypeTag{
		// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#koppel-type
		// Fixed value 'tight' for a View or Download service
		SVCouplingType: iso1911x.CodeListValueTag{
			CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#SV_CouplingType",
			CodeListValue: "tight",
			Value:         "tight",
		},
	}

	entry.Metadata.IdentificationInfo.ServiceIdentification.ContainsOperations = []iso1911x.OperationMetadataTag{
		{
			OperationMetadata: iso1911x.SVOperationMetadata{
				OperationName: iso1911x.OperationNameTag{
					// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#operatie-naam
					// Name of the operation, i.e. GetCapabilities
					CharacterString: protocol.ServiceAccessPointOperation,
				},
				DCP: iso1911x.DCPTag{
					// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#DCP
					DCPList: iso1911x.CodeListValueTag{
						CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#DCPList",
						CodeListValue: "WebServices",
						Value:         "WebServices",
					},
				},
				ConnectPoint: iso1911x.ConnectPointTag{
					// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#connectie-url
					// The accessPoint of the service, which includes the operations and endpoints.
					// For OGC services, this is the URL to the capabilities document
					OnlineResource: iso1911x.CIOnlineResource{
						Linkage: iso1911x.URLTag{
							URL: config.AccessPoint,
						},
					},
				},
			},
		},
	}

	// OperatesOn
	expectedSize := 5

	operatesOn := make([]iso1911x.OperatesOn, 0, expectedSize)
	for _, linkedDataset := range config.GetLinkedDatasets() {
		// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#gekoppelde-bron
		// The attribute xlink:href must contain a URI pointing to the MD_DataIdentification section in the XML of the dataset's metadata record
		operatesOn = append(operatesOn, iso1911x.OperatesOn{
			UUIDRef: linkedDataset,
			Href:    "https://nationaalgeoregister.nl/geonetwork/srv/dut/csw?service=CSW&request=GetRecordById&version=2.0.2&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full&id=" + linkedDataset + "#MD_DataIdentification",
		})
	}

	if len(operatesOn) > 0 {
		entry.Metadata.IdentificationInfo.ServiceIdentification.OperatesOn = operatesOn
	}

	return nil
}

func (g *ISO19119Generator) setDistributionInfo() error {
	entry, err := g.CurrentEntry()
	if err != nil {
		return err
	}

	config := entry.Config

	protocol, ok := g.Codelist.GetProtocolDetailsByProtocol(config.Type)
	if !ok {
		return fmt.Errorf("no protcol found for service type: %s", config.Type)
	}

	entry.Metadata.DistributionInfo = iso1911x.DistributionInfo{
		Distribution: iso1911x.Distribution{
			TransferOptions: iso1911x.TransferOptions{
				DigitalTransferOptions: iso1911x.DigitalTransferOptions{
					Online: iso1911x.OnlineResourceWrapper{
						Resource: iso1911x.CIOnlineResource{
							Linkage: iso1911x.URLTag{
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#url
								// Reference to the Capabilities document of the service
								URL: config.AccessPoint,
							},
							Protocol: &iso1911x.ProtocolTag{
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#protocol
								Anchor: iso1911x.AnchorTag{
									Href:  protocol.ServiceProtocolURL,
									Value: protocol.ServiceProtocol,
								},
							},
							Description: &iso1911x.DescriptionTag{
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#omschrijving
								Anchor: iso1911x.AnchorTag{
									Href:  "http://inspire.ec.europa.eu/metadata-codelist/OnLineDescriptionCode/accessPoint",
									Value: "accessPoint",
								},
							},
						},
					},
				},
			},
		},
	}

	return nil
}

//nolint:funlen,maintidx
func (g *ISO19119Generator) setDataQualityInfo() error {
	entry, err := g.CurrentEntry()
	if err != nil {
		return err
	}

	config := entry.Config

	entry.Metadata.DataQualityInfo = iso1911x.DataQualityInfo{
		DataQuality: iso1911x.DataQuality{
			Scope: iso1911x.ScopeTag{
				Scope: iso1911x.ScopeDetails{
					Level: iso1911x.LevelTag{
						// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#niveau-kwaliteitsbeschrijving
						MDScopeCode: iso1911x.CodeListValueTag{
							CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#MD_ScopeCode",
							CodeListValue: "service",
							Value:         "service",
						},
					},
					LevelDescription: iso1911x.LevelDescriptionTag{
						ScopeDescription: iso1911x.ScopeDescriptionTag{
							Other: iso1911x.CharacterStringTag{
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#niveau-kwaliteitsbeschrijving-naam
								CharacterString: "service",
							},
						},
					},
				},
			},
		},
	}

	inspireServiceType, ok := g.Codelist.GetInspireServiceTypeByServiceType(config.Type)
	if !ok {
		return fmt.Errorf("no INSPIRE service type found for type: %s", config.Type)
	}

	// https://docs.geostandaarden.nl/eu/INSPIRE-handreiking/#invulinstructie-service-metadata
	if config.ServiceInspireType != nil && *config.ServiceInspireType == NetworkService {
		entry.Metadata.DataQualityInfo.DataQuality.Report = []iso1911x.ReportTag{
			{
				DomainConsistency: &iso1911x.DomainConsistencyTag{
					Result: iso1911x.ConformanceResultTag{
						DQConformanceResult: iso1911x.DQConformanceResult{
							Specification: iso1911x.Citation{
								CICitation: iso1911x.CICitation{
									Title: iso1911x.TitleTag{
										// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#specificatie
										Anchor: &iso1911x.AnchorTag{
											Href:  "https://data.europa.eu/eli/reg/2009/976",
											Value: "VERORDENING (EG) Nr. 976/2009 VAN DE COMMISSIE van 19 oktober 2009 tot uitvoering van Richtlijn 2007/2/EG van het Europees Parlement en de Raad wat betreft de netwerkdiensten",
										},
									},
									Dates: []iso1911x.CIDateTag{
										{
											CIDate: iso1911x.CIDate{
												Date: iso1911x.DateTag{
													// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#specificatiedatum
													Date: "2009-10-19",
												},
												DateType: iso1911x.DateTypeTag{
													// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#specificatiedatum-type
													CIDateTypeCode: iso1911x.CodeListValueTag{
														CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#CI_DateTypeCode",
														CodeListValue: "publication",
														Value:         "publicatie",
													},
												},
											},
										},
									},
								},
							},
							Explanation: iso1911x.CharacterStringTag{
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#verklaring
								CharacterString: "Conform verordening",
							},
							Pass: iso1911x.BooleanTag{
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#conformiteit-indicatie-met-de-specificatie
								Value: true,
							},
						},
					},
				},
			},
			{
				DomainConsistency: &iso1911x.DomainConsistencyTag{
					Result: iso1911x.ConformanceResultTag{
						DQConformanceResult: iso1911x.DQConformanceResult{
							Specification: iso1911x.Citation{
								CICitation: iso1911x.CICitation{
									Title: iso1911x.TitleTag{
										// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#specificatie
										Anchor: &iso1911x.AnchorTag{
											Href:  inspireServiceType.InspireTechnicalGuidance,
											Value: "Technical Guidance for the implementation of INSPIRE " + inspireServiceType.InspireServiceType + " Services",
										},
									},
									Dates: []iso1911x.CIDateTag{
										{
											CIDate: iso1911x.CIDate{
												Date: iso1911x.DateTag{
													// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#specificatiedatum
													Date: inspireServiceType.InspireTechnicalGuidanceDate,
												},
												DateType: iso1911x.DateTypeTag{
													// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#specificatiedatum-type
													CIDateTypeCode: iso1911x.CodeListValueTag{
														CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#CI_DateTypeCode",
														CodeListValue: "publication",
														Value:         "publicatie",
													},
												},
											},
										},
									},
								},
							},
							Explanation: iso1911x.CharacterStringTag{
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#verklaring
								CharacterString: "Conform technische specificatie",
							},
							Pass: iso1911x.BooleanTag{
								// https://docs.geostandaarden.nl/md/mdprofiel-iso19119/#conformiteit-indicatie-met-de-specificatie
								Value: true,
							},
						},
					},
				},
			},
		}
	}

	// https://docs.geostandaarden.nl/eu/INSPIRE-handreiking/#invulinstructie-invocable-sds-metadata
	// https://docs.geostandaarden.nl/eu/INSPIRE-handreiking/#invulinstructie-interoperable-sds-metadata
	if config.ServiceInspireType != nil &&
		(*config.ServiceInspireType == Invocable || *config.ServiceInspireType == Interoperable) {
		SDSServiceCategory, ok := g.Codelist.GetSDSServiceCategoryBySDSCategory(
			string(*config.ServiceInspireType),
		)
		if !ok {
			return fmt.Errorf("no INSPIRE service type found for type: %s", config.Type)
		}

		protocol, ok := g.Codelist.GetProtocolDetailsByProtocol(config.Type)
		if !ok {
			return fmt.Errorf("no protcol found for service type: %s", config.Type)
		}

		entry.Metadata.DataQualityInfo.DataQuality.Report = []iso1911x.ReportTag{
			{
				DomainConsistency: &iso1911x.DomainConsistencyTag{
					Result: iso1911x.ConformanceResultTag{
						DQConformanceResult: iso1911x.DQConformanceResult{
							Specification: iso1911x.Citation{
								CICitation: iso1911x.CICitation{
									Title: iso1911x.TitleTag{
										Anchor: &iso1911x.AnchorTag{
											Href:  "https://data.europa.eu/eli/reg/2010/1089",
											Value: "VERORDENING (EU) Nr. 1089/2010 VAN DE COMMISSIE van 23 november 2010 ter uitvoering van Richtlijn 2007/2/EG van het Europees Parlement en de Raad betreffende de interoperabiliteit van verzamelingen ruimtelijke gegevens en van diensten met betrekking tot ruimtelijke gegevens",
										},
									},
									Dates: []iso1911x.CIDateTag{
										{
											CIDate: iso1911x.CIDate{
												Date: iso1911x.DateTag{
													Date: "2010-12-08",
												},
												DateType: iso1911x.DateTypeTag{
													CIDateTypeCode: iso1911x.CodeListValueTag{
														CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#CI_DateTypeCode",
														CodeListValue: "publication",
														Value:         "publicatie",
													},
												},
											},
										},
									},
								},
							},
							Explanation: iso1911x.CharacterStringTag{
								CharacterString: "Conform verordening",
							},
							Pass: iso1911x.BooleanTag{
								Value: true,
							},
						},
					},
				},
			},
			{
				DomainConsistency: &iso1911x.DomainConsistencyTag{
					Result: iso1911x.ConformanceResultTag{
						DQConformanceResult: iso1911x.DQConformanceResult{
							Specification: iso1911x.Citation{
								CICitation: iso1911x.CICitation{
									Title: iso1911x.TitleTag{
										Anchor: &iso1911x.AnchorTag{
											Href:  SDSServiceCategory.URI,
											Value: SDSServiceCategory.Value,
										},
									},
									Dates: []iso1911x.CIDateTag{
										{
											CIDate: iso1911x.CIDate{
												Date: iso1911x.DateTag{
													Date: "2014-12-11",
												},
												DateType: iso1911x.DateTypeTag{
													CIDateTypeCode: iso1911x.CodeListValueTag{
														CodeList:      "http://www.isotc211.org/2005/resources/codeList.xml#CI_DateTypeCode",
														CodeListValue: "publication",
														Value:         "publicatie",
													},
												},
											},
										},
									},
								},
							},
							Explanation: iso1911x.CharacterStringTag{
								CharacterString: "De service voldoet aan de requirements van de " + SDSServiceCategory.Value + " conformance class",
							},
							Pass: iso1911x.BooleanTag{
								Value: true,
							},
						},
					},
				},
			},
			{
				DomainConsistency: &iso1911x.DomainConsistencyTag{
					Result: iso1911x.ConformanceResultTag{
						DQConformanceResult: iso1911x.DQConformanceResult{
							Specification: iso1911x.Citation{
								CICitation: iso1911x.CICitation{
									Title: iso1911x.TitleTag{
										Anchor: &iso1911x.AnchorTag{
											Href:  protocol.ServiceProtocolURL,
											Value: protocol.ServiceProtocolName,
										},
									},
									Dates: []iso1911x.CIDateTag{
										{
											CIDate: iso1911x.CIDate{
												Date: iso1911x.DateTag{
													Date: protocol.ProtocolReleaseDate,
												},
												DateType: iso1911x.DateTypeTag{
													CIDateTypeCode: iso1911x.CodeListValueTag{
														CodeList:      "https://standards.iso.org/iso/19139/resources/gmxCodelists.xml#CI_DateTypeCode",
														CodeListValue: "publication",
														Value:         "publicatie",
													},
												},
											},
										},
									},
								},
							},
							Explanation: iso1911x.CharacterStringTag{
								CharacterString: "is conform " + protocol.ServiceProtocolName + " specificatie",
							},
							Pass: iso1911x.BooleanTag{
								Value: true,
							},
						},
					},
				},
			},
		}

		if *config.ServiceInspireType == Interoperable {
			qosAvailabilityReport := iso1911x.ReportTag{
				ConceptualConsistency: &iso1911x.ConceptualConsistencyTag{
					NameOfMeasure: iso1911x.AnchorTag{
						Href:  "http://inspire.ec.europa.eu/metadata-codelist/QualityOfServiceCriteria/availability",
						Value: "beschikbaarheid",
					},
					MeasureDescription: iso1911x.CharacterStringTag{
						CharacterString: "Beschikbaarheid op jaarbasis, uitgedrukt in percentage in tijd",
					},
					Result: iso1911x.QuantitativeResult{
						DQQuantitativeResult: iso1911x.DQQuantitativeResult{
							ValueUnit: iso1911x.ValueUnitTag{
								Href: "urn:ogc:def:uom:OGC::percent",
							},
							Value: iso1911x.RecordTag{
								Type:  "xs:double",
								Value: config.GetQosAvailability(),
							},
						},
					},
				},
			}
			entry.Metadata.DataQualityInfo.DataQuality.Report = append(
				entry.Metadata.DataQualityInfo.DataQuality.Report,
				qosAvailabilityReport,
			)

			qosPerformanceReport := iso1911x.ReportTag{
				ConceptualConsistency: &iso1911x.ConceptualConsistencyTag{
					NameOfMeasure: iso1911x.AnchorTag{
						Href:  "http://inspire.ec.europa.eu/metadata-codelist/QualityOfServiceCriteria/performance",
						Value: "performance",
					},
					MeasureDescription: iso1911x.CharacterStringTag{
						CharacterString: "Gemiddelde response tijd, uitgedrukt in seconden",
					},
					Result: iso1911x.QuantitativeResult{
						DQQuantitativeResult: iso1911x.DQQuantitativeResult{
							ValueUnit: iso1911x.ValueUnitTag{
								Href: "http://www.opengis.net/def/uom/SI/second",
							},
							Value: iso1911x.RecordTag{
								Type:  "xs:double",
								Value: config.GetQosPerformance(),
							},
						},
					},
				},
			}
			entry.Metadata.DataQualityInfo.DataQuality.Report = append(
				entry.Metadata.DataQualityInfo.DataQuality.Report,
				qosPerformanceReport,
			)

			qosCapacityReport := iso1911x.ReportTag{
				ConceptualConsistency: &iso1911x.ConceptualConsistencyTag{
					NameOfMeasure: iso1911x.AnchorTag{
						Href:  "http://inspire.ec.europa.eu/metadata-codelist/QualityOfServiceCriteria/capacity",
						Value: "capaciteit",
					},
					MeasureDescription: iso1911x.CharacterStringTag{
						CharacterString: "Maximum aantal gelijktijdige requests per seconde die aan de performance criteria voldoen, uitgedrukt als aantal requests per seconde",
					},
					Result: iso1911x.QuantitativeResult{
						DQQuantitativeResult: iso1911x.DQQuantitativeResult{
							ValueUnit: iso1911x.ValueUnitTag{
								Href: "http://www.opengis.net/def/uom/OGC/1.0/unity",
							},
							Value: iso1911x.RecordTag{
								Type:  "xs:integer",
								Value: config.GetQosCapacity(),
							},
						},
					},
				},
			}
			entry.Metadata.DataQualityInfo.DataQuality.Report = append(
				entry.Metadata.DataQualityInfo.DataQuality.Report,
				qosCapacityReport,
			)
		}
	}

	return nil
}
