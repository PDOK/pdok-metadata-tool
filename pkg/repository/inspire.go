package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pdok/pdok-metadata-tool/internal/common"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
)

// InspireRepository is used for retrieving INSPIRE related information.
type InspireRepository struct {
	localCachePath string
}

// NewInspireRepository creates a new instance of an InspireRepository.
func NewInspireRepository(localCachePath string) *InspireRepository {
	return &InspireRepository{
		localCachePath: localCachePath,
	}
}

// Download will later download from the INSPIRE endpoints. These are not working correctly, so for now we use local json.
func (ir *InspireRepository) Download(kind inspire.InspireRegisterKind) error {
	dutchURL := inspire.GetInspireEndpoint(kind, inspire.Dutch)
	englishURL := inspire.GetInspireEndpoint(kind, inspire.English)

	dutchFilePath := inspire.GetInspirePath(kind, inspire.Dutch)
	englishFilePath := inspire.GetInspirePath(kind, inspire.English)

	dutchFullPath := filepath.Join(ir.localCachePath, dutchFilePath)
	englishFullPath := filepath.Join(ir.localCachePath, englishFilePath)

	perm := 0755

	if err := os.MkdirAll(ir.localCachePath, os.FileMode(perm)); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := downloadFile(dutchURL, dutchFullPath); err != nil {
		return fmt.Errorf("failed to download Dutch file: %w", err)
	}

	if err := downloadFile(englishURL, englishFullPath); err != nil {
		return fmt.Errorf("failed to download English file: %w", err)
	}

	return nil
}

// GetThemes retrieves the INSPIRE themes.
func (ir *InspireRepository) GetThemes() ([]inspire.InspireTheme, error) {
	englishData, dutchData, err := ir.getKind(inspire.Theme)
	if err != nil {
		return nil, fmt.Errorf("failed to get theme data: %w", err)
	}

	return ir.parseThemes(englishData, dutchData)
}

// GetLayers retrieves the INSPIRE layers.
func (ir *InspireRepository) GetLayers() ([]inspire.InspireLayer, error) {
	englishData, dutchData, err := ir.getKind(inspire.Layer)
	if err != nil {
		return nil, fmt.Errorf("failed to get layer data: %w", err)
	}

	return ir.parseLayers(englishData, dutchData)
}

