package app

var HvdEndpoint = "https://op.europa.eu/o/opportal-service/euvoc-download-handler?cellarURI=http%3A%2F%2Fpublications.europa.eu%2Fresource%2Fdistribution%2Fhigh-value-dataset-category%2F20241002-0%2Frdf%2Fskos_core%2Fhigh-value-dataset-category.rdf&fileName=high-value-dataset-category.rdf"
var HvdLocalRDFPath = "./high-value-dataset-category.rdf"
var HvdLocalCSVPath = "./high-value-dataset-category.csv"

var InspireSourceEnpoints = map[string]string{
	"themes": "https://inspire.ec.europa.eu/themes",
	"layers": "https://inspire.ec.europa.eu/layers",
}

func GetInspireSourceNames() []string {
	sourceNames := make([]string, 0, len(InspireSourceEnpoints))
	for name := range InspireSourceEnpoints {
		sourceNames = append(sourceNames, name)
	}
	return sourceNames
}
