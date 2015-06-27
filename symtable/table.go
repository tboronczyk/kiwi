package symtable

type table map[string]interface{}

// SymTable represents a symbol table.
type SymTable struct {
	table []table   // symbol table data
	parent *SymTable // parent of table
}

// New allocates a new SymTable.
func New() *SymTable {
	s := new(SymTable)
	s.table = make([]table, len(stypes))
	for i := 0; i < len(stypes); i++ {
		s.table[i] = make(table, 0)
	}
	return s
}

// Set places symbol sym in type t's symbol table with key k. The table keeps
// types separate; e.g., VAR:foo is a different symbol than FUNC:foo.
func (s *SymTable) Set(k string, t SymType, sym interface{}) {
	s.table[t][k] = sym
}

// Get retrieves the symbol from type t's symbol table identified by key k.
// If the symbol is not found, ok is false.
func (s *SymTable) Get(k string, t SymType) (sym interface{}, ok bool) {
	cur := s
	for {
		if sym, ok = cur.table[t][k]; ok {
			return sym, ok
		}
		if cur.parent == nil {
			return nil, false
		}
		cur = cur.parent
	}
}

// NewScope allocates and returns a new SymTable with s as its parent.
func NewScope(s *SymTable) *SymTable {
	t := New()
	t.parent = s
	return t
}

// Parent returns the parent of SymTable s.
func (s *SymTable) Parent() *SymTable {
	return s.parent
}
