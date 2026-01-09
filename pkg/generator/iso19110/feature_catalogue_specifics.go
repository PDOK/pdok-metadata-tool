package iso19110

import (
	"os"

	"gopkg.in/yaml.v3"
)

// FeatureCatalogueSpecifics struct for unmarshalling the input for feature catalogue metadata generation.
type FeatureCatalogueSpecifics struct {
	Globals           GlobalConfig             `json:"globals,omitempty"           yaml:"globals,omitempty"`
	FeatureCatalogues []FeatureCatalogueConfig `json:"featureCatalogues,omitempty" yaml:"featureCatalogues,omitempty"`
}

// GlobalConfig struct for unmarshalling service specifics input.
type GlobalConfig struct {
}

type FeatureCatalogueConfig struct {
	ID                      string             `json:"id"                      yaml:"id"`
	Name                    string             `json:"name"                    yaml:"name"`
	VersionNumber           string             `json:"versionNumber"           yaml:"versionNumber"`
	VersionDate             string             `json:"versionDate"             yaml:"versionDate"`
	ContactOrganisationName string             `json:"contactOrganisationName" yaml:"contactOrganisationName"`
	ContactEmail            string             `json:"contactEmail"            yaml:"contactEmail"`
	ContactURL              string             `json:"contactUrl"              yaml:"contactUrl"`
	TypeName                string             `json:"typeName"                yaml:"typeName"`
	Code                    *CodeTag           `json:"code,omitempty"          yaml:"code,omitempty"`
	Definition              string             `json:"definition"              yaml:"definition"`
	FeatureAttributes       []FeatureAttribute `json:"featureAttributes"       yaml:"featureAttributes"`
}

type CodeTag struct {
	Href  string `json:"href"  yaml:"href"`
	Value string `json:"value" yaml:"value"`
}

type FeatureAttribute struct {
	MemberName string `json:"memberName" yaml:"memberName"`
	Definition string `json:"definition" yaml:"definition"`
}

func (c FeatureCatalogueConfig) GetID() string { return c.ID }

// LoadFromYamlOrJson unmarshalls the input for the given input file.
func (f *FeatureCatalogueSpecifics) LoadFromYamlOrJson(filename string) error {
	//nolint:gosec
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(yamlFile, f); err != nil {
		return err
	}

	return nil
}
