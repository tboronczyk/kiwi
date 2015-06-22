package runtime

import (
	"fmt"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/symtable"
	"github.com/tboronczyk/kiwi/token"
	"github.com/tboronczyk/kiwi/util"
	"math"
	"strconv"
	"strings"
)

type stackEntry struct {
	Type  symtable.DataType
	Value interface{}
}

type Runtime struct {
	Return bool
	funcs  util.Stack
	vars   util.Stack
	stack  util.Stack
}

func New() *Runtime {
	r := &Runtime{
		Return: false,
		funcs:  util.NewStack(),
		vars:   util.NewStack(),
		stack:  util.NewStack(),
	}
	r.vars.Push(symtable.New())
	r.funcs.Push(symtable.New())
	// register builtin functions
	for _, n := range funcSigs {
		r.funcSet(n.Name, n, symtable.BUILTIN)
	}
	return r
}

func (r *Runtime) funcSet(n string, v interface{}, t symtable.DataType) {
	s := r.funcs.Peek().(symtable.SymTable)
	s.Set(n, v, t)
}

func (r *Runtime) funcGet(n string) (interface{}, symtable.DataType, bool) {
	return r.funcs.Peek().(symtable.SymTable).Get(n)
}

func (r *Runtime) varSet(n string, v interface{}, t symtable.DataType) {
	s := r.vars.Peek().(symtable.SymTable)
	s.Set(n, v, t)
}

func (r *Runtime) varGet(n string) (interface{}, symtable.DataType, bool) {
	return r.vars.Peek().(symtable.SymTable).Get(n)
}

func (r *Runtime) enterScope(s symtable.SymTable) {
	r.vars.Push(s)
}

func (r *Runtime) leaveScope() {
	r.vars.Pop()
}

func (r *Runtime) pushStack(v interface{}, t symtable.DataType) {
	r.stack.Push(stackEntry{Value: v, Type: t})
}

func (r *Runtime) popStack() stackEntry {
	return r.stack.Pop().(stackEntry)
}

func (r *Runtime) VisitValueExpr(n ast.ValueExpr) {
	switch n.Type {
	case token.NUMBER:
		val, _ := strconv.ParseFloat(n.Value, 64)
		r.pushStack(val, symtable.NUMBER)
		break
	case token.BOOL:
		r.pushStack(strings.ToUpper(n.Value) == "TRUE", symtable.BOOL)
		break
	case token.STRING:
		r.pushStack(n.Value, symtable.STRING)
		break
	}
}

func (r *Runtime) VisitCastExpr(n ast.CastExpr) {
	n.Expr.Accept(r)
	expr := r.popStack()
	switch strings.ToUpper(n.Cast) {
	case "STRING":
		switch expr.Type {
		case symtable.STRING:
			r.pushStack(expr.Value, symtable.STRING)
			break
		case symtable.NUMBER:
			val := fmt.Sprintf("%f", expr.Value.(float64))
			val = strings.TrimRight(val, "0")
			val = strings.TrimRight(val, ".")
			r.pushStack(val, symtable.STRING)
			break
		case symtable.BOOL:
			r.pushStack(strconv.FormatBool(expr.Value.(bool)), symtable.STRING)
			break
		}
		break
	case "NUMBER":
		switch expr.Type {
		case symtable.STRING:
			val, err := strconv.ParseFloat(expr.Value.(string), 64)
			if err != nil {
				val = 0.0
			}
			r.pushStack(val, symtable.NUMBER)
			break
		case symtable.NUMBER:
			r.pushStack(expr.Value, symtable.NUMBER)
			break
		case symtable.BOOL:
			val := 0.0
			if expr.Value.(bool) {
				val = 1.0
			}
			r.pushStack(val, symtable.NUMBER)
			break
		}
		break
	case "BOOL":
		switch expr.Type {
		case symtable.STRING:
			val := strings.ToUpper(expr.Value.(string)) != "FALSE" &&
				strings.TrimSpace(expr.Value.(string)) != ""
			r.pushStack(val, symtable.BOOL)
			break
		case symtable.NUMBER:
			r.pushStack(expr.Value != 0.0, symtable.BOOL)
			break
		case symtable.BOOL:
			r.pushStack(expr.Value, symtable.BOOL)
			break
		}
		break
	default:
		panic("Invalid cast")
	}
}

func (r *Runtime) VisitVariableExpr(n ast.VariableExpr) {
	val, dtype, ok := r.varGet(n.Name)
	if !ok {
		panic("Variable is not set")
	}
	r.pushStack(val, dtype)
}

func (r *Runtime) VisitUnaryExpr(n ast.UnaryExpr) {
	n.Right.Accept(r)
	right := r.popStack()

	switch right.Type {
	case symtable.NUMBER:
		switch n.Op {
		case token.ADD:
			r.pushStack(math.Abs(right.Value.(float64)), symtable.NUMBER)
			break
		case token.SUBTRACT:
			r.pushStack(-right.Value.(float64), symtable.NUMBER)
			break
		}
	case symtable.BOOL:
		switch n.Op {
		case token.NOT:
			r.pushStack(right.Value.(bool) == false, symtable.BOOL)
			break
		}
	}
}

