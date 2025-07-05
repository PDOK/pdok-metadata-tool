package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"pdok-metadata-tool/pkg/model"
	"time"
)

type InspireRepository struct {
	localCachePath string
}

func NewInspireRepository(localCachePath string) *InspireRepository {
	return &InspireRepository{
		localCachePath: localCachePath,
	}
}

func (ir *InspireRepository) Download(kind model.InspireRegisterKind) error {
	// Get Dutch and English URLs
	dutchURL := model.GetInspireEndpoint(kind, model.InspireDutch)
	englishURL := model.GetInspireEndpoint(kind, model.InspireEnglish)

	// Get Dutch and English file paths
	dutchFilePath := model.GetInspirePath(kind, model.InspireDutch)
	englishFilePath := model.GetInspirePath(kind, model.InspireEnglish)

	// Create full paths for storing files
	dutchFullPath := filepath.Join(ir.localCachePath, dutchFilePath)
	englishFullPath := filepath.Join(ir.localCachePath, englishFilePath)

	// Ensure directory exists
	if err := os.MkdirAll(ir.localCachePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Download Dutch file
	if err := downloadFile(dutchURL, dutchFullPath); err != nil {
		return fmt.Errorf("failed to download Dutch file: %w", err)
	}

	// Download English file
	if err := downloadFile(englishURL, englishFullPath); err != nil {
		return fmt.Errorf("failed to download English file: %w", err)
	}

	return nil
}

// downloadFile downloads a file from a URL and saves it to a local path
func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func (ir *InspireRepository) getKind(kind model.InspireRegisterKind) ([]byte, []byte, error) {
	// Get Dutch and English file paths
	dutchFilePath := model.GetInspirePath(kind, model.InspireDutch)
	englishFilePath := model.GetInspirePath(kind, model.InspireEnglish)

	// Create full paths for storing files
	dutchFullPath := filepath.Join(ir.localCachePath, dutchFilePath)
	englishFullPath := filepath.Join(ir.localCachePath, englishFilePath)

	// Ensure directory exists
	if err := os.MkdirAll(ir.localCachePath, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Check Dutch file
	dutchNeedsDownload := true
	if fileInfo, err := os.Stat(dutchFullPath); err == nil {
		// File exists, check if it's older than 3 days
		if time.Since(fileInfo.ModTime()) < 3*24*time.Hour {
			dutchNeedsDownload = false
		}
	}

	// Check English file
	englishNeedsDownload := true
	if fileInfo, err := os.Stat(englishFullPath); err == nil {
		// File exists, check if it's older than 3 days
		if time.Since(fileInfo.ModTime()) < 3*24*time.Hour {
			englishNeedsDownload = false
		}
	}

	// Download Dutch file if needed
	if dutchNeedsDownload {
		dutchURL := model.GetInspireEndpoint(kind, model.InspireDutch)
		if err := downloadFile(dutchURL, dutchFullPath); err != nil {
			return nil, nil, fmt.Errorf("failed to download Dutch file: %w", err)
		}
	}

	// Download English file if needed
	if englishNeedsDownload {
		englishURL := model.GetInspireEndpoint(kind, model.InspireEnglish)
		if err := downloadFile(englishURL, englishFullPath); err != nil {
			return nil, nil, fmt.Errorf("failed to download English file: %w", err)
		}
	}

	// Read Dutch file
	dutchData, err := ioutil.ReadFile(dutchFullPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read Dutch file: %w", err)
	}

	// Read English file
	englishData, err := ioutil.ReadFile(englishFullPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read English file: %w", err)
	}

	return englishData, dutchData, nil
}

func (ir *InspireRepository) parseKind(kind model.InspireRegisterKind) ([]interface{}, error) {
	// todo: use generic types

	// Get English and Dutch data
	englishData, dutchData, err := ir.getKind(kind)
	if err != nil {
		return nil, fmt.Errorf("failed to get kind data: %w", err)
	}

	// Parse based on kind
	switch kind {
	case model.InspireKindTheme:
		return ir.parseThemes(englishData, dutchData)
	case model.InspireKindLayer:
		return ir.parseLayers(englishData, dutchData)
	default:
		return nil, fmt.Errorf("unknown kind: %s", kind)
	}
}

func (ir *InspireRepository) parseThemes(englishData, dutchData []byte) ([]interface{}, error) {

	// todo: add raw Struct to model
	var englishThemes []struct {
		Id    string `json:"id"`
		Order int    `json:"order"`
		Label string `json:"label"`
		URL   string `json:"url"`
	}
	var dutchThemes []struct {
		Id    string `json:"id"`
		Order int    `json:"order"`
		Label string `json:"label"`
		URL   string `json:"url"`
	}

	// Parse English data
	if err := json.Unmarshal(englishData, &englishThemes); err != nil {
		return nil, fmt.Errorf("failed to parse English themes: %w", err)
	}

	// Parse Dutch data
	if err := json.Unmarshal(dutchData, &dutchThemes); err != nil {
		return nil, fmt.Errorf("failed to parse Dutch themes: %w", err)
	}

	// Create combined themes
	var themes []interface{}
	for _, englishTheme := range englishThemes {
		// Find matching Dutch theme
		var dutchLabel string
		for _, dutchTheme := range dutchThemes {
			if dutchTheme.Id == englishTheme.Id {
				dutchLabel = dutchTheme.Label
				break
			}
		}

		// Create combined theme
		theme := model.InspireTheme{
			Id:           englishTheme.Id,
			Order:        englishTheme.Order,
			LabelDutch:   dutchLabel,
			LabelEnglish: englishTheme.Label,
			URL:          englishTheme.URL,
		}
		themes = append(themes, theme)
	}

	return themes, nil
}

func (ir *InspireRepository) parseLayers(englishData, dutchData []byte) ([]interface{}, error) {
	var englishLayers []struct {
		Id    string `json:"id"`
		Label string `json:"label"`
	}
	var dutchLayers []struct {
		Id    string `json:"id"`
		Label string `json:"label"`
	}

	// Parse English data
	if err := json.Unmarshal(englishData, &englishLayers); err != nil {
		return nil, fmt.Errorf("failed to parse English layers: %w", err)
	}

	// Parse Dutch data
	if err := json.Unmarshal(dutchData, &dutchLayers); err != nil {
		return nil, fmt.Errorf("failed to parse Dutch layers: %w", err)
	}

	// Create combined layers
	var layers []interface{}
	for _, englishLayer := range englishLayers {
		// Find matching Dutch layer
		var dutchLabel string
		for _, dutchLayer := range dutchLayers {
			if dutchLayer.Id == englishLayer.Id {
				dutchLabel = dutchLayer.Label
				break
			}
		}

		// Create combined layer
		layer := model.InspireLayer{
			Id:           englishLayer.Id,
			LabelDutch:   dutchLabel,
			LabelEnglish: englishLayer.Label,
		}
		layers = append(layers, layer)
	}

	return layers, nil
}

func (ir *InspireRepository) GetAllThemes() ([]interface{}, error) {
	// todo: use generic types
	// Call parseKind to retrieve all INSPIRE themes
	return ir.parseKind(model.InspireKindTheme)
}

func (ir *InspireRepository) GetAllLayers() ([]interface{}, error) {
	// todo: use generic types
	// Call parseKind to retrieve all INSPIRE layers
	return ir.parseKind(model.InspireKindLayer)
}
