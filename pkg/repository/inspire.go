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

func (hvd *InspireRepository) getKind() ([]interface{}, error) {
	// todo: make this func generic so
	//  - when we ask model.InspireTheme we Download model.InspireKindTheme and parse english and dutch results into []model.InspireTheme and return this
	//  - when we ask model.InspireLayer we Download model.InspireKindLayer and parse english and dutch results into []model.InspireLayer and return this

	// todo: for every file check if it exists in localCachePath
	// todo: If not, download it
	// todo: If file exists, check if it is older than 3 days -> If so, download it
	// todo: If file exists and is not older than 3 days, read it
	// todo: Parse the file and return the result as []T

	return nil, nil
}

func (hvd *InspireRepository) GetAllThemes() ([]interface{}, error) {
	// Call getThemes to retrieve all INSPIRE themes
	return hvd.getKind()
}

func (hvd *InspireRepository) GetAllLayers() ([]interface{}, error) {
	// Call getLayers to retrieve all INSPIRE layers
	return hvd.getKind()
}
