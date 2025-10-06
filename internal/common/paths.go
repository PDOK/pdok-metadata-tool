// Package common provides shared util functions.
package common

import (
	"path/filepath"
	"runtime"
)

var (
	// CachePath is a local path for the cache.
	CachePath = filepath.Join(GetProjectRoot(), "cache")
	// HvdLocalRDFPath is a local path for the HVD RDF file.
	HvdLocalRDFPath = filepath.Join(CachePath, "high-value-dataset-category.rdf")
	// HvdLocalCSVPath is a local path for the HVD CSV file.
	HvdLocalCSVPath = CachePath + "/high-value-dataset-category.csv"
	// InspireLocalPath is a local path for the INSPIRE files.
	InspireLocalPath = CachePath
)

// GetProjectRoot returns the root of the project as a string.
func GetProjectRoot() string {
	// Get the full path of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get caller info")
	}

	// This file is in /internal/common/, so go up two levels to reach the project root
	return filepath.Join(filepath.Dir(filename), "..", "..")
}
