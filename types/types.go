package types

import (
	"strconv"
)

// DataType represents the data type of an expression or runtime value.
type DataType uint8

const (
	UNKNOWN DataType = iota
	BUILTIN
	BOOL
	FUNC
	NUMBER
	STRING
)

var dtypes = []string{
	UNKNOWN: "UNKNOWN",
	BUILTIN: "BUILTIN",
	BOOL:    "BOOL",
	FUNC:    "FUNC",
	NUMBER:  "NUMBER",
	STRING:  "STRING",
}

// String returns the string representation of a type.
func (t DataType) String() string {
	str := ""
	if t >= 0 && t < DataType(len(dtypes)) {
		str = dtypes[t]
	}
	if str == "" {
		str = "DataType(" + strconv.Itoa(int(t)) + ")"
	}
	return str
}
