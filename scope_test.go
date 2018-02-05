package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScope(t *testing.T) {
	t.Parallel()

	t.Run("Test variable scope", func(t *testing.T) {
		s := NewScope()
		s.SetVar("foo", ScopeEntry{TypNumber, 42})
		entry, ok := s.GetVar("foo")
		assert.Equal(t, 42, entry.Value)
		assert.True(t, ok)

		s = NewScopeWithParent(s)
		_, ok = s.GetVar("foo")
		assert.False(t, ok)
	})

	t.Run("Test function scope", func(t *testing.T) {
		s := NewScope()
		s.SetFunc("foo", ScopeEntry{TypNumber, 42})
		entry, ok := s.GetFunc("foo")
		assert.Equal(t, 42, entry.Value)
		assert.True(t, ok)

		s = NewScopeWithParent(s)
		entry, ok = s.GetFunc("foo")
		assert.Equal(t, 42, entry.Value)
		assert.True(t, ok)
	})
}
