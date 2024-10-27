package models

//go:generate gotots

type Simple int64

// Composite contains field that refer to other defined types
type Composite struct {
	SimpleField  Simple
	PointerField *Simple
	SimpleArray  []Simple
	SimpleMap    map[int]Simple
	One, Two     int
}
