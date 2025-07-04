package repository

import "pdok-metadata-tool/pkg/model"

type InspireRepository struct {
	localCachePath string
}

func NewInspireRepository(localCachePath string) *InspireRepository {
	return &InspireRepository{
		localCachePath: localCachePath,
	}
}

func (hvd *InspireRepository) Download(kind model.InspireRegisterKind) error {

	// todo: get dutch url from model.GetInspireEndpoint(kind, model.InspireDutch)
	// todo: get english url from model.GetInspireEndpoint(kind, model.InspireEnglish)
	// todo: download both english and dutch file
	// todo: get dutch file path from model.GetInspirePath(kind, model.InspireDutch)
	// todo: get english file path from model.GetInspirePath(kind, model.InspireEnglish)
	// todo: store both files in localCachePath

	return nil
}

func (hvd *InspireRepository) getThemes() ([]model.InspireTheme, error) {
	// todo: Check if file exists in localCachePath
	// todo: If not, download it
	// todo: If file exists, check if it is older than 3 days -> If so, download it
	// todo: If file exists and is not older than 3 days, read it
	// todo: Parse the file and return the result as []model.InspireTheme

	return nil, nil
}

func (hvd *InspireRepository) getLayers() ([]model.InspireLayer, error) {
	// todo: Check if file exists in localCachePath
	// todo: If not, download it
	// todo: If file exists, check if it is older than 3 days -> If so, download it
	// todo: If file exists and is not older than 3 days, read it
	// todo: Parse the file and return the result as []model.InspireLayer

	return nil, nil
}

func (hvd *InspireRepository) GetAllThemes() ([]model.InspireTheme, error) {
	// Call getThemes to retrieve all INSPIRE themes
	return hvd.getThemes()
}

func (hvd *InspireRepository) GetAllLayers() ([]model.InspireLayer, error) {
	// Call getLayers to retrieve all INSPIRE layers
	return hvd.getLayers()
}
