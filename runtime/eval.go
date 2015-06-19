package runtime

import (
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/symtable"
	"github.com/tboronczyk/kiwi/token"
	"strconv"
	"fmt"
	"strings"
)

func Eval(node ast.Node, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	switch node.(type) {
	case ast.AssignStmt:
		return evalAssignStmt(node.(ast.AssignStmt), varTable, funTable)
	case ast.FuncDef:
		return evalFuncDef(node.(ast.FuncDef), varTable, funTable)
	case ast.FuncCall:
		return evalFuncCall(node.(ast.FuncCall), varTable, funTable)
	case ast.IfStmt:
		return evalIfStmt(node.(ast.IfStmt), varTable, funTable)
	case ast.ReturnStmt:
		return evalReturnStmt(node.(ast.ReturnStmt), varTable, funTable)
	case ast.WhileStmt:
		return evalWhileStmt(node.(ast.WhileStmt), varTable, funTable)
	case ast.VariableExpr:
		return evalVariableExpr(node.(ast.VariableExpr), varTable, funTable)
	case ast.ValueExpr:
		return evalValueExpr(node.(ast.ValueExpr), varTable, funTable)
	case ast.UnaryExpr:
		return evalUnaryExpr(node.(ast.UnaryExpr), varTable, funTable)
	case ast.BinaryExpr:
		return evalBinaryExpr(node.(ast.BinaryExpr), varTable, funTable)
	case ast.CastExpr:
		return evalCastExpr(node.(ast.CastExpr), varTable, funTable)
	}
	panic("Whoops!")
}

func evalAssignStmt(node ast.AssignStmt, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	value, dtype, _ := Eval(node.Expr, varTable, funTable)
	varTable.Set(node.Name, value, dtype)
	return value, dtype, false
}

func evalFuncDef(node ast.FuncDef, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	funTable.Set(node.Name, node, symtable.USRFUNC)
	return true, symtable.BOOL, false
}

func evalFuncCall(node ast.FuncCall, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	f, t, ok := funTable.Get(node.Name)
	if !ok {
		panic("Function not defined")
	}
	if len(f.(ast.FuncDef).Args) != len(node.Args) {
		panic("Function arity mis-match")
	}
	s := symtable.New()
	for i, name := range f.(ast.FuncDef).Args {
		value, dtype, _ := Eval(node.Args[i], varTable, funTable)
		s.Set(name, value, dtype)
	}
	switch t {
	case symtable.USRFUNC:
		for _, stmt := range f.(ast.FuncDef).Body {
			value, dtype, isReturn := Eval(stmt, s, funTable)
			if isReturn {
				return value, dtype, isReturn
			}
		}
	case symtable.BUILTIN:
		return builtins[node.Name](s, funTable)
	}

	return true, symtable.BOOL, false
}

func evalIfStmt(n ast.IfStmt, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	value, dtype, _ := Eval(n.Condition, varTable, funTable)
	if dtype != symtable.BOOL {
		panic("Invalid data type")
	}
	if value.(bool) {
		var isReturn bool
		for _, stmt := range n.Body {
			value, dtype, isReturn = Eval(stmt, varTable, funTable)
			if isReturn {
				return value, dtype, true
			}
		}
	}
	return true, symtable.BOOL, false
}

func evalReturnStmt(node ast.ReturnStmt, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	value, dtype, _ := Eval(node.Expr, varTable, funTable)
	return value, dtype, true
}

func evalWhileStmt(node ast.WhileStmt, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	for {
		value, dtype, _ := Eval(node.Condition, varTable, funTable)
		if dtype != symtable.BOOL {
			panic("Invalid data type")
		}
		if !value.(bool) {
			return true, symtable.BOOL, false
		}
		var isReturn bool
		for _, stmt := range node.Body {
			value, dtype, isReturn = Eval(stmt, varTable, funTable)
			if isReturn {
				return value, dtype, true
			}
		}
	}
}

func evalValueExpr(node ast.ValueExpr, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	switch node.Type {
	case token.NUMBER:
		value, _ := strconv.ParseFloat(node.Value, 64)
		return value, symtable.NUMBER, false
	case token.BOOL:
		return node.Value == "TRUE", symtable.BOOL, false
	}
	return node.Value, symtable.STRING, false
}

func evalVariableExpr(node ast.VariableExpr, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	value, dtype, ok := varTable.Get(node.Name)
	if !ok {
		panic("Variable not set")
	}
	return value, dtype, false
}

func evalUnaryExpr(node ast.UnaryExpr, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	rVal, rType, _ := Eval(node.Right, varTable, funTable)
	switch rType {
	case symtable.NUMBER:
		switch node.Op {
		case token.ADD:
			return +rVal.(float64), symtable.NUMBER, false
		case token.SUBTRACT:
			return -rVal.(float64), symtable.NUMBER, false
		}
	case symtable.BOOL:
		switch node.Op {
		case token.NOT:
			return !rVal.(bool), symtable.BOOL, false
		}
	}
	panic("Invalid data type")
}

func evalBinaryExpr(node ast.BinaryExpr, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	lVal, lType, _ := Eval(node.Left, varTable, funTable)
	if lType == symtable.BOOL && lVal.(bool) && node.Op == token.OR {
		return true, symtable.BOOL, false
	}

	rVal, rType, _ := Eval(node.Right, varTable, funTable)
	if lType != rType {
		panic("Data types do not match")
	}

	switch lType {
	case symtable.NUMBER:
		switch node.Op {
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
		switch node.Op {
		case token.ADD:
			return lVal.(string) + rVal.(string), symtable.STRING, false
		case token.EQUAL:
			return lVal.(string) == rVal.(string), symtable.BOOL, false
		}
	case symtable.BOOL:
		switch node.Op {
		case token.AND:
			return lVal.(bool) && rVal.(bool), symtable.BOOL, false
		case token.OR:
			return lVal.(bool) || rVal.(bool), symtable.BOOL, false
		}
	}
	panic("Invalid data type")
}

func evalCastExpr(node ast.CastExpr, varTable, funTable symtable.SymTable) (interface{}, symtable.DataType, bool) {
	value, dtype, isReturn := Eval(node.Expr, varTable, funTable)
	switch strings.ToUpper(node.Cast) {
	case "STRING":
		switch dtype {
		case symtable.STRING:
			return value, symtable.STRING, isReturn
		case symtable.NUMBER:
			return fmt.Sprintf("%f", value.(float64)), symtable.STRING, isReturn
		case symtable.BOOL:
			return strconv.FormatBool(value.(bool)), symtable.STRING, isReturn
		}
	case "NUMBER":
		switch dtype {
		case symtable.STRING:
			val, err := strconv.ParseFloat(value.(string), 64)
			if err != nil {
				val = 0.0
			}
			return val, symtable.NUMBER, isReturn
		case symtable.NUMBER:
			return value.(float64), symtable.NUMBER, isReturn
		case symtable.BOOL:
			if value.(bool) {
			return 1.0, symtable.NUMBER, isReturn
			}
			return 0.0, symtable.NUMBER, isReturn
		}
	case "BOOL":
		switch dtype {
		case symtable.STRING:
			val := strings.ToUpper(value.(string)) != "FALSE" &&
				strings.TrimSpace(value.(string)) != ""
			return val, symtable.BOOL, isReturn
		case symtable.NUMBER:
			return value.(float64) != 0.0, symtable.BOOL, isReturn
		case symtable.BOOL:
			return value.(bool), symtable.BOOL, isReturn
		}
	}
	panic("Invalid cast")
}
