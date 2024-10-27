package models

//go:generate gotots

// JsonTags contains fields decorated with JSON tags
type JsonTags struct {
	OtherName   string `json:"name_from_json"`
	NotExported string `json:"-"`
	Optional    string `json:",omitempty"`
}
