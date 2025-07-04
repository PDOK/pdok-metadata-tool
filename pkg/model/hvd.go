package model

import (
	"encoding/xml"
)

// RDF XML structures
type RDF struct {
	XMLName      xml.Name      `xml:"RDF"`
	Descriptions []Description `xml:"Description"`
}

type Description struct {
	About      string  `xml:"about,attr"`
	Type       Type    `xml:"type"`
	Order      Order   `xml:"order"`
	Broader    Broader `xml:"broader"`
	PrefLabels []Label `xml:"prefLabel"`
	Identifier string  `xml:"identifier"`
}

type Type struct {
	Resource string `xml:"resource,attr"`
}

type Order struct {
	Value string `xml:",chardata"`
}

type Broader struct {
	Resource string `xml:"resource,attr"`
}

type Label struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// HVDCategory represents a High Value Dataset category
type HVDCategory struct {
	// Id is the primary key, must be unique
	Id string `json:"id" validate:"required,max=10"`

	// Parent is a foreign key to another HVDCategory.Id
	Parent string `json:"parent" validate:"max=10"`

	// Order represents the display order
	Order string `json:"order" validate:"max=6"`

	// LabelDutch is the Dutch language label
	LabelDutch string `json:"labelDutch"`

	// LabelEnglish is the English language label
	LabelEnglish string `json:"labelEnglish"`
}
