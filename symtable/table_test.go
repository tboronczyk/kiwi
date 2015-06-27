package symtable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymbolTableSetGet(t *testing.T) {
	s := New()
	s.Set("foo", VAR, 42)

	value, ok := s.Get("foo", VAR)
	assert.Equal(t, 42, value)
	assert.True(t, ok)

	_, ok = s.Get("foo", FUNC)
	assert.False(t, ok)
}

func TestSymbolTableNewScopeSetGet(t *testing.T) {
	s := New()
	s.Set("foo", VAR, 42)
	s = NewScope(s)

	s.Set("foo", VAR, 73)

	value, _ := s.Get("foo", VAR)
	assert.Equal(t, 73, value)

	s = s.Parent()
	value, _ = s.Get("foo", VAR)
	assert.Equal(t, 42, value)
}

func TestSymbolTableParentSearch(t *testing.T) {
	s := New()
	s.Set("foo", VAR, 42)
	s = NewScope(s)

	value, _ := s.Get("foo", VAR)
	assert.Equal(t, 42, value)
}
