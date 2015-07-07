// Package runtime provides runtime support for executing Kiwi programs.
package runtime

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/symtable"
	"github.com/tboronczyk/kiwi/token"
	"github.com/tboronczyk/kiwi/util"
)

// Runtime provides an execution context for a program.
type (
	Runtime struct {
		stack     util.Stack
		varTable  *symtable.SymTable
		funcTable *symtable.SymTable
	}

	valueEntry struct {
		value interface{}
		dtype DataType
	}

	params []valueEntry
)

// New returns a new runtime instance.
func New() *Runtime {
	r := &Runtime{
		stack:     util.NewStack(),
		varTable:  symtable.New(),
		funcTable: symtable.New(),
	}
	for n, f := range builtins {
		r.funcTable.Set(n, valueEntry{dtype: BUILTIN, value: f})
	}
	return r
}

// VisitAssignNode evaluates the assignment node n.
func (r *Runtime) VisitAssignNode(n *ast.AssignNode) {
	n.Expr.Accept(r)
	v := r.stack.Pop().(valueEntry)
	// preserve datatype if the variable is already set
	e, ok := r.varTable.Get(n.Name)
	if ok {
		if e.(valueEntry).dtype != v.dtype {
			panic("value type does not match variable type")
		}
	}
	r.varTable.Set(n.Name, v)
}

// VisitBinaryNode evaluates the binary operator expression node n.
func (r *Runtime) VisitBinaryOpNode(n *ast.BinaryOpNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(valueEntry)
	// short-circuit logic operators
	if left.dtype == BOOL {
		if n.Op == token.OR && left.value.(bool) {
			r.stack.Push(valueEntry{value: true, dtype: BOOL})
			return
		}
		if n.Op == token.AND && !left.value.(bool) {
			r.stack.Push(valueEntry{value: false, dtype: BOOL})
			return
		}
	}

	n.Right.Accept(r)
	right := r.stack.Pop().(valueEntry)
	if left.dtype != right.dtype {
		panic("mis-matched types")
	}

	switch left.dtype {
	case NUMBER:
		switch n.Op {
		case token.ADD:
			r.stack.Push(valueEntry{
				value: left.value.(float64) + right.value.(float64),
				dtype: NUMBER,
			})
			return
		case token.SUBTRACT:
			r.stack.Push(valueEntry{
				value: left.value.(float64) - right.value.(float64),
				dtype: NUMBER,
			})
			return
		case token.MULTIPLY:
			r.stack.Push(valueEntry{
				value: left.value.(float64) * right.value.(float64),
				dtype: NUMBER,
			})
			return
		case token.DIVIDE:
			r.stack.Push(valueEntry{
				value: left.value.(float64) / right.value.(float64),
				dtype: NUMBER,
			})
			return
		case token.MODULO:
			r.stack.Push(valueEntry{
				value: math.Mod(left.value.(float64), right.value.(float64)),
				dtype: NUMBER,
			})
			return
		case token.EQUAL:
			r.stack.Push(valueEntry{
				value: left.value.(float64) == right.value.(float64),
				dtype: BOOL,
			})
			return
		case token.NOT_EQUAL:
			r.stack.Push(valueEntry{
				value: left.value.(float64) != right.value.(float64),
				dtype: BOOL,
			})
			return
		case token.LESS:
			r.stack.Push(valueEntry{
				value: left.value.(float64) < right.value.(float64),
				dtype: BOOL,
			})
			return
		case token.LESS_EQ:
			r.stack.Push(valueEntry{
				value: left.value.(float64) <= right.value.(float64),
				dtype: BOOL,
			})
			return
		case token.GREATER:
			r.stack.Push(valueEntry{
				value: left.value.(float64) > right.value.(float64),
				dtype: BOOL,
			})
			return
		case token.GREATER_EQ:
			r.stack.Push(valueEntry{
				value: left.value.(float64) >= right.value.(float64),
				dtype: BOOL,
			})
			return
		}
		break
	case STRING:
		switch n.Op {
		case token.ADD:
			r.stack.Push(valueEntry{
				value: left.value.(string) + right.value.(string),
				dtype: STRING,
			})
			return
		case token.EQUAL:
			r.stack.Push(valueEntry{
				value: left.value.(string) == right.value.(string),
				dtype: BOOL,
			})
			return
		case token.NOT_EQUAL:
			r.stack.Push(valueEntry{
				value: left.value.(string) != right.value.(string),
				dtype: BOOL,
			})
			return
		}
		break
	case BOOL:
		switch n.Op {
		case token.AND:
			r.stack.Push(valueEntry{
				value: left.value.(bool) && right.value.(bool),
				dtype: BOOL,
			})
			return
		case token.OR:
			r.stack.Push(valueEntry{
				value: left.value.(bool) || right.value.(bool),
				dtype: BOOL,
			})
			return
		case token.EQUAL:
			r.stack.Push(valueEntry{
				value: left.value.(bool) == right.value.(bool),
				dtype: BOOL,
			})
			return
		case token.NOT_EQUAL:
			r.stack.Push(valueEntry{
				value: left.value.(bool) != right.value.(bool),
				dtype: BOOL,
			})
			return
		}
		break
	}
	panic("operation not permitted on type")
}

