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

package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/token"
	"testing"
)

type tokenPair struct {
	token token.Token
	value string
}

type mockScanner struct {
	i  uint8
	tp []tokenPair
}

func NewMockScanner() *mockScanner {
	return &mockScanner{i: 0}
}

func (s *mockScanner) reset(pairs []tokenPair) {
	s.i = 0
	s.tp = pairs
}

func (s *mockScanner) Scan() (token.Token, string) {
	t := s.tp[s.i].token
	v := s.tp[s.i].value
	s.i++
	return t, v
}

func TestSkipComment(t *testing.T) {
	s := NewMockScanner()
	s.reset([]tokenPair{
		{token.COMMENT, "//"},
		{token.STRING, ""},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	assert.Equal(t, token.STRING, p.tkn)
}

func TestParseStmtList(t *testing.T) {
	s := NewMockScanner()
	// a := true ; b := false ; c := 1 ;
	s.reset([]tokenPair{
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ":="},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "b"},
		{token.ASSIGN, ":="},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "c"},
		{token.ASSIGN, ":="},
		{token.NUMBER, "1"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseStmtList()
	assert.Equal(t, token.NUMBER, node.Children[0].Children[1].Children[1].Token)
	assert.Equal(t, token.FALSE, node.Children[0].Children[0].Children[1].Children[1].Token)
	assert.Equal(t, token.TRUE, node.Children[0].Children[0].Children[0].Children[1].Children[1].Token)
	assert.Nil(t, err)
}

func TestParseExprList(t *testing.T) {
	s := NewMockScanner()
	// a, 1, "foo"
	s.reset([]tokenPair{
		{token.IDENTIFIER, "a"},
		{token.COMMA, ","},
		{token.NUMBER, "1"},
		{token.COMMA, ","},
		{token.STRING, "foo"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseExprList()
	assert.Equal(t, token.IDENTIFIER, node.Children[0].Children[0].Children[1].Token)
	assert.Equal(t, token.NUMBER, node.Children[0].Children[1].Token)
	assert.Equal(t, token.STRING, node.Children[1].Token)
	assert.Nil(t, err)
}

func TestParseIfStmt(t *testing.T) {
	s := NewMockScanner()
	// if true { a := false; }
	s.reset([]tokenPair{
		{token.IF, "if"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ":="},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseStmt()
	assert.Equal(t, token.TRUE, node.Children[0].Token)
	assert.Equal(t, token.IF, node.Token)
	assert.Equal(t, token.ASSIGN, node.Children[1].Token)
	assert.Nil(t, err)
}

func TestParseWhileStmt(t *testing.T) {
	s := NewMockScanner()
	// while true { a := false; }
	s.reset([]tokenPair{
		{token.WHILE, "while"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ":="},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseStmt()
	assert.Equal(t, token.TRUE, node.Children[0].Token)
	assert.Equal(t, token.WHILE, node.Token)
	assert.Equal(t, token.ASSIGN, node.Children[1].Token)
	assert.Nil(t, err)
}

func TestParseAssignStmt(t *testing.T) {
	s := NewMockScanner()
	// a := 4;
	s.reset([]tokenPair{
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ":="},
		{token.NUMBER, "4"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseStmt()
	assert.Equal(t, token.IDENTIFIER, node.Children[0].Token)
	assert.Equal(t, token.ASSIGN, node.Token)
	assert.Equal(t, token.NUMBER, node.Children[1].Token)
	assert.Nil(t, err)
}

func TestParseExpr(t *testing.T) {
	s := NewMockScanner()
	// true && true
	s.reset([]tokenPair{
		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.TRUE, "true"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseExpr()
	assert.Equal(t, token.TRUE, node.Children[0].Token)
	assert.Equal(t, token.AND, node.Token)
	assert.Equal(t, token.TRUE, node.Children[1].Token)
	assert.Nil(t, err)
}

func TestParseExprError(t *testing.T) {
	s := NewMockScanner()
	// true && EOF
	s.reset([]tokenPair{
		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.ParseExpr()
	assert.NotNil(t, err)
}

func TestParseRelation(t *testing.T) {
	s := NewMockScanner()
	s.reset([]tokenPair{
		// 2 < 4
		{token.NUMBER, "2"},
		{token.LESS, "<"},
		{token.NUMBER, "4"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseRelation()
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, node.Children[0].Token)
	assert.Equal(t, token.LESS, node.Token)
	assert.Equal(t, token.NUMBER, node.Children[1].Token)
}

func TestParseRelationError(t *testing.T) {
	s := NewMockScanner()
	// 2 < EOF
	s.reset([]tokenPair{
		{token.NUMBER, "2"},
		{token.LESS, "<"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.ParseRelation()
	assert.NotNil(t, err)
}

func TestParseSimpleExpr(t *testing.T) {
	s := NewMockScanner()
	// 2 + 4
	s.reset([]tokenPair{
		{token.NUMBER, "2"},
		{token.ADD, "+"},
		{token.NUMBER, "4"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseSimpleExpr()
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, node.Children[0].Token)
	assert.Equal(t, token.ADD, node.Token)
	assert.Equal(t, token.NUMBER, node.Children[1].Token)
}

func TestParseSimpleExprError(t *testing.T) {
	s := NewMockScanner()
	// 2 + EOF
	s.reset([]tokenPair{
		{token.NUMBER, "2"},
		{token.ADD, "+"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.ParseSimpleExpr()
	assert.NotNil(t, err)
}

func TestParseTerm(t *testing.T) {
	s := NewMockScanner()
	// 2 * 4
	s.reset([]tokenPair{
		{token.NUMBER, "2"},
		{token.MULTIPLY, "*"},
		{token.NUMBER, "4"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseTerm()
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, node.Children[0].Token)
	assert.Equal(t, token.MULTIPLY, node.Token)
	assert.Equal(t, token.NUMBER, node.Children[1].Token)
}

func TestParseTermError(t *testing.T) {
	s := NewMockScanner()
	// 2 * EOF
	s.reset([]tokenPair{
		{token.NUMBER, "2"},
		{token.MULTIPLY, "*"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.ParseTerm()
	assert.NotNil(t, err)
}

func TestParseFactorParens(t *testing.T) {
	s := NewMockScanner()
	// ( X )
	s.reset([]tokenPair{
		{token.LPAREN, "("},
		{token.IDENTIFIER, "X"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseFactor()
	assert.Nil(t, err)
	assert.Equal(t, token.IDENTIFIER, node.Token)
}

func TestParseFactorParensExprError(t *testing.T) {
	s := NewMockScanner()
	// ( EOF
	s.reset([]tokenPair{
		{token.LPAREN, "("},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.ParseFactor()
	assert.NotNil(t, err)
}

func TestParseFactorParensCloseError(t *testing.T) {
	s := NewMockScanner()
	// ( X EOF
	s.reset([]tokenPair{
		{token.LPAREN, "("},
		{token.IDENTIFIER, "X"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.ParseFactor()
	assert.NotNil(t, err)
}

func TestParseFactorNot(t *testing.T) {
	s := NewMockScanner()
	// ~ true
	s.reset([]tokenPair{
		{token.NOT, "~"},
		{token.TRUE, "true"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseFactor()
	assert.Nil(t, err)
	assert.Equal(t, token.TRUE, node.Children[0].Token)
	assert.Equal(t, token.NOT, node.Token)
}

func TestParseFactorNotError(t *testing.T) {
	s := NewMockScanner()
	// ~ EOF
	s.reset([]tokenPair{
		{token.NOT, "~"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.ParseFactor()
	assert.NotNil(t, err)
}

func TestParseFactorSigned(t *testing.T) {
	s := NewMockScanner()
	// -1
	s.reset([]tokenPair{
		{token.SUBTRACT, "-"},
		{token.NUMBER, "1"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseFactor()
	assert.Nil(t, err)
	assert.Equal(t, token.SUBTRACT, node.Token)
	assert.Equal(t, token.NUMBER, node.Children[0].Token)
}

func TestParseFactorSignedError(t *testing.T) {
	s := NewMockScanner()
	// -EOF
	s.reset([]tokenPair{
		{token.SUBTRACT, "-"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.ParseFactor()
	assert.NotNil(t, err)
}

func TestParseTerminalFuncCall(t *testing.T) {
	s := NewMockScanner()
	// f()
	s.reset([]tokenPair{
		{token.IDENTIFIER, "f"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseTerminal()
	assert.Equal(t, token.IDENTIFIER, node.Children[0].Token)
	assert.Nil(t, err)
}

func TestParseTerminalFuncCallWithArgs(t *testing.T) {
	s := NewMockScanner()
	// f(1,2)
	s.reset([]tokenPair{
		{token.IDENTIFIER, "f"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "a"},
		{token.COMMA, ","},
		{token.NUMBER, "1"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseTerminal()
	assert.Equal(t, token.IDENTIFIER, node.Children[0].Token)
	assert.Equal(t, token.IDENTIFIER, node.Children[1].Children[0].Children[1].Token)
	assert.Equal(t, token.NUMBER, node.Children[1].Children[1].Token)
	assert.Nil(t, err)
}

func TestParseFuncCallStmt(t *testing.T) {
	s := NewMockScanner()
	// f(1,2)
	s.reset([]tokenPair{
		{token.IDENTIFIER, "f"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "a"},
		{token.COMMA, ","},
		{token.NUMBER, "1"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseStmt()
	assert.Equal(t, token.IDENTIFIER, node.Children[0].Token)
	assert.Equal(t, token.IDENTIFIER, node.Children[1].Children[0].Children[1].Token)
	assert.Equal(t, token.NUMBER, node.Children[1].Children[1].Token)
	assert.Nil(t, err)
}

func TestParseFullExpression(t *testing.T) {
	s := NewMockScanner()
	// ~((-1 < 0) && (X > 2 + 4 * 1))
	s.reset([]tokenPair{
		{token.NOT, "~"},
		{token.LPAREN, "("},
		{token.LPAREN, "("},
		{token.SUBTRACT, "-"},
		{token.NUMBER, "1"},
		{token.LESS, "<"},
		{token.NUMBER, "0"},
		{token.RPAREN, ")"},
		{token.AND, "&&"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "X"},
		{token.GREATER, ">"},
		{token.NUMBER, "2"},
		{token.ADD, "+"},
		{token.NUMBER, "4"},
		{token.MULTIPLY, "*"},
		{token.NUMBER, "1"},
		{token.RPAREN, ")"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.ParseFactor()
	assert.Nil(t, err)
	assert.Equal(t, token.NOT, node.Token)
	assert.Equal(t, token.AND, node.Children[0].Token)
	assert.Equal(t, token.LESS, node.Children[0].Children[0].Token)
	assert.Equal(t, token.SUBTRACT, node.Children[0].Children[0].Children[0].Token)
	assert.Equal(t, token.NUMBER, node.Children[0].Children[0].Children[0].Children[0].Token)
	assert.Equal(t, token.NUMBER, node.Children[0].Children[0].Children[1].Token)
	assert.Equal(t, token.GREATER, node.Children[0].Children[1].Token)
	assert.Equal(t, token.IDENTIFIER, node.Children[0].Children[1].Children[0].Token)
	assert.Equal(t, token.ADD, node.Children[0].Children[1].Children[1].Token)
	assert.Equal(t, token.NUMBER, node.Children[0].Children[1].Children[1].Children[0].Token)
	assert.Equal(t, token.MULTIPLY, node.Children[0].Children[1].Children[1].Children[1].Token)
	assert.Equal(t, token.NUMBER, node.Children[0].Children[1].Children[1].Children[1].Children[0].Token)
	assert.Equal(t, token.NUMBER, node.Children[0].Children[1].Children[1].Children[1].Children[1].Token)
}
