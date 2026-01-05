package generator

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/pdok/pdok-metadata-tool/internal/common"
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

func TestGetTitle(t *testing.T) {
	globalConfig := GlobalConfig{
		OverrideableFields: OverrideableFields{
			Title: common.Ptr("Title"),
		},
	}

	globalConfigWithAtomPostfix := GlobalConfig{
		OverrideableFields: OverrideableFields{
			Title: common.Ptr("Title ATOM"),
		},
	}

	var tests = []struct {
		description   string
		serviceConfig ServiceConfig
		expectedTitle string
	}{

		{
			description: "Global title with WMS postfix",
			serviceConfig: ServiceConfig{
				Type:    "WMS",
				Globals: &globalConfig,
			},
			expectedTitle: "Title WMS",
		},
		{
			description: "Global title with WFS postfix",
			serviceConfig: ServiceConfig{
				Type:    "wFs",
				Globals: &globalConfig,
			},
			expectedTitle: "Title WFS",
		},
		{
			description: "Global title with ATOM postfix",
			serviceConfig: ServiceConfig{
				Type:    "Atom",
				Globals: &globalConfig,
			},
			expectedTitle: "Title ATOM",
		},
		{
			description: "Global title with OGC API Features postfix",
			serviceConfig: ServiceConfig{
				Type:    "OAF",
				Globals: &globalConfig,
			},
			expectedTitle: "Title OGC API Features",
		},
		{
			description: "Global title with OGC API (Vector) Tiles postfix",
			serviceConfig: ServiceConfig{
				Type:    "OAT",
				Globals: &globalConfig,
			},
			expectedTitle: "Title OGC API (Vector) Tiles",
		},
		{
			description: "Specific title on service level without postfix",
			serviceConfig: ServiceConfig{
				Type:    "WMS",
				Globals: &globalConfig,
				OverrideableFields: OverrideableFields{
					Title: common.Ptr("A specific title without postfix"),
				},
			},
			expectedTitle: "A specific title without postfix",
		},
		{
			description: "Empty string if no title info is available",
			serviceConfig: ServiceConfig{
				Type:    "WMS",
				Globals: &GlobalConfig{},
			},
			expectedTitle: "",
		},
		{
			description: "Global title with pre existing ATOM postfix",
			serviceConfig: ServiceConfig{
				Type:    "Atom",
				Globals: &globalConfigWithAtomPostfix,
			},
			expectedTitle: "Title ATOM",
		},
	}
	for _, test := range tests {
		title := test.serviceConfig.GetTitle()
		assert.Equal(t, test.expectedTitle, title)
	}
}
