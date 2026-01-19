package iso1911x

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const serviceMetadataStandard = "ISO19119"
const datasetMetadataStandard = "ISO19115"

func TestMDMetadata_GetKeywords(t *testing.T) {
	tests := []struct {
		name         string
		standard     string
		filename     string
		wantKeywords []string
	}{
		{
			name:     "Service INSPIRE and HVD",
			standard: serviceMetadataStandard,
			filename: "39d03482-fef0-4706-8f66-16ffb2617155.xml",
			wantKeywords: []string{
				"Kwetsbaar gebied",
				"Richtlijn 91/271/EEG",
				"31991L0271",
				"Agglomeraties",
			},
		},
		{
			name:     "Dataset INSPIRE and HVD",
			standard: datasetMetadataStandard,
			filename: "5951efa2-1ff3-4763-a966-a2f5497679ee.xml",
			wantKeywords: []string{
				"vervoersnetwerken",
				"waterwegen",
				"transport networks",
				"water",
				"transport",
				"haven",
				"veerverbinding",
				"Nationaal",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join("..", "..", "..", "examples", tt.standard, tt.filename)
			md := loadMDMetadataFromXML(t, filename)

			if gotKeywords := md.GetKeywords(); !reflect.DeepEqual(gotKeywords, tt.wantKeywords) {
				t.Errorf("GetKeywords() = %v, want %v", gotKeywords, tt.wantKeywords)
			}
		})
	}
}

func loadMDMetadataFromXML(t *testing.T, path string) MDMetadata {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	var md MDMetadata
	//nolint:musttag
	if err := xml.Unmarshal(data, &md); err != nil {
		t.Fatalf("unmarshal %s: %v", path, err)
	}

	return md
}

func TestMDMetadata_GetLicenseURL(t *testing.T) {
	tests := []struct {
		name           string
		standard       string
		filename       string
		wantLicenseURL string
	}{
		{
			name:           "Service Creative Commons license",
			standard:       serviceMetadataStandard,
			filename:       "39d03482-fef0-4706-8f66-16ffb2617155.xml",
			wantLicenseURL: "https://creativecommons.org/publicdomain/zero/1.0/deed.nl",
		},
		{
			name:           "Dataset Creative Commons license",
			standard:       datasetMetadataStandard,
			filename:       "5951efa2-1ff3-4763-a966-a2f5497679ee.xml",
			wantLicenseURL: "http://creativecommons.org/publicdomain/mark/1.0/deed.nl",
		},
		{
			name:           "Service Geo Gedeeld license",
			standard:       serviceMetadataStandard,
			filename:       "392e6a4e-5274-11ea-954f-080027325297.xml",
			wantLicenseURL: "https://www.routedatabank.nl/uitleveringsbeleid/",
		},
		{
			name:           "Dataset Geo Gedeeld license",
			standard:       datasetMetadataStandard,
			filename:       "25d77eb3-c4f6-4e6a-b974-8a93a1ace20a.xml",
			wantLicenseURL: "https://www.routedatabank.nl/uitleveringsbeleid/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join("..", "..", "..", "examples", tt.standard, tt.filename)
			md := loadMDMetadataFromXML(t, filename)

			if gotLicenseURL := md.GetLicenseURL(); !reflect.DeepEqual(
				gotLicenseURL,
				tt.wantLicenseURL,
			) {
				t.Errorf("GetLicenseURL() = %v, want %v", gotLicenseURL, tt.wantLicenseURL)
			}
		})
	}
}
