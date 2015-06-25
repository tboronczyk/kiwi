package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/token"
)

func TestVisitValueNodeBool(t *testing.T) {
	node := &ast.ValueNode{Type: token.BOOL}
	a := New()
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, BOOL, actual)
}

func TestVisitValueNodeNumber(t *testing.T) {
	node := &ast.ValueNode{Type: token.NUMBER}
	a := New()
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, NUMBER, actual)
}

func TestVisitValueNodeString(t *testing.T) {
	node := &ast.ValueNode{Type: token.STRING}
	a := New()
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, STRING, actual)
}

func TestVisitValueNodeUnknown(t *testing.T) {
	node := &ast.ValueNode{Type: token.IDENTIFIER}
	a := New()

	assert.Panics(t, func() {
		node.Accept(a)
	})
}

func TestVisitUnaryNodeBoolNot(t *testing.T) {
	node := &ast.UnaryOpNode{
		Op:   token.NOT,
		Expr: &ast.ValueNode{Type: token.BOOL},
	}
	a := New()
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, BOOL, actual)
}

func TestVisitUnaryNodeNumberPos(t *testing.T) {
	node := &ast.UnaryOpNode{
		Op:   token.ADD,
		Expr: &ast.ValueNode{Type: token.NUMBER},
	}
	a := New()
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, NUMBER, actual)
}

func TestVisitUnaryNodeNeg(t *testing.T) {
	node := &ast.UnaryOpNode{
		Op:   token.SUBTRACT,
		Expr: &ast.ValueNode{Type: token.NUMBER},
	}
	a := New()
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, NUMBER, actual)
}

func TestVisitUnaryNodeInvalid(t *testing.T) {
	node := &ast.UnaryOpNode{
		Op:   token.ADD,
		Expr: &ast.ValueNode{Type: token.BOOL},
	}
	a := New()

	assert.Panics(t, func() {
		node.Accept(a)
	})
}

func TestVisitVariableNodeAssigned(t *testing.T) {
	node := &ast.VariableNode{Name: "foo"}
	a := New()
	a.symTable.Set("foo", STRING)
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, STRING, actual)
}

func TestVisitVariableNodeNotAssigned(t *testing.T) {
	node := &ast.VariableNode{Name: "foo"}
	a := New()

	assert.Panics(t, func() {
		node.Accept(a)
	})
}

func TestVisitAssignNode(t *testing.T) {
	node := &ast.AssignNode{
		Name: "foo",
		Expr: &ast.ValueNode{Type: token.STRING},
	}
	a := New()
	node.Accept(a)

	expr, ok := a.symTable.Get("foo")
	assert.Equal(t, STRING, expr)
	assert.True(t, ok)
}

func TestVisitFuncDefNode(t *testing.T) {
	node := &ast.FuncDefNode{
		Name: "foo",
		Args: []string{"bar", "baz"},
		Body: []ast.Node{
			&ast.AssignNode{
				Name: "qux",
				Expr: &ast.ValueNode{Type: token.STRING},
			},
		},
	}
	a := New()
	node.Accept(a)

	expr, ok := a.symTable.Get("foo")
	assert.Equal(t, expr, UNKNOWN)
	assert.True(t, ok)

	value, _ := node.SymTable.Table["bar"]
	assert.Equal(t, value, UNKNOWN)

	value, _ = node.SymTable.Table["baz"]
	assert.Equal(t, value, UNKNOWN)

	value, _ = node.SymTable.Table["qux"]
	assert.Equal(t, value, STRING)

	_, ok = node.SymTable.Table["norf"]
	assert.False(t, ok)
}

func TestVisitBinaryOpNodeSameType(t *testing.T) {
	node := &ast.BinaryOpNode{
		Op:    token.ADD,
		Left:  &ast.ValueNode{Type: token.NUMBER},
		Right: &ast.ValueNode{Type: token.NUMBER},
	}
	a := New()
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, NUMBER, actual)
}

