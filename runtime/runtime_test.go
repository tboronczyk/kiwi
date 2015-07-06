package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/token"
)

func TestEvalValueNode(t *testing.T) {
	nodeData := []struct {
		valVal    string
		valType   token.Token
		expctVal  interface{}
		expctType DataType
	}{
		{"42", token.NUMBER, 42.0, NUMBER},
		{"foo", token.STRING, "foo", STRING},
		{"true", token.BOOL, true, BOOL},
	}
	for _, d := range nodeData {
		n := &ast.ValueNode{
			Value: d.valVal,
			Type:  d.valType,
		}
		r := New()
		n.Accept(r)

		e := r.stack.Pop().(ValueEntry)
		assert.Equal(t, d.expctVal, e.Value)
		assert.Equal(t, d.expctType, e.Type)
	}
}

func TestEvalValueNodeInvalid(t *testing.T) {
	n := &ast.ValueNode{Value: "foo", Type: token.UNKNOWN}
	assert.Panics(t, func() {
		n.Accept(New())
	})
}

func TestEvalAssignNode(t *testing.T) {
	n := &ast.AssignNode{
		Name: "foo",
		Expr: &ast.ValueNode{Value: "bar", Type: token.STRING},
	}
	r := New()
	n.Accept(r)

	e, _ := r.varTable.Get("foo")
	assert.Equal(t, "bar", e.(ValueEntry).Value)
	assert.Equal(t, STRING, e.(ValueEntry).Type)
}

