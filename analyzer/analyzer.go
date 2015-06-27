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
		stack      util.Stack
		symTable   *symtable.SymTable
		funcStack  util.Stack
	}

	StackEntry struct {
		Value interface{}
		Type  DataType
	}
)

const (
	UNKNOWN DataType = iota
	BOOL
	FUNC
	NUMBER
	STRING
)

func New() *Analyzer {
	a := &Analyzer{
		stack:     util.NewStack(),
		symTable:  symtable.New(),
		funcStack: util.NewStack(),
	}

	return a
}

func (a *Analyzer) VisitAssignNode(n *ast.AssignNode) {
	n.Expr.Accept(a)
	expr := a.stack.Pop()
	a.symTable.Set(n.Name, symtable.VARTYPE, expr)
	n.SymTable = a.symTable
}

func (a *Analyzer) VisitBinaryOpNode(n *ast.BinaryOpNode) {
	n.Left.Accept(a)
	left := a.stack.Pop()

	n.Right.Accept(a)
	right := a.stack.Pop()

	if left != right {
		panic("Type mis-match")
	}

	if (left == NUMBER || right == NUMBER) &&
		(n.Op == token.ADD || n.Op == token.SUBTRACT ||
			n.Op == token.MULTIPLY || n.Op == token.DIVIDE ||
			n.Op == token.MODULO) {
		a.stack.Push(NUMBER)
		return
	}
	if (left == STRING || right == STRING) &&
		(n.Op == token.ADD) {
		a.stack.Push(STRING)
		return
	}

	if n.Op == token.EQUAL || n.Op == token.NOT_EQUAL ||
		n.Op == token.LESS || n.Op == token.LESS_EQ ||
		n.Op == token.GREATER || n.Op == token.GREATER_EQ ||
		n.Op == token.AND || n.Op == token.OR {
		a.stack.Push(BOOL)
		return
	}

	a.stack.Push(UNKNOWN)
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
	dtype, ok := a.symTable.Get(n.Name, symtable.FUNCTYPE)
	if !ok {
		panic("Function not defined")
	}
	n.SymTable = a.symTable
	a.stack.Push(dtype)
}

func (a *Analyzer) VisitFuncDefNode(n *ast.FuncDefNode) {
	a.funcStack.Push(n.Name)
	a.symTable.Set(n.Name, symtable.FUNCTYPE, UNKNOWN)
	a.symTable = symtable.NewScope(a.symTable)

	for _, arg := range n.Args {
		a.symTable.Set(arg, symtable.VARTYPE, UNKNOWN)
	}
	for _, stmt := range n.Body {
		stmt.Accept(a)
		// set function's type from return type
		if a.stack.Size() > 0 {
			newType := a.stack.Pop()
			oldType, _ := a.symTable.P.Get(n.Name, symtable.FUNCTYPE)
			// multiple return statements may be present
			if oldType != UNKNOWN && newType != oldType {
				panic("Inconsistent return types")
			}
			a.symTable.P.Set(n.Name, symtable.FUNCTYPE, newType)
		}
	}
	n.SymTable = a.symTable
	a.symTable = a.symTable.Parent()
	a.funcStack.Pop()
}

func (a *Analyzer) VisitIfNode(n *ast.IfNode) {
	n.Condition.Accept(a)
	cond := a.stack.Pop()
	if cond != BOOL {
		panic("Non-bool expression used as condition")
	}
	for _, stmt := range n.Body {
		stmt.Accept(a)
		// if in a function, set its type from return
		if a.stack.Size() > 0 {
			newType := a.stack.Pop()
			oldType, _ := a.symTable.P.Get(
				a.funcStack.Peek().(string),
				symtable.FUNCTYPE,
			)
			// multiple return statements may be present
			if oldType != UNKNOWN && newType != oldType {
				panic("Inconsistent return types")
			}
			a.symTable.P.Set(
				a.funcStack.Peek().(string),
				symtable.FUNCTYPE,
				 newType,
			)
		}
	}
	if n.Else != nil {
		n.Else.Accept(a)
	}
}

func (a *Analyzer) VisitReturnNode(n *ast.ReturnNode) {
	n.Expr.Accept(a)
}

func (a *Analyzer) VisitUnaryOpNode(n *ast.UnaryOpNode) {
	n.Expr.Accept(a)
	t := a.stack.Pop()
	if (t == NUMBER && n.Op != token.ADD && n.Op != token.SUBTRACT) ||
		(t == BOOL && n.Op != token.NOT) {
		panic("Invalid type for operator")
	}
	a.stack.Push(t)
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
	dtype, ok := a.symTable.Get(n.Name, symtable.VARTYPE)
	if !ok {
		panic("Variable not defined")
	}
	n.SymTable = a.symTable
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
		// if in a function, set its type from return
		if a.stack.Size() > 0 {
			newType := a.stack.Pop()
			oldType, _ := a.symTable.P.Get(
				a.funcStack.Peek().(string),
				symtable.FUNCTYPE,
			)
			// multiple return statements may be present
			if oldType != UNKNOWN && newType != oldType {
				panic("Inconsistent return types")
			}
			a.symTable.P.Set(
				a.funcStack.Peek().(string),
				symtable.FUNCTYPE,
				 newType,
			)
		}
	}
}
