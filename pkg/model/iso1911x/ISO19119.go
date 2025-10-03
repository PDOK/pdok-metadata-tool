package iso1911x

import (
	"encoding/xml"
)

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

type CharacterStringTag struct {
	CharacterString string `xml:"gco:CharacterString"`
}

type CodeListValueTag struct {
	CodeList      string `xml:"codeList,attr"`
	CodeListValue string `xml:"codeListValue,attr"`
	Value         string `xml:",chardata"`
}

type LanguageTag struct {
	LanguageCode CodeListValueTag `xml:"gmd:LanguageCode"`
}

type CharacterSetTag struct {
	MDCharacterSetCode CodeListValueTag `xml:"gmd:MD_CharacterSetCode"`
}

type HierarchyLevelTag struct {
	MDScopeCode CodeListValueTag `xml:"gmd:MD_ScopeCode"`
}

type ContactTag struct {
	ResponsibleParty ResponsibleParty `xml:"gmd:CI_ResponsibleParty"`
}

type ResponsibleParty struct {
	OrganisationName OrganisationNameTag `xml:"gmd:organisationName"`
	ContactInfo      ContactInfoTag      `xml:"gmd:contactInfo"`
	Role             RoleTag             `xml:"gmd:role"`
}

type OrganisationNameTag struct {
	Anchor AnchorTag `xml:"gmx:Anchor"`
}

type AnchorTag struct {
	Href  string `xml:"xlink:href,attr"`
	Value string `xml:",chardata"`
}

type ContactInfoTag struct {
	Contact ContactDetails `xml:"gmd:CI_Contact"`
}

type ContactDetails struct {
	Address        AddressTag        `xml:"gmd:address"`
	OnlineResource OnlineResourceTag `xml:"gmd:onlineResource"`
}

type AddressTag struct {
	CIAddress CIAddressTag `xml:"gmd:CI_Address"`
}

type CIAddressTag struct {
	Email CharacterStringTag `xml:"gmd:electronicMailAddress"`
}

type OnlineResourceTag struct {
	CIOnlineResource CIOnlineResourceTag `xml:"gmd:CI_OnlineResource"`
}

type CIOnlineResourceTag struct {
	Linkage URLTag `xml:"gmd:linkage"`
}

type URLTag struct {
	URL string `xml:"gmd:URL"`
}

type RoleTag struct {
	CIRoleCode CodeListValueTag `xml:"gmd:CI_RoleCode"`
}

type GraphicOverviewTag struct {
	BrowseGraphic BrowseGraphic `xml:"gmd:MD_BrowseGraphic"`
}

type BrowseGraphic struct {
	FileName        CharacterStringTag  `xml:"gmd:fileName"`
	FileDescription CharacterStringTag  `xml:"gmd:fileDescription"`
	FileType        *CharacterStringTag `xml:"gmd:fileType,omitempty"`
}

type DescriptiveKeywordsTag struct {
	Keywords *MDKeywords `xml:"gmd:MD_Keywords"`
}

type MDKeywords struct {
	Keyword       []KeywordTag    `xml:"gmd:keyword"`
	Type          *KeywordTypeTag `xml:"gmd:type,omitempty"`
	ThesaurusName *Citation       `xml:"gmd:thesaurusName,omitempty"`
}

type KeywordTag struct {
	Anchor          *AnchorTag `xml:"gmx:Anchor,omitempty"`
	CharacterString *string    `xml:"gco:CharacterString,omitempty"`
}

type KeywordTypeTag struct {
	Code CodeListValueTag `xml:"gmd:MD_KeywordTypeCode"`
}

type DateTag struct {
	Date string `xml:"gco:Date"`
}

type IdentificationInfo struct {
	ServiceIdentification ServiceIdentification `xml:"srv:SV_ServiceIdentification"`
}

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

type Citation struct {
	CICitation CICitation `xml:"gmd:CI_Citation"`
}

