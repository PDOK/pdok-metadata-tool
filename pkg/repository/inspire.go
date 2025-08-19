package repository

import (
	"encoding/json"
	"fmt"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

func (ir *InspireRepository) Download(kind inspire.InspireRegisterKind) error {
	dutchURL := inspire.GetInspireEndpoint(kind, inspire.InspireDutch)
	englishURL := inspire.GetInspireEndpoint(kind, inspire.InspireEnglish)

	dutchFilePath := inspire.GetInspirePath(kind, inspire.InspireDutch)
	englishFilePath := inspire.GetInspirePath(kind, inspire.InspireEnglish)

	dutchFullPath := filepath.Join(ir.localCachePath, dutchFilePath)
	englishFullPath := filepath.Join(ir.localCachePath, englishFilePath)

	if err := os.MkdirAll(ir.localCachePath, 0755); err != nil {
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

// downloadFile downloads a file from a URL and saves it to a local path
func downloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func (ir *InspireRepository) getKind(kind inspire.InspireRegisterKind) ([]byte, []byte, error) {
	dutchFilePath := inspire.GetInspirePath(kind, inspire.InspireDutch)
	englishFilePath := inspire.GetInspirePath(kind, inspire.InspireEnglish)

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
		dutchURL := inspire.GetInspireEndpoint(kind, inspire.InspireDutch)
		if err := downloadFile(dutchURL, dutchFullPath); err != nil {
			return nil, nil, fmt.Errorf("failed to download Dutch file: %w", err)
		}
	}

	if englishNeedsDownload {
		englishURL := inspire.GetInspireEndpoint(kind, inspire.InspireEnglish)
		if err := downloadFile(englishURL, englishFullPath); err != nil {
			return nil, nil, fmt.Errorf("failed to download English file: %w", err)
		}
	}

	// Read Dutch file
	dutchData, err := os.ReadFile(dutchFullPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read Dutch file: %w", err)
	}

	// Read English file
	englishData, err := os.ReadFile(englishFullPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read English file: %w", err)
	}

	return englishData, dutchData, nil
}

func (ir *InspireRepository) parseThemes(englishData, dutchData []byte) ([]inspire.InspireTheme, error) {
	var englishThemeRaw inspire.InspireThemeRaw
	var dutchThemeRaw inspire.InspireThemeRaw

	if err := json.Unmarshal(englishData, &englishThemeRaw); err != nil {
		return nil, fmt.Errorf("failed to parse English themes: %w", err)
	}

	if err := json.Unmarshal(dutchData, &dutchThemeRaw); err != nil {
		return nil, fmt.Errorf("failed to parse Dutch themes: %w", err)
	}

	// Create combined themes
	var themes []inspire.InspireTheme
	for _, englishItem := range englishThemeRaw.Register.Containeditems {
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
		if englishItem.Theme.Themenumber != "" {
			if orderVal, err := strconv.Atoi(englishItem.Theme.Themenumber); err == nil {
				order = orderVal
			}
		}

		// Find matching Dutch theme
		var dutchLabel string
		for _, dutchItem := range dutchThemeRaw.Register.Containeditems {
			if dutchItem.Theme.Id == englishItem.Theme.Id {
				dutchLabel = dutchItem.Theme.Label.Text
				break
			}
		}

		// Create combined theme
		theme := inspire.InspireTheme{
			Id:           shortID,
			Order:        order,
			LabelDutch:   dutchLabel,
			LabelEnglish: englishItem.Theme.Label.Text,
			URL:          englishItem.Theme.Id,
		}
		themes = append(themes, theme)
	}

	return themes, nil
}

func (ir *InspireRepository) parseLayers(englishData, dutchData []byte) ([]inspire.InspireLayer, error) {
	var englishLayerRaw inspire.InspireLayerRaw
	var dutchLayerRaw inspire.InspireLayerRaw

	if err := json.Unmarshal(englishData, &englishLayerRaw); err != nil {
		return nil, fmt.Errorf("failed to parse English layers: %w", err)
	}

	if err := json.Unmarshal(dutchData, &dutchLayerRaw); err != nil {
		return nil, fmt.Errorf("failed to parse Dutch layers: %w", err)
	}

	// Create combined layers
	var layers []inspire.InspireLayer
	for _, englishItem := range englishLayerRaw.Register.Containeditems {
		// Extract the short ID from the URL (e.g., "GE.ActiveWell" from "http://inspire.ec.europa.eu/layer/GE.ActiveWell")
		englishLayerURL := englishItem.Layer.Id
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
		for _, dutchItem := range dutchLayerRaw.Register.Containeditems {
			if dutchItem.Layer.Id == englishItem.Layer.Id {
				dutchLabel = dutchItem.Layer.Label.Text
				break
			}
		}

		// Create combined layer
		layer := inspire.InspireLayer{
			Id:           shortID,
			LabelDutch:   dutchLabel,
			LabelEnglish: englishItem.Layer.Label.Text,
		}
		layers = append(layers, layer)
	}

	return layers, nil
}

func (ir *InspireRepository) GetThemes() ([]inspire.InspireTheme, error) {
	englishData, dutchData, err := ir.getKind(inspire.InspireKindTheme)
	if err != nil {
		return nil, fmt.Errorf("failed to get theme data: %w", err)
	}

	return ir.parseThemes(englishData, dutchData)
}

func (ir *InspireRepository) GetLayers() ([]inspire.InspireLayer, error) {
	englishData, dutchData, err := ir.getKind(inspire.InspireKindLayer)
	if err != nil {
		return nil, fmt.Errorf("failed to get layer data: %w", err)
	}

	return ir.parseLayers(englishData, dutchData)
}
