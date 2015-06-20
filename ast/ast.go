package ast

import (
	"github.com/tboronczyk/kiwi/token"
)

type (
	Node interface {
		print(string)
	}

	CastExpr struct {
		Cast string
		Expr Node
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
		Right Node
	}

	BinaryExpr struct {
		Op    token.Token
		Left  Node
		Right Node
	}

	FuncCall struct {
		Name string
		Args []Node
	}

	AssignStmt struct {
		Name string
		Expr Node
	}

	FuncDef struct {
		Name string
		Args []string
		Body []Node
	}

	IfStmt struct {
		Condition Node
		Body      []Node
		Else      Node
	}

	ReturnStmt struct {
		Expr Node
	}

	WhileStmt struct {
		Condition Node
		Body      []Node
	}
)
