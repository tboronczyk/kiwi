package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokens(t *testing.T) {
	t.Parallel()

	t.Run("Test IsAddOp", func(t *testing.T) {
		for i := 0; i < int(endTokens); i++ {
			tkn := Token(i)
			if tkn == TkAdd || tkn == TkSubtract {
				assert.True(t, tkn.IsAddOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsAddOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsMulOp", func(t *testing.T) {
		for i := 0; i < int(endTokens); i++ {
			tkn := Token(i)
			if tkn == TkMultiply || tkn == TkDivide || tkn == TkModulo {
				assert.True(t, tkn.IsMulOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsMulOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsCmpOp", func(t *testing.T) {
		for i := 0; i < int(endTokens); i++ {
			tkn := Token(i)
			if tkn == TkEqual || tkn == TkNotEqual ||
				tkn == TkLess || tkn == TkLessEq ||
				tkn == TkGreater || tkn == TkGreaterEq {
				assert.True(t, tkn.IsCmpOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsCmpOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsLogOp", func(t *testing.T) {
		for i := 0; i < int(endTokens); i++ {
			tkn := Token(i)
			if tkn == TkAnd || tkn == TkOr || tkn == TkNot {
				assert.True(t, tkn.IsLogOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsLogOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsBinOp", func(t *testing.T) {
		for i := 0; i < int(endTokens); i++ {
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
		for i := 0; i < int(endTokens); i++ {
			tkn := Token(i)
			if tkn == TkAdd || tkn == TkSubtract || tkn == TkNot {
				assert.True(t, tkn.IsUnaryOp(), tkn.String())
			} else {
				assert.False(t, tkn.IsUnaryOp(), tkn.String())
			}
		}
	})

	t.Run("Test IsStmtKeyword", func(t *testing.T) {
		for i := 0; i < int(endTokens); i++ {
			tkn := Token(i)
			if tkn == TkIf || tkn == TkWhile || tkn == TkFunc || tkn == TkReturn {
				assert.True(t, tkn.IsStmtKeyword(), tkn.String())
			} else {
				assert.False(t, tkn.IsStmtKeyword(), tkn.String())
			}
		}
	})

	t.Run("Test IsLiteral", func(t *testing.T) {
		for i := 0; i < int(endTokens); i++ {
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
		assert.Panics(t, func() {
			Precedence(TkIf)
		})
	})
}
