package models

//go:generate go run ../main.go

type Single []int

type Basic struct {
	First  []int
	Second *Single `json:"name_from_json"`
	hidden int
}

type BasicTypes struct {
	BoolType bool

	StringType string

	// Number types
	IntType   int
	Int8Type  int8
	Int16Type int16
	Int32Type int32
	Int64Type int64

	UintType   uint
	Uint8Type  uint8
	Uint16Type uint16
	Uint32Type uint32
	Uint64Type uint64

	Float32Type float32
	Float64Type float64

	ByteType byte

	RuneType rune

	// Complex number are not supported (complex64 and complex128)
}
