package main

type (
	Entry struct {
		Value interface{}
		DataType
	}

	table map[string]Entry

	Scope struct {
		Parent *Scope
		vars   table
		funcs  table
	}
)

func NewScope() *Scope {
	s := &Scope{
		vars:  make(table, 0),
		funcs: make(table, 0),
	}
	return s
}

func NewScopeWithParent(p *Scope) *Scope {
	s := NewScope()
	s.Parent = p
	return s
}

func NewScopeClone(s *Scope) *Scope {
	s2 := NewScope()
	s2.funcs = s.funcs
	s2.Parent = s.Parent
	return s2
}

func (s *Scope) SetVar(k string, e Entry) {
	s.vars[k] = e
}

func (s *Scope) GetVar(k string) (e Entry, ok bool) {
	e, ok = s.vars[k]
	return
}

func (s *Scope) SetFunc(k string, e Entry) {
	s.funcs[k] = e
}

func (s *Scope) GetFunc(k string) (e Entry, ok bool) {
	cur := s
	for {
		if e, ok = cur.funcs[k]; ok || cur.Parent == nil {
			return
		}
		cur = cur.Parent
	}
}
