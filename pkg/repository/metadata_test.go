package repository

import (
	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
	"github.com/stretchr/testify/assert"
	"testing"
)

const cswHost = "https://nationaalgeoregister.nl/"
const cswPath = "/geonetwork/srv/dut/csw"

func TestMetadataRepository_GetDatasetMetadata(t *testing.T) {
	logPrefix := "UNITTEST_GetDatasetMetadata"

	mr, err := NewMetadataRepository(cswHost, cswPath)
	assert.Nil(t, err)

	metadataRecords, err := mr.GetDatasetMetadata(logPrefix, 123)
	assert.Nil(t, err)
	assert.Len(t, metadataRecords, 123)
	for _, metadataRecord := range metadataRecords {
		assert.NotNil(t, metadataRecord.MetadataId)
	}
}

func TestMetadataRepository_GetDatasetMetadataById(t *testing.T) {
	logPrefix := "UNITTEST_GetDatasetMetadataById"

	type args struct {
		id        string
		logPrefix string
	}

	tests := []struct {
		name                 string
		args                 args
		wantErr              bool
		wantTitle            string
		wantAbstractContains string
		wantContactName      string
		wantContactEmail     string
		wantKeywords         []string
		wantLicenceUrl       string
		wantThumbnailUrl     string
		wantInspireVariant   inspire.InspireVariant
		wantInspireThemes    []string
		wantHVDCategories    []hvd.HVDCategory
	}{
		{
			name: "GetDatasetMetadataById for BAG AsIs",
			args: args{
				id:        "aa3b5e6e-7baa-40c0-8972-3353e927ec2f",
				logPrefix: logPrefix,
			},
			wantErr:              false,
			wantTitle:            "Basisregistratie Adressen en gebouwen (BAG)",
			wantAbstractContains: "De gegevens bestaan uit BAG-panden en een deelselectie van BAG-gegevens van deze panden en de zich daarin bevindende verblijfsobjecten.",
			wantContactName:      "Klantcontactcenter",
			wantContactEmail:     "bag@kadaster.nl",
			wantKeywords:         []string{"Basisregistraties Adressen en Gebouwen", "BAG", "adres", "gebouw", "pand", "verblijfsobject", "ligplaats", "standplaats", "nummeraanduiding", "woonplaats", "basisset NOVEX", "Adressen", "Gebouwen", "Nationaal", "HVD", "Geospatiale data"},
			wantLicenceUrl:       "http://creativecommons.org/publicdomain/mark/1.0/deed.nl",
			wantThumbnailUrl:     "https://github.com/kadaster/metagegevens-voorbeelden/raw/master/Terugmeldingen_BAG.jpg",
			wantInspireVariant:   inspire.AsIs,
			wantInspireThemes:    []string{"ad", "bu"},
			wantHVDCategories:    []hvd.HVDCategory{{Id: "c_ac64a52d"}},
		},
		{
			name: "GetDatasetMetadataById for BAG Gebouwen Harmonized",
			args: args{
				id:        "b4ae622c-6201-49d8-bd2e-f7fce9206a1e",
				logPrefix: logPrefix,
			},
			wantErr:              false,
			wantTitle:            "Gebouwen - Buildings (INSPIRE geharmoniseerd)",
			wantAbstractContains: "INSPIRE Gebouwen (Buildings)",
			wantContactName:      "Klantcontactcenter",
			wantContactEmail:     "kcc@kadaster.nl",
			wantKeywords:         []string{"Gebouwen", "gebouwen", "buildings", "Nationaal", "HVD", "Geospatiale data"},
			wantLicenceUrl:       "http://creativecommons.org/publicdomain/zero/1.0/deed.nl",
			wantThumbnailUrl:     "",
			wantInspireVariant:   inspire.Harmonized,
			wantInspireThemes:    []string{"bu"},
			wantHVDCategories:    []hvd.HVDCategory{{Id: "c_ac64a52d"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mr, err := NewMetadataRepository(cswHost, cswPath)
			assert.Nil(t, err)

			metadataRecord, err := mr.GetDatasetMetadataById(tt.args.id, tt.args.logPrefix)
			if !tt.wantErr {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.wantTitle, metadataRecord.Title)
			assert.Contains(t, metadataRecord.Abstract, tt.wantAbstractContains)
			assert.Equal(t, tt.wantContactName, metadataRecord.ContactName)
			assert.Equal(t, tt.wantContactEmail, metadataRecord.ContactEmail)
			assert.Equal(t, tt.wantKeywords, metadataRecord.Keywords)
			assert.Equal(t, tt.wantLicenceUrl, metadataRecord.LicenceUrl)
			assert.Equal(t, &tt.wantThumbnailUrl, metadataRecord.ThumbnailUrl)
			assert.Equal(t, tt.wantInspireVariant, *metadataRecord.InspireVariant)
			assert.Equal(t, tt.wantInspireThemes, metadataRecord.InspireThemes)
			assert.Equal(t, tt.wantHVDCategories, metadataRecord.HVDCategories)
		})
	}
}

func TestMetadataRepository_HarvestSummaryRecords(t *testing.T) {
	logPrefix := "UNITTEST_HarvestSummaryRecords"

	mr, err := NewMetadataRepository(cswHost, cswPath)
	assert.Nil(t, err)
	records, err := mr.HarvestSummaryRecords(csw.Dataset, 50, logPrefix)

	assert.Nil(t, err)

	for _, record := range records {
		assert.NotNil(t, record.Title)
		assert.NotNil(t, record.Identifier)
	}
}

func TestMetadataRepository_HarvestMDMetadata(t *testing.T) {
	logPrefix := "UNITTEST_HarvestMDMetadata"

	mr, err := NewMetadataRepository(cswHost, cswPath)
	assert.Nil(t, err)
	records, err := mr.HarvestMDMetadata(csw.Dataset, 50, logPrefix)

	assert.Nil(t, err)
	assert.Len(t, records, 50)
	for _, record := range records {
		assert.NotNil(t, record.Title)
		assert.NotNil(t, record.MetadataId)
	}
}
