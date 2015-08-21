package ast

import (
	"github.com/tboronczyk/kiwi/scope"
	"github.com/tboronczyk/kiwi/token"
)

type (
	Node interface {
		Accept(Visitor)
	}

	AssignNode struct {
		Name string
		Expr Node
	}

	BinOpNode struct {
		Op    token.Token
		Left  Node
		Right Node
	}

	CastNode struct {
		Cast string
		Term Node
	}

	FuncCallNode struct {
		Name string
		Args []Node
	}

	FuncDefNode struct {
		*scope.Scope
		Name string
		Args []string
		Body []Node
	}

	IfNode struct {
		Cond Node
		Body []Node
		Else Node
	}

	ProgramNode struct {
		*scope.Scope
		Stmts []Node
	}

	ReturnNode struct {
		Expr Node
	}

	UnaryOpNode struct {
		Op   token.Token
		Term Node
	}

	ValueNode struct {
		Value string
		Type  token.Token
	}

	VariableNode struct {
		Name string
	}

	WhileNode struct {
		Cond Node
		Body []Node
	}
)

func (n *AssignNode) Accept(v Visitor) {
	v.VisitAssignNode(n)
}

func (n *BinOpNode) Accept(v Visitor) {
	v.VisitBinOpNode(n)
}

func (n *CastNode) Accept(v Visitor) {
	v.VisitCastNode(n)
}

func (n *FuncCallNode) Accept(v Visitor) {
	v.VisitFuncCallNode(n)
}

func (n *FuncDefNode) Accept(v Visitor) {
	v.VisitFuncDefNode(n)
}

func (n *IfNode) Accept(v Visitor) {
	v.VisitIfNode(n)
}

func (n *ProgramNode) Accept(v Visitor) {
	v.VisitProgramNode(n)
}

func (n *ReturnNode) Accept(v Visitor) {
	v.VisitReturnNode(n)
}

func (n *UnaryOpNode) Accept(v Visitor) {
	v.VisitUnaryOpNode(n)
}

func (n *ValueNode) Accept(v Visitor) {
	v.VisitValueNode(n)
}

func (n *VariableNode) Accept(v Visitor) {
	v.VisitVariableNode(n)
}

func (n *WhileNode) Accept(v Visitor) {
	v.VisitWhileNode(n)
}
