// Package symtable provides an implementation of a symbol table to manage
// symbols during analysis and runtime.
package symtable

type table map[string]interface{}

// SymTable represents a symbol table.
type SymTable struct {
	parent *SymTable // parent of table
	table            // symbol table data
}

// New allocates a new SymTable.
func New() *SymTable {
	s := new(SymTable)
	s.table = make(table, 0)
	return s
}

// Set places symbol sym in the table at key k.
func (s *SymTable) Set(k string, sym interface{}) {
	s.table[k] = sym
}

// Get retrieves the symbol from key k. If the symbol is not found, ok is
// false.
func (s *SymTable) Get(k string) (sym interface{}, ok bool) {
	cur := s
	for {
		if sym, ok = cur.table[k]; ok {
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
