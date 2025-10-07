package hvd

import "encoding/xml"

// RDF XML structures.
type RDF struct {
	XMLName      xml.Name      `xml:"RDF"`
	Descriptions []Description `xml:"Description"`
}

// Description struct for parsing RDF.
type Description struct {
	About      string  `xml:"about,attr"`
	Type       Type    `xml:"type"`
	Order      Order   `xml:"order"`
	Broader    Broader `xml:"broader"`
	PrefLabels []Label `xml:"prefLabel"`
	Identifier string  `xml:"identifier"`
}

// Type struct for parsing RDF.
type Type struct {
	Resource string `xml:"resource,attr"`
}

// Order struct for parsing RDF.
type Order struct {
	Value string `xml:",chardata"`
}

// Broader struct for parsing RDF.
type Broader struct {
	Resource string `xml:"resource,attr"`
}

// Label struct for parsing RDF.
type Label struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}
