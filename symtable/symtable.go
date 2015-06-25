package symtable

type (
	SymbolType uint8

	Table map[string]interface{}

	SymTable struct {
		Table
		parent *SymTable
	}
)

const (
	UNKNOWN SymbolType = iota
	VARIABLE
	FUNCTION
)

func New() *SymTable {
	return &SymTable{Table: make(Table, 0)}
}

func (s *SymTable) Set(name string, st SymbolType, symbol interface{}) {
	s.Table[name] = symbol
}

func (s *SymTable) Get(name string, st SymbolType) (interface{}, bool) {
	cur := s
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

func ScopeEnter(s *SymTable) *SymTable {
	t := New()
	t.parent = s
	return t
}

func ScopeLeave(s *SymTable) *SymTable {
	return s.parent
}
