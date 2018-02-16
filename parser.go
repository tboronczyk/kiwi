package main

import (
	"errors"
	"strconv"
	"strings"
)

type Parser struct {
	curToken Token
	curValue string
	scanner  *Scanner
	scope    *Scope
}

func NewParser(s *Scanner) *Parser {
	p := &Parser{
		scanner: s,
		scope:   NewScope(),
	}
	p.advance()
	return p
}

// advance retrieves the next token/value pair from the scanner. COMMENT tokens
// are skipped as whitespace.
func (p *Parser) advance() {
	for {
		p.curToken, p.curValue = p.scanner.Scan()
		if p.curToken != TkComment {
			return
		}
	}
}

// match returns bool indicating whether the current token matches one of the
// specified tokens.
func (p Parser) match(tokens ...Token) bool {
	for _, t := range tokens {
		if p.curToken == t {
			return true
		}
	}
	return false
}

// consume advances to the next token/value pair when the current token matches
// one in t, otherwise it panics.
func (p *Parser) consume(t Token) {
	if !p.match(t) {
		panic("unexpected lexeme " + p.curToken.String())
	}
	p.advance()
}

// Parse consumes the token stream and returns the parsed program as an AST
// (ProgramNode). err is nil for a successful parse.
func (p *Parser) Parse() (prog *AstProgramNode, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(e.(string))
		}
	}()

	prog = &AstProgramNode{Scope: p.scope}
	for {
		if p.curToken == TkEOF {
			return prog, nil
		}
		prog.Stmts = append(prog.Stmts, p.stmt())
	}
}

// expr = cmp-expr [log-op expr]
func (p *Parser) expr() AstNode {
	node := p.cmpExpr()
	switch p.curToken {
	case TkAnd:
		p.advance()
		return &AstAndNode{Left: node, Right: p.expr()}
	case TkOr:
		p.advance()
		return &AstOrNode{Left: node, Right: p.expr()}
	}
	return node
}

// cmp-expr = add-expr [cmp-op cmp-expr]
func (p *Parser) cmpExpr() AstNode {
	node := p.addExpr()

	switch p.curToken {
	case TkEqual:
		p.advance()
		return &AstEqualNode{Left: node, Right: p.cmpExpr()}
	case TkNotEqual:
		p.advance()
		return &AstNotEqualNode{Left: node, Right: p.cmpExpr()}
	case TkGreater:
		p.advance()
		return &AstGreaterNode{Left: node, Right: p.cmpExpr()}
	case TkGreaterEq:
		p.advance()
		return &AstGreaterEqualNode{Left: node, Right: p.cmpExpr()}
	case TkLess:
		p.advance()
		return &AstLessNode{Left: node, Right: p.cmpExpr()}
	case TkLessEq:
		p.advance()
		return &AstLessEqualNode{Left: node, Right: p.cmpExpr()}
	}
	return node
}

// add-expr = mul-expr [add-op add-expr]
func (p *Parser) addExpr() AstNode {
	node := p.mulExpr()
	switch p.curToken {
	case TkAdd:
		p.advance()
		return &AstAddNode{Left: node, Right: p.addExpr()}
	case TkSubtract:
		p.advance()
		return &AstSubtractNode{Left: node, Right: p.addExpr()}
	}
	return node
}

// mul-expr = cast-expr [mul-op mul-expr]
func (p *Parser) mulExpr() AstNode {
	node := p.castExpr()
	switch p.curToken {
	case TkMultiply:
		p.advance()
		return &AstMultiplyNode{Left: node, Right: p.mulExpr()}
	case TkDivide:
		p.advance()
		return &AstDivideNode{Left: node, Right: p.mulExpr()}
	case TkModulo:
		p.advance()
		return &AstModuloNode{Left: node, Right: p.mulExpr()}
	}
	return node
}

// cast-expr = term [":" ident]
func (p *Parser) castExpr() AstNode {
	node := p.term()
	if p.curToken == TkColon {
		p.advance()
		return &AstCastNode{Cast: p.ident(), Term: node}
	}
	return node
}

// term = "(" expr ")" / ("+" / "-" / "~") term / boolean / number / string /
//        func-call / ident
func (p *Parser) term() AstNode {
	switch p.curToken {
	case TkLParen:
		p.advance()
		node := p.expr()
		p.consume(TkRParen)
		return node
	case TkAdd:
		p.advance()
		return &AstPositiveNode{p.term()}
	case TkSubtract:
		p.advance()
		return &AstNegativeNode{p.term()}
	case TkIf:
		p.advance()
		return &AstNotNode{p.term()}
	case TkBool:
		node := &AstBoolNode{strings.ToLower(p.curValue) == "true"}
		p.advance()
		return node
	case TkNumber:
		val, _ := strconv.ParseFloat(p.curValue, 64)
		node := &AstNumberNode{val}
		p.advance()
		return node
	case TkString:
		node := &AstStringNode{p.curValue}
		p.advance()
		return node
	case TkIdentifier:
		name := p.ident()
		if p.match(TkLParen) {
			return &AstFuncCallNode{Name: name, Args: p.parenExprList()}
		}
		return &AstVariableNode{Name: name}
	}

	panic("whoops?")
}

