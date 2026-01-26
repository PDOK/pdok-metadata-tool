// Package iso1911x provides logic for generating iso1911x metadata.
package iso1911x

import (
	"encoding/xml"
)

// ISO19119 struct for XML marshalling.
type ISO19119 struct {
	XMLName            xml.Name           `xml:"gmd:MD_Metadata"`
	XmlnsGmd           string             `xml:"xmlns:gmd,attr"`
	XmlnsGco           string             `xml:"xmlns:gco,attr"`
	XmlnsSrv           string             `xml:"xmlns:srv,attr"`
	XmlnsGml           string             `xml:"xmlns:gml,attr"`
	XmlnsXsi           string             `xml:"xmlns:xsi,attr"`
	XmlnsXs            string             `xml:"xmlns:xs,attr"`
	XmlnsCsw           string             `xml:"xmlns:csw,attr"`
	XmlnsGmx           string             `xml:"xmlns:gmx,attr"`
	XmlnsGts           string             `xml:"xmlns:gts,attr"`
	XmlnsXlink         string             `xml:"xmlns:xlink,attr"`
	XsiSchemaLocation  string             `xml:"xsi:schemaLocation,attr"`
	FileIdentifier     CharacterStringTag `xml:"gmd:fileIdentifier"`
	Language           LanguageTag        `xml:"gmd:language"`
	CharacterSet       CharacterSetTag    `xml:"gmd:characterSet"`
	HierarchyLevel     HierarchyLevelTag  `xml:"gmd:hierarchyLevel"`
	HierarchyLevelName CharacterStringTag `xml:"gmd:hierarchyLevelName"`
	Contact            ContactTag         `xml:"gmd:contact"`

	DateStamp               DateTag            `xml:"gmd:dateStamp"`
	MetadataStandardName    CharacterStringTag `xml:"gmd:metadataStandardName"`
	MetadataStandardVersion CharacterStringTag `xml:"gmd:metadataStandardVersion"`

	IdentificationInfo IdentificationInfo `xml:"gmd:identificationInfo"`
	DistributionInfo   DistributionInfo   `xml:"gmd:distributionInfo"`
	DataQualityInfo    DataQualityInfo    `xml:"gmd:dataQualityInfo"`
}

// CodeListValueTag struct for XML marshalling.
type CodeListValueTag struct {
	CodeList      string `xml:"codeList,attr"`
	CodeListValue string `xml:"codeListValue,attr"`
	Value         string `xml:",chardata"`
}

// LanguageTag struct for XML marshalling.
type LanguageTag struct {
	LanguageCode CodeListValueTag `xml:"gmd:LanguageCode"`
}

// CharacterSetTag struct for XML marshalling.
type CharacterSetTag struct {
	MDCharacterSetCode CodeListValueTag `xml:"gmd:MD_CharacterSetCode"`
}

// HierarchyLevelTag struct for XML marshalling.
type HierarchyLevelTag struct {
	MDScopeCode CodeListValueTag `xml:"gmd:MD_ScopeCode"`
}

// ContactTag struct for XML marshalling.
type ContactTag struct {
	ResponsibleParty ResponsibleParty `xml:"gmd:CI_ResponsibleParty"`
}

// ResponsibleParty struct for XML marshalling.
type ResponsibleParty struct {
	OrganisationName AnchorOrCharacterStringTag `xml:"gmd:organisationName"`
	ContactInfo      ContactInfoTag             `xml:"gmd:contactInfo"`
	Role             RoleTag                    `xml:"gmd:role"`
}

// ContactInfoTag struct for XML marshalling.
type ContactInfoTag struct {
	Contact ContactDetails `xml:"gmd:CI_Contact"`
}

// ContactDetails struct for XML marshalling.
type ContactDetails struct {
	Address        AddressTag        `xml:"gmd:address"`
	OnlineResource OnlineResourceTag `xml:"gmd:onlineResource"`
}

// AddressTag struct for XML marshalling.
type AddressTag struct {
	CIAddress CIAddressTag `xml:"gmd:CI_Address"`
}

// CIAddressTag struct for XML marshalling.
type CIAddressTag struct {
	Email CharacterStringTag `xml:"gmd:electronicMailAddress"`
}

// OnlineResourceTag struct for XML marshalling.
type OnlineResourceTag struct {
	CIOnlineResource CIOnlineResourceTag `xml:"gmd:CI_OnlineResource"`
}

// CIOnlineResourceTag struct for XML marshalling.
type CIOnlineResourceTag struct {
	Linkage URLTag `xml:"gmd:linkage"`
}

