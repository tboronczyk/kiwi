package symtable

type (
	SymTable map[string]entry

	entry struct {
		value interface{}
		dtype DataType
	}
)

func New() SymTable {
	return make(SymTable, 0)
}

func (s *SymTable) Set(k string, v interface{}, t DataType) {
	(*s)[k] = entry{value: v, dtype: t}
}

func (s SymTable) Get(k string) (interface{}, DataType, bool) {
	e, ok := s[k]
	if !ok {
		return nil, UNKNOWN, ok
	}
	return e.value, e.dtype, ok
}
