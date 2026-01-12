package iso19110

import (
	"path/filepath"
	"testing"

	"github.com/pdok/pdok-metadata-tool/v2/pkg/generator/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const inputPath = "testdata/input/"
const outputFolder = "testdata/output"
const expectedPath = "testdata/expected"

func TestGenerateMetadataISO19110(t *testing.T) {
	var tests = []struct {
		configFileName string
		fileOutput     map[string]string
	}{
		{
			configFileName: filepath.Join(inputPath, "voorbeeld_geonovum.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000001.xml": "voorbeeld_geonovum.xml",
			},
		},
		{
			configFileName: filepath.Join(inputPath, "nwb_wegen.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000002.xml": "nwb_wegen_hectopunten.xml",
			},
		},
	}

	for _, test := range tests {
		var featureCatalogueSpecifics FeatureCatalogueSpecifics

		err := featureCatalogueSpecifics.LoadFromYamlOrJson(test.configFileName)
		require.NoError(t, err)

		generator, err := NewGenerator(featureCatalogueSpecifics, outputFolder)
		require.NoError(t, err)

		err = generator.Generate()
		require.NoError(t, err)

		for createdOutput, expectedOutput := range test.fileOutput {
			xml1, err := utils.CanonicalizeXML(filepath.Join(outputFolder, createdOutput))
			require.NoError(t, err)

			xml2, err := utils.CanonicalizeXML(filepath.Join(expectedPath, expectedOutput))
			require.NoError(t, err)

			assert.Equal(t, xml1, xml2, "Canonicalized XML files should be equal")
		}
	}
}
