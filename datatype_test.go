package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataType(t *testing.T) {
	t.Parallel()

	t.Run("Test type to string", func(t *testing.T) {
		types := map[DataType]string{
			TypUnknown:    "TypUnknown",
			TypBuiltin:    "TypBuiltin",
			TypBool:       "TypBool",
			TypFunc:       "TypFunc",
			TypNumber:     "TypNumber",
			TypString:     "TypString",
			DataType(255): "DataType(255)",
		}

		for typ, str := range types {
			assert.Equal(t, str, typ.String())
		}
	})
}
