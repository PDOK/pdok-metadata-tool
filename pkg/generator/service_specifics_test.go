package generator

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceSpecificsLoadFromYAMLTestServiceSpecificsValidate(t *testing.T) {
	var tests = []struct {
		filename                 string
		expectedValid            bool
		expectedValidationErrors []string
	}{
		// Valid specifics
		{filename: "regular.yaml", expectedValid: true, expectedValidationErrors: nil},
		{filename: "hvd_simple.yaml", expectedValid: true, expectedValidationErrors: nil},
		{filename: "hvd_complex.yaml", expectedValid: true, expectedValidationErrors: nil},
		{filename: "inspire_asis.yaml", expectedValid: true, expectedValidationErrors: nil},
		{filename: "inspire_harmonised.yaml", expectedValid: true, expectedValidationErrors: nil},
		{
			filename:                 "inspire_hvd_complex.yaml",
			expectedValid:            true,
			expectedValidationErrors: nil,
		},
		{filename: "oaf.yaml", expectedValid: true, expectedValidationErrors: nil},
		{filename: "oat.yaml", expectedValid: true, expectedValidationErrors: nil},
		{filename: "regular.json", expectedValid: true, expectedValidationErrors: nil},

		// Invalid specifics
		{
			filename:      "invalid_empty_values.yaml",
			expectedValid: false,
			expectedValidationErrors: []string{
				"id is required",
				"title is required",
				"contactEmail is required",
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
			filename:      "invalid_inspire_type_no_themes.yaml",
			expectedValid: false,
			expectedValidationErrors: []string{
				"inspireThemes are required when inspireType is set",
			},
		},
	}

	for _, test := range tests {
		var serviceSpecifics ServiceSpecifics

		err := serviceSpecifics.LoadFromYamlOrJson(inputPath + test.filename)
		require.NoError(t, err)

		err = serviceSpecifics.Validate()
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

func TestParseAnnotations(t *testing.T) {
	var serviceSpecifics ServiceSpecifics

	err := serviceSpecifics.LoadFromYamlOrJson(inputPath + "annotations_minimum.json")
	require.NoError(t, err)

	out, err := json.MarshalIndent(serviceSpecifics, "", "  ")
	require.NoError(t, err)

	// Read expected JSON from file
	expectedJSONBytes, err := os.ReadFile(expectedPath + "/annotations_minimum.json")
	require.NoError(t, err)
	require.JSONEq(
		t,
		string(expectedJSONBytes),
		string(out),
		"Generated JSON does not match expected JSON",
	)
}
