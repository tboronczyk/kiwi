package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scope"
	"github.com/tboronczyk/kiwi/token"
	"github.com/tboronczyk/kiwi/types"
)

func TestEvalValueNode(t *testing.T) {
	nodeData := []struct {
		valVal    string
		valType   token.Token
		expctVal  interface{}
		expctType types.DataType
	}{
		{"42", token.NUMBER, 42.0, types.NUMBER},
		{"foo", token.STRING, "foo", types.STRING},
		{"true", token.BOOL, true, types.BOOL},
	}
	for _, d := range nodeData {
		n := &ast.ValueNode{
			Value: d.valVal,
			Type:  d.valType,
		}
		r := New()
		n.Accept(r)

		e := r.stack.Pop().(scope.Entry)
		assert.Equal(t, d.expctVal, e.Value)
		assert.Equal(t, d.expctType, e.DataType)
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

	e, _ := r.curScope.GetVar("foo")
	assert.Equal(t, "bar", e.Value)
	assert.Equal(t, types.STRING, e.DataType)
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
	r.curScope.SetVar("foo", scope.Entry{
		Value:    42.0,
		DataType: types.NUMBER,
	})
	n.Accept(r)

	e, _ := r.curScope.GetVar("foo")
	assert.Equal(t, 42.0, e.Value)
	assert.Equal(t, types.NUMBER, e.DataType)
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
		termVal   string
		termType  token.Token
		expctVal  interface{}
		expctType types.DataType
	}{
		{token.NOT, "false", token.BOOL, true, types.BOOL},
		{token.ADD, "-42", token.NUMBER, 42.0, types.NUMBER},
		{token.SUBTRACT, "42", token.NUMBER, -42.0, types.NUMBER},
	}
	for _, d := range nodeData {
		n := &ast.UnaryOpNode{
			Op:   d.op,
			Term: &ast.ValueNode{Value: d.termVal, Type: d.termType},
		}
		r := New()
		n.Accept(r)

		e := r.stack.Pop().(scope.Entry)
		assert.Equal(t, d.expctVal, e.Value)
		assert.Equal(t, d.expctType, e.DataType)
	}
}

func TestEvalUnaryOpInvalid(t *testing.T) {
	n := &ast.UnaryOpNode{
		Op:   token.NOT,
		Term: &ast.ValueNode{Value: "foo", Type: token.STRING},
	}
	assert.Panics(t, func() {
		n.Accept(New())
	})
}

func TestCastNode(t *testing.T) {
	nodeData := []struct {
		c         string
		termVal   string
		termType  token.Token
		expctVal  interface{}
		expctType types.DataType
	}{
		{"str", "foo", token.STRING, "foo", types.STRING},
		{"str", "42", token.NUMBER, "42", types.STRING},
		{"str", "true", token.BOOL, "true", types.STRING},
		{"num", "foo", token.STRING, 0.0, types.NUMBER},
		{"num", "42", token.NUMBER, 42.0, types.NUMBER},
		{"num", "true", token.BOOL, 1.0, types.NUMBER},
		{"bool", "foo", token.STRING, true, types.BOOL},
		{"bool", "42", token.NUMBER, true, types.BOOL},
		{"bool", "true", token.BOOL, true, types.BOOL},
		{"bool", "", token.STRING, false, types.BOOL},
		{"bool", "0", token.NUMBER, false, types.BOOL},
		{"bool", "false", token.BOOL, false, types.BOOL},
	}
	for _, d := range nodeData {
		n := &ast.CastNode{
			Cast: d.c,
			Term: &ast.ValueNode{Value: d.termVal, Type: d.termType},
		}
		r := New()
		n.Accept(r)

		e := r.stack.Pop().(scope.Entry)
		assert.Equal(t, d.expctVal, e.Value)
		assert.Equal(t, d.expctType, e.DataType)
	}
}

