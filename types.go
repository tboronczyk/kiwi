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

// String returns the string representation of a type.
func (dt DataType) String() string {
	str := ""
	if dt >= 0 && dt < DataType(len(types)) {
		str = types[dt]
	}
	if str == "" {
		str = "DataType(" + strconv.Itoa(int(dt)) + ")"
	}
	return str
}
