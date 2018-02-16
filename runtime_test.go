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
				Left:  &AstNumberNode{42},
				Right: &AstNumberNode{73},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, 115.0, e.Value)
			assert.Equal(t, TypNumber, e.DataType)
		})

		t.Run("Evaluate AddNode with strings", func(t *testing.T) {
			n := &AstAddNode{
				Left:  &AstStringNode{"foo"},
				Right: &AstStringNode{"bar"},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, "foobar", e.Value)
			assert.Equal(t, TypString, e.DataType)
		})

		t.Run("Evaluate AddNode with type error", func(t *testing.T) {
			n := &AstAddNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test AndNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate AndNode", func(t *testing.T) {
			n := &AstAndNode{
				Left:  &AstBoolNode{true},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})

		t.Run("Evaluate AndNode short circuited", func(t *testing.T) {
			n := &AstAndNode{
				Left: &AstBoolNode{false},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, false, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})

		t.Run("Evaluate AndNode with type error", func(t *testing.T) {
			n := &AstAndNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
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
				Expr: &AstStringNode{"bar"},
			}
			r := NewRuntime(nil)
			n.Accept(r)

			e, _ := r.currScope.GetVar("foo")
			assert.Equal(t, "bar", e.Value)
			assert.Equal(t, TypString, e.DataType)
		})

		t.Run("Evaluate AssignNode with type error", func(t *testing.T) {
			n := &AstAssignNode{
				Name: "foo",
				Expr: &AstStringNode{"bar"},
			}
			r := NewRuntime(nil)
			n.Accept(r)

			n.Expr = &AstNumberNode{42}
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
				{"str", &AstStringNode{"foo"}, "foo", TypString},
				{"str", &AstNumberNode{42}, "42", TypString},
				{"str", &AstBoolNode{true}, "true", TypString},
				{"num", &AstStringNode{"foo"}, 0.0, TypNumber},
				{"num", &AstNumberNode{42}, 42.0, TypNumber},
				{"num", &AstBoolNode{true}, 1.0, TypNumber},
				{"bool", &AstStringNode{"foo"}, true, TypBool},
				{"bool", &AstNumberNode{42}, true, TypBool},
				{"bool", &AstBoolNode{true}, true, TypBool},
				{"bool", &AstStringNode{""}, false, TypBool},
				{"bool", &AstNumberNode{0}, false, TypBool},
				{"bool", &AstBoolNode{false}, false, TypBool},
			}
			for _, d := range nodeData {
				n := &AstCastNode{
					Cast: d.cast,
					Term: d.term,
				}
				r := NewRuntime(nil)
				n.Accept(r)

				e := r.stack.Pop().(ScopeEntry)
				assert.Equal(t, d.expctVal, e.Value)
				assert.Equal(t, d.expctType, e.DataType)
			}
		})
	})

	t.Run("Test DivideNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate DivideNode", func(t *testing.T) {
			n := &AstDivideNode{
				Left:  &AstNumberNode{110},
				Right: &AstNumberNode{4},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, 27.5, e.Value)
			assert.Equal(t, TypNumber, e.DataType)
		})

		t.Run("Evaluate DivideNode with type error", func(t *testing.T) {
			n := &AstDivideNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test EqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate EqualNode", func(t *testing.T) {
			n := &AstEqualNode{
				Left:  &AstBoolNode{true},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})

		t.Run("Evaluate EqualNode with type error", func(t *testing.T) {
			n := &AstEqualNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test GreaterEqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate GreaterEqualNode", func(t *testing.T) {
			n := &AstGreaterEqualNode{
				Left:  &AstNumberNode{1984},
				Right: &AstNumberNode{1776},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})

		t.Run("Evaluate GreaterEqualNode with type error", func(t *testing.T) {
			n := &AstGreaterEqualNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test GreaterNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate GreaterNode", func(t *testing.T) {
			n := &AstGreaterNode{
				Left:  &AstNumberNode{1984},
				Right: &AstNumberNode{1776},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})

		t.Run("Evaluate GreaterNode with type error", func(t *testing.T) {
			n := &AstGreaterNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(&RuntimeEnv{nil, nil, nil})
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test LessEqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate LessEqualNode", func(t *testing.T) {
			n := &AstLessEqualNode{
				Left:  &AstNumberNode{1984},
				Right: &AstNumberNode{1776},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, false, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})

		t.Run("Evaluate LessEqualNode with type error", func(t *testing.T) {
			n := &AstLessEqualNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test LessNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate LessNode", func(t *testing.T) {
			n := &AstLessNode{
				Left:  &AstNumberNode{1984},
				Right: &AstNumberNode{1776},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, false, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})

		t.Run("Evaluate LessNode with type error", func(t *testing.T) {
			n := &AstLessNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test ModuloNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate ModuloNode", func(t *testing.T) {
			n := &AstModuloNode{
				Left:  &AstNumberNode{73},
				Right: &AstNumberNode{42},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, 31.0, e.Value)
			assert.Equal(t, TypNumber, e.DataType)
		})

		t.Run("Evaluate ModuloNode with type error", func(t *testing.T) {
			n := &AstModuloNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test MultiplyNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate MultiplyNode", func(t *testing.T) {
			n := &AstMultiplyNode{
				Left:  &AstNumberNode{21},
				Right: &AstNumberNode{2},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, 42.0, e.Value)
			assert.Equal(t, TypNumber, e.DataType)
		})

		t.Run("Evaluate MultiplyNode with type error", func(t *testing.T) {
			n := &AstMultiplyNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test NegativeNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate NegativeNode", func(t *testing.T) {
			n := &AstNegativeNode{&AstNumberNode{42}}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, -42.0, e.Value)
			assert.Equal(t, TypNumber, e.DataType)
		})

		t.Run("Evaluate NegativeNode with type error", func(t *testing.T) {
			n := &AstNegativeNode{&AstBoolNode{true}}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test NotEqualNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate NotEqualNode", func(t *testing.T) {
			n := &AstNotEqualNode{
				Left:  &AstBoolNode{true},
				Right: &AstBoolNode{false},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})

		t.Run("Evaluate NotEqualNode with type error", func(t *testing.T) {
			n := &AstNotEqualNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test NotNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate NotNode", func(t *testing.T) {
			n := &AstNotNode{&AstBoolNode{false}}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})

		t.Run("Evaluate NotNode with type error", func(t *testing.T) {
			n := &AstNotNode{&AstNumberNode{42}}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test OrNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate OrNode", func(t *testing.T) {
			n := &AstOrNode{
				Left:  &AstBoolNode{false},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})

		t.Run("Evaluate OrNode short circuited", func(t *testing.T) {
			n := &AstOrNode{
				Left: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})

		t.Run("Evaluate OrNode with type error", func(t *testing.T) {
			n := &AstOrNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test PositiveNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate PositiveNode", func(t *testing.T) {
			n := &AstPositiveNode{&AstNumberNode{-42}}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, 42.0, e.Value)
			assert.Equal(t, TypNumber, e.DataType)
		})

		t.Run("Evaluate PositiveNode with type error", func(t *testing.T) {
			n := &AstPositiveNode{&AstBoolNode{true}}
			r := NewRuntime(nil)
			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})

	t.Run("Test ReturnNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate ReturnNode", func(t *testing.T) {
			n := &AstReturnNode{
				Expr: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, true, e.Value)
			assert.Equal(t, TypBool, e.DataType)
		})
	})

	t.Run("Test SubtractNode", func(t *testing.T) {
		t.Parallel()

		t.Run("Evaluate SubtractNode", func(t *testing.T) {
			n := &AstSubtractNode{
				Left:  &AstNumberNode{73},
				Right: &AstNumberNode{42},
			}
			r := NewRuntime(nil)
			n.Accept(r)
			e := r.stack.Pop().(ScopeEntry)
			assert.Equal(t, 31.0, e.Value)
			assert.Equal(t, TypNumber, e.DataType)
		})

		t.Run("Eval SubtractNode with type error", func(t *testing.T) {
			n := &AstSubtractNode{
				Left:  &AstNumberNode{42},
				Right: &AstBoolNode{true},
			}
			r := NewRuntime(nil)
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
			r := NewRuntime(nil)
			r.currScope.SetVar("foo", ScopeEntry{TypString, "bar"})
			n.Accept(r)

			e, _ := r.currScope.GetVar("foo")
			assert.Equal(t, "bar", e.Value)
			assert.Equal(t, TypString, e.DataType)
		})

		t.Run("Evaluate VariableNode undefined", func(t *testing.T) {
			n := &AstVariableNode{
				Name: "foo",
			}
			r := NewRuntime(nil)

			assert.Panics(t, func() {
				n.Accept(r)
			})
		})
	})
}
