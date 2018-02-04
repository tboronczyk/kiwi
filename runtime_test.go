package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuntime(t *testing.T) {
	t.Parallel()

	t.Run("Test AddNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate AddNode with numbers", func(t *testing.T) {
			n := &AstAddNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstNumberNode{Value: 73},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, 115.0, e.Value)
			assert.Equal(t, NUMBER, e.DataType)
		})

		t.Run("Evaluate AddNode with strings", func(t *testing.T) {
			n := &AstAddNode{
				Left:  &AstStringNode{Value: "foo"},
				Right: &AstStringNode{Value: "bar"},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, "foobar", e.Value)
			assert.Equal(t, STRING, e.DataType)
		})

		t.Run("Evaluate AddNode with type error", func(t *testing.T) {
			n := &AstAddNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test AndNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate AndNode", func(t *testing.T) {
			n := &AstAndNode{
				Left:  &AstBoolNode{Value: true},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})

		t.Run("Evaluate AndNode short circuited", func(t *testing.T) {
			n := &AstAndNode{
				Left: &AstBoolNode{Value: false},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, false, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})

		t.Run("Evaluate AndNode with type error", func(t *testing.T) {
			n := &AstAndNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test AssignNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate AssignNode", func(t *testing.T) {
			n := &AstAssignNode{
				Name: "foo",
				Expr: &AstStringNode{Value: "bar"},
			}
			r := NewRuntime()
			n.Accept(r)

			e, _ := r.curScope.GetVar("foo")
			assert.Equal(t, "bar", e.Value)
			assert.Equal(t, STRING, e.DataType)
		})

		t.Run("Evaluate AssignNode with type error", func(t *testing.T) {
			n := &AstAssignNode{
				Name: "foo",
				Expr: &AstStringNode{Value: "bar"},
			}
			r := NewRuntime()
			n.Accept(r)

			n.Expr = &AstNumberNode{Value: 42}
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
				term      AstNode
				expctVal  interface{}
				expctType DataType
			}{
				{"str", &AstStringNode{Value: "foo"}, "foo", STRING},
				{"str", &AstNumberNode{Value: 42}, "42", STRING},
				{"str", &AstBoolNode{Value: true}, "true", STRING},
				{"num", &AstStringNode{Value: "foo"}, 0.0, NUMBER},
				{"num", &AstNumberNode{Value: 42}, 42.0, NUMBER},
				{"num", &AstBoolNode{Value: true}, 1.0, NUMBER},
				{"bool", &AstStringNode{Value: "foo"}, true, BOOL},
				{"bool", &AstNumberNode{Value: 42}, true, BOOL},
				{"bool", &AstBoolNode{Value: true}, true, BOOL},
				{"bool", &AstStringNode{Value: ""}, false, BOOL},
				{"bool", &AstNumberNode{Value: 0}, false, BOOL},
				{"bool", &AstBoolNode{Value: false}, false, BOOL},
			}
			for _, d := range nodeData {
				n := &AstCastNode{
					Cast: d.cast,
					Term: d.term,
				}
				r := NewRuntime()
				n.Accept(r)

				e := r.stack.Pop().(Entry)
				assert.Equal(t, d.expctVal, e.Value)
				assert.Equal(t, d.expctType, e.DataType)
			}
		})
	})

	t.Run("Test DivideNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate DivideNode", func(t *testing.T) {
			n := &AstDivideNode{
				Left:  &AstNumberNode{Value: 110},
				Right: &AstNumberNode{Value: 4},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, 27.5, e.Value)
			assert.Equal(t, NUMBER, e.DataType)
		})

		t.Run("Evaluate DivideNode with type error", func(t *testing.T) {
			n := &AstDivideNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test EqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate EqualNode", func(t *testing.T) {
			n := &AstEqualNode{
				Left:  &AstBoolNode{Value: true},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})

		t.Run("Evaluate EqualNode with type error", func(t *testing.T) {
			n := &AstEqualNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test GreaterEqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate GreaterEqualNode", func(t *testing.T) {
			n := &AstGreaterEqualNode{
				Left:  &AstNumberNode{Value: 1984},
				Right: &AstNumberNode{Value: 1776},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})

		t.Run("Evaluate GreaterEqualNode with type error", func(t *testing.T) {
			n := &AstGreaterEqualNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test GreaterNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate GreaterNode", func(t *testing.T) {
			n := &AstGreaterNode{
				Left:  &AstNumberNode{Value: 1984},
				Right: &AstNumberNode{Value: 1776},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})

		t.Run("Evaluate GreaterNode with type error", func(t *testing.T) {
			n := &AstGreaterNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test LessEqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate LessEqualNode", func(t *testing.T) {
			n := &AstLessEqualNode{
				Left:  &AstNumberNode{Value: 1984},
				Right: &AstNumberNode{Value: 1776},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, false, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})

		t.Run("Evaluate LessEqualNode with type error", func(t *testing.T) {
			n := &AstLessEqualNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test LessNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate LessNode", func(t *testing.T) {
			n := &AstLessNode{
				Left:  &AstNumberNode{Value: 1984},
				Right: &AstNumberNode{Value: 1776},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, false, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})

		t.Run("Evaluate LessNode with type error", func(t *testing.T) {
			n := &AstLessNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test ModuloNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate ModuloNode", func(t *testing.T) {
			n := &AstModuloNode{
				Left:  &AstNumberNode{Value: 73},
				Right: &AstNumberNode{Value: 42},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, 31.0, e.Value)
			assert.Equal(t, NUMBER, e.DataType)
		})

		t.Run("Evaluate ModuloNode with type error", func(t *testing.T) {
			n := &AstModuloNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test MultiplyNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate MultiplyNode", func(t *testing.T) {
			n := &AstMultiplyNode{
				Left:  &AstNumberNode{Value: 21},
				Right: &AstNumberNode{Value: 2},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, 42.0, e.Value)
			assert.Equal(t, NUMBER, e.DataType)
		})

		t.Run("Evaluate MultiplyNode with type error", func(t *testing.T) {
			n := &AstMultiplyNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test NegativeNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate NegativeNode", func(t *testing.T) {
			n := &AstNegativeNode{
				Term: &AstNumberNode{Value: 42},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, -42.0, e.Value)
			assert.Equal(t, NUMBER, e.DataType)
		})

		t.Run("Evaluate NegativeNode with type error", func(t *testing.T) {
			n := &AstNegativeNode{
				Term: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test NotEqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate NotEqualNode", func(t *testing.T) {
			n := &AstNotEqualNode{
				Left:  &AstBoolNode{Value: true},
				Right: &AstBoolNode{Value: false},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})

		t.Run("Evaluate NotEqualNode with type error", func(t *testing.T) {
			n := &AstNotEqualNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test NotNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate NotNode", func(t *testing.T) {
			n := &AstNotNode{
				Term: &AstBoolNode{Value: false},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})

		t.Run("Evaluate NotNode with type error", func(t *testing.T) {
			n := &AstNotNode{
				Term: &AstNumberNode{Value: 42},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test OrNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate OrNode", func(t *testing.T) {
			n := &AstOrNode{
				Left:  &AstBoolNode{Value: false},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})

		t.Run("Evaluate OrNode short circuited", func(t *testing.T) {
			n := &AstOrNode{
				Left: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})

		t.Run("Evaluate OrNode with type error", func(t *testing.T) {
			n := &AstOrNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test PositiveNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate PositiveNode", func(t *testing.T) {
			n := &AstPositiveNode{
				Term: &AstNumberNode{Value: -42},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, 42.0, e.Value)
			assert.Equal(t, NUMBER, e.DataType)
		})

		t.Run("Evaluate PositiveNode with type error", func(t *testing.T) {
			n := &AstPositiveNode{
				Term: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test ReturnNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate ReturnNode", func(t *testing.T) {
			n := &AstReturnNode{
				Expr: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, BOOL, e.DataType)
		})
	})

	t.Run("Test SubtractNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate SubtractNode", func(t *testing.T) {
			n := &AstSubtractNode{
				Left:  &AstNumberNode{Value: 73},
				Right: &AstNumberNode{Value: 42},
			}
			r := NewRuntime()
			n.Accept(r)
			e := r.stack.Pop().(Entry)
			assert.Equal(t, 31.0, e.Value)
			assert.Equal(t, NUMBER, e.DataType)
		})

		t.Run("Eval SubtractNode with type error", func(t *testing.T) {
			n := &AstSubtractNode{
				Left:  &AstNumberNode{Value: 42},
				Right: &AstBoolNode{Value: true},
			}
			r := NewRuntime()
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test VariableNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate VariableNode", func(t *testing.T) {
			n := &AstVariableNode{
				Name: "foo",
			}
			r := NewRuntime()
			r.curScope.SetVar("foo", Entry{
				Value:    "bar",
				DataType: STRING,
			})
			n.Accept(r)

			e, _ := r.curScope.GetVar("foo")
			assert.Equal(t, "bar", e.Value)
			assert.Equal(t, STRING, e.DataType)
		})

		t.Run("Evaluate VariableNode undefined", func(t *testing.T) {
			n := &AstVariableNode{
				Name: "foo",
			}
			r := NewRuntime()

			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})
}
