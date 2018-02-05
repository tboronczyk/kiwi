package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type (
	Runtime struct {
		stack      Stack
		scopeStack Stack
		currScope  *Scope
		stdin      *os.File
		stdout     *os.File
		stderr     *os.File
	}

	params []ScopeEntry
)

func NewRuntime(stdin, stdout, stderr *os.File) *Runtime {
	r := &Runtime{
		NewStack(),
		NewStack(),
		NewScope(),
		stdin, stdout, stderr,
	}

	for name, fn := range builtins {
		r.currScope.SetFunc(name, ScopeEntry{TypBuiltin, fn})
	}
	return r
}

func (r *Runtime) VisitAddNode(n *AstAddNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == TypNumber && right.DataType == TypNumber {
		r.stack.Push(ScopeEntry{
			TypNumber,
			left.Value.(float64) + right.Value.(float64),
		})
		return
	}
	if left.DataType == TypString && right.DataType == TypString {
		r.stack.Push(ScopeEntry{
			TypString,
			left.Value.(string) + right.Value.(string),
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitAndNode(n *AstAndNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)
	// short-circuit if false
	if left.DataType == TypBool && !left.Value.(bool) {
		r.stack.Push(ScopeEntry{TypBool, false})
		return
	}

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == TypBool && right.DataType == TypBool {
		r.stack.Push(ScopeEntry{
			TypBool,
			left.Value.(bool) == right.Value.(bool),
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitAssignNode(n *AstAssignNode) {
	n.Expr.Accept(r)
	v := r.stack.Pop().(ScopeEntry)

	// preserve datatype if the variable is already set
	e, ok := r.currScope.GetVar(n.Name)
	if ok {
		if e.DataType != v.DataType {
			panic("value type does not match variable type")
		}
	}
	r.currScope.SetVar(n.Name, v)
}

func (r *Runtime) VisitBoolNode(n *AstBoolNode) {
	r.stack.Push(ScopeEntry{TypBool, n.Value})
}

func (r *Runtime) VisitCastNode(n *AstCastNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(ScopeEntry)
	switch strings.ToUpper(n.Cast) {
	case "STR":
		switch e.DataType {
		case TypString:
			break
		case TypNumber:
			val := fmt.Sprintf("%f", e.Value.(float64))
			val = strings.TrimRight(val, "0")
			val = strings.TrimRight(val, ".")
			e.Value = val
			break
		case TypBool:
			e.Value = strconv.FormatBool(e.Value.(bool))
			break
		}
		e.DataType = TypString
		break
	case "NUM":
		switch e.DataType {
		case TypString:
			e.Value, _ = strconv.ParseFloat(e.Value.(string), 64)
			break
		case TypNumber:
			break
		case TypBool:
			val := 0.0
			if e.Value.(bool) {
				val = 1.0
			}
			e.Value = val
			break
		}
		e.DataType = TypNumber
		break
	case "BOOL":
		switch e.DataType {
		case TypString:
			value := strings.ToUpper(e.Value.(string)) != "FALSE" &&
				strings.TrimSpace(e.Value.(string)) != ""
			e.Value = value
			break
		case TypNumber:
			e.Value = e.Value.(float64) != 0.0
			break
		case TypBool:
			break
		}
		e.DataType = TypBool
		break
	}
	r.stack.Push(e)
}

func (r *Runtime) VisitDivideNode(n *AstDivideNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == TypNumber && right.DataType == TypNumber {
		r.stack.Push(ScopeEntry{
			TypNumber,
			left.Value.(float64) / right.Value.(float64),
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitEqualNode(n *AstEqualNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == right.DataType {
		r.stack.Push(ScopeEntry{
			TypBool,
			left.Value == right.Value,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitFuncCallNode(n *AstFuncCallNode) {
	e, ok := r.currScope.GetFunc(n.Name)
	if !ok {
		panic("Function not defined")
	}

	var p params
	for _, arg := range n.Args {
		arg.Accept(r)
		p = append(p, r.stack.Pop().(ScopeEntry))
	}

	if e.DataType == TypBuiltin {
		builtins[n.Name](&r.stack, p, r.stdin, r.stdout, r.stderr)
		return
	}

	f := e.Value.(*AstFuncDefNode)
	if len(n.Args) != len(f.Args) {
		panic("wrong number of arguments in function call")
	}

	r.scopeStack.Push(r.currScope)
	r.currScope = f.Scope.EmptyVarCopy()
	for i, arg := range f.Args {
		r.currScope.SetVar(arg, p[i])
	}
	for _, stmt := range f.Body {
		stmt.Accept(r)
		if r.stack.Size() > 0 {
			break
		}
	}
	r.currScope = r.scopeStack.Pop().(*Scope)
}

func (r *Runtime) VisitFuncDefNode(n *AstFuncDefNode) {
	// nothing to do
}

func (r *Runtime) VisitGreaterEqualNode(n *AstGreaterEqualNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == TypNumber && right.DataType == TypNumber {
		r.stack.Push(ScopeEntry{
			TypBool,
			left.Value.(float64) >= right.Value.(float64),
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitGreaterNode(n *AstGreaterNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == TypNumber && right.DataType == TypNumber {
		r.stack.Push(ScopeEntry{
			TypBool,
			left.Value.(float64) > right.Value.(float64),
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitIfNode(n *AstIfNode) {
	n.Cond.Accept(r)
	cond := r.stack.Pop().(ScopeEntry)
	if cond.DataType != TypBool {
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
	left := r.stack.Pop().(ScopeEntry)

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == TypNumber && right.DataType == TypNumber {
		r.stack.Push(ScopeEntry{
			TypBool,
			left.Value.(float64) <= right.Value.(float64),
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitLessNode(n *AstLessNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == TypNumber && right.DataType == TypNumber {
		r.stack.Push(ScopeEntry{
			TypBool,
			left.Value.(float64) < right.Value.(float64),
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitModuloNode(n *AstModuloNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == TypNumber && right.DataType == TypNumber {
		r.stack.Push(ScopeEntry{
			TypNumber,
			math.Mod(left.Value.(float64), right.Value.(float64)),
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitMultiplyNode(n *AstMultiplyNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == TypNumber && right.DataType == TypNumber {
		r.stack.Push(ScopeEntry{
			TypNumber,
			left.Value.(float64) * right.Value.(float64),
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNegativeNode(n *AstNegativeNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(ScopeEntry)

	if e.DataType == TypNumber {
		r.stack.Push(ScopeEntry{TypNumber, -e.Value.(float64)})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNotEqualNode(n *AstNotEqualNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == right.DataType {
		r.stack.Push(ScopeEntry{
			TypBool,
			left.Value != right.Value,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNotNode(n *AstNotNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(ScopeEntry)

	if e.DataType == TypBool {
		r.stack.Push(ScopeEntry{TypBool, !e.Value.(bool)})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNumberNode(n *AstNumberNode) {
	r.stack.Push(ScopeEntry{TypNumber, n.Value})
}

func (r *Runtime) VisitOrNode(n *AstOrNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)
	// short-circuit if true
	if left.DataType == TypBool && left.Value.(bool) {
		r.stack.Push(ScopeEntry{TypBool, true})
		return
	}

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == TypBool && right.DataType == TypBool {
		r.stack.Push(ScopeEntry{
			TypBool,
			left.Value.(bool) || right.Value.(bool),
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitPositiveNode(n *AstPositiveNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(ScopeEntry)

	if e.DataType == TypNumber {
		r.stack.Push(ScopeEntry{TypNumber, math.Abs(e.Value.(float64))})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitProgramNode(n *AstProgramNode) {
	n.Scope.parent = r.currScope
	r.scopeStack.Push(r.currScope)
	r.currScope = n.Scope

	for _, stmt := range n.Stmts {
		stmt.Accept(r)
	}

	r.currScope = r.scopeStack.Pop().(*Scope)
}

func (r *Runtime) VisitReturnNode(n *AstReturnNode) {
	n.Expr.Accept(r)
}

func (r *Runtime) VisitStringNode(n *AstStringNode) {
	r.stack.Push(ScopeEntry{TypString, n.Value})
}

func (r *Runtime) VisitSubtractNode(n *AstSubtractNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ScopeEntry)

	n.Right.Accept(r)
	right := r.stack.Pop().(ScopeEntry)

	if left.DataType == TypNumber && right.DataType == TypNumber {
		r.stack.Push(ScopeEntry{
			TypNumber,
			left.Value.(float64) - right.Value.(float64),
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitVariableNode(n *AstVariableNode) {
	expr, ok := r.currScope.GetVar(n.Name)
	if !ok {
		panic("variable is not defined")
	}
	r.stack.Push(expr)
}

func (r *Runtime) VisitWhileNode(n *AstWhileNode) {
	for {
		n.Cond.Accept(r)
		cond := r.stack.Pop().(ScopeEntry)
		if cond.DataType != TypBool {
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
