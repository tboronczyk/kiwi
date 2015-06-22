package symtable

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataTypeToString(t *testing.T) {
	dtypes := []struct{ actual, expected string }{
		{UNKNOWN.String(), "UNKNOWN"},
		{BOOL.String(), "BOOL"},
		{NUMBER.String(), "NUMBER"},
		{STRING.String(), "STRING"},
		{DataType(254).String(), "DataType(254)"},
	}

	for _, dt := range dtypes {
		assert.Equal(t, dt.expected, dt.actual)
	}
}
