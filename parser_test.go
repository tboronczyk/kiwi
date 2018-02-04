package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newParser(s string) *Parser {
	return NewParser(NewScanner(bytes.NewReader([]byte(s))))
}

func TestParser(t *testing.T) {
	t.Parallel()

	t.Run("Skip comment", func(t *testing.T) {
		p := newParser("//\n \tfoo")
		str := p.ident()
		assert.Equal(t, "foo", str)
	})

	t.Run("Parse nil", func(t *testing.T) {
		p := newParser("")
		result, _ := p.Parse()
		assert.Equal(t, 0, len(result.Stmts))
	})

	t.Run("Test parse recovery", func(t *testing.T) {
		p := newParser("42")
		_, err := p.Parse()
		assert.NotNil(t, err)
	})

	t.Run("Parse identifier", func(t *testing.T) {
		p := newParser("foo")
		str := p.ident()
		assert.Equal(t, "foo", str)
	})

	t.Run("Parse identifier error", func(t *testing.T) {
		p := newParser("42")
		assert.Panics(t, func() {
			p.ident()
		})
	})

	t.Run("Parse parenthesized term", func(t *testing.T) {
		p := newParser("(42)")
		node := p.term().(*AstNumberNode)
		assert.Equal(t, 42.0, node.Value)
	})

	t.Run("Parse signed term", func(t *testing.T) {
		p := newParser("-42")
		node := p.term().(*AstNegativeNode)
		assert.Equal(t, 42.0, node.Term.(*AstNumberNode).Value)
	})

	t.Run("Parse cast", func(t *testing.T) {
		p := newParser("foo:string")
		node := p.castExpr().(*AstCastNode)
		assert.Equal(t, "string", node.Cast)
		assert.Equal(t, "foo", node.Term.(*AstVariableNode).Name)
	})

	t.Run("Parse func call term", func(t *testing.T) {
		p := newParser("foo()")
		node := p.term().(*AstFuncCallNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, 0, len(node.Args))
	})

	t.Run("Parse term func call with args", func(t *testing.T) {
		p := newParser("foo(bar, 42, \"baz\")")
		node := p.term().(*AstFuncCallNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, "bar", node.Args[0].(*AstVariableNode).Name)
		assert.Equal(t, 42.0, node.Args[1].(*AstNumberNode).Value)
		assert.Equal(t, "baz", node.Args[2].(*AstStringNode).Value)
	})

	t.Run("Parse empty braced statement list", func(t *testing.T) {
		p := newParser("{}")
		node := p.braceStmtList()
		assert.Equal(t, 0, len(node))
	})

	t.Run("Parse braced statement list", func(t *testing.T) {
		p := newParser("{foo := 42 bar := 73}")
		node := p.braceStmtList()
		assert.Equal(t, "foo", node[0].(*AstAssignNode).Name)
		assert.Equal(t, 42.0, node[0].(*AstAssignNode).Expr.(*AstNumberNode).Value)
		assert.Equal(t, "bar", node[1].(*AstAssignNode).Name)
		assert.Equal(t, 73.0, node[1].(*AstAssignNode).Expr.(*AstNumberNode).Value)
	})

	t.Run("Parse braced statement list with statement error", func(t *testing.T) {
		p := newParser("{foo := 42\n bar 73")
		assert.Panics(t, func() {
			p.braceStmtList()
		})
	})

	t.Run("Parse braced statement list with brace error", func(t *testing.T) {
		p := newParser("{foo := 42\n")
		assert.Panics(t, func() {
			p.braceStmtList()
		})
	})

	t.Run("Parse function def", func(t *testing.T) {
		p := newParser("func foo {}")
		node := p.stmt().(*AstFuncDefNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, 0, len(node.Args))
		assert.Equal(t, 0, len(node.Body))
	})

	t.Run("Parse function def with one parameter", func(t *testing.T) {
		p := newParser("func foo bar {}")
		node := p.stmt().(*AstFuncDefNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, "bar", node.Args[0])
	})

	t.Run("Parse function def with many parameters", func(t *testing.T) {
		p := newParser("func foo bar baz {}")
		node := p.stmt().(*AstFuncDefNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, "bar", node.Args[0])
		assert.Equal(t, "baz", node.Args[1])
	})

	t.Run("Parse if statement", func(t *testing.T) {
		p := newParser("if true {foo := 42}")
		node := p.stmt().(*AstIfNode)
		assert.Equal(t, true, node.Cond.(*AstBoolNode).Value)
		assert.Equal(t, "foo", node.Body[0].(*AstAssignNode).Name)
	})

	t.Run("Parse if statement with expression error", func(t *testing.T) {
		p := newParser("if foo = {")
		assert.Panics(t, func() {
			p.stmt()
		})
	})

	t.Run("Parse if statement with brace error", func(t *testing.T) {
		p := newParser("if true foo :=")
		assert.Panics(t, func() {
			p.stmt()
		})
	})

	t.Run("Parse if statement with else", func(t *testing.T) {
		p := newParser("if false {} else false {} else {}")
		node := p.stmt().(*AstIfNode)
		assert.Equal(t, false, node.Cond.(*AstBoolNode).Value)
		assert.Equal(t, false, node.Else[0].(*AstIfNode).Cond.(*AstBoolNode).Value)
	})

	t.Run("Parse return statement", func(t *testing.T) {
		p := newParser("return 42\n")
		node := p.stmt().(*AstReturnNode)
		assert.Equal(t, 42.0, node.Expr.(*AstNumberNode).Value)
	})

	t.Run("Parse return statement without expression", func(t *testing.T) {
		p := newParser("return }")
		node := p.stmt()
		assert.Nil(t, node.(*AstReturnNode).Expr)
	})

	t.Run("Parse while statement", func(t *testing.T) {
		p := newParser("while foo = true {bar := 42}")
		node := p.stmt().(*AstWhileNode)
		assert.Equal(t, "foo", node.Cond.(*AstEqualNode).Left.(*AstVariableNode).Name)
		assert.Equal(t, true, node.Cond.(*AstEqualNode).Right.(*AstBoolNode).Value)
		assert.Equal(t, "bar", node.Body[0].(*AstAssignNode).Name)
	})

	t.Run("Parse while statement with expression error", func(t *testing.T) {
		p := newParser("while foo = {")
		assert.Panics(t, func() {
			p.stmt()
		})
	})

	t.Run("Test parse statement error", func(t *testing.T) {
		p := newParser("\n")
		assert.Panics(t, func() {
			p.stmt()
		})
	})

	t.Run("Parse assignment statement", func(t *testing.T) {
		p := newParser("foo := 42 + 73\n")
		node := p.stmt().(*AstAssignNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, 42.0, node.Expr.(*AstAddNode).Left.(*AstNumberNode).Value)
		assert.Equal(t, 73.0, node.Expr.(*AstAddNode).Right.(*AstNumberNode).Value)
	})

	t.Run("Parse assignment statement with expression error", func(t *testing.T) {
		p := newParser("foo := 42 +")
		assert.Panics(t, func() {
			p.stmt()
		})
	})

	t.Run("Parse function call", func(t *testing.T) {
		p := newParser("foo()\n")
		node := p.stmt().(*AstFuncCallNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, 0, len(node.Args))
	})
}
