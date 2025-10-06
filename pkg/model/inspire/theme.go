package inspire

// InspireTheme represents an INSPIRE theme with both English and Dutch labels.
//
//nolint:revive
type InspireTheme struct {
	ID           string `json:"id"`           // Primary Key, Unique, 2 characters
	Order        int    `json:"order"`        // Order number
	LabelDutch   string `json:"labelDutch"`   // Dutch label
	LabelEnglish string `json:"labelEnglish"` // English label
	URL          string `json:"url"`          // URL for the theme
}

// InspireThemeRaw represents the raw structure of an INSPIRE theme as received from the API.
//
//nolint:revive
type InspireThemeRaw struct {
	Register struct {
		ContainedItems []struct {
			Theme struct {
				Id          string `json:"id"`          // URL of the theme
				ThemeNumber string `json:"themenumber"` // Theme number/order
				Label       struct {
					Text string `json:"text"` // The actual label text
					Lang string `json:"lang"` // The language of the label
				} `json:"label"` // Theme label in a specific language
			} `json:"theme"`
		} `json:"containeditems"`
	} `json:"register"`
}
