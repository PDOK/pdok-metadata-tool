package client

import (
	"github.com/pdok/pdok-metadata-tool/pkg/model/ngr"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNgrClient_GetRecordTags(t *testing.T) {
	mockedNGRServer := preTestSetup()
	ngrClient := getNgrClient(t, mockedNGRServer)

	type args struct {
		uuid string
	}
	tests := []struct {
		name         string
		args         args
		wantNrOfTags int
		want         ngr.RecordTagsResponse
		wantErr      bool
	}{
		{
			name: "GetRecordTags_INSPIRE_Dataset",
			args: args{
				uuid: "b4ae622c-6201-49d8-bd2e-f7fce9206a1e",
			},
			wantNrOfTags: 1,
			want: ngr.RecordTagsResponse{
				{
					ID:   224342,
					Name: "inspire",
					Label: map[string]string{
						"dut": "Inspire",
						"eng": "Inspire",
					},
				},
			},
			wantErr: false,
		}, {
			name: "GetRecordTags_tagless_Dataset",
			args: args{
				uuid: "c4bda1aa-d6e6-482c-a6f1-bd519e3202d4",
			},
			wantNrOfTags: 0,
			want:         ngr.RecordTagsResponse{},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ngrClient.GetRecordTags(tt.args.uuid)

			assert.Nil(t, err)
			assert.Len(t, got, tt.wantNrOfTags)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNgrClient_updateServiceMetadataRecord(t *testing.T) {
	//categoryId := "224342"
	mockedNGRServer := preTestSetup()
	ngrClient := getNgrClient(t, mockedNGRServer)

	type args struct {
		pathRecordXml string
		uuid          string
		categoryId    *string
		groupId       *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test_01",
			wantErr: false,
			args: args{
				uuid:          "689c413e-a057-11f0-8de9-0242ac120002",
				pathRecordXml: "testdata/nwbwegen222-wms.xml",
				//categoryId:    &categoryId,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataToBeCreated, err := os.ReadFile(tt.args.pathRecordXml)
			assert.Nil(t, err)
			recordToBeCreated := string(dataToBeCreated)
			timeString := time.Now().Format("2006-01-02 15:04:05")
			// create unpublished record
			recordToBeCreated = strings.Replace(recordToBeCreated, "NWB - wegen22222", "NWB - wegen33333 "+timeString, -1)
			err = ngrClient.createOrUpdateServiceMetadataRecord(recordToBeCreated, tt.args.categoryId, tt.args.groupId, &ngrClient.NgrConfig, false)
			assert.Nil(t, err)
			recordCreated, getStatus, err := ngrClient.getRecord(tt.args.uuid, &ngrClient.NgrConfig)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, getStatus)
			assert.Contains(t, recordCreated, timeString)

			// publish record and update title

			//time.Sleep(60 * time.Second)
			//dataToBeUpdated, err := os.ReadFile(tt.args.pathRecordXml)
			//assert.Nil(t, err)
			//recordToBeUpdated := string(dataToBeUpdated)
			//recordToBeUpdated = strings.Replace(recordToBeUpdated, "NWB - wegen33333", "NWB - wegen00000", -1)
			//err = ngrClient.createOrUpdateServiceMetadataRecord(recordToBeUpdated, tt.args.uuid, tt.args.categoryId, tt.args.groupId, &ngrClient.NgrConfig, false)
			//assert.Nil(t, err)
			//recordUpdated, getStatus, err := ngrClient.getRecord(tt.args.uuid, &ngrClient.NgrConfig)
			//assert.Nil(t, err)
			//assert.Equal(t, http.StatusOK, getStatus)
			//assert.Contains(t, recordUpdated, "NWB - wegen00000")

			//time.Sleep(60 * time.Second)
			//deleteStatus, err := ngrClient.deleteRecord(tt.args.uuid, &ngrClient.NgrConfig)
			//assert.Nil(t, err)
			//assert.Equal(t, http.StatusNoContent, deleteStatus)
			//delRecord, getDelStatus, err := ngrClient.getRecord(tt.args.uuid, &ngrClient.NgrConfig)
			//assert.NotNil(t, err)
			//assert.Equal(t, getDelStatus, http.StatusNotFound)
			//assert.Equal(t, "", delRecord)

			addTagStatus, err := ngrClient.addTagToRecord(tt.args.uuid, &ngrClient.NgrConfig, INSPIRE_TAG)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusCreated, addTagStatus)
			tagsList, getStatus, err := ngrClient.getTagsByRecord(tt.args.uuid, &ngrClient.NgrConfig)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, getStatus)
			assert.Contains(t, tagsList, "224342")
		})
	}
}

func getNgrClient(t *testing.T, mockedNGRServer *httptest.Server) *NgrClient {
	//hostURL, err := url.Parse(mockedNGRServer.URL)
	accUrl, err := url.Parse("https://ngr.acceptatie.nationaalgeoregister.nl")
	if err != nil {
		t.Fatalf("Failed to parse mocked NGR server URL: %v", err)
	}

	config := NgrConfig{
		NgrUrl:      accUrl.String(),
		NgrUserName: NGR_USER_NAME,
		NgrPassword: NGR_PASSWORD,
	}
	ngrClient := NewNgrClient(config)
	return &ngrClient
}
