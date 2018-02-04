package scope

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/types"
)

func TestScope(t *testing.T) {
	t.Parallel()

	t.Run("Test variable scope", func(t *testing.T) {
		s := New()
		s.SetVar("foo", Entry{Value: 42, DataType: types.NUMBER})
		entry, ok := s.GetVar("foo")
		assert.Equal(t, 42, entry.Value)
		assert.True(t, ok)

		s = NewWithParent(s)
		_, ok = s.GetVar("foo")
		assert.False(t, ok)
	})

	t.Run("Test function scope", func(t *testing.T) {
		s := New()
		s.SetFunc("foo", Entry{Value: 42, DataType: types.NUMBER})
		entry, ok := s.GetFunc("foo")
		assert.Equal(t, 42, entry.Value)
		assert.True(t, ok)

		s = NewWithParent(s)
		entry, ok = s.GetFunc("foo")
		assert.Equal(t, 42, entry.Value)
		assert.True(t, ok)
	})
}
