package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scope"
	"github.com/tboronczyk/kiwi/types"
)

func TestEvalAddNodeNumber(t *testing.T) {
	n := &ast.AddNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.NumberNode{Value: 73},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, 115.0, e.Value)
	assert.Equal(t, types.NUMBER, e.DataType)
}

func TestEvalAddNodeString(t *testing.T) {
	n := &ast.AddNode{
		Left:  &ast.StringNode{Value: "foo"},
		Right: &ast.StringNode{Value: "bar"},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, "foobar", e.Value)
	assert.Equal(t, types.STRING, e.DataType)
}

func TestEvalAddNodeTypeError(t *testing.T) {
	n := &ast.AddNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalAndNode(t *testing.T) {
	n := &ast.AndNode{
		Left:  &ast.BoolNode{Value: true},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, true, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalAndNodeShortCircuit(t *testing.T) {
	n := &ast.AndNode{
		Left: &ast.BoolNode{Value: false},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, false, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalAndNodeTypeError(t *testing.T) {
	n := &ast.AndNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalAssignNode(t *testing.T) {
	n := &ast.AssignNode{
		Name: "foo",
		Expr: &ast.StringNode{Value: "bar"},
	}
	r := New()
	n.Accept(r)

	e, _ := r.curScope.GetVar("foo")
	assert.Equal(t, "bar", e.Value)
	assert.Equal(t, types.STRING, e.DataType)
}

func TestEvalAssignNodeTypeError(t *testing.T) {
	n := &ast.AssignNode{
		Name: "foo",
		Expr: &ast.StringNode{Value: "bar"},
	}
	r := New()
	n.Accept(r)

	n.Expr = &ast.NumberNode{Value: 42}
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestCastNode(t *testing.T) {
	nodeData := []struct {
		cast      string
		term      ast.Node
		expctVal  interface{}
		expctType types.DataType
	}{
		{"str", &ast.StringNode{Value: "foo"}, "foo", types.STRING},
		{"str", &ast.NumberNode{Value: 42}, "42", types.STRING},
		{"str", &ast.BoolNode{Value: true}, "true", types.STRING},
		{"num", &ast.StringNode{Value: "foo"}, 0.0, types.NUMBER},
		{"num", &ast.NumberNode{Value: 42}, 42.0, types.NUMBER},
		{"num", &ast.BoolNode{Value: true}, 1.0, types.NUMBER},
		{"bool", &ast.StringNode{Value: "foo"}, true, types.BOOL},
		{"bool", &ast.NumberNode{Value: 42}, true, types.BOOL},
		{"bool", &ast.BoolNode{Value: true}, true, types.BOOL},
		{"bool", &ast.StringNode{Value: ""}, false, types.BOOL},
		{"bool", &ast.NumberNode{Value: 0}, false, types.BOOL},
		{"bool", &ast.BoolNode{Value: false}, false, types.BOOL},
	}
	for _, d := range nodeData {
		n := &ast.CastNode{
			Cast: d.cast,
			Term: d.term,
		}
		r := New()
		n.Accept(r)

		e := r.stack.Pop().(scope.Entry)
		assert.Equal(t, d.expctVal, e.Value)
		assert.Equal(t, d.expctType, e.DataType)
	}
}

func TestEvalDivideNode(t *testing.T) {
	n := &ast.DivideNode{
		Left:  &ast.NumberNode{Value: 110},
		Right: &ast.NumberNode{Value: 4},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, 27.5, e.Value)
	assert.Equal(t, types.NUMBER, e.DataType)
}

func TestEvalDivideNodeTypeError(t *testing.T) {
	n := &ast.DivideNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalEqualNode(t *testing.T) {
	n := &ast.EqualNode{
		Left:  &ast.BoolNode{Value: true},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, true, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalEqualNodeTypeError(t *testing.T) {
	n := &ast.EqualNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalGreaterEqualNode(t *testing.T) {
	n := &ast.GreaterEqualNode{
		Left:  &ast.NumberNode{Value: 1984},
		Right: &ast.NumberNode{Value: 1776},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, true, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalGreaterEqualNodeTypeError(t *testing.T) {
	n := &ast.GreaterEqualNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalGreaterNode(t *testing.T) {
	n := &ast.GreaterNode{
		Left:  &ast.NumberNode{Value: 1984},
		Right: &ast.NumberNode{Value: 1776},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, true, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalGreaterNodeTypeError(t *testing.T) {
	n := &ast.GreaterNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalLessEqualNode(t *testing.T) {
	n := &ast.LessEqualNode{
		Left:  &ast.NumberNode{Value: 1984},
		Right: &ast.NumberNode{Value: 1776},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, false, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalLessEqualNodeTypeError(t *testing.T) {
	n := &ast.LessEqualNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalLessNode(t *testing.T) {
	n := &ast.LessNode{
		Left:  &ast.NumberNode{Value: 1984},
		Right: &ast.NumberNode{Value: 1776},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, false, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalLessNodeTypeError(t *testing.T) {
	n := &ast.LessNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalModuloNode(t *testing.T) {
	n := &ast.ModuloNode{
		Left:  &ast.NumberNode{Value: 73},
		Right: &ast.NumberNode{Value: 42},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, 31.0, e.Value)
	assert.Equal(t, types.NUMBER, e.DataType)
}

func TestEvalModuloNodeTypeError(t *testing.T) {
	n := &ast.ModuloNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalMultiplyNode(t *testing.T) {
	n := &ast.MultiplyNode{
		Left:  &ast.NumberNode{Value: 21},
		Right: &ast.NumberNode{Value: 2},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, 42.0, e.Value)
	assert.Equal(t, types.NUMBER, e.DataType)
}

func TestEvalMultiplyNodeTypeError(t *testing.T) {
	n := &ast.MultiplyNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalNegativeNode(t *testing.T) {
	n := &ast.NegativeNode{
		Term: &ast.NumberNode{Value: 42},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, -42.0, e.Value)
	assert.Equal(t, types.NUMBER, e.DataType)
}

func TestEvalNegativeNodeTypeError(t *testing.T) {
	n := &ast.NegativeNode{
		Term: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalNotEqualNode(t *testing.T) {
	n := &ast.NotEqualNode{
		Left:  &ast.BoolNode{Value: true},
		Right: &ast.BoolNode{Value: false},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, true, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalNotEqualNodeTypeError(t *testing.T) {
	n := &ast.NotEqualNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalNotNode(t *testing.T) {
	n := &ast.NotNode{
		Term: &ast.BoolNode{Value: false},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, true, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalNotNodeTypeError(t *testing.T) {
	n := &ast.NotNode{
		Term: &ast.NumberNode{Value: 42},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalOrNode(t *testing.T) {
	n := &ast.OrNode{
		Left:  &ast.BoolNode{Value: false},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, true, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalOrNodeShortCircuit(t *testing.T) {
	n := &ast.OrNode{
		Left: &ast.BoolNode{Value: true},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, true, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalOrNodeTypeError(t *testing.T) {
	n := &ast.OrNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalPositiveNode(t *testing.T) {
	n := &ast.PositiveNode{
		Term: &ast.NumberNode{Value: -42},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, 42.0, e.Value)
	assert.Equal(t, types.NUMBER, e.DataType)
}

func TestEvalPositiveNodeTypeError(t *testing.T) {
	n := &ast.PositiveNode{
		Term: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalReturnNode(t *testing.T) {
	n := &ast.ReturnNode{
		Expr: &ast.BoolNode{Value: true},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, true, e.Value)
	assert.Equal(t, types.BOOL, e.DataType)
}

func TestEvalSubtractNode(t *testing.T) {
	n := &ast.SubtractNode{
		Left:  &ast.NumberNode{Value: 73},
		Right: &ast.NumberNode{Value: 42},
	}
	r := New()
	n.Accept(r)
	e := r.stack.Pop().(scope.Entry)
	assert.Equal(t, 31.0, e.Value)
	assert.Equal(t, types.NUMBER, e.DataType)
}

func TestEvalSubtractNodeTypeError(t *testing.T) {
	n := &ast.SubtractNode{
		Left:  &ast.NumberNode{Value: 42},
		Right: &ast.BoolNode{Value: true},
	}
	r := New()
	assert.Panics(t, func() {
		n.Accept(r)
	})
}

func TestEvalVariableNode(t *testing.T) {
	n := &ast.VariableNode{
		Name: "foo",
	}
	r := New()
	r.curScope.SetVar("foo", scope.Entry{
		Value:    "bar",
		DataType: types.STRING,
	})
	n.Accept(r)

	e, _ := r.curScope.GetVar("foo")
	assert.Equal(t, "bar", e.Value)
	assert.Equal(t, types.STRING, e.DataType)
}

func TestEvalVariableNodeNotDefined(t *testing.T) {
	n := &ast.VariableNode{
		Name: "foo",
	}
	r := New()

	assert.Panics(t, func() {
		n.Accept(r)
	})
}
