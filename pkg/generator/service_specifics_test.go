package generator

import (
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
		0: {filename: "regular.yaml", expectedValid: true, expectedValidationErrors: nil},
		1: {filename: "hvd_simple.yaml", expectedValid: true, expectedValidationErrors: nil},
		2: {filename: "hvd_complex.yaml", expectedValid: true, expectedValidationErrors: nil},
		3: {filename: "inspire.yaml", expectedValid: true, expectedValidationErrors: nil},
		4: {
			filename:                 "inspire_hvd_complex.yaml",
			expectedValid:            true,
			expectedValidationErrors: nil,
		},
		5: {filename: "oaf.yaml", expectedValid: true, expectedValidationErrors: nil},
		6: {filename: "oat.yaml", expectedValid: true, expectedValidationErrors: nil},
		7: {filename: "regular.json", expectedValid: true, expectedValidationErrors: nil},

		// Invalid specifics
		8: {
			filename:      "invalid_empty_values.yaml",
			expectedValid: false,
			expectedValidationErrors: []string{
				"id is required",
				"title is required",
				"contactEmail is required",
			},
		},
		9: {
			filename:                 "invalid_id_not_uuid.yaml",
			expectedValid:            false,
			expectedValidationErrors: []string{"id is not a valid UUID"},
		},
		10: {
			filename:                 "invalid_id_duplicates.yaml",
			expectedValid:            false,
			expectedValidationErrors: []string{"id is duplicate"},
		},
		11: {
			filename:      "invalid_inspire_type_no_themes.yaml",
			expectedValid: false,
			expectedValidationErrors: []string{
				"inspireThemes are required when inspireType is set",
			},
		},
	}

	for _, test := range tests {
		var serviceSpecifics ServiceSpecifics

		err := serviceSpecifics.LoadFromYAML(inputPath + test.filename)
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
