package models

//go:generate go run ../main.go

type Simple int64

type Composite struct {
	SimpleField  Simple
	PointerField *Simple
	SimpleArray  []Simple
	SimpleMap    map[int]Simple
	One, Two     int
}
