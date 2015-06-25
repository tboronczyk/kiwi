package runtime

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/tboronczyk/kiwi/analyzer"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/token"
	"github.com/tboronczyk/kiwi/util"
)

type (
	Runtime struct {
		stack util.Stack
	}

	stackEntry struct {
		Value interface{}
		Type  analyzer.DataType
	}
)

func New() *Runtime {
	return &Runtime{
		stack: util.NewStack(),
	}
}

func (r *Runtime) VisitAssignNode(n *ast.AssignNode) {
	n.Expr.Accept(r)
	n.SymTable.Set(n.Name, r.stack.Pop())
}

func (r *Runtime) VisitBinaryOpNode(n *ast.BinaryOpNode) {
	n.Left.Accept(r)
	left := r.stack.Pop().(stackEntry)

	// short-circuit logic operators
	if left.Type == analyzer.BOOL {
		if n.Op == token.OR && left.Value.(bool) {
			r.stack.Push(stackEntry{Value: true, Type: analyzer.BOOL})
			return
		}
		if n.Op == token.AND && !left.Value.(bool) {
			r.stack.Push(stackEntry{Value: false, Type: analyzer.BOOL})
			return
		}
	}

	n.Right.Accept(r)
	right := r.stack.Pop().(stackEntry)

	e := stackEntry{}
	switch left.Type {
	case analyzer.NUMBER:
		switch n.Op {
		case token.ADD:
			e.Value = left.Value.(float64) + right.Value.(float64)
			e.Type = analyzer.NUMBER
			break
		case token.SUBTRACT:
			e.Value = left.Value.(float64) - right.Value.(float64)
			e.Type = analyzer.NUMBER
			break
		case token.MULTIPLY:
			e.Value = left.Value.(float64) * right.Value.(float64)
			e.Type = analyzer.NUMBER
			break
		case token.DIVIDE:
			e.Value = left.Value.(float64) / right.Value.(float64)
			e.Type = analyzer.NUMBER
			break
		case token.MODULO:
			e.Value = math.Mod(left.Value.(float64), right.Value.(float64))
			e.Type = analyzer.NUMBER
			break
		case token.EQUAL:
			e.Value = left.Value.(float64) == right.Value.(float64)
			e.Type = analyzer.BOOL
			break
		case token.NOT_EQUAL:
			e.Value = left.Value.(float64) != right.Value.(float64)
			e.Type = analyzer.BOOL
			break
		case token.LESS:
			e.Value = left.Value.(float64) < right.Value.(float64)
			e.Type = analyzer.BOOL
			break
		case token.LESS_EQ:
			e.Value = left.Value.(float64) <= right.Value.(float64)
			e.Type = analyzer.BOOL
			break
		case token.GREATER:
			e.Value = left.Value.(float64) > right.Value.(float64)
			e.Type = analyzer.BOOL
			break
		case token.GREATER_EQ:
			e.Value = left.Value.(float64) >= right.Value.(float64)
			e.Type = analyzer.BOOL
			break
		}
		break
	case analyzer.STRING:
		switch n.Op {
		case token.ADD:
			e.Value = left.Value.(string) + right.Value.(string)
			e.Type = analyzer.STRING
			break
		case token.EQUAL:
			e.Value = left.Value.(string) == right.Value.(string)
			e.Type = analyzer.BOOL
			break
		case token.NOT_EQUAL:
			e.Value = left.Value.(string) != right.Value.(string)
			e.Type = analyzer.BOOL
		}
		break
	case analyzer.BOOL:
		switch n.Op {
		case token.AND:
			e.Value = left.Value.(bool) && right.Value.(bool)
			e.Type = analyzer.BOOL
			break
		case token.OR:
			e.Value = left.Value.(bool) || right.Value.(bool)
			e.Type = analyzer.BOOL
			break
		case token.EQUAL:
			e.Value = left.Value.(bool) == right.Value.(bool)
			e.Type = analyzer.BOOL
			break
		case token.NOT_EQUAL:
			e.Value = left.Value.(bool) != right.Value.(bool)
			e.Type = analyzer.BOOL
			break
		}
		break
	}
	r.stack.Push(e)
}

