package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scanner"
	"github.com/tboronczyk/kiwi/token"
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
	assert.Nil(t, result)
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

func TestParseExprOpPrecedenceLower(t *testing.T) {
	p := newParser("42 + 73 * 101")
	node := p.expr().(*ast.BinOpNode)
	assert.Equal(t, "42", node.Left.(*ast.ValueNode).Value)
	assert.Equal(t, token.ADD, node.Op)
	assert.Equal(t, token.MULTIPLY, node.Right.(*ast.BinOpNode).Op)
	assert.Equal(t, "73", node.Right.(*ast.BinOpNode).Left.(*ast.ValueNode).Value)
	assert.Equal(t, "101", node.Right.(*ast.BinOpNode).Right.(*ast.ValueNode).Value)
}

func TestParseExprOpPrecedenceHigher(t *testing.T) {
	p := newParser("42 * 73 + 101")
	node := p.expr().(*ast.BinOpNode)
	assert.Equal(t, "42", node.Left.(*ast.BinOpNode).Left.(*ast.ValueNode).Value)
	assert.Equal(t, token.MULTIPLY, node.Left.(*ast.BinOpNode).Op)
	assert.Equal(t, "73", node.Left.(*ast.BinOpNode).Right.(*ast.ValueNode).Value)
	assert.Equal(t, token.ADD, node.Op)
	assert.Equal(t, "101", node.Right.(*ast.ValueNode).Value)
}

func TestParseTermParens(t *testing.T) {
	p := newParser("(42)")
	node := p.term().(*ast.ValueNode)
	assert.Equal(t, "42", node.Value)
	assert.Equal(t, token.NUMBER, node.Type)
}

func TestParseTermSigned(t *testing.T) {
	p := newParser("-42")
	node := p.term().(*ast.UnaryOpNode)
	assert.Equal(t, token.SUBTRACT, node.Op)
	assert.Equal(t, "42", node.Term.(*ast.ValueNode).Value)
	assert.Equal(t, token.NUMBER, node.Term.(*ast.ValueNode).Type)
}

func TestParseCast(t *testing.T) {
	p := newParser("foo:string")
	node := p.term().(*ast.CastNode)
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
	assert.Equal(t, token.NUMBER, node.Args[1].(*ast.ValueNode).Type)
	assert.Equal(t, token.STRING, node.Args[2].(*ast.ValueNode).Type)
}

func TestParseBraceStmtListEmpty(t *testing.T) {
	p := newParser("{}")
	node := p.braceStmtList()
	assert.Equal(t, 0, len(node))
}

func TestParseBraceStmtList(t *testing.T) {
	p := newParser("{foo := 42\nbar := 73}")
	node := p.braceStmtList()
	assert.Equal(t, "foo", node[0].(*ast.AssignNode).Name)
	assert.Equal(t, "42", node[0].(*ast.AssignNode).Expr.(*ast.ValueNode).Value)
	assert.Equal(t, "bar", node[1].(*ast.AssignNode).Name)
	assert.Equal(t, "73", node[1].(*ast.AssignNode).Expr.(*ast.ValueNode).Value)
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
	assert.Equal(t, "TRUE", node.Condition.(*ast.ValueNode).Value)
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
	assert.Equal(t, "FALSE", node.Condition.(*ast.ValueNode).Value)
	assert.Equal(t, "FALSE", node.Else.(*ast.IfNode).Condition.(*ast.ValueNode).Value)
	assert.Equal(t, "TRUE", node.Else.(*ast.IfNode).Else.(*ast.IfNode).Condition.(*ast.ValueNode).Value)
}

func TestParseReturnStmt(t *testing.T) {
	p := newParser("return 42\n")
	node := p.stmt().(*ast.ReturnNode)
	assert.Equal(t, "42", node.Expr.(*ast.ValueNode).Value)
}

func TestParseReturnStmtNoExpr(t *testing.T) {
	p := newParser("return\n")
	node := p.stmt()
	assert.Nil(t, node.(*ast.ReturnNode).Expr)
}

func TestParseReturnStmtError(t *testing.T) {
	p := newParser("return }")
	assert.Panics(t, func() {
		p.stmt()
	})
}

func TestParseWhileStmt(t *testing.T) {
	p := newParser("while foo = true {bar := 42}")
	node := p.stmt().(*ast.WhileNode)
	assert.Equal(t, token.EQUAL, node.Condition.(*ast.BinOpNode).Op)
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
	assert.Equal(t, token.ADD, node.Expr.(*ast.BinOpNode).Op)
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

func TestParseStmtTerminationError(t *testing.T) {
	p := newParser("foo := 42")
	assert.Panics(t, func() {
		p.stmt()
	})
}
