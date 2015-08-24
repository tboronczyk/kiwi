// Package runtime provides runtime support for executing Kiwi programs.
package runtime

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scope"
	"github.com/tboronczyk/kiwi/types"
	"github.com/tboronczyk/kiwi/util"
)

type (
	Runtime struct {
		stack      util.Stack
		curScope   *scope.Scope
		scopeStack util.Stack
	}

	params []scope.Entry
)

func New() *Runtime {
	r := &Runtime{
		stack:      util.NewStack(),
		curScope:   scope.New(),
		scopeStack: util.NewStack(),
	}

	for n, f := range builtins {
		r.curScope.SetFunc(n, scope.Entry{Value: f, DataType: types.BUILTIN})
	}
	return r
}

func (r *Runtime) VisitAddNode(n *ast.AddNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == types.NUMBER && right.DataType == types.NUMBER {
		r.stack.Push(scope.Entry{
			Value:    left.Value.(float64) + right.Value.(float64),
			DataType: types.NUMBER,
		})
		return
	}
	if left.DataType == types.STRING && right.DataType == types.STRING {
		r.stack.Push(scope.Entry{
			Value:    left.Value.(string) + right.Value.(string),
			DataType: types.STRING,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitAndNode(n *ast.AndNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)
	// short-circuit if false
	if left.DataType == types.BOOL && !left.Value.(bool) {
		r.stack.Push(scope.Entry{
			Value:    false,
			DataType: types.BOOL,
		})
		return
	}

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == types.BOOL && right.DataType == types.BOOL {
		r.stack.Push(scope.Entry{
			Value:    left.Value.(bool) == right.Value.(bool),
			DataType: types.BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitAssignNode(n *ast.AssignNode) {
	n.Expr.Accept(r)
	v := r.stack.Pop().(scope.Entry)

	// preserve datatype if the variable is already set
	e, ok := r.curScope.GetVar(n.Name)
	if ok {
		if e.DataType != v.DataType {
			panic("value type does not match variable type")
		}
	}
	r.curScope.SetVar(n.Name, v)
}

func (r *Runtime) VisitBoolNode(n *ast.BoolNode) {
	r.stack.Push(scope.Entry{
		Value:    n.Value,
		DataType: types.BOOL,
	})
}

func (r *Runtime) VisitCastNode(n *ast.CastNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	switch strings.ToUpper(n.Cast) {
	case "STR":
		switch e.DataType {
		case types.STRING:
			break
		case types.NUMBER:
			val := fmt.Sprintf("%f", e.Value.(float64))
			val = strings.TrimRight(val, "0")
			val = strings.TrimRight(val, ".")
			e.Value = val
			break
		case types.BOOL:
			e.Value = strconv.FormatBool(e.Value.(bool))
			break
		}
		e.DataType = types.STRING
		break
	case "NUM":
		switch e.DataType {
		case types.STRING:
			e.Value, _ = strconv.ParseFloat(e.Value.(string), 64)
			break
		case types.NUMBER:
			break
		case types.BOOL:
			val := 0.0
			if e.Value.(bool) {
				val = 1.0
			}
			e.Value = val
			break
		}
		e.DataType = types.NUMBER
		break
	case "BOOL":
		switch e.DataType {
		case types.STRING:
			value := strings.ToUpper(e.Value.(string)) != "FALSE" &&
				strings.TrimSpace(e.Value.(string)) != ""
			e.Value = value
			break
		case types.NUMBER:
			e.Value = e.Value.(float64) != 0.0
			break
		case types.BOOL:
			break
		}
		e.DataType = types.BOOL
		break
	}
	r.stack.Push(e)
}

func (r *Runtime) VisitDivideNode(n *ast.DivideNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == types.NUMBER && right.DataType == types.NUMBER {
		r.stack.Push(scope.Entry{
			Value:    left.Value.(float64) / right.Value.(float64),
			DataType: types.NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitEqualNode(n *ast.EqualNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == right.DataType {
		r.stack.Push(scope.Entry{
			Value:    left.Value == right.Value,
			DataType: types.BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitFuncCallNode(n *ast.FuncCallNode) {
	e, ok := r.curScope.GetFunc(n.Name)
	if !ok {
		panic("Function not defined")
	}

	var p params
	for _, arg := range n.Args {
		arg.Accept(r)
		p = append(p, r.stack.Pop().(scope.Entry))
	}

	if e.DataType == types.BUILTIN {
		builtins[n.Name](&r.stack, p)
		return
	}

	f := e.Value.(*ast.FuncDefNode)
	if len(n.Args) != len(f.Args) {
		panic("wrong number of arguments in function call")
	}

	r.scopeStack.Push(r.curScope)
	r.curScope = scope.CleanClone(f.Scope)
	for i, arg := range f.Args {
		r.curScope.SetVar(arg, p[i])
	}
	for _, stmt := range f.Body {
		stmt.Accept(r)
		if r.stack.Size() > 0 {
			break
		}
	}
	r.curScope = r.scopeStack.Pop().(*scope.Scope)
}

func (r *Runtime) VisitFuncDefNode(n *ast.FuncDefNode) {
	// nothing to do
}

func (r *Runtime) VisitGreaterEqualNode(n *ast.GreaterEqualNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == types.NUMBER && right.DataType == types.NUMBER {
		r.stack.Push(scope.Entry{
			Value:    left.Value.(float64) >= right.Value.(float64),
			DataType: types.BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitGreaterNode(n *ast.GreaterNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == types.NUMBER && right.DataType == types.NUMBER {
		r.stack.Push(scope.Entry{
			Value:    left.Value.(float64) > right.Value.(float64),
			DataType: types.BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitIfNode(n *ast.IfNode) {
	n.Cond.Accept(r)
	cond := r.stack.Pop().(scope.Entry)
	if cond.DataType != types.BOOL {
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

func (r *Runtime) VisitLessEqualNode(n *ast.LessEqualNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == types.NUMBER && right.DataType == types.NUMBER {
		r.stack.Push(scope.Entry{
			Value:    left.Value.(float64) <= right.Value.(float64),
			DataType: types.BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitLessNode(n *ast.LessNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == types.NUMBER && right.DataType == types.NUMBER {
		r.stack.Push(scope.Entry{
			Value:    left.Value.(float64) < right.Value.(float64),
			DataType: types.BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitModuloNode(n *ast.ModuloNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == types.NUMBER && right.DataType == types.NUMBER {
		r.stack.Push(scope.Entry{
			Value:    math.Mod(left.Value.(float64), right.Value.(float64)),
			DataType: types.NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitMultiplyNode(n *ast.MultiplyNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == types.NUMBER && right.DataType == types.NUMBER {
		r.stack.Push(scope.Entry{
			Value:    left.Value.(float64) * right.Value.(float64),
			DataType: types.NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNegativeNode(n *ast.NegativeNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(scope.Entry)

	if e.DataType == types.NUMBER {
		r.stack.Push(scope.Entry{
			Value:    -e.Value.(float64),
			DataType: types.NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNotEqualNode(n *ast.NotEqualNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == right.DataType {
		r.stack.Push(scope.Entry{
			Value:    left.Value != right.Value,
			DataType: types.BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNotNode(n *ast.NotNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(scope.Entry)

	if e.DataType == types.BOOL {
		r.stack.Push(scope.Entry{
			Value:    !e.Value.(bool),
			DataType: types.BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitNumberNode(n *ast.NumberNode) {
	r.stack.Push(scope.Entry{
		Value:    n.Value,
		DataType: types.NUMBER,
	})
}

func (r *Runtime) VisitOrNode(n *ast.OrNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)
	// short-circuit if true
	if left.DataType == types.BOOL && left.Value.(bool) {
		r.stack.Push(scope.Entry{
			Value:    true,
			DataType: types.BOOL,
		})
		return
	}

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == types.BOOL && right.DataType == types.BOOL {
		r.stack.Push(scope.Entry{
			Value:    left.Value.(bool) || right.Value.(bool),
			DataType: types.BOOL,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitPositiveNode(n *ast.PositiveNode) {
	n.Term.Accept(r)
	e := r.stack.Pop().(scope.Entry)

	if e.DataType == types.NUMBER {
		r.stack.Push(scope.Entry{
			Value:    math.Abs(e.Value.(float64)),
			DataType: types.NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitProgramNode(n *ast.ProgramNode) {
	n.Scope.Parent = r.curScope
	r.scopeStack.Push(r.curScope)
	r.curScope = n.Scope

	for _, stmt := range n.Stmts {
		stmt.Accept(r)
	}

	r.curScope = r.scopeStack.Pop().(*scope.Scope)
}

func (r *Runtime) VisitReturnNode(n *ast.ReturnNode) {
	n.Expr.Accept(r)
}

func (r *Runtime) VisitStringNode(n *ast.StringNode) {
	r.stack.Push(scope.Entry{
		Value:    n.Value,
		DataType: types.STRING,
	})
}

func (r *Runtime) VisitSubtractNode(n *ast.SubtractNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)

	if left.DataType == types.NUMBER && right.DataType == types.NUMBER {
		r.stack.Push(scope.Entry{
			Value:    left.Value.(float64) - right.Value.(float64),
			DataType: types.NUMBER,
		})
		return
	}
	panic("operation not permitted with type")
}

func (r *Runtime) VisitVariableNode(n *ast.VariableNode) {
	expr, ok := r.curScope.GetVar(n.Name)
	if !ok {
		panic("variable is not defined")
	}
	r.stack.Push(expr)
}

func (r *Runtime) VisitWhileNode(n *ast.WhileNode) {
	for {
		n.Cond.Accept(r)
		cond := r.stack.Pop().(scope.Entry)
		if cond.DataType != types.BOOL {
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