// URLTag struct for XML marshalling.
type URLTag struct {
	URL string `xml:"gmd:URL"`
}

// RoleTag struct for XML marshalling.
type RoleTag struct {
	CIRoleCode CodeListValueTag `xml:"gmd:CI_RoleCode"`
}

// GraphicOverviewTag struct for XML marshalling.
type GraphicOverviewTag struct {
	BrowseGraphic BrowseGraphic `xml:"gmd:MD_BrowseGraphic"`
}

// BrowseGraphic struct for XML marshalling.
type BrowseGraphic struct {
	FileName        CharacterStringTag  `xml:"gmd:fileName"`
	FileDescription CharacterStringTag  `xml:"gmd:fileDescription"`
	FileType        *CharacterStringTag `xml:"gmd:fileType,omitempty"`
}

// DescriptiveKeywordsTag struct for XML marshalling.
type DescriptiveKeywordsTag struct {
	Keywords *MDKeywords `xml:"gmd:MD_Keywords"`
}

// MDKeywords struct for XML marshalling.
type MDKeywords struct {
	Keyword       []KeywordTag    `xml:"gmd:keyword"`
	Type          *KeywordTypeTag `xml:"gmd:type,omitempty"`
	ThesaurusName *Citation       `xml:"gmd:thesaurusName,omitempty"`
}

// KeywordTag struct for XML marshalling.
type KeywordTag struct {
	Anchor          *AnchorTag `xml:"gmx:Anchor,omitempty"`
	CharacterString *string    `xml:"gco:CharacterString,omitempty"`
}

// KeywordTypeTag struct for XML marshalling.
type KeywordTypeTag struct {
	Code CodeListValueTag `xml:"gmd:MD_KeywordTypeCode"`
}

// IdentificationInfo struct for XML marshalling.
type IdentificationInfo struct {
	ServiceIdentification ServiceIdentification `xml:"srv:SV_ServiceIdentification"`
}

// ServiceIdentification struct for XML marshalling.
type ServiceIdentification struct {
	Citation            Citation                 `xml:"gmd:citation"`
	Abstract            CharacterStringTag       `xml:"gmd:abstract"`
	PointOfContact      ContactTag               `xml:"gmd:pointOfContact"`
	GraphicOverview     []GraphicOverviewTag     `xml:"gmd:graphicOverview"`
	DescriptiveKeywords []DescriptiveKeywordsTag `xml:"gmd:descriptiveKeywords"`
	ResourceConstraints []ResourceConstraint     `xml:"gmd:resourceConstraints"`
	ServiceType         ServiceTypeTag           `xml:"srv:serviceType"`
	Extent              ExtentTag                `xml:"srv:extent"`
	CouplingType        CouplingTypeTag          `xml:"srv:couplingType"`
	ContainsOperations  []OperationMetadataTag   `xml:"srv:containsOperations"`
	OperatesOn          []OperatesOn             `xml:"srv:operatesOn"`
}

// Citation struct for XML marshalling.
type Citation struct {
	CICitation CICitation `xml:"gmd:CI_Citation"`
}

// CICitation struct for XML marshalling.
type CICitation struct {
	Title      TitleTag       `xml:"gmd:title"`
	Dates      []CIDateTag    `xml:"gmd:date"`
	Identifier *IdentifierTag `xml:"gmd:identifier,omitempty"`
}

// TitleTag struct for XML marshalling.
type TitleTag struct {
	CharacterString *string    `xml:"gco:CharacterString,omitempty"`
	Anchor          *AnchorTag `xml:"gmx:Anchor,omitempty"`
}

// CIDateTag struct for XML marshalling.
type CIDateTag struct {
	CIDate CIDate `xml:"gmd:CI_Date"`
}

// CIDate struct for XML marshalling.
type CIDate struct {
	Date     DateTag     `xml:"gmd:date"`
	DateType DateTypeTag `xml:"gmd:dateType"`
}

// DateTypeTag struct for XML marshalling.
type DateTypeTag struct {
	CIDateTypeCode CodeListValueTag `xml:"gmd:CI_DateTypeCode"`
}

// IdentifierTag struct for XML marshalling.
type IdentifierTag struct {
	MDIdentifier MDIdentifier `xml:"gmd:MD_Identifier"`
}

// MDIdentifier struct for XML marshalling.
type MDIdentifier struct {
	Code CodeTag `xml:"gmd:code"`
}

