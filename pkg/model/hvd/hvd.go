package hvd

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
