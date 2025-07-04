package repository

type HVDRepository struct {
	thesaurusEndpoint       string
	thesaurusLocalCachePath string
}

func NewHVDRepository(thesaurusEndpoint string, thesaurusLocalCachePath string) *HVDRepository {
	return &HVDRepository{
		thesaurusEndpoint:       thesaurusEndpoint,
		thesaurusLocalCachePath: thesaurusLocalCachePath,
	}
}

func (hvd *HVDRepository) Download() error {
	return nil
}

func (hvd *HVDRepository) GetAllHVDCategories() error {
	return nil
}
