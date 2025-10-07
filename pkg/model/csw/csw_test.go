package csw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRecordsOgcFilter_ToRequestBody(t *testing.T) {
	expectedRequestBody := `<csw:GetRecords xmlns:csw="http://www.opengis.net/cat/csw/2.0.2" xmlns:ogc="http://www.opengis.net/ogc" service="CSW" version="2.0.2" resultType="results" startPosition="1" maxRecords="5" outputFormat="application/xml" outputSchema="http://www.opengis.net/cat/csw/2.0.2" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.opengis.net/cat/csw/2.0.2 http://schemas.opengis.net/csw/2.0.2/CSW-discovery.xsd">
	<csw:Query typeNames="csw:Record">
		<csw:ElementSetName>summary</csw:ElementSetName>
                <csw:Constraint version="1.1.0">
                <ogc:Filter>
                    <ogc:And>
    <ogc:PropertyIsEqualTo>
    <ogc:PropertyName>dc:type</ogc:PropertyName>
    <ogc:Literal>service</ogc:Literal>
</ogc:PropertyIsEqualTo>
<ogc:Or>
    <ogc:PropertyIsLike escapeChar="" wildCard="%" singleChar="_">
    <ogc:PropertyName>dc:title</ogc:PropertyName>
    <ogc:Literal>%title%</ogc:Literal>
</ogc:PropertyIsLike>
<ogc:PropertyIsEqualTo>
    <ogc:PropertyName>dc:identifier</ogc:PropertyName>
    <ogc:Literal>C2DFBDBC-5092-11E0-BA8E-B62DE0D72086</ogc:Literal>
</ogc:PropertyIsEqualTo>
</ogc:Or>
</ogc:And>
                </ogc:Filter>
            </csw:Constraint>
	</csw:Query>
</csw:GetRecords>`

	title := "title"
	identifier := "C2DFBDBC-5092-11E0-BA8E-B62DE0D72086"
	filter := GetRecordsOgcFilter{
		MetadataType: Service,
		Title:        &title,
		Identifier:   &identifier,
	}

	requestBody, err := filter.ToRequestBody()
	require.NoError(t, err)
	assert.Equal(t, expectedRequestBody, requestBody)
}

func TestGetRecordsOgcFilter_getPropertyIsLikeClause(t *testing.T) {
	expectedClause := `<ogc:PropertyIsLike escapeChar="" wildCard="%" singleChar="_">
    <ogc:PropertyName>property</ogc:PropertyName>
    <ogc:Literal>%value%</ogc:Literal>
</ogc:PropertyIsLike>`
	filter := GetRecordsOgcFilter{}

	clause := filter.getPropertyIsLikeClause("property", "value")
	assert.Equal(t, expectedClause, clause)
}

func TestGetRecordsOgcFilter_getPropertyIsEqualToClause(t *testing.T) {
	expectedClause := `<ogc:PropertyIsEqualTo>
    <ogc:PropertyName>property</ogc:PropertyName>
    <ogc:Literal>value</ogc:Literal>
</ogc:PropertyIsEqualTo>`

	filter := GetRecordsOgcFilter{}
	clause := filter.getPropertyIsEqualToClause("property", "value")
	assert.Equal(t, expectedClause, clause)
}
