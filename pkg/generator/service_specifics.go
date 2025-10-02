package generator

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pdok/pdok-metadata-tool/internal/common"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
	"time"
)

type ServiceSpecifics struct {
	Globals  GlobalConfig    `yaml:"globals"`
	Services []ServiceConfig `yaml:"services"`
}

type GlobalConfig struct {
	OverrideableFields `yaml:",inline"`
}

type ServiceConfig struct {
	OverrideableFields `yaml:",inline"`
	Type               string       `yaml:"type"`
	ID                 string       `yaml:"id"`
	AccessPoint        string       `yaml:"accessPoint"`
	InspireType        *InspireType `yaml:"inspireType"`

	// Pointer to globals
	Globals *GlobalConfig `yaml:"-"`
}

type InspireType string

const (
	Harmonised    InspireType = "HARMONISED"
	Interoperable InspireType = "INTEROPERABLE"
	Invocable     InspireType = "INVOCABLE"
)

func (iv *InspireType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	switch strings.ToUpper(s) {
	case "HARMONIZED", "HARMONISED":
		*iv = Harmonised
	case "INTEROPERABLE":
		*iv = Interoperable
	case "INVOCABLE":
		*iv = Invocable
	}

	return nil
}

type OverrideableFields struct {
	Title                     *string      `yaml:"title,omitempty"`
	CreationDate              *string      `yaml:"creationDate,omitempty"`
	Abstract                  *string      `yaml:"abstract,omitempty"`
	Keywords                  []string     `yaml:"keywords,omitempty"`
	ContactOrganisationName   *string      `yaml:"contactOrganisationName,omitempty"`
	ContactOrganisationUri    *string      `yaml:"contactOrganisationUri,omitempty"`
	ContactEmail              *string      `yaml:"contactEmail,omitempty"`
	ContactUrl                *string      `yaml:"contactUrl,omitempty"`
	InspireThemes             []string     `yaml:"inspireThemes,omitempty"`
	HvdCategories             []string     `yaml:"hvdCategories,omitempty"`
	ServiceLicense            *string      `yaml:"serviceLicense,omitempty"`
	UseLimitation             *string      `yaml:"useLimitation,omitempty"`
	BoundingBox               *BoundingBox `yaml:"boundingBox,omitempty"`
	LinkedDatasets            []string     `yaml:"linkedDatasets,omitempty"`
	CoordinateReferenceSystem *string      `yaml:"coordinateReferenceSystem,omitempty"`
	Thumbnails                []Thumbnail  `yaml:"thumbnails,omitempty"`
	QosAvailability           *float64     `yaml:"qosAvailability,omitempty"`
	QosPerformance            *float64     `yaml:"qosPerformance,omitempty"`
	QosCapacity               *int         `yaml:"qosCapacity,omitempty"`
}

type BoundingBox struct {
	MinX string `yaml:"minX"`
	MaxX string `yaml:"maxX"`
	MinY string `yaml:"minY"`
	MaxY string `yaml:"maxY"`
}

type Thumbnail struct {
	File        string `yaml:"file"`
	Description string `yaml:"description"`
	Filetype    string `yaml:"filetype"`
}

func (s *ServiceSpecifics) LoadFromYAML(filename string) error {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(yamlFile, s); err != nil {
		return err
	}

	// Setup pointer to Globals for each service
	for i := range s.Services {
		s.Services[i].Globals = &s.Globals
	}

	// Set default INSPIRE type
	for _, service := range s.Services {
		if len(service.GetInspireThemes()) > 0 && service.InspireType == nil {
			service.InspireType = common.Ptr(Harmonised)
		}
	}

	return nil
}

func (s *ServiceSpecifics) Validate() error {
	var validationErrors []string
	seenIDs := make(map[string]bool)

	for i, service := range s.Services {
		// Check for duplicate ID
		if seenIDs[service.ID] {
			validationErrors = append(validationErrors, fmt.Sprintf("Service[%d]: id is duplicate '%s'", i, service.ID))
		} else {
			seenIDs[service.ID] = true
		}

		// Validate individual service
		if err := service.Validate(); err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("Service[%d] (%s): %v", i, service.ID, err))
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation failed:\n%s", strings.Join(validationErrors, "\n"))
	}
	return nil
}

func (sc *ServiceConfig) Validate() error {
	var errors []string

	if sc.ID == "" {
		errors = append(errors, "id is required")
	} else {
		if _, err := uuid.Parse(sc.ID); err != nil {
			errors = append(errors, fmt.Sprintf("id is not a valid UUID: %s", sc.ID))
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

	if sc.GetAbstract() == "" {
		errors = append(errors, "abstract is required (either local or global)")
	}
	if len(sc.GetKeywords()) == 0 {
		errors = append(errors, "at least one keyword is required (either local or global)")
	}

	if sc.GetContactOrganisationName() == "" {
		errors = append(errors, "contactOrganisationName is required (either local or global)")
	}
	if sc.GetContactOrganisationUri() == "" {
		errors = append(errors, "contactOrganisationName is required (either local or global)")
	}
	if sc.GetContactEmail() == "" {
		errors = append(errors, "contactEmail is required (either local or global)")
	}
	if sc.GetContactUrl() == "" {
		errors = append(errors, "contactUrl is required (either local or global)")
	}
	if sc.GetServiceLicense() == "" {
		errors = append(errors, "serviceLicense is required (either local or global)")
	}
	if sc.GetCoordinateReferenceSystem() == "" {
		errors = append(errors, "serviceLicense is required (either local or global)")
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

	if sc.InspireType != nil && len(sc.GetInspireThemes()) == 0 {
		errors = append(errors, "inspireThemes are required when inspireType is set")
	}

	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "; "))
	}
	return nil
}

