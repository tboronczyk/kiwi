package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokens(t *testing.T) {
	t.Parallel()

	t.Run("Test token to string", func(t *testing.T) {
		tokens := []struct{ actual, expected string }{
			{T_UNKNOWN.String(), "UNKNOWN"},
			{T_EOF.String(), "EOF"},
			{T_ADD.String(), "+"},
			{T_SUBTRACT.String(), "-"},
			{T_MULTIPLY.String(), "*"},
			{T_DIVIDE.String(), "/"},
			{T_MODULO.String(), "%"},
			{T_EQUAL.String(), "="},
			{T_NOT_EQUAL.String(), "~="},
			{T_LESS.String(), "<"},
			{T_LESS_EQ.String(), "<="},
			{T_GREATER.String(), ">"},
			{T_GREATER_EQ.String(), ">="},
			{T_AND.String(), "&&"},
			{T_OR.String(), "||"},
			{T_NOT.String(), "~"},
			{T_IF.String(), "if"},
			{T_FUNC.String(), "func"},
			{T_RETURN.String(), "return"},
			{T_WHILE.String(), "while"},
			{T_BOOL.String(), "BOOL"},
			{T_IDENTIFIER.String(), "IDENTIFIER"},
			{T_NUMBER.String(), "NUMBER"},
			{T_STRING.String(), "STRING"},
			{T_ASSIGN.String(), ":="},
			{T_LBRACE.String(), "{"},
			{T_RBRACE.String(), "}"},
			{T_COLON.String(), ":"},
			{T_COMMA.String(), ","},
			{T_COMMENT.String(), "COMMENT"},
			{T_ELSE.String(), "else"},
			{T_LPAREN.String(), "("},
			{T_RPAREN.String(), ")"},
			{Token(254).String(), "Token(254)"},
		}

		for _, tkn := range tokens {
			assert.Equal(t, tkn.actual, tkn.expected)
		}
	})

	t.Run("Test IsAddOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == T_ADD || tkn == T_SUBTRACT {
				assert.True(t, tkn.IsAddOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsAddOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsMulOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == T_MULTIPLY || tkn == T_DIVIDE || tkn == T_MODULO {
				assert.True(t, tkn.IsMulOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsMulOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsCmpOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == T_EQUAL || tkn == T_NOT_EQUAL || tkn == T_LESS ||
				tkn == T_LESS_EQ || tkn == T_GREATER || tkn == T_GREATER_EQ {
				assert.True(t, tkn.IsCmpOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsCmpOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsLogOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == T_AND || tkn == T_OR || tkn == T_NOT {
				assert.True(t, tkn.IsLogOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsLogOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsBinOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == T_ADD || tkn == T_SUBTRACT ||
				tkn == T_MULTIPLY || tkn == T_DIVIDE || tkn == T_MODULO ||
				tkn == T_EQUAL || tkn == T_NOT_EQUAL || tkn == T_LESS ||
				tkn == T_LESS_EQ || tkn == T_GREATER || tkn == T_GREATER_EQ ||
				tkn == T_AND || tkn == T_OR {
				assert.True(t, tkn.IsBinOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsBinOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsUnaryOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == T_ADD || tkn == T_SUBTRACT || tkn == T_NOT {
				assert.True(t, tkn.IsUnaryOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsUnaryOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsStmtKeyword", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == T_IF || tkn == T_WHILE || tkn == T_FUNC || tkn == T_RETURN {
				assert.True(t, tkn.IsStmtKeyword(), tkn.String())
			} else {
				assert.False(t, tkn.IsStmtKeyword(), tkn.String())
			}
		}
	})

	t.Run("Test IsLiteral", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == T_IDENTIFIER || tkn == T_BOOL || tkn == T_NUMBER ||
				tkn == T_STRING {
				assert.True(t, tkn.IsLiteral(), tkn.String())
			} else {
				assert.False(t, tkn.IsLiteral(), tkn.String())
			}
		}
	})

	t.Run("Test precedence", func(t *testing.T) {
		tokens := []struct{ t1, t2 Token }{
			{T_MULTIPLY, T_ADD},
			{T_ADD, T_LESS},
			{T_LESS, T_AND},
		}
		for _, tkns := range tokens {
			p := Precedence(tkns.t1) > Precedence(tkns.t2)
			assert.True(t, p, tkns.t1.String()+" and "+tkns.t2.String())
		}
	})

	t.Run("Test precedence of non-operator", func(t *testing.T) {
		i := Precedence(T_IF)
		assert.Equal(t, 0, i)
	})
}
