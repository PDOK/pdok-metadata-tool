package codelist

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetReferenceSystemByEPSGCode(t *testing.T) {
	codelistLookupService, err := NewCodelist()
	require.NoError(t, err)

	referenceSystem, ok := codelistLookupService.GetReferenceSystemByEPSGCode("EPSG:28992")
	assert.True(t, ok)
	assert.NotNil(t, referenceSystem)
	assert.Equal(t, "Amersfoort / RD New", referenceSystem.Name)

	_, ok = codelistLookupService.GetReferenceSystemByEPSGCode("28992")
	assert.False(t, ok)
}

func TestGetINSPIREThemeLabelByURI(t *testing.T) {
	codelistLookupService, err := NewCodelist()
	require.NoError(t, err)

	themeLabel, ok := codelistLookupService.GetINSPIREThemeLabelByURI(
		"https://www.eionet.europa.eu/gemet/nl/inspire-theme/hy",
	)
	assert.True(t, ok)
	assert.NotNil(t, themeLabel)
	assert.Equal(t, "Hydrografie", *themeLabel)
}

func TestGetProtocolDetailsByProtocol(t *testing.T) {
	codelistLookupService, err := NewCodelist()
	require.NoError(t, err)

	protocol, ok := codelistLookupService.GetProtocolDetailsByProtocol("wfs")
	assert.True(t, ok)
	assert.NotNil(t, protocol)
	assert.Equal(t, "OGC:WFS", protocol.ServiceProtocol)
	assert.Equal(t, "infoFeatureAccessService", protocol.SpatialDataserviceCategory)
	assert.Equal(t, "2010-11-02", protocol.ProtocolReleaseDate)
}

func TestGetInspireServiceTypeByServiceType(t *testing.T) {
	codelistLookupService, err := NewCodelist()
	require.NoError(t, err)

	inspireServiceType, ok := codelistLookupService.GetInspireServiceTypeByServiceType("wMs")
	assert.True(t, ok)
	assert.NotNil(t, inspireServiceType)
	assert.Equal(t, "view", inspireServiceType.InspireServiceType)
}

func TestGetSDSServiceCategoryBySDSCategory(t *testing.T) {
	codelistLookupService, err := NewCodelist()
	require.NoError(t, err)

	serviceCategory, ok := codelistLookupService.GetSDSServiceCategoryBySDSCategory("invocable")
	assert.True(t, ok)
	assert.NotNil(t, serviceCategory)

	assert.Equal(
		t,
		"http://inspire.ec.europa.eu/id/ats/metadata/2.0/sds-invocable",
		serviceCategory.URI,
	)
	assert.Equal(t, "invocable", serviceCategory.Value)
	assert.Equal(t, "Aanroepbare datadienst", serviceCategory.Description)
}

func TestGetDataLicenseByLicenseURI(t *testing.T) {
	codelistLookupService, err := NewCodelist()
	require.NoError(t, err)

	for _, dataLicense := range codelistLookupService.DataLicenses {
		_, err = regexp.Compile(dataLicense.URIRegex)
		require.NoError(t, err)
	}

	dataLicense1, ok := codelistLookupService.GetDataLicenseByURI(
		"https://creativecommons.org/licenses/by/4.0/deed.nl",
	)
	assert.True(t, ok)
	assert.NotNil(t, dataLicense1)
	assert.Equal(t, "Open data (CC-BY)", dataLicense1.Value)
	assert.Equal(t, "Naamsvermelding verplicht, organisatienaam", dataLicense1.Description)

	dataLicense2, ok := codelistLookupService.GetDataLicenseByURI(
		"https://creativecommons.org/licenses/by-sa/4.0/deed.nl",
	)
	assert.True(t, ok)
	assert.NotNil(t, dataLicense2)
	assert.Equal(t, "Gebruiksvoorwaarden (CC-by-sa)", dataLicense2.Value)
	assert.Equal(
		t,
		"Gelijk Delen, Naamsvermelding verplicht, organisatienaam",
		dataLicense2.Description,
	)

	dataLicense3, ok := codelistLookupService.GetDataLicenseByURI(
		"https://creativecommons.org/publicdomain/zero/1.0/deed.nl",
	)
	assert.True(t, ok)
	assert.NotNil(t, dataLicense3)
	assert.Equal(t, "Open data (CCO)", dataLicense3.Value)
	assert.Equal(t, "Geen beperkingen", dataLicense3.Description)
}