func (r *Runtime) VisitBinaryExpr(n ast.BinaryExpr) {
	n.Left.Accept(r)
	left := r.popStack()
	if left.Type == symtable.BOOL {
		if left.Value.(bool) && n.Op == token.OR {
			r.pushStack(true, symtable.BOOL)
			return
		}
		if !left.Value.(bool) && n.Op == token.AND {
			r.pushStack(false, symtable.BOOL)
			return
		}
	}

	n.Right.Accept(r)
	right := r.popStack()
	if left.Type != right.Type {
		panic("Data types do not match")
	}

	switch left.Type {
	case symtable.NUMBER:
		switch n.Op {
		case token.ADD:
			r.pushStack(left.Value.(float64)+right.Value.(float64), symtable.NUMBER)
			break
		case token.SUBTRACT:
			r.pushStack(left.Value.(float64)-right.Value.(float64), symtable.NUMBER)
			break
		case token.MULTIPLY:
			r.pushStack(left.Value.(float64)*right.Value.(float64), symtable.NUMBER)
			break
		case token.DIVIDE:
			r.pushStack(left.Value.(float64)/right.Value.(float64), symtable.NUMBER)
			break
		case token.MODULO:
			r.pushStack(math.Mod(left.Value.(float64), right.Value.(float64)), symtable.NUMBER)
			break
		case token.EQUAL:
			r.pushStack(left.Value.(float64) == right.Value.(float64), symtable.BOOL)
			break
		case token.LESS:
			r.pushStack(left.Value.(float64) < right.Value.(float64), symtable.BOOL)
			break
		case token.LESS_EQ:
			r.pushStack(left.Value.(float64) <= right.Value.(float64), symtable.BOOL)
			break
		case token.GREATER:
			r.pushStack(left.Value.(float64) > right.Value.(float64), symtable.BOOL)
			break
		case token.GREATER_EQ:
			r.pushStack(left.Value.(float64) >= right.Value.(float64), symtable.BOOL)
			break
		}
	case symtable.STRING:
		switch n.Op {
		case token.ADD:
			r.pushStack(left.Value.(string)+right.Value.(string), symtable.STRING)
			break
		case token.EQUAL:
			r.pushStack(left.Value.(string) == right.Value.(string), symtable.BOOL)
			break
		}
	case symtable.BOOL:
		switch n.Op {
		case token.AND:
			r.pushStack(left.Value.(bool) && right.Value.(bool), symtable.BOOL)
			break
		case token.OR:
			r.pushStack(left.Value.(bool) || right.Value.(bool), symtable.BOOL)
			break
		}
	}
}

func (r *Runtime) VisitFuncCall(n ast.FuncCall) {
	defer func() {
		r.Return = false
	}()

	fun, dtype, ok := r.funcGet(n.Name)
	if !ok {
		panic("Function not defined")
	}
	if len(fun.(ast.FuncDef).Args) != len(n.Args) {
		panic("Function arity mis-match")
	}

	s := symtable.New()
	for i, name := range fun.(ast.FuncDef).Args {
		n.Args[i].Accept(r)
		arg := r.popStack()
		s.Set(name, arg.Value, arg.Type)
	}
	r.enterScope(s)

	switch dtype {
	case symtable.USRFUNC:
		for _, stmt := range fun.(ast.FuncDef).Body {
			stmt.Accept(r)
			if r.Return {
				break
			}
		}
		break
	case symtable.BUILTIN:
		builtinFuncs[n.Name](r)
		break
	}
	r.leaveScope()
}

func (r *Runtime) VisitAssignStmt(n ast.AssignStmt) {
	n.Expr.Accept(r)
	expr := r.popStack()
	r.varSet(n.Name, expr.Value, expr.Type)
}

func (r *Runtime) VisitFuncDef(n ast.FuncDef) {
	r.funcSet(n.Name, n, symtable.USRFUNC)
}

func (r *Runtime) VisitIfStmt(n ast.IfStmt) {
	n.Condition.Accept(r)
	cond := r.popStack()
	if cond.Type != symtable.BOOL {
		panic("Non-bool result used for condition")
	}
	if cond.Value.(bool) {
		for _, stmt := range n.Body {
			stmt.Accept(r)
			if r.Return {
				return
			}
		}
	} else if n.Else != nil {
		n.Else.Accept(r)
	}
}

func (r *Runtime) VisitReturnStmt(n ast.ReturnStmt) {
	n.Expr.Accept(r)
	r.Return = true
}

func (r *Runtime) VisitWhileStmt(n ast.WhileStmt) {
	for {
		n.Condition.Accept(r)
		cond := r.popStack()
		if cond.Type != symtable.BOOL {
			panic("Non-bool result used for condition")
		}
		if !cond.Value.(bool) {
			return
		}
		for _, stmt := range n.Body {
			stmt.Accept(r)
			if r.Return {
				return
			}
		}
	}
}
