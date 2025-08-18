package client

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"pdok-metadata-tool/pkg/model/csw"
	"testing"
)

const cswHost = "https://nationaalgeoregister.nl/"
const cswPath = "/geonetwork/srv/dut/csw"

func TestCswClient_GetRecords(t *testing.T) {

	type args struct {
		constraint csw.GetRecordsConstraint
		offset     int
		logPrefix  string
	}

	dataset := csw.Dataset
	service := csw.Service

	tests := []struct {
		name            string
		args            args
		wantErr         bool
		wantNrOfRecords int
		wantNextRecord  int
	}{
		{
			name: "GetRecords for Datasets",
			args: args{
				constraint: csw.GetRecordsConstraint{
					MetadataType:     &dataset,
					OrganisationName: nil,
				},
				offset:    1,
				logPrefix: "UNITTEST_GetRecords_Datasets",
			},
			wantErr:         false,
			wantNrOfRecords: 10,
			wantNextRecord:  11,
		},
		{
			name: "GetRecords for Services",
			args: args{
				constraint: csw.GetRecordsConstraint{
					MetadataType:     &service,
					OrganisationName: nil,
				},
				offset:    11,
				logPrefix: "UNITTEST_GetRecords_Services",
			},
			wantErr:         false,
			wantNrOfRecords: 10,
			wantNextRecord:  21,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cswClient, _ := getCswClient(cswHost, cswPath)

			mdRecords, nextRecord, err := cswClient.GetRecords(&tt.args.constraint, tt.args.offset, tt.args.logPrefix)
			if !tt.wantErr {
				assert.Nil(t, err)
			}
			assert.Len(t, mdRecords, tt.wantNrOfRecords)
			assert.Equal(t, nextRecord, tt.wantNextRecord)
			for _, record := range mdRecords {
				assert.Equal(t, string(*tt.args.constraint.MetadataType), record.Type)
			}
		})
	}

}

func TestCswClient_GetRecordById(t *testing.T) {
	type args struct {
		id        string
		logPrefix string
	}
	tests := []struct {
		name             string
		args             args
		wantErr          bool
		wantMetadataType csw.MetadataType
	}{
		{
			name: "GetRecordById for Dataset",
			args: args{
				id:        "2350b86b-3efd-47e4-883e-519bfa8d0abd",
				logPrefix: "UNITTEST_GetRecordById_Dataset",
			},
			wantErr:          false,
			wantMetadataType: csw.Dataset,
		},
		{
			name: "GetRecordById for Service",
			args: args{
				id:        "b7b9859e-c1ca-465d-83c2-f24c2d2567b4",
				logPrefix: "UNITTEST_GetRecordById_Service",
			},
			wantErr:          false,
			wantMetadataType: csw.Service,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cswClient, _ := getCswClient(cswHost, cswPath)

			MDMetadata, err := cswClient.GetRecordById(tt.args.id, tt.args.logPrefix)
			if !tt.wantErr {
				assert.Nil(t, err)
			}
			assert.Equal(t, MDMetadata.UUID, tt.args.id)

			if tt.wantMetadataType == csw.Dataset {
				assert.NotNil(t, MDMetadata.IdentificationInfo.MDDataIdentification)
				assert.NotEmpty(t, MDMetadata.IdentificationInfo.MDDataIdentification.Title)
			}

			if tt.wantMetadataType == csw.Service {
				assert.NotNil(t, MDMetadata.IdentificationInfo.SVServiceIdentification)
				assert.NotEmpty(t, MDMetadata.IdentificationInfo.SVServiceIdentification.Title)
			}
		})
	}
}

func getCswClient(host string, path string) (*CswClient, error) {
	h, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	h.Path = path

	cswClient := NewCswClient(h)
	return &cswClient, nil
}
