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
	tkn token.Token
	val string
	s   scanner.Scanner
}

func NewParser() *Parser {
	return &Parser{}
}

func (p Parser) match(t token.Token) bool {
	return p.tkn == t
}

func (p *Parser) advance() {
	for {
		p.tkn, p.val = p.s.Scan()
		if p.tkn != token.COMMENT {
			break
		}
	}
}

func (p *Parser) InitScanner(s scanner.Scanner) {
	p.s = s
	p.advance()
}

func (p *Parser) Parse() (*ast.Node, error) {
	if p.tkn == token.EOF {
		return ast.NewNode(p.tkn, p.val, 2), nil
	}
	return p.ParseStmt()
}

func (p *Parser) ParseStmtList() (*ast.Node, error) {
	var err error
	node := ast.NewNode(nil, nil, 2)
	for p.tkn == token.IF || p.tkn == token.WHILE || p.tkn == token.IDENTIFIER {
		node.Children[1], err = p.ParseStmt()
		if err != nil {
			return node, err
		}
		n := ast.NewNode(nil, nil, 2)
		n.Children[0] = node
		node = n
	}
	return node, err
}

func (p *Parser) ParseExprList() (*ast.Node, error) {
	var err error = nil
	n := ast.NewNode(nil, nil, 2)
	n.Children[1], err = p.ParseExpr()
	if err != nil || p.tkn != token.COMMA {
		return n, err
	}
	var node *ast.Node
	for p.tkn == token.COMMA {
		p.advance()
		node = ast.NewNode(nil, nil, 2)
		node.Children[0] = n
		node.Children[1], err = p.ParseExpr()
		if err != nil {
			return node, err
		}
		n = node
	}
	return n, err
}

func (p *Parser) ParseStmt() (*ast.Node, error) {
	switch p.tkn {
	case token.IF:
		return p.ParseIfStmt()
	case token.WHILE:
		return p.ParseWhileStmt()
	case token.IDENTIFIER:
		return p.ParseAssignOrFuncCallStmt()
	}
	return &ast.Node{}, errors.New("Expected if, while, or an identifier " + "but saw " + p.tkn.String())
}

func (p *Parser) ParseIfStmt() (*ast.Node, error) {
	node := ast.NewNode(p.tkn, p.val, 2)
	p.advance()

	var err error
	node.Children[0], err = p.ParseExpr()
	if err != nil {
		return node, err
	}

	if p.tkn != token.LBRACE {
		return node, errors.New("Expected " + token.LBRACE.String() + " but saw " + p.tkn.String())
	}
	p.advance()

	node.Children[1], err = p.ParseStmt()
	if err != nil {
		return node, err
	}

	if p.tkn != token.RBRACE {
		return node, errors.New("Expected " + token.RBRACE.String() + " but saw " + p.tkn.String())
	}
	p.advance()

	return node, nil
}

func (p *Parser) ParseWhileStmt() (*ast.Node, error) {
	node := ast.NewNode(p.tkn, p.val, 2)
	p.advance()

	var err error
	node.Children[0], err = p.ParseExpr()
	if err != nil {
		return node, err
	}

	if p.tkn != token.LBRACE {
		return node, errors.New("Expected " + token.LBRACE.String() + " but saw " + p.tkn.String())
	}
	p.advance()

	node.Children[1], err = p.ParseStmt()
	if err != nil {
		return node, err
	}

	if p.tkn != token.RBRACE {
		return node, errors.New("Expected " + token.RBRACE.String() + " but saw " + p.tkn.String())
	}
	p.advance()

	return node, nil
}

func (p *Parser) ParseAssignOrFuncCallStmt() (*ast.Node, error) {
	var err error
	node := ast.NewNode(nil, nil, 2)
	node.Children[0], err = p.ParseIdentifier()
	if err != nil {
		return node, err
	}

	if p.tkn == token.ASSIGN {
		node.Token = token.ASSIGN
		node.Value = token.ASSIGN.String()
		p.advance()

		node.Children[1], err = p.ParseExpr()
		if err != nil {
			return node, err
		}
	} else if p.tkn == token.LPAREN {
		node.Children[1], err = p.ParseFuncCallArgs()
		if err != nil {
			return node, err
		}
	} else {
		return node, errors.New("Expected " + token.ASSIGN.String() + " or " + token.LPAREN.String() + " but saw " + p.tkn.String())
	}

	if p.tkn != token.SEMICOLON {
		return node, errors.New("Expected " + token.SEMICOLON.String() + " but saw " + p.tkn.String())
	}
	p.advance()

	return node, nil
}

