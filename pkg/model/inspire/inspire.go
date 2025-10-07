// Package inspire provides the models used for retrieving INSPIRE data.
package inspire

import (
	"fmt"
)

// InspireEndpoint is the endpoint for INSPIRE registry downloads (currently not used since it's not complete).
const InspireEndpoint = "https://inspire.ec.europa.eu"

// InspireRegisterKind values.
type InspireRegisterKind string

// Values for InspireRegisterKind.
const (
	Theme InspireRegisterKind = "theme"
	Layer InspireRegisterKind = "layer"
)

// InspireRegisterLanguage values.
type InspireRegisterLanguage string

// Values for InspireRegisterLanguage.
const (
	Dutch   InspireRegisterLanguage = "nl"
	English InspireRegisterLanguage = "en"
)

// InspireRegisterKinds holds the INSPIRE registry kinds.
var InspireRegisterKinds = []InspireRegisterKind{Theme, Layer}

// InspireVariant values.
type InspireVariant string

// Values for InspireVariant.
const (
	Harmonised InspireVariant = "HARMONISED"
	AsIs       InspireVariant = "ASIS"
)

// InspireItem is an interface that both InspireLayer and InspireTheme implement.
type InspireItem interface {
	GetId() string
	GetLabelDutch() string
	GetLabelEnglish() string
}

// GetInspireEndpoint returns the INSPIRE endpoint for a given kind and language.
func GetInspireEndpoint(kind InspireRegisterKind, language InspireRegisterLanguage) string {
	return fmt.Sprintf("%s/%s/%s.%s.json", InspireEndpoint, kind, kind, language)
}

// GetInspirePath returns the JSON path for a given kind and language.
func GetInspirePath(kind InspireRegisterKind, language InspireRegisterLanguage) string {
	return fmt.Sprintf("%s.%s.json", kind, language)
}

// GetInspireThemeIDForDutchLabel returns the INSPIRE theme id for a given Dutch label.
//
//nolint:cyclop,funlen
func GetInspireThemeIDForDutchLabel(labelDutch string) (id string) {
	//nolint:misspell
	switch labelDutch {
	case "Administratieve eenheden":
		id = "au"
	case "Adressen":
		id = "ad"
	case "Atmosferische omstandigheden":
		id = "ac"
	case "Beschermde gebieden":
		id = "ps"
	case "Biogeografische gebieden":
		id = "br"
	case "Bodem":
		id = "so"
	case "Bodemgebruik":
		id = "lc"
	case "Energiebronnen":
		id = "er"
	case "Faciliteiten voor landbouw en aquacultuur":
		id = "af"
	case "Faciliteiten voor productive en industrie":
		id = "pf"
	case "Gebieden met natuurrisico's":
		id = "nz"
	case "Gebiedsbeheer, gebieden waar beperkingen gelden, gereguleerde gebieden en rapportage-eenheden":
		id = "am"
	case "Gebouwen":
		id = "bu"
	case "Geografische namen":
		id = "gn"
	case "Geografisch rastersysteem":
		id = "gg"
	case "Geologie":
		id = "ge"
	case "Habitats en biotopen":
		id = "hb"
	case "Hoogte":
		id = "el"
	case "Hydrografie":
		id = "hy"
	case "Kadastrale percelen":
		id = "cp"
	case "Landgebruik":
		id = "lu"
	case "Menselijke gezondheid en veiligheid":
		id = "hh"
	case "Meteorologische geografische kenmerken":
		id = "mf"
	case "Milieubewakingsvoorzieningen":
		id = "ef"
	case "Minerale bronnen":
		id = "mr"
	case "Nutsdiensten en overheidsdiensten":
		id = "us"
	case "Oceanografische geografische kenmerken":
		id = "of"
	case "Orthobeeldvorming":
		id = "oi"
	case "Spreiding van de bevolking — demografie":
		id = "pd"
	case "Spreiding van soorten":
		id = "sd"
	case "Statistische eenheden":
		id = "su"
	case "Systemen voor verwijzing door middel van coördinaten":
		id = "rs"
	case "Vervoersnetwerken":
		id = "tn"
	case "Zeegebieden":
		id = "sr"
	}

	return id
}
