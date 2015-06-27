// Package symtable provides the symbol table implementation for managing
// symbols during analysis.
package symtable

// Type represents the type of value placed in the table. The symbol table
// keeps types separate; that is, VAR:foo is treated as a different symbol
// than FUNC:foo.
type Type uint8
const (
	UNKNOWN Type = iota
	VAR
	VARTYPE
	FUNC
	FUNCTYPE
)
const numTypes = 5

type table map[string]interface{}

// SymTable represents a symbol table.
type SymTable struct {
	T []table    // symbol table data
	P *SymTable  // parent table
}

// New allocates a new SymTable
func New() *SymTable {
	s := new(SymTable)
	s.T = make([]table, numTypes)
	for i := 0; i < numTypes; i++ {
		s.T[i] = make(table, 0)
	}
	return s
}

// Set places symbol sym in type t's symbol table with key k.
func (s *SymTable) Set(k string, t Type, sym interface{}) {
	s.T[t][k] = sym
}

// Get retrieves the symbol from type t's symbol table identified by key k.
// If no symbol is stored, ok is false.
func (s *SymTable) Get(k string, t Type) (sym interface{}, ok bool) {
	cur := s
	for {
		if sym, ok = cur.T[t][k]; ok {
			return sym, ok
		}
		if cur.P == nil {
			return nil, false
		}
		cur = cur.P
	}
}

// NewScope allocates and returns a new SymTable with s as its parent.
func NewScope(s *SymTable) *SymTable {
	t := New()
	t.P = s
	return t
}

// Parent returns the parent of SymTable s.
func (s *SymTable) Parent() *SymTable {
	return s.P
}
