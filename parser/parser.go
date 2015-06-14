package parser

import (
	"errors"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/scanner"
	"github.com/tboronczyk/kiwi/token"
)

type (
	Parser interface {
		InitScanner(scanner.Scanner)
		Parse() (ast.StmtNode, error)
	}

	parser struct {
		token   token.Token
		value   string
		scanner scanner.Scanner
	}
)

func New() *parser {
	return &parser{}
}

func (p parser) match(tkn token.Token) bool {
	return p.token == tkn
}

func (p *parser) advance() {
	for {
		p.token, p.value = p.scanner.Scan()
		if p.token != token.COMMENT {
			break
		}
	}
}

func (p *parser) consume(t token.Token) {
	if p.token != t {
		panic(p.expected(t))
	}
	p.advance()
}

func (p parser) expected(value interface{}) string {
	var str string
	switch value.(type) {
	case token.Token:
		str = value.(token.Token).String()
	default:
		str = value.(string)
	}
	return "Expected " + str + " but saw " + p.token.String()
}

func (p *parser) InitScanner(scnr scanner.Scanner) {
	p.scanner = scnr
	p.advance()
}

func (p *parser) Parse() (node ast.StmtNode, err error) {
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

func (p *parser) expr() ast.ExprNode {
	n := p.relation()
	if !p.token.IsLogOp() {
		return n
	}

	node := ast.BinaryExpr{Op: p.token, Left: n}
	p.advance()
	node.Right = p.expr()

	return node
}

func (p *parser) relation() ast.ExprNode {
	n := p.simpleExpr()
	if !p.token.IsCmpOp() {
		return n
	}

	node := ast.BinaryExpr{Op: p.token, Left: n}
	p.advance()
	node.Right = p.relation()

	return node
}

func (p *parser) simpleExpr() ast.ExprNode {
	n := p.term()
	if !p.token.IsAddOp() {
		return n
	}

	node := ast.BinaryExpr{Op: p.token, Left: n}
	p.advance()
	node.Right = p.simpleExpr()

	return node
}

func (p *parser) term() ast.ExprNode {
	n := p.factor()
	if !p.token.IsMulOp() {
		return n
	}

	node := ast.BinaryExpr{Op: p.token, Left: n}
	p.advance()
	node.Right = p.term()

	return node
}

func (p *parser) factor() ast.ExprNode {
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

func (p *parser) terminal() ast.ExprNode {
	if p.token == token.BOOL || p.token == token.NUMBER ||
		p.token == token.STRING {
		defer p.advance()
		return ast.ValueExpr{Value: p.value, Type: p.token}
	}

	name := p.identifier()
	if p.token != token.LPAREN {
		return ast.VariableExpr{Name: name}
	}
	return ast.FuncCall{Name: name, Args: p.parenExprList()}
}

func (p *parser) parenExprList() []ast.ExprNode {
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

func (p *parser) identList() []string {
	var list []string
	for {
		list = append(list, p.identifier())
		if p.token != token.COMMA {
			return list
		}
		p.advance()
	}
}

func (p *parser) stmt() ast.StmtNode {
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

func (p *parser) funcDef() ast.FuncDef {
	p.consume(token.FUNC)
	node := ast.FuncDef{Name: p.identifier()}
	if p.token != token.LBRACE {
		node.Args = p.identList()
	}
	node.Body = p.braceStmtList()
	return node
}

func (p *parser) ifStmt() ast.IfStmt {
	p.consume(token.IF)
	return ast.IfStmt{Condition: p.expr(), Body: p.braceStmtList()}
}

func (p *parser) returnStmt() ast.ReturnStmt {
	defer p.consume(token.DOT)
	p.consume(token.RETURN)
	node := ast.ReturnStmt{}
	if p.token != token.DOT {
		node.Expr = p.expr()
	}
	return node
}

func (p *parser) braceStmtList() []ast.StmtNode {
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

func (p *parser) whileStmt() ast.WhileStmt {
	p.consume(token.WHILE)
	return ast.WhileStmt{Condition: p.expr(), Body: p.braceStmtList()}
}

func (p *parser) assignStmtOrFuncCall() ast.StmtNode {
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

func (p *parser) identifier() string {
	defer p.consume(token.IDENTIFIER)
	return p.value
}