func TestEvalAssignNodeBadType(t *testing.T) {
	n := &ast.AssignNode{
		Name: "foo",
		Expr: &ast.ValueNode{Value: "bar", Type: token.STRING},
	}
	r := New()
	n.Accept(r)

	n.Expr = &ast.ValueNode{Value: "42", Type: token.NUMBER}
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalVariableNode(t *testing.T) {
	n := &ast.VariableNode{Name: "foo"}
	r := New()
	r.varTable.Set("foo", ValueEntry{
		Value: 42.0,
		Type:  NUMBER,
	})
	n.Accept(r)

	e, _ := r.varTable.Get("foo")
	assert.Equal(t, 42.0, e.(ValueEntry).Value)
	assert.Equal(t, NUMBER, e.(ValueEntry).Type)
}

func TestEvalVariableNodeNoExist(t *testing.T) {
	n := &ast.VariableNode{Name: "foo"}
	assert.Panics(t, func() {
		n.Accept(New())
	})
}

func TestEvalUnaryOpNode(t *testing.T) {
	nodeData := []struct {
		op        token.Token
		exprVal   string
		exprType  token.Token
		expctVal  interface{}
		expctType DataType
	}{
		{token.NOT, "false", token.BOOL, true, BOOL},
		{token.ADD, "-42", token.NUMBER, 42.0, NUMBER},
		{token.SUBTRACT, "42", token.NUMBER, -42.0, NUMBER},
	}
	for _, d := range nodeData {
		n := &ast.UnaryOpNode{
			Op:   d.op,
			Expr: &ast.ValueNode{Value: d.exprVal, Type: d.exprType},
		}
		r := New()
		n.Accept(r)

		e := r.stack.Pop().(ValueEntry)
		assert.Equal(t, d.expctVal, e.Value)
		assert.Equal(t, d.expctType, e.Type)
	}
}

func TestEvalUnaryOpInvalid(t *testing.T) {
	n := &ast.UnaryOpNode{
		Op:   token.NOT,
		Expr: &ast.ValueNode{Value: "foo", Type: token.STRING},
	}
	assert.Panics(t, func() {
		n.Accept(New())
	})
}

func TestCastNode(t *testing.T) {
	nodeData := []struct {
		c         string
		exprVal   string
		exprType  token.Token
		expctVal  interface{}
		expctType DataType
	}{
		{"str", "foo", token.STRING, "foo", STRING},
		{"str", "42", token.NUMBER, "42", STRING},
		{"str", "true", token.BOOL, "true", STRING},
		{"num", "foo", token.STRING, 0.0, NUMBER},
		{"num", "42", token.NUMBER, 42.0, NUMBER},
		{"num", "true", token.BOOL, 1.0, NUMBER},
		{"bool", "foo", token.STRING, true, BOOL},
		{"bool", "42", token.NUMBER, true, BOOL},
		{"bool", "true", token.BOOL, true, BOOL},
		{"bool", "", token.STRING, false, BOOL},
		{"bool", "0", token.NUMBER, false, BOOL},
		{"bool", "false", token.BOOL, false, BOOL},
	}
	for _, d := range nodeData {
		n := &ast.CastNode{
			Cast: d.c,
			Expr: &ast.ValueNode{Value: d.exprVal, Type: d.exprType},
		}
		r := New()
		n.Accept(r)

		e := r.stack.Pop().(ValueEntry)
		assert.Equal(t, d.expctVal, e.Value)
		assert.Equal(t, d.expctType, e.Type)
	}
}

func TestEvalBinaryOpNode(t *testing.T) {
	nodeData := []struct {
		op        token.Token
		ex1Val    string
		ex1Type   token.Token
		ex2Val    string
		ex2Type   token.Token
		expctVal  interface{}
		expctType DataType
	}{
		{token.ADD, "11", token.NUMBER, "7", token.NUMBER, 18.0, NUMBER},
		{token.SUBTRACT, "11", token.NUMBER, "7", token.NUMBER, 4.0, NUMBER},
		{token.MULTIPLY, "11", token.NUMBER, "7", token.NUMBER, 77.0, NUMBER},
		{token.DIVIDE, "11", token.NUMBER, "7", token.NUMBER, 11 / 7.0, NUMBER},
		{token.MODULO, "11", token.NUMBER, "7", token.NUMBER, 4.0, NUMBER},
		{token.EQUAL, "11", token.NUMBER, "7", token.NUMBER, false, BOOL},
		{token.NOT_EQUAL, "11", token.NUMBER, "7", token.NUMBER, true, BOOL},
		{token.LESS, "11", token.NUMBER, "7", token.NUMBER, false, BOOL},
		{token.LESS_EQ, "11", token.NUMBER, "7", token.NUMBER, false, BOOL},
		{token.GREATER, "11", token.NUMBER, "7", token.NUMBER, true, BOOL},
		{token.GREATER_EQ, "11", token.NUMBER, "7", token.NUMBER, true, BOOL},
		{token.ADD, "foo", token.STRING, "bar", token.STRING, "foobar", STRING},
		{token.EQUAL, "foo", token.STRING, "bar", token.STRING, false, BOOL},
		{token.NOT_EQUAL, "foo", token.STRING, "bar", token.STRING, true, BOOL},
		{token.AND, "true", token.BOOL, "false", token.BOOL, false, BOOL},
		{token.OR, "false", token.BOOL, "true", token.BOOL, true, BOOL},
		{token.EQUAL, "true", token.BOOL, "false", token.BOOL, false, BOOL},
		{token.NOT_EQUAL, "true", token.BOOL, "false", token.BOOL, true, BOOL},
		// short circuit
		{token.AND, "false", token.BOOL, "true", token.BOOL, false, BOOL},
		{token.OR, "true", token.BOOL, "false", token.BOOL, true, BOOL},
	}
	for _, d := range nodeData {
		n := &ast.BinaryOpNode{
			Op:    d.op,
			Left:  &ast.ValueNode{Value: d.ex1Val, Type: d.ex1Type},
			Right: &ast.ValueNode{Value: d.ex2Val, Type: d.ex2Type},
		}
		r := New()
		n.Accept(r)

		e := r.stack.Pop().(ValueEntry)
		assert.Equal(t, d.expctVal, e.Value)
		assert.Equal(t, d.expctType, e.Type)
	}
}

func TestEvalBinaryOpNodeMismatch(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.EQUAL,
		Left:  &ast.ValueNode{Value: "foo", Type: token.STRING},
		Right: &ast.ValueNode{Value: "42", Type: token.NUMBER},
	}
	assert.Panics(t, func() {
		n.Accept(New())
	})
}

func TestEvalBinaryOpNodeNotPermitted(t *testing.T) {
	nodeData := []struct {
		op     token.Token
		exVal  string
		exType token.Token
	}{
		{token.NOT, "42", token.NUMBER},
		{token.GREATER, "foo", token.STRING},
		{token.ADD, "true", token.BOOL},
	}
	for _, d := range nodeData {
		n := &ast.BinaryOpNode{
			Op:    d.op,
			Left:  &ast.ValueNode{Value: d.exVal, Type: d.exType},
			Right: &ast.ValueNode{Value: d.exVal, Type: d.exType},
		}
		assert.Panics(t, func() {
			n.Accept(New())
		})
	}
}

