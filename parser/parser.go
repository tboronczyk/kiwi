// Package parser provides the parser implementation that constructs an
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
	token   token.Token
	value   string
	scanner *scanner.Scanner
}

// New returns a new Parser initialized to read from s.
func New(s *scanner.Scanner) *Parser {
	p := &Parser{scanner: s}
	p.advance()
	return p
}

// match returns a bool indicating the current token in memory is the same
// as t.
func (p Parser) match(t token.Token) bool {
	return p.token == t
}

// advance retrieves the next token/value pair from the token stream and sets
// them as the current token and value. COMMENT tokens are skipped.
func (p *Parser) advance() {
	for {
		p.token, p.value = p.scanner.Scan()
		if p.token != token.COMMENT {
			break
		}
	}
}

// consume sets the new current token/value pair if the existing current token
// matches t, otherwise will panic.
func (p *Parser) consume(t token.Token) {
	if p.token != t {
		panic(p.expected(t))
	}
	p.advance()
}

// expected formats v as an error string. If v is a Token then its value is
// displayed as its string value, otherwise v is cast directly to string.
func (p Parser) expected(v interface{}) string {
	var str string
	switch v.(type) {
	case token.Token:
		str = v.(token.Token).String()
	default:
		str = v.(string)
	}
	return "Expected " + str + " but saw " + p.token.String()
}

// Parse consumes the token stream and returns an AST of the production. err
// is nil for successful parses.
func (p *Parser) Parse() (node ast.Node, err error) {
	if p.token == token.EOF {
		return nil, nil
	}
	defer func() {
		if e := recover(); e != nil {
			node = nil
			err = errors.New(e.(string))
		}
	}()
	return p.stmt(), nil
}

// expr := relation (log-op expr)? .
func (p *Parser) expr() ast.Node {
	n := p.relation()
	if !p.token.IsLogOp() {
		return n
	}

	node := &ast.BinaryOpNode{Op: p.token, Left: n}
	p.advance()
	node.Right = p.expr()

	return node
}

// relation := simple-expr (cmp-op relation)? .
func (p *Parser) relation() ast.Node {
	n := p.simpleExpr()
	if !p.token.IsCmpOp() {
		return n
	}

	node := &ast.BinaryOpNode{Op: p.token, Left: n}
	p.advance()
	node.Right = p.relation()

	return node
}

// simple-expr := term (add-op simple-expr)? .
func (p *Parser) simpleExpr() ast.Node {
	n := p.term()
	if !p.token.IsAddOp() {
		return n
	}

	node := &ast.BinaryOpNode{Op: p.token, Left: n}
	p.advance()
	node.Right = p.simpleExpr()

	return node
}

// term =: factor (mul-op term)? .
func (p *Parser) term() ast.Node {
	n := p.factor()
	if !p.token.IsMulOp() {
		return n
	}

	node := &ast.BinaryOpNode{Op: p.token, Left: n}
	p.advance()
	node.Right = p.term()

	return node
}

// factor =: '(' expr ')' | expr-op expr | cast .
func (p *Parser) factor() ast.Node {
	if p.match(token.LPAREN) {
		defer p.consume(token.RPAREN)
		p.advance()
		return p.expr()
	}
	if p.token.IsExprOp() {
		node := &ast.UnaryOpNode{Op: p.token}
		p.advance()
		node.Right = p.factor()
		return node
	}
	return p.cast()
}

// cast =: terminal ('!' IDENT)? .
func (p *Parser) cast() ast.Node {
	node := p.terminal()
	if p.token != token.CAST {
		return node
	}
	p.advance()
	return &ast.CastNode{Cast: p.identifier(), Expr: node}
}

