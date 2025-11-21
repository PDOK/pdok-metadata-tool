package generator

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ucarion/c14n"
)

const inputPath = "testdata/input/"
const outputFolder = "testdata/output"
const expectedPath = "testdata/expected"

func TestGenerateMetadata(t *testing.T) {
	var tests = []struct {
		configFileName string
		fileOutput     map[string]string
	}{
		{
			configFileName: filepath.Join(inputPath, "regular.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000001.xml": "regular_wfs.xml",
				"00000000-0000-0000-0000-000000000002.xml": "regular_wms.xml",
			},
		},
		{
			configFileName: filepath.Join(inputPath, "hvd_simple.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000003.xml": "hvd_simple_wms.xml",
			},
		},
		{
			configFileName: filepath.Join(inputPath, "hvd_complex.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000004.xml": "hvd_complex_wms.xml",
			},
		},

		{
			configFileName: filepath.Join(inputPath, "oaf.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000005.xml": "oaf_oaf.xml",
			},
		},
		{
			configFileName: filepath.Join(inputPath, "oat.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000006.xml": "oat_oat.xml",
			},
		},
		{
			configFileName: filepath.Join(inputPath, "inspire_asis.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000007.xml": "inspire_asis_wms.xml",
				"00000000-0000-0000-0000-000000000008.xml": "inspire_asis_wfs.xml",
				"00000000-0000-0000-0000-000000000009.xml": "inspire_asis_atom.xml",
			},
		},
		{
			configFileName: filepath.Join(inputPath, "inspire_harmonised.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000010.xml": "inspire_harmonised_wms.xml",
				"00000000-0000-0000-0000-000000000011.xml": "inspire_harmonised_wfs.xml",
				"00000000-0000-0000-0000-000000000012.xml": "inspire_harmonised_atom.xml",
			},
		},
		{
			configFileName: filepath.Join(inputPath, "inspire_hvd_complex.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000013.xml": "inspire_hvd_complex_atom.xml",
				"00000000-0000-0000-0000-000000000014.xml": "inspire_hvd_complex_wfs_invocable.xml",
				"00000000-0000-0000-0000-000000000015.xml": "inspire_hvd_complex_wfs_interoperable.xml",
				"00000000-0000-0000-0000-000000000016.xml": "inspire_hvd_complex_oaf_interoperable.xml",
			},
		},
	}

	for _, test := range tests {
		var serviceSpecifics ServiceSpecifics

		err := serviceSpecifics.LoadFromYamlOrJson(test.configFileName)
		require.NoError(t, err)

		err = serviceSpecifics.Validate()
		require.NoError(t, err)

		generator, err := NewISO19119Generator(serviceSpecifics, outputFolder, nil, nil)
		require.NoError(t, err)

		err = generator.Generate()
		require.NoError(t, err)

		for createdOutput, expectedOutput := range test.fileOutput {
			xml1, err := canonicalizeXML(filepath.Join(outputFolder, createdOutput))
			require.NoError(t, err)

			xml2, err := canonicalizeXML(filepath.Join(expectedPath, expectedOutput))
			require.NoError(t, err)

			assert.Equal(t, xml1, xml2, "Canonicalized XML files should be equal")
		}

		generatedMetadata, err := generator.GenerateAsStrings()
		require.NoError(t, err)

		// The generated metadata content has already been compared, see above
		// For GenerateAsStrings we only need to check if the numbers match
		assert.Len(t, generatedMetadata, len(test.fileOutput))

		for id, metadata := range generatedMetadata {
			assert.NotEmpty(t, metadata)

			_, ok := test.fileOutput[id+".xml"]
			assert.True(t, ok)
		}
	}
}

func canonicalizeXML(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	decoder := xml.NewDecoder(strings.NewReader(string(data)))

	canonical, err := c14n.Canonicalize(decoder)
	if err != nil {
		return "", err
	}

	return string(canonical), nil
}
