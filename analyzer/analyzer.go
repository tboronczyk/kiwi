package analyzer

import (
	"strings"

	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/symtable"
	"github.com/tboronczyk/kiwi/token"
	"github.com/tboronczyk/kiwi/util"
)

type (
	DataType uint8

	Analyzer struct {
		stack    util.Stack
		symtable *symtable.SymTable
	}
)

const (
	UNKNOWN DataType = iota
	ANY
	BOOL
	FUNCTION
	NUMBER
	STRING
)

func New() *Analyzer {
	return &Analyzer{
		stack:    util.NewStack(),
		symtable: symtable.New(),
	}
}

func (a *Analyzer) VisitAssignNode(n *ast.AssignNode) {
	n.Expr.Accept(a)
	expr := a.stack.Pop()
	a.symtable.Set(n.Name, symtable.VARIABLE, expr)
}

func (a *Analyzer) VisitBinaryOpNode(n *ast.BinaryOpNode) {
	n.Left.Accept(a)
	left := a.stack.Pop()
	n.Right.Accept(a)
	right := a.stack.Pop()
	if left != ANY && right != ANY && left != right {
		panic("Type mis-match")
	}
	a.stack.Push(left)
}

func (a *Analyzer) VisitCastNode(n *ast.CastNode) {
	t := UNKNOWN
	switch strings.ToUpper(n.Cast) {
	case "BOOL":
		t = BOOL
		break
	case "NUMBER":
		t = NUMBER
		break
	case "STRING":
		t = STRING
		break
	}
	a.stack.Push(t)
}

func (a *Analyzer) VisitFuncCallNode(n *ast.FuncCallNode) {
	dtype, ok := a.symtable.Get(n.Name, symtable.FUNCTION)
	if !ok {
		panic("Function not defined")
	}
	a.stack.Push(dtype)

}

func (a *Analyzer) VisitFuncDefNode(n *ast.FuncDefNode) {
	a.symtable.Set(n.Name, symtable.FUNCTION, UNKNOWN)
	a.symtable.Enter()
	for _, arg := range n.Args {
		a.symtable.Set(arg, symtable.VARIABLE, ANY)
	}
	for _, stmt := range n.Body {
		stmt.Accept(a)
	}
	n.Scope = a.symtable.Current()
	a.symtable.Leave()
}

func (a *Analyzer) VisitIfNode(n *ast.IfNode) {
	n.Condition.Accept(a)
	cond := a.stack.Pop()
	if cond != BOOL {
		panic("Non-bool expression used as condition")
	}
	for _, stmt := range n.Body {
		stmt.Accept(a)
	}
	if n.Else != nil {
		n.Else.Accept(a)
	}
}

func (a *Analyzer) VisitReturnNode(n *ast.ReturnNode) {
}

func (a *Analyzer) VisitUnaryOpNode(n *ast.UnaryOpNode) {
	n.Right.Accept(a)
	// right := a.stack.Pop()
	// a.stack.Push(right)
}

func (a *Analyzer) VisitValueNode(n *ast.ValueNode) {
	t := UNKNOWN
	switch n.Type {
	case token.NUMBER:
		t = NUMBER
		break
	case token.STRING:
		t = STRING
		break
	case token.BOOL:
		t = BOOL
		break
	default:
		panic("Unknown value type")
	}
	a.stack.Push(t)
}

func (a *Analyzer) VisitVariableNode(n *ast.VariableNode) {
	dtype, ok := a.symtable.Get(n.Name, symtable.VARIABLE)
	if !ok {
		panic("Variable not defined")
	}
	a.stack.Push(dtype)
}

func (a *Analyzer) VisitWhileNode(n *ast.WhileNode) {
	n.Condition.Accept(a)
	cond := a.stack.Pop()
	if cond != BOOL {
		panic("Non-bool expression used as condition")
	}
	for _, stmt := range n.Body {
		stmt.Accept(a)
	}
}