// VisitCastNode evaluates the cast node n.
func (r *Runtime) VisitCastNode(n *ast.CastNode) {
	n.Expr.Accept(r)
	expr := r.stack.Pop().(valueEntry)
	switch strings.ToUpper(n.Cast) {
	case "STR":
		switch expr.dtype {
		case STRING:
			break
		case NUMBER:
			val := fmt.Sprintf("%f", expr.value.(float64))
			val = strings.TrimRight(val, "0")
			val = strings.TrimRight(val, ".")
			expr.value = val
			break
		case BOOL:
			expr.value = strconv.FormatBool(expr.value.(bool))
			break
		}
		expr.dtype = STRING
		break
	case "NUM":
		switch expr.dtype {
		case STRING:
			expr.value, _ = strconv.ParseFloat(expr.value.(string), 64)
			break
		case NUMBER:
			break
		case BOOL:
			val := 0.0
			if expr.value.(bool) {
				val = 1.0
			}
			expr.value = val
			break
		}
		expr.dtype = NUMBER
		break
	case "BOOL":
		switch expr.dtype {
		case STRING:
			value := strings.ToUpper(expr.value.(string)) != "FALSE" &&
				strings.TrimSpace(expr.value.(string)) != ""
			expr.value = value
			break
		case NUMBER:
			expr.value = expr.value.(float64) != 0.0
			break
		case BOOL:
			break
		}
		expr.dtype = BOOL
		break
	}
	r.stack.Push(expr)
}

// VisitFuncCallNode evaluates the function call node n.
func (r *Runtime) VisitFuncCallNode(n *ast.FuncCallNode) {

	e, ok := r.funcTable.Get(n.Name)
	if !ok {
		panic("Function not defined")
	}

	var p params
	for _, arg := range n.Args {
		arg.Accept(r)
		p = append(p, r.stack.Pop().(valueEntry))
	}

	if e.(valueEntry).dtype == BUILTIN {
		builtins[n.Name](&r.stack, p)
		return
	}

	f := e.(valueEntry).value.(*ast.FuncDefNode)
	if len(n.Args) != len(f.Args) {
		panic("wrong number of arguments in function call")
	}

	r.varTable = symtable.NewScope(r.varTable)
	for i, arg := range f.Args {
		r.varTable.Set(arg, p[i])
	}
	r.funcTable = symtable.NewScope(r.funcTable)
	for _, stmt := range f.Body {
		stmt.Accept(r)
		if r.stack.Size() > 0 {
			break
		}
	}
	r.varTable = r.varTable.Parent()
	r.funcTable = r.funcTable.Parent()
}

// VisitFuncDefNode evaluates the function definition node n.
func (r *Runtime) VisitFuncDefNode(n *ast.FuncDefNode) {
	_, ok := r.funcTable.Get(n.Name)
	if ok {
		panic("function is already defined")
	}
	r.funcTable.Set(n.Name, valueEntry{dtype: FUNC, value: n})
}

// VisitIfNode evaluates the if construct node n.
func (r *Runtime) VisitIfNode(n *ast.IfNode) {
	n.Condition.Accept(r)
	cond := r.stack.Pop().(valueEntry)
	if cond.dtype != BOOL {
		panic("non-bool expression used as condition")
	}
	if cond.value.(bool) {
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
	n.Expr.Accept(r)
	expr := r.stack.Pop().(valueEntry)

	switch expr.dtype {
	case NUMBER:
		switch n.Op {
		case token.ADD:
			expr.value = math.Abs(expr.value.(float64))
			r.stack.Push(expr)
			return
		case token.SUBTRACT:
			expr.value = 0.0 - expr.value.(float64)
			r.stack.Push(expr)
			return
		}
	case BOOL:
		switch n.Op {
		case token.NOT:
			expr.value = !expr.value.(bool)
			r.stack.Push(expr)
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
		r.stack.Push(valueEntry{
			value: value,
			dtype: NUMBER,
		})
		break
	case token.BOOL:
		r.stack.Push(valueEntry{
			value: strings.ToUpper(n.Value) == "TRUE",
			dtype: BOOL,
		})
		break
	case token.STRING:
		r.stack.Push(valueEntry{
			value: n.Value,
			dtype: STRING,
		})
		break
	default:
		panic("literal is of unknown type")
	}
}

// VisitVariableNode evaluates the variable expression node n.
func (r *Runtime) VisitVariableNode(n *ast.VariableNode) {
	expr, ok := r.varTable.Get(n.Name)
	if !ok {
		panic("variable is not defined")
	}
	r.stack.Push(expr)
}

// VisitWhileNode evaluates the while construct node n.
func (r *Runtime) VisitWhileNode(n *ast.WhileNode) {
	for {
		n.Condition.Accept(r)
		cond := r.stack.Pop().(valueEntry)
		if cond.dtype != BOOL {
			panic("non-bool expression used as condition")
		}
		if !cond.value.(bool) {
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
