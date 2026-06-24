// Package hvd provides the model for retrieving HVD categories.
package hvd

// HvdEndpoint is the endpoint for the HVD RDF file.
const HvdEndpoint = "https://publications.europa.eu/resource/distribution/high-value-dataset-category/20241002-0/rdf/skos_core/high-value-dataset-category.rdf"

// HVDCategory represents a High Value Dataset category.
type HVDCategory struct {
	ID           string `json:"id"           validate:"required,max=10"` // ID is the primary key, must be unique
	Parent       string `json:"parent"       validate:"max=10"`          // Parent is a foreign key to another HVDCategory.ID
	Order        string `json:"order"        validate:"max=6"`
	LabelDutch   string `json:"labelDutch"`
	LabelEnglish string `json:"labelEnglish"`
}

// CategoryProvider abstracts a provider that can resolve a category by its code.
// Implemented by repository.HVDRepository.
type CategoryProvider interface {
	GetHVDCategoryByCode(code string) (*HVDCategory, error)
}
