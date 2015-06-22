package symtable

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSymbolTableSet(t *testing.T) {
	s := New()
	s.Set("foo", 42, NUMBER)

	actual := s["foo"]
	assert.Equal(t, 42, actual.value)
	assert.Equal(t, NUMBER, actual.dtype)
}

func TestSymbolTableGet(t *testing.T) {
	s := New()
	s.Set("foo", 42, NUMBER)

	value, dtype, _ := s.Get("foo")
	assert.Equal(t, 42, value)
	assert.Equal(t, NUMBER, dtype)
}

func TestSymbolTableGetNoExist(t *testing.T) {
	s := New()
	_, _, ok := s.Get("foo")
	assert.False(t, ok)
}
