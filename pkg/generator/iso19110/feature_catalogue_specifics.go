package iso19110

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
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
	ID                      string             `json:"id"                                yaml:"id"`
	Name                    string             `json:"name"                              yaml:"name"`
	VersionNumber           string             `json:"versionNumber"                     yaml:"versionNumber"`
	VersionDate             string             `json:"versionDate"                       yaml:"versionDate"`
	Scope                   *string            `json:"scope,omitempty"                   yaml:"scope,omitempty"`
	FieldOfApplication      *string            `json:"fieldOfApplication,omitempty"      yaml:"fieldOfApplication,omitempty"`
	ContactIndividualName   *string            `json:"contactIndividualName,omitempty"   yaml:"contactIndividualName,omitempty"`
	ContactOrganisationName *string            `json:"contactOrganisationName,omitempty" yaml:"contactOrganisationName,omitempty"`
	ContactEmail            *string            `json:"contactEmail,omitempty"            yaml:"contactEmail,omitempty"`
	ContactURL              *string            `json:"contactUrl,omitempty"              yaml:"contactUrl,omitempty"`
	TypeName                string             `json:"typeName"                          yaml:"typeName"`
	Code                    *CodeTag           `json:"code,omitempty"                    yaml:"code,omitempty"`
	Definition              string             `json:"definition"                        yaml:"definition"`
	IsAbstract              *bool              `json:"isAbstract,omitempty"              yaml:"isAbstract,omitempty"`
	Aliases                 []string           `json:"aliases"                           yaml:"aliases"`
	ConstrainedBy           []string           `json:"constrainedBy"                     yaml:"constrainedBy"`
	FeatureAttributes       []FeatureAttribute `json:"featureAttributes"                 yaml:"featureAttributes"`
}

type CodeTag struct {
	Href  string `json:"href"  yaml:"href"`
	Value string `json:"value" yaml:"value"`
}

type FeatureAttribute struct {
	MemberName           string                `json:"memberName"                     yaml:"memberName"`
	Definition           string                `json:"definition"                     yaml:"definition"`
	Cardinality          *Cardinality          `json:"cardinality,omitempty"          yaml:"cardinality,omitempty"`
	ValueMeasurementUnit *ValueMeasurementUnit `json:"valueMeasurementUnit,omitempty" yaml:"valueMeasurementUnit,omitempty"`
	ValueType            *string               `json:"valueType,omitempty"            yaml:"valueType,omitempty"`
	ListedValues         []ListedValue         `json:"listedValues"                   yaml:"listedValues"`
}

type Cardinality struct {
	Lower int `json:"lower" yaml:"lower"`
	Upper int `json:"upper" yaml:"upper"`
}

type ValueMeasurementUnit struct {
	UnitDefinitionId string `json:"unitDefinitionId" yaml:"unitDefinitionId"`
	Codespace        string `json:"codespace"        yaml:"codespace"`
	Identifier       string `json:"identifier"       yaml:"identifier"`
	Name             string `json:"name"             yaml:"name"`
	CatalogSymbol    string `json:"catalogSymbol"    yaml:"catalogSymbol"`
}

type ListedValue struct {
	Label      string `json:"label"      yaml:"label"`
	Code       string `json:"code"       yaml:"code"`
	Definition string `json:"definition" yaml:"definition"`
}

func (fc FeatureCatalogueConfig) GetID() string { return fc.ID }

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

// Validate the FeatureCatalogueSpecifics on a global level, also calls Validate on feature catalogue level.
func (f *FeatureCatalogueSpecifics) Validate() error {
	var validationErrors []string

	seenIDs := make(map[string]bool)

	for i, featureCatalogue := range f.FeatureCatalogues {
		// Check for duplicate ID
		if seenIDs[featureCatalogue.ID] {
			validationErrors = append(
				validationErrors,
				fmt.Sprintf("FeatureCatalogue[%d]: id is duplicate '%s'", i, featureCatalogue.ID),
			)
		} else {
			seenIDs[featureCatalogue.ID] = true
		}

		// Validate individual featureCatalogue
		if err := featureCatalogue.Validate(); err != nil {
			validationErrors = append(
				validationErrors,
				fmt.Sprintf("FeatureCatalogue[%d] (%s): %v", i, featureCatalogue.ID, err),
			)
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation failed:\n%s", strings.Join(validationErrors, "\n"))
	}

	return nil
}

// Validate the FeatureCatalogueSpecifics on service level.
func (fc FeatureCatalogueConfig) Validate() error {
	var errors []string

	if fc.ID == "" {
		errors = append(errors, "id is required")
	} else {
		if _, err := uuid.Parse(fc.ID); err != nil {
			errors = append(errors, "id is not a valid UUID: "+fc.ID)
		}
	}

	if fc.TypeName == "" {
		errors = append(errors, "typeName is required")
	}

	if fc.Definition == "" {
		errors = append(errors, "definition is required")
	}

	hasMissingMemberName := false
	hasMissingDefinition := false

	for _, attribute := range fc.FeatureAttributes {
		if attribute.MemberName == "" && !hasMissingMemberName {
			errors = append(errors, "memberName is required for all attributes")
			hasMissingMemberName = true
		}

		if attribute.Definition == "" && !hasMissingDefinition {
			errors = append(errors, "definition is required for all attributes.")
			hasMissingDefinition = true
		}

		if attribute.ValueMeasurementUnit != nil {
			vmu := attribute.ValueMeasurementUnit
			if vmu.Codespace == "" {
				errors = append(errors, "codespace is required if a valueMeasurementUnit is set.")
			}

			if vmu.Name == "" {
				errors = append(errors, "name is required if a valueMeasurementUnit is set.")
			}

			if vmu.CatalogSymbol == "" {
				errors = append(
					errors,
					"catalogSymbol is required if a valueMeasurementUnit is set.",
				)
			}

			if vmu.UnitDefinitionId == "" {
				errors = append(
					errors,
					"unitDefinitionId is required if a valueMeasurementUnit is set.",
				)
			}

			if vmu.Identifier == "" {
				errors = append(errors, "Identifier is required if a valueMeasurementUnit is set.")
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "; "))
	}

	return nil
}
