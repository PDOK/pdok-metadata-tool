package common

import (
	"path/filepath"
	"runtime"
)

var (
	CachePath       = filepath.Join(GetProjectRoot(), "cache")
	HvdLocalRDFPath = filepath.Join(CachePath, "high-value-dataset-category.rdf")
	HvdLocalCSVPath = CachePath + "/high-value-dataset-category.csv"

	InspireLocalPath = CachePath
)

func GetProjectRoot() string {
	// Get the full path of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get caller info")
	}

	// This file is in /internal/common/, so go up two levels to reach the project root
	return filepath.Join(filepath.Dir(filename), "..", "..")
}
