package symtable

type (
	SymbolType uint8

	Table map[string]interface{}

	Scope struct {
		Table
		parent *Scope
	}

	SymTable struct {
		scope *Scope
	}
)

const (
	UNKNOWN SymbolType = iota
	VARIABLE
	FUNCTION
)

func New() *SymTable {
	s := &SymTable{
		scope: &Scope{Table: make(Table, 0)},
	}
	return s
}

func (s *SymTable) Set(name string, st SymbolType, symbol interface{}) {
	s.scope.Table[name] = symbol
}

func (s *SymTable) Get(name string, st SymbolType) (interface{}, bool) {
	cur := s.scope
	for {
		if sym, ok := cur.Table[name]; ok {
			return sym, true
		}
		if cur.parent == nil {
			return nil, false
		}
		cur = cur.parent
	}
}

func (s *SymTable) Enter() {
	s.scope = &Scope{
		Table:  make(Table, 0),
		parent: s.scope,
	}
}

func (s *SymTable) Leave() {
	s.scope = s.scope.parent
}

func (s SymTable) Current() *Scope {
	return s.scope
}
