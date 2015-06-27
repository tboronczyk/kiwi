package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataTypesToString(t *testing.T) {
	dtypes := []struct{ actual, expected string }{
		{UNKNOWN.String(), "UNKNOWN"},
		{BUILTIN.String(), "BUILTIN"},
		{BOOL.String(), "BOOL"},
		{FUNC.String(), "FUNC"},
		{NUMBER.String(), "NUMBER"},
		{STRING.String(), "STRING"},
		{DataType(254).String(), "DataType(254)"},
	}

	for _, dtype := range dtypes {
		assert.Equal(t, dtype.actual, dtype.expected)
	}
}
