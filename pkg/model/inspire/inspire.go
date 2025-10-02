package inspire

import (
	"fmt"
)

var InspireEndpoint = "https://inspire.ec.europa.eu"

type InspireRegisterKind string

var InspireKindTheme InspireRegisterKind = "theme"
var InspireKindLayer InspireRegisterKind = "layer"

type InspireRegisterLanguage string

var InspireDutch InspireRegisterLanguage = "nl"
var InspireEnglish InspireRegisterLanguage = "en"

var InspireRegisterKinds = []InspireRegisterKind{InspireKindTheme, InspireKindLayer}
var InspireRegisterLanguages = []InspireRegisterLanguage{InspireDutch, InspireEnglish}

type InspireVariant string

const (
	Harmonised InspireVariant = "HARMONISED"
	AsIs       InspireVariant = "ASIS"
)

// InspireItem is an interface that both InspireLayer and InspireTheme implement
type InspireItem interface {
	GetId() string
	GetLabelDutch() string
	GetLabelEnglish() string
}

func GetInspireEndpoint(kind InspireRegisterKind, language InspireRegisterLanguage) string {
	return fmt.Sprintf("%s/%s/%s.%s.json", InspireEndpoint, kind, kind, language)
}

func GetInspirePath(kind InspireRegisterKind, language InspireRegisterLanguage) string {
	return fmt.Sprintf("%s.%s.json", kind, language)
}
