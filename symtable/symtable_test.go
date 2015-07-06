package symtable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymbolTableSetGet(t *testing.T) {
	s := New()
	s.Set("foo", 42)

	value, ok := s.Get("foo")
	assert.Equal(t, 42, value)
	assert.True(t, ok)
}

func TestSymbolTableNewScopeSetGet(t *testing.T) {
	s := New()
	s.Set("foo", 42)

	s = NewScope(s)
	s.Set("foo", 73)

	value, _ := s.Get("foo")
	assert.Equal(t, 73, value)

	s = s.Parent()

	value, _ = s.Get("foo")
	assert.Equal(t, 42, value)
}

func TestSymbolTableParentSearch(t *testing.T) {
	s := New()
	s.Set("foo", 42)
	s = NewScope(s)

	value, _ := s.Get("foo")
	assert.Equal(t, 42, value)
}
