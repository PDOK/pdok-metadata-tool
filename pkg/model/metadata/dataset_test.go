package metadata

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/pdok/pdok-metadata-tool/v2/internal/common"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/v2/pkg/model/iso1911x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Static expectations built by analysing the XML examples under examples/ISO19115.
func TestNewNLDatasetMetadataFromMDMetadataWithHVDRepo_StaticExamples(t *testing.T) {
	type expected struct {
		Metadata NLDatasetMetadata
		File     string
	}

	root := common.GetProjectRoot()
	examples := filepath.Join(root, "examples", "ISO19115")

	cases := []expected{
		{
			File: filepath.Join(examples, "500d396f-5ec6-4e4b-a151-5fb3cddd8082.xml"),
			Metadata: NLDatasetMetadata{
				MetadataID:              "500d396f-5ec6-4e4b-a151-5fb3cddd8082",
				SourceID:                "440c4a06-6924-4f9c-a9e2-6f61340f711b",
				Title:                   "Gemeten Zwaveldioxide concentraties in buitenlucht.",
				Abstract:                "Ruwe ongevalideerde uurwaarden zwaveldioxide (SO2) op grondniveau in de buitenlucht gemeten in het Landelijk Meetnet Luchtkwaliteit (LML). Zwaveldioxide is een kleurloos gas. Het wordt voornamelijk gevormd het gebruik van zwavelhoudende brandstoffen. Belangrijke bronnen zijn kolengestookte energiecentrales, raffinaderijen en het verkeer (de laatste jaren is voornamelijk de internationale scheepvaart van belang). De concentraties zwaveldioxide zijn in Nederland sterk gedaald door maatregelen op de belangrijkste bronnen. Sinds de jaren 90 van de vorige eeuw zijn er geen normoverschrijdingen meer geweest. Bij hoge concentraties heeft zwaveldioxide negatieve effecten op de menselijke gezondheid en draagt het bij aan de verzuring van ecosystemen. Zwaveldioxide wordt in de lucht gedeeltelijk omgezet in sulfaatdeeltjes en heeft zo een bijdrage aan fijn stof.",
				ContactName:             "",
				ContactEmail:            "geodata@rivm.nl",
				ContactURL:              "",
				ContactOrganisationName: "RIVM",
				Keywords: []string{
					"Zwaveldioxide",
					"Vegetatie",
					"Verzuring",
					"SO2",
					"Luchtkwaliteit",
					"Landelijk Meetnet Luchtkwalteit",
					"LML",
					"Buitenlucht",
					"Kwaliteitsmetingen en modelleringsgegevens (Richtlijn Luchtkwaliteit)",
				},
				LicenceURL:     "http://creativecommons.org/publicdomain/mark/1.0/deed.nl",
				UseLimitation:  "Geen",
				ThumbnailURL:   "http://inspire.rivm.nl/sos/eaq/#map",
				InspireVariant: "ASIS",
				InspireThemes: []string{
					"ef",
					"hh",
				},
				HVDCategories: nil,
				CreationDate:  "2015-11-01",
				BoundingBox: &BoundingBox{
					WestBoundLongitude: "3.37087",
					EastBoundLongitude: "7.21097",
					SouthBoundLatitude: "50.7539",
					NorthBoundLatitude: "53.4658",
				},
				DigitalTransferOptions: []OnLine{
					{
						URL:      "",
						Protocol: "",
					},
					{
						URL:      "https://inspire.rivm.nl/sos/eaq/service?service=AQD&version=1.0.0&request=GetCapabilities",
						Protocol: "OGC:SOS",
					},
				},
			},
		},
		{
			File: filepath.Join(examples, "5951efa2-1ff3-4763-a966-a2f5497679ee.xml"),
			Metadata: NLDatasetMetadata{
				MetadataID:              "5951efa2-1ff3-4763-a966-a2f5497679ee",
				SourceID:                "2482250f-3b00-4439-9f93-f3118229b226",
				Title:                   "Vervoersnetwerken: Waterwegen - Transport Networks: Water (INSPIRE geharmoniseerd)",
				Abstract:                "INSPIRE Vervoersnetwerken: Waterwegen (Transport Networks: Water) themalaag, geharmoniseerd, gevuld met relevante objecten uit TOP10NL (onderdeel van de Basisregistratie Topografie BRT), geproduceerd en beheerd door het Kadaster.",
				ContactName:             "Klantcontactcenter",
				ContactEmail:            "kcc@kadaster.nl",
				ContactURL:              "https://www.kadaster.nl",
				ContactOrganisationName: "Kadaster",
				Keywords: []string{
					"vervoersnetwerken",
					"waterwegen",
					"transport networks",
					"water",
					"transport",
					"haven",
					"veerverbinding",
					"Nationaal",
				},
				LicenceURL:     "http://creativecommons.org/publicdomain/mark/1.0/deed.nl",
				UseLimitation:  "Geen gebruiksbeperkingen",
				ThumbnailURL:   "https://github.com/kadaster/top10nl/raw/master/TOP10NL.JPG",
				InspireVariant: "HARMONISED",
				InspireThemes: []string{
					"tn",
				},
				HVDCategories: []hvd.HVDCategory{
					{
						ID:         "c_b79e35eb",
						LabelDutch: "Mobiliteit",
					},
				},
				CreationDate: "",
				BoundingBox: &BoundingBox{
					WestBoundLongitude: "3.30",
					EastBoundLongitude: "7.24",
					SouthBoundLatitude: "50.73",
					NorthBoundLatitude: "53.60",
				},
			},
		},
		{
			File: filepath.Join(examples, "a90027f8-7323-45d6-86a7-9374d0de05bf.xml"),
			Metadata: NLDatasetMetadata{
				MetadataID:              "a90027f8-7323-45d6-86a7-9374d0de05bf",
				SourceID:                "948874aa-c599-4c0f-b0c2-e6b357e73566",
				Title:                   "Emissies naar het riool vanuit de industrie (2019 - heden) (INSPIRE)",
				Abstract:                "Emissies naar het riool vanuit de industrie worden via het e-MJV (elektronisch Milieujaarverslag) gerapporteerd wanneer bedrijven bepaalde drempelwaarden overschrijden, zoals vastgelegd in het EPRTR-protocol (European Pollutant Release and Transfer Register). Bij lozingen op het riool gaat het om stoffen die via industriële processen in het bedrijfsafvalwater terechtkomen en via het gemeentelijk riool naar een rioolwaterzuiveringsinstallatie (RWZI) worden afgevoerd. Bedrijven moeten deze emissies rapporteren als ze onder de reikwijdte van de E-PRTR-verordening vallen én als de emissies van bepaalde stoffen boven de rapportagedrempels uitkomen.",
				ContactName:             "",
				ContactEmail:            "emissieregistratie@rivm.nl",
				ContactURL:              "",
				ContactOrganisationName: "Rijksinstituut voor Volksgezondheid en Milieu",
				Keywords: []string{
					"Nationaal",
					"Emissies (Richtlijn Industriële emissies)",
					"Inrichtingen (Europees register inzake uitstoot en overbrenging van verontreinigende stoffen)",
					"Verordening (EG) 166/2006",
					"menselijke gezondheid",
					"milieubeleid",
					"Trefwoord zonder thesaurus",
					"Tweede trefwoord zonder thesaurus",
				},
				LicenceURL:     "https://creativecommons.org/publicdomain/mark/1.0/deed.nl",
				UseLimitation:  "Geen beperkingen",
				ThumbnailURL:   "URL naar voorbeeldweergave van de dataset", // This is legit the value they filled in
				InspireVariant: "ASIS",
				InspireThemes: []string{
					"us",
					"pf",
				},
				HVDCategories: []hvd.HVDCategory{
					{
						ID:         "c_dd313021",
						LabelDutch: "Aardobservatie en milieu",
					},
					{
						ID:         "c_4ba9548e",
						LabelDutch: "Emissies",
					},
				},
				CreationDate: "2025-02-06",
				BoundingBox: &BoundingBox{
					WestBoundLongitude: "3.37",
					EastBoundLongitude: "7.21",
					SouthBoundLatitude: "50.75",
					NorthBoundLatitude: "53.47",
				},
			},
		},
		{
			File: filepath.Join(examples, "F646DFB9-5BF6-EAB9-042B-CAB6FF2DC275.xml"),
			Metadata: NLDatasetMetadata{
				MetadataID:              "F646DFB9-5BF6-EAB9-042B-CAB6FF2DC275",
				SourceID:                "23c5bc1b-212b-49b5-8375-846ccabd544d",
				Title:                   "BRO - Digitaal Geologisch Model (DGM) as-is",
				Abstract:                "Het Digitaal Geologisch Model (DGM) is een driedimensionaal lagenmodel van de Nederlandse ondergrond tot een diepte van ongeveer 500 m onder NAP, met lokaal uitschieters tot 1200 m. De ondergrondlagen in dit deel van de ondergrond bestaan hoofdzakelijk uit onverharde sedimenten, waarin de grondsoorten klei, zand, grind en veen voorkomen. De lagen worden op basis van verschillen in lithologie en andere eigenschappen ingedeeld in lithostratigrafische eenheden. DGM is een model van de opbouw en de samenhang (geometrie) van deze lithostratigrafische eenheden. De hoogteligging van de onder- en bovenkant en de dikte van de eenheden worden vastgelegd in gridbestanden (rasters) met een celgrootte van 100 bij 100 m. Behalve de laaginformatie bevat DGM ook de geïnterpreteerde boorbeschrijvingen die bij het maken van het model gebruikt zijn. Het modelgebied van DGM bestaat uit het vasteland van Nederland. De ondergrond van het Nederlandse deel van het Continentaal Plat is niet in DGM opgenomen. DGM is een regionaal model. Het is niet geschikt voor gebruik op lokale schaal; voor het maken van een lokaal ondergrondmodel zullen altijd aanvullende gegevens nodig zijn. Voor verdere informatie wordt verwezen naar de website van de BRO: https://basisregistratieondergrond.nl/",
				ContactName:             "",
				ContactEmail:            "support@broservicedesk.nl",
				ContactURL:              "https://www.basisregistratieondergrond.nl",
				ContactOrganisationName: "TNO Geologische Dienst Nederland",
				Keywords: []string{
					"Digitaal Geologisch Model",
					"DGM",
					"humanGeographicViewer",
					"Boringen",
					"Formatie",
					"Nederland",
					"Bodem",
					"basisset NOVEX",
					"Nationaal",
				},
				LicenceURL:     "http://creativecommons.org/publicdomain/zero/1.0/deed.nl",
				UseLimitation:  "Geen gebruiksbeperkingen",
				ThumbnailURL:   "",
				InspireVariant: "ASIS",
				InspireThemes: []string{
					"ge",
				},
				HVDCategories: []hvd.HVDCategory{
					{
						ID:         "c_dd313021",
						LabelDutch: "Aardobservatie en milieu",
					},
					{
						ID:         "c_e3f55603",
						LabelDutch: "Geologie",
					},
				},
				CreationDate: "",
				BoundingBox: &BoundingBox{
					WestBoundLongitude: "3.358",
					EastBoundLongitude: "7.227",
					SouthBoundLatitude: "50.750",
					NorthBoundLatitude: "53.576",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(filepath.Base(tc.File), func(t *testing.T) {
			b, err := os.ReadFile(tc.File)
			require.NoError(t, err)

			var md iso1911x.MDMetadata
			require.NoError(t, xml.Unmarshal(b, &md)) //nolint

			flat := NewNLDatasetMetadataFromMDMetadataWithHVDRepo(&md, nil)
			require.NotNil(t, flat)

			assert.Equal(t, tc.Metadata.MetadataID, flat.MetadataID)
			assert.Equal(t, tc.Metadata.SourceID, flat.SourceID)
			assert.Equal(t, tc.Metadata.Title, flat.Title)
			assert.Equal(t, tc.Metadata.Abstract, flat.Abstract)
			assert.Equal(t, tc.Metadata.ContactName, flat.ContactName)
			assert.Equal(t, tc.Metadata.ContactEmail, flat.ContactEmail)
			assert.Equal(t, tc.Metadata.ContactURL, flat.ContactURL)
			assert.Equal(t, tc.Metadata.Keywords, flat.Keywords)
			assert.Equal(t, tc.Metadata.LicenceURL, flat.LicenceURL)
			assert.Equal(t, tc.Metadata.UseLimitation, flat.UseLimitation)
			assert.Equal(t, tc.Metadata.ThumbnailURL, flat.ThumbnailURL)
			assert.Equal(t, tc.Metadata.InspireVariant, flat.InspireVariant)
			assert.Equal(t, tc.Metadata.CreationDate, flat.CreationDate)
			assert.Equal(t, tc.Metadata.InspireThemes, flat.InspireThemes)

			if tc.Metadata.HVDCategories != nil {
				assert.NotEmpty(t, flat.HVDCategories)
				assert.Equal(t, tc.Metadata.HVDCategories, flat.HVDCategories)
			}

			if tc.Metadata.BoundingBox != nil {
				if assert.NotNil(t, flat.BoundingBox) {
					assert.Equal(
						t,
						tc.Metadata.BoundingBox.WestBoundLongitude,
						flat.BoundingBox.WestBoundLongitude,
					)
					assert.Equal(
						t,
						tc.Metadata.BoundingBox.EastBoundLongitude,
						flat.BoundingBox.EastBoundLongitude,
					)
					assert.Equal(
						t,
						tc.Metadata.BoundingBox.SouthBoundLatitude,
						flat.BoundingBox.SouthBoundLatitude,
					)
					assert.Equal(
						t,
						tc.Metadata.BoundingBox.NorthBoundLatitude,
						flat.BoundingBox.NorthBoundLatitude,
					)
				}
			}

			// Test temp fields for validation script
			assert.Equal(t, tc.Metadata.ContactOrganisationName, flat.ContactOrganisationName)

			if tc.Metadata.DigitalTransferOptions != nil {
				assert.NotEmpty(t, flat.DigitalTransferOptions)
				assert.Equal(t, tc.Metadata.DigitalTransferOptions, flat.DigitalTransferOptions)
			}
		})
	}
}
