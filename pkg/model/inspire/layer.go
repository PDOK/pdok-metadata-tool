package inspire

// InspireLayer represents an INSPIRE layer with both English and Dutch labels.
type InspireLayer struct {
	ID           string `json:"id"`           // Primary Key, Unique, up to 100 characters
	LabelDutch   string `json:"labelDutch"`   // Dutch label
	LabelEnglish string `json:"labelEnglish"` // English label
}

// InspireLayerRaw represents the raw structure of an INSPIRE layer as received from the API.
type InspireLayerRaw struct {
	Register struct {
		Registry struct {
			Label struct {
				Text string `json:"text"` // The registry label text
				Lang string `json:"lang"` // The language of the label
			} `json:"label"`
			ID string `json:"id"` // The registry ID
		} `json:"registry"`
		Label struct {
			Text string `json:"text"` // The register label text
			Lang string `json:"lang"` // The language of the label
		} `json:"label"`
		ID             string `json:"id"` // The register ID
		ContainedItems []struct {
			Layer struct {
				ID    string `json:"id"` // URL of the layer
				Label struct {
					Text string `json:"text"` // The actual label text
					Lang string `json:"lang"` // The language of the label
				} `json:"label"`
				LayerName struct {
					Text string `json:"text"` // The layer name
					Lang string `json:"lang"` // The language of the layer name
				} `json:"layername"`
				Theme struct {
					Label struct {
						Text string `json:"text"` // The theme label text
						Lang string `json:"lang"` // The language of the label
					} `json:"label"`
					URI string `json:"uri"` // The theme URI
				} `json:"theme"`
				Status struct {
					Label struct {
						Text string `json:"text"` // The status label text
						Lang string `json:"lang"` // The language of the label
					} `json:"label"`
					ID string `json:"id"` // The status ID
				} `json:"status"`
			} `json:"layer"`
		} `json:"containeditems"`
	} `json:"register"`
}
