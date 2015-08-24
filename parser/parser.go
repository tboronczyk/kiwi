// Package parser provides the implementation of a parser that constructs an
// abstract syntax tree from a stream of tokens to represent a Kiwi program in
// memory.
package parser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scanner"
	"github.com/tboronczyk/kiwi/scope"
	"github.com/tboronczyk/kiwi/token"
	"github.com/tboronczyk/kiwi/types"
)

// Parser produces ASTs from tokens.
type Parser struct {
	curToken token.Token // the current token
	prvToken token.Token // the previous token
	curValue string      // the current lexeme value
	prvValue string      // the previous lexeme value
	braces   int         // pseudo-stack counter tracking brace nesting
	scanner  *scanner.Scanner
	scope    *scope.Scope
}

// New returns a new Parser instance that is initialized to read from s.
func New(s *scanner.Scanner) *Parser {
	p := &Parser{scanner: s, scope: scope.New()}
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
func (p *Parser) Parse() (prog *ast.ProgramNode, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(e.(string))
		}
	}()

	prog = &ast.ProgramNode{Scope: p.scope}
	for {
		if p.curToken == token.EOF {
			return prog, nil
		}
		prog.Stmts = append(prog.Stmts, p.stmt())
	}
}

// expr = cmp-expr [log-op expr]
func (p *Parser) expr() ast.Node {
	node := p.cmpExpr()
	switch p.curToken {
	case token.AND:
		p.advance()
		return &ast.AndNode{Left: node, Right: p.expr()}
	case token.OR:
		p.advance()
		return &ast.OrNode{Left: node, Right: p.expr()}
	}
	return node
}

// cmp-expr = add-expr [cmp-op cmp-expr]
func (p *Parser) cmpExpr() ast.Node {
	node := p.addExpr()
	switch p.curToken {
	case token.EQUAL:
		p.advance()
		return &ast.EqualNode{Left: node, Right: p.cmpExpr()}
	case token.NOT_EQUAL:
		p.advance()
		return &ast.NotEqualNode{Left: node, Right: p.cmpExpr()}
	case token.GREATER:
		p.advance()
		return &ast.GreaterNode{Left: node, Right: p.cmpExpr()}
	case token.GREATER_EQ:
		p.advance()
		return &ast.GreaterEqualNode{Left: node, Right: p.cmpExpr()}
	case token.LESS:
		p.advance()
		return &ast.LessNode{Left: node, Right: p.cmpExpr()}
	case token.LESS_EQ:
		p.advance()
		return &ast.LessEqualNode{Left: node, Right: p.cmpExpr()}
	}
	return node
}

// add-expr = mul-expr [add-op add-expr]
func (p *Parser) addExpr() ast.Node {
	node := p.mulExpr()
	switch p.curToken {
	case token.ADD:
		p.advance()
		return &ast.AddNode{Left: node, Right: p.addExpr()}
	case token.SUBTRACT:
		p.advance()
		return &ast.SubtractNode{Left: node, Right: p.addExpr()}
	}
	return node
}

// mul-expr = cast-expr [mul-op mul-expr]
func (p *Parser) mulExpr() ast.Node {
	node := p.castExpr()
	switch p.curToken {
	case token.MULTIPLY:
		p.advance()
		return &ast.MultiplyNode{Left: node, Right: p.mulExpr()}
	case token.DIVIDE:
		p.advance()
		return &ast.DivideNode{Left: node, Right: p.mulExpr()}
	case token.MODULO:
		p.advance()
		return &ast.ModuloNode{Left: node, Right: p.mulExpr()}
	}
	return node
}

// cast-expr = term [":" ident]
func (p *Parser) castExpr() ast.Node {
	node := p.term()
	if p.curToken == token.COLON {
		p.advance()
		return &ast.CastNode{Cast: p.ident(), Term: node}
	}
	return node
}

// term = "(" expr ")" / ("+" / "-" / "~") term / boolean / number / string /
//        func-call / ident
func (p *Parser) term() ast.Node {
	switch p.curToken {
	case token.LPAREN:
		p.advance()
		node := p.expr()
		p.consume(token.RPAREN)
		return node
	case token.ADD:
		p.advance()
		return &ast.PositiveNode{Term: p.term()}
	case token.SUBTRACT:
		p.advance()
		return &ast.NegativeNode{Term: p.term()}
	case token.NOT:
		p.advance()
		return &ast.NotNode{Term: p.term()}
	case token.BOOL:
		node := &ast.BoolNode{Value: strings.ToLower(p.curValue) == "true"}
		p.advance()
		return node
	case token.NUMBER:
		val, _ := strconv.ParseFloat(p.curValue, 64)
		node := &ast.NumberNode{Value: val}
		p.advance()
		return node
	case token.STRING:
		node := &ast.StringNode{Value: p.curValue}
		p.advance()
		return node
	case token.IDENTIFIER:
		name := p.ident()
		if p.match(token.LPAREN) {
			return &ast.FuncCallNode{Name: name, Args: p.parenExprList()}
		}
		return &ast.VariableNode{Name: name}
	}

	panic("whoops?")
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
	node := &ast.IfNode{Cond: p.expr(), Body: p.braceStmtList()}
	if p.match(token.ELSE) {
		p.advance()
		if p.match(token.LBRACE) {
			node.Else = p.braceStmtList()
		} else {
			node.Else = append(node.Else, p.elseClause())
		}
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
	node := &ast.IfNode{Cond: p.expr(), Body: p.braceStmtList()}
	if p.match(token.ELSE) {
		p.advance()
		if p.match(token.LBRACE) {
			node.Else = p.braceStmtList()
		} else {
			node.Else = append(node.Else, p.elseClause())
		}
	}
	return node
}

// while-stmt = "while" expr brace-stmt-list
func (p *Parser) whileStmt() *ast.WhileNode {
	p.consume(token.WHILE)
	return &ast.WhileNode{Cond: p.expr(), Body: p.braceStmtList()}
}

// func-def = "func" ident *ident brace-stmt-list
func (p *Parser) funcDef() *ast.FuncDefNode {
	p.consume(token.FUNC)

	node := &ast.FuncDefNode{
		Name: p.ident(),
	}
	p.scope.SetFunc(node.Name, scope.Entry{Value: node, DataType: types.FUNC})
	node.Scope = scope.NewWithParent(p.scope)
	p.scope = node.Scope

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

	p.scope = node.Scope.Parent
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
