package models

//go:generate go run ../main.go

type Single []int

type Basic struct {
	First  []int
	Second *Single `json:"name_from_json"`
	hidden int
	What   map[string]int
}

type notExported struct {
	Wee int64
}
