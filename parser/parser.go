// Package parser provides a parser implementation that constructs an
// abstract syntax tree from a stream of tokens.
package parser

import (
	"errors"

	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scanner"
	"github.com/tboronczyk/kiwi/token"
)

// Parser produces ASTs from tokens.
type Parser struct {
	curToken token.Token
	prvToken token.Token
	curValue string
	prvValue string
	braces   int
	scanner  *scanner.Scanner
}

// New returns a new Parser instance initialized to read from s.
func New(s *scanner.Scanner) *Parser {
	p := &Parser{scanner: s}
	p.advance()
	return p
}

// match returns a bool indicating whether the current token matches one of the
// provided tokens.
func (p Parser) match(tokens ...token.Token) bool {
	for _, t := range tokens {
		if p.curToken == t {
			return true
		}
	}
	return false
}

// advance retrieves the next token/value pair from the scanner, keeping track
// of the previous pair. COMMENT tokens are always passed over. Multiple
// NEWLINE tokens are consolidated but never appear as the current pair. This
// allows the parser to treat arbirary newlines as whitespace but to still be
// aware of their presense when they have semantic meaning.
func (p *Parser) advance() {
	p.prvToken, p.prvValue = p.curToken, p.curValue
	for {
		p.curToken, p.curValue = p.scanner.Scan()
		if p.curToken != token.COMMENT && p.curToken != token.NEWLINE {
			return
		}
		if p.curToken == token.NEWLINE {
			p.prvToken, p.prvValue = p.curToken, p.curValue
		}
	}
}

// newline returns whether the previous token was NEWLINE.
func (p Parser) newline() bool {
	return p.prvToken == token.NEWLINE
}

// consume sets the new current token/value pair if the existing current token
// matches t, otherwise will panic.
func (p *Parser) consume(t token.Token) {
	if !p.match(t) {
		panic("unexpected " + p.curToken.String())
	}
	p.advance()
}

// stmtEnd is called by statement-parsing methods to ensure their statements
// terminate with NEWLINE when they aren't the last statement in a block.
func (p *Parser) stmtEnd() {
	if !(p.newline() || (p.curToken == token.RBRACE && p.braces > 0)) {
		panic("unexpected " + p.curToken.String())
	}
}

// Parse consumes the token stream and returns an AST of the production. err
// is nil for successful parses.
func (p *Parser) Parse() (node ast.Node, err error) {
	defer func() {
		if e := recover(); e != nil {
			node = nil
			err = errors.New(e.(string))
		}
	}()
	if p.curToken == token.EOF {
		return nil, nil
	}
	return p.stmt(), nil
}

// expr = term [expr-op expr]
func (p *Parser) expr() ast.Node {
	term := p.term()
	if !p.curToken.IsExprOp() {
		return term
	}

	op := p.curToken
	p.advance()
	expr := p.expr()

	node := &ast.BinaryOpNode{Op: op, Left: term}
	switch expr.(type) {
	case *ast.BinaryOpNode:
		prec, _ := token.Precedence(op, expr.(*ast.BinaryOpNode).Op)
		if prec {
			// adjust tree for higher precedence of expr's op
			node.Right = expr.(*ast.BinaryOpNode).Left
			expr.(*ast.BinaryOpNode).Left = node
			return expr
		}
	}
	node.Right = expr
	return node
}

// term = "(" expr ")" / term-op expr / cast
func (p *Parser) term() ast.Node {
	if p.match(token.LPAREN) {
		defer p.consume(token.RPAREN)
		p.advance()
		return p.expr()
	}
	if p.curToken.IsTermOp() {
		node := &ast.UnaryOpNode{Op: p.curToken}
		p.advance()
		node.Expr = p.term()
		return node
	}
	return p.cast()
}

// cast =: terminal [":" ident]
func (p *Parser) cast() ast.Node {
	node := p.terminal()
	if !p.match(token.COLON) {
		return node
	}
	p.advance()
	return &ast.CastNode{Cast: p.ident(), Expr: node}
}

// terminal := boolean / number / STRING / IDENT / func-call
func (p *Parser) terminal() ast.Node {
	if p.match(token.BOOL, token.NUMBER, token.STRING) {
		defer p.advance()
		return &ast.ValueNode{Value: p.curValue, Type: p.curToken}
	}

	name := p.ident()
	if !p.match(token.LPAREN) {
		return &ast.VariableNode{Name: name}
	}
	return &ast.FuncCallNode{Name: name, Args: p.parenExprList()}
}

