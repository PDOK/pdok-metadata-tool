package client

import (
	"github.com/pdok/pdok-metadata-tool/pkg/model/ngr"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
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
			recordCreated, err := ngrClient.getRecord(tt.args.uuid)
			assert.Nil(t, err)
			assert.Contains(t, recordCreated, "NWB - wegen22222")

			time.Sleep(10 * time.Second)
			recordToBeUpdated := strings.Replace(recordToBeCreated, "NWB - wegen22222", "NWB - wegen33333", -1)
			err = ngrClient.createOrUpdateServiceMetadataRecord(recordToBeUpdated, tt.args.categoryId, tt.args.groupId, true)
			assert.Nil(t, err)
			recordUpdated, err := ngrClient.getRecord(tt.args.uuid)
			assert.Nil(t, err)
			assert.Contains(t, recordUpdated, "NWB - wegen33333")

			time.Sleep(10 * time.Second)
			err = ngrClient.addTagToRecord(tt.args.uuid, INSPIRE_TAG)
			assert.Nil(t, err)

			time.Sleep(10 * time.Second)
			tagsList, err := ngrClient.GetRecordTags(tt.args.uuid)
			assert.Nil(t, err)
			assert.Equal(t, tagsList[0].ID, 224342)

			time.Sleep(10 * time.Second)
			err = ngrClient.deleteRecord(tt.args.uuid)
			assert.Nil(t, err)
			delRecord, err := ngrClient.getRecord(tt.args.uuid)
			assert.NotNil(t, err)
			assert.Equal(t, "", delRecord)
		})
	}
}

func getNgrClient(t *testing.T, mockedNGRServer *httptest.Server) *NgrClient {
	config := NgrConfig{
		//NgrUrl: "https://ngr.acceptatie.nationaalgeoregister.nl",
		//NgrUserName: NGR_USER_NAME,
		//NgrPassword: NGR_PASSWORD,
		NgrUrl:      mockedNGRServer.URL,
		NgrUserName: "NGR_USER_NAME",
		NgrPassword: "NGR_PASSWORD",
	}
	ngrClient := NewNgrClient(config)
	return &ngrClient
}
