package main

// DataType represents the data type of an expression or runtime value.
type DataType uint

const (
	TypUnknown DataType = iota
	TypBuiltin
	TypBool
	TypFunc
	TypNumber
	TypString
)
