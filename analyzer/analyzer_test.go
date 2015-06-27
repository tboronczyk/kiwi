package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/symtable"
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
	a.symTable.Set("foo", symtable.VARTYPE, STRING)
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

	expr, ok := a.symTable.Get("foo", symtable.VARTYPE)
	assert.Equal(t, STRING, expr)
	assert.True(t, ok)
}

func TestVisitFuncDefNodeNoReturn(t *testing.T) {
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

	expr, ok := a.symTable.Get("foo", symtable.FUNCTYPE)
	assert.Equal(t, expr, UNKNOWN)
	assert.True(t, ok)

	value, _ := node.SymTable.T[symtable.VARTYPE]["bar"]
	assert.Equal(t, value, UNKNOWN)

	value, _ = node.SymTable.T[symtable.VARTYPE]["baz"]
	assert.Equal(t, value, UNKNOWN)

	value, _ = node.SymTable.T[symtable.VARTYPE]["qux"]
	assert.Equal(t, value, STRING)

	_, ok = node.SymTable.T[symtable.VARTYPE]["norf"]
	assert.False(t, ok)
}

func TestVisitFuncDefNodeSingleReturn(t *testing.T) {
	node := &ast.FuncDefNode{
		Name: "foo",
		Args: []string{},
		Body: []ast.Node{
			&ast.ReturnNode{
				Expr: &ast.ValueNode{Type: token.STRING},
			},
		},
	}
	a := New()
	node.Accept(a)

	expr, ok := a.symTable.Get("foo", symtable.FUNCTYPE)
	assert.Equal(t, expr, STRING)
	assert.True(t, ok)
}

func TestVisitFuncDefNodeMultipleReturn(t *testing.T) {
	node := &ast.FuncDefNode{
		Name: "foo",
		Args: []string{},
		Body: []ast.Node{
			&ast.ReturnNode{
				Expr: &ast.ValueNode{Type: token.STRING},
			},
			&ast.ReturnNode{
				Expr: &ast.ValueNode{Type: token.STRING},
			},
		},
	}
	a := New()
	node.Accept(a)

	expr, ok := a.symTable.Get("foo", symtable.FUNCTYPE)
	assert.Equal(t, expr, STRING)
	assert.True(t, ok)
}

func TestVisitFuncDefNodeMultipleReturnBad(t *testing.T) {
	node := &ast.FuncDefNode{
		Name: "foo",
		Args: []string{},
		Body: []ast.Node{
			&ast.ReturnNode{
				Expr: &ast.ValueNode{Type: token.STRING},
			},
			&ast.ReturnNode{
				Expr: &ast.ValueNode{Type: token.NUMBER},
			},
		},
	}
	a := New()
	assert.Panics(t, func() {
		node.Accept(a)
	})
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
	a.symTable.Set("foo", symtable.VARTYPE, UNKNOWN)
	a.symTable.Set("bar", symtable.VARTYPE, UNKNOWN)
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
	a.symTable.Set("foo", symtable.VARTYPE, STRING)
	a.symTable.Set("bar", symtable.VARTYPE, NUMBER)

	assert.Panics(t, func() {
		node.Accept(a)
	})
}

func TestVisitFuncCallNode(t *testing.T) {
	node := &ast.FuncCallNode{Name: "foo"}
	a := New()
	a.symTable.Set("foo", symtable.FUNCTYPE, UNKNOWN)
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

	_, ok := a.symTable.Get("foo", symtable.VARTYPE)
	assert.True(t, ok)

	_, ok = a.symTable.Get("bar", symtable.VARTYPE)
	assert.True(t, ok)
}

func TestVisitIfNodeInFuncDefReturn(t *testing.T) {
	node := &ast.FuncDefNode{
		Name: "foo",
		Args: []string{},
		Body: []ast.Node{
			&ast.IfNode{
				Condition: &ast.ValueNode{Type: token.BOOL},
				Body: []ast.Node{
					&ast.ReturnNode{
						Expr: &ast.ValueNode{Type: token.NUMBER},
					},
				},
			},
		},
	}
	a := New()
	node.Accept(a)
	dtype, _ := a.symTable.Get("foo", symtable.FUNCTYPE)
	assert.Equal(t, NUMBER, dtype)
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

	_, ok := a.symTable.Get("foo", symtable.VARTYPE)
	assert.True(t, ok)
}

func TestVisitWhileNodeInFuncDefReturn(t *testing.T) {
	node := &ast.FuncDefNode{
		Name: "foo",
		Args: []string{},
		Body: []ast.Node{
			&ast.WhileNode{
				Condition: &ast.ValueNode{Type: token.BOOL},
				Body: []ast.Node{
					&ast.ReturnNode{
						Expr: &ast.ValueNode{Type: token.NUMBER},
					},
				},
			},
		},
	}
	a := New()
	node.Accept(a)

	dtype, _ := a.symTable.Get("foo", symtable.FUNCTYPE)
	assert.Equal(t, NUMBER, dtype)
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
