package ast

import (
	"github.com/tboronczyk/kiwi/symtable"
	"github.com/tboronczyk/kiwi/token"
	"strconv"
)

func (n AssignStmt) Eval(varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	value, dtype, _ := n.Expr.Eval(varTable, funTable)
	varTable.Set(n.Name, value, dtype)
	return value, dtype, false
}

func (n FuncDef) Eval(varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	funTable.Set(n.Name, n, symtable.USRFUNC)
	return true, symtable.BOOL, false
}

func (n FuncCall) Eval(varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	f, t, ok := funTable.Get(n.Name)
	if !ok {
		panic("Function not defined")
	}
	if len(f.(FuncDef).Args) != len(n.Args) {
		panic("Function arity mis-match")
	}
	s := symtable.New()
	for i, name := range f.(FuncDef).Args {
		value, dtype, _ := n.Args[i].Eval(varTable, funTable)
		s.Set(name, value, dtype)
	}
	switch t {
	case symtable.USRFUNC:
		for _, stmt := range f.(FuncDef).Body {
			value, dtype, isReturn := stmt.Eval(s, funTable)
			if isReturn {
				return value, dtype, isReturn
			}
		}
	case symtable.BUILTIN:
		return builtins[n.Name](s, funTable)
	}

	return true, symtable.BOOL, false
}

func (n IfStmt) Eval(varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	value, dtype, _ := n.Condition.Eval(varTable, funTable)
	if dtype != symtable.BOOL {
		panic("Invalid data type")
	}
	if value.(bool) {
		var isReturn bool
		for _, stmt := range n.Body {
			value, dtype, isReturn = stmt.Eval(varTable, funTable)
			if isReturn {
				return value, dtype, true
			}
		}
	}
	return true, symtable.BOOL, false
}

func (n ReturnStmt) Eval(varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	value, dtype, _ := n.Expr.Eval(varTable, funTable)
	return value, dtype, true
}

func (n WhileStmt) Eval(varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	for {
		value, dtype, _ := n.Condition.Eval(varTable, funTable)
		if dtype != symtable.BOOL {
			panic("Invalid data type")
		}
		if !value.(bool) {
			return true, symtable.BOOL, false
		}
		var isReturn bool
		for _, stmt := range n.Body {
			value, dtype, isReturn = stmt.Eval(varTable, funTable)
			if isReturn {
				return value, dtype, true
			}
		}
	}
}

func (n ValueExpr) Eval(varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	switch n.Type {
	case token.NUMBER:
		value, _ := strconv.ParseFloat(n.Value, 64)
		return value, symtable.NUMBER, false
	case token.BOOL:
		return n.Value == "TRUE", symtable.BOOL, false
	}
	return n.Value, symtable.STRING, false
}

func (n VariableExpr) Eval(varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	value, dtype, ok := varTable.Get(n.Name)
	if !ok {
		panic("Variable not set")
	}
	return value, dtype, false
}

func (n UnaryExpr) Eval(varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	rVal, rType, _ := n.Right.Eval(varTable, funTable)
	switch rType {
	case symtable.NUMBER:
		switch n.Op {
		case token.ADD:
			return +rVal.(float64), symtable.NUMBER, false
		case token.SUBTRACT:
			return -rVal.(float64), symtable.NUMBER, false
		}
	case symtable.BOOL:
		switch n.Op {
		case token.NOT:
			return !rVal.(bool), symtable.BOOL, false
		}
	}
	panic("Invalid data type")
}

func (n BinaryExpr) Eval(varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	lVal, lType, _ := n.Left.Eval(varTable, funTable)
	if lType == symtable.BOOL && lVal.(bool) && n.Op == token.OR {
		return true, symtable.BOOL, false
	}

	rVal, rType, _ := n.Right.Eval(varTable, funTable)
	if lType != rType {
		panic("Data types do not match")
	}

	switch lType {
	case symtable.NUMBER:
		switch n.Op {
		case token.ADD:
			return lVal.(float64) + rVal.(float64), symtable.NUMBER, false
		case token.SUBTRACT:
			return lVal.(float64) - rVal.(float64), symtable.NUMBER, false
		case token.MULTIPLY:
			return lVal.(float64) * rVal.(float64), symtable.NUMBER, false
		case token.DIVIDE:
			return lVal.(float64) / rVal.(float64), symtable.NUMBER, false
		case token.EQUAL:
			return lVal.(float64) == rVal.(float64), symtable.BOOL, false
		case token.LESS:
			return lVal.(float64) < rVal.(float64), symtable.BOOL, false
		case token.LESS_EQ:
			return lVal.(float64) <= rVal.(float64), symtable.BOOL, false
		case token.GREATER:
			return lVal.(float64) > rVal.(float64), symtable.BOOL, false
		case token.GREATER_EQ:
			return lVal.(float64) >= rVal.(float64), symtable.BOOL, false
		}
	case symtable.STRING:
		switch n.Op {
		case token.ADD:
			return lVal.(string) + rVal.(string), symtable.STRING, false
		case token.EQUAL:
			return lVal.(string) == rVal.(string), symtable.BOOL, false
		}
	case symtable.BOOL:
		switch n.Op {
		case token.AND:
			return lVal.(bool) && rVal.(bool), symtable.BOOL, false
		case token.OR:
			return lVal.(bool) || rVal.(bool), symtable.BOOL, false
		}
	}
	panic("Invalid data type")
}
