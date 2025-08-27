package repository

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

func preTestSetup() *httptest.Server {
	ngrServer := buildMockWebserverNgr()
	return ngrServer
}

func buildMockWebserverNgr() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		test := req.URL.String()
		fmt.Print(test)
		switch req.URL.String() {
		case "/?service=CSW&request=GetRecordById&version=2.0.2&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full&id=C2DFBDBC-5092-11E0-BA8E-B62DE0D72085":

			responsePath := "../../examples/ISO19115/Voorbeeld_Metadata_Dataset_2022_max.xml"
			metadataResponse, err := readFileToString(responsePath)
			if err != nil {
				log.Errorf("%v", err)
			}

			rw.Header().Set("Content-Type", "application/xml")
			rw.WriteHeader(http.StatusOK)

			getRecordByIdResponse := wrapAsGetRecordByIdResponse(metadataResponse)
			_, err = fmt.Fprint(rw, getRecordByIdResponse)
			if err != nil {
				log.Errorf("%v", err)
			}
		case "/?service=CSW&request=GetRecords&version=2.0.2&typeNames=gmd:MD_Metadata&resultType=results":

			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
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
					log.Errorf("%v", err)
				}

				rw.Header().Set("Content-Type", "application/xml")
				rw.WriteHeader(http.StatusOK)

				_, err = fmt.Fprint(rw, metadataResponse)
				if err != nil {
					log.Errorf("%v", err)
				}
			}

		default:
			log.Infof("no handler for request %s in test setup", req.URL.String())
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

func wrapAsGetRecordByIdResponse(metadata string) string {
	return `<csw:GetRecordByIdResponse xmlns:csw="http://www.opengis.net/cat/csw/2.0.2">` +
		metadata +
		`</csw:GetRecordByIdResponse>`
}
