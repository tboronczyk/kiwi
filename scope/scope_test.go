package scope

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/types"
)

func TestScopeVar(t *testing.T) {
	s := New()
	s.SetVar("foo", Entry{Value: 42, DataType: types.NUMBER})
	entry, ok := s.GetVar("foo")
	assert.Equal(t, 42, entry.Value)
	assert.True(t, ok)

	s = NewWithParent(s)
	_, ok = s.GetVar("foo")
	assert.False(t, ok)
}

func TestScopeFunc(t *testing.T) {
	s := New()
	s.SetFunc("foo", Entry{Value: 42, DataType: types.NUMBER})
	entry, ok := s.GetFunc("foo")
	assert.Equal(t, 42, entry.Value)
	assert.True(t, ok)

	s = NewWithParent(s)
	entry, ok = s.GetFunc("foo")
	assert.Equal(t, 42, entry.Value)
	assert.True(t, ok)
}
