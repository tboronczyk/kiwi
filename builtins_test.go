package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuiltins(t *testing.T) {
	t.Parallel()

	t.Run("strlen", func(t *testing.T) {
		s := &Stack{}
		p := []ScopeEntry{
			ScopeEntry{TypString, "hello world"},
		}

		builtins["strlen"](s, p, nil)
		result := s.Pop().(ScopeEntry)
		assert.Equal(t, 11, result.Value)
		assert.Equal(t, TypNumber, result.DataType)
	})

	t.Run("write", func(t *testing.T) {
		s := &Stack{}
		p := []ScopeEntry{
			ScopeEntry{TypString, "hello world"},
		}
		out := bytes.NewBuffer([]byte{})
		builtins["write"](s, p, &RuntimeEnv{nil, out, nil})
		assert.Equal(t, "hello world", out.String())
	})

	t.Run("read", func(t *testing.T) {
		s := &Stack{}
		p := []ScopeEntry{
			ScopeEntry{TypString, "hello world"},
		}
		in := strings.NewReader("hello world")
		builtins["read"](s, p, &RuntimeEnv{in, nil, nil})
		result := s.Pop().(ScopeEntry)
		assert.Equal(t, "hello world", result.Value)
		assert.Equal(t, TypString, result.DataType)
	})
}
