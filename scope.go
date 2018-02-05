package main

type (
	Scope struct {
		parent *Scope
		vars   ScopeTable
		funcs  ScopeTable
	}

	ScopeTable map[string]ScopeEntry

	ScopeEntry struct {
		DataType
		Value interface{}
	}
)

func NewScope() *Scope {
	s := &Scope{
		vars:  make(ScopeTable, 0),
		funcs: make(ScopeTable, 0),
	}
	return s
}

func NewScopeWithParent(p *Scope) *Scope {
	s := NewScope()
	s.parent = p
	return s
}

// EmptyVarCopy returns a copy of the scope with a empty var table but funcs
// still defined.
func (s *Scope) EmptyVarCopy() *Scope {
	s2 := NewScope()
	s2.parent = s.parent
	s2.funcs = s.funcs
	return s2
}

func (s *Scope) SetVar(key string, entry ScopeEntry) {
	s.vars[key] = entry
}

func (s *Scope) GetVar(key string) (ScopeEntry, bool) {
	entry, ok := s.vars[key]
	return entry, ok
}

func (s *Scope) SetFunc(key string, entry ScopeEntry) {
	s.funcs[key] = entry
}

func (s *Scope) GetFunc(key string) (ScopeEntry, bool) {
	cur := s
	for {
		if entry, ok := cur.funcs[key]; ok || cur.parent == nil {
			return entry, ok
		}
		cur = cur.parent
	}
}
