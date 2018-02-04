package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	t.Parallel()

	t.Run("Test data types to string", func(t *testing.T) {
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
	})
}
