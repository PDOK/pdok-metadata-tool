package iso1911x

// CharacterStringTag struct for XML marshalling.
type CharacterStringTag struct {
	CharacterString string `xml:"gco:CharacterString"`
}

// DateTag struct for XML marshalling.
type DateTag struct {
	Date string `xml:"gco:Date"`
}

// BooleanTag struct for XML marshalling.
type BooleanTag struct {
	Value bool `xml:"gco:Boolean"`
}

// IntegerTag struct for XML marshalling.
type IntegerTag struct {
	Value int `xml:"gco:Integer"`
}

// OrganisationNameTag struct for XML marshalling.
type OrganisationNameTag struct {
	Anchor          *AnchorTag `xml:"gmx:Anchor,omitempty"`
	CharacterString *string    `xml:"gco:CharacterString,omitempty"`
}

// CodeTag struct for XML marshalling.
type CodeTag struct {
	Anchor AnchorTag `xml:"gmx:Anchor"`
}

// AnchorTag struct for XML marshalling.
type AnchorTag struct {
	Href  string `xml:"xlink:href,attr"`
	Value string `xml:",chardata"`
}
