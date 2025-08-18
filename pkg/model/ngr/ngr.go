package ngr

type RecordTagsResponse []Tag

type Tag struct {
	Id    int               `json:"id"`
	Name  string            `json:"name"`
	Label map[string]string `json:"label"`
}
