package models

//go:generate go run ../main.go

type Single []int

type Basic struct {
	First  []string
	Second *Single `json:"wee"`
	hidden int
}
