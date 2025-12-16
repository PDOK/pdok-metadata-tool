// Package csw holds models for handling CSW requests.
package csw

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"

	"github.com/pdok/pdok-metadata-tool/pkg/model/iso1911x"
)

// Note: Use iso1911x.MetadataType directly for type and constants (Service/Dataset).

// GetRecordByIDResponse struct for unmarshalling a CSW GetRecordByID response.
type GetRecordByIDResponse struct {
	XMLName    xml.Name            `xml:"GetRecordByIdResponse"`
	Text       string              `xml:",chardata"`
	Csw        string              `xml:"csw,attr"`
	MDMetadata iso1911x.MDMetadata `xml:"MD_Metadata"`
}

// GetRecordsResponse struct for unmarshalling a CSW GetRecords response.
type GetRecordsResponse struct {
	XMLName       xml.Name `xml:"GetRecordsResponse"`
	SearchResults struct {
		NumberOfRecordsMatched string          `xml:"numberOfRecordsMatched,attr"`
		NextRecord             string          `xml:"nextRecord,attr"`
		SummaryRecords         []SummaryRecord `xml:"SummaryRecord"`
	} `xml:"SearchResults"`
}

// GetRecordsCQLConstraint struct for creating a CQL constraint.
type GetRecordsCQLConstraint struct {
	MetadataType     *iso1911x.MetadataType
	OrganisationName *string
}

// ToQueryParameter returns a query parameter string based on the constraints.
func (c *GetRecordsCQLConstraint) ToQueryParameter() (constraint string) {
	var constraints []string

	if c.MetadataType != nil {
		constraints = append(constraints, fmt.Sprintf("type='%s'", *c.MetadataType))
	}

	if c.OrganisationName != nil {
		constraints = append(constraints, fmt.Sprintf("OrganisationName='%s'", *c.OrganisationName))
	}

	if len(constraints) == 0 {
		return constraint
	}

	constraint += "&constraintLanguage=CQL_TEXT"
	constraint += "&constraint_language_version=1.1.0"

	// Build the raw CQL expression and URL-encode it as a single query parameter value
	var (
		expr     string
		exprSb67 strings.Builder
	)

	for i, c := range constraints {
		exprSb67.WriteString(c)

		if i < len(constraints)-1 {
			exprSb67.WriteString(" AND ")
		}
	}

	expr += exprSb67.String()

	constraint += "&constraint=" + url.QueryEscape(expr)

	return constraint
}

// GetRecordsOgcFilter struct for creating an OgcFilter for a CSW GetRecords request.
type GetRecordsOgcFilter struct {
	MetadataType iso1911x.MetadataType
	Title        *string
	Identifier   *string
}

// ToRequestBody Returns a request body string for a CSW GetRecords request.
func (f *GetRecordsOgcFilter) ToRequestBody() (ogcFilter string, err error) {
	template := `<csw:GetRecords xmlns:csw="http://www.opengis.net/cat/csw/2.0.2" xmlns:ogc="http://www.opengis.net/ogc" service="CSW" version="2.0.2" resultType="results" startPosition="1" maxRecords="5" outputFormat="application/xml" outputSchema="http://www.opengis.net/cat/csw/2.0.2" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.opengis.net/cat/csw/2.0.2 http://schemas.opengis.net/csw/2.0.2/CSW-discovery.xsd">
	<csw:Query typeNames="csw:Record">
		<csw:ElementSetName>summary</csw:ElementSetName>
                <csw:Constraint version="1.1.0">
                <ogc:Filter>
                    %s
                </ogc:Filter>
            </csw:Constraint>
	</csw:Query>
</csw:GetRecords>`

	metadataTypeFilter := f.getPropertyIsEqualToClause("dc:type", f.MetadataType.String())

	var filters []string
	if f.Title != nil && len(*f.Title) > 0 {
		filters = append(filters, f.getPropertyIsLikeClause("dc:title", *f.Title))
	}

	if f.Identifier != nil && len(*f.Identifier) > 0 {
		filters = append(filters, f.getPropertyIsEqualToClause("dc:identifier", *f.Identifier))
	}

	var filter string

	switch len(filters) {
	case 0:
		filter = metadataTypeFilter
	case 1:
		filter = wrapFiltersInAndOperator([]string{metadataTypeFilter, filters[0]})
	default:
		filter = wrapFiltersInAndOperator(
			[]string{metadataTypeFilter, wrapFiltersInOrOperator(filters)},
		)
	}

	requestBody := fmt.Sprintf(template, filter)

	return requestBody, nil
}

func wrapFiltersInOrOperator(filters []string) string {
	template := `<ogc:Or>
    %s
</ogc:Or>`

	return fmt.Sprintf(template, strings.Join(filters, "\n"))
}

func wrapFiltersInAndOperator(filters []string) string {
	template := `<ogc:And>
    %s
</ogc:And>`

	return fmt.Sprintf(template, strings.Join(filters, "\n"))
}

func (f *GetRecordsOgcFilter) getPropertyIsLikeClause(
	property string,
	value string,
) (clause string) {
	template := `<ogc:PropertyIsLike escapeChar="" wildCard="%%" singleChar="_">
    <ogc:PropertyName>%s</ogc:PropertyName>
    <ogc:Literal>%%%s%%</ogc:Literal>
</ogc:PropertyIsLike>`

	return fmt.Sprintf(template, property, value)
}

func (f *GetRecordsOgcFilter) getPropertyIsEqualToClause(
	property string,
	value string,
) (clause string) {
	template := `<ogc:PropertyIsEqualTo>
    <ogc:PropertyName>%s</ogc:PropertyName>
    <ogc:Literal>%s</ogc:Literal>
</ogc:PropertyIsEqualTo>`

	return fmt.Sprintf(template, property, value)
}

// Helper methods for MDMetadata live on the underlying iso1911x.MDMetadata type.

// SummaryRecord struct for unmarshalling the CSW response.
type SummaryRecord struct {
	Identifier string `json:"identifier" xml:"identifier"`
	Title      string `json:"title"      xml:"title"`
	Type       string `json:"type"       xml:"type"`
}
