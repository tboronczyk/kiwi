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
			{T_ADD, "+"},
			{T_SUBTRACT, "-"},
			{T_MULTIPLY, "*"},
			{T_DIVIDE, "/"},
			{T_MODULO, "%"},
			{T_ASSIGN, ":="},
			{T_COLON, ":"},
			{T_EQUAL, "="},
			{T_LESS, "<"},
			{T_LESS_EQ, "<="},
			{T_GREATER, ">"},
			{T_GREATER_EQ, ">="},
			{T_AND, "&&"},
			{T_UNKNOWN, "&"},
			{T_OR, "||"},
			{T_UNKNOWN, "|"},
			{T_NOT, "~"},
			{T_NOT_EQUAL, "~="},
			{T_LPAREN, "("},
			{T_RPAREN, ")"},
			{T_LBRACE, "{"},
			{T_RBRACE, "}"},
			{T_COMMA, ","},
			{T_UNKNOWN, "?"},
			{T_EOF, ""},
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
			{T_FUNC, "func"},
			{T_IF, "if"},
			{T_ELSE, "else"},
			{T_RETURN, "return"},
			{T_WHILE, "while"},
			{T_BOOL, "TRUE"},
			{T_BOOL, "FALSE"},
			{T_IDENTIFIER, "if"},
			{T_IDENTIFIER, "ident"},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, expected.token, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})

	t.Run("Test scan strings", func(t *testing.T) {
		str := "\"abc\"" +
			"\"\"" +
			"\"\\\\\\\"\\r\\n\\t\\x\"" +
			"\"broken"
		s := NewScanner(strings.NewReader(str))

		tokens := []struct {
			token Token
			value string
		}{
			{T_STRING, "abc"},
			{T_STRING, ""},
			{T_STRING, "\\\"\r\n\t\\x"},
			{T_UNKNOWN, "broken"},
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
			assert.Equal(t, T_COMMENT, actual1)
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
			{T_COMMENT, "/**/"},
			{T_COMMENT, "/* a /* nested */ comment */"},
			{T_UNKNOWN, "/* broken"},
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
			{T_NUMBER, "123"},
			{T_NUMBER, "0.123"},
			{T_NUMBER, "1."},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, expected.token, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})
}
