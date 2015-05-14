/*
 * Copyright (c) 2012, 2015 Timothy Boronczyk
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
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

package scanner

import (
	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/token"
	"strings"
	"testing"
)

func TestScanSimpleTokens(t *testing.T) {
	str := "+ - * / % := : = < <= > >= && & || | ~ ~= ( ) { } ; , ?"
	s := NewScanner(strings.NewReader(str))

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
		{token.MALFORMED, ":"},
		{token.EQUAL, "="},
		{token.LESS, "<"},
		{token.LESS_EQ, "<="},
		{token.GREATER, ">"},
		{token.GREATER_EQ, ">="},
		{token.AND, "&&"},
		{token.MALFORMED, "&"},
		{token.OR, "||"},
		{token.MALFORMED, "|"},
		{token.NOT, "~"},
		{token.NOT_EQUAL, "~="},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.COMMA, ","},
		{token.UNKNOWN, "?"},
		{token.EOF, ""},
	}

	for _, expected := range tokens {
		actual1, actual2 := s.Scan()
		assert.Equal(t, expected.token, actual1)
		assert.Equal(t, expected.value, actual2)
	}
}

func TestScanIdentifiers(t *testing.T) {
	str := "if while true false `if ident"
	s := NewScanner(strings.NewReader(str))

	tokens := []struct {
		token token.Token
		value string
	}{
		{token.IF, "if"},
		{token.WHILE, "while"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.IDENTIFIER, "if"},
		{token.IDENTIFIER, "ident"},
	}

	for _, expected := range tokens {
		actual1, actual2 := s.Scan()
		assert.Equal(t, expected.token, actual1)
		assert.Equal(t, expected.value, actual2)
	}
}

func TestScanStrings(t *testing.T) {
	str := "\"abc\"" +
		"\"\"" +
		"\"\\\"\"" +
		"\"broken"
	s := NewScanner(strings.NewReader(str))

	tokens := []struct {
		token token.Token
		value string
	}{
		{token.STRING, "abc"},
		{token.STRING, ""},
		{token.STRING, "\\\""},
		{token.MALFORMED, "broken"},
	}

	for _, expected := range tokens {
		actual1, actual2 := s.Scan()
		assert.Equal(t, expected.token, actual1)
		assert.Equal(t, expected.value, actual2)
	}
}

func TestScanLineComments(t *testing.T) {
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
		assert.Equal(t, token.COMMENT, actual1)
		assert.Equal(t, expected.value, actual2)
	}
}

func TestScanMultiLineComments(t *testing.T) {
	str := "/**/" +
		"/* a /* nested */ comment */" +
		"/* broken"
	s := NewScanner(strings.NewReader(str))

	tokens := []struct {
		token token.Token
		value string
	}{
		{token.COMMENT, "/**/"},
		{token.COMMENT, "/* a /* nested */ comment */"},
		{token.MALFORMED, "/* broken"},
	}

	for _, expected := range tokens {
		actual1, actual2 := s.Scan()
		assert.Equal(t, expected.token, actual1)
		assert.Equal(t, expected.value, actual2)
	}
}

func TestScanNumbers(t *testing.T) {
	str := "123"
	s := NewScanner(strings.NewReader(str))

	tokens := []struct {
		token token.Token
		value string
	}{
		{token.NUMBER, "123"},
	}

	for _, expected := range tokens {
		actual1, actual2 := s.Scan()
		assert.Equal(t, expected.token, actual1)
		assert.Equal(t, expected.value, actual2)
	}
}