type CICitation struct {
	Title      TitleTag       `xml:"gmd:title"`
	Dates      []CIDateTag    `xml:"gmd:date"`
	Identifier *IdentifierTag `xml:"gmd:identifier,omitempty"`
}

type TitleTag struct {
	CharacterString *string    `xml:"gco:CharacterString,omitempty"`
	Anchor          *AnchorTag `xml:"gmx:Anchor,omitempty"`
}

type CIDateTag struct {
	CIDate CIDate `xml:"gmd:CI_Date"`
}

type CIDate struct {
	Date     DateTag     `xml:"gmd:date"`
	DateType DateTypeTag `xml:"gmd:dateType"`
}

type DateTypeTag struct {
	CIDateTypeCode CodeListValueTag `xml:"gmd:CI_DateTypeCode"`
}

type IdentifierTag struct {
	MDIdentifier MDIdentifier `xml:"gmd:MD_Identifier"`
}

type MDIdentifier struct {
	Code CodeTag `xml:"gmd:code"`
}

type CodeTag struct {
	Anchor AnchorTag `xml:"gmx:Anchor"`
}

type DistributionInfo struct {
	Distribution Distribution `xml:"gmd:MD_Distribution"`
}

type Distribution struct {
	TransferOptions TransferOptions `xml:"gmd:transferOptions"`
}

type TransferOptions struct {
	DigitalTransferOptions DigitalTransferOptions `xml:"gmd:MD_DigitalTransferOptions"`
}

type DigitalTransferOptions struct {
	Online OnlineResourceWrapper `xml:"gmd:onLine"`
}

type OnlineResourceWrapper struct {
	Resource CIOnlineResource `xml:"gmd:CI_OnlineResource"`
}

type CIOnlineResource struct {
	Linkage     URLTag          `xml:"gmd:linkage"`
	Protocol    *ProtocolTag    `xml:"gmd:protocol,omitempty"`
	Description *DescriptionTag `xml:"gmd:description,omitempty"`
}

type ProtocolTag struct {
	Anchor AnchorTag `xml:"gmx:Anchor"`
}

type DescriptionTag struct {
	Anchor AnchorTag `xml:"gmx:Anchor"`
}

type DataQualityInfo struct {
	DataQuality DataQuality `xml:"gmd:DQ_DataQuality"`
}

type DataQuality struct {
	Scope  ScopeTag    `xml:"gmd:scope"`
	Report []ReportTag `xml:"gmd:report"`
}

type ScopeTag struct {
	Scope ScopeDetails `xml:"gmd:DQ_Scope"`
}

type ScopeDetails struct {
	Level            LevelTag            `xml:"gmd:level"`
	LevelDescription LevelDescriptionTag `xml:"gmd:levelDescription"`
}

type LevelTag struct {
	MDScopeCode CodeListValueTag `xml:"gmd:MD_ScopeCode"`
}

type LevelDescriptionTag struct {
	ScopeDescription ScopeDescriptionTag `xml:"gmd:MD_ScopeDescription"`
}

type ScopeDescriptionTag struct {
	Other CharacterStringTag `xml:"gmd:other"`
}

type ReportTag struct {
	DomainConsistency     *DomainConsistencyTag     `xml:"gmd:DQ_DomainConsistency,omitempty"`
	ConceptualConsistency *ConceptualConsistencyTag `xml:"gmd:DQ_ConceptualConsistency,omitempty"`
}

type DomainConsistencyTag struct {
	Result ConformanceResultTag `xml:"gmd:result"`
}

type ConformanceResultTag struct {
	DQConformanceResult DQConformanceResult `xml:"gmd:DQ_ConformanceResult"`
}

type DQConformanceResult struct {
	Specification Citation           `xml:"gmd:specification"`
	Explanation   CharacterStringTag `xml:"gmd:explanation"`
	Pass          BooleanTag         `xml:"gmd:pass"`
}

type ConceptualConsistencyTag struct {
	NameOfMeasure      AnchorTag          `xml:"gmd:nameOfMeasure"`
	MeasureDescription CharacterStringTag `xml:"gmd:measureDescription"`
	Result             QuantitativeResult `xml:"gmd:result"`
}