func TestEvalWhileNode(t *testing.T) {
	n := &ast.WhileNode{
		Condition: &ast.BinaryOpNode{
			Op:    token.LESS,
			Left:  &ast.VariableNode{Name: "foo"},
			Right: &ast.ValueNode{Value: "10", Type: token.NUMBER},
		},
		Body: []ast.Node{
			&ast.AssignNode{
				Name: "foo",
				Expr: &ast.BinaryOpNode{
					Op:   token.ADD,
					Left: &ast.VariableNode{Name: "foo"},
					Right: &ast.ValueNode{
						Value: "1",
						Type:  token.NUMBER,
					},
				},
			},
		},
	}
	r := New()
	r.varTable.Set("foo", ValueEntry{Value: 0.0, Type: NUMBER})
	n.Accept(r)

	e, _ := r.varTable.Get("foo")
	assert.Equal(t, 10.0, e.(ValueEntry).Value)
}

func TestEvalWhileNodeNonBool(t *testing.T) {
	n := &ast.WhileNode{
		Condition: &ast.ValueNode{Value: "foo", Type: token.STRING},
		Body:      []ast.Node{},
	}
	assert.Panics(t, func() {
		n.Accept(New())
	})
}

func TestEvalIfNode(t *testing.T) {
	nodeData := []struct {
		cVal  string
		exVal bool
	}{
		{"true", false},
		{"false", true},
	}
	for _, d := range nodeData {
		n := &ast.IfNode{
			Condition: &ast.ValueNode{Value: d.cVal, Type: token.BOOL},
			Body: []ast.Node{
				&ast.ReturnNode{
					Expr: &ast.ValueNode{Value: "false", Type: token.BOOL},
				},
			},
			Else: &ast.IfNode{
				Condition: &ast.ValueNode{Value: "true", Type: token.BOOL},
				Body: []ast.Node{
					&ast.ReturnNode{
						Expr: &ast.ValueNode{Value: "true", Type: token.BOOL},
					},
				},
			},
		}
		r := New()
		n.Accept(r)
		e := r.stack.Pop().(ValueEntry)
		assert.Equal(t, d.exVal, e.Value)
	}
}

func TestEvalIfNodeNonBool(t *testing.T) {
	n := &ast.IfNode{
		Condition: &ast.ValueNode{Value: "foo", Type: token.STRING},
		Body:      []ast.Node{},
	}
	assert.Panics(t, func() {
		n.Accept(New())
	})
}

func TestEvalFuncDefNode(t *testing.T) {
	n := &ast.FuncDefNode{
		Name: "foo",
		Args: []string{},
		Body: []ast.Node{},
	}
	r := New()
	n.Accept(r)

	e, _ := r.funcTable.Get("foo")
	assert.Equal(t, "foo", e.(ValueEntry).Value.(*ast.FuncDefNode).Name)
}

func TestEvalFuncDefNodeExists(t *testing.T) {
	r := New()
	r.funcTable.Set("foo", &ast.FuncDefNode{})
	n := &ast.FuncDefNode{
		Name: "foo",
		Args: []string{},
		Body: []ast.Node{},
	}
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalFuncCallNode(t *testing.T) {
	r := New()
	f := &ast.FuncDefNode{
		Name: "foo",
		Args: []string{"bar"},
		Body: []ast.Node{
			&ast.ReturnNode{
				Expr: &ast.VariableNode{Name: "bar"},
			},
		},
	}
	f.Accept(r)

	n := &ast.FuncCallNode{
		Name: "foo",
		Args: []ast.Node{
			&ast.ValueNode{Value: "bar", Type: token.STRING},
		},
	}
	n.Accept(r)

	e := r.stack.Pop().(ValueEntry)
	assert.Equal(t, "bar", e.Value)
}

func TestEvalFuncCallNodeBuiltin(t *testing.T) {
	n := &ast.FuncCallNode{
		Name: "strlen",
		Args: []ast.Node{
			&ast.ValueNode{Value: "foo", Type: token.STRING},
		},
	}
	r := New()
	assert.NotPanics(t, func() {
		n.Accept(r)
	})
}

func TestEvalFuncCallNodeNotDefined(t *testing.T) {
	n := &ast.FuncCallNode{
		Name: "foo",
		Args: []ast.Node{},
	}
	assert.Panics(t, func() {
		n.Accept(New())
	})
}

func TestEvalFuncCallNodeBadCount(t *testing.T) {
	r := New()
	r.funcTable.Set("foo", ValueEntry{
		Value: &ast.FuncDefNode{},
		Type:  FUNC,
	})
	n := &ast.FuncCallNode{
		Name: "foo",
		Args: []ast.Node{
			&ast.ValueNode{Value: "bar", Type: token.STRING},
		},
	}
	assert.Panics(t, func() {
		n.Accept(r)
	})
}
