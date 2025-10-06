// Package ngr provides structs for handling NGR responses.
package ngr

// RecordTagsResponse for retrieving tags from NGR.
type RecordTagsResponse []Tag

// Tag struct for retrieving tags from NGR.
type Tag struct {
	ID    int               `json:"id"`
	Name  string            `json:"name"`
	Label map[string]string `json:"label"`
}
