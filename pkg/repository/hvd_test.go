package repository

import (
	"path"
	"testing"

	"github.com/pdok/pdok-metadata-tool/internal/common"
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHvdRepository_GetAllHvdCategories(t *testing.T) {
	hvdRepo := getNewHVDRepository()

	result, err := hvdRepo.GetAllHVDCategories()

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHvdRepository_GetFilteredHvdCategories(t *testing.T) {
	var tests = []struct {
		filterCodes   []string
		expectedCodes []string
	}{
		// Valid specifics
		0: {
			filterCodes:   []string{"c_f76b01e6"},
			expectedCodes: []string{"c_b79e35eb", "c_b151a0ba", "c_f76b01e6"},
		},
		1: {
			filterCodes:   []string{"c_f76b01e6", "c_407951ff"},
			expectedCodes: []string{"c_b79e35eb", "c_b151a0ba", "c_f76b01e6", "c_407951ff"},
		},
		2: {
			filterCodes: []string{"c_f76b01e6", "c_315692ad"},
			expectedCodes: []string{
				"c_dd313021",
				"c_315692ad",
				"c_b79e35eb",
				"c_b151a0ba",
				"c_f76b01e6",
			},
		},
	}
	for _, test := range tests {
		hvdRepo := getNewHVDRepository()
		filteredCategories, err := hvdRepo.GetFilteredHvdCategories(test.filterCodes)
		require.NoError(t, err)

		for i, code := range test.expectedCodes {
			assert.Equal(t, code, filteredCategories[i].ID)
		}
	}
}

func getNewHVDRepository() *HVDRepository {
	hvdCachePath := path.Join(common.GetProjectRoot(), common.HvdLocalRDFPath)

	return NewHVDRepository(hvd.HvdEndpoint, hvdCachePath)
}