// downloadFile downloads a file from a URL and saves it to a local path.
func downloadFile(url string, filepath string) error {
	//nolint:gosec
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer common.SafeClose(out)

	//nolint:gosec,bodyclose // We use common.SafeClose to handle closing the response body
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer common.SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

func (ir *InspireRepository) getKind(kind inspire.InspireRegisterKind) ([]byte, []byte, error) {
	dutchFilePath := inspire.GetInspirePath(kind, inspire.Dutch)
	englishFilePath := inspire.GetInspirePath(kind, inspire.English)

	dutchFullPath := filepath.Join(ir.localCachePath, dutchFilePath)
	englishFullPath := filepath.Join(ir.localCachePath, englishFilePath)

	dutchNeedsDownload := true

	if fileInfo, err := os.Stat(dutchFullPath); err == nil {
		// File exists, check if it's older than 3 days
		if time.Since(fileInfo.ModTime()) < 3*24*time.Hour {
			dutchNeedsDownload = false
		}
	}

	englishNeedsDownload := true

	if fileInfo, err := os.Stat(englishFullPath); err == nil {
		// File exists, check if it's older than 3 days
		if time.Since(fileInfo.ModTime()) < 3*24*time.Hour {
			englishNeedsDownload = false
		}
	}

	if dutchNeedsDownload {
		dutchURL := inspire.GetInspireEndpoint(kind, inspire.Dutch)
		if err := downloadFile(dutchURL, dutchFullPath); err != nil {
			return nil, nil, fmt.Errorf("failed to download Dutch file: %w", err)
		}
	}

	if englishNeedsDownload {
		englishURL := inspire.GetInspireEndpoint(kind, inspire.English)
		if err := downloadFile(englishURL, englishFullPath); err != nil {
			return nil, nil, fmt.Errorf("failed to download English file: %w", err)
		}
	}

	// Read Dutch file
	//nolint:gosec
	dutchData, err := os.ReadFile(dutchFullPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read Dutch file: %w", err)
	}

	// Read English file
	//nolint:gosec
	englishData, err := os.ReadFile(englishFullPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read English file: %w", err)
	}

	return englishData, dutchData, nil
}

func (ir *InspireRepository) parseThemes(
	englishData, dutchData []byte,
) ([]inspire.InspireTheme, error) {
	var (
		englishThemeRaw inspire.InspireThemeRaw
		dutchThemeRaw   inspire.InspireThemeRaw
	)

	if err := json.Unmarshal(englishData, &englishThemeRaw); err != nil {
		return nil, fmt.Errorf("failed to parse English themes: %w", err)
	}

	if err := json.Unmarshal(dutchData, &dutchThemeRaw); err != nil {
		return nil, fmt.Errorf("failed to parse Dutch themes: %w", err)
	}

	// Create combined themes

	// OperatesOn
	expectedSize := 34
	themes := make([]inspire.InspireTheme, 0, expectedSize)

	for _, englishItem := range englishThemeRaw.Register.ContainedItems {
		// Extract the short ID from the URL (e.g., "ad" from "http://inspire.ec.europa.eu/theme/ad")
		englishThemeURL := englishItem.Theme.Id

		shortID := englishThemeURL
		if len(englishThemeURL) > 0 {
			// Find the last part of the URL after the last "/"
			lastSlashIndex := -1

			for i := len(englishThemeURL) - 1; i >= 0; i-- {
				if englishThemeURL[i] == '/' {
					lastSlashIndex = i

					break
				}
			}

			if lastSlashIndex != -1 && lastSlashIndex < len(englishThemeURL)-1 {
				shortID = englishThemeURL[lastSlashIndex+1:]
			}
		}

		// Convert themenumber to int
		order := 0

		if englishItem.Theme.ThemeNumber != "" {
			if orderVal, err := strconv.Atoi(englishItem.Theme.ThemeNumber); err == nil {
				order = orderVal
			}
		}

		// Find matching Dutch theme
		var dutchLabel string

		for _, dutchItem := range dutchThemeRaw.Register.ContainedItems {
			if dutchItem.Theme.Id == englishItem.Theme.Id {
				dutchLabel = dutchItem.Theme.Label.Text

				break
			}
		}

		// Create combined theme
		theme := inspire.InspireTheme{
			ID:           shortID,
			Order:        order,
			LabelDutch:   dutchLabel,
			LabelEnglish: englishItem.Theme.Label.Text,
			URL:          englishItem.Theme.Id,
		}
		themes = append(themes, theme)
	}

	return themes, nil
}

func (ir *InspireRepository) parseLayers(
	englishData, dutchData []byte,
) ([]inspire.InspireLayer, error) {
	var (
		englishLayerRaw inspire.InspireLayerRaw
		dutchLayerRaw   inspire.InspireLayerRaw
	)

	if err := json.Unmarshal(englishData, &englishLayerRaw); err != nil {
		return nil, fmt.Errorf("failed to parse English layers: %w", err)
	}

	if err := json.Unmarshal(dutchData, &dutchLayerRaw); err != nil {
		return nil, fmt.Errorf("failed to parse Dutch layers: %w", err)
	}

	// Create combined layers
	expectedSize := 100
	layers := make([]inspire.InspireLayer, 0, expectedSize)

	for _, englishItem := range englishLayerRaw.Register.ContainedItems {
		// Extract the short ID from the URL (e.g., "GE.ActiveWell" from "http://inspire.ec.europa.eu/layer/GE.ActiveWell")
		englishLayerURL := englishItem.Layer.ID

		shortID := englishLayerURL
		if len(englishLayerURL) > 0 {
			// Find the last part of the URL after the last "/"
			lastSlashIndex := -1

			for i := len(englishLayerURL) - 1; i >= 0; i-- {
				if englishLayerURL[i] == '/' {
					lastSlashIndex = i

					break
				}
			}

			if lastSlashIndex != -1 && lastSlashIndex < len(englishLayerURL)-1 {
				shortID = englishLayerURL[lastSlashIndex+1:]
			}
		}

		// Find matching Dutch layer
		var dutchLabel string

		for _, dutchItem := range dutchLayerRaw.Register.ContainedItems {
			if dutchItem.Layer.ID == englishItem.Layer.ID {
				dutchLabel = dutchItem.Layer.Label.Text

				break
			}
		}

		// Create combined layer
		layer := inspire.InspireLayer{
			ID:           shortID,
			LabelDutch:   dutchLabel,
			LabelEnglish: englishItem.Layer.Label.Text,
		}
		layers = append(layers, layer)
	}

	return layers, nil
}
