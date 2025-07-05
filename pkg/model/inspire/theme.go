package inspire

// InspireTheme represents an INSPIRE theme with both English and Dutch labels
type InspireTheme struct {
	Id           string `json:"id"`           // Primary Key, Unique, 2 characters
	Order        int    `json:"order"`        // Order number
	LabelDutch   string `json:"labelDutch"`   // Dutch label
	LabelEnglish string `json:"labelEnglish"` // English label
	URL          string `json:"url"`          // URL for the theme
}

// GetId returns the ID of the theme
func (t InspireTheme) GetId() string {
	return t.Id
}

// GetLabelDutch returns the Dutch label of the theme
func (t InspireTheme) GetLabelDutch() string {
	return t.LabelDutch
}

// GetLabelEnglish returns the English label of the theme
func (t InspireTheme) GetLabelEnglish() string {
	return t.LabelEnglish
}
