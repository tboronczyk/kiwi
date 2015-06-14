package symtable

import (
	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/token"
	"testing"
)

func TestSet(t *testing.T) {
	s := New()
	s.Set("foo", 42, token.NUMBER)

	expected := entry{v: 42, t: token.NUMBER}
	actual := s.table["foo"]
	assert.Equal(t, expected, actual)
}

func TestGet(t *testing.T) {
	s := New()
	s.table["foo"] = entry{v: 42, t: token.NUMBER}

	expected := 42
	actual, _, _ := s.Get("foo")
	assert.Equal(t, expected, actual)
}