type QuantitativeResult struct {
	DQQuantitativeResult DQQuantitativeResult `xml:"gmd:DQ_QuantitativeResult"`
}

type DQQuantitativeResult struct {
	ValueUnit ValueUnitTag `xml:"gmd:valueUnit"`
	Value     RecordTag    `xml:"gmd:value"`
}

type ValueUnitTag struct {
	Href string `xml:"xlink:href,attr"`
}

//type RecordTag struct {
//	Type       string   `xml:"xsi:type,attr"`
//	FloatValue *float64 `xml:",chardata,omitempty"`
//	IntValue   *int     `xml:",chardata,omitempty"`
//}

type RecordTag struct {
	//XMLName xml.Name `xml:"gco:Record"`
	Type  string `xml:"xsi:type,attr"`
	Value string `xml:",chardata"`
}

type BooleanTag struct {
	Value bool `xml:"gco:Boolean"`
}

type ResourceConstraint struct {
	MDConstraints      *MDConstraints      `xml:"gmd:MD_Constraints,omitempty"`
	MDLegalConstraints *MDLegalConstraints `xml:"gmd:MD_LegalConstraints,omitempty"`
}

type MDConstraints struct {
	UseLimitation CharacterStringTag `xml:"gmd:useLimitation"`
}

type MDLegalConstraints struct {
	AccessConstraints []AccessConstraintTag `xml:"gmd:accessConstraints"`
	OtherConstraints  []OtherConstraintTag  `xml:"gmd:otherConstraints"`
}

type AccessConstraintTag struct {
	MDRestrictionCode CodeListValueTag `xml:"gmd:MD_RestrictionCode"`
}

type OtherConstraintTag struct {
	Anchor AnchorTag `xml:"gmx:Anchor"`
}

type ServiceTypeTag struct {
	LocalName LocalNameTag `xml:"gco:LocalName"`
}

type LocalNameTag struct {
	CodeSpace string `xml:"codeSpace,attr"`
	Value     string `xml:",chardata"`
}

type ExtentTag struct {
	EXExtent EXExtentTag `xml:"gmd:EX_Extent"`
}

type EXExtentTag struct {
	GeographicElement GeographicElementTag `xml:"gmd:geographicElement"`
}

type GeographicElementTag struct {
	GeographicBoundingBox GeographicBoundingBoxTag `xml:"gmd:EX_GeographicBoundingBox"`
}

type GeographicBoundingBoxTag struct {
	WestBoundLongitude DecimalTag `xml:"gmd:westBoundLongitude"`
	EastBoundLongitude DecimalTag `xml:"gmd:eastBoundLongitude"`
	SouthBoundLatitude DecimalTag `xml:"gmd:southBoundLatitude"`
	NorthBoundLatitude DecimalTag `xml:"gmd:northBoundLatitude"`
}

type DecimalTag struct {
	Value string `xml:"gco:Decimal"`
}

type CouplingTypeTag struct {
	SVCouplingType CodeListValueTag `xml:"srv:SV_CouplingType"`
}

type OperationMetadataTag struct {
	OperationMetadata SVOperationMetadata `xml:"srv:SV_OperationMetadata"`
}

type SVOperationMetadata struct {
	OperationName OperationNameTag `xml:"srv:operationName"`
	DCP           DCPTag           `xml:"srv:DCP"`
	ConnectPoint  ConnectPointTag  `xml:"srv:connectPoint"`
}

type OperationNameTag struct {
	CharacterString string `xml:"gco:CharacterString"`
}

type DCPTag struct {
	DCPList CodeListValueTag `xml:"srv:DCPList"`
}

type ConnectPointTag struct {
	OnlineResource CIOnlineResource `xml:"gmd:CI_OnlineResource"`
}

type OperatesOn struct {
	UUIDRef string `xml:"uuidref,attr"`
	Href    string `xml:"xlink:href,attr"`
}
