package repository

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	model "github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
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

func (hvd *HVDRepository) Download() ([]byte, error) {
	// Download the thesaurus from thesaurusEndpoint and store it in thesaurusLocalCachePath
	resp, err := http.Get(hvd.thesaurusEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to download thesaurus: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download thesaurus: status code %d", resp.StatusCode)
	}

	// Read all bytes from response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read thesaurus response body: %w", err)
	}

	// Write to the cache file
	if err := os.WriteFile(hvd.thesaurusLocalCachePath, body, 0644); err != nil {
		return nil, fmt.Errorf("failed to write thesaurus to local cache: %w", err)
	}

	return body, nil
}

func (hvd *HVDRepository) getThesaurus() ([]byte, error) {
	// Check if thesaurusLocalCachePath exists and that it is not older than 3 days
	fileInfo, err := os.Stat(hvd.thesaurusLocalCachePath)

	// If file doesn't exist or there's an error, download it
	if os.IsNotExist(err) || err != nil {
		_, err = hvd.Download()
		if err != nil {
			return nil, fmt.Errorf("failed to download thesaurus: %w", err)
		}
	} else {
		// Check if file is older than 3 days
		threedays := time.Hour * 24 * 3
		if time.Since(fileInfo.ModTime()) > threedays {
			_, err = hvd.Download()
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

// hvdCategoriesFromRDF converts an RDF document into a sorted list of HVDCategory
func hvdCategoriesFromRDF(rdf *model.RDF) []model.HVDCategory {
	result := make([]model.HVDCategory, 0, len(rdf.Descriptions))
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

	sort.Slice(result, func(i, j int) bool {
		return result[i].Order < result[j].Order
	})

	return result
}

// GetAllHVDCategoriesFromContent parses the provided RDF XML bytes into HVD categories.
// This allows usage of the parser without needing an HVDRepository instance.
func GetAllHVDCategoriesFromContent(content []byte) ([]model.HVDCategory, error) {
	var rdf model.RDF
	if err := xml.Unmarshal(content, &rdf); err != nil {
		return nil, fmt.Errorf("failed to parse thesaurus: %w", err)
	}
	return hvdCategoriesFromRDF(&rdf), nil
}

func (hvd *HVDRepository) GetAllHVDCategories() (result []model.HVDCategory, err error) {
	rdf, err := hvd.parseThesaurus()
	if err != nil {
		return nil, err
	}

	result = hvdCategoriesFromRDF(rdf)
	return result, nil
}

func (hvd *HVDRepository) GetHVDCategoryByCode(code string) (*model.HVDCategory, error) {
	allCategories, err := hvd.GetAllHVDCategories()
	if err != nil {
		return nil, err
	}

	for _, category := range allCategories {
		if category.Id == code {
			return &category, err
		}
	}

	return nil, fmt.Errorf("no HVD category found for code: %s", code)
}

func (hvd *HVDRepository) GetFilteredHvdCategories(filterCategories []string) ([]model.HVDCategory, error) {
	allCategories, err := hvd.GetAllHVDCategories()
	if err != nil {
		return nil, err
	}

	// Make sure all parent codes are present in the category filter
	filterCategoriesIncludingParents := map[string]bool{}

	for _, filterCategory := range filterCategories {
		category, err := hvd.GetHVDCategoryByCode(filterCategory)
		if err != nil {
			return nil, err
		}

		// Check for 1st level parent
		if category.Parent != "" {
			firstParent, err := hvd.GetHVDCategoryByCode(category.Parent)
			if err != nil {
				return nil, err
			}

			// Check for 2nd level parent
			if firstParent.Parent != "" {
				secondParent, err := hvd.GetHVDCategoryByCode(firstParent.Parent)
				if err != nil {
					return nil, err
				}
				filterCategoriesIncludingParents[secondParent.Id] = true
			}
			filterCategoriesIncludingParents[firstParent.Id] = true
		}
		filterCategoriesIncludingParents[category.Id] = true
	}

	// Filter all HVD categories while keeping the order unchanged
	var result []model.HVDCategory
	for _, category := range allCategories {
		_, ok := filterCategoriesIncludingParents[category.Id]
		if ok {
			result = append(result, category)
		}
	}
	return result, nil
}
