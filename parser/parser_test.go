package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/token"
	"testing"
)

type tokenPair struct {
	token token.Token
	value string
}

type mockScanner struct {
	i  uint8
	tp []tokenPair
}

func NewMockScanner() *mockScanner {
	return &mockScanner{i: 0}
}

func (s *mockScanner) reset(pairs []tokenPair) {
	s.i = 0
	s.tp = pairs
}

func (s *mockScanner) Scan() (token.Token, string) {
	t := s.tp[s.i].token
	v := s.tp[s.i].value
	s.i++
	return t, v
}

func TestSkipComment(t *testing.T) {
	s := NewMockScanner()
	s.reset([]tokenPair{
		{token.COMMENT, "//"},
		{token.STRING, ""},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Equal(t, token.STRING, p.token)
}

func TestParseIdentifier(t *testing.T) {
	s := NewMockScanner()
	// foo
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	str := p.identifier()
	assert.Equal(t, "foo", str)
}

func TestParseIdentifierError(t *testing.T) {
	s := NewMockScanner()
	// foo
	s.reset([]tokenPair{
		{token.NUMBER, "42"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.identifier()
	})
}

func TestParseTerm(t *testing.T) {
	s := NewMockScanner()
	// 42 * 73
	s.reset([]tokenPair{
		{token.NUMBER, "42"},
		{token.MULTIPLY, "*"},
		{token.NUMBER, "73"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.term()
	assert.Equal(t, token.MULTIPLY, node.(ast.BinaryExpr).Op)
	assert.Equal(t, token.NUMBER, node.(ast.BinaryExpr).Left.(ast.ValueExpr).Type)
	assert.Equal(t, token.NUMBER, node.(ast.BinaryExpr).Right.(ast.ValueExpr).Type)
}

func TestParseTermError(t *testing.T) {
	s := NewMockScanner()
	// 42 *
	s.reset([]tokenPair{
		{token.NUMBER, "42"},
		{token.MULTIPLY, "*"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.term()
	})
}

func TestParseSimpleExpr(t *testing.T) {
	s := NewMockScanner()
	// 42 + 73
	s.reset([]tokenPair{
		{token.NUMBER, "42"},
		{token.ADD, "+"},
		{token.NUMBER, "73"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.simpleExpr()
	assert.Equal(t, token.ADD, node.(ast.BinaryExpr).Op)
	assert.Equal(t, token.NUMBER, node.(ast.BinaryExpr).Left.(ast.ValueExpr).Type)
	assert.Equal(t, token.NUMBER, node.(ast.BinaryExpr).Right.(ast.ValueExpr).Type)
}

func TestParseSimpleExprError(t *testing.T) {
	s := NewMockScanner()
	// 42 +
	s.reset([]tokenPair{
		{token.NUMBER, "42"},
		{token.ADD, "+"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.simpleExpr()
	})
}

func TestParseRelation(t *testing.T) {
	s := NewMockScanner()
	s.reset([]tokenPair{
		// 42 < 73
		{token.NUMBER, "42"},
		{token.LESS, "<"},
		{token.NUMBER, "73"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.relation()
	assert.Equal(t, token.LESS, node.(ast.BinaryExpr).Op)
	assert.Equal(t, token.NUMBER, node.(ast.BinaryExpr).Left.(ast.ValueExpr).Type)
	assert.Equal(t, token.NUMBER, node.(ast.BinaryExpr).Right.(ast.ValueExpr).Type)
}

func TestParseRelationError(t *testing.T) {
	s := NewMockScanner()
	// 42 <
	s.reset([]tokenPair{
		{token.NUMBER, "42"},
		{token.LESS, "<"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.relation()
	})
}

func TestParseExpr(t *testing.T) {
	s := NewMockScanner()
	// true && true
	s.reset([]tokenPair{
		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.TRUE, "true"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.expr()
	assert.Equal(t, token.AND, node.(ast.BinaryExpr).Op)
	assert.Equal(t, token.TRUE, node.(ast.BinaryExpr).Left.(ast.ValueExpr).Type)
	assert.Equal(t, token.TRUE, node.(ast.BinaryExpr).Right.(ast.ValueExpr).Type)
}

func TestParseExprError(t *testing.T) {
	s := NewMockScanner()
	// true &&
	s.reset([]tokenPair{
		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.expr()
	})
}

func TestParseFactorParens(t *testing.T) {
	s := NewMockScanner()
	// ( 42 )
	s.reset([]tokenPair{
		{token.LPAREN, "("},
		{token.NUMBER, "42"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.factor()
	assert.Equal(t, token.NUMBER, node.(ast.ValueExpr).Type)
}

func TestParseFactorParensExprError(t *testing.T) {
	s := NewMockScanner()
	// ( EOF
	s.reset([]tokenPair{
		{token.LPAREN, "("},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.factor()
	})
}

func TestParseFactorParensCloseError(t *testing.T) {
	s := NewMockScanner()
	// ( 42 EOF
	s.reset([]tokenPair{
		{token.LPAREN, "("},
		{token.NUMBER, "42"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.factor()
	})
}

func TestParseFactorSigned(t *testing.T) {
	s := NewMockScanner()
	// -42
	s.reset([]tokenPair{
		{token.SUBTRACT, "-"},
		{token.NUMBER, "42"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.factor()
	assert.Equal(t, token.SUBTRACT, node.(ast.UnaryExpr).Op)
	assert.Equal(t, token.NUMBER, node.(ast.UnaryExpr).Right.(ast.ValueExpr).Type)
}

func TestParseTerminalVariable(t *testing.T) {
	s := NewMockScanner()
	// foo
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.terminal()
	assert.Equal(t, "foo", node.(ast.VariableExpr).Name)
}

func TestParseTerminalCall(t *testing.T) {
	s := NewMockScanner()
	// foo()
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.terminal()
	assert.Equal(t, "foo", node.(ast.FuncCall).Name)
}

func TestParseTerminalCallWithArgs(t *testing.T) {
	s := NewMockScanner()
	// foo(bar, 42, "baz")
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "bar"},
		{token.COMMA, ","},
		{token.NUMBER, "42"},
		{token.COMMA, ","},
		{token.STRING, "baz"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.terminal()
	assert.Equal(t, "foo", node.(ast.FuncCall).Name)
	assert.Equal(t, "bar", node.(ast.FuncCall).Args[0].(ast.VariableExpr).Name)
	assert.Equal(t, token.NUMBER, node.(ast.FuncCall).Args[1].(ast.ValueExpr).Type)
	assert.Equal(t, token.STRING, node.(ast.FuncCall).Args[2].(ast.ValueExpr).Type)
}

func TestParseTerminalCallWithArgsExprError(t *testing.T) {
	s := NewMockScanner()
	// foo(bar, 42 < ,
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "bar"},
		{token.COMMA, ","},
		{token.NUMBER, "42"},
		{token.LESS, "<"},
		{token.COMMA, ","},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.terminal()
	})
}

func TestParseTerminalFuncCallArgsListError(t *testing.T) {
	s := NewMockScanner()
	// foo(bar 42
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "bar"},
		{token.IDENTIFIER, "42"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.terminal()
	})
}

func TestParseBraceStmtListEmpty(t *testing.T) {
	s := NewMockScanner()
	// { }
	s.reset([]tokenPair{
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.braceStmtList()
	assert.Equal(t, 0, len(node))
}

func TestParseBraceStmtList(t *testing.T) {
	s := NewMockScanner()
	// { foo := 42. bar := 73. }
	s.reset([]tokenPair{
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "foo"},
		{token.ASSIGN, ":="},
		{token.NUMBER, "42"},
		{token.DOT, "."},
		{token.IDENTIFIER, "bar"},
		{token.ASSIGN, ":="},
		{token.NUMBER, "73"},
		{token.DOT, "."},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.braceStmtList()
	assert.Equal(t, "foo", node[0].(ast.AssignStmt).Name)
	assert.Equal(t, "bar", node[1].(ast.AssignStmt).Name)
}

func TestParseBraceStmtListStmtError(t *testing.T) {
	s := NewMockScanner()
	// { foo := 42. bar 73
	s.reset([]tokenPair{
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "foo"},
		{token.ASSIGN, ":="},
		{token.NUMBER, "42"},
		{token.DOT, "."},
		{token.IDENTIFIER, "bar"},
		{token.NUMBER, "73"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.braceStmtList()
	})
}

func TestParseBraceStmtListBraceError(t *testing.T) {
	s := NewMockScanner()
	// { foo := 42.
	s.reset([]tokenPair{
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "foo"},
		{token.ASSIGN, ":="},
		{token.NUMBER, "42"},
		{token.DOT, "."},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.braceStmtList()
	})
}

func TestParseFuncDef(t *testing.T) {
	s := NewMockScanner()
	// func foo {}
	s.reset([]tokenPair{
		{token.FUNC, "func"},
		{token.IDENTIFIER, "foo"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.stmt()
	assert.Equal(t, "foo", node.(ast.FuncDef).Name)
}

func TestParseFuncDefOneParam(t *testing.T) {
	s := NewMockScanner()
	// func foo bar {}
	s.reset([]tokenPair{
		{token.FUNC, "func"},
		{token.IDENTIFIER, "foo"},
		{token.IDENTIFIER, "bar"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.stmt()
	assert.Equal(t, "foo", node.(ast.FuncDef).Name)
	assert.Equal(t, "bar", node.(ast.FuncDef).Args[0])
}

func TestParseFuncDefManyParams(t *testing.T) {
	s := NewMockScanner()
	// func foo bar baz {}
	s.reset([]tokenPair{
		{token.FUNC, "func"},
		{token.IDENTIFIER, "foo"},
		{token.IDENTIFIER, "bar"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "baz"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.stmt()
	assert.Equal(t, "foo", node.(ast.FuncDef).Name)
	assert.Equal(t, "bar", node.(ast.FuncDef).Args[0])
	assert.Equal(t, "baz", node.(ast.FuncDef).Args[1])
}

func TestParseIfStmt(t *testing.T) {
	s := NewMockScanner()
	// if true { foo := 42. }
	s.reset([]tokenPair{
		{token.IF, "if"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "foo"},
		{token.ASSIGN, ":="},
		{token.NUMBER, "42"},
		{token.DOT, "."},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.stmt()
	assert.Equal(t, token.TRUE, node.(ast.IfStmt).Condition.(ast.ValueExpr).Type)
	assert.Equal(t, "foo", node.(ast.IfStmt).Body[0].(ast.AssignStmt).Name)
}

func TestParseIfStmtExprError(t *testing.T) {
	s := NewMockScanner()
	// if foo = {
	s.reset([]tokenPair{
		{token.IF, "if"},
		{token.IDENTIFIER, "foo"},
		{token.EQUAL, "="},
		{token.LBRACE, "{"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.stmt()
	})
}

func TestParseIfStmtBraceError(t *testing.T) {
	s := NewMockScanner()
	// if true foo :=
	s.reset([]tokenPair{
		{token.IF, "if"},
		{token.TRUE, "true"},
		{token.IDENTIFIER, "foo"},
		{token.ASSIGN, ":="},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.stmt()
	})
}

func TestParseReturnStmt(t *testing.T) {
	s := NewMockScanner()
	// return 42.
	s.reset([]tokenPair{
		{token.RETURN, "return"},
		{token.NUMBER, "42"},
		{token.DOT, "."},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.stmt()
	assert.Equal(t, token.NUMBER, node.(ast.ReturnStmt).Expr.(ast.ValueExpr).Type)
}

func TestParseReturnStmtNoExpr(t *testing.T) {
	s := NewMockScanner()
	// return.
	s.reset([]tokenPair{
		{token.RETURN, "return"},
		{token.DOT, "."},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.stmt()
	assert.Nil(t, node.(ast.ReturnStmt).Expr)
}

func TestParseWhileStmt(t *testing.T) {
	s := NewMockScanner()
	// while foo = true { bar := 42. }
	s.reset([]tokenPair{
		{token.WHILE, "while"},
		{token.IDENTIFIER, "foo"},
		{token.EQUAL, "="},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "bar"},
		{token.ASSIGN, ":="},
		{token.NUMBER, "42"},
		{token.DOT, "."},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.stmt()
	assert.Equal(t, token.EQUAL, node.(ast.WhileStmt).Condition.(ast.BinaryExpr).Op)
	assert.Equal(t, "bar", node.(ast.WhileStmt).Body[0].(ast.AssignStmt).Name)
}

func TestParseWhileStmtExprError(t *testing.T) {
	s := NewMockScanner()
	// while foo = {
	s.reset([]tokenPair{
		{token.WHILE, "while"},
		{token.IDENTIFIER, "foo"},
		{token.EQUAL, "="},
		{token.LBRACE, "{"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.stmt()
	})
}

func TestParseStmtError(t *testing.T) {
	s := NewMockScanner()
	// .
	s.reset([]tokenPair{
		{token.DOT, "."},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.stmt()
	})
}

func TestParseAssignStmt(t *testing.T) {
	s := NewMockScanner()
	// foo := 42 + 73.
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.ASSIGN, ":="},
		{token.NUMBER, "42"},
		{token.ADD, "+"},
		{token.NUMBER, "73"},
		{token.DOT, "."},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node := p.stmt()
	assert.Equal(t, "foo", node.(ast.AssignStmt).Name)
	assert.Equal(t, token.ADD, node.(ast.AssignStmt).Expr.(ast.BinaryExpr).Op)
}

func TestParseAssignSmtExprError(t *testing.T) {
	s := NewMockScanner()
	// foo := 42 +
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.ASSIGN, ":="},
		{token.DOT, "."},
		{token.NUMBER, "42"},
		{token.ADD, "+"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.stmt()
	})
}

func TestFuncCall(t *testing.T) {
	s := NewMockScanner()
	// foo().
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.DOT, "."},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)
	node := p.stmt()
	assert.Equal(t, "foo", node.(ast.FuncCall).Name)
	assert.Equal(t, 0, len(node.(ast.FuncCall).Args))
}

func TestStmtTerminationError(t *testing.T) {
	s := NewMockScanner()
	// foo := 42
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.ASSIGN, ":="},
		{token.NUMBER, "42"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Panics(t, func() {
		p.stmt()
	})
}
