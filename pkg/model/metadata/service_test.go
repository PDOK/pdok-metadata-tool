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

// Static expectations built by analysing the XML examples under examples/ISO19119.
func TestNewNLServiceMetadataFromMDMetadataWithHVDRepo_StaticExamples(t *testing.T) {
	type expected struct {
		Metadata NLServiceMetadata
		File     string
	}

	root := common.GetProjectRoot()
	examples := filepath.Join(root, "examples", "ISO19119")

	cases := []expected{
		{
			File: filepath.Join(examples, "39d03482-fef0-4706-8f66-16ffb2617155.xml"),
			Metadata: NLServiceMetadata{
				MetadataID:       "39d03482-fef0-4706-8f66-16ffb2617155",
				Title:            "Gebiedsbeheer eenheden - Kwetsbaar gebied - Agglomeraties - RSA (INSPIRE geharmoniseerd) WFS",
				Abstract:         "Dit is de web feature service van INSPIRE thema Gebiedsbeheer geharmoniseerde agglomeraties zoals gerapporteerd naar de Europese Commissie tbv EU rapportage Stedelijk Afvalwater 2020.",
				OrganisationName: "Beheer PDOK",
				Keywords: []string{
					"Kwetsbaar gebied",
					"Richtlijn 91/271/EEG",
					"31991L0271",
					"Agglomeraties",
				},
				ServiceType: "other",
				OperatesOn: []string{
					"2350b86b-3efd-47e4-883e-519bfa8d0abd",
				},
				Endpoints: []iso1911x.ServiceEndpoint{
					{
						URL:      "https://service.pdok.nl/rws/gebiedsbeheer/kwetsbaargebied-agglomeraties/wfs/v1_0?request=GetCapabilities&service=WFS",
						Protocol: "OGC:WFS",
					},
				},
				LicenceURL:   "https://creativecommons.org/publicdomain/zero/1.0/deed.nl",
				ThumbnailURL: "https://www.nationaalgeoregister.nl/geonetwork/srv/api/records/39d03482-fef0-4706-8f66-16ffb2617155/attachments/map%20(1).png",
				CreationDate: "2024-04-25",
				RevisionDate: "2024-11-22",
				InspireThemes: []string{
					"am",
				},
				HVDCategories: []hvd.HVDCategory{
					{
						ID:         "c_dd313021",
						LabelDutch: "Aardobservatie en milieu",
					},
				},
			},
		},
		{
			File: filepath.Join(examples, "1761ab61-c41d-4897-8ee3-a575e717d765.xml"),
			Metadata: NLServiceMetadata{
				MetadataID:       "1761ab61-c41d-4897-8ee3-a575e717d765",
				Title:            "Luchtfoto Landelijke Voorziening Beeldmateriaal (2015) WMS",
				Abstract:         "De orthofotomozaieken zijn een samenstelling van afzonderlijke orthofoto's, in principe van de centrale gedeelten van iedere orthofoto. Daardoor is de omvalling in de mozaieken zo klein mogelijk gehouden. De orthomozaïeken zijn landsdekkend. Binnen de service worden 3 lagen per jaargang aangeboden, - Hoge resolutie orthofoto onder de naam Ortho10 - is beschikbaar, landsdekkend - Lage resolutie orthofoto (RGB) onder de naam Ortho25 - is beschikbaar, landsdekkend (de opnames van noord-oost Groningen zijn van 2014) - Lage resolutie orthofoto (CIR) onder de naam Ortho25IR - is beschikbaar, landsdekkend (de opnames van noord-oost Groningen zijn van 2014)",
				OrganisationName: "Beheer PDOK",
				Keywords: []string{
					"Landelijke voorziening beeldmateriaal",
					"LVB",
					"Beeldmateriaal",
					"Orthofotomozaiek",
					"Orthofoto",
					"Ortho",
					"Luchtfoto",
					"Mozaiek",
					"Hoge resolutie",
					"HR",
					"Luchtbeelden",
					"Lage resolutie",
					"Infrarood",
				},
				ServiceType: "view",
				OperatesOn: []string{
					"574e27b1-76c4-4e3e-a941-d29e1549e401",
					"18a6fdd8-cc71-4242-8c29-cd9f1d46d7b2",
					"bc695b6e-d6a9-4c90-a04f-6fc87c767857",
				},
				Endpoints: []iso1911x.ServiceEndpoint{
					{
						URL:      "https://secure.geodata2.nationaalgeoregister.nl/lv-beeldmateriaal/2015/wms?",
						Protocol: "OGC:WMS",
					},
				},
				LicenceURL:    "https://creativecommons.org/publicdomain/mark/1.0/deed.nl",
				ThumbnailURL:  "",
				CreationDate:  "",
				RevisionDate:  "2016-01-13",
				InspireThemes: nil,
				HVDCategories: nil,
			},
		},
		{
			File: filepath.Join(examples, "dae8f9e3-99af-4d21-9feb-29f2a1693077.xml"),
			Metadata: NLServiceMetadata{
				MetadataID:       "dae8f9e3-99af-4d21-9feb-29f2a1693077",
				Title:            "Vervoersnetwerken (INSPIRE geharmoniseerd) WMS",
				Abstract:         "INSPIRE transportnetwerken, geharmoniseerd, gevuld met relevante objecten uit TOP10NL (onderdeel van de Basisregistreatie Topografie BRT), geproduceerd en beheerd door het Kadaster.",
				OrganisationName: "Beheer PDOK",
				Keywords: []string{
					"Transport Networks",
				},
				ServiceType: "view",
				OperatesOn: []string{
					"31de946d-85d4-4c93-bb97-e25f4ef1401a",
					"5951efa2-1ff3-4763-a966-a2f5497679ee",
					"6c06740d-058f-4a12-bb3f-bf68efd03d09",
					"31de946d-85d4-4c93-bb97-e25f4ef1401a",
					"31de946d-85d4-4c93-bb97-e25f4ef1401a",
					"3a7dd0a6-d130-4c4c-b0ba-24365cf036e2",
					"3a7dd0a6-d130-4c4c-b0ba-24365cf036e2",
					"5951efa2-1ff3-4763-a966-a2f5497679ee",
					"8f45b8ef-0ce8-463a-9059-5efdcecb785c",
				},
				Endpoints: []iso1911x.ServiceEndpoint{
					{
						URL:      "https://service.pdok.nl/kadaster/tn/wms/v1_0?request=GetCapabilities&service=WMS",
						Protocol: "OGC:WMS",
					},
				},
				LicenceURL:   "http://creativecommons.org/publicdomain/mark/1.0/deed.nl",
				ThumbnailURL: "https://www.nationaalgeoregister.nl/geonetwork/srv/api/records/dae8f9e3-99af-4d21-9feb-29f2a1693077/attachments/vervoers.jpg",
				CreationDate: "2021-12-03",
				RevisionDate: "2025-12-09",
				InspireThemes: []string{
					"tn",
				},
				HVDCategories: []hvd.HVDCategory{
					{
						ID:         "c_b79e35eb",
						LabelDutch: "Mobiliteit",
					},
				},
			},
		},
		{
			File: filepath.Join(examples, "0017219b-fb75-47aa-a6bf-496f2514e545.xml"),
			Metadata: NLServiceMetadata{
				MetadataID:       "0017219b-fb75-47aa-a6bf-496f2514e545",
				Title:            "Aardkundige Waarden - Provincies (INSPIRE geharmoniseerd) ATOM",
				Abstract:         "Deze nationale dataset bevat de Aardkundige waarden. De dataset Aardkundige waarden valt binnen het INSPIRE-thema Beschermde gebieden.",
				OrganisationName: "Beheer PDOK",
				Keywords: []string{
					"Nationaal",
				},
				ServiceType: "download",
				OperatesOn: []string{
					"f002bfc5-7d87-46b6-819e-8415422b65c9",
				},
				Endpoints: []iso1911x.ServiceEndpoint{
					{
						URL:      "https://service.pdok.nl/provincies/aardkundige-waarden/atom",
						Protocol: "INSPIRE Atom",
					},
				},
				LicenceURL:    "https://creativecommons.org/licenses/by/4.0/deed.nl",
				ThumbnailURL:  "https://www.nationaalgeoregister.nl/geonetwork/srv/api/records/0017219b-fb75-47aa-a6bf-496f2514e545/attachments/AardkundigeWaarden.png",
				CreationDate:  "2022-05-12",
				RevisionDate:  "2025-07-14",
				InspireThemes: []string{"ps"},
				HVDCategories: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(filepath.Base(tc.File), func(t *testing.T) {
			b, err := os.ReadFile(tc.File)
			require.NoError(t, err)

			var md iso1911x.MDMetadata
			require.NoError(t, xml.Unmarshal(b, &md)) //nolint

			flat := NewNLServiceMetadataFromMDMetadataWithHVDRepo(&md, nil)
			require.NotNil(t, flat)

			assert.Equal(t, tc.Metadata.MetadataID, flat.MetadataID)
			assert.Equal(t, tc.Metadata.Title, flat.Title)
			assert.Equal(t, tc.Metadata.Abstract, flat.Abstract)
			assert.Equal(t, tc.Metadata.OrganisationName, flat.OrganisationName)
			assert.Equal(t, tc.Metadata.Keywords, flat.Keywords)
			assert.Equal(t, tc.Metadata.ServiceType, flat.ServiceType)
			assert.Equal(t, tc.Metadata.OperatesOn, flat.OperatesOn)
			assert.Equal(t, tc.Metadata.Endpoints, flat.Endpoints)
			assert.Equal(t, tc.Metadata.LicenceURL, flat.LicenceURL)
			assert.Equal(t, tc.Metadata.ThumbnailURL, flat.ThumbnailURL)
			assert.Equal(t, tc.Metadata.CreationDate, flat.CreationDate)
			assert.Equal(t, tc.Metadata.RevisionDate, flat.RevisionDate)
			assert.Equal(t, tc.Metadata.InspireThemes, flat.InspireThemes)

			if tc.Metadata.HVDCategories != nil {
				assert.NotEmpty(t, flat.HVDCategories)
				assert.Equal(t, tc.Metadata.HVDCategories, flat.HVDCategories)
			}
		})
	}
}
