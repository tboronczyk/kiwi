package symtable

import (
	"strconv"
)

// SymType represents the type of value placed in the symbol table.
type SymType uint8

const (
	UNKNOWN SymType = iota
	VAR
	FUNC
)

var stypes = []string{
	UNKNOWN: "UNKNOWN",
	VAR:     "VAR",
	FUNC:    "FUNC",
}

// String returns the string representation of a type.
func (t SymType) String() string {
	str := ""
	if t >= 0 && t < SymType(len(stypes)) {
		str = stypes[t]
	}
	if str == "" {
		str = "SymType(" + strconv.Itoa(int(t)) + ")"
	}
	return str
}
