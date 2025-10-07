package client

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pdok/pdok-metadata-tool/pkg/model/csw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCswClient_GetRecords(t *testing.T) {
	mockedNGRServer := preTestSetup()
	cswClient := getCswClient(t, mockedNGRServer)

	type args struct {
		constraint csw.GetRecordsCQLConstraint
		offset     int
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
				constraint: csw.GetRecordsCQLConstraint{
					MetadataType:     &dataset,
					OrganisationName: nil,
				},
				offset: 1,
			},
			wantErr:         false,
			wantNrOfRecords: 10,
			wantNextRecord:  11,
		},
		{
			name: "GetRecords for Services",
			args: args{
				constraint: csw.GetRecordsCQLConstraint{
					MetadataType:     &service,
					OrganisationName: nil,
				},
				offset: 11,
			},
			wantErr:         false,
			wantNrOfRecords: 10,
			wantNextRecord:  21,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mdRecords, nextRecord, err := cswClient.GetRecords(&tt.args.constraint, tt.args.offset)
			if !tt.wantErr {
				require.NoError(t, err)
			}

			assert.Len(t, mdRecords, tt.wantNrOfRecords)
			assert.Equal(t, tt.wantNextRecord, nextRecord)

			for _, record := range mdRecords {
				assert.Equal(t, string(*tt.args.constraint.MetadataType), record.Type)
			}
		})
	}
}

func TestCswClient_GetRecordsWithOGCFilter(t *testing.T) {
	mockedNGRServer := preTestSetup()
	cswClient := getCswClient(t, mockedNGRServer)

	type args struct {
		filter csw.GetRecordsOgcFilter
	}

	tests := []struct {
		name            string
		args            args
		wantErr         bool
		wantNrOfRecords int
	}{
		{
			name: "GetRecordsWithOGCFilter for Datasets",
			args: args{
				filter: csw.GetRecordsOgcFilter{
					MetadataType: csw.Dataset,
					Title:        nil,
					Identifier:   nil,
				},
			},
			wantErr:         false,
			wantNrOfRecords: 10,
		},
		{
			name: "GetRecordsWithOGCFilter for Services",
			args: args{
				filter: csw.GetRecordsOgcFilter{
					MetadataType: csw.Service,
					Title:        nil,
					Identifier:   nil,
				},
			},
			wantErr:         false,
			wantNrOfRecords: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mdRecords, err := cswClient.GetRecordsWithOGCFilter(&tt.args.filter)
			if !tt.wantErr {
				require.NoError(t, err)
			}

			assert.Len(t, mdRecords, tt.wantNrOfRecords)

			for _, record := range mdRecords {
				assert.Equal(t, tt.args.filter.MetadataType.String(), record.Type)
			}
		})
	}
}

func TestCswClient_GetRecordById(t *testing.T) {
	mockedNGRServer := preTestSetup()
	cswClient := getCswClient(t, mockedNGRServer)

	type args struct {
		id string
	}

	tests := []struct {
		name             string
		args             args
		wantErr          bool
		wantMetadataType csw.MetadataType
	}{
		{
			name: "GetRecordByID for Dataset",
			args: args{
				id: "C2DFBDBC-5092-11E0-BA8E-B62DE0D72085",
			},
			wantErr:          false,
			wantMetadataType: csw.Dataset,
		},
		{
			name: "GetRecordByID for Service",
			args: args{
				id: "C2DFBDBC-5092-11E0-BA8E-B62DE0D72086",
			},
			wantErr:          false,
			wantMetadataType: csw.Service,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MDMetadata, err := cswClient.GetRecordByID(tt.args.id)
			if !tt.wantErr {
				require.NoError(t, err)
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

func getCswClient(t *testing.T, mockedNGRServer *httptest.Server) *CswClient {
	t.Helper()

	hostURL, err := url.Parse(mockedNGRServer.URL)
	if err != nil {
		t.Fatalf("Failed to parse mocked NGR server URL: %v", err)
	}

	cswClient := NewCswClient(hostURL)

	return &cswClient
}
