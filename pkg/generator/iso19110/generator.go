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
		XsiSchemaLocation: "http://www.isotc211.org/2005/gfc/gfc.xsd",

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
			// The organisation which is responsible for the feature catalogue
			OrganisationName: iso1911x.OrganisationNameTag{
				CharacterString: &config.ContactOrganisationName,
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
			Code: &iso1911x.CodeTag{
				Anchor: iso1911x.AnchorTag{
					Href:  config.Code.Href,
					Value: config.Code.Value,
				},
			},
			Definition: iso1911x.CharacterStringTag{
				CharacterString: config.Definition,
			},
		},
	}

	carrier := iso1911x.CarrierOfCharacteristicsTag{
		FeatureAttributes: []iso1911x.FeatureAttribute{},
	}

	for _, attribute := range config.FeatureAttributes {
		attr := iso1911x.FeatureAttribute{
			MemberName: iso1911x.MemberNameTag{
				LocalName: attribute.MemberName,
			},
			Definition: iso1911x.CharacterStringTag{
				CharacterString: attribute.Definition,
			},
		}
		carrier.FeatureAttributes = append(carrier.FeatureAttributes, attr)
	}

	entry.Metadata.FeatureType.FeatureType.CarrierOfCharacteristics = carrier

	return nil
}
