package token

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testData struct {
	actual   string
	expected string
}

func TestTokenToString(t *testing.T) {
	tokens := []testData{
		{UNKNOWN.String(), "UNKNOWN"},
		{MALFORMED.String(), "MALFORMED"},
		{EOF.String(), "EOF"},
		{ADD.String(), "+"},
		{SUBTRACT.String(), "-"},
		{MULTIPLY.String(), "*"},
		{DIVIDE.String(), "/"},
		{MODULO.String(), "%"},
		{ASSIGN.String(), ":="},
		{EQUAL.String(), "="},
		{NOT_EQUAL.String(), "~="},
		{LESS.String(), "<"},
		{LESS_EQ.String(), "<="},
		{GREATER.String(), ">"},
		{GREATER_EQ.String(), ">="},
		{AND.String(), "&&"},
		{OR.String(), "||"},
		{NOT.String(), "~"},
		{LPAREN.String(), "("},
		{RPAREN.String(), ")"},
		{LBRACE.String(), "{"},
		{RBRACE.String(), "}"},
		{SEMICOLON.String(), ";"},
		{COMMA.String(), ","},
		{COMMENT.String(), "COMMENT"},
		{IF.String(), "if"},
		{WHILE.String(), "while"},
		{FUNC.String(), "func"},
		{RETURN.String(), "return"},
		{TRUE.String(), "true"},
		{FALSE.String(), "false"},
		{NUMBER.String(), "NUMBER"},
		{STRING.String(), "STRING"},
		{IDENTIFIER.String(), "IDENTIFIER"},
		{Token(254).String(), "Token(254)"},
	}

	for _, tkn := range tokens {
		assert.Equal(t, tkn.actual, tkn.expected)
	}
}

func TestIsAddOp(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tkn := Token(i)
		if tkn == ADD || tkn == SUBTRACT {
			assert.True(t, tkn.IsAddOp(), tkn.String())
		} else {
			assert.False(t, tkn.IsAddOp(), tkn.String())
		}
	}
}

func TestIsMulOp(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tkn := Token(i)
		if tkn == MULTIPLY || tkn == DIVIDE || tkn == MODULO {
			assert.True(t, tkn.IsMulOp(), tkn.String())
		} else {
			assert.False(t, tkn.IsMulOp(), tkn.String())
		}
	}
}

func TestIsCmpOp(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tkn := Token(i)
		if tkn == EQUAL || tkn == NOT_EQUAL || tkn == LESS ||
			tkn == LESS_EQ || tkn == GREATER || tkn == GREATER_EQ {
			assert.True(t, tkn.IsCmpOp(), tkn.String())
		} else {
			assert.False(t, tkn.IsCmpOp(), tkn.String())
		}
	}
}

func TestIsLogOp(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tkn := Token(i)
		if tkn == AND || tkn == OR || tkn == NOT {
			assert.True(t, tkn.IsLogOp(), tkn.String())
		} else {
			assert.False(t, tkn.IsLogOp(), tkn.String())
		}
	}
}

func TestIsExprOp(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tkn := Token(i)
		if tkn == ADD || tkn == SUBTRACT || tkn == NOT {
			assert.True(t, tkn.IsExprOp(), tkn.String())
		} else {
			assert.False(t, tkn.IsExprOp(), tkn.String())
		}
	}
}

func TestIsLiteral(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tkn := Token(i)
		if tkn == IDENTIFIER || tkn == TRUE || tkn == FALSE ||
			tkn == NUMBER || tkn == STRING {
			assert.True(t, tkn.IsLiteral(), tkn.String())
		} else {
			assert.False(t, tkn.IsLiteral(), tkn.String())
		}
	}
}

func TestIsStmtKeyword(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tkn := Token(i)
		if tkn == IF || tkn == WHILE || tkn == FUNC || tkn == RETURN {
			assert.True(t, tkn.IsStmtKeyword(), tkn.String())
		} else {
			assert.False(t, tkn.IsStmtKeyword(), tkn.String())
		}
	}
}
