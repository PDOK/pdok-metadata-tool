package repository

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

const GetRecordByID = "/?service=CSW&request=GetRecordById"

func preTestSetup() *httptest.Server {
	ngrServer := buildMockWebserverNgr()

	return ngrServer
}

func buildMockWebserverNgr() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch url := req.URL.String(); {
		case strings.HasPrefix(url, GetRecordByID):
			responsePath := ""

			switch path := strings.TrimPrefix(url, GetRecordByID); path {
			case "&version=2.0.2&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full&id=C2DFBDBC-5092-11E0-BA8E-B62DE0D72085":
				responsePath = "../../examples/ISO19115/Voorbeeld_Metadata_Dataset_2022_max.xml"
			case "&version=2.0.2&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full&id=3703b249-a0eb-484e-ba7a-10e31a55bcec":
				responsePath = "../../examples/ISO19115/Invasieve_Exoten_INSPIRE_geharmoniseerd.xml"
			case "&version=2.0.2&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full&id=07575774-57a1-4419-bab4-6c88fdeb02b2":
				responsePath = "../../examples/ISO19115/Waterschappen_Hydrografie_INSPIRE_geharmoniseerd.xml"
			case "&version=2.0.2&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full&id=19165027-a13a-4c19-9013-ec1fd191019d":
				responsePath = "../../examples/ISO19115/Wetlands_INSPIRE_geharmoniseerd.xml"
			default:
				slog.Info("no handler for request in test setup", "url", req.URL.String())
				rw.WriteHeader(http.StatusNotFound)
			}

			metadataResponse, err := readFileToString(responsePath)
			if err != nil {
				slog.Error("error reading file", "err", err)
			}

			rw.Header().Set("Content-Type", "application/xml")
			rw.WriteHeader(http.StatusOK)

			getRecordByIDResponse := wrapAsGetRecordByIDResponse(metadataResponse)

			_, _ = fmt.Fprint(rw, getRecordByIDResponse)
		case url == "/":
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				//nolint:forbidigo
				http.Error(rw, "Error reading body", http.StatusInternalServerError)

				return
			}
			defer req.Body.Close()

			requestBody := string(bodyBytes)

			if strings.Contains(requestBody, "<ogc:PropertyName>dc:type</ogc:PropertyName>") &&
				strings.Contains(requestBody, "<ogc:Literal>dataset</ogc:Literal>") {
				responsePath := "../client/testdata/CSW_GetRecordsResponse_Dataset.xml"

				metadataResponse, err := readFileToString(responsePath)
				if err != nil {
					slog.Error("error reading file", "err", err)
				}

				rw.Header().Set("Content-Type", "application/xml")
				rw.WriteHeader(http.StatusOK)

				_, _ = fmt.Fprint(rw, metadataResponse)
			}
		default:
			slog.Info("no handler for request in test setup", "url", req.URL.String())
			rw.WriteHeader(http.StatusNotFound)
		}
	}))
}

func readFileToString(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func wrapAsGetRecordByIDResponse(metadata string) string {
	return `<csw:GetRecordByIdResponse xmlns:csw="http://www.opengis.net/cat/csw/2.0.2">` +
		metadata +
		`</csw:GetRecordByIdResponse>`
}
