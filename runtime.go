// Package runtime provides runtime support for executing Kiwi programs.
package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type (
	Runtime struct {
		stack      Stack
		curScope   *Scope
		scopeStack Stack
	}

	params []Entry
)

func NewRuntime() *Runtime {
	r := &Runtime{
		stack:      NewStack(),
		curScope:   NewScope(),
		scopeStack: NewStack(),
	}

	for n, f := range builtins {
		r.curScope.SetFunc(n, Entry{Value: f, DataType: BUILTIN})
	}
	return r
}

func (r *Runtime) VisitAddNode(n *AstAddNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == NUMBER && right.DataType == NUMBER {
		r.stack.Push(Entry{
			Value:    left.Value.(float64) + right.Value.(float64),
			DataType: NUMBER,
		})
		return
	}
	if left.DataType == STRING && right.DataType == STRING {
		r.stack.Push(Entry{
			Value:    left.Value.(string) + right.Value.(string),
			DataType: STRING,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitAndNode(n *AstAndNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)
	// short-circuit if false
	if left.DataType == BOOL && !left.Value.(bool) {
		r.stack.Push(Entry{
			Value:    false,
			DataType: BOOL,
		})
		return
	}

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == BOOL && right.DataType == BOOL {
		r.stack.Push(Entry{
			Value:    left.Value.(bool) == right.Value.(bool),
			DataType: BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitAssignNode(n *AstAssignNode) {
	n.Expr.Accept(r)
	v := r.stack.Pop().(Entry)

	// preserve datatype if the variable is already set
	e, ok := r.curScope.GetVar(n.Name)
	if ok {
		if e.DataType != v.DataType {
			panic("value type does not match variable type")
		}
	}
	r.curScope.SetVar(n.Name, v)
}

func (r *Runtime) VisitBoolNode(n *AstBoolNode) {
	r.stack.Push(Entry{
		Value:    n.Value,
		DataType: BOOL,
	})
}

func (r *Runtime) VisitCastNode(n *AstCastNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(Entry)
	switch strings.ToUpper(n.Cast) {
	case "STR":
		switch e.DataType {
		case STRING:
			break
		case NUMBER:
			val := fmt.Sprintf("%f", e.Value.(float64))
			val = strings.TrimRight(val, "0")
			val = strings.TrimRight(val, ".")
			e.Value = val
			break
		case BOOL:
			e.Value = strconv.FormatBool(e.Value.(bool))
			break
		}
		e.DataType = STRING
		break
	case "NUM":
		switch e.DataType {
		case STRING:
			e.Value, _ = strconv.ParseFloat(e.Value.(string), 64)
			break
		case NUMBER:
			break
		case BOOL:
			val := 0.0
			if e.Value.(bool) {
				val = 1.0
			}
			e.Value = val
			break
		}
		e.DataType = NUMBER
		break
	case "BOOL":
		switch e.DataType {
		case STRING:
			value := strings.ToUpper(e.Value.(string)) != "FALSE" &&
				strings.TrimSpace(e.Value.(string)) != ""
			e.Value = value
			break
		case NUMBER:
			e.Value = e.Value.(float64) != 0.0
			break
		case BOOL:
			break
		}
		e.DataType = BOOL
		break
	}
	r.stack.Push(e)
}

func (r *Runtime) VisitDivideNode(n *AstDivideNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == NUMBER && right.DataType == NUMBER {
		r.stack.Push(Entry{
			Value:    left.Value.(float64) / right.Value.(float64),
			DataType: NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitEqualNode(n *AstEqualNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == right.DataType {
		r.stack.Push(Entry{
			Value:    left.Value == right.Value,
			DataType: BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitFuncCallNode(n *AstFuncCallNode) {
	e, ok := r.curScope.GetFunc(n.Name)
	if !ok {
		panic("Function not defined")
	}

	var p params
	for _, arg := range n.Args {
		arg.Accept(r)
		p = append(p, r.stack.Pop().(Entry))
	}

	if e.DataType == BUILTIN {
		builtins[n.Name](&r.stack, p)
		return
	}

	f := e.Value.(*AstFuncDefNode)
	if len(n.Args) != len(f.Args) {
		panic("wrong number of arguments in function call")
	}

	r.scopeStack.Push(r.curScope)
	r.curScope = NewScopeClone(f.Scope)
	for i, arg := range f.Args {
		r.curScope.SetVar(arg, p[i])
	}
	for _, stmt := range f.Body {
		stmt.Accept(r)
		if r.stack.Size() > 0 {
			break
		}
	}
	r.curScope = r.scopeStack.Pop().(*Scope)
}

func (r *Runtime) VisitFuncDefNode(n *AstFuncDefNode) {
	// nothing to do
}

func (r *Runtime) VisitGreaterEqualNode(n *AstGreaterEqualNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == NUMBER && right.DataType == NUMBER {
		r.stack.Push(Entry{
			Value:    left.Value.(float64) >= right.Value.(float64),
			DataType: BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitGreaterNode(n *AstGreaterNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == NUMBER && right.DataType == NUMBER {
		r.stack.Push(Entry{
			Value:    left.Value.(float64) > right.Value.(float64),
			DataType: BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitIfNode(n *AstIfNode) {
	n.Cond.Accept(r)
	cond := r.stack.Pop().(Entry)
	if cond.DataType != BOOL {
		panic("non-bool expression used as condition")
	}
	if cond.Value.(bool) {
		for _, stmt := range n.Body {
			stmt.Accept(r)
			if r.stack.Size() > 0 {
				break
			}
		}
	} else if n.Else != nil {
		for _, stmt := range n.Else {
			stmt.Accept(r)
			if r.stack.Size() > 0 {
				break
			}
		}
	}
}

func (r *Runtime) VisitLessEqualNode(n *AstLessEqualNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == NUMBER && right.DataType == NUMBER {
		r.stack.Push(Entry{
			Value:    left.Value.(float64) <= right.Value.(float64),
			DataType: BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitLessNode(n *AstLessNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == NUMBER && right.DataType == NUMBER {
		r.stack.Push(Entry{
			Value:    left.Value.(float64) < right.Value.(float64),
			DataType: BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitModuloNode(n *AstModuloNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == NUMBER && right.DataType == NUMBER {
		r.stack.Push(Entry{
			Value:    math.Mod(left.Value.(float64), right.Value.(float64)),
			DataType: NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitMultiplyNode(n *AstMultiplyNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == NUMBER && right.DataType == NUMBER {
		r.stack.Push(Entry{
			Value:    left.Value.(float64) * right.Value.(float64),
			DataType: NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNegativeNode(n *AstNegativeNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(Entry)

	if e.DataType == NUMBER {
		r.stack.Push(Entry{
			Value:    -e.Value.(float64),
			DataType: NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNotEqualNode(n *AstNotEqualNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == right.DataType {
		r.stack.Push(Entry{
			Value:    left.Value != right.Value,
			DataType: BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNotNode(n *AstNotNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(Entry)

	if e.DataType == BOOL {
		r.stack.Push(Entry{
			Value:    !e.Value.(bool),
			DataType: BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNumberNode(n *AstNumberNode) {
	r.stack.Push(Entry{
		Value:    n.Value,
		DataType: NUMBER,
	})
}

func (r *Runtime) VisitOrNode(n *AstOrNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)
	// short-circuit if true
	if left.DataType == BOOL && left.Value.(bool) {
		r.stack.Push(Entry{
			Value:    true,
			DataType: BOOL,
		})
		return
	}

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == BOOL && right.DataType == BOOL {
		r.stack.Push(Entry{
			Value:    left.Value.(bool) || right.Value.(bool),
			DataType: BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitPositiveNode(n *AstPositiveNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(Entry)

	if e.DataType == NUMBER {
		r.stack.Push(Entry{
			Value:    math.Abs(e.Value.(float64)),
			DataType: NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitProgramNode(n *AstProgramNode) {
	n.Scope.Parent = r.curScope
	r.scopeStack.Push(r.curScope)
	r.curScope = n.Scope

	for _, stmt := range n.Stmts {
		stmt.Accept(r)
	}

	r.curScope = r.scopeStack.Pop().(*Scope)
}

func (r *Runtime) VisitReturnNode(n *AstReturnNode) {
	n.Expr.Accept(r)
}

func (r *Runtime) VisitStringNode(n *AstStringNode) {
	r.stack.Push(Entry{
		Value:    n.Value,
		DataType: STRING,
	})
}

func (r *Runtime) VisitSubtractNode(n *AstSubtractNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(Entry)

	if left.DataType == NUMBER && right.DataType == NUMBER {
		r.stack.Push(Entry{
			Value:    left.Value.(float64) - right.Value.(float64),
			DataType: NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitVariableNode(n *AstVariableNode) {
	expr, ok := r.curScope.GetVar(n.Name)
	if !ok {
		panic("variable is not defined")
	}
	r.stack.Push(expr)
}

func (r *Runtime) VisitWhileNode(n *AstWhileNode) {
	for {
		n.Cond.Accept(r)
		cond := r.stack.Pop().(Entry)
		if cond.DataType != BOOL {
			panic("non-bool expression used as condition")
		}
		if !cond.Value.(bool) {
			return
		}
		for _, stmt := range n.Body {
			stmt.Accept(r)
			if r.stack.Size() > 0 {
				return
			}
		}
	}
}