func (sc *ServiceConfig) GetTitle() string {
	if sc.Title != nil {
		return *sc.Title
	}
	if sc.Globals.Title != nil {
		postfix := ""
		switch sc.Type {
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

		return *sc.Globals.Title + postfix
	}
	return ""
}

func (sc *ServiceConfig) GetCreationDate() string {
	if sc.CreationDate != nil {
		return *sc.CreationDate
	}
	if sc.Globals.CreationDate != nil {
		return *sc.Globals.CreationDate
	}
	return ""
}

func (sc *ServiceConfig) GetAbstract() string {
	if sc.Abstract != nil {
		return *sc.Abstract
	}
	if sc.Globals.Abstract != nil {
		return *sc.Globals.Abstract
	}
	return ""
}

func (sc *ServiceConfig) GetKeywords() []string {
	if len(sc.Keywords) > 0 {
		return sc.Keywords
	}
	return sc.Globals.Keywords
}

func (sc *ServiceConfig) GetContactOrganisationName() string {
	if sc.ContactOrganisationName != nil {
		return *sc.ContactOrganisationName
	}
	if sc.Globals.ContactOrganisationName != nil {
		return *sc.Globals.ContactOrganisationName
	}
	return ""
}

func (sc *ServiceConfig) GetContactOrganisationUri() string {
	if sc.ContactOrganisationUri != nil {
		return *sc.ContactOrganisationUri
	}
	if sc.Globals.ContactOrganisationUri != nil {
		return *sc.Globals.ContactOrganisationUri
	}
	return ""
}

func (sc *ServiceConfig) GetContactEmail() string {
	if sc.ContactEmail != nil {
		return *sc.ContactEmail
	}
	if sc.Globals.ContactEmail != nil {
		return *sc.Globals.ContactEmail
	}
	return ""
}

func (sc *ServiceConfig) GetContactUrl() string {
	if sc.ContactUrl != nil {
		return *sc.ContactUrl
	}
	if sc.Globals.ContactUrl != nil {
		return *sc.Globals.ContactUrl
	}
	return ""
}

func (sc *ServiceConfig) GetInspireThemes() []string {
	if len(sc.InspireThemes) > 0 {
		return sc.InspireThemes
	}
	return sc.Globals.InspireThemes
}

func (sc *ServiceConfig) GetHvdCategories() []string {
	if len(sc.HvdCategories) > 0 {
		return sc.HvdCategories
	}
	return sc.Globals.HvdCategories
}

func (sc *ServiceConfig) GetServiceLicense() string {
	if sc.ServiceLicense != nil {
		return *sc.ServiceLicense
	}
	if sc.Globals.ServiceLicense != nil {
		return *sc.Globals.ServiceLicense
	}
	return ""
}

func (sc *ServiceConfig) GetUseLimitation() string {
	if sc.UseLimitation != nil {
		return *sc.UseLimitation
	}
	if sc.Globals.UseLimitation != nil {
		return *sc.Globals.UseLimitation
	}
	return "Geen beperkingen"
}

func (sc *ServiceConfig) GetBoundingBox() *BoundingBox {
	if sc.BoundingBox != nil {
		return sc.BoundingBox
	}
	return sc.Globals.BoundingBox
}

func (sc *ServiceConfig) GetLinkedDatasets() []string {
	if len(sc.LinkedDatasets) > 0 {
		return sc.LinkedDatasets
	}
	return sc.Globals.LinkedDatasets
}

func (sc *ServiceConfig) GetCoordinateReferenceSystem() string {
	if sc.CoordinateReferenceSystem != nil {
		return *sc.CoordinateReferenceSystem
	}
	if sc.Globals.CoordinateReferenceSystem != nil {
		return *sc.Globals.CoordinateReferenceSystem
	}
	return ""
}

func (sc *ServiceConfig) GetThumbnails() []Thumbnail {
	if len(sc.Thumbnails) > 0 {
		return sc.Thumbnails
	}
	return sc.Globals.Thumbnails
}

func (sc *ServiceConfig) GetQosAvailability() string {
	var value float64 = -999
	if sc.QosAvailability != nil {
		value = *sc.QosAvailability
	} else if sc.Globals.QosAvailability != nil {
		value = *sc.Globals.QosAvailability
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func (sc *ServiceConfig) GetQosPerformance() string {
	var value float64 = -999
	if sc.QosPerformance != nil {
		value = *sc.QosPerformance
	} else if sc.Globals.QosPerformance != nil {
		value = *sc.Globals.QosPerformance
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func (sc *ServiceConfig) GetQosCapacity() string {
	var value int = -999
	if sc.QosCapacity != nil {
		value = *sc.QosCapacity
	} else if sc.Globals.QosCapacity != nil {
		value = *sc.Globals.QosCapacity
	}
	return strconv.Itoa(value)
}
