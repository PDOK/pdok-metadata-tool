// Package codelist provides lookup access to the values in the codelists.
package codelist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
)

// Codelist is used for unmarshalling the JSON codelists.
type Codelist struct {
	ReferenceSystems    map[string]ReferenceSystem    `json:"referenceSystems"`
	InspireThemes       map[string]string             `json:"inspireThemes"`
	HvdCategories       map[string]string             `json:"hvdCategories"`
	Protocol            map[string]ProtocolDetails    `json:"protocols"`
	InspireServiceTypes []InspireServiceType          `json:"inspireServiceTypes"`
	SDSServiceCategory  map[string]SDSServiceCategory `json:"sdsServiceCategories"`
	DataLicenses        []DataLicense                 `json:"dataLicenses"`
}

// ReferenceSystem is used for unmarshalling the JSON codelists.
type ReferenceSystem struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// ProtocolDetails is used for unmarshalling the JSON codelists.
type ProtocolDetails struct {
	ServiceProtocolURL              string `json:"serviceProtocolUrl"`
	ServiceProtocol                 string `json:"serviceProtocol"`
	ProtocolReleaseDate             string `json:"protocolReleaseDate"`
	ProtocolVersion                 string `json:"protocolVersion"`
	ServiceProtocolName             string `json:"serviceProtocolName"`
	SpatialDataserviceCategory      string `json:"spatialDataserviceCategory"`
	SpatialDataserviceCategoryURI   string `json:"spatialDataserviceCategoryUri"`
	SpatialDataserviceCategoryLabel string `json:"spatialDataserviceCategoryLabel"`
	ServiceAccessPointOperation     string `json:"serviceAccessPointOperation"`
}

// InspireServiceType is used for unmarshalling the JSON codelists.
type InspireServiceType struct {
	OGCServiceTypes              []string `json:"ogcServiceTypes"`
	InspireURI                   string   `json:"inspireUri"`
	InspireServiceType           string   `json:"inspireServiceType"`
	InspireTechnicalGuidance     string   `json:"inspireTechnicalGuidance"`
	InspireTechnicalGuidanceDate string   `json:"inspireTechnicalGuidanceDate"`
}

// DataLicense is used for unmarshalling the JSON codelists.
type DataLicense struct {
	URI         string `json:"uri"`
	URIRegex    string `json:"uriRegex"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

// SDSServiceCategory is used for unmarshalling the JSON codelists.
type SDSServiceCategory struct {
	URI         string `json:"uri"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

// NewCodelist creates a new instance of Codelist.
func NewCodelist() (*Codelist, error) {
	data, err := os.ReadFile(getCodelistPath())
	if err != nil {
		return nil, fmt.Errorf("failed to read codelist file: %w", err)
	}

	var codelist Codelist
	if err := json.Unmarshal(data, &codelist); err != nil {
		return nil, fmt.Errorf("failed to parse codelist JSON: %w", err)
	}

	return &codelist, nil
}

func getCodelistPath() string {
	//nolint:dogsled
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	return filepath.Join(dir, "codelists.json")
}

// GetReferenceSystemByEPSGCode returns a ReferenceSystem for a given EPSG code.
func (cs *Codelist) GetReferenceSystemByEPSGCode(epsgCode string) (*ReferenceSystem, bool) {
	epsgCode = strings.ToUpper(epsgCode)
	rs, ok := cs.ReferenceSystems[epsgCode]

	return &rs, ok
}

// GetINSPIREThemeLabelByURI returns INSPIRE Theme label string for a given URI.
func (cs *Codelist) GetINSPIREThemeLabelByURI(uri string) (*string, bool) {
	uri = strings.ToLower(uri)

	if strings.HasPrefix(uri, "http://") {
		uri = strings.Replace(uri, "http://", "https://", 1)
	}

	it, ok := cs.InspireThemes[uri]

	return &it, ok
}

// GetProtocolDetailsByProtocol returns ProtocolDetails for a given protocol.
func (cs *Codelist) GetProtocolDetailsByProtocol(protocol string) (*ProtocolDetails, bool) {
	protocol = strings.ToLower(protocol)
	p, ok := cs.Protocol[protocol]

	return &p, ok
}

// GetInspireServiceTypeByServiceType returns a InspireServiceType for a given serviceType.
func (cs *Codelist) GetInspireServiceTypeByServiceType(
	serviceType string,
) (*InspireServiceType, bool) {
	serviceType = strings.ToLower(serviceType)
	for _, item := range cs.InspireServiceTypes {
		if slices.Contains(item.OGCServiceTypes, serviceType) {
			return &item, true
		}
	}

	return nil, false
}

// GetSDSServiceCategoryBySDSCategory returns a SDSServiceCategory for a given category.
func (cs *Codelist) GetSDSServiceCategoryBySDSCategory(
	category string,
) (*SDSServiceCategory, bool) {
	category = strings.ToLower(category)
	if category == "harmonized" {
		category = "harmonised"
	}

	sc, ok := cs.SDSServiceCategory[category]

	return &sc, ok
}

// GetDataLicenseByURI returns a DataLicense for a given URI.
func (cs *Codelist) GetDataLicenseByURI(uri string) (*DataLicense, bool) {
	uri = strings.ToLower(uri)

	for _, dataLicense := range cs.DataLicenses {
		re, err := regexp.Compile(dataLicense.URIRegex)
		if err != nil {
			return nil, false
		}

		if re.MatchString(uri) {
			return &dataLicense, true
		}
	}

	return nil, false
}
