package scanner

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/token"
)

func TestScanner(t *testing.T) {
	t.Parallel()

	t.Run("Test scan simple tokens", func(t *testing.T) {
		str := "+ - * / % := : = < <= > >= && & || | ~ ~= ( ) { } , ?"
		s := New(strings.NewReader(str))

		tokens := []struct {
			token token.Token
			value string
		}{
			{token.ADD, "+"},
			{token.SUBTRACT, "-"},
			{token.MULTIPLY, "*"},
			{token.DIVIDE, "/"},
			{token.MODULO, "%"},
			{token.ASSIGN, ":="},
			{token.COLON, ":"},
			{token.EQUAL, "="},
			{token.LESS, "<"},
			{token.LESS_EQ, "<="},
			{token.GREATER, ">"},
			{token.GREATER_EQ, ">="},
			{token.AND, "&&"},
			{token.UNKNOWN, "&"},
			{token.OR, "||"},
			{token.UNKNOWN, "|"},
			{token.NOT, "~"},
			{token.NOT_EQUAL, "~="},
			{token.LPAREN, "("},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RBRACE, "}"},
			{token.COMMA, ","},
			{token.UNKNOWN, "?"},
			{token.EOF, ""},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, expected.token, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})

	t.Run("Test scan identifiers", func(t *testing.T) {
		str := "func if else return while true false `if ident"
		s := New(strings.NewReader(str))

		tokens := []struct {
			token token.Token
			value string
		}{
			{token.FUNC, "func"},
			{token.IF, "if"},
			{token.ELSE, "else"},
			{token.RETURN, "return"},
			{token.WHILE, "while"},
			{token.BOOL, "TRUE"},
			{token.BOOL, "FALSE"},
			{token.IDENTIFIER, "if"},
			{token.IDENTIFIER, "ident"},
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
		s := New(strings.NewReader(str))

		tokens := []struct {
			token token.Token
			value string
		}{
			{token.STRING, "abc"},
			{token.STRING, ""},
			{token.STRING, "\\\"\r\n\t\\x"},
			{token.UNKNOWN, "broken"},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, expected.token, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})

	t.Run("Test scan line comments", func(t *testing.T) {
		str := "// single1\n// single2"
		s := New(strings.NewReader(str))

		tokens := []struct {
			value string
		}{
			{"// single1"},
			{"// single2"},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, token.COMMENT, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})

	t.Run("Test scan multiline comments", func(t *testing.T) {
		str := "/**/" +
			"/* a /* nested */ comment */" +
			"/* broken"
		s := New(strings.NewReader(str))

		tokens := []struct {
			token token.Token
			value string
		}{
			{token.COMMENT, "/**/"},
			{token.COMMENT, "/* a /* nested */ comment */"},
			{token.UNKNOWN, "/* broken"},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, expected.token, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})

	t.Run("Test scan numbers", func(t *testing.T) {
		str := "123 0.123 1."
		s := New(strings.NewReader(str))

		tokens := []struct {
			token token.Token
			value string
		}{
			{token.NUMBER, "123"},
			{token.NUMBER, "0.123"},
			{token.NUMBER, "1."},
		}

		for _, expected := range tokens {
			actual1, actual2 := s.Scan()
			assert.Equal(t, expected.token, actual1)
			assert.Equal(t, expected.value, actual2)
		}
	})
}
