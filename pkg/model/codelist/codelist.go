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

type Codelist struct {
	ReferenceSystems    map[string]ReferenceSystem    `json:"reference_systems"`
	InspireThemes       map[string]string             `json:"inspire_themes"`
	HvdCategories       map[string]string             `json:"hvd_categories"`
	Protocol            map[string]ProtocolDetails    `json:"protocols"`
	InspireServiceTypes []InspireServiceType          `json:"inspire_service_types"`
	SDSServiceCategory  map[string]SDSServiceCategory `json:"sds_service_categories"`
	DataLicenses        []DataLicense                 `json:"data_licenses"`
}

type Codelists struct {
	ReferenceSystems map[string]ReferenceSystem `json:"reference_systems"`
	InspireThemes    map[string]string          `json:"inspire_themes"`
	// HvdCategories       map[string]string             `json:"hvd_categories"`
	Protocol            map[string]ProtocolDetails    `json:"protocols"`
	InspireServiceTypes []InspireServiceType          `json:"inspire_service_types"`
	SDSServiceCategory  map[string]SDSServiceCategory `json:"sds_service_categories"`
	DataLicenses        []DataLicense                 `json:"data_licenses"`
}

type ReferenceSystem struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type ProtocolDetails struct {
	ServiceProtocolURL              string `json:"service_protocol_url"`
	ServiceProtocol                 string `json:"service_protocol"`
	ProtocolReleaseDate             string `json:"protocol_release_date"`
	ProtocolVersion                 string `json:"protocol_version"`
	ServiceProtocolName             string `json:"service_protocol_name"`
	SpatialDataserviceCategory      string `json:"spatial_dataservice_category"`
	SpatialDataserviceCategoryURI   string `json:"spatial_dataservice_category_uri"`
	SpatialDataserviceCategoryLabel string `json:"spatial_dataservice_category_label"`
	ServiceAccessPointOperation     string `json:"service_access_point_operation"`
}

type InspireServiceType struct {
	OGCServiceTypes              []string `json:"ogc_service_types"`
	InspireURI                   string   `json:"inspire_uri"`
	InspireServiceType           string   `json:"inspire_servicetype"`
	InspireTechnicalGuidance     string   `json:"inspire_technicalguidance"`
	InspireTechnicalGuidanceDate string   `json:"inspire_technicalguidance_date"`
}

type DataLicense struct {
	URI         string `json:"uri"`
	URIRegex    string `json:"uri_regex"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type SDSServiceCategory struct {
	URI         string `json:"uri"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

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
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "codelists.json")
}

func (cs *Codelist) GetReferenceSystemByEPSGCode(EPSGCode string) (*ReferenceSystem, bool) {
	EPSGCode = strings.ToUpper(EPSGCode)
	rs, ok := cs.ReferenceSystems[EPSGCode]
	return &rs, ok
}

func (cs *Codelist) GetINSPIREThemeLabelByURI(URI string) (*string, bool) {
	URI = strings.ToLower(URI)

	if strings.HasPrefix(URI, "http://") {
		URI = strings.Replace(URI, "http://", "https://", 1)
	}

	it, ok := cs.InspireThemes[URI]
	return &it, ok
}

func (cs *Codelist) GetProtocolDetailsByProtocol(protocol string) (*ProtocolDetails, bool) {
	protocol = strings.ToLower(protocol)
	p, ok := cs.Protocol[protocol]
	return &p, ok
}

func (cs *Codelist) GetInspireServiceTypeByServiceType(serviceType string) (*InspireServiceType, bool) {
	serviceType = strings.ToLower(serviceType)
	for _, item := range cs.InspireServiceTypes {
		if slices.Contains(item.OGCServiceTypes, serviceType) {
			return &item, true
		}
	}
	return nil, false
}

func (cs *Codelist) GetSDSServiceCategoryBySDSCategory(category string) (*SDSServiceCategory, bool) {
	category = strings.ToLower(category)
	if category == "harmonized" {
		category = "harmonised"
	}
	sc, ok := cs.SDSServiceCategory[category]
	return &sc, ok
}

func (cs *Codelist) GetDataLicenseByURI(URI string) (*DataLicense, bool) {
	URI = strings.ToLower(URI)
	for _, dataLicense := range cs.DataLicenses {

		re, err := regexp.Compile(dataLicense.URIRegex)
		if err != nil {
			return nil, false
		}

		if re.MatchString(URI) {
			return &dataLicense, true
		}
	}

	return nil, false
}
