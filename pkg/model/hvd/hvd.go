package hvd

var HvdEndpoint = "https://op.europa.eu/o/opportal-service/euvoc-download-handler?cellarURI=http%3A%2F%2Fpublications.europa.eu%2Fresource%2Fdistribution%2Fhigh-value-dataset-category%2F20241002-0%2Frdf%2Fskos_core%2Fhigh-value-dataset-category.rdf&fileName=high-value-dataset-category.rdf"

// HVDCategory represents a High Value Dataset category
type HVDCategory struct {
	// Id is the primary key, must be unique
	Id string `json:"id" validate:"required,max=10"`
	// Parent is a foreign key to another HVDCategory.Id
	Parent       string `json:"parent" validate:"max=10"`
	Order        string `json:"order" validate:"max=6"`
	LabelDutch   string `json:"labelDutch"`
	LabelEnglish string `json:"labelEnglish"`
}
