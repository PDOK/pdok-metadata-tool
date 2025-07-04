package model

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