// DistributionInfo struct for XML marshalling.
type DistributionInfo struct {
	Distribution Distribution `xml:"gmd:MD_Distribution"`
}

// Distribution struct for XML marshalling.
type Distribution struct {
	TransferOptions TransferOptions `xml:"gmd:transferOptions"`
}

// TransferOptions struct for XML marshalling.
type TransferOptions struct {
	DigitalTransferOptions DigitalTransferOptions `xml:"gmd:MD_DigitalTransferOptions"`
}

// DigitalTransferOptions struct for XML marshalling.
type DigitalTransferOptions struct {
	Online OnlineResourceWrapper `xml:"gmd:onLine"`
}

// OnlineResourceWrapper struct for XML marshalling.
type OnlineResourceWrapper struct {
	Resource CIOnlineResource `xml:"gmd:CI_OnlineResource"`
}

// CIOnlineResource struct for XML marshalling.
type CIOnlineResource struct {
	Linkage     URLTag          `xml:"gmd:linkage"`
	Protocol    *ProtocolTag    `xml:"gmd:protocol,omitempty"`
	Description *DescriptionTag `xml:"gmd:description,omitempty"`
}

// ProtocolTag struct for XML marshalling.
type ProtocolTag struct {
	Anchor AnchorTag `xml:"gmx:Anchor"`
}

// DescriptionTag struct for XML marshalling.
type DescriptionTag struct {
	Anchor AnchorTag `xml:"gmx:Anchor"`
}

// DataQualityInfo struct for XML marshalling.
type DataQualityInfo struct {
	DataQuality DataQuality `xml:"gmd:DQ_DataQuality"`
}

// DataQuality struct for XML marshalling.
type DataQuality struct {
	Scope  ScopeTag    `xml:"gmd:scope"`
	Report []ReportTag `xml:"gmd:report"`
}

// ScopeTag struct for XML marshalling.
type ScopeTag struct {
	Scope ScopeDetails `xml:"gmd:DQ_Scope"`
}

// ScopeDetails struct for XML marshalling.
type ScopeDetails struct {
	Level            LevelTag            `xml:"gmd:level"`
	LevelDescription LevelDescriptionTag `xml:"gmd:levelDescription"`
}

// LevelTag struct for XML marshalling.
type LevelTag struct {
	MDScopeCode CodeListValueTag `xml:"gmd:MD_ScopeCode"`
}

// LevelDescriptionTag struct for XML marshalling.
type LevelDescriptionTag struct {
	ScopeDescription ScopeDescriptionTag `xml:"gmd:MD_ScopeDescription"`
}

// ScopeDescriptionTag struct for XML marshalling.
type ScopeDescriptionTag struct {
	Other CharacterStringTag `xml:"gmd:other"`
}

// ReportTag struct for XML marshalling.
type ReportTag struct {
	DomainConsistency     *DomainConsistencyTag     `xml:"gmd:DQ_DomainConsistency,omitempty"`
	ConceptualConsistency *ConceptualConsistencyTag `xml:"gmd:DQ_ConceptualConsistency,omitempty"`
}

// DomainConsistencyTag struct for XML marshalling.
type DomainConsistencyTag struct {
	Result ConformanceResultTag `xml:"gmd:result"`
}

// ConformanceResultTag struct for XML marshalling.
type ConformanceResultTag struct {
	DQConformanceResult DQConformanceResult `xml:"gmd:DQ_ConformanceResult"`
}

// DQConformanceResult struct for XML marshalling.
type DQConformanceResult struct {
	Specification Citation           `xml:"gmd:specification"`
	Explanation   CharacterStringTag `xml:"gmd:explanation"`
	Pass          BooleanTag         `xml:"gmd:pass"`
}

// ConceptualConsistencyTag struct for XML marshalling.
type ConceptualConsistencyTag struct {
	NameOfMeasure      NameOfMeasureTag   `xml:"gmd:nameOfMeasure"`
	MeasureDescription CharacterStringTag `xml:"gmd:measureDescription"`
	Result             QuantitativeResult `xml:"gmd:result"`
}

// NameOfMeasureTag struct for XML marshalling.
type NameOfMeasureTag struct {
	Anchor AnchorTag `xml:"gmx:Anchor"`
}

// QuantitativeResult struct for XML marshalling.
type QuantitativeResult struct {
	DQQuantitativeResult DQQuantitativeResult `xml:"gmd:DQ_QuantitativeResult"`
}

