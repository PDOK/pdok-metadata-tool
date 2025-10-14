package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	log "github.com/sirupsen/logrus" //nolint:depguard
)

const GetRecordByID = "/?service=CSW&request=GetRecordByID"
const GetRecords = "/?service=CSW&request=GetRecords"

const (
	ContentTypeJSON = "application/json"
	ContentTypeXML  = "application/xml"
)

func preTestSetup() *httptest.Server {
	ngrServer := buildMockWebserverNgr()
	return ngrServer
}

var getCreated = true

func buildMockWebserverNgr() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch url := req.URL.String(); {
		case strings.EqualFold(url, "/geonetwork/srv/dut/info?type=me"):
			writeForbiddenResponse(rw, ContentTypeJSON)
		case strings.EqualFold(url, "/geonetwork/srv/api/records/b4ae622c-6201-49d8-bd2e-f7fce9206a1e/tags"):
			writeOkResponse("./testdata/API_Records_Tags_Inspire.json", rw, ContentTypeJSON)
		case strings.EqualFold(url, "/geonetwork/srv/api/records/689c413e-a057-11f0-8de9-0242ac120002/tags"):
			writeOkResponse("./testdata/API_Records_Tags_Inspire.json", rw, ContentTypeJSON)
		case strings.EqualFold(url, "/geonetwork/srv/api/records/c4bda1aa-d6e6-482c-a6f1-bd519e3202d4/tags"):
			writeOkResponse("./testdata/API_Records_Tags_Empty.json", rw, ContentTypeJSON)
		case strings.EqualFold(url, "/geonetwork/srv/api/records/?metadataType=METADATA&uuidProcessing=REMOVE_AND_REPLACE&publishToAll=false"):
			writeOkResponse("./testdata/API_Records_Tags_Empty.json", rw, ContentTypeJSON)
		case strings.EqualFold(url, "/geonetwork/srv/api/records/?metadataType=METADATA&uuidProcessing=REMOVE_AND_REPLACE&publishToAll=true"):
			writeOkResponse("./testdata/API_Records_Tags_Empty.json", rw, ContentTypeJSON)
		case strings.EqualFold(url, "/geonetwork/srv/api/records/689c413e-a057-11f0-8de9-0242ac120002/tags?id=224342"):
			writeOkResponse("./testdata/API_Records_Tags_Inspire.json", rw, ContentTypeJSON)

		case strings.EqualFold(url, "/geonetwork/srv/api/records/689c413e-a057-11f0-8de9-0242ac120002"):
			if getCreated {
				writeOkResponse("./testdata/nwbwegen222-wms.xml", rw, ContentTypeXML)
			} else {
				writeOkResponse("./testdata/nwbwegen333-wms.xml", rw, ContentTypeXML)
			}

			getCreated = !getCreated
		case strings.HasPrefix(url, GetRecordByID):
			var responsePath string

			switch path := strings.TrimPrefix(url, GetRecordByID); path {
			case "&version=2.0.2&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full&id=C2DFBDBC-5092-11E0-BA8E-B62DE0D72085":
				responsePath = "../../examples/ISO19115/Voorbeeld_Metadata_Dataset_2022_max.xml"
			case "&version=2.0.2&outputSchema=http://www.isotc211.org/2005/gmd&elementSetName=full&id=C2DFBDBC-5092-11E0-BA8E-B62DE0D72086":
				responsePath = "../../examples/ISO19119/Voorbeeld_Metadata_Services_2019_max.xml"
			default:
				log.Infof("no handler for request %s in test setup", req.URL.String())
				rw.WriteHeader(http.StatusNotFound)

				return
			}

			metadataResponse, err := readFileToString(responsePath)
			if err != nil {
				log.Errorf("%v", err)
			}

			rw.Header().Set("Content-Type", "application/xml")
			rw.WriteHeader(http.StatusOK)

			getRecordByIDResponse := wrapAsGetRecordByIDResponse(metadataResponse)

			_, err = fmt.Fprint(rw, getRecordByIDResponse)
			if err != nil {
				log.Errorf("%v", err)
			}
		case strings.HasPrefix(url, GetRecords):
			switch path := strings.TrimPrefix(url, GetRecords); path {
			case "&version=2.0.2&typeNames=gmd:MD_Metadata&resultType=results&startPosition=1&constraintLanguage=CQL_TEXT&constraint_language_version=1.1.0&constraint=type='dataset'":
				writeOkResponse("./testdata/CSW_GetRecordsResponse_Dataset.xml", rw, ContentTypeXML)
			case "&version=2.0.2&typeNames=gmd:MD_Metadata&resultType=results&startPosition=11&constraintLanguage=CQL_TEXT&constraint_language_version=1.1.0&constraint=type='service'":
				writeOkResponse("./testdata/CSW_GetRecordsResponse_Service.xml", rw, ContentTypeXML)
			default:
				log.Infof("no handler for request %s in test setup", req.URL.String())
				rw.WriteHeader(http.StatusNotFound)

				return
			}
		case url == "/":
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				//nolint:forbidigo
				http.Error(rw, "Error reading body", http.StatusInternalServerError)

				return
			}
			defer req.Body.Close()

			requestBody := string(bodyBytes)
			bodyMatched := false

			if strings.Contains(requestBody, "<ogc:PropertyName>dc:type</ogc:PropertyName>") {
				if strings.Contains(requestBody, "<ogc:Literal>service</ogc:Literal>") {
					writeOkResponse(
						"./testdata/CSW_GetRecordsResponse_Service.xml",
						rw,
						ContentTypeXML,
					)

					bodyMatched = true
				} else if strings.Contains(requestBody, "<ogc:Literal>dataset</ogc:Literal>") {
					writeOkResponse("./testdata/CSW_GetRecordsResponse_Dataset.xml", rw, ContentTypeXML)

					bodyMatched = true
				}
			}

			if !bodyMatched {
				log.Infof("no handler for request %s in test setup", req.URL.String())
				rw.WriteHeader(http.StatusNotFound)
			}
		default:
			log.Infof("no handler for request %s in test setup", req.URL.String())
			rw.WriteHeader(http.StatusNotFound)

			return
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

func writeOkResponse(responsePath string, rw http.ResponseWriter, contentType string) {
	response, err := readFileToString(responsePath)
	if err != nil {
		log.Errorf("%v", err)
	}

	rw.Header().Set("Content-Type", contentType)
	rw.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(rw, response)
}

func writeForbiddenResponse(rw http.ResponseWriter, contentType string) {
	rw.Header().Set("Content-Type", contentType)
	rw.Header().Set("Set-Cookie", "XSRF-TOKEN=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
	rw.WriteHeader(http.StatusForbidden)
}
