package symtable

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testData struct {
	actual   string
	expected string
}

func TestSet(t *testing.T) {
	s := New()
	s.Set("foo", 42, NUMBER)

	expected := entry{v: 42, t: NUMBER}
	actual := s.table["foo"]
	assert.Equal(t, expected, actual)
}

func TestGet(t *testing.T) {
	s := New()
	s.table["foo"] = entry{v: 42, t: NUMBER}

	expected := 42
	actual, _, _ := s.Get("foo")
	assert.Equal(t, expected, actual)
}

func TestGetNoExist(t *testing.T) {
	s := New()
	_, _, found := s.Get("foo")
	assert.False(t, found)
}

func TestDataTypeToString(t *testing.T) {
	tokens := []testData{
		{UNKNOWN.String(), "UNKNOWN"},
		{BOOL.String(), "BOOL"},
		{NUMBER.String(), "NUMBER"},
		{STRING.String(), "STRING"},
		{DataType(254).String(), "DataType(254)"},
	}

	for _, tkn := range tokens {
		assert.Equal(t, tkn.actual, tkn.expected)
	}
}
