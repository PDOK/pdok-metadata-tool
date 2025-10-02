package repository

import (
	"github.com/pdok/pdok-metadata-tool/pkg/client"
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestMetadataRepository_GetDatasetMetadataById1(t *testing.T) {

	mockedNGRServer := preTestSetup()
	mr, _ := NewMetadataRepository("", "")
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
		wantMetadataId     string
		wantSourceId       string
		wantTitle          string
		wantAbstract       *string
		wantContactName    *string
		wantContactEmail   *string
		wantKeywords       []string
		wantLicenceUrl     *string
		wantThumbnailUrl   *string
		wantInspireVariant *inspire.InspireVariant
		wantInspireThemes  []string
		wantHVDCategories  []hvd.HVDCategory
	}{
		{
			name: "Voorbeeld metadata Dataset",
			args: args{
				id: "C2DFBDBC-5092-11E0-BA8E-B62DE0D72085",
			},
			wantMetadataId:   "C2DFBDBC-5092-11E0-BA8E-B62DE0D72085",
			wantSourceId:     "C2DFBDBC-5092-11E0-BA8E-B62DE0D72085",
			wantTitle:        "Naam van de dataset (*)",
			wantAbstract:     ptr("Samenvatting (*)"),
			wantContactName:  ptr("persoon verantwoordelijk voor de dataset"),
			wantContactEmail: ptr("Email@organisatie.nl"),
			wantKeywords: []string{
				"Beschermde gebieden",
				"Habitats en biotopen",
				"Nationaal",
				"Verspreidingsgebied van habitattypen (Habitatrichtlijn)",
				"Trefwoorden uit een andere thesaurus",
				"Trefwoord zonder thesaurus",
				"Tweede trefwoord zonder thesaurus",
			},
			wantLicenceUrl:     ptr("https://creativecommons.org/publicdomain/mark/*/deed.nl"),
			wantThumbnailUrl:   ptr("URL naar voorbeeldweergave van de dataset"),
			wantInspireVariant: ptr(inspire.Harmonized),
			wantInspireThemes:  []string{"ps", "hb"},
			wantHVDCategories:  nil,
		},
		{
			name: "Invasieve Exoten Dataset",
			args: args{
				id: "3703b249-a0eb-484e-ba7a-10e31a55bcec",
			},
			wantMetadataId:     "3703b249-a0eb-484e-ba7a-10e31a55bcec",
			wantSourceId:       "3703b249-a0eb-484e-ba7a-10e31a55bcec",
			wantTitle:          "Invasieve Exoten (INSPIRE Geharmoniseerd)",
			wantInspireVariant: ptr(inspire.Harmonized),
			wantInspireThemes:  []string{"sd"},
			wantHVDCategories:  []hvd.HVDCategory{{Id: "c_dd313021"}},
		},
		{
			name: "Waterschappen Hydrografie Dataset",
			args: args{
				id: "07575774-57a1-4419-bab4-6c88fdeb02b2",
			},
			wantMetadataId:     "07575774-57a1-4419-bab4-6c88fdeb02b2",
			wantSourceId:       "07575774-57a1-4419-bab4-6c88fdeb02b2",
			wantTitle:          "Waterschappen Hydrografie INSPIRE (geharmoniseerd)",
			wantInspireVariant: ptr(inspire.Harmonized),
			wantInspireThemes:  []string{"hy"},
			wantHVDCategories:  []hvd.HVDCategory{{Id: "c_dd313021"}},
		},
		{
			name: "Wetlands Dataset",
			args: args{
				id: "19165027-a13a-4c19-9013-ec1fd191019d",
			},
			wantMetadataId:     "19165027-a13a-4c19-9013-ec1fd191019d",
			wantSourceId:       "19165027-a13a-4c19-9013-ec1fd191019d",
			wantTitle:          "Wetlands (INSPIRE Geharmoniseerd)",
			wantThumbnailUrl:   ptr("https://geodata.nationaalgeoregister.nl/wetlands/ows?LAYERS=wetlands&TRANSPARENT=true&FORMAT=image%2Fpng&SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&STYLES=&EXCEPTIONS=application%2Fvnd.ogc.se_inimage&SRS=EPSG%3A28992&BBOX=-42621.76,303655.36,446379.2,686856.64&WIDTH=284&HEIGHT=223"),
			wantInspireVariant: ptr(inspire.Harmonized),
			wantInspireThemes:  []string{"ps"},
			wantHVDCategories:  []hvd.HVDCategory{{Id: "c_dd313021"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			metadataRecord, err := mr.GetDatasetMetadataById(tt.args.id)
			assert.Nil(t, err)

			assert.Equal(t, tt.wantMetadataId, metadataRecord.MetadataId)
			assert.Equal(t, tt.wantSourceId, metadataRecord.SourceId)
			assert.Equal(t, tt.wantTitle, metadataRecord.Title)

			if tt.wantAbstract != nil {
				assert.Equal(t, *tt.wantAbstract, metadataRecord.Abstract)
			}
			if tt.wantContactName != nil {
				assert.Equal(t, *tt.wantContactName, metadataRecord.ContactName)
			}
			if tt.wantContactEmail != nil {
				assert.Equal(t, *tt.wantContactEmail, metadataRecord.ContactEmail)
			}
			if tt.wantKeywords != nil {
				assert.Equal(t, tt.wantKeywords, metadataRecord.Keywords)
			}
			if tt.wantLicenceUrl != nil {
				assert.Equal(t, *tt.wantLicenceUrl, metadataRecord.LicenceUrl)
			}
			assert.Equal(t, tt.wantThumbnailUrl, metadataRecord.ThumbnailUrl)
			assert.Equal(t, tt.wantInspireVariant, metadataRecord.InspireVariant)
			assert.Equal(t, tt.wantInspireThemes, metadataRecord.InspireThemes)
			assert.Equal(t, tt.wantHVDCategories, metadataRecord.HVDCategories)

		})
	}
}

func TestMetadataRepository_SearchDatasetMetadata(t *testing.T) {
	mockedNGRServer := preTestSetup()
	mr, err := NewMetadataRepository("", "")
	if mr == nil {
		assert.FailNow(t, "NewMetadataRepository is nil")
		return
	}
	mr.CswClient = getCswClient(t, mockedNGRServer)

	title := "ataset titl"
	summaryRecords, err := mr.SearchDatasetMetadata(&title, nil)

	assert.Nil(t, err)
	assert.NotNil(t, summaryRecords)

	for _, summaryRecord := range summaryRecords {
		assert.Contains(t, summaryRecord.Title, title)
		assert.Equal(t, summaryRecord.Type, "dataset")
	}

}

func getCswClient(t *testing.T, mockedNGRServer *httptest.Server) *client.CswClient {
	hostURL, err := url.Parse(mockedNGRServer.URL)
	if err != nil {
		t.Fatalf("Failed to parse mocked NGR server URL: %v", err)
	}
	cswClient := client.NewCswClient(hostURL)
	return &cswClient
}

func ptr[T any](v T) *T {
	return &v
}
