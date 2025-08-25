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

func TestMetadataRepository_GetDatasetMetadataById(t *testing.T) {
	mockedNGRServer := preTestSetup()
	mr, err := NewMetadataRepository("", "")
	if mr == nil {
		assert.FailNow(t, "NewMetadataRepository is nil")
		return
	}

	mr.CswClient = getCswClient(t, mockedNGRServer)

	datasetId := "C2DFBDBC-5092-11E0-BA8E-B62DE0D72085"
	logPrefix := "UNITTEST_GetDatasetMetadataById"

	metadataRecord, err := mr.GetDatasetMetadataById(datasetId, logPrefix)
	assert.Nil(t, err)

	assert.Equal(t, "Naam van de dataset (*)", metadataRecord.Title)
	assert.Equal(t, "Samenvatting (*)", metadataRecord.Abstract)
	assert.Equal(t, "persoon verantwoordelijk voor de dataset", metadataRecord.ContactName)
	assert.Equal(t, "Email@organisatie.nl", metadataRecord.ContactEmail)
	assert.Equal(t, []string{
		"Beschermde gebieden",
		"Habitats en biotopen",
		"Nationaal",
		"Verspreidingsgebied van habitattypen (Habitatrichtlijn)",
		"Trefwoorden uit een andere thesaurus",
		"Trefwoord zonder thesaurus",
		"Tweede trefwoord zonder thesaurus",
	}, metadataRecord.Keywords)
	assert.Equal(t, "https://creativecommons.org/publicdomain/mark/*/deed.nl", metadataRecord.LicenceUrl)
	assert.Equal(t, "URL naar voorbeeldweergave van de dataset", *metadataRecord.ThumbnailUrl)

	expectedInspireVariant := (*inspire.InspireVariant)(nil)
	assert.Equal(t, expectedInspireVariant, metadataRecord.InspireVariant)

	var expectedKeywords []string = nil
	assert.Equal(t, expectedKeywords, metadataRecord.InspireThemes)

	var expectedHVDCategories []hvd.HVDCategory = nil
	assert.Equal(t, expectedHVDCategories, metadataRecord.HVDCategories)

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
	logPrefix := "UNITTEST_GetDatasetMetadataById"

	summaryRecords, err := mr.SearchDatasetMetadata(&title, nil, logPrefix)

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
