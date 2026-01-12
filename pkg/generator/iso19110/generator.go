// Package iso19110 holds the logic for generating iso19110 metadata.
package iso19110

import (
	"github.com/pdok/pdok-metadata-tool/pkg/generator/core"
	"github.com/pdok/pdok-metadata-tool/pkg/model/iso1911x"
)

type Entry = core.MetadataEntry[iso1911x.ISO19110, FeatureCatalogueConfig]

type Generator struct {
	*core.Generator[iso1911x.ISO19110, FeatureCatalogueConfig]
}

func NewGenerator(
	spec FeatureCatalogueSpecifics,
	outputDir string,
) (*Generator, error) {
	// Setup holder for featureCatalogue-config and output
	holder := make(map[string]*Entry)
	for _, featureCatalogueConfig := range spec.FeatureCatalogues {
		holder[featureCatalogueConfig.ID] = &Entry{
			Config: featureCatalogueConfig,
		}
	}

	// Setup base generator
	base := &core.Generator[iso1911x.ISO19110, FeatureCatalogueConfig]{
		MetadataHolder: holder,
		OutputDir:      outputDir,
	}

	return &Generator{Generator: base}, nil
}

// Generate generates metadata and writes a file for each entry in the metadata holder.
func (g *Generator) Generate() error {
	if err := g.generateMetadataEntries(); err != nil {
		return err
	}

	for id := range g.MetadataHolder {
		g.CurrentID = &id

		if err := g.WriteToFile(); err != nil {
			return err
		}

		g.CurrentID = nil
	}

	return nil
}

// SetMetadata sets all the values for the metadata.
func (g *Generator) SetMetadata() error {
	if err := g.setGeneralInfo(); err != nil {
		return err
	}

	if err := g.setProducerInfo(); err != nil {
		return err
	}

	if err := g.setFeatureTypeInfo(); err != nil {
		return err
	}

	return nil
}

// generateMetadataEntries generates the metadata for each entry in the metadata holder.
func (g *Generator) generateMetadataEntries() error {
	for id := range g.MetadataHolder {
		g.CurrentID = &id
		if err := g.SetMetadata(); err != nil {
			return err
		}

		if err := g.CreateXML(); err != nil {
			return err
		}

		g.CurrentID = nil
	}

	return nil
}

func (g *Generator) setGeneralInfo() error {
	entry, err := g.CurrentEntry()
	if err != nil {
		return err
	}

	config := entry.Config

	entry.Metadata = iso1911x.ISO19110{
		XmlnsGfc:          "http://www.isotc211.org/2005/gfc",
		XmlnsGco:          "http://www.isotc211.org/2005/gco",
		XmlnsGmd:          "http://www.isotc211.org/2005/gmd",
		XmlnsGmx:          "http://www.isotc211.org/2005/gmx",
		XmlnsXsi:          "http://www.w3.org/2001/XMLSchema-instance",
		XmlnsXlink:        "http://www.w3.org/1999/xlink",
		XsiSchemaLocation: "http://www.isotc211.org/2005/gfc http://www.isotc211.org/2005/gfc/gfc.xsd",
		Uuid:              config.ID,

		Name: iso1911x.CharacterStringTag{
			CharacterString: config.Name,
		},
		VersionNumber: iso1911x.CharacterStringTag{
			CharacterString: config.VersionNumber,
		},
		VersionDate: iso1911x.DateTag{
			Date: config.VersionDate,
		},
	}

	if config.Scope != nil {
		entry.Metadata.Scope = &iso1911x.CharacterStringTag{
			CharacterString: *config.Scope,
		}
	}

	if config.FieldOfApplication != nil {
		entry.Metadata.FieldOfApplication = &iso1911x.CharacterStringTag{
			CharacterString: *config.FieldOfApplication,
		}
	}

	return nil
}

func (g *Generator) setProducerInfo() error {
	entry, err := g.CurrentEntry()
	if err != nil {
		return err
	}

	config := entry.Config

	entry.Metadata.Producer = iso1911x.ProducerTag{
		CIResponsibleParty: iso1911x.CIResponsibleParty{
			// The individual which is responsible for the feature catalogue
			IndividualName: iso1911x.AnchorOrCharacterStringTag{
				CharacterString: config.ContactIndividualName,
			},
			// The organisation which is responsible for the feature catalogue
			OrganisationName: iso1911x.AnchorOrCharacterStringTag{
				CharacterString: config.ContactOrganisationName,
			},

			Role: iso1911x.RoleTag{
				CIRoleCode: iso1911x.CodeListValueTag{
					CodeList:      "CI_RoleCode",
					CodeListValue: "pointOfContact",
				},
			},
		},
	}

	return nil
}

