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

type (
	Runtime struct {
		stack    util.Stack
		symTable *symtable.SymTable
	}

	ValueEntry struct {
		Value interface{}
		Type  DataType
	}

	Params []ValueEntry
)

func New() *Runtime {
	r := &Runtime{
		stack:    util.NewStack(),
		symTable: symtable.New(),
	}
	for n, f := range builtins {
		r.symTable.Set(n, symtable.FUNC,
			ValueEntry{Type: BUILTIN, Value: f})
	}
	return r
}

func (r *Runtime) VisitAssignNode(n *ast.AssignNode) {
	n.Expr.Accept(r)
	r.symTable.Set(n.Name, symtable.VAR, r.stack.Pop())
	n.SymTable = r.symTable
}

func (r *Runtime) VisitBinaryOpNode(n *ast.BinaryOpNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(ValueEntry)
	// short-circuit logic operators
	if left.Type == BOOL {
		if n.Op == token.OR && left.Value.(bool) {
			r.stack.Push(ValueEntry{Value: true, Type: BOOL})
			return
		}
		if n.Op == token.AND && !left.Value.(bool) {
			r.stack.Push(ValueEntry{Value: false, Type: BOOL})
			return
		}
	}

	n.Right.Accept(r)
	right := r.stack.Pop().(ValueEntry)
	if left.Type != right.Type {
		panic("mis-matched types")
	}

	switch left.Type {
	case NUMBER:
		switch n.Op {
		case token.ADD:
			r.stack.Push(ValueEntry{
				Value: left.Value.(float64) + right.Value.(float64),
				Type:  NUMBER,
			})
			return
		case token.SUBTRACT:
			r.stack.Push(ValueEntry{
				Value: left.Value.(float64) - right.Value.(float64),
				Type:  NUMBER,
			})
			return
		case token.MULTIPLY:
			r.stack.Push(ValueEntry{
				Value: left.Value.(float64) * right.Value.(float64),
				Type:  NUMBER,
			})
			return
		case token.DIVIDE:
			r.stack.Push(ValueEntry{
				Value: left.Value.(float64) / right.Value.(float64),
				Type:  NUMBER,
			})
			return
		case token.MODULO:
			r.stack.Push(ValueEntry{
				Value: math.Mod(left.Value.(float64), right.Value.(float64)),
				Type:  NUMBER,
			})
			return
		case token.EQUAL:
			r.stack.Push(ValueEntry{
				Value: left.Value.(float64) == right.Value.(float64),
				Type:  BOOL,
			})
			return
		case token.NOT_EQUAL:
			r.stack.Push(ValueEntry{
				Value: left.Value.(float64) != right.Value.(float64),
				Type:  BOOL,
			})
			return
		case token.LESS:
			r.stack.Push(ValueEntry{
				Value: left.Value.(float64) < right.Value.(float64),
				Type:  BOOL,
			})
			return
		case token.LESS_EQ:
			r.stack.Push(ValueEntry{
				Value: left.Value.(float64) <= right.Value.(float64),
				Type:  BOOL,
			})
			return
		case token.GREATER:
			r.stack.Push(ValueEntry{
				Value: left.Value.(float64) > right.Value.(float64),
				Type:  BOOL,
			})
			return
		case token.GREATER_EQ:
			r.stack.Push(ValueEntry{
				Value: left.Value.(float64) >= right.Value.(float64),
				Type:  BOOL,
			})
			return
		}
		break
	case STRING:
		switch n.Op {
		case token.ADD:
			r.stack.Push(ValueEntry{
				Value: left.Value.(string) + right.Value.(string),
				Type:  STRING,
			})
			return
		case token.EQUAL:
			r.stack.Push(ValueEntry{
				Value: left.Value.(string) == right.Value.(string),
				Type:  BOOL,
			})
			return
		case token.NOT_EQUAL:
			r.stack.Push(ValueEntry{
				Value: left.Value.(string) != right.Value.(string),
				Type:  BOOL,
			})
			return
		}
		break
	case BOOL:
		switch n.Op {
		case token.AND:
			r.stack.Push(ValueEntry{
				Value: left.Value.(bool) && right.Value.(bool),
				Type:  BOOL,
			})
			return
		case token.OR:
			r.stack.Push(ValueEntry{
				Value: left.Value.(bool) || right.Value.(bool),
				Type:  BOOL,
			})
			return
		case token.EQUAL:
			r.stack.Push(ValueEntry{
				Value: left.Value.(bool) == right.Value.(bool),
				Type:  BOOL,
			})
			return
		case token.NOT_EQUAL:
			r.stack.Push(ValueEntry{
				Value: left.Value.(bool) != right.Value.(bool),
				Type:  BOOL,
			})
			return
		}
		break
	}
	panic("operation not permitted on type")
}

