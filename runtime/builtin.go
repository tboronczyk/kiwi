package runtime

import (
	"fmt"
	"github.com/tboronczyk/kiwi/ast"
)

var funcSigs = []ast.FuncDefNode{
	ast.FuncDefNode{
		Name: "write",
		Args: []string{"value"},
	},
}

var builtinFuncs = map[string]func(*Runtime){
	"write": func(r *Runtime) {
		value, _, _ := r.varGet("value")
		fmt.Print(value)
	},
}
