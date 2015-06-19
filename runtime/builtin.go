package runtime

import (
	"fmt"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/symtable"
)

var defs = []ast.FuncDef{
	ast.FuncDef{
		Name: "write",
		Args: []string{"value"},
	},
}

var builtins = map[string]func(symtable.SymTable, symtable.SymTable) (interface{}, symtable.DataType, bool){
	"write": func(varTable symtable.SymTable, funcTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
		value, _, _ := varTable.Get("value")
		fmt.Print(value)
		return true, symtable.BOOL, false
	},
}

func LoadBuiltins(s symtable.SymTable) {
	for _, f := range defs {
		s.Set(f.Name, f, symtable.BUILTIN)
	}
}