// paren-expr-list = "(" [expr *("," expr)] ")"
func (p *Parser) parenExprList() []AstNode {
	defer p.consume(TkRParen)
	p.consume(TkLParen)

	var list []AstNode
	if p.match(TkRParen) {
		return list
	}
	for {
		list = append(list, p.expr())
		if !p.match(TkComma) {
			return list
		}
		p.advance()
	}
}

// stmt = if-stmt / while-stmt / func-def / return-stmt / assign-stmt /
//        func-call
func (p *Parser) stmt() (node AstNode) {
	switch p.curToken {
	case TkIf:
		return p.ifStmt()
	case TkWhile:
		return p.whileStmt()
	case TkFunc:
		return p.funcDef()
	case TkReturn:
		return p.returnStmt()
	case TkIdentifier:
		return p.assignStmtOrFuncCall()
	}
	panic("unexpected lexeme " + p.curToken.String())
}

// if-stmt = "if" expr brace-stmt-list [else-clause]
func (p *Parser) ifStmt() *AstIfNode {
	p.consume(TkIf)
	node := &AstIfNode{Cond: p.expr(), Body: p.braceStmtList()}
	if p.match(TkElse) {
		p.advance()
		if p.match(TkLBrace) {
			node.Else = p.braceStmtList()
		} else {
			node.Else = append(node.Else, p.elseClause())
		}
	}
	return node
}

// brace-stmt-list = "{" *stmt "}"
func (p *Parser) braceStmtList() (list []AstNode) {
	p.consume(TkLBrace)
	for {
		if !(p.curToken.IsStmtKeyword() || p.match(TkIdentifier)) {
			break
		}
		list = append(list, p.stmt())
	}
	p.consume(TkRBrace)
	return list
}

// else-clause = "else" (brace-stmt-list / expr brace-stmt-list else-clause)
// Note: an else with an expression becomes an if-stmt within a default else
// clause.
func (p *Parser) elseClause() *AstIfNode {
	node := &AstIfNode{Cond: p.expr(), Body: p.braceStmtList()}
	if p.match(TkElse) {
		p.advance()
		if p.match(TkLBrace) {
			node.Else = p.braceStmtList()
		} else {
			node.Else = append(node.Else, p.elseClause())
		}
	}
	return node
}

// while-stmt = "while" expr brace-stmt-list
func (p *Parser) whileStmt() *AstWhileNode {
	p.consume(TkWhile)
	return &AstWhileNode{Cond: p.expr(), Body: p.braceStmtList()}
}

// func-def = "func" ident *ident brace-stmt-list
func (p *Parser) funcDef() *AstFuncDefNode {
	p.consume(TkFunc)

	node := &AstFuncDefNode{
		Name: p.ident(),
	}
	p.scope.SetFunc(node.Name, ScopeEntry{TypFunc, node})
	node.Scope = NewScopeWithParent(p.scope)
	p.scope = node.Scope

	if !p.match(TkLBrace) {
		var list []string
		for p.match(TkIdentifier) {
			list = append(list, p.ident())
		}
		node.Args = list
	}
	node.Body = p.braceStmtList()

	p.scope = node.Scope.parent
	return node
}

// return-stmt = "return" [expr]
func (p *Parser) returnStmt() *AstReturnNode {
	p.consume(TkReturn)
	node := &AstReturnNode{}
	if p.match(TkLParen, TkAdd, TkSubtract, TkIf, TkBool, TkNumber, TkString,
		TkIdentifier) {
		node.Expr = p.expr()
	}
	return node
}

// assign-stmt = ident ":=" expr
// func-call   = ident paren-expr-list
func (p *Parser) assignStmtOrFuncCall() AstNode {
	name := p.ident()
	if p.match(TkAssign) {
		p.advance()
		return &AstAssignNode{Name: name, Expr: p.expr()}
	}
	if p.match(TkLParen) {
		return &AstFuncCallNode{Name: name, Args: p.parenExprList()}
	}
	panic("unexpected " + p.curToken.String())
}

// identifier returns the lexeme value of the current identifier.
func (p *Parser) ident() string {
	defer p.consume(TkIdentifier)
	return p.curValue
}
