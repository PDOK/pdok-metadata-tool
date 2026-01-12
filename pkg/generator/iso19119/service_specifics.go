package iso19119

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/google/uuid"
	"github.com/pdok/pdok-metadata-tool/v2/internal/common"
)

// ServiceSpecifics struct for unmarshalling the input for service metadata generation.
type ServiceSpecifics struct {
	Globals  GlobalConfig    `json:"globals,omitempty"  yaml:"globals,omitempty"`
	Services []ServiceConfig `json:"services,omitempty" yaml:"services,omitempty"`
}

// GlobalConfig struct for unmarshalling service specifics input.
type GlobalConfig struct {
	OverrideableFields `json:",inline,omitempty" yaml:",inline,omitempty"`

	InspireDatasetType *InspireDatasetType `json:"inspireDatasetType,omitempty" yaml:"inspireDatasetType,omitempty"`
}

// ServiceConfig struct for unmarshalling service specifics input.
type ServiceConfig struct {
	OverrideableFields `json:",inline,omitempty" yaml:",inline,omitempty"`

	Type               string              `json:"type,omitempty"               yaml:"type,omitempty"`
	ID                 string              `json:"id,omitempty"                 yaml:"id,omitempty"`
	AccessPoint        string              `json:"accessPoint,omitempty"        yaml:"accessPoint,omitempty"`
	ServiceInspireType *InspireServiceType `json:"serviceInspireType,omitempty" yaml:"serviceInspireType,omitempty"`

	// Pointer to globals
	Globals *GlobalConfig `json:"globals,omitempty" yaml:"globals,omitempty"`
}

func (sc ServiceConfig) GetID() string { return sc.ID }

// InspireDatasetType struct for unmarshalling service specifics input.
type InspireDatasetType string

// Values for InspireDatasetType.
const (
	Harmonised InspireDatasetType = "HARMONISED"
	AsIs       InspireDatasetType = "AS-IS"
)

// UnmarshalYAML unmarshalls the expected string for INSPIRE types.
func (st *InspireDatasetType) UnmarshalYAML(unmarshal func(any) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	normalized := strings.ToUpper(strings.ReplaceAll(s, "-", " "))
	normalized = strings.TrimSpace(normalized)

	inspireMap := map[string]InspireDatasetType{
		"ASIS":       AsIs,
		"AS IS":      AsIs,
		"HARMONISED": Harmonised,
		"HARMONIZED": Harmonised,
	}

	if val, ok := inspireMap[normalized]; ok {
		*st = val

		return nil
	}

	return fmt.Errorf("invalid InspireDatasetType: %s", s)
}

// InspireServiceType struct for unmarshalling service specifics input.
type InspireServiceType string

// Values for InspireServiceType.
const (
	NetworkService InspireServiceType = "NETWORKSERVICE"
	Interoperable  InspireServiceType = "INTEROPERABLE"
	Invocable      InspireServiceType = "INVOCABLE"
)

// UnmarshalYAML unmarshalls the expected string for INSPIRE types.
func (st *InspireServiceType) UnmarshalYAML(unmarshal func(any) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	normalized := strings.ToUpper(strings.ReplaceAll(s, "-", " "))
	normalized = strings.TrimSpace(normalized)

	inspireMap := map[string]InspireServiceType{
		"NETWORKSERVICE":   NetworkService,
		"NETWORK SERVICE":  NetworkService,
		"NETWORKSERVICES":  NetworkService,
		"NETWORK SERVICES": NetworkService,
		"INTEROPERABLE":    Interoperable,
		"INVOCABLE":        Invocable,
	}

	if val, ok := inspireMap[normalized]; ok {
		*st = val

		return nil
	}

	return fmt.Errorf("invalid InspireServiceType: %s", s)
}

