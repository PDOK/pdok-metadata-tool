package iso19110

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceSpecificsLoadFromYAMLAndValidate(t *testing.T) {
	var tests = []struct {
		filename                 string
		expectedValid            bool
		expectedValidationErrors []string
	}{
		// Valid specifics
		{filename: "voorbeeld_geonovum.yaml", expectedValid: true, expectedValidationErrors: nil},
		{filename: "nwb_wegen.yaml", expectedValid: true, expectedValidationErrors: nil},
		// Invalid specifics
		{
			filename:      "invalid_empty_values.yaml",
			expectedValid: false,
			expectedValidationErrors: []string{
				"id is required",
				"typeName is required",
				"definition is required",
			},
		},
		{
			filename:                 "invalid_id_not_uuid.yaml",
			expectedValid:            false,
			expectedValidationErrors: []string{"id is not a valid UUID"},
		},
		{
			filename:                 "invalid_id_duplicates.yaml",
			expectedValid:            false,
			expectedValidationErrors: []string{"id is duplicate"},
		},
		{
			filename:      "invalid_feature_attribute_empty_values.yaml",
			expectedValid: false,
			expectedValidationErrors: []string{
				"definition is required for all attributes",
				".; memberName is required for all attributes",
			},
		},
		{
			filename:      "invalid_empty_value_measurement_unit.yaml",
			expectedValid: false,
			expectedValidationErrors: []string{
				"codespace is required if a valueMeasurementUnit is set.",
				"catalogSymbol is required if a valueMeasurementUnit is set.",
				"unitDefinitionId is required if a valueMeasurementUnit is set.",
				"Identifier is required if a valueMeasurementUnit is set.",
			},
		},
	}

	for _, test := range tests {
		var fcSpecifics FeatureCatalogueSpecifics

		err := fcSpecifics.LoadFromYamlOrJson(inputPath + test.filename)
		require.NoError(t, err)

		err = fcSpecifics.Validate()
		if test.expectedValid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
			validationError := err.Error()

			for _, expectedError := range test.expectedValidationErrors {
				assert.Contains(t, validationError, expectedError)
			}
		}
	}
}
