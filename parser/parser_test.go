package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scanner"
)

func newParser(s string) *Parser {
	return New(scanner.New(bytes.NewReader([]byte(s))))
}

func TestSkipComment(t *testing.T) {
	p := newParser("//\n \tfoo")
	str := p.ident()
	assert.Equal(t, "foo", str)
}

func TestParseNil(t *testing.T) {
	p := newParser("")
	result, _ := p.Parse()
	assert.Equal(t, 0, len(result.Stmts))
}

func TestParseRecovery(t *testing.T) {
	p := newParser("42")
	_, err := p.Parse()
	assert.NotNil(t, err)
}

func TestParseIdentifier(t *testing.T) {
	p := newParser("foo")
	str := p.ident()
	assert.Equal(t, "foo", str)
}

func TestParseIdentifierError(t *testing.T) {
	p := newParser("42")
	assert.Panics(t, func() {
		p.ident()
	})
}

func TestParseTermParens(t *testing.T) {
	p := newParser("(42)")
	node := p.term().(*ast.NumberNode)
	assert.Equal(t, 42.0, node.Value)
}

func TestParseTermSigned(t *testing.T) {
	p := newParser("-42")
	node := p.term().(*ast.NegativeNode)
	assert.Equal(t, 42.0, node.Term.(*ast.NumberNode).Value)
}

func TestParseCast(t *testing.T) {
	p := newParser("foo:string")
	node := p.castExpr().(*ast.CastNode)
	assert.Equal(t, "string", node.Cast)
	assert.Equal(t, "foo", node.Term.(*ast.VariableNode).Name)
}

func TestParseTermFuncCall(t *testing.T) {
	p := newParser("foo()")
	node := p.term().(*ast.FuncCallNode)
	assert.Equal(t, "foo", node.Name)
	assert.Equal(t, 0, len(node.Args))
}

func TestParseTerminalCallWithArgs(t *testing.T) {
	p := newParser("foo(bar, 42, \"baz\")")
	node := p.term().(*ast.FuncCallNode)
	assert.Equal(t, "foo", node.Name)
	assert.Equal(t, "bar", node.Args[0].(*ast.VariableNode).Name)
	assert.Equal(t, 42.0, node.Args[1].(*ast.NumberNode).Value)
	assert.Equal(t, "baz", node.Args[2].(*ast.StringNode).Value)
}

func TestParseBraceStmtListEmpty(t *testing.T) {
	p := newParser("{}")
	node := p.braceStmtList()
	assert.Equal(t, 0, len(node))
}

func TestParseBraceStmtList(t *testing.T) {
	p := newParser("{foo := 42 bar := 73}")
	node := p.braceStmtList()
	assert.Equal(t, "foo", node[0].(*ast.AssignNode).Name)
	assert.Equal(t, 42.0, node[0].(*ast.AssignNode).Expr.(*ast.NumberNode).Value)
	assert.Equal(t, "bar", node[1].(*ast.AssignNode).Name)
	assert.Equal(t, 73.0, node[1].(*ast.AssignNode).Expr.(*ast.NumberNode).Value)
}

func TestParseBraceStmtListStmtError(t *testing.T) {
	p := newParser("{foo := 42\n bar 73")
	assert.Panics(t, func() {
		p.braceStmtList()
	})
}

func TestParseBraceStmtListBraceError(t *testing.T) {
	p := newParser("{foo := 42\n")
	assert.Panics(t, func() {
		p.braceStmtList()
	})
}

func TestParseFuncDef(t *testing.T) {
	p := newParser("func foo {}")
	node := p.stmt().(*ast.FuncDefNode)
	assert.Equal(t, "foo", node.Name)
	assert.Equal(t, 0, len(node.Args))
	assert.Equal(t, 0, len(node.Body))
}

func TestParseFuncDefOneParam(t *testing.T) {
	p := newParser("func foo bar {}")
	node := p.stmt().(*ast.FuncDefNode)
	assert.Equal(t, "foo", node.Name)
	assert.Equal(t, "bar", node.Args[0])
}

func TestParseFuncDefManyParams(t *testing.T) {
	p := newParser("func foo bar baz {}")
	node := p.stmt().(*ast.FuncDefNode)
	assert.Equal(t, "foo", node.Name)
	assert.Equal(t, "bar", node.Args[0])
	assert.Equal(t, "baz", node.Args[1])
}

func TestParseIfStmt(t *testing.T) {
	p := newParser("if true {foo := 42}")
	node := p.stmt().(*ast.IfNode)
	assert.Equal(t, true, node.Cond.(*ast.BoolNode).Value)
	assert.Equal(t, "foo", node.Body[0].(*ast.AssignNode).Name)
}

func TestParseIfStmtExprError(t *testing.T) {
	p := newParser("if foo = {")
	assert.Panics(t, func() {
		p.stmt()
	})
}

func TestParseIfStmtBraceError(t *testing.T) {
	p := newParser("if true foo :=")
	assert.Panics(t, func() {
		p.stmt()
	})
}

func TestParseIfStmtWithElse(t *testing.T) {
	p := newParser("if false {} else false {} else {}")
	node := p.stmt().(*ast.IfNode)
	assert.Equal(t, false, node.Cond.(*ast.BoolNode).Value)
	assert.Equal(t, false, node.Else[0].(*ast.IfNode).Cond.(*ast.BoolNode).Value)
}

func TestParseReturnStmt(t *testing.T) {
	p := newParser("return 42\n")
	node := p.stmt().(*ast.ReturnNode)
	assert.Equal(t, 42.0, node.Expr.(*ast.NumberNode).Value)
}

func TestParseReturnStmtNoExpr(t *testing.T) {
	p := newParser("return }")
	node := p.stmt()
	assert.Nil(t, node.(*ast.ReturnNode).Expr)
}

func TestParseWhileStmt(t *testing.T) {
	p := newParser("while foo = true {bar := 42}")
	node := p.stmt().(*ast.WhileNode)
	assert.Equal(t, "foo", node.Cond.(*ast.EqualNode).Left.(*ast.VariableNode).Name)
	assert.Equal(t, true, node.Cond.(*ast.EqualNode).Right.(*ast.BoolNode).Value)
	assert.Equal(t, "bar", node.Body[0].(*ast.AssignNode).Name)
}

func TestParseWhileStmtExprError(t *testing.T) {
	p := newParser("while foo = {")
	assert.Panics(t, func() {
		p.stmt()
	})
}

func TestParseStmtError(t *testing.T) {
	p := newParser("\n")
	assert.Panics(t, func() {
		p.stmt()
	})
}

func TestParseAssignStmt(t *testing.T) {
	p := newParser("foo := 42 + 73\n")
	node := p.stmt().(*ast.AssignNode)
	assert.Equal(t, "foo", node.Name)
	assert.Equal(t, 42.0, node.Expr.(*ast.AddNode).Left.(*ast.NumberNode).Value)
	assert.Equal(t, 73.0, node.Expr.(*ast.AddNode).Right.(*ast.NumberNode).Value)
}

func TestParseAssignSmtExprError(t *testing.T) {
	p := newParser("foo := 42 +")
	assert.Panics(t, func() {
		p.stmt()
	})
}

func TestParseFuncCall(t *testing.T) {
	p := newParser("foo()\n")
	node := p.stmt().(*ast.FuncCallNode)
	assert.Equal(t, "foo", node.Name)
	assert.Equal(t, 0, len(node.Args))
}
