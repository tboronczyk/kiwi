package main

import "strconv"

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

var types = []string{
	TypUnknown: "Unknown",
	TypBuiltin: "Builtin",
	TypBool:    "Bool",
	TypFunc:    "Func",
	TypNumber:  "Number",
	TypString:  "String",
}

func (dt DataType) String() string {
	if dt < DataType(len(types)) {
		return types[dt]
	}
	return "DataType(" + strconv.Itoa(int(dt)) + ")"
}
