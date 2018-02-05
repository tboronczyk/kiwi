package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokens(t *testing.T) {
	t.Parallel()

	t.Run("Test token to string", func(t *testing.T) {
		tokens := []struct{ actual, expected string }{
			{TkUnknown.String(), "Unknown"},
			{TkEof.String(), "eof"},
			{TkAdd.String(), "+"},
			{TkSubtract.String(), "-"},
			{TkMultiply.String(), "*"},
			{TkDivide.String(), "/"},
			{TkModulo.String(), "%"},
			{TkEqual.String(), "="},
			{TkNotEqual.String(), "~="},
			{TkLess.String(), "<"},
			{TkLessEq.String(), "<="},
			{TkGreater.String(), ">"},
			{TkGreaterEq.String(), ">="},
			{TkAnd.String(), "&&"},
			{TkOr.String(), "||"},
			{TkNot.String(), "~"},
			{TkIf.String(), "if"},
			{TkFunc.String(), "func"},
			{TkReturn.String(), "return"},
			{TkWhile.String(), "while"},
			{TkBool.String(), "Bool"},
			{TkIdentifier.String(), "Identifier"},
			{TkNumber.String(), "Number"},
			{TkString.String(), "String"},
			{TkAssign.String(), ":="},
			{TkLBrace.String(), "{"},
			{TkRBrace.String(), "}"},
			{TkColon.String(), ":"},
			{TkComma.String(), ","},
			{TkComment.String(), "Comment"},
			{TkElse.String(), "else"},
			{TkLParen.String(), "("},
			{TkRParent.String(), ")"},
			{Token(254).String(), "Token(254)"},
		}

		for _, tkn := range tokens {
			assert.Equal(t, tkn.actual, tkn.expected)
		}
	})

	t.Run("Test IsAddOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == TkAdd || tkn == TkSubtract {
				assert.True(t, tkn.IsAddOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsAddOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsMulOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == TkMultiply || tkn == TkDivide || tkn == TkModulo {
				assert.True(t, tkn.IsMulOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsMulOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsCmpOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == TkEqual || tkn == TkNotEqual || tkn == TkLess ||
				tkn == TkLessEq || tkn == TkGreater || tkn == TkGreaterEq {
				assert.True(t, tkn.IsCmpOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsCmpOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsLogOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == TkAnd || tkn == TkOr || tkn == TkNot {
				assert.True(t, tkn.IsLogOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsLogOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsBinOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == TkAdd || tkn == TkSubtract ||
				tkn == TkMultiply || tkn == TkDivide || tkn == TkModulo ||
				tkn == TkEqual || tkn == TkNotEqual ||
				tkn == TkLess || tkn == TkLessEq ||
				tkn == TkGreater || tkn == TkGreaterEq ||
				tkn == TkAnd || tkn == TkOr {
				assert.True(t, tkn.IsBinOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsBinOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsUnaryOp", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == TkAdd || tkn == TkSubtract || tkn == TkNot {
				assert.True(t, tkn.IsUnaryOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsUnaryOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsStmtKeyword", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == TkIf || tkn == TkWhile || tkn == TkFunc || tkn == TkReturn {
				assert.True(t, tkn.IsStmtKeyword(), tkn.String())
			} else {
				assert.False(t, tkn.IsStmtKeyword(), tkn.String())
			}
		}
	})

	t.Run("Test IsLiteral", func(t *testing.T) {
		for i := 0; i < len(tokens); i++ {
			tkn := Token(i)
			if tkn == TkIdentifier || tkn == TkBool || tkn == TkNumber ||
				tkn == TkString {
				assert.True(t, tkn.IsLiteral(), tkn.String())
			} else {
				assert.False(t, tkn.IsLiteral(), tkn.String())
			}
		}
	})

	t.Run("Test precedence", func(t *testing.T) {
		tokens := []struct{ t1, t2 Token }{
			{TkMultiply, TkAdd},
			{TkAdd, TkLess},
			{TkLess, TkAnd},
		}
		for _, tkns := range tokens {
			p := Precedence(tkns.t1) > Precedence(tkns.t2)
			assert.True(t, p, tkns.t1.String()+" and "+tkns.t2.String())
		}
	})

	t.Run("Test precedence of non-operator", func(t *testing.T) {
		i := Precedence(TkIf)
		assert.Equal(t, 0, i)
	})
}