// terminal := boolean | number | STRING | IDENT | func-call .
func (p *Parser) terminal() ast.Node {
	if p.token == token.BOOL || p.token == token.NUMBER ||
		p.token == token.STRING {
		defer p.advance()
		return &ast.ValueNode{Value: p.value, Type: p.token}
	}

	name := p.identifier()
	if p.token != token.LPAREN {
		return &ast.VariableNode{Name: name}
	}
	return &ast.FuncCallNode{Name: name, Args: p.parenExprList()}
}

// paren-expr-list := '(' ')' | '(' expr (',' expr)* ')' .
func (p *Parser) parenExprList() []ast.Node {
	defer p.consume(token.RPAREN)
	p.consume(token.LPAREN)

	var list []ast.Node
	if p.token == token.RPAREN {
		return list
	}
	for {
		list = append(list, p.expr())
		if p.token != token.COMMA {
			return list
		}
		p.advance()
	}
}

// stmt := if-stmt | while-stmt | func-def | return-stmt | assign-stmt |
//         func-call .
func (p *Parser) stmt() ast.Node {
	switch p.token {
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
	panic(p.expected("statement"))
}

// if-stmt := 'if' expr brace-stmt-list (else-clause)? .
func (p *Parser) ifStmt() *ast.IfNode {
	p.consume(token.IF)
	node := &ast.IfNode{Condition: p.expr(), Body: p.braceStmtList()}
	if p.token == token.ELSE {
		node.Else = p.elseClause()
	}
	return node
}

// brace-stmt-list := '{' (stmt)* '}'
func (p *Parser) braceStmtList() []ast.Node {
	defer p.consume(token.RBRACE)
	p.consume(token.LBRACE)

	var list []ast.Node
	for {
		if !(p.token.IsStmtKeyword() || p.token == token.IDENTIFIER) {
			return list
		}
		list = append(list, p.stmt())
	}
}

// else-clause := 'else' (brace-stmt-list | expr brace-stmt-list else-clause) .
// An else with an expression becomes an if-stmt within default else clause.
func (p *Parser) elseClause() *ast.IfNode {
	p.consume(token.ELSE)
	node := &ast.IfNode{}
	isFinal := false
	if p.token == token.LBRACE {
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

// while-stmt := 'while' expr brace-stmt-list .
func (p *Parser) whileStmt() *ast.WhileNode {
	p.consume(token.WHILE)
	return &ast.WhileNode{Condition: p.expr(), Body: p.braceStmtList()}
}

// func-def := 'func' (ident-list)? brace-stmt-list .
func (p *Parser) funcDef() *ast.FuncDefNode {
	p.consume(token.FUNC)
	node := &ast.FuncDefNode{Name: p.identifier()}
	if p.token != token.LBRACE {
		node.Args = p.identList()
	}
	node.Body = p.braceStmtList()
	return node
}

// ident-list := IDENT (',' IDENT)? .
func (p *Parser) identList() []string {
	var list []string
	for {
		list = append(list, p.identifier())
		if p.token != token.COMMA {
			return list
		}
		p.advance()
	}
}

// return-stmt := 'return' (expr)? '.' .
func (p *Parser) returnStmt() *ast.ReturnNode {
	defer p.consume(token.DOT)
	p.consume(token.RETURN)
	node := &ast.ReturnNode{}
	if p.token != token.DOT {
		node.Expr = p.expr()
	}
	return node
}

// assign-stmt := IDENT ':=' expr '.' .
//   func-call := IDENT paren-expr-list .
func (p *Parser) assignStmtOrFuncCall() ast.Node {
	defer p.consume(token.DOT)

	name := p.identifier()
	if p.token == token.ASSIGN {
		p.advance()
		return &ast.AssignNode{Name: name, Expr: p.expr()}
	}
	if p.token == token.LPAREN {
		return &ast.FuncCallNode{Name: name, Args: p.parenExprList()}
	}
	panic(p.expected(
		token.ASSIGN.String() + " or " + token.LPAREN.String()))
}

// identifier returns the identifier's value.
func (p *Parser) identifier() string {
	defer p.consume(token.IDENTIFIER)
	return p.value
}
