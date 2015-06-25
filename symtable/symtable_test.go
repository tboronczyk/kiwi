package symtable

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSymbolTableSetGet(t *testing.T) {
	s := New()
	s.Set("foo", VARIABLE, 42)

	value, ok := s.Get("foo", VARIABLE)
	assert.Equal(t, 42, value)
	assert.True(t, ok)
}

func TestSymbolTableGetNoExist(t *testing.T) {
	s := New()

	value, ok := s.Get("foo", VARIABLE)
	assert.Nil(t, value)
	assert.False(t, ok)
}

func TestSymbolTableNewScopeGet(t *testing.T) {
	s := New()
	s.Set("foo", VARIABLE, 42)
	s = ScopeEnter(s)

	value, ok := s.Get("foo", VARIABLE)
	assert.Equal(t, 42, value)
	assert.True(t, ok)
}

func TestSymbolTableNewScopeGetNoExit(t *testing.T) {
	s := New()
	s = ScopeEnter(s)

	value, ok := s.Get("foo", VARIABLE)
	assert.Nil(t, value)
	assert.False(t, ok)
}

func TestSymbolTableNewScopeSet(t *testing.T) {
	s := New()
	s.Set("foo", VARIABLE, 42)
	s = ScopeEnter(s)
	s.Set("foo", VARIABLE, 73)

	value, ok := s.Get("foo", VARIABLE)
	assert.Equal(t, 73, value)
	assert.True(t, ok)

	s = ScopeLeave(s)
	value, ok = s.Get("foo", VARIABLE)
	assert.Equal(t, 42, value)
	assert.True(t, ok)
}
