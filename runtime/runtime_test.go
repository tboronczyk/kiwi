package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scope"
	"github.com/tboronczyk/kiwi/types"
)

func TestRuntime(t *testing.T) {
	t.Parallel()

	t.Run("Test AddNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate AddNode with numbers", func(t *testing.T) {
			n := &ast.AddNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.NumberNode{Value: 73},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, 115.0, e.Value)
			assert.Equal(t, types.NUMBER, e.DataType)
		})

		t.Run("Evaluate AddNode with strings", func(t *testing.T) {
			n := &ast.AddNode{
				Left:  &ast.StringNode{Value: "foo"},
				Right: &ast.StringNode{Value: "bar"},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, "foobar", e.Value)
			assert.Equal(t, types.STRING, e.DataType)
		})

		t.Run("Evaluate AddNode with type error", func(t *testing.T) {
			n := &ast.AddNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test AndNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate AndNode", func(t *testing.T) {
			n := &ast.AndNode{
				Left:  &ast.BoolNode{Value: true},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})

		t.Run("Evaluate AndNode short circuited", func(t *testing.T) {
			n := &ast.AndNode{
				Left: &ast.BoolNode{Value: false},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, false, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})

		t.Run("Evaluate AndNode with type error", func(t *testing.T) {
			n := &ast.AndNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test AssignNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate AssignNode", func(t *testing.T) {
			n := &ast.AssignNode{
				Name: "foo",
				Expr: &ast.StringNode{Value: "bar"},
			}
			r := New()
			n.Accept(r)

			e, _ := r.curScope.GetVar("foo")
			assert.Equal(t, "bar", e.Value)
			assert.Equal(t, types.STRING, e.DataType)
		})

		t.Run("Evaluate AssignNode with type error", func(t *testing.T) {
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
		})
	})

	t.Run("Test CastNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate CastNode", func(t *testing.T) {
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
		})
	})

	t.Run("Test DivideNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate DivideNode", func(t *testing.T) {
			n := &ast.DivideNode{
				Left:  &ast.NumberNode{Value: 110},
				Right: &ast.NumberNode{Value: 4},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, 27.5, e.Value)
			assert.Equal(t, types.NUMBER, e.DataType)
		})

		t.Run("Evaluate DivideNode with type error", func(t *testing.T) {
			n := &ast.DivideNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test EqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate EqualNode", func(t *testing.T) {
			n := &ast.EqualNode{
				Left:  &ast.BoolNode{Value: true},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})

		t.Run("Evaluate EqualNode with type error", func(t *testing.T) {
			n := &ast.EqualNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test GreaterEqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate GreaterEqualNode", func(t *testing.T) {
			n := &ast.GreaterEqualNode{
				Left:  &ast.NumberNode{Value: 1984},
				Right: &ast.NumberNode{Value: 1776},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})

		t.Run("Evaluate GreaterEqualNode with type error", func(t *testing.T) {
			n := &ast.GreaterEqualNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test GreaterNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate GreaterNode", func(t *testing.T) {
			n := &ast.GreaterNode{
				Left:  &ast.NumberNode{Value: 1984},
				Right: &ast.NumberNode{Value: 1776},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})

		t.Run("Evaluate GreaterNode with type error", func(t *testing.T) {
			n := &ast.GreaterNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test LessEqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate LessEqualNode", func(t *testing.T) {
			n := &ast.LessEqualNode{
				Left:  &ast.NumberNode{Value: 1984},
				Right: &ast.NumberNode{Value: 1776},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, false, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})

		t.Run("Evaluate LessEqualNode with type error", func(t *testing.T) {
			n := &ast.LessEqualNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test LessNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate LessNode", func(t *testing.T) {
			n := &ast.LessNode{
				Left:  &ast.NumberNode{Value: 1984},
				Right: &ast.NumberNode{Value: 1776},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, false, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})

		t.Run("Evaluate LessNode with type error", func(t *testing.T) {
			n := &ast.LessNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test ModuloNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate ModuloNode", func(t *testing.T) {
			n := &ast.ModuloNode{
				Left:  &ast.NumberNode{Value: 73},
				Right: &ast.NumberNode{Value: 42},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, 31.0, e.Value)
			assert.Equal(t, types.NUMBER, e.DataType)
		})

		t.Run("Evaluate ModuloNode with type error", func(t *testing.T) {
			n := &ast.ModuloNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test MultiplyNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate MultiplyNode", func(t *testing.T) {
			n := &ast.MultiplyNode{
				Left:  &ast.NumberNode{Value: 21},
				Right: &ast.NumberNode{Value: 2},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, 42.0, e.Value)
			assert.Equal(t, types.NUMBER, e.DataType)
		})

		t.Run("Evaluate MultiplyNode with type error", func(t *testing.T) {
			n := &ast.MultiplyNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test NegativeNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate NegativeNode", func(t *testing.T) {
			n := &ast.NegativeNode{
				Term: &ast.NumberNode{Value: 42},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, -42.0, e.Value)
			assert.Equal(t, types.NUMBER, e.DataType)
		})

		t.Run("Evaluate NegativeNode with type error", func(t *testing.T) {
			n := &ast.NegativeNode{
				Term: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test NotEqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate NotEqualNode", func(t *testing.T) {
			n := &ast.NotEqualNode{
				Left:  &ast.BoolNode{Value: true},
				Right: &ast.BoolNode{Value: false},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})

		t.Run("Evaluate NotEqualNode with type error", func(t *testing.T) {
			n := &ast.NotEqualNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test NotNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate NotNode", func(t *testing.T) {
			n := &ast.NotNode{
				Term: &ast.BoolNode{Value: false},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})

		t.Run("Evaluate NotNode with type error", func(t *testing.T) {
			n := &ast.NotNode{
				Term: &ast.NumberNode{Value: 42},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test OrNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate OrNode", func(t *testing.T) {
			n := &ast.OrNode{
				Left:  &ast.BoolNode{Value: false},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})

		t.Run("Evaluate OrNode short circuited", func(t *testing.T) {
			n := &ast.OrNode{
				Left: &ast.BoolNode{Value: true},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})

		t.Run("Evaluate OrNode with type error", func(t *testing.T) {
			n := &ast.OrNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test PositiveNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate PositiveNode", func(t *testing.T) {
			n := &ast.PositiveNode{
				Term: &ast.NumberNode{Value: -42},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, 42.0, e.Value)
			assert.Equal(t, types.NUMBER, e.DataType)
		})

		t.Run("Evaluate PositiveNode with type error", func(t *testing.T) {
			n := &ast.PositiveNode{
				Term: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test ReturnNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate ReturnNode", func(t *testing.T) {
			n := &ast.ReturnNode{
				Expr: &ast.BoolNode{Value: true},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, types.BOOL, e.DataType)
		})
	})

	t.Run("Test SubtractNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate SubtractNode", func(t *testing.T) {
			n := &ast.SubtractNode{
				Left:  &ast.NumberNode{Value: 73},
				Right: &ast.NumberNode{Value: 42},
			}
			r := New()
			n.Accept(r)
			e := r.stack.Pop().(scope.Entry)
			assert.Equal(t, 31.0, e.Value)
			assert.Equal(t, types.NUMBER, e.DataType)
		})

		t.Run("Eval SubtractNode with type error", func(t *testing.T) {
			n := &ast.SubtractNode{
				Left:  &ast.NumberNode{Value: 42},
				Right: &ast.BoolNode{Value: true},
			}
			r := New()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test VariableNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate VariableNode", func(t *testing.T) {
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
		})

		t.Run("Evaluate VariableNode undefined", func(t *testing.T) {
			n := &ast.VariableNode{
				Name: "foo",
			}
			r := New()

			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})
}