func TestVisitBinaryOpNodeAnyType(t *testing.T) {
	node := &ast.BinaryOpNode{
		Op:    token.ADD,
		Left:  &ast.VariableNode{Name: "foo"},
		Right: &ast.VariableNode{Name: "bar"},
	}
	a := New()
	a.symTable.Set("foo", UNKNOWN)
	a.symTable.Set("bar", UNKNOWN)
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, UNKNOWN, actual)
}

func TestVisitBinaryOpNodeTypeFail(t *testing.T) {
	node := &ast.BinaryOpNode{
		Op:    token.ADD,
		Left:  &ast.VariableNode{Name: "foo"},
		Right: &ast.VariableNode{Name: "bar"},
	}
	a := New()
	a.symTable.Set("foo", STRING)
	a.symTable.Set("bar", NUMBER)

	assert.Panics(t, func() {
		node.Accept(a)
	})
}

func TestVisitFuncCallNode(t *testing.T) {
	node := &ast.FuncCallNode{Name: "foo"}
	a := New()
	a.symTable.Set("foo", UNKNOWN)
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, UNKNOWN, actual)
}

func TestVisitFuncCallNodeNoExist(t *testing.T) {
	node := &ast.FuncCallNode{Name: "foo"}
	a := New()
	assert.Panics(t, func() {
		node.Accept(a)
	})
}

func TestVisitCastNodeBool(t *testing.T) {
	node := &ast.CastNode{
		Cast: "bool",
		Expr: &ast.ValueNode{Type: token.NUMBER},
	}
	a := New()
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, BOOL, actual)
}

func TestVisitCastNodeNumber(t *testing.T) {
	node := &ast.CastNode{
		Cast: "number",
		Expr: &ast.ValueNode{Type: token.BOOL},
	}
	a := New()
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, NUMBER, actual)
}

func TestVisitCastNodeString(t *testing.T) {
	node := &ast.CastNode{
		Cast: "string",
		Expr: &ast.ValueNode{Type: token.NUMBER},
	}
	a := New()
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, STRING, actual)
}

func TestVisitCastNodeUnknown(t *testing.T) {
	node := &ast.CastNode{
		Cast: "foo",
		Expr: &ast.ValueNode{Type: token.NUMBER},
	}
	a := New()
	node.Accept(a)

	actual := a.stack.Pop()
	assert.Equal(t, UNKNOWN, actual)
}

func TestVisitIfNode(t *testing.T) {
	node := &ast.IfNode{
		Condition: &ast.ValueNode{Type: token.BOOL},
		Body: []ast.Node{
			&ast.AssignNode{
				Name: "foo",
				Expr: &ast.ValueNode{Type: token.NUMBER},
			},
		},
		Else: &ast.IfNode{
			Condition: &ast.ValueNode{Type: token.BOOL},
			Body: []ast.Node{
				&ast.AssignNode{
					Name: "bar",
					Expr: &ast.ValueNode{Type: token.NUMBER},
				},
			},
		},
	}
	a := New()
	node.Accept(a)

	_, ok := a.symTable.Get("foo")
	assert.True(t, ok)

	_, ok = a.symTable.Get("bar")
	assert.True(t, ok)
}

func TestVisitIfNodeBadCondition(t *testing.T) {
	node := &ast.IfNode{
		Condition: &ast.ValueNode{Type: token.NUMBER},
	}
	a := New()
	assert.Panics(t, func() {
		node.Accept(a)
	})
}

func TestVisitWhileNode(t *testing.T) {
	node := &ast.WhileNode{
		Condition: &ast.ValueNode{Type: token.BOOL},
		Body: []ast.Node{
			&ast.AssignNode{
				Name: "foo",
				Expr: &ast.ValueNode{Type: token.NUMBER},
			},
		},
	}
	a := New()
	node.Accept(a)

	_, ok := a.symTable.Get("foo")
	assert.True(t, ok)
}

func TestVisitWhileNodeBadCondition(t *testing.T) {
	node := &ast.WhileNode{
		Condition: &ast.ValueNode{Type: token.NUMBER},
	}
	a := New()
	assert.Panics(t, func() {
		node.Accept(a)
	})
}
