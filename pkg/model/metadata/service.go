package metadata

import (
	"github.com/pdok/pdok-metadata-tool/pkg/model/hvd"
	"github.com/pdok/pdok-metadata-tool/pkg/model/inspire"
)

// NLDatasetMetadata is used for retrieving the relevant fields from dataset metadata.
type NLServiceMetadata struct {
	MetadataID     string
	SourceID       string
	Title          string
	Abstract       string
	ContactName    string
	ContactEmail   string
	ContactURL     string
	Keywords       []string
	LicenceURL     string
	UseLimitation  string
	ThumbnailURL   *string
	InspireVariant *inspire.InspireVariant
	InspireThemes  []string
	HVDCategories  []hvd.HVDCategory
	BoundingBox    *BoundingBox
}
