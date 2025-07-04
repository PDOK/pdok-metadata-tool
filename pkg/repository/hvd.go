package repository

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"pdok-metadata-tool/pkg/model"
	"strings"
	"time"
)

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
	// Download the thesaurus from thesaurusEndpoint and store it in thesaurusLocalCachePath
	resp, err := http.Get(hvd.thesaurusEndpoint)
	if err != nil {
		return fmt.Errorf("failed to download thesaurus: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download thesaurus: status code %d", resp.StatusCode)
	}

	// Create the file
	file, err := os.Create(hvd.thesaurusLocalCachePath)
	if err != nil {
		return fmt.Errorf("failed to create local cache file: %w", err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write thesaurus to local cache: %w", err)
	}

	return nil
}

func (hvd *HVDRepository) getThesaurus() ([]byte, error) {
	// Check if thesaurusLocalCachePath exists and that it is not older than 3 days
	fileInfo, err := os.Stat(hvd.thesaurusLocalCachePath)

	// If file doesn't exist or there's an error, download it
	if os.IsNotExist(err) || err != nil {
		err = hvd.Download()
		if err != nil {
			return nil, fmt.Errorf("failed to download thesaurus: %w", err)
		}
	} else {
		// Check if file is older than 3 days
		threedays := time.Hour * 24 * 3
		if time.Since(fileInfo.ModTime()) > threedays {
			err = hvd.Download()
			if err != nil {
				return nil, fmt.Errorf("failed to download thesaurus: %w", err)
			}
		}
	}

	content, err := os.ReadFile(hvd.thesaurusLocalCachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read thesaurus file: %w", err)
	}

	return content, nil
}

func (hvd *HVDRepository) parseThesaurus() (*model.RDF, error) {
	content, err := hvd.getThesaurus()
	if err != nil {
		return nil, fmt.Errorf("failed to get thesaurus: %w", err)
	}

	var rdf model.RDF
	err = xml.Unmarshal(content, &rdf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse thesaurus: %w", err)
	}

	return &rdf, nil
}

func (hvd *HVDRepository) GetAllHVDCategories() (result []model.HVDCategory, err error) {
	rdf, err := hvd.parseThesaurus()
	if err != nil {
		return nil, err
	}

	result = make([]model.HVDCategory, 0, len(rdf.Descriptions))
	for _, desc := range rdf.Descriptions {
		// Skip if not a Concept
		if desc.Type.Resource != "http://www.w3.org/2004/02/skos/core#Concept" {
			continue
		}

		// Extract ID from the "about" attribute
		id := desc.Identifier
		if id == "" {
			// If identifier is empty, extract from about URL
			parts := strings.Split(desc.About, "/")
			if len(parts) > 0 {
				id = parts[len(parts)-1]
			}
		}

		// Find Dutch and English labels
		var labelDutch, labelEnglish string
		for _, label := range desc.PrefLabels {
			if label.Lang == "nl" {
				labelDutch = label.Value
			} else if label.Lang == "en" {
				labelEnglish = label.Value
			}
		}

		// Extract parent from broader
		parent := ""
		if desc.Broader.Resource != "" {
			parts := strings.Split(desc.Broader.Resource, "/")
			if len(parts) > 0 {
				parent = parts[len(parts)-1]
			}
		}

		category := model.HVDCategory{
			Id:           id,
			Parent:       parent,
			Order:        desc.Order.Value,
			LabelDutch:   labelDutch,
			LabelEnglish: labelEnglish,
		}

		result = append(result, category)
	}

	return result, nil
}