// paren-expr-list = "(" [expr *("," expr)] ")"
func (p *Parser) parenExprList() []ast.Node {
	defer p.consume(token.RPAREN)
	p.consume(token.LPAREN)

	var list []ast.Node
	if p.match(token.RPAREN) {
		return list
	}
	for {
		list = append(list, p.expr())
		if !p.match(token.COMMA) {
			return list
		}
		p.advance()
	}
}

// stmt = if-stmt / while-stmt / func-def / return-stmt / assign-stmt /
//        func-call
func (p *Parser) stmt() ast.Node {
	switch p.curToken {
	case token.IF:
		return p.ifStmt()
	case token.WHILE:
		return p.whileStmt()
	case token.FUNC:
		return p.funcDef()
	case token.RETURN:
		return p.returnStmt()
	case token.IDENTIFIER:
		return p.assignStmtOrFuncCall()
	}
	panic("unexpected " + p.curToken.String())
}

// if-stmt = "if" expr brace-stmt-list [else-clause]
func (p *Parser) ifStmt() *ast.IfNode {
	p.consume(token.IF)
	node := &ast.IfNode{Condition: p.expr(), Body: p.braceStmtList()}
	if p.match(token.ELSE) {
		node.Else = p.elseClause()
	}
	return node
}

// brace-stmt-list = "{" *stmt "}"
func (p *Parser) braceStmtList() []ast.Node {
	defer func() {
		p.consume(token.RBRACE)
		p.braces--
	}()

	p.consume(token.LBRACE)
	p.braces++

	var list []ast.Node
	for {
		if !(p.curToken.IsStmtKeyword() || p.match(token.IDENTIFIER)) {
			break
		}
		list = append(list, p.stmt())
	}
	return list
}

// else-clause = "else" (brace-stmt-list / expr brace-stmt-list else-clause)
// An else with an expression becomes an if-stmt within default else clause.
func (p *Parser) elseClause() *ast.IfNode {
	p.consume(token.ELSE)
	node := &ast.IfNode{}
	isFinal := false
	if p.match(token.LBRACE) {
		// a final clause without an expression defaults to an
		// expression that evaluates true to make evaluation of the
		// AST easier.
		isFinal = true
		node.Condition = &ast.ValueNode{Value: "TRUE", Type: token.BOOL}
	} else {
		node.Condition = p.expr()
	}
	node.Body = p.braceStmtList()
	if !isFinal {
		node.Else = p.elseClause()
	}
	return node
}

// while-stmt = "while" expr brace-stmt-list
func (p *Parser) whileStmt() *ast.WhileNode {
	p.consume(token.WHILE)
	return &ast.WhileNode{Condition: p.expr(), Body: p.braceStmtList()}
}

// func-def = "func" [ident-list] brace-stmt-list
func (p *Parser) funcDef() *ast.FuncDefNode {
	p.consume(token.FUNC)
	node := &ast.FuncDefNode{Name: p.ident()}
	if !p.match(token.LBRACE) {
		node.Args = p.identList()
	}
	node.Body = p.braceStmtList()
	return node
}

// ident-list = ident *("," ident)
func (p *Parser) identList() []string {
	var list []string
	for {
		list = append(list, p.ident())
		if !p.match(token.COMMA) {
			return list
		}
		p.advance()
	}
}

// return-stmt = "return" [expr] LF
func (p *Parser) returnStmt() *ast.ReturnNode {
	defer p.stmtEnd()
	p.consume(token.RETURN)
	node := &ast.ReturnNode{}
	if !p.newline() {
		node.Expr = p.expr()
	}
	return node
}

// assign-stmt = ident ":=" expr LF
// func-call   = ident paren-expr-list
func (p *Parser) assignStmtOrFuncCall() ast.Node {
	defer p.stmtEnd()
	name := p.ident()
	if p.match(token.ASSIGN) {
		p.advance()
		return &ast.AssignNode{Name: name, Expr: p.expr()}
	}
	if p.match(token.LPAREN) {
		return &ast.FuncCallNode{Name: name, Args: p.parenExprList()}
	}
	panic("unexpected " + p.curToken.String())
}

// identifier returns the lexeme value of the current identifier.
func (p *Parser) ident() string {
	defer p.consume(token.IDENTIFIER)
	return p.curValue
}
