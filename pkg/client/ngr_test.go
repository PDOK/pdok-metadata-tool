package client

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"pdok-metadata-tool/pkg/model/ngr"
	"testing"
)

const ngrHost = "https://nationaalgeoregister.nl/"

func TestNgrClient_GetRecordTags(t *testing.T) {
	type args struct {
		uuid      string
		logPrefix string
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
				uuid:      "b4ae622c-6201-49d8-bd2e-f7fce9206a1e",
				logPrefix: "UNITTEST_GetRecordTags_INSPIRE_Dataset",
			},
			wantNrOfTags: 1,
			want: ngr.RecordTagsResponse{
				{
					Id:   224342,
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
				uuid:      "c4bda1aa-d6e6-482c-a6f1-bd519e3202d4",
				logPrefix: "GetRecordTags_tagless_Dataset",
			},
			wantNrOfTags: 0,
			want:         ngr.RecordTagsResponse{},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ngrClient, _ := getNgrClient(ngrHost)
			got, err := ngrClient.GetRecordTags(tt.args.uuid, tt.args.logPrefix)

			assert.Nil(t, err)
			assert.Len(t, got, tt.wantNrOfTags)
			assert.Equal(t, tt.want, got)
		})
	}
}

func getNgrClient(host string) (*NgrClient, error) {
	h, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	ngrClient := NewNgrClient(h)
	return &ngrClient, nil
}
