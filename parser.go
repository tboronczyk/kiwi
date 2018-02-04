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
	p := &Parser{scanner: s, scope: NewScope()}
	p.advance()
	return p
}

// advance retrieves the next token/value pair from the scanner. COMMENT tokens
// are skipped as whitespace.
func (p *Parser) advance() {
	for {
		p.curToken, p.curValue = p.scanner.Scan()
		if p.curToken != T_COMMENT {
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
		if p.curToken == T_EOF {
			return prog, nil
		}
		prog.Stmts = append(prog.Stmts, p.stmt())
	}
}

// expr = cmp-expr [log-op expr]
func (p *Parser) expr() AstNode {
	node := p.cmpExpr()
	switch p.curToken {
	case T_AND:
		p.advance()
		return &AstAndNode{Left: node, Right: p.expr()}
	case T_OR:
		p.advance()
		return &AstOrNode{Left: node, Right: p.expr()}
	}
	return node
}

// cmp-expr = add-expr [cmp-op cmp-expr]
func (p *Parser) cmpExpr() AstNode {
	node := p.addExpr()
	switch p.curToken {
	case T_EQUAL:
		p.advance()
		return &AstEqualNode{Left: node, Right: p.cmpExpr()}
	case T_NOT_EQUAL:
		p.advance()
		return &AstNotEqualNode{Left: node, Right: p.cmpExpr()}
	case T_GREATER:
		p.advance()
		return &AstGreaterNode{Left: node, Right: p.cmpExpr()}
	case T_GREATER_EQ:
		p.advance()
		return &AstGreaterEqualNode{Left: node, Right: p.cmpExpr()}
	case T_LESS:
		p.advance()
		return &AstLessNode{Left: node, Right: p.cmpExpr()}
	case T_LESS_EQ:
		p.advance()
		return &AstLessEqualNode{Left: node, Right: p.cmpExpr()}
	}
	return node
}

// add-expr = mul-expr [add-op add-expr]
func (p *Parser) addExpr() AstNode {
	node := p.mulExpr()
	switch p.curToken {
	case T_ADD:
		p.advance()
		return &AstAddNode{Left: node, Right: p.addExpr()}
	case T_SUBTRACT:
		p.advance()
		return &AstSubtractNode{Left: node, Right: p.addExpr()}
	}
	return node
}

// mul-expr = cast-expr [mul-op mul-expr]
func (p *Parser) mulExpr() AstNode {
	node := p.castExpr()
	switch p.curToken {
	case T_MULTIPLY:
		p.advance()
		return &AstMultiplyNode{Left: node, Right: p.mulExpr()}
	case T_DIVIDE:
		p.advance()
		return &AstDivideNode{Left: node, Right: p.mulExpr()}
	case T_MODULO:
		p.advance()
		return &AstModuloNode{Left: node, Right: p.mulExpr()}
	}
	return node
}

// cast-expr = term [":" ident]
func (p *Parser) castExpr() AstNode {
	node := p.term()
	if p.curToken == T_COLON {
		p.advance()
		return &AstCastNode{Cast: p.ident(), Term: node}
	}
	return node
}

// term = "(" expr ")" / ("+" / "-" / "~") term / boolean / number / string /
//        func-call / ident
func (p *Parser) term() AstNode {
	switch p.curToken {
	case T_LPAREN:
		p.advance()
		node := p.expr()
		p.consume(T_RPAREN)
		return node
	case T_ADD:
		p.advance()
		return &AstPositiveNode{Term: p.term()}
	case T_SUBTRACT:
		p.advance()
		return &AstNegativeNode{Term: p.term()}
	case T_NOT:
		p.advance()
		return &AstNotNode{Term: p.term()}
	case T_BOOL:
		node := &AstBoolNode{Value: strings.ToLower(p.curValue) == "true"}
		p.advance()
		return node
	case T_NUMBER:
		val, _ := strconv.ParseFloat(p.curValue, 64)
		node := &AstNumberNode{Value: val}
		p.advance()
		return node
	case T_STRING:
		node := &AstStringNode{Value: p.curValue}
		p.advance()
		return node
	case T_IDENTIFIER:
		name := p.ident()
		if p.match(T_LPAREN) {
			return &AstFuncCallNode{Name: name, Args: p.parenExprList()}
		}
		return &AstVariableNode{Name: name}
	}

	panic("whoops?")
}

// paren-expr-list = "(" [expr *("," expr)] ")"
func (p *Parser) parenExprList() []AstNode {
	defer p.consume(T_RPAREN)
	p.consume(T_LPAREN)

	var list []AstNode
	if p.match(T_RPAREN) {
		return list
	}
	for {
		list = append(list, p.expr())
		if !p.match(T_COMMA) {
			return list
		}
		p.advance()
	}
}

// stmt = if-stmt / while-stmt / func-def / return-stmt / assign-stmt /
//        func-call
func (p *Parser) stmt() (node AstNode) {
	switch p.curToken {
	case T_IF:
		return p.ifStmt()
	case T_WHILE:
		return p.whileStmt()
	case T_FUNC:
		return p.funcDef()
	case T_RETURN:
		return p.returnStmt()
	case T_IDENTIFIER:
		return p.assignStmtOrFuncCall()
	}
	panic("unexpected lexeme " + p.curToken.String())
}

// if-stmt = "if" expr brace-stmt-list [else-clause]
func (p *Parser) ifStmt() *AstIfNode {
	p.consume(T_IF)
	node := &AstIfNode{Cond: p.expr(), Body: p.braceStmtList()}
	if p.match(T_ELSE) {
		p.advance()
		if p.match(T_LBRACE) {
			node.Else = p.braceStmtList()
		} else {
			node.Else = append(node.Else, p.elseClause())
		}
	}
	return node
}

// brace-stmt-list = "{" *stmt "}"
func (p *Parser) braceStmtList() (list []AstNode) {
	p.consume(T_LBRACE)
	for {
		if !(p.curToken.IsStmtKeyword() || p.match(T_IDENTIFIER)) {
			break
		}
		list = append(list, p.stmt())
	}
	p.consume(T_RBRACE)
	return list
}

// else-clause = "else" (brace-stmt-list / expr brace-stmt-list else-clause)
// Note: an else with an expression becomes an if-stmt within a default else
// clause.
func (p *Parser) elseClause() *AstIfNode {
	node := &AstIfNode{Cond: p.expr(), Body: p.braceStmtList()}
	if p.match(T_ELSE) {
		p.advance()
		if p.match(T_LBRACE) {
			node.Else = p.braceStmtList()
		} else {
			node.Else = append(node.Else, p.elseClause())
		}
	}
	return node
}

// while-stmt = "while" expr brace-stmt-list
func (p *Parser) whileStmt() *AstWhileNode {
	p.consume(T_WHILE)
	return &AstWhileNode{Cond: p.expr(), Body: p.braceStmtList()}
}

// func-def = "func" ident *ident brace-stmt-list
func (p *Parser) funcDef() *AstFuncDefNode {
	p.consume(T_FUNC)

	node := &AstFuncDefNode{
		Name: p.ident(),
	}
	p.scope.SetFunc(node.Name, Entry{Value: node, DataType: FUNC})
	node.Scope = NewScopeWithParent(p.scope)
	p.scope = node.Scope

	if !p.match(T_LBRACE) {
		var list []string
		for {
			list = append(list, p.ident())
			if !p.match(T_IDENTIFIER) {
				break
			}
		}
		node.Args = list
	}
	node.Body = p.braceStmtList()

	p.scope = node.Scope.Parent
	return node
}

// return-stmt = "return" [expr]
func (p *Parser) returnStmt() *AstReturnNode {
	p.consume(T_RETURN)
	node := &AstReturnNode{}
	if p.match(T_LPAREN, T_ADD, T_SUBTRACT, T_NOT, T_BOOL, T_NUMBER, T_STRING,
		T_IDENTIFIER) {
		node.Expr = p.expr()
	}
	return node
}

// assign-stmt = ident ":=" expr
// func-call   = ident paren-expr-list
func (p *Parser) assignStmtOrFuncCall() AstNode {
	name := p.ident()
	if p.match(T_ASSIGN) {
		p.advance()
		return &AstAssignNode{Name: name, Expr: p.expr()}
	}
	if p.match(T_LPAREN) {
		return &AstFuncCallNode{Name: name, Args: p.parenExprList()}
	}
	panic("unexpected " + p.curToken.String())
}

// identifier returns the lexeme value of the current identifier.
func (p *Parser) ident() string {
	defer p.consume(T_IDENTIFIER)
	return p.curValue
}