//nolint:funlen
func (g *Generator) setFeatureTypeInfo() error {
	entry, err := g.CurrentEntry()
	if err != nil {
		return err
	}

	config := entry.Config

	entry.Metadata.FeatureType = iso1911x.FeatureTypeTag{
		FeatureType: iso1911x.FeatureType{
			TypeName: iso1911x.TypeNameTag{
				LocalName: config.TypeName,
			},
			Definition: iso1911x.CharacterStringTag{
				CharacterString: config.Definition,
			},
		},
	}

	if config.Code != nil {
		entry.Metadata.FeatureType.FeatureType.Code = &iso1911x.CodeTag{
			Anchor: iso1911x.AnchorTag{
				Href:  config.Code.Href,
				Value: config.Code.Value,
			},
		}
	}

	if config.IsAbstract != nil {
		entry.Metadata.FeatureType.FeatureType.IsAbstract = &iso1911x.BooleanTag{
			Value: *config.IsAbstract,
		}
	}

	if len(config.Aliases) > 0 {
		entry.Metadata.FeatureType.FeatureType.Aliases = &iso1911x.Aliases{}

		for _, alias := range config.Aliases {
			val := iso1911x.LocalNameValue{
				Value: alias,
			}
			entry.Metadata.FeatureType.FeatureType.Aliases.LocalNameValues = append(
				entry.Metadata.FeatureType.FeatureType.Aliases.LocalNameValues, val)
		}
	}

	if len(config.ConstrainedBy) > 0 {
		entry.Metadata.FeatureType.FeatureType.ConstrainedBy = &iso1911x.ConstrainedBy{}
		for _, constraint := range config.ConstrainedBy {
			val := iso1911x.Constraint{
				Description: iso1911x.CharacterStringTag{
					CharacterString: constraint,
				},
			}
			entry.Metadata.FeatureType.FeatureType.ConstrainedBy.Constraints = append(
				entry.Metadata.FeatureType.FeatureType.ConstrainedBy.Constraints, val)
		}
	}

	// carrier := iso1911x.CarrierOfCharacteristicsTag{
	//	FeatureAttribute: []iso1911x.FeatureAttribute{},
	//}

	for _, attribute := range config.FeatureAttributes {
		carrier := iso1911x.CarrierOfCharacteristicsTag{
			FeatureAttribute: iso1911x.FeatureAttribute{

				FeatureType: &struct{}{},

				MemberName: iso1911x.MemberNameTag{
					LocalName: attribute.MemberName,
				},
				Definition: iso1911x.CharacterStringTag{
					CharacterString: attribute.Definition,
				},
			},
		}
		if attribute.Cardinality != nil {
			carrier.FeatureAttribute.Cardinality = &iso1911x.Cardinality{
				Multiplicity: iso1911x.Multiplicity{
					Range: iso1911x.RangeTag{
						MultiplicityRange: iso1911x.MultiplicityRange{
							Lower: iso1911x.LowerTag{
								Value: attribute.Cardinality.Lower,
							},
							Upper: iso1911x.UnlimitedIntegerHolder{
								Unlimited: iso1911x.UnlimitedInteger{
									Value: &attribute.Cardinality.Upper,
								},
							},
						},
					},
				},
			}
		}

		if attribute.ValueType != nil {
			carrier.FeatureAttribute.ValueType = &iso1911x.ValueTypeTag{
				TypeName: iso1911x.TypeName{
					AName: iso1911x.CharacterStringTag{
						CharacterString: *attribute.ValueType,
					},
				},
			}
		}

		if len(attribute.ListedValues) > 0 {
			carrier.FeatureAttribute.ListedValues = []iso1911x.ListedValue{}

			for _, listedValue := range attribute.ListedValues {
				val := iso1911x.ListedValue{
					FCListedValues: iso1911x.FCListedValue{
						Label: iso1911x.CharacterStringTag{
							CharacterString: listedValue.Label,
						},
						Code: iso1911x.CharacterStringTag{
							CharacterString: listedValue.Code,
						},
						Definition: iso1911x.CharacterStringTag{
							CharacterString: listedValue.Definition,
						},
					},
				}
				carrier.FeatureAttribute.ListedValues = append(
					carrier.FeatureAttribute.ListedValues, val)
			}
		}

		entry.Metadata.FeatureType.FeatureType.CarrierOfCharacteristics = append(
			entry.Metadata.FeatureType.FeatureType.CarrierOfCharacteristics, carrier)
	}

	return nil
}