func (r *Runtime) VisitCastNode(n *ast.CastNode) {
	n.Expr.Accept(r)
	expr := r.stack.Pop().(stackEntry)
	switch strings.ToUpper(n.Cast) {
	case "STRING":
		switch expr.Type {
		case analyzer.STRING:
			break
		case analyzer.NUMBER:
			val := fmt.Sprintf("%f", expr.Value.(float64))
			val = strings.TrimRight(val, "0")
			val = strings.TrimRight(val, ".")
			expr.Value = val
			break
		case analyzer.BOOL:
			expr.Value = strconv.FormatBool(expr.Value.(bool))
			break
		}
		expr.Type = analyzer.STRING
		break
	case "NUMBER":
		switch expr.Type {
		case analyzer.STRING:
			expr.Value, _ = strconv.ParseFloat(expr.Value.(string), 64)
			break
		case analyzer.NUMBER:
			break
		case analyzer.BOOL:
			val := 0.0
			if expr.Value.(bool) {
				val = 1.0
			}
			expr.Value = val
			break
		}
		expr.Type = analyzer.NUMBER
		break
	case "BOOL":
		switch expr.Type {
		case analyzer.STRING:
			value := strings.ToUpper(expr.Value.(string)) != "FALSE" &&
				strings.TrimSpace(expr.Value.(string)) != ""
			expr.Value = value
			break
		case analyzer.NUMBER:
			expr.Value = expr.Value.(float64) != 0.0
			break
		case analyzer.BOOL:
			break
		}
		expr.Type = analyzer.BOOL
		break
	}
	r.stack.Push(expr)
}

func (r *Runtime) VisitFuncCallNode(n *ast.FuncCallNode) {
}

func (r *Runtime) VisitFuncDefNode(n *ast.FuncDefNode) {
}

func (r *Runtime) VisitIfNode(n *ast.IfNode) {
	n.Condition.Accept(r)
	cond := r.stack.Pop().(stackEntry)
	if cond.Value.(bool) {
		for _, stmt := range n.Body {
			stmt.Accept(r)
		}
	} else if n.Else != nil {
		n.Else.Accept(r)
	}
}

func (r *Runtime) VisitReturnNode(n *ast.ReturnNode) {
}

func (r *Runtime) VisitUnaryOpNode(n *ast.UnaryOpNode) {
	n.Expr.Accept(r)
	expr := r.stack.Pop().(stackEntry)
	switch expr.Type {
	case analyzer.NUMBER:
		switch n.Op {
		case token.ADD:
			expr.Value = math.Abs(expr.Value.(float64))
			break
		case token.SUBTRACT:
			expr.Value = 0.0 - expr.Value.(float64)
			break
		}
		break
	case analyzer.BOOL:
		switch n.Op {
		case token.NOT:
			expr.Value = !expr.Value.(bool)
			break
		}
		break
	}
	r.stack.Push(expr)
}

func (r *Runtime) VisitValueNode(n *ast.ValueNode) {
	switch n.Type {
	case token.NUMBER:
		value, _ := strconv.ParseFloat(n.Value, 64)
		r.stack.Push(stackEntry{
			Value: value,
			Type:  analyzer.NUMBER,
		})
		break
	case token.BOOL:
		r.stack.Push(stackEntry{
			Value: strings.ToUpper(n.Value) == "TRUE",
			Type:  analyzer.BOOL,
		})
		break
	case token.STRING:
		r.stack.Push(stackEntry{
			Value: n.Value,
			Type:  analyzer.STRING,
		})
		break
	}
}

func (r *Runtime) VisitVariableNode(n *ast.VariableNode) {
	value, _ := n.SymTable.Get(n.Name)
	r.stack.Push(value)
}

func (r *Runtime) VisitWhileNode(n *ast.WhileNode) {
	for {
		n.Condition.Accept(r)
		cond := r.stack.Pop().(stackEntry)
		if !cond.Value.(bool) {
			return
		}
		for _, stmt := range n.Body {
			stmt.Accept(r)
		}
	}
}
