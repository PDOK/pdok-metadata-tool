package iso1911x

import "encoding/xml"

// ISO19110 struct for XML marshalling.
type ISO19110 struct {
	XMLName           xml.Name `xml:"gfc:FC_FeatureCatalogue"`
	XmlnsGfc          string   `xml:"xmlns:gfc,attr"`
	XmlnsGco          string   `xml:"xmlns:gco,attr"`
	XmlnsGmd          string   `xml:"xmlns:gmd,attr"`
	XmlnsGmx          string   `xml:"xmlns:gmx,attr"`
	XmlnsXsi          string   `xml:"xmlns:xsi,attr"`
	XmlnsXlink        string   `xml:"xmlns:xlink,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`

	Name               CharacterStringTag  `xml:"gfc:name"`
	Scope              *CharacterStringTag `xml:"gfc:scope,omitempty"`
	FieldOfApplication *CharacterStringTag `xml:"gfc:fieldOfApplication,omitempty"`
	VersionNumber      CharacterStringTag  `xml:"gfc:versionNumber"`
	VersionDate        DateTag             `xml:"gfc:versionDate"`
	Producer           ProducerTag         `xml:"gfc:producer"`
	FeatureType        FeatureTypeTag      `xml:"gfc:featureType"`
}

// ProducerTag struct for XML marshalling.
type ProducerTag struct {
	CIResponsibleParty CIResponsibleParty `xml:"gmd:CI_ResponsibleParty"`
}

// CIResponsibleParty struct for XML marshalling.
type CIResponsibleParty struct {
	OrganisationName OrganisationNameTag `xml:"gmd:organisationName"`
	Role             RoleTag             `xml:"gmd:role"`
}

// FeatureTypeTag struct for XML marshalling.
type FeatureTypeTag struct {
	FeatureType FeatureType `xml:"gfc:FC_FeatureType"`
}

// FeatureType struct for XML marshalling.
type FeatureType struct {
	TypeName                 TypeNameTag                 `xml:"gfc:typeName"`
	Code                     *CodeTag                    `xml:"gfc:code,omitempty"`
	Definition               CharacterStringTag          `xml:"gfc:definition"`
	IsAbstract               *BooleanTag                 `xml:"gfc:isAbstract,omitempty"`
	Aliases                  *Aliases                    `xml:"gfc:aliases"`
	FeatureCatalogue         *struct{}                   `xml:"gfc:featureCatalogue"`
	ConstrainedBy            *ConstrainedBy              `xml:"gfc:constrainedBy,omitempty"`
	CarrierOfCharacteristics CarrierOfCharacteristicsTag `xml:"gfc:carrierOfCharacteristics"`
}

// TypeNameTag struct for XML marshalling.
type TypeNameTag struct {
	LocalName string `xml:"gco:LocalName"`
}

// Aliases struct for XML marshalling.
type Aliases struct {
	LocalNameValue []LocalNameValue `xml:"gco:LocalName"`
}

// LocalNameValue struct for XML marshalling.
type LocalNameValue struct {
	Value string `xml:",chardata"`
}

// ConstrainedBy struct for XML marshalling.
type ConstrainedBy struct {
	Constraints []Constraint `xml:"gfc:FC_Constraint"`
}

// Constraint struct for XML marshalling.
type Constraint struct {
	Description CharacterStringTag `xml:"gfc:description"`
}

// CarrierOfCharacteristicsTag  struct for XML marshalling.
type CarrierOfCharacteristicsTag struct {
	FeatureAttributes []FeatureAttribute `xml:"gfc:FC_FeatureAttribute"`
}

// FeatureAttribute struct for XML marshalling.
type FeatureAttribute struct {
	FeatureType          *struct{}          `xml:"gfc:featureType"`
	MemberName           MemberNameTag      `xml:"gfc:memberName"`
	Definition           CharacterStringTag `xml:"gfc:definition"`
	Cardinality          *Cardinality       `xml:"gfc:cardinality,omitempty"`
	ValueMeasurementUnit *UnitDefinition    `xml:"gfc:valueMeasurementUnit,omitempty"`
	ValueType            *ValueTypeTag      `xml:"gfc:valueType,omitempty"`
	ListedValues         *ListedValues      `xml:"gfc:listedValue,omitempty"`
}

// MemberNameTag struct for XML marshalling.
type MemberNameTag struct {
	LocalName string `xml:"gco:LocalName"`
}

// Cardinality struct for XML marshalling.
type Cardinality struct {
	Multiplicity Multiplicity `xml:"gco:Multiplicity"`
}

// Multiplicity struct for XML marshalling.
type Multiplicity struct {
	Range RangeTag `xml:"gco:range"`
}

// RangeTag struct for XML marshalling.
type RangeTag struct {
	MultiplicityRange MultiplicityRange `xml:"gco:MultiplicityRange"`
}

// MultiplicityRange struct for XML marshalling.
type MultiplicityRange struct {
	Lower LowerTag               `xml:"gco:lower"`
	Upper UnlimitedIntegerHolder `xml:"gco:upper"`
}

// LowerTag struct for XML marshalling.
type LowerTag struct {
	Integer IntegerTag `xml:"gco:Integer"`
}

// UnlimitedIntegerHolder struct for XML marshalling.
type UnlimitedIntegerHolder struct {
	Unlimited UnlimitedInteger `xml:"gco:UnlimitedInteger"`
}

// UnlimitedInteger struct for XML marshalling.
type UnlimitedInteger struct {
	IsInfinite bool `xml:"isInfinite,attr"`
	Nil        bool `xml:"xsi:nil,attr,omitempty"`
	Value      *int `xml:",chardata"`
}

// ValueTypeTag struct for XML marshalling.
type ValueTypeTag struct {
	TypeName TypeName `xml:"gco:TypeName"`
}

// TypeName struct for XML marshalling.
type TypeName struct {
	TypeNameWithAName TypeNameWithAName `xml:"gfc:valueType"`
}

// TypeNameWithAName struct for XML marshalling.
type TypeNameWithAName struct {
	AName CharacterStringTag `xml:"gco:aName"`
}

// ListedValues struct for XML marshalling.
type ListedValues struct {
	FCListedValues []ListedValue `xml:"gfc:FC_ListedValue"`
}

// ListedValue struct for XML marshalling.
type ListedValue struct {
	Label      CharacterStringTag `xml:"gfc:label"`
	Code       CharacterStringTag `xml:"gfc:code"`
	Definition CharacterStringTag `xml:"gfc:definition"`
}

// ValueMeasurementUnit struct for XML marshalling.
type ValueMeasurementUnit struct {
	*UnitDefinition `xml:"gml:UnitDefinition"`
}

// UnitDefinition struct for XML marshalling.
type UnitDefinition struct {
	XmlnsGml    string     `xml:"xmlns:gml,attr,omitempty"`
	ID          string     `xml:"gml:id,attr,omitempty"`
	Description *string    `xml:"gml:description"`
	Identifier  Identifier `xml:"gml:identifier"`
}

// Identifier struct for XML marshalling.
type Identifier struct {
	CodeSpace string `xml:"codeSpace,attr,omitempty"`
	Value     string `xml:",chardata"`
}