func (p *Parser) ParseExpr() (*ast.Node, error) {
	n, err := p.ParseRelation()
	if err != nil || !p.tkn.IsLogOp() {
		return n, err
	}

	node := ast.NewNode(p.tkn, p.val, 2)
	node.Children[0] = n
	p.advance()
	node.Children[1], err = p.ParseExpr()
	return node, err
}

func (p *Parser) ParseRelation() (*ast.Node, error) {
	n, err := p.ParseSimpleExpr()
	if err != nil || !p.tkn.IsCmpOp() {
		return n, err
	}

	node := ast.NewNode(p.tkn, p.val, 2)
	node.Children[0] = n
	p.advance()
	node.Children[1], err = p.ParseRelation()
	return node, err
}

func (p *Parser) ParseSimpleExpr() (*ast.Node, error) {
	n, err := p.ParseTerm()
	if err != nil || !p.tkn.IsAddOp() {
		return n, err
	}

	node := ast.NewNode(p.tkn, p.val, 2)
	node.Children[0] = n
	p.advance()
	node.Children[1], err = p.ParseSimpleExpr()
	return node, err
}

func (p *Parser) ParseTerm() (*ast.Node, error) {
	n, err := p.ParseFactor()
	if err != nil || !p.tkn.IsMulOp() {
		return n, err
	}

	node := ast.NewNode(p.tkn, p.val, 2)
	node.Children[0] = n
	p.advance()
	node.Children[1], err = p.ParseTerm()
	return node, err
}

func (p *Parser) ParseFactor() (*ast.Node, error) {
	var n *ast.Node
	var err error
	if p.match(token.LPAREN) {
		p.advance()
		n, err = p.ParseExpr()
		if err == nil {
			if !p.match(token.RPAREN) {
				err = errors.New("xcvxcvcExpected " + token.RPAREN.String() + " but saw " + p.tkn.String())
			} else {
				p.advance()
			}
		}
		return n, err
	}
	if p.match(token.NOT) || p.match(token.ADD) || p.match(token.SUBTRACT) {
		node := ast.NewNode(p.tkn, p.val, 1)
		p.advance()
		n, err = p.ParseFactor()
		node.Children[0] = n
		return node, err
	}
	return p.ParseTerminal()
}

func (p *Parser) ParseTerminal() (*ast.Node, error) {
	n := ast.NewNode(p.tkn, p.val, 0)
	if !p.tkn.IsLiteral() {
		return n, errors.New("Expected a value or identifier but saw " + p.tkn.String())
	}
	p.advance()
	if n.Token != token.IDENTIFIER || p.tkn != token.LPAREN {
		return n, nil
	}

	var err error = nil
	node := ast.NewNode(nil, nil, 2)
	node.Children[0] = n
	node.Children[1], err = p.ParseFuncCallArgs()
	return node, err
}

func (p *Parser) ParseIdentifier() (*ast.Node, error) {
	n := ast.NewNode(p.tkn, p.val, 0)
	if p.tkn != token.IDENTIFIER {
		return n, errors.New("Expected an identifier but saw " + p.tkn.String())
	}
	p.advance()
	return n, nil
}

func (p *Parser) ParseFuncCallArgs() (*ast.Node, error) {
	var err error = nil
	n := ast.NewNode(nil, nil, 2)
	if p.tkn != token.LPAREN {
		return n, errors.New("Expected " + token.LPAREN.String() + " but saw " + p.tkn.String())
	}
	p.advance()
	if p.tkn == token.RPAREN {
		p.advance()
		return n, nil
	}
	n, err = p.ParseExprList()
	if err != nil {
		return n, err
	}
	if p.tkn != token.RPAREN {
		return n, errors.New("Expected " + token.LPAREN.String() + " but saw " + p.tkn.String())
	}
	p.advance()
	return n, nil
}
