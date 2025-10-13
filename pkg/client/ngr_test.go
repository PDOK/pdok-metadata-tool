package client

import (
	"github.com/pdok/pdok-metadata-tool/pkg/model/ngr"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"testing"
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
			// create unpublished record
			err = ngrClient.createOrUpdateServiceMetadataRecord(recordToBeCreated, tt.args.categoryId, tt.args.groupId, false)
			assert.Nil(t, err)
			//recordCreated, err := ngrClient.getRecord(tt.args.uuid)
			//assert.Nil(t, err)
			//assert.Contains(t, recordCreated, "NWB - wegen22222")

			//time.Sleep(60 * time.Second)
			//deleteStatus, err := ngrClient.deleteRecord(tt.args.uuid, &ngrClient.NgrConfig)
			//assert.Nil(t, err)
			//assert.Equal(t, http.StatusNoContent, deleteStatus)
			//delRecord, getDelStatus, err := ngrClient.getRecord(tt.args.uuid, &ngrClient.NgrConfig)
			//assert.NotNil(t, err)
			//assert.Equal(t, getDelStatus, http.StatusNotFound)
			//assert.Equal(t, "", delRecord)

			//err = ngrClient.addTagToRecord(tt.args.uuid, INSPIRE_TAG)
			//assert.Nil(t, err)
			//tagsList, err := ngrClient.GetRecordTags(tt.args.uuid)
			//assert.Nil(t, err)
			//assert.Contains(t, tagsList, "224342")
		})
	}
}

func getNgrClient(t *testing.T, mockedNGRServer *httptest.Server) *NgrClient {
	config := NgrConfig{
		//NgrUrl: "https://ngr.acceptatie.nationaalgeoregister.nl",
		NgrUrl:      mockedNGRServer.URL,
		NgrUserName: NGR_USER_NAME,
		NgrPassword: NGR_PASSWORD,
	}
	ngrClient := NewNgrClient(config)
	return &ngrClient
}

//  gebruik het volgende functie als je wil testen met de ngr.acceptatie omgeving
//func getNgrClient(t *testing.T, mockedNGRServer *httptest.Server) *NgrClient {
//	//hostURL, err := url.Parse(mockedNGRServer.URL)
//	accUrl, err := url.Parse("https://ngr.acceptatie.nationaalgeoregister.nl")
//	if err != nil {
//		t.Fatalf("Failed to parse mocked NGR server URL: %v", err)
//	}
//
//	config := NgrConfig{
//		NgrUrl:      accUrl.String(),
//		NgrUserName: "NGR_USER_NAME",
//		NgrPassword: "NGR_PASSWORD",
//	}
//	ngrClient := NewNgrClient(config)
//	return &ngrClient
//}
