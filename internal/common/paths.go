// Package common provides shared util functions.
package common //nolint:revive,nolintlint

import (
	"path/filepath"
	"runtime"
	"strings"
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
	// MetadataCachePath is a local path for the metadata records cache.
	MetadataCachePath = filepath.Join(CachePath, "records")
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

// NormalizeForFilename converts an arbitrary string into a safe, filesystem-friendly slug.
// Rules:
// - trim spaces
// - lowercase
// - replace spaces with '-'
// - replace any non [a-z0-9-_] rune with '-'
// - collapse multiple '-'
// - trim leading/trailing '-'
// - when input is empty (or normalizes to empty), return "all"
func NormalizeForFilename(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "all"
	}
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")

	// replace any non allowed char with '-'
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		allowed := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_'
		if !allowed {
			// write dash (will be collapsed)
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
			continue
		}
		if r == '-' {
			if !prevDash {
				b.WriteRune(r)
				prevDash = true
			}
			continue
		}
		b.WriteRune(r)
		prevDash = false
	}

	res := strings.Trim(b.String(), "-")
	if res == "" {
		return "all"
	}
	return res
}
