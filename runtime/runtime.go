// Package runtime provides runtime support for executing Kiwi programs.
package runtime

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scope"
	"github.com/tboronczyk/kiwi/token"
	"github.com/tboronczyk/kiwi/types"
	"github.com/tboronczyk/kiwi/util"
)

// Runtime provides an execution context for a program.
type (
	Runtime struct {
		stack      util.Stack
		curScope   *scope.Scope
		scopeStack util.Stack
	}

	params []scope.Entry
)

// New returns a new runtime instance.
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

func (r *Runtime) VisitProgramNode(n *ast.ProgramNode) {
	n.Scope.Parent = r.curScope
	r.scopeStack.Push(r.curScope)
	r.curScope = n.Scope

	for _, stmt := range n.Stmts {
		stmt.Accept(r)
	}

	r.curScope = r.scopeStack.Pop().(*scope.Scope)
}

// VisitAssignNode evaluates the assignment node n.
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

// VisitBinaryNode evaluates the binary operator expression node n.
func (r *Runtime) VisitBinOpNode(n *ast.BinOpNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(scope.Entry)
	// short-circuit logic operators
	if left.DataType == types.BOOL {
		if n.Op == token.OR && left.Value.(bool) {
			r.stack.Push(scope.Entry{Value: true, DataType: types.BOOL})
			return
		}
		if n.Op == token.AND && !left.Value.(bool) {
			r.stack.Push(scope.Entry{Value: false, DataType: types.BOOL})
			return
		}
	}

	n.Right.Accept(r)
	right := r.stack.Pop().(scope.Entry)
	if left.DataType != right.DataType {
		panic("mis-matched types")
	}

	switch left.DataType {
	case types.NUMBER:
		switch n.Op {
		case token.ADD:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(float64) + right.Value.(float64),
				DataType: types.NUMBER,
			})
			return
		case token.SUBTRACT:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(float64) - right.Value.(float64),
				DataType: types.NUMBER,
			})
			return
		case token.MULTIPLY:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(float64) * right.Value.(float64),
				DataType: types.NUMBER,
			})
			return
		case token.DIVIDE:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(float64) / right.Value.(float64),
				DataType: types.NUMBER,
			})
			return
		case token.MODULO:
			r.stack.Push(scope.Entry{
				Value:    math.Mod(left.Value.(float64), right.Value.(float64)),
				DataType: types.NUMBER,
			})
			return
		case token.EQUAL:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(float64) == right.Value.(float64),
				DataType: types.BOOL,
			})
			return
		case token.NOT_EQUAL:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(float64) != right.Value.(float64),
				DataType: types.BOOL,
			})
			return
		case token.LESS:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(float64) < right.Value.(float64),
				DataType: types.BOOL,
			})
			return
		case token.LESS_EQ:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(float64) <= right.Value.(float64),
				DataType: types.BOOL,
			})
			return
		case token.GREATER:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(float64) > right.Value.(float64),
				DataType: types.BOOL,
			})
			return
		case token.GREATER_EQ:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(float64) >= right.Value.(float64),
				DataType: types.BOOL,
			})
			return
		}
		break
	case types.STRING:
		switch n.Op {
		case token.ADD:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(string) + right.Value.(string),
				DataType: types.STRING,
			})
			return
		case token.EQUAL:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(string) == right.Value.(string),
				DataType: types.BOOL,
			})
			return
		case token.NOT_EQUAL:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(string) != right.Value.(string),
				DataType: types.BOOL,
			})
			return
		}
		break
	case types.BOOL:
		switch n.Op {
		case token.AND:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(bool) && right.Value.(bool),
				DataType: types.BOOL,
			})
			return
		case token.OR:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(bool) || right.Value.(bool),
				DataType: types.BOOL,
			})
			return
		case token.EQUAL:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(bool) == right.Value.(bool),
				DataType: types.BOOL,
			})
			return
		case token.NOT_EQUAL:
			r.stack.Push(scope.Entry{
				Value:    left.Value.(bool) != right.Value.(bool),
				DataType: types.BOOL,
			})
			return
		}
		break
	}
	panic("operation not permitted on type")
}

