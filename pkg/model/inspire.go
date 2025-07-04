package model

import "fmt"

var InspireEndpoint = "https://inspire.ec.europa.eu"

type InspireRegisterKind string

var InspireKindTheme InspireRegisterKind = "theme"
var InspireKindLayer InspireRegisterKind = "layer"

type InspireRegisterLanguage string

var InspireDutch InspireRegisterLanguage = "nl"
var InspireEnglish InspireRegisterLanguage = "en"

var InspireRegisterKinds = []InspireRegisterKind{InspireKindTheme, InspireKindLayer}
var InspireRegisterLanguages = []InspireRegisterLanguage{InspireDutch, InspireEnglish}

func GetInspireEndpoint(kind InspireRegisterKind, language InspireRegisterLanguage) string {
	return fmt.Sprintf("%s/%s/%s.%s.json", InspireEndpoint, kind, kind, language)
}

func GetInspirePath(kind InspireRegisterKind, language InspireRegisterLanguage) string {
	return fmt.Sprintf("%s.%s.json", kind, language)
}

// InspireTheme represents an INSPIRE theme with both English and Dutch labels
type InspireTheme struct {
	Id           string `json:"id"`           // Primary Key, Unique, 2 characters
	Order        int    `json:"order"`        // Order number
	LabelDutch   string `json:"labelDutch"`   // Dutch label
	LabelEnglish string `json:"labelEnglish"` // English label
	URL          string `json:"url"`          // URL for the theme
}

// InspireLayer represents an INSPIRE layer with both English and Dutch labels
type InspireLayer struct {
	Id           string `json:"id"`           // Primary Key, Unique, up to 100 characters
	LabelDutch   string `json:"labelDutch"`   // Dutch label
	LabelEnglish string `json:"labelEnglish"` // English label
}