func (r *Runtime) VisitCastNode(n *ast.CastNode) {
	n.Expr.Accept(r)
	expr := r.stack.Pop().(ValueEntry)
	switch strings.ToUpper(n.Cast) {
	case "STR":
		switch expr.Type {
		case STRING:
			break
		case NUMBER:
			val := fmt.Sprintf("%f", expr.Value.(float64))
			val = strings.TrimRight(val, "0")
			val = strings.TrimRight(val, ".")
			expr.Value = val
			break
		case BOOL:
			expr.Value = strconv.FormatBool(expr.Value.(bool))
			break
		}
		expr.Type = STRING
		break
	case "NUM":
		switch expr.Type {
		case STRING:
			expr.Value, _ = strconv.ParseFloat(expr.Value.(string), 64)
			break
		case NUMBER:
			break
		case BOOL:
			val := 0.0
			if expr.Value.(bool) {
				val = 1.0
			}
			expr.Value = val
			break
		}
		expr.Type = NUMBER
		break
	case "BOOL":
		switch expr.Type {
		case STRING:
			value := strings.ToUpper(expr.Value.(string)) != "FALSE" &&
				strings.TrimSpace(expr.Value.(string)) != ""
			expr.Value = value
			break
		case NUMBER:
			expr.Value = expr.Value.(float64) != 0.0
			break
		case BOOL:
			break
		}
		expr.Type = BOOL
		break
	}
	r.stack.Push(expr)
}

func (r *Runtime) VisitFuncCallNode(n *ast.FuncCallNode) {

	e, ok := r.symTable.Get(n.Name, symtable.FUNC)
	if !ok {
		panic("Function not defined")
	}

	var p Params
	for _, arg := range n.Args {
		arg.Accept(r)
		p = append(p, r.stack.Pop().(ValueEntry))
	}

	if e.(ValueEntry).Type == BUILTIN {
		builtins[n.Name](&r.stack, p)
		return
	}

	f := e.(ValueEntry).Value.(*ast.FuncDefNode)
	if len(n.Args) != len(f.Args) {
		panic("wrong number of arguments in function call")
	}

	r.symTable = symtable.NewScope(r.symTable)
	for i, arg := range f.Args {
		r.symTable.Set(arg, symtable.VAR, p[i])
	}
	for _, stmt := range f.Body {
		stmt.Accept(r)
		if r.stack.Size() > 0 {
			break
		}
	}
	r.symTable = r.symTable.Parent()
}

func (r *Runtime) VisitFuncDefNode(n *ast.FuncDefNode) {
	_, ok := r.symTable.Get(n.Name, symtable.FUNC)
	if ok {
		panic("function is already defined")
	}
	n.SymTable = r.symTable
	r.symTable.Set(n.Name, symtable.FUNC, ValueEntry{Type: FUNC, Value: n})
}

func (r *Runtime) VisitIfNode(n *ast.IfNode) {
	n.Condition.Accept(r)
	cond := r.stack.Pop().(ValueEntry)
	if cond.Type != BOOL {
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

func (r *Runtime) VisitReturnNode(n *ast.ReturnNode) {
	n.Expr.Accept(r)
}

func (r *Runtime) VisitUnaryOpNode(n *ast.UnaryOpNode) {
	n.Expr.Accept(r)
	expr := r.stack.Pop().(ValueEntry)

	switch expr.Type {
	case NUMBER:
		switch n.Op {
		case token.ADD:
			expr.Value = math.Abs(expr.Value.(float64))
			r.stack.Push(expr)
			return
		case token.SUBTRACT:
			expr.Value = 0.0 - expr.Value.(float64)
			r.stack.Push(expr)
			return
		}
	case BOOL:
		switch n.Op {
		case token.NOT:
			expr.Value = !expr.Value.(bool)
			r.stack.Push(expr)
			return
		}
	}
	panic("invalid cast")
}

func (r *Runtime) VisitValueNode(n *ast.ValueNode) {
	switch n.Type {
	case token.NUMBER:
		value, _ := strconv.ParseFloat(n.Value, 64)
		r.stack.Push(ValueEntry{
			Value: value,
			Type:  NUMBER,
		})
		break
	case token.BOOL:
		r.stack.Push(ValueEntry{
			Value: strings.ToUpper(n.Value) == "TRUE",
			Type:  BOOL,
		})
		break
	case token.STRING:
		r.stack.Push(ValueEntry{
			Value: n.Value,
			Type:  STRING,
		})
		break
	default:
		panic("literal is of unknown type")
	}
}

func (r *Runtime) VisitVariableNode(n *ast.VariableNode) {
	expr, ok := r.symTable.Get(n.Name, symtable.VAR)
	if !ok {
		panic("variable is not defined")
	}
	r.stack.Push(expr)
	n.SymTable = r.symTable
}

func (r *Runtime) VisitWhileNode(n *ast.WhileNode) {
	for {
		n.Condition.Accept(r)
		cond := r.stack.Pop().(ValueEntry)
		if cond.Type != BOOL {
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
