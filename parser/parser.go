// Package parser provides the implementation of a parser that constructs an
// abstract syntax tree from a stream of tokens to represent a Kiwi program in
// memory.
package parser

import (
	"errors"

	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scanner"
	"github.com/tboronczyk/kiwi/token"
)

// Parser produces ASTs from tokens.
type Parser struct {
	curToken token.Token // the current token
	prvToken token.Token // the previous token
	curValue string      // the current lexeme value
	prvValue string      // the previous lexeme value
	braces   int         // pseudo-stack counter tracking brace nesting
	scanner  *scanner.Scanner
}

// New returns a new Parser instance that is initialized to read from s.
func New(s *scanner.Scanner) *Parser {
	p := &Parser{scanner: s}
	p.advance()
	return p
}

// match returns a bool to indicate whether the current token matches one of
// specified tokens.
func (p Parser) match(tokens ...token.Token) bool {
	for _, t := range tokens {
		if p.curToken == t {
			return true
		}
	}
	return false
}

// advance retrieves the next token/value pair from the scanner, keeping track
// of the previous pair. NOTE: COMMENT tokens are always passed over. Multiple
// NEWLINE tokens are consolidated. A NEWLINE will never appear as the current
// pair (only previous). This allows the parser to treat arbirary newlines as
// whitespace and still be aware of their presense when they have semantic
// meaning.
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

// newline returns a bool to indicate whether the previous token was NEWLINE.
func (p Parser) newline() bool {
	return p.prvToken == token.NEWLINE
}

// consume advances the parser to the next token/value pair if the current token
// matches t, otherwise it will panic.
func (p *Parser) consume(t token.Token) {
	if !p.match(t) {
		panic("unexpected " + p.curToken.String())
	}
	p.advance()
}

// stmtEnd may be called by a statement-parsing method to ensure a statement
// is terminated by a NEWLINE when it isn't the last statement in a block
// (the terminating NEWLINE is optional for final statements in blocks).
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

// expr = term [bin-op expr]
func (p *Parser) expr() ast.Node {
	term := p.term()
	if !p.curToken.IsBinOp() {
		return term
	}

	op := p.curToken
	p.advance()
	expr := p.expr()

	node := &ast.BinOpNode{Op: op, Left: term}
	switch expr.(type) {
	case *ast.BinOpNode:
		if token.Precedence(op) > token.Precedence(expr.(*ast.BinOpNode).Op) {
			// adjust tree for higher precedence of expr's op
			node.Right = expr.(*ast.BinOpNode).Left
			expr.(*ast.BinOpNode).Left = node
			return expr
		}
	}
	node.Right = expr
	return node
}

// term = "(" expr ")" / unary-op term / boolean / number / string / ident /
//        func-call / term [":" ident]
func (p *Parser) term() ast.Node {
	var node ast.Node
	if p.match(token.LPAREN) {
		p.advance()
		node = p.expr()
		p.consume(token.RPAREN)
	} else if p.curToken.IsUnaryOp() {
		op := p.curToken
		p.advance()
		node = &ast.UnaryOpNode{Op: op, Term: p.term()}
	} else if p.match(token.BOOL, token.NUMBER, token.STRING) {
		node = &ast.ValueNode{Value: p.curValue, Type: p.curToken}
		p.advance()
	} else {
		name := p.ident()
		if p.match(token.LPAREN) {
			node = &ast.FuncCallNode{Name: name, Args: p.parenExprList()}
		} else {
			node = &ast.VariableNode{Name: name}
		}
	}
	if p.match(token.COLON) {
		p.advance()
		node = &ast.CastNode{Cast: p.ident(), Term: node}
	}
	return node
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
// Note: an else with an expression becomes an if-stmt within a default else
// clause.
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

// func-def = "func" ident *ident brace-stmt-list
func (p *Parser) funcDef() *ast.FuncDefNode {
	p.consume(token.FUNC)
	node := &ast.FuncDefNode{Name: p.ident()}
	if !p.match(token.LBRACE) {
		var list []string
		for {
			list = append(list, p.ident())
			if !p.match(token.IDENTIFIER) {
				break
			}
		}
		node.Args = list
	}
	node.Body = p.braceStmtList()
	return node
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
