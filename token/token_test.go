/*
 * Copyright (c) 2012, 2015 Timothy Boronczyk
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met.String(),
 *
 *  1. Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *
 *  2. Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *
 *  3. The names of the authors may not be used to endorse or promote products
 *     derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED "AS IS" AND WITHOUT ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, WITHOUT LIMITATION, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE.
 */

package token

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenToString(t *testing.T) {
	tokens := []struct {
		actual   string
		expected string
	}{
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
		{COMMENT.String(), "COMMENT"},
		{IF.String(), "if"},
		{WHILE.String(), "while"},
		{TRUE.String(), "true"},
		{FALSE.String(), "false"},
		{NUMBER.String(), "NUMBER"},
		{STRING.String(), "STRING"},
		{IDENTIFIER.String(), "IDENTIFIER"},
		{Token(254).String(), "Token(254)"},
	}

	for _, tok := range tokens {
		assert.Equal(t, tok.actual, tok.expected)
	}
}

func TestIsAddOp(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tok := Token(i)
		if tok == ADD || tok == SUBTRACT {
			assert.True(t, tok.IsAddOp(), tok.String())
		} else {
			assert.False(t, tok.IsAddOp(), tok.String())
		}
	}
}

func TestIsMulOp(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tok := Token(i)
		if tok == MULTIPLY || tok == DIVIDE || tok == MODULO {
			assert.True(t, tok.IsMulOp(), tok.String())
		} else {
			assert.False(t, tok.IsMulOp(), tok.String())
		}
	}
}

func TestIsCmpOp(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tok := Token(i)
		if tok == EQUAL || tok == NOT_EQUAL || tok == LESS || tok == LESS_EQ ||
			tok == GREATER || tok == GREATER_EQ {
			assert.True(t, tok.IsCmpOp(), tok.String())
		} else {
			assert.False(t, tok.IsCmpOp(), tok.String())
		}
	}
}

func TestIsLogOp(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tok := Token(i)
		if tok == AND || tok == OR || tok == NOT {
			assert.True(t, tok.IsLogOp(), tok.String())
		} else {
			assert.False(t, tok.IsLogOp(), tok.String())
		}
	}
}

func TestIsLiteral(t *testing.T) {
	for i := 0; i < len(tokens); i++ {
		tok := Token(i)
		if tok == TRUE || tok == FALSE || tok == NUMBER || tok == STRING || tok == IDENTIFIER {
			assert.True(t, tok.IsLiteral(), tok.String())
		} else {
			assert.False(t, tok.IsLiteral(), tok.String())
		}
	}
}