// VisitCastNode evaluates the cast node n.
func (r *Runtime) VisitCastNode(n *ast.CastNode) {
	n.Term.Accept(r)
	term := r.stack.Pop().(scope.Entry)
	switch strings.ToUpper(n.Cast) {
	case "STR":
		switch term.DataType {
		case types.STRING:
			break
		case types.NUMBER:
			val := fmt.Sprintf("%f", term.Value.(float64))
			val = strings.TrimRight(val, "0")
			val = strings.TrimRight(val, ".")
			term.Value = val
			break
		case types.BOOL:
			term.Value = strconv.FormatBool(term.Value.(bool))
			break
		}
		term.DataType = types.STRING
		break
	case "NUM":
		switch term.DataType {
		case types.STRING:
			term.Value, _ = strconv.ParseFloat(term.Value.(string), 64)
			break
		case types.NUMBER:
			break
		case types.BOOL:
			val := 0.0
			if term.Value.(bool) {
				val = 1.0
			}
			term.Value = val
			break
		}
		term.DataType = types.NUMBER
		break
	case "BOOL":
		switch term.DataType {
		case types.STRING:
			value := strings.ToUpper(term.Value.(string)) != "FALSE" &&
				strings.TrimSpace(term.Value.(string)) != ""
			term.Value = value
			break
		case types.NUMBER:
			term.Value = term.Value.(float64) != 0.0
			break
		case types.BOOL:
			break
		}
		term.DataType = types.BOOL
		break
	}
	r.stack.Push(term)
}

// VisitFuncCallNode evaluates the function call node n.
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

// VisitFuncDefNode evaluates the function definition node n.
func (r *Runtime) VisitFuncDefNode(n *ast.FuncDefNode) {
	// nothing to do
}

// VisitIfNode evaluates the if construct node n.
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
		n.Else.Accept(r)
	}
}

// VisitReturnNode evaluates the reutnr statment node n.
func (r *Runtime) VisitReturnNode(n *ast.ReturnNode) {
	n.Expr.Accept(r)
}

// VisitUnaryOpNode evaluates the unary operator expression node n.
func (r *Runtime) VisitUnaryOpNode(n *ast.UnaryOpNode) {
	n.Term.Accept(r)
	term := r.stack.Pop().(scope.Entry)

	switch term.DataType {
	case types.NUMBER:
		switch n.Op {
		case token.ADD:
			term.Value = math.Abs(term.Value.(float64))
			r.stack.Push(term)
			return
		case token.SUBTRACT:
			term.Value = 0.0 - term.Value.(float64)
			r.stack.Push(term)
			return
		}
	case types.BOOL:
		switch n.Op {
		case token.NOT:
			term.Value = !term.Value.(bool)
			r.stack.Push(term)
			return
		}
	}
	panic("invalid cast")
}

// VisitValueNode evaluates the value expression node n.
func (r *Runtime) VisitValueNode(n *ast.ValueNode) {
	switch n.Type {
	case token.NUMBER:
		value, _ := strconv.ParseFloat(n.Value, 64)
		r.stack.Push(scope.Entry{
			Value:    value,
			DataType: types.NUMBER,
		})
		break
	case token.BOOL:
		r.stack.Push(scope.Entry{
			Value:    strings.ToUpper(n.Value) == "TRUE",
			DataType: types.BOOL,
		})
		break
	case token.STRING:
		r.stack.Push(scope.Entry{
			Value:    n.Value,
			DataType: types.STRING,
		})
		break
	default:
		panic("literal is of unknown type")
	}
}

// VisitVariableNode evaluates the variable expression node n.
func (r *Runtime) VisitVariableNode(n *ast.VariableNode) {
	expr, ok := r.curScope.GetVar(n.Name)
	if !ok {
		panic("variable is not defined")
	}
	r.stack.Push(expr)
}

// VisitWhileNode evaluates the while construct node n.
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