// OverrideableFields struct for unmarshalling service specifics input.
type OverrideableFields struct {
	Title                     *string      `json:"title,omitempty"                     yaml:"title,omitempty"`
	CreationDate              *string      `json:"creationDate,omitempty"              yaml:"creationDate,omitempty"`
	RevisionDate              *string      `json:"revisionDate,omitempty"              yaml:"revisionDate,omitempty"`
	Abstract                  *string      `json:"abstract,omitempty"                  yaml:"abstract,omitempty"`
	Keywords                  []string     `json:"keywords,omitempty"                  yaml:"keywords,omitempty"`
	ContactOrganisationName   *string      `json:"contactOrganisationName,omitempty"   yaml:"contactOrganisationName,omitempty"`
	ContactOrganisationURI    *string      `json:"contactOrganisationUri,omitempty"    yaml:"contactOrganisationUri,omitempty"`
	ContactEmail              *string      `json:"contactEmail,omitempty"              yaml:"contactEmail,omitempty"`
	ContactURL                *string      `json:"contactUrl,omitempty"                yaml:"contactUrl,omitempty"`
	InspireThemes             []string     `json:"inspireThemes,omitempty"             yaml:"inspireThemes,omitempty"`
	HvdCategories             []string     `json:"hvdCategories,omitempty"             yaml:"hvdCategories,omitempty"`
	ServiceLicense            *string      `json:"serviceLicense,omitempty"            yaml:"serviceLicense,omitempty"`
	UseLimitation             *string      `json:"useLimitation,omitempty"             yaml:"useLimitation,omitempty"`
	BoundingBox               *BoundingBox `json:"boundingBox,omitempty"               yaml:"boundingBox,omitempty"`
	LinkedDatasets            []string     `json:"linkedDatasets,omitempty"            yaml:"linkedDatasets,omitempty"`
	CoordinateReferenceSystem *string      `json:"coordinateReferenceSystem,omitempty" yaml:"coordinateReferenceSystem,omitempty"`
	Thumbnails                []Thumbnail  `json:"thumbnails,omitempty"                yaml:"thumbnails,omitempty"`
	QosAvailability           *float64     `json:"qosAvailability,omitempty"           yaml:"qosAvailability,omitempty"`
	QosPerformance            *float64     `json:"qosPerformance,omitempty"            yaml:"qosPerformance,omitempty"`
	QosCapacity               *int         `json:"qosCapacity,omitempty"               yaml:"qosCapacity,omitempty"`
}

// BoundingBox struct for unmarshalling service specifics input.
type BoundingBox struct {
	MinX string `json:"minX,omitempty" yaml:"minX,omitempty"`
	MaxX string `json:"maxX,omitempty" yaml:"maxX,omitempty"`
	MinY string `json:"minY,omitempty" yaml:"minY,omitempty"`
	MaxY string `json:"maxY,omitempty" yaml:"maxY,omitempty"`
}

