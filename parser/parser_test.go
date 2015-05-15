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

	assert.Equal(t, token.STRING, p.token)
}

func TestParseIdentifier(t *testing.T) {
	s := NewMockScanner()
	// foo
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseIdentifier()
	assert.Equal(t, token.IDENTIFIER, node.Token)
	assert.Nil(t, err)
}

func TestParseIdentifierError(t *testing.T) {
	s := NewMockScanner()
	// foo
	s.reset([]tokenPair{
		{token.NUMBER, "42"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseIdentifier()
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

	node, err := p.parseTerm()
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, node.Children[0].Token)
	assert.Equal(t, token.MULTIPLY, node.Token)
	assert.Equal(t, token.NUMBER, node.Children[1].Token)
}

func TestParseTermError(t *testing.T) {
	s := NewMockScanner()
	// 2 *
	s.reset([]tokenPair{
		{token.NUMBER, "2"},
		{token.MULTIPLY, "*"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseTerm()
	assert.NotNil(t, err)
}

func TestParseSimpleExpr(t *testing.T) {
	s := NewMockScanner()
	// 42 + 73
	s.reset([]tokenPair{
		{token.NUMBER, "42"},
		{token.ADD, "+"},
		{token.NUMBER, "73"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseSimpleExpr()
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, node.Children[0].Token)
	assert.Equal(t, token.ADD, node.Token)
	assert.Equal(t, token.NUMBER, node.Children[1].Token)
}

func TestParseSimpleExprError(t *testing.T) {
	s := NewMockScanner()
	// 42 +
	s.reset([]tokenPair{
		{token.NUMBER, "42"},
		{token.ADD, "+"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseSimpleExpr()
	assert.NotNil(t, err)
}

func TestParseRelation(t *testing.T) {
	s := NewMockScanner()
	s.reset([]tokenPair{
		// 1776 < 2001
		{token.NUMBER, "1776"},
		{token.LESS, "<"},
		{token.NUMBER, "2001"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseRelation()
	assert.Nil(t, err)
	assert.Equal(t, token.NUMBER, node.Children[0].Token)
	assert.Equal(t, token.LESS, node.Token)
	assert.Equal(t, token.NUMBER, node.Children[1].Token)
}

func TestParseRelationError(t *testing.T) {
	s := NewMockScanner()
	// 1776 <
	s.reset([]tokenPair{
		{token.NUMBER, "1776"},
		{token.LESS, "<"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseRelation()
	assert.NotNil(t, err)
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

	node, err := p.parseExpr()
	assert.Equal(t, token.TRUE, node.Children[0].Token)
	assert.Equal(t, token.AND, node.Token)
	assert.Equal(t, token.TRUE, node.Children[1].Token)
	assert.Nil(t, err)
}

func TestParseExprError(t *testing.T) {
	s := NewMockScanner()
	// true &&
	s.reset([]tokenPair{
		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseExpr()
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

	node, err := p.parseFactor()
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

	_, err := p.parseFactor()
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

	_, err := p.parseFactor()
	assert.NotNil(t, err)
}

func TestParseFactorSigned(t *testing.T) {
	s := NewMockScanner()
	// -101
	s.reset([]tokenPair{
		{token.SUBTRACT, "-"},
		{token.NUMBER, "101"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseFactor()
	assert.Nil(t, err)
	assert.Equal(t, token.SUBTRACT, node.Token)
	assert.Equal(t, token.NUMBER, node.Children[0].Token)
}

func TestParseTerminalIdentifier(t *testing.T) {
	s := NewMockScanner()
	// foo
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseTerminal()
	assert.Equal(t, token.IDENTIFIER, node.Token)
	assert.Nil(t, err)
}

func TestParseTerminalFuncCall(t *testing.T) {
	s := NewMockScanner()
	// foo()
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseTerminal()
	assert.Equal(t, token.IDENTIFIER, node.Children[0].Token)
	assert.Nil(t, err)
}

func TestParseTerminalFuncCallWithArgs(t *testing.T) {
	s := NewMockScanner()
	// foo(x, 42, "hello")
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.NUMBER, "42"},
		{token.COMMA, ","},
		{token.STRING, "hello"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseTerminal()
	assert.Equal(t, token.IDENTIFIER, node.Children[0].Token)
	assert.Equal(t, token.IDENTIFIER, node.Children[1].Children[0].Children[0].Children[1].Token)
	assert.Equal(t, token.NUMBER, node.Children[1].Children[0].Children[1].Token)
	assert.Equal(t, token.STRING, node.Children[1].Children[1].Token)
	assert.Nil(t, err)
}

func TestParseTerminalFuncCallWithArgsExprError(t *testing.T) {
	s := NewMockScanner()
	// foo(bar, 2001 < ,
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "bar"},
		{token.COMMA, ","},
		{token.NUMBER, "2001"},
		{token.LESS, "<"},
		{token.COMMA, ","},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseTerminal()
	assert.NotNil(t, err)
}

func TestParseTerminalFuncCallExprListError(t *testing.T) {
	s := NewMockScanner()
	// foo(bar 2001
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "bar"},
		{token.IDENTIFIER, "2001"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseTerminal()
	assert.NotNil(t, err)
}

func TestParseBraceStmtListEmpty(t *testing.T) {
	s := NewMockScanner()
	// { }
	s.reset([]tokenPair{
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseBraceStmtList()
	assert.Nil(t, node.Children[0])
	assert.Nil(t, err)
}

func TestParseBraceStmtList(t *testing.T) {
	s := NewMockScanner()
	// { a := true; b := false; }
	s.reset([]tokenPair{
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ":="},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "b"},
		{token.ASSIGN, ":="},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseBraceStmtList()
	assert.Equal(t, token.ASSIGN, node.Children[0].Children[0].Children[1].Token)
	assert.Equal(t, token.ASSIGN, node.Children[0].Children[1].Token)
	assert.Nil(t, err)
}

func TestParseBraceStmtListStmtError(t *testing.T) {
	s := NewMockScanner()
	// { a := true; b false
	s.reset([]tokenPair{
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ":="},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "b"},
		{token.FALSE, "false"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseBraceStmtList()
	assert.NotNil(t, err)
}

func TestParseBraceStmtListBraceError(t *testing.T) {
	s := NewMockScanner()
	// { a := true;
	s.reset([]tokenPair{
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ":="},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseBraceStmtList()
	assert.NotNil(t, err)
}

func TestParseIfStmt(t *testing.T) {
	s := NewMockScanner()
	// if a = true { b := false; }
	s.reset([]tokenPair{
		{token.IF, "if"},
		{token.IDENTIFIER, "a"},
		{token.EQUAL, "="},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "b"},
		{token.ASSIGN, ":="},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseStmt()
	assert.Equal(t, token.EQUAL, node.Children[0].Token)
	assert.Equal(t, token.ASSIGN, node.Children[1].Children[0].Children[1].Token)
	assert.Nil(t, err)
}

func TestParseIfStmtExprError(t *testing.T) {
	s := NewMockScanner()
	// if a = {
	s.reset([]tokenPair{
		{token.IF, "if"},
		{token.IDENTIFIER, "a"},
		{token.EQUAL, "="},
		{token.LBRACE, "{"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseStmt()
	assert.NotNil(t, err)
}

func TestParseIfStmtBraceError(t *testing.T) {
	s := NewMockScanner()
	// if a = true b :=
	s.reset([]tokenPair{
		{token.IF, "if"},
		{token.IDENTIFIER, "a"},
		{token.EQUAL, "="},
		{token.TRUE, "true"},
		{token.IDENTIFIER, "b"},
		{token.ASSIGN, ":="},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseStmt()
	assert.NotNil(t, err)
}

func TestParseWhileStmt(t *testing.T) {
	s := NewMockScanner()
	// while a = true { b := false; }
	s.reset([]tokenPair{
		{token.WHILE, "while"},
		{token.IDENTIFIER, "a"},
		{token.EQUAL, "="},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "b"},
		{token.ASSIGN, ":="},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseStmt()
	assert.Equal(t, token.EQUAL, node.Children[0].Token)
	assert.Equal(t, token.ASSIGN, node.Children[1].Children[0].Children[1].Token)
	assert.Nil(t, err)
}

func TestParseWhileStmtExprError(t *testing.T) {
	s := NewMockScanner()
	// while a = {
	s.reset([]tokenPair{
		{token.WHILE, "if"},
		{token.IDENTIFIER, "a"},
		{token.EQUAL, "="},
		{token.LBRACE, "{"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseStmt()
	assert.NotNil(t, err)
}

func TestParseStmtError(t *testing.T) {
	s := NewMockScanner()
	// ;
	s.reset([]tokenPair{
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseStmt()
	assert.NotNil(t, err)
}

func TestParseAssignStmt(t *testing.T) {
	s := NewMockScanner()
	// a := 2 + 4;
	s.reset([]tokenPair{
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ":="},
		{token.NUMBER, "2"},
		{token.ADD, "+"},
		{token.NUMBER, "4"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseStmt()
	assert.Equal(t, token.ASSIGN, node.Token)
	assert.Equal(t, token.IDENTIFIER, node.Children[0].Token)
	assert.Equal(t, token.ADD, node.Children[1].Token)
	assert.Nil(t, err)
}

func TestParseAssignSmtExprError(t *testing.T) {
	s := NewMockScanner()
	// a := 2 +
	s.reset([]tokenPair{
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ":="},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "2"},
		{token.ADD, "+"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseStmt()
	assert.NotNil(t, err)
}

func TestFuncCallStmt(t *testing.T) {
	s := NewMockScanner()
	// foo();
	s.reset([]tokenPair{
		{token.IDENTIFIER, "foo"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	node, err := p.parseStmt()
	assert.Equal(t, token.IDENTIFIER, node.Children[0].Token)
	assert.Nil(t, node.Children[1].Children)
	assert.Nil(t, err)
}

func TestStmtSemicolonError(t *testing.T) {
	s := NewMockScanner()
	// a := true
	s.reset([]tokenPair{
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ":="},
		{token.TRUE, "true"},
		{token.EOF, ""},
	})
	p := NewParser()
	p.InitScanner(s)

	_, err := p.parseStmt()
	assert.NotNil(t, err)
}
