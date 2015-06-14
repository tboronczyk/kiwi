package symtable

import (
	"github.com/tboronczyk/kiwi/token"
)

type (
	symTable struct {
		table  map[string]entry
	}

	entry struct {
		v interface{}
		t token.Token
	}

	SymTable interface {
		Set(string, interface{}, token.Token)
		Get(string) (interface{}, token.Token, bool)
	}
)

func New() *symTable {
	return &symTable{table: make(map[string]entry)}
}

func (s symTable) Set(key string, val interface{}, t token.Token) {
	s.table[key] = entry{v: val, t: t}
	return
}

func (s symTable) Get(key string) (interface{}, token.Token, bool) {
	e, ok := s.table[key]
	if !ok {
		return nil, token.UNKNOWN, ok
	}
	return e.v, e.t, ok
}
