package ast

import (
	"github.com/tboronczyk/kiwi/token"
)

type (
	Node interface {
		print(string)
	}

	ExprNode interface {
		Node
	}

	StmtNode interface {
		Node
	}

	ValueExpr struct {
		Value string
		Type  token.Token
	}

	VariableExpr struct {
		Name string
	}

	UnaryExpr struct {
		Op    token.Token
		Right ExprNode
	}

	BinaryExpr struct {
		Op    token.Token
		Left  ExprNode
		Right ExprNode
	}

	FuncCall struct {
		Name string
		Args []ExprNode
	}

	AssignStmt struct {
		Name string
		Expr ExprNode
	}

	FuncDef struct {
		Name string
		Args []string
		Body []StmtNode
	}

	IfStmt struct {
		Condition ExprNode
		Body      []StmtNode
	}

	ReturnStmt struct {
		Expr ExprNode
	}

	WhileStmt struct {
		Condition ExprNode
		Body      []StmtNode
	}
)
