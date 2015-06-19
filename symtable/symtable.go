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

var dtypes = [...]string{
	UNKNOWN: "UNKNOWN",
	BUILTIN: "BUILTIN",
	USRFUNC: "USRFUNC",
	BOOL:    "BOOL",
	NUMBER:  "NUMBER",
	STRING:  "STRING",
}

type (
	symTable struct {
		table map[string]entry
	}

	entry struct {
		v interface{}
		t DataType
	}

	SymTable interface {
		Set(string, interface{}, DataType)
		Get(string) (interface{}, DataType, bool)
	}
)

func New() *symTable {
	return &symTable{table: make(map[string]entry)}
}

func (s symTable) Set(key string, val interface{}, t DataType) {
	s.table[key] = entry{v: val, t: t}
	return
}

func (s symTable) Get(key string) (interface{}, DataType, bool) {
	e, ok := s.table[key]
	if !ok {
		return nil, UNKNOWN, ok
	}
	return e.v, e.t, ok
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
