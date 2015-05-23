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
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scanner"
	"github.com/tboronczyk/kiwi/token"
)

type Parser struct {
	token   token.Token
	value   string
	scanner scanner.Scanner
}

func NewParser() *Parser {
	return &Parser{}
}

func (p Parser) match(tkn token.Token) bool {
	return p.token == tkn
}

func (p *Parser) advance() {
	for {
		p.token, p.value = p.scanner.Scan()
		if p.token != token.COMMENT {
			break
		}
	}
}

func (p Parser) expected(value interface{}) string {
	var str string
	switch value.(type) {
	case token.Token:
		str = value.(token.Token).String()
	default:
		str = value.(string)
	}
	return "Expected " + str + " but saw " + p.token.String()
}

func (p *Parser) InitScanner(scnr scanner.Scanner) {
	p.scanner = scnr
	p.advance()
}

func (p *Parser) Parse() ast.Node {
	if p.token == token.EOF {
		return ast.Operator{Op: p.token}
	}
	return p.parseStmt()
}

func (p *Parser) parseExpr() ast.Node {
	n := p.parseRelation()
	if !p.token.IsLogOp() {
		return n
	}

	node := ast.Operator{Op: p.token}
	p.advance()

	node.Left = n
	node.Right = p.parseExpr()
	return node
}

func (p *Parser) parseRelation() ast.Node {
	n := p.parseSimpleExpr()
	if !p.token.IsCmpOp() {
		return n
	}

	node := ast.Operator{Op: p.token}
	p.advance()

	node.Left = n
	node.Right = p.parseRelation()
	return node
}

func (p *Parser) parseSimpleExpr() ast.Node {
	n := p.parseTerm()
	if !p.token.IsAddOp() {
		return n
	}

	node := ast.Operator{Op: p.token}
	p.advance()

	node.Left = n
	node.Right = p.parseSimpleExpr()
	return node
}

func (p *Parser) parseTerm() ast.Node {
	n := p.parseFactor()
	if !p.token.IsMulOp() {
		return n
	}

	node := ast.Operator{Op: p.token}
	p.advance()

	node.Left = n
	node.Right = p.parseTerm()
	return node
}

func (p *Parser) parseFactor() ast.Node {
	if p.match(token.LPAREN) {
		p.advance()

		node := p.parseExpr()
		if !p.match(token.RPAREN) {
			panic(p.expected(token.RPAREN))
		}
		p.advance()

		return node
	}
	if p.token.IsExprOp() {
		node := ast.Operator{Op: p.token}
		p.advance()

		n := p.parseFactor()
		node.Left = n
		return node
	}
	return p.parseTerminal()
}

func (p *Parser) parseTerminal() ast.Node {
	if p.token == token.TRUE || p.token == token.FALSE ||
		p.token == token.NUMBER || p.token == token.STRING {
		node := ast.Literal{Type: p.token, Value: p.value}
		p.advance()
		return node
	}

	n := p.parseIdentifier()
	if p.token != token.LPAREN {
		return n
	}
	node := ast.FuncCall{Name: n.Value}
	node.Body = p.parseParenExprList()
	return node
}

func (p *Parser) parseParenExprList() ast.Node {
	// method only called when already p.token == token.LPAREN
	/*
		if p.token != token.LPAREN {
			panic(p.expected(token.LPAREN))
		}
	*/
	p.advance()

	if p.token == token.RPAREN {
		p.advance()
		return nil
	}

	node := p.parseExprList()

	if p.token != token.RPAREN {
		panic(p.expected(token.RPAREN))
	}
	p.advance()
	return node
}

func (p *Parser) parseExprList() ast.Node {
	n := p.parseExpr()
	if p.token != token.COMMA {
		return n
	}

	node := ast.List{}
	node.Node = n
	for p.token == token.COMMA {
		p.advance()

		next := ast.List{}
		next.Next = node
		next.Node = p.parseExpr()
		node = next
	}
	return node
}

func (p *Parser) parseStmt() ast.Node {
	switch p.token {
	case token.IF:
		return p.parseIfStmt()
	case token.WHILE:
		return p.parseWhileStmt()
	case token.IDENTIFIER:
		return p.parseAssignOrCallStmt()
	}
	panic(p.expected(
		token.IF.String() + ", " + token.WHILE.String() + ", or " +
			token.IDENTIFIER.String()))
}

func (p *Parser) parseIfStmt() ast.If {
	node := ast.If{}
	p.advance()

	n := p.parseExpr()
	node.Condition = n

	node.Body = p.parseBraceStmtList()
	return node
}

func (p *Parser) parseBraceStmtList() ast.List {
	if p.token != token.LBRACE {
		panic(p.expected(token.LBRACE))
	}
	p.advance()

	node := p.parseStmtList()
	if p.token != token.RBRACE {
		panic(p.expected(token.RBRACE))
	}
	p.advance()

	return node
}

func (p *Parser) parseStmtList() ast.List {
	node := ast.List{}
	for p.token.IsStmtKeyword() || p.token == token.IDENTIFIER {
		n := p.parseStmt()
		node.Node = n
		next := ast.List{}
		next.Next = node
		node = next
	}
	return node
}

func (p *Parser) parseWhileStmt() ast.While {
	node := ast.While{}
	p.advance()

	n := p.parseExpr()
	node.Condition = n
	node.Body = p.parseBraceStmtList()
	return node
}

func (p *Parser) parseAssignOrCallStmt() ast.Node {
	n := p.parseIdentifier()
	if p.token == token.ASSIGN {
		node := ast.Operator{Op: p.token}
		p.advance()
		node.Left = n
		node.Right = p.parseExpr()
		if p.token != token.SEMICOLON {
			panic(p.expected(token.SEMICOLON))
		}
		p.advance()
		return node
	}
	if p.token == token.LPAREN {
		node := ast.FuncCall{Name: n.Value}
		node.Body = p.parseParenExprList()
		if p.token != token.SEMICOLON {
			panic(p.expected(token.SEMICOLON))
		}
		p.advance()
		return node
	}
	panic(p.expected(
		token.ASSIGN.String() + " or " + token.LPAREN.String()))
}

func (p *Parser) parseIdentifier() ast.Literal {
	node := ast.Literal{Type: p.token, Value: p.value}
	if p.token != token.IDENTIFIER {
		panic(p.expected(token.IDENTIFIER))
	}
	p.advance()
	return node
}
