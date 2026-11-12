package iso19119

import (
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdok/pdok-metadata-tool/v2/internal/common"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/generator/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const inputPath = "testdata/input/"
const outputFolder = "testdata/output"
const expectedPath = "testdata/expected"

func TestGenerateMetadataISO19119(t *testing.T) {
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
		{
			configFileName: filepath.Join(inputPath, "geo_gedeeld.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000017.xml": "geo_gedeeld_wms.xml",
				"00000000-0000-0000-0000-000000000018.xml": "geo_gedeeld_wfs.xml",
			},
		},
		{
			configFileName: filepath.Join(inputPath, "oam.yaml"),
			fileOutput: map[string]string{
				"00000000-0000-0000-0000-000000000019.xml": "oam_oam.xml",
			},
		},
	}

	hvdCachePath := path.Join(common.GetProjectRoot(), common.HvdLocalRDFPath)

	for _, test := range tests {
		var serviceSpecifics ServiceSpecifics

		err := serviceSpecifics.LoadFromYamlOrJson(test.configFileName)
		require.NoError(t, err)

		err = serviceSpecifics.Validate()
		require.NoError(t, err)

		generator, err := NewGenerator(serviceSpecifics, outputFolder, nil, &hvdCachePath)
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

func TestGenerateMetadataISO19119DeduplicatesKeywords(t *testing.T) {
	tests := []struct {
		configFileName string
		id             string
		reserved       string
	}{
		{filepath.Join(inputPath, "keywords_dedupe_hvd.yaml"), "00000000-0000-0000-0000-000000000020", "Mobiliteit"},
		{filepath.Join(inputPath, "keywords_dedupe_inspire.yaml"), "00000000-0000-0000-0000-000000000021", "Beschermde gebieden"},
	}

	hvdCachePath := path.Join(common.GetProjectRoot(), common.HvdLocalRDFPath)

	for _, test := range tests {
		var serviceSpecifics ServiceSpecifics

		require.NoError(t, serviceSpecifics.LoadFromYamlOrJson(test.configFileName))
		require.NoError(t, serviceSpecifics.Validate())

		generator, err := NewGenerator(serviceSpecifics, outputFolder, nil, &hvdCachePath)
		require.NoError(t, err)

		generated, err := generator.GenerateAsStrings()
		require.NoError(t, err)

		xml := generated[test.id]
		require.NotEmpty(t, xml)

		// Plain keywords are CharacterStrings, INSPIRE/HVD keywords are Anchors
		assert.Equal(t, 1, strings.Count(xml, "<gco:CharacterString>A</gco:CharacterString>"), "plain keyword A once")
		assert.Equal(t, 1, strings.Count(xml, "<gco:CharacterString>B</gco:CharacterString>"), "plain keyword B once")
		assert.NotContains(t, xml, "<gco:CharacterString>"+test.reserved+"</gco:CharacterString>")
		assert.NotContains(t, xml, "<gco:CharacterString>ps</gco:CharacterString>")
		assert.NotContains(t, xml, "<gco:CharacterString>HVD</gco:CharacterString>")
		assert.Equal(t, 1, strings.Count(xml, ">"+test.reserved+"</gmx:Anchor>"), "category keyword still present once")
	}
}
