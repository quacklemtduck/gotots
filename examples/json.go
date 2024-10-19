package models

//go:generate go run ../main.go

type JsonTags struct {
	OtherName   string `json:"name_from_json"`
	NotExported string `json:"-"`
	Optional    string `json:",omitempty"`
}
