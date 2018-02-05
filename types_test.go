package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	t.Parallel()

	t.Run("Test data types to string", func(t *testing.T) {
		types := []struct{ actual, expected string }{
			{TypUnknown.String(), "Unknown"},
			{TypBuiltin.String(), "Builtin"},
			{TypBool.String(), "Bool"},
			{TypFunc.String(), "Func"},
			{TypNumber.String(), "Number"},
			{TypString.String(), "String"},
			{DataType(254).String(), "DataType(254)"},
		}

		for _, dt := range types {
			assert.Equal(t, dt.actual, dt.expected)
		}
	})
}
