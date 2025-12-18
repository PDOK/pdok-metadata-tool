package repository

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pdok/pdok-metadata-tool/pkg/client"
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
	"github.com/pdok/pdok-metadata-tool/pkg/model/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetadataRepository_GetDatasetMetadataById(t *testing.T) {
	mockedNGRServer := preTestSetup()

	mr, _ := NewMetadataRepository("")
	if mr == nil {
		assert.FailNow(t, "NewMetadataRepository is nil")

		return
	}

	mr.CswClient = getCswClient(t, mockedNGRServer)

	type args struct {
		id string
	}

	tests := []struct {
		name               string
		args               args
		wantMetadataID     string
		wantSourceID       string
		wantTitle          string
		wantContactURL     string
		wantAbstract       string
		wantContactName    string
		wantContactEmail   string
		wantKeywords       []string
		wantLicenceURL     string
		wantUseLimitation  string
		wantThumbnailURL   string
		wantInspireVariant inspire.InspireVariant
		wantInspireThemes  []string
		wantHVDCategories  []hvd.HVDCategory
		wantBoundingBox    *metadata.BoundingBox
	}{
		{
			name: "Voorbeeld metadata Dataset",
			args: args{
				id: "C2DFBDBC-5092-11E0-BA8E-B62DE0D72085",
			},
			wantMetadataID:   "C2DFBDBC-5092-11E0-BA8E-B62DE0D72085",
			wantSourceID:     "1234",
			wantTitle:        "Naam van de dataset (*)",
			wantAbstract:     "Samenvatting (*)",
			wantContactName:  "persoon verantwoordelijk voor de dataset",
			wantContactEmail: "Email@organisatie.nl",
			wantContactURL:   "https://www.geonovum.nl/",
			wantKeywords: []string{
				"Nationaal", //nolint:misspell
				"Verspreidingsgebied van habitattypen (Habitatrichtlijn)",
				"Trefwoorden uit een andere thesaurus",
				"Trefwoord zonder thesaurus",
				"Tweede trefwoord zonder thesaurus",
			},
			wantLicenceURL:     "https://creativecommons.org/publicdomain/mark/*/deed.nl",
			wantUseLimitation:  "Gebruiksbeperkingen (*), Toepassingen waarvoor de data niet geschikt is.",
			wantThumbnailURL:   "URL naar voorbeeldweergave van de dataset",
			wantInspireVariant: inspire.AsIs,
			wantInspireThemes:  []string{"ps", "hb"},
			wantHVDCategories:  nil,
			wantBoundingBox: &metadata.BoundingBox{
				WestBoundLongitude: "3.37087",
				EastBoundLongitude: "7.21097",
				SouthBoundLatitude: "50.7539",
				NorthBoundLatitude: "53.4658",
			},
		},
		{
			name: "Invasieve Exoten Dataset",
			args: args{
				id: "3703b249-a0eb-484e-ba7a-10e31a55bcec",
			},
			wantMetadataID:     "3703b249-a0eb-484e-ba7a-10e31a55bcec",
			wantSourceID:       "3703b249-a0eb-484e-ba7a-10e31a55bcec",
			wantTitle:          "Invasieve Exoten (INSPIRE Geharmoniseerd)",
			wantInspireVariant: inspire.Harmonised,
			wantInspireThemes:  []string{"sd"},
			wantHVDCategories: []hvd.HVDCategory{
				{ID: "c_dd313021", LabelDutch: "Aardobservatie en milieu"},
			},
			wantBoundingBox: &metadata.BoundingBox{
				WestBoundLongitude: "-3.5879",
				EastBoundLongitude: "13.5757",
				SouthBoundLatitude: "49.1241",
				NorthBoundLatitude: "54.9991",
			},
		},
		{
			name: "Waterschappen Hydrografie Dataset",
			args: args{
				id: "07575774-57a1-4419-bab4-6c88fdeb02b2",
			},
			wantMetadataID:     "07575774-57a1-4419-bab4-6c88fdeb02b2",
			wantSourceID:       "07575774-57a1-4419-bab4-6c88fdeb02b2",
			wantTitle:          "Waterschappen Hydrografie INSPIRE (geharmoniseerd)",
			wantUseLimitation:  "Niet te gebruiken voor navigatie. Niet te gebruiken voor juridische bewijsvoering.",
			wantInspireVariant: inspire.Harmonised,
			wantInspireThemes:  []string{"hy"},
			wantHVDCategories: []hvd.HVDCategory{
				{ID: "c_dd313021", LabelDutch: "Aardobservatie en milieu"},
			},
			wantBoundingBox: &metadata.BoundingBox{
				WestBoundLongitude: "2.65899516",
				EastBoundLongitude: "7.83057492",
				SouthBoundLatitude: "50.58707771",
				NorthBoundLatitude: "53.73639341",
			},
		},
		{
			name: "Wetlands Dataset",
			args: args{
				id: "19165027-a13a-4c19-9013-ec1fd191019d",
			},
			wantMetadataID:     "19165027-a13a-4c19-9013-ec1fd191019d",
			wantSourceID:       "84487381-957b-4bd6-a9c9-47c6b6037223",
			wantTitle:          "Wetlands (INSPIRE Geharmoniseerd)",
			wantUseLimitation:  "Geen gebruiksbeperkingen",
			wantThumbnailURL:   "https://geodata.nationaalgeoregister.nl/wetlands/ows?LAYERS=wetlands&TRANSPARENT=true&FORMAT=image%2Fpng&SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&STYLES=&EXCEPTIONS=application%2Fvnd.ogc.se_inimage&SRS=EPSG%3A28992&BBOX=-42621.76,303655.36,446379.2,686856.64&WIDTH=284&HEIGHT=223",
			wantInspireVariant: inspire.Harmonised,
			wantInspireThemes:  []string{"ps"},
			wantHVDCategories: []hvd.HVDCategory{
				{ID: "c_dd313021", LabelDutch: "Aardobservatie en milieu"},
			},
			wantBoundingBox: &metadata.BoundingBox{
				WestBoundLongitude: "2.1339",
				EastBoundLongitude: "8.16",
				SouthBoundLatitude: "50.5591",
				NorthBoundLatitude: "53.7509",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadataRecord, err := mr.GetDatasetMetadataByID(tt.args.id)
			require.NoError(t, err)

			assert.Equal(t, tt.wantMetadataID, metadataRecord.MetadataID)
			assert.Equal(t, tt.wantSourceID, metadataRecord.SourceID)
			assert.Equal(t, tt.wantTitle, metadataRecord.Title)

			if tt.wantAbstract != "" {
				assert.Equal(t, tt.wantAbstract, metadataRecord.Abstract)
			}

			if tt.wantContactName != "" {
				assert.Equal(t, tt.wantContactName, metadataRecord.ContactName)
			}

			if tt.wantContactEmail != "" {
				assert.Equal(t, tt.wantContactEmail, metadataRecord.ContactEmail)
			}

			if tt.wantContactURL != "" {
				assert.Equal(t, tt.wantContactURL, metadataRecord.ContactURL)
			}

			if tt.wantKeywords != nil {
				assert.Equal(t, tt.wantKeywords, metadataRecord.Keywords)
			}

			if tt.wantLicenceURL != "" {
				assert.Equal(t, tt.wantLicenceURL, metadataRecord.LicenceURL)
			}

			if tt.wantUseLimitation != "" {
				assert.Equal(t, tt.wantUseLimitation, metadataRecord.UseLimitation)
			}

			assert.Equal(t, tt.wantThumbnailURL, metadataRecord.ThumbnailURL)
			assert.Equal(t, tt.wantInspireVariant, metadataRecord.InspireVariant)
			assert.Equal(t, tt.wantInspireThemes, metadataRecord.InspireThemes)
			assert.Equal(t, tt.wantHVDCategories, metadataRecord.HVDCategories)
			assert.Equal(t, tt.wantBoundingBox, metadataRecord.BoundingBox)
		})
	}
}

func TestMetadataRepository_SearchDatasetMetadata(t *testing.T) {
	mockedNGRServer := preTestSetup()

	mr, _ := NewMetadataRepository("")
	if mr == nil {
		assert.FailNow(t, "NewMetadataRepository is nil")

		return
	}

	mr.CswClient = getCswClient(t, mockedNGRServer)

	title := "ataset titl"
	summaryRecords, err := mr.SearchDatasetMetadata(&title, nil)

	require.NoError(t, err)
	assert.NotNil(t, summaryRecords)

	for _, summaryRecord := range summaryRecords {
		assert.Contains(t, summaryRecord.Title, title)
		assert.Equal(t, "dataset", summaryRecord.Type)
	}
}

func getCswClient(t *testing.T, mockedNGRServer *httptest.Server) *client.CswClient {
	t.Helper()

	hostURL, err := url.Parse(mockedNGRServer.URL)
	if err != nil {
		t.Fatalf("Failed to parse mocked NGR server URL: %v", err)
	}

	cswClient := client.NewCswClient(hostURL)

	return &cswClient
}
