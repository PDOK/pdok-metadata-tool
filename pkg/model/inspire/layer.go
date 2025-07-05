package inspire

// InspireLayer represents an INSPIRE layer with both English and Dutch labels
type InspireLayer struct {
	Id           string `json:"id"`           // Primary Key, Unique, up to 100 characters
	LabelDutch   string `json:"labelDutch"`   // Dutch label
	LabelEnglish string `json:"labelEnglish"` // English label
}

// GetId returns the ID of the layer
func (l InspireLayer) GetId() string {
	return l.Id
}

// GetLabelDutch returns the Dutch label of the layer
func (l InspireLayer) GetLabelDutch() string {
	return l.LabelDutch
}

// GetLabelEnglish returns the English label of the layer
func (l InspireLayer) GetLabelEnglish() string {
	return l.LabelEnglish
}