// Thumbnail struct for unmarshalling service specifics input.
type Thumbnail struct {
	File        string `json:"file,omitempty"        yaml:"file,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Filetype    string `json:"filetype,omitempty"    yaml:"filetype,omitempty"`
}

// LoadFromYamlOrJson unmarshalls the input for the given input file.
func (s *ServiceSpecifics) LoadFromYamlOrJson(filename string) error {
	//nolint:gosec
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(yamlFile, s); err != nil {
		return err
	}

	s.InitializeFields()

	return nil
}

// InitializeFields Sets pointers and inferred values
func (s *ServiceSpecifics) InitializeFields() {
	// Setup pointer to Globals for each service
	for i := range s.Services {
		s.Services[i].Globals = &s.Globals
	}

	s.setInspireTypes()
}

// Validate the ServiceSpecifics on a global level, also calls Validate on service level.
func (s *ServiceSpecifics) Validate() error {
	var validationErrors []string

	seenIDs := make(map[string]bool)

	for i, service := range s.Services {
		// Check for duplicate ID
		if seenIDs[service.ID] {
			validationErrors = append(
				validationErrors,
				fmt.Sprintf("Service[%d]: id is duplicate '%s'", i, service.ID),
			)
		} else {
			seenIDs[service.ID] = true
		}

		// Validate individual service
		if err := service.Validate(); err != nil {
			validationErrors = append(
				validationErrors,
				fmt.Sprintf("Service[%d] (%s): %v", i, service.ID, err),
			)
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation failed:\n%s", strings.Join(validationErrors, "\n"))
	}

	return nil
}

// Validate the ServiceSpecifics on service level.
//
//nolint:cyclop
func (sc ServiceConfig) Validate() error {
	var errors []string

	if sc.ID == "" {
		errors = append(errors, "id is required")
	} else {
		if _, err := uuid.Parse(sc.ID); err != nil {
			errors = append(errors, "id is not a valid UUID: "+sc.ID)
		}
	}

	if sc.AccessPoint == "" {
		errors = append(errors, "accessPoint is required")
	}

	if sc.GetTitle() == "" {
		errors = append(errors, "title is required (either local or global)")
	}

	if sc.GetCreationDate() == "" {
		errors = append(errors, "creationDate is required (either local or global)")
	} else {
		_, err := time.Parse("2006-01-02", sc.GetCreationDate())
		if err != nil {
			errors = append(errors, "creationDate does not match the date format 'YYYY-MM-DD'")
		}
	}

	if sc.GetRevisionDate() == "" {
		errors = append(errors, "revisionDate is required (either local or global)")
	} else {
		_, err := time.Parse("2006-01-02", sc.GetRevisionDate())
		if err != nil {
			errors = append(errors, "revisionDate does not match the date format 'YYYY-MM-DD'")
		}
	}

	if sc.GetAbstract() == "" {
		errors = append(errors, "abstract is required (either local or global)")
	}

	if len(sc.GetKeywords()) == 0 {
		errors = append(errors, "at least one keyword is required (either local or global)")
	}

	if sc.GetContactOrganisationName() == "" {
		errors = append(errors, "contactOrganisationName is required (either local or global)")
	}

	if sc.GetContactOrganisationURI() == "" {
		errors = append(errors, "GetContactOrganisationURI is required (either local or global)")
	}

	if sc.GetContactEmail() == "" {
		errors = append(errors, "contactEmail is required (either local or global)")
	}

	if sc.GetContactURL() == "" {
		errors = append(errors, "contactUrl is required (either local or global)")
	}

	if sc.GetServiceLicense() == "" {
		errors = append(errors, "serviceLicense is required (either local or global)")
	}

	if sc.GetBoundingBox() == nil {
		errors = append(errors, "boundingBox is required (either local or global)")
	}

	if sc.GetQosAvailability() == "-999" {
		errors = append(errors, "qosAvailability is required (either local or global)")
	}

	if sc.GetQosPerformance() == "-999" {
		errors = append(errors, "qosPerformance is required (either local or global)")
	}

	if sc.GetQosCapacity() == "-999" {
		errors = append(errors, "qosCapacity is required (either local or global)")
	}

	if sc.ServiceInspireType != nil && len(sc.GetInspireThemes()) == 0 {
		errors = append(errors, "inspireThemes are required when inspireType is set")
	}

	if sc.ServiceInspireType == nil && len(sc.GetInspireThemes()) > 0 {
		errors = append(errors, "inspireType is required when inspireThemes are set")
	}

	if sc.Globals.InspireDatasetType != nil &&
		*sc.Globals.InspireDatasetType == Harmonised &&
		len(sc.GetInspireThemes()) != 1 {
		errors = append(errors,
			"exactly 1 inspireTheme must be set if InspireDatasetType is 'harmonised'")
	}

	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "; "))
	}

	return nil
}

// GetTitle returns the (overrideable) title, and possibly adds a postfix.
func (sc ServiceConfig) GetTitle() string {
	if sc.Title != nil {
		return *sc.Title
	}

	if sc.Globals.Title != nil {
		postfix := ""

		switch strings.ToLower(sc.Type) {
		case "wms":
			postfix = " WMS"
		case "wfs":
			postfix = " WFS"
		case "atom":
			postfix = " ATOM"
		case "oaf":
			postfix = " OGC API Features"
		case "oat":
			postfix = " OGC API (Vector) Tiles"
		}

		// Only add the postfix if it's not already in the title
		if !strings.HasSuffix(strings.ToLower(*sc.Globals.Title), strings.ToLower(postfix)) {
			return *sc.Globals.Title + postfix
		}

		return *sc.Globals.Title
	}

	return ""
}

// GetCreationDate returns the (overrideable) creation date.
func (sc ServiceConfig) GetCreationDate() string {
	if sc.CreationDate != nil {
		return *sc.CreationDate
	}

	if sc.Globals.CreationDate != nil {
		return *sc.Globals.CreationDate
	}

	return ""
}

// GetRevisionDate returns the (overrideable) revision date.
func (sc ServiceConfig) GetRevisionDate() string {
	if sc.RevisionDate != nil {
		return *sc.RevisionDate
	}

	if sc.Globals.RevisionDate != nil {
		return *sc.Globals.RevisionDate
	}

	return ""
}

// GetAbstract returns the (overrideable) abstract.
func (sc ServiceConfig) GetAbstract() string {
	if sc.Abstract != nil {
		return *sc.Abstract
	}

	if sc.Globals.Abstract != nil {
		return *sc.Globals.Abstract
	}

	return ""
}

// GetKeywords returns the (overrideable) keywords.
func (sc ServiceConfig) GetKeywords() []string {
	if len(sc.Keywords) > 0 {
		return sc.Keywords
	}

	return sc.Globals.Keywords
}

// GetContactOrganisationName returns the (overrideable) contact organisation name.
func (sc ServiceConfig) GetContactOrganisationName() string {
	if sc.ContactOrganisationName != nil {
		return *sc.ContactOrganisationName
	}

	if sc.Globals.ContactOrganisationName != nil {
		return *sc.Globals.ContactOrganisationName
	}

	return ""
}

// GetContactOrganisationURI returns the (overrideable) contact organisation URI.
func (sc ServiceConfig) GetContactOrganisationURI() string {
	if sc.ContactOrganisationURI != nil {
		return *sc.ContactOrganisationURI
	}

	if sc.Globals.ContactOrganisationURI != nil {
		return *sc.Globals.ContactOrganisationURI
	}

	return ""
}

// GetContactEmail returns the (overrideable) contact email.
func (sc ServiceConfig) GetContactEmail() string {
	if sc.ContactEmail != nil {
		return *sc.ContactEmail
	}

	if sc.Globals.ContactEmail != nil {
		return *sc.Globals.ContactEmail
	}

	return ""
}

// GetContactURL returns the (overrideable) contact URL.
func (sc ServiceConfig) GetContactURL() string {
	if sc.ContactURL != nil {
		return *sc.ContactURL
	}

	if sc.Globals.ContactURL != nil {
		return *sc.Globals.ContactURL
	}

	return ""
}

// GetInspireThemes returns the (overrideable) INSPIRE themes.
func (sc ServiceConfig) GetInspireThemes() []string {
	if len(sc.InspireThemes) > 0 {
		return sc.InspireThemes
	}

	return sc.Globals.InspireThemes
}

// GetHvdCategories returns the (overrideable) HVD categories.
func (sc ServiceConfig) GetHvdCategories() []string {
	if len(sc.HvdCategories) > 0 {
		return sc.HvdCategories
	}

	return sc.Globals.HvdCategories
}

// GetServiceLicense returns the (overrideable) service license.
func (sc ServiceConfig) GetServiceLicense() string {
	if sc.ServiceLicense != nil {
		return *sc.ServiceLicense
	}

	if sc.Globals.ServiceLicense != nil {
		return *sc.Globals.ServiceLicense
	}

	return ""
}

// GetUseLimitation returns the (overrideable) use limitation.
func (sc ServiceConfig) GetUseLimitation() string {
	if sc.UseLimitation != nil {
		return *sc.UseLimitation
	}

	if sc.Globals.UseLimitation != nil {
		return *sc.Globals.UseLimitation
	}

	return "Geen beperkingen"
}

// GetBoundingBox returns the (overrideable) bounding box.
func (sc ServiceConfig) GetBoundingBox() *BoundingBox {
	if sc.BoundingBox != nil {
		return sc.BoundingBox
	}

	if sc.Globals.BoundingBox != nil {
		return sc.Globals.BoundingBox
	}

	return nil
}

// GetLinkedDatasets returns the (overrideable) linked datasets.
func (sc ServiceConfig) GetLinkedDatasets() []string {
	if len(sc.LinkedDatasets) > 0 {
		return sc.LinkedDatasets
	}

	return sc.Globals.LinkedDatasets
}

// GetCoordinateReferenceSystem returns the (overrideable) coordinate reference system.
func (sc ServiceConfig) GetCoordinateReferenceSystem() string {
	if sc.CoordinateReferenceSystem != nil {
		return *sc.CoordinateReferenceSystem
	}

	if sc.Globals.CoordinateReferenceSystem != nil {
		return *sc.Globals.CoordinateReferenceSystem
	}

	return ""
}

// GetThumbnails returns the (overrideable) thumbnails.
func (sc ServiceConfig) GetThumbnails() []Thumbnail {
	if len(sc.Thumbnails) > 0 {
		return sc.Thumbnails
	}

	return sc.Globals.Thumbnails
}

// GetQosAvailability returns the (overrideable) availability.
func (sc ServiceConfig) GetQosAvailability() string {
	var value float64 = -999
	if sc.QosAvailability != nil {
		value = *sc.QosAvailability
	} else if sc.Globals.QosAvailability != nil {
		value = *sc.Globals.QosAvailability
	}

	return strconv.FormatFloat(value, 'f', -1, 64)
}

// GetQosPerformance returns the (overrideable) performance.
func (sc ServiceConfig) GetQosPerformance() string {
	var value float64 = -999
	if sc.QosPerformance != nil {
		value = *sc.QosPerformance
	} else if sc.Globals.QosPerformance != nil {
		value = *sc.Globals.QosPerformance
	}

	return strconv.FormatFloat(value, 'f', -1, 64)
}

// GetQosCapacity returns the (overrideable) capacity.
func (sc ServiceConfig) GetQosCapacity() string {
	var value = -999
	if sc.QosCapacity != nil {
		value = *sc.QosCapacity
	} else if sc.Globals.QosCapacity != nil {
		value = *sc.Globals.QosCapacity
	}

	return strconv.Itoa(value)
}

// setInspireTypes sets INSPIRE Service types based on INSPIRE Dataset type
func (s *ServiceSpecifics) setInspireTypes() {
	inspireDatasetType := s.Globals.InspireDatasetType
	if inspireDatasetType == nil {
		return
	}

	// INSPIRE type mapping
	typeMap := map[string]InspireServiceType{}

	switch *inspireDatasetType {
	case AsIs:
		typeMap = map[string]InspireServiceType{
			"wms":  NetworkService,
			"atom": NetworkService,
			"wfs":  Invocable,
			"oaf":  Invocable,
			"oat":  Invocable,
		}
	case Harmonised:
		typeMap = map[string]InspireServiceType{
			"wms":  NetworkService,
			"atom": NetworkService,
			"wfs":  Interoperable,
			"oaf":  Interoperable,
			"oat":  Interoperable,
		}
	}

	for i := range s.Services {
		service := &s.Services[i]

		if service.ServiceInspireType != nil {
			// Don't touch ServiceInspireType if it is already set through an override
			continue
		}

		serviceType := strings.ToLower(service.Type)
		if inspireType, ok := typeMap[serviceType]; ok {
			service.ServiceInspireType = common.Ptr(inspireType)
		}
	}
}

func (sc ServiceConfig) isInspireSDS() bool {
	if sc.ServiceInspireType == nil {
		return false
	}

	return *sc.ServiceInspireType == Interoperable || *sc.ServiceInspireType == Invocable
}

func (sc ServiceConfig) isInspireNetworkService() bool {
	if sc.ServiceInspireType == nil {
		return false
	}

	return *sc.ServiceInspireType == NetworkService
}
