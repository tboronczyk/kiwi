package symtable

import (
	"strconv"
)

type DataType uint8

const (
	UNKNOWN DataType = iota
	BUILTIN
	USRFUNC
	BOOL
	NUMBER
	STRING
)

var dtypes = []string{
	UNKNOWN: "UNKNOWN",
	BUILTIN: "BUILTIN",
	USRFUNC: "USRFUNC",
	BOOL:    "BOOL",
	NUMBER:  "NUMBER",
	STRING:  "STRING",
}

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
