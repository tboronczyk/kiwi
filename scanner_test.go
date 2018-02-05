package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	t.Parallel()

	t.Run("Test scan simple tokens", func(t *testing.T) {
		str := "+ - * / % := : = < <= > >= && & || | ~ ~= ( ) { } , ?"
		s := NewScanner(strings.NewReader(str))

		tokens := []struct {
			token Token
			value string
		}{
			{TkAdd, "+"},
			{TkSubtract, "-"},
			{TkMultiply, "*"},
			{TkDivide, "/"},
			{TkModulo, "%"},
			{TkAssign, ":="},
			{TkColon, ":"},
			{TkEqual, "="},
			{TkLess, "<"},
			{TkLessEq, "<="},
			{TkGreater, ">"},
			{TkGreaterEq, ">="},
			{TkAnd, "&&"},
			{TkUnknown, "&"},
			{TkOr, "||"},
			{TkUnknown, "|"},
			{TkIf, "~"},
			{TkNotEqual, "~="},
			{TkLParen, "("},
			{TkRParent, ")"},
			{TkLBrace, "{"},
			{TkRBrace, "}"},
			{TkComma, ","},
			{TkUnknown, "?"},
			{TkEof, ""},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, expected.token, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})

	t.Run("Test scan identifiers", func(t *testing.T) {
		str := "func if else return while true false `if ident"
		s := NewScanner(strings.NewReader(str))

		tokens := []struct {
			token Token
			value string
		}{
			{TkFunc, "func"},
			{TkIf, "if"},
			{TkElse, "else"},
			{TkReturn, "return"},
			{TkWhile, "while"},
			{TkBool, "TRUE"},
			{TkBool, "FALSE"},
			{TkIdentifier, "if"},
			{TkIdentifier, "ident"},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, expected.token, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})

	t.Run("Test scan strings", func(t *testing.T) {
		str := `"abc"` +
			`""` +
			`"\\\"` + "\r\n\t" + `\x"` +
			`"broken`
		s := NewScanner(strings.NewReader(str))

		tokens := []struct {
			token Token
			value string
		}{
			{TkString, "abc"},
			{TkString, ""},
			{TkString, "\\\"\r\n\t\\x"},
			{TkUnknown, "broken"},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, expected.token, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})

	t.Run("Test scan line comments", func(t *testing.T) {
		str := "// single1\n// single2"
		s := NewScanner(strings.NewReader(str))

		tokens := []struct {
			value string
		}{
			{"// single1"},
			{"// single2"},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, TkComment, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})

	t.Run("Test scan multiline comments", func(t *testing.T) {
		str := "/**/" +
			"/* a /* nested */ comment */" +
			"/* broken"
		s := NewScanner(strings.NewReader(str))

		tokens := []struct {
			token Token
			value string
		}{
			{TkComment, "/**/"},
			{TkComment, "/* a /* nested */ comment */"},
			{TkUnknown, "/* broken"},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, expected.token, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})

	t.Run("Test scan numbers", func(t *testing.T) {
		str := "123 0.123 1."
		s := NewScanner(strings.NewReader(str))

		tokens := []struct {
			token Token
			value string
		}{
			{TkNumber, "123"},
			{TkNumber, "0.123"},
			{TkNumber, "1."},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, expected.token, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})
}
