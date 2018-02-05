package main

import (
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

		builtins["strlen"](s, p)
		result := s.Pop().(ScopeEntry)
		assert.Equal(t, 11, result.Value)
		assert.Equal(t, TypNumber, result.DataType)
	})
}
