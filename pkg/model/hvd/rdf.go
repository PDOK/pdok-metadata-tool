package hvd

import "encoding/xml"

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
