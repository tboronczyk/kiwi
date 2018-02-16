package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRuntimeEnv(str string) *RuntimeEnv {
	return &RuntimeEnv{
		strings.NewReader(str),
		bytes.NewBuffer([]byte{}),
		bytes.NewBuffer([]byte{}),
	}
}

func TestBuiltins(t *testing.T) {
	t.Parallel()

	greeting := ScopeEntry{TypString, "hello world"}

	t.Run("strlen", func(t *testing.T) {
		s := &Stack{}
		p := []ScopeEntry{greeting}
		env := testRuntimeEnv("")

		builtins["strlen"](s, p, env)
		assert.Equal(t, ScopeEntry{TypNumber, 11}, s.Pop().(ScopeEntry))
	})

	t.Run("write", func(t *testing.T) {
		s := &Stack{}
		p := []ScopeEntry{greeting}
		env := testRuntimeEnv(greeting.Value.(string))

		builtins["write"](s, p, env)
		assert.Equal(t, greeting.Value.(string), env.stdout.(*bytes.Buffer).String())
	})

	t.Run("read", func(t *testing.T) {
		s := &Stack{}
		p := []ScopeEntry{greeting}
		env := testRuntimeEnv(greeting.Value.(string))

		builtins["read"](s, p, env)
		assert.Equal(t, greeting, s.Pop().(ScopeEntry))
	})
}