func TestEvalBinOpNode(t *testing.T) {
	nodeData := []struct {
		op        token.Token
		ex1Val    string
		ex1Type   token.Token
		ex2Val    string
		ex2Type   token.Token
		expctVal  interface{}
		expctType types.DataType
	}{
		{token.ADD, "11", token.NUMBER, "7", token.NUMBER, 18.0, types.NUMBER},
		{token.SUBTRACT, "11", token.NUMBER, "7", token.NUMBER, 4.0, types.NUMBER},
		{token.MULTIPLY, "11", token.NUMBER, "7", token.NUMBER, 77.0, types.NUMBER},
		{token.DIVIDE, "11", token.NUMBER, "7", token.NUMBER, 11 / 7.0, types.NUMBER},
		{token.MODULO, "11", token.NUMBER, "7", token.NUMBER, 4.0, types.NUMBER},
		{token.EQUAL, "11", token.NUMBER, "7", token.NUMBER, false, types.BOOL},
		{token.NOT_EQUAL, "11", token.NUMBER, "7", token.NUMBER, true, types.BOOL},
		{token.LESS, "11", token.NUMBER, "7", token.NUMBER, false, types.BOOL},
		{token.LESS_EQ, "11", token.NUMBER, "7", token.NUMBER, false, types.BOOL},
		{token.GREATER, "11", token.NUMBER, "7", token.NUMBER, true, types.BOOL},
		{token.GREATER_EQ, "11", token.NUMBER, "7", token.NUMBER, true, types.BOOL},
		{token.ADD, "foo", token.STRING, "bar", token.STRING, "foobar", types.STRING},
		{token.EQUAL, "foo", token.STRING, "bar", token.STRING, false, types.BOOL},
		{token.NOT_EQUAL, "foo", token.STRING, "bar", token.STRING, true, types.BOOL},
		{token.AND, "true", token.BOOL, "false", token.BOOL, false, types.BOOL},
		{token.OR, "false", token.BOOL, "true", token.BOOL, true, types.BOOL},
		{token.EQUAL, "true", token.BOOL, "false", token.BOOL, false, types.BOOL},
		{token.NOT_EQUAL, "true", token.BOOL, "false", token.BOOL, true, types.BOOL},
		// short circuit
		{token.AND, "false", token.BOOL, "true", token.BOOL, false, types.BOOL},
		{token.OR, "true", token.BOOL, "false", token.BOOL, true, types.BOOL},
	}
	for _, d := range nodeData {
		n := &ast.BinOpNode{
			Op:    d.op,
			Left:  &ast.ValueNode{Value: d.ex1Val, Type: d.ex1Type},
			Right: &ast.ValueNode{Value: d.ex2Val, Type: d.ex2Type},
		}
		r := New()
		n.Accept(r)

		e := r.stack.Pop().(scope.Entry)
		assert.Equal(t, d.expctVal, e.Value)
		assert.Equal(t, d.expctType, e.DataType)
	}
}

func TestEvalBinOpNodeMismatch(t *testing.T) {
	n := &ast.BinOpNode{
		Op:    token.EQUAL,
		Left:  &ast.ValueNode{Value: "foo", Type: token.STRING},
		Right: &ast.ValueNode{Value: "42", Type: token.NUMBER},
	}
	assert.Panics(t, func() {
		n.Accept(New())
	})
}

func TestEvalBinOpNodeNotPermitted(t *testing.T) {
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
		n := &ast.BinOpNode{
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
		Cond: &ast.BinOpNode{
			Op:    token.LESS,
			Left:  &ast.VariableNode{Name: "foo"},
			Right: &ast.ValueNode{Value: "10", Type: token.NUMBER},
		},
		Body: []ast.Node{
			&ast.AssignNode{
				Name: "foo",
				Expr: &ast.BinOpNode{
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
	r.curScope.SetVar("foo", scope.Entry{Value: 0.0, DataType: types.NUMBER})
	n.Accept(r)

	e, _ := r.curScope.GetVar("foo")
	assert.Equal(t, 10.0, e.Value)
}

func TestEvalWhileNodeNonBool(t *testing.T) {
	n := &ast.WhileNode{
		Cond: &ast.ValueNode{Value: "foo", Type: token.STRING},
		Body: []ast.Node{},
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
			Cond: &ast.ValueNode{Value: d.cVal, Type: token.BOOL},
			Body: []ast.Node{
				&ast.ReturnNode{
					Expr: &ast.ValueNode{Value: "false", Type: token.BOOL},
				},
			},
			Else: &ast.IfNode{
				Cond: &ast.ValueNode{Value: "true", Type: token.BOOL},
				Body: []ast.Node{
					&ast.ReturnNode{
						Expr: &ast.ValueNode{Value: "true", Type: token.BOOL},
					},
				},
			},
		}
		r := New()
		n.Accept(r)
		e := r.stack.Pop().(scope.Entry)
		assert.Equal(t, d.exVal, e.Value)
	}
}

func TestEvalIfNodeNonBool(t *testing.T) {
	n := &ast.IfNode{
		Cond: &ast.ValueNode{Value: "foo", Type: token.STRING},
		Body: []ast.Node{},
	}
	assert.Panics(t, func() {
		n.Accept(New())
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
		Scope: r.curScope,
	}
	r.curScope.SetFunc("foo", scope.Entry{
		Value: f, DataType: types.FUNC,
	})

	n := &ast.FuncCallNode{
		Name: "foo",
		Args: []ast.Node{
			&ast.ValueNode{Value: "bar", Type: token.STRING},
		},
	}
	n.Accept(r)

	e := r.stack.Pop().(scope.Entry)
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
	r.curScope.SetFunc("foo", scope.Entry{
		Value:    &ast.FuncDefNode{Name: "foo", Scope: r.curScope},
		DataType: types.FUNC,
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
