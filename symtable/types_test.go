package symtable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymTypesToString(t *testing.T) {
	stypes := []struct{ actual, expected string }{
		{UNKNOWN.String(), "UNKNOWN"},
		{VAR.String(), "VAR"},
		{FUNC.String(), "FUNC"},
		{SymType(254).String(), "SymType(254)"},
	}

	for _, stype := range stypes {
		assert.Equal(t, stype.actual, stype.expected)
	}
}
