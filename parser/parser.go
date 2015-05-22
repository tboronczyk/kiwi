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
	"errors"
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

func (p Parser) expectedError(tkn token.Token) error {
	return errors.New("Expected " + tkn.String() + " but saw " + p.token.String())
}

func (p Parser) expectedErrorStr(str string) error {
	return errors.New("Expected " + str + " but saw " + p.token.String())
}

func (p *Parser) InitScanner(scnr scanner.Scanner) {
	p.scanner = scnr
	p.advance()
}

func (p *Parser) Parse() (ast.Node, error) {
	if p.token == token.EOF {
		return ast.NewOperator(p.token), nil
	}
	return p.parseStmt()
}

func (p *Parser) parseExpr() (ast.Node, error) {
	n, err := p.parseRelation()
	if err != nil || !p.token.IsLogOp() {
		return n, err
	}

	node := ast.NewOperator(p.token)
	p.advance()

	node.Left = n
	node.Right, err = p.parseExpr()
	return node, err
}

func (p *Parser) parseRelation() (ast.Node, error) {
	n, err := p.parseSimpleExpr()
	if err != nil || !p.token.IsCmpOp() {
		return n, err
	}

	node := ast.NewOperator(p.token)
	p.advance()

	node.Left = n
	node.Right, err = p.parseRelation()
	return node, err
}

func (p *Parser) parseSimpleExpr() (ast.Node, error) {
	n, err := p.parseTerm()
	if err != nil || !p.token.IsAddOp() {
		return n, err
	}

	node := ast.NewOperator(p.token)
	p.advance()

	node.Left = n
	node.Right, err = p.parseSimpleExpr()
	return node, err
}

func (p *Parser) parseTerm() (ast.Node, error) {
	n, err := p.parseFactor()
	if err != nil || !p.token.IsMulOp() {
		return n, err
	}

	node := ast.NewOperator(p.token)
	p.advance()

	node.Left = n
	node.Right, err = p.parseTerm()
	return node, err
}

func (p *Parser) parseFactor() (ast.Node, error) {
	if p.match(token.LPAREN) {
		p.advance()

		node, err := p.parseExpr()
		if err != nil {
			return node, err
		}

		if !p.match(token.RPAREN) {
			return node, p.expectedError(token.RPAREN)
		}
		p.advance()

		return node, err
	}
	if p.token.IsExprOp() {
		node := ast.NewOperator(p.token)
		p.advance()

		n, err := p.parseFactor()
		node.Left = n
		return node, err
	}
	return p.parseTerminal()
}

func (p *Parser) parseTerminal() (ast.Node, error) {
	if p.token == token.TRUE || p.token == token.FALSE ||
		p.token == token.NUMBER || p.token == token.STRING {
		node := ast.NewLiteral(p.token, p.value)
		p.advance()
		return node, nil
	}

	n, err := p.parseIdentifier()
	if err != nil || p.token != token.LPAREN {
		return n, err
	}
	node := ast.NewFuncCall(n.Value)
	node.Body, err = p.parseParenExprList()
	return node, err
}

func (p *Parser) parseParenExprList() (ast.Node, error) {
	// method only called when already p.token == token.LPAREN
	/*
		if p.token != token.LPAREN {
			return &ast.Node{}, p.expectedError(token.LPAREN)
		}
	*/
	p.advance()

	if p.token == token.RPAREN {
		p.advance()
		return nil, nil
	}

	node, err := p.parseExprList()
	if err != nil {
		return node, err
	}

	if p.token != token.RPAREN {
		return node, p.expectedError(token.RPAREN)
	}
	p.advance()
	return node, nil
}

func (p *Parser) parseExprList() (ast.Node, error) {
	n, err := p.parseExpr()
	if err != nil || p.token != token.COMMA {
		return n, err
	}

	node := ast.NewList()
	node.Node = n
	for p.token == token.COMMA {
		p.advance()

		next := ast.NewList()
		next.Next = node
		next.Node, err = p.parseExpr()
		if err != nil {
			return next, err
		}
		node = next
	}
	return node, err
}

func (p *Parser) parseStmt() (ast.Node, error) {
	switch p.token {
	case token.IF:
		return p.parseIfStmt()
	case token.WHILE:
		return p.parseWhileStmt()
	case token.IDENTIFIER:
		return p.parseAssignOrCallStmt()
	}
	return nil, p.expectedErrorStr(
		token.IF.String() + ", " + token.WHILE.String() + ", or " +
			token.IDENTIFIER.String())
}

func (p *Parser) parseIfStmt() (ast.Node, error) {
	node := ast.NewIf()
	p.advance()

	n, err := p.parseExpr()
	node.Condition = n
	if err != nil {
		return node, err
	}

	node.Body, err = p.parseBraceStmtList()
	return node, err
}

func (p *Parser) parseBraceStmtList() (ast.Node, error) {
	if p.token != token.LBRACE {
		return nil, p.expectedError(token.LBRACE)
	}
	p.advance()

	node, err := p.parseStmtList()
	if err != nil {
		return node, err
	}

	if p.token != token.RBRACE {
		return node, p.expectedError(token.RBRACE)
	}
	p.advance()

	return node, nil
}

func (p *Parser) parseStmtList() (ast.Node, error) {
	node := ast.NewList()
	for p.token.IsStmtKeyword() || p.token == token.IDENTIFIER {
		n, err := p.parseStmt()
		node.Node = n
		if err != nil {
			return node, err
		}
		next := ast.NewList()
		next.Next = node
		node = next
	}
	return node, nil
}

func (p *Parser) parseWhileStmt() (ast.Node, error) {
	node := ast.NewWhile()
	p.advance()

	n, err := p.parseExpr()
	node.Condition = n
	if err != nil {
		return node, err
	}

	node.Body, err = p.parseBraceStmtList()
	return node, err
}

func (p *Parser) parseAssignOrCallStmt() (ast.Node, error) {
	n, err := p.parseIdentifier()
	if p.token == token.ASSIGN {
		node := ast.NewOperator(p.token)
		p.advance()
		node.Left = n
		node.Right, err = p.parseExpr()
		if err != nil {
			return node, err
		}
		if p.token != token.SEMICOLON {
			return node, p.expectedError(token.SEMICOLON)
		}
		p.advance()
		return node, nil
	}
	if p.token == token.LPAREN {
		node := ast.NewFuncCall(n.Value)
		node.Body, err = p.parseParenExprList()
		if err != nil {
			return node, err
		}
		if p.token != token.SEMICOLON {
			return node, p.expectedError(token.SEMICOLON)
		}
		p.advance()
		return node, nil
	}
	return nil, p.expectedErrorStr(
		token.ASSIGN.String() + " or " + token.LPAREN.String())
}

func (p *Parser) parseIdentifier() (ast.Literal, error) {
	node := ast.NewLiteral(p.token, p.value)
	if p.token != token.IDENTIFIER {
		return node, p.expectedError(token.IDENTIFIER)
	}
	p.advance()
	return node, nil
}
