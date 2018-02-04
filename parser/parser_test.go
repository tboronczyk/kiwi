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
		node := p.term().(*ast.NumberNode)
		assert.Equal(t, 42.0, node.Value)
	})

	t.Run("Parse signed term", func(t *testing.T) {
		p := newParser("-42")
		node := p.term().(*ast.NegativeNode)
		assert.Equal(t, 42.0, node.Term.(*ast.NumberNode).Value)
	})

	t.Run("Parse cast", func(t *testing.T) {
		p := newParser("foo:string")
		node := p.castExpr().(*ast.CastNode)
		assert.Equal(t, "string", node.Cast)
		assert.Equal(t, "foo", node.Term.(*ast.VariableNode).Name)
	})

	t.Run("Parse func call term", func(t *testing.T) {
		p := newParser("foo()")
		node := p.term().(*ast.FuncCallNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, 0, len(node.Args))
	})

	t.Run("Parse term func call with args", func(t *testing.T) {
		p := newParser("foo(bar, 42, \"baz\")")
		node := p.term().(*ast.FuncCallNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, "bar", node.Args[0].(*ast.VariableNode).Name)
		assert.Equal(t, 42.0, node.Args[1].(*ast.NumberNode).Value)
		assert.Equal(t, "baz", node.Args[2].(*ast.StringNode).Value)
	})

	t.Run("Parse empty braced statement list", func(t *testing.T) {
		p := newParser("{}")
		node := p.braceStmtList()
		assert.Equal(t, 0, len(node))
	})

	t.Run("Parse braced statement list", func(t *testing.T) {
		p := newParser("{foo := 42 bar := 73}")
		node := p.braceStmtList()
		assert.Equal(t, "foo", node[0].(*ast.AssignNode).Name)
		assert.Equal(t, 42.0, node[0].(*ast.AssignNode).Expr.(*ast.NumberNode).Value)
		assert.Equal(t, "bar", node[1].(*ast.AssignNode).Name)
		assert.Equal(t, 73.0, node[1].(*ast.AssignNode).Expr.(*ast.NumberNode).Value)
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
		node := p.stmt().(*ast.FuncDefNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, 0, len(node.Args))
		assert.Equal(t, 0, len(node.Body))
	})

	t.Run("Parse function def with one parameter", func(t *testing.T) {
		p := newParser("func foo bar {}")
		node := p.stmt().(*ast.FuncDefNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, "bar", node.Args[0])
	})

	t.Run("Parse function def with many parameters", func(t *testing.T) {
		p := newParser("func foo bar baz {}")
		node := p.stmt().(*ast.FuncDefNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, "bar", node.Args[0])
		assert.Equal(t, "baz", node.Args[1])
	})

	t.Run("Parse if statement", func(t *testing.T) {
		p := newParser("if true {foo := 42}")
		node := p.stmt().(*ast.IfNode)
		assert.Equal(t, true, node.Cond.(*ast.BoolNode).Value)
		assert.Equal(t, "foo", node.Body[0].(*ast.AssignNode).Name)
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
		node := p.stmt().(*ast.IfNode)
		assert.Equal(t, false, node.Cond.(*ast.BoolNode).Value)
		assert.Equal(t, false, node.Else[0].(*ast.IfNode).Cond.(*ast.BoolNode).Value)
	})

	t.Run("Parse return statement", func(t *testing.T) {
		p := newParser("return 42\n")
		node := p.stmt().(*ast.ReturnNode)
		assert.Equal(t, 42.0, node.Expr.(*ast.NumberNode).Value)
	})

	t.Run("Parse return statement without expression", func(t *testing.T) {
		p := newParser("return }")
		node := p.stmt()
		assert.Nil(t, node.(*ast.ReturnNode).Expr)
	})

	t.Run("Parse while statement", func(t *testing.T) {
		p := newParser("while foo = true {bar := 42}")
		node := p.stmt().(*ast.WhileNode)
		assert.Equal(t, "foo", node.Cond.(*ast.EqualNode).Left.(*ast.VariableNode).Name)
		assert.Equal(t, true, node.Cond.(*ast.EqualNode).Right.(*ast.BoolNode).Value)
		assert.Equal(t, "bar", node.Body[0].(*ast.AssignNode).Name)
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
		node := p.stmt().(*ast.AssignNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, 42.0, node.Expr.(*ast.AddNode).Left.(*ast.NumberNode).Value)
		assert.Equal(t, 73.0, node.Expr.(*ast.AddNode).Right.(*ast.NumberNode).Value)
	})

	t.Run("Parse assignment statement with expression error", func(t *testing.T) {
		p := newParser("foo := 42 +")
		assert.Panics(t, func() {
			p.stmt()
		})
	})

	t.Run("Parse function call", func(t *testing.T) {
		p := newParser("foo()\n")
		node := p.stmt().(*ast.FuncCallNode)
		assert.Equal(t, "foo", node.Name)
		assert.Equal(t, 0, len(node.Args))
	})
}