// DQQuantitativeResult struct for XML marshalling.
type DQQuantitativeResult struct {
	ValueUnit ValueUnitTag `xml:"gmd:valueUnit"`
	Value     ValueTag     `xml:"gmd:value"`
}

// ValueUnitTag struct for XML marshalling.
type ValueUnitTag struct {
	Href string `xml:"xlink:href,attr"`
}

// ValueTag struct for XML marshalling.
type ValueTag struct {
	Record RecordTag `xml:"gco:Record"`
}

// RecordTag struct for XML marshalling.
type RecordTag struct {
	Type  string `xml:"xsi:type,attr"`
	Value string `xml:",chardata"`
}

// ResourceConstraint struct for XML marshalling.
type ResourceConstraint struct {
	MDConstraints      *MDConstraints      `xml:"gmd:MD_Constraints,omitempty"`
	MDLegalConstraints *MDLegalConstraints `xml:"gmd:MD_LegalConstraints,omitempty"`
}

// MDConstraints struct for XML marshalling.
type MDConstraints struct {
	UseLimitation CharacterStringTag `xml:"gmd:useLimitation"`
}

// MDLegalConstraints struct for XML marshalling.
type MDLegalConstraints struct {
	AccessConstraints []AccessConstraintTag `xml:"gmd:accessConstraints"`
	OtherConstraints  []OtherConstraintTag  `xml:"gmd:otherConstraints"`
}

// AccessConstraintTag struct for XML marshalling.
type AccessConstraintTag struct {
	MDRestrictionCode CodeListValueTag `xml:"gmd:MD_RestrictionCode"`
}

// OtherConstraintTag struct for XML marshalling.
type OtherConstraintTag struct {
	Anchor AnchorTag `xml:"gmx:Anchor"`
}

// ServiceTypeTag struct for XML marshalling.
type ServiceTypeTag struct {
	LocalName LocalNameTag `xml:"gco:LocalName"`
}

// LocalNameTag struct for XML marshalling.
type LocalNameTag struct {
	CodeSpace string `xml:"codeSpace,attr"`
	Value     string `xml:",chardata"`
}

// ExtentTag struct for XML marshalling.
type ExtentTag struct {
	EXExtent EXExtentTag `xml:"gmd:EX_Extent"`
}

// EXExtentTag struct for XML marshalling.
type EXExtentTag struct {
	GeographicElement GeographicElementTag `xml:"gmd:geographicElement"`
}

// GeographicElementTag struct for XML marshalling.
type GeographicElementTag struct {
	GeographicBoundingBox GeographicBoundingBoxTag `xml:"gmd:EX_GeographicBoundingBox"`
}

// GeographicBoundingBoxTag struct for XML marshalling.
type GeographicBoundingBoxTag struct {
	WestBoundLongitude DecimalTag `xml:"gmd:westBoundLongitude"`
	EastBoundLongitude DecimalTag `xml:"gmd:eastBoundLongitude"`
	SouthBoundLatitude DecimalTag `xml:"gmd:southBoundLatitude"`
	NorthBoundLatitude DecimalTag `xml:"gmd:northBoundLatitude"`
}

// DecimalTag struct for XML marshalling.
type DecimalTag struct {
	Value string `xml:"gco:Decimal"`
}

// CouplingTypeTag struct for XML marshalling.
type CouplingTypeTag struct {
	SVCouplingType CodeListValueTag `xml:"srv:SV_CouplingType"`
}

// OperationMetadataTag struct for XML marshalling.
type OperationMetadataTag struct {
	OperationMetadata SVOperationMetadata `xml:"srv:SV_OperationMetadata"`
}

// SVOperationMetadata struct for XML marshalling.
type SVOperationMetadata struct {
	OperationName OperationNameTag `xml:"srv:operationName"`
	DCP           DCPTag           `xml:"srv:DCP"`
	ConnectPoint  ConnectPointTag  `xml:"srv:connectPoint"`
}

// OperationNameTag struct for XML marshalling.
type OperationNameTag struct {
	CharacterString string `xml:"gco:CharacterString"`
}

// DCPTag struct for XML marshalling.
type DCPTag struct {
	DCPList CodeListValueTag `xml:"srv:DCPList"`
}

// ConnectPointTag struct for XML marshalling.
type ConnectPointTag struct {
	OnlineResource CIOnlineResource `xml:"gmd:CI_OnlineResource"`
}

// OperatesOn struct for XML marshalling.
type OperatesOn struct {
	UUIDRef string `xml:"uuidref,attr"`
	Href    string `xml:"xlink:href,attr"`
}
