package symtable

import (
	"strconv"
	"github.com/tboronczyk/kiwi/util"
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
	entry struct {
		v interface{}
		t DataType
	}

	table map[string]entry

	symTable struct {
		stack *util.Stack
	}

	SymTable interface {
		Set(string, interface{}, DataType)
		Get(string) (interface{}, DataType, bool)
	}
)

func New() *symTable {
	s := &symTable{stack: util.NewStack()}
	s.stack.Push(make(table))
	return s
}

func (s symTable) Set(key string, val interface{}, t DataType) {
	s.stack.Peek().(table)[key] = entry{v: val, t: t}
	return
}

func (s symTable) Get(key string) (interface{}, DataType, bool) {
	e, ok := s.stack.Peek().(table)[key]
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
