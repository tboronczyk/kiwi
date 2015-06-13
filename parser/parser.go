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

func (p *Parser) consume(t token.Token) {
	if p.token != t {
		panic(p.expected(t))
	}
	p.advance()
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

func (p *Parser) Parse() (node ast.StmtNode, err error) {
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

func (p *Parser) expr() ast.ExprNode {
	n := p.relation()
	if !p.token.IsLogOp() {
		return n
	}

	node := ast.BinaryExpr{Op: p.token, Left: n}
	p.advance()
	node.Right = p.expr()

	return node
}

func (p *Parser) relation() ast.ExprNode {
	n := p.simpleExpr()
	if !p.token.IsCmpOp() {
		return n
	}

	node := ast.BinaryExpr{Op: p.token, Left: n}
	p.advance()
	node.Right = p.relation()

	return node
}

func (p *Parser) simpleExpr() ast.ExprNode {
	n := p.term()
	if !p.token.IsAddOp() {
		return n
	}

	node := ast.BinaryExpr{Op: p.token, Left: n}
	p.advance()
	node.Right = p.simpleExpr()

	return node
}

func (p *Parser) term() ast.ExprNode {
	n := p.factor()
	if !p.token.IsMulOp() {
		return n
	}

	node := ast.BinaryExpr{Op: p.token, Left: n}
	p.advance()
	node.Right = p.term()

	return node
}

func (p *Parser) factor() ast.ExprNode {
	if p.match(token.LPAREN) {
		defer p.consume(token.RPAREN)
		p.advance()
		return p.expr()
	}
	if p.token.IsExprOp() {
		node := ast.UnaryExpr{Op: p.token}
		p.advance()
		node.Right = p.factor()
		return node
	}
	return p.terminal()
}

func (p *Parser) terminal() ast.ExprNode {
	if p.token == token.TRUE || p.token == token.FALSE ||
		p.token == token.NUMBER || p.token == token.STRING {
		defer p.advance()
		return ast.ValueExpr{Value: p.value, Type: p.token}
	}

	name := p.identifier()
	if p.token != token.LPAREN {
		return ast.VariableExpr{Name: name}
	}
	return ast.FuncCall{Name: name, Args: p.parenExprList()}
}

func (p *Parser) parenExprList() []ast.ExprNode {
	defer p.consume(token.RPAREN)
	p.consume(token.LPAREN)

	var list []ast.ExprNode
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

func (p *Parser) stmt() ast.StmtNode {
	switch p.token {
	case token.FUNC:
		return p.funcDef()
	case token.IF:
		return p.ifStmt()
	case token.RETURN:
		return p.returnStmt()
	case token.WHILE:
		return p.whileStmt()
	case token.IDENTIFIER:
		return p.assignStmtOrFuncCall()
	}
	panic(p.expected("statement"))
}

func (p *Parser) funcDef() ast.FuncDef {
	p.consume(token.FUNC)
	node := ast.FuncDef{Name: p.identifier()}
	if p.token != token.LBRACE {
		node.Args = p.identList()
	}
	node.Body = p.braceStmtList()
	return node
}

func (p *Parser) ifStmt() ast.IfStmt {
	p.consume(token.IF)
	return ast.IfStmt{Condition: p.expr(), Body: p.braceStmtList()}
}

func (p *Parser) returnStmt() ast.ReturnStmt {
	defer p.consume(token.DOT)
	p.consume(token.RETURN)
	node := ast.ReturnStmt{}
	if p.token != token.DOT {
		node.Expr = p.expr()
	}
	return node
}

func (p *Parser) braceStmtList() []ast.StmtNode {
	defer p.consume(token.RBRACE)
	p.consume(token.LBRACE)

	var list []ast.StmtNode
	for {
		if !(p.token.IsStmtKeyword() || p.token == token.IDENTIFIER) {
			return list
		}
		list = append(list, p.stmt())
	}
}

func (p *Parser) whileStmt() ast.WhileStmt {
	p.consume(token.WHILE)
	return ast.WhileStmt{Condition: p.expr(), Body: p.braceStmtList()}
}

func (p *Parser) assignStmtOrFuncCall() ast.StmtNode {
	defer p.consume(token.DOT)

	name := p.identifier()
	if p.token == token.ASSIGN {
		p.advance()
		return ast.AssignStmt{Name: name, Expr: p.expr()}
	}
	if p.token == token.LPAREN {
		return ast.FuncCall{Name: name, Args: p.parenExprList()}
	}
	panic(p.expected(
		token.ASSIGN.String() + " or " + token.LPAREN.String()))
}

func (p *Parser) identifier() string {
	defer p.consume(token.IDENTIFIER)
	return p.value
}
