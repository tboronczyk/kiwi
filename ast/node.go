package ast

import (
	"github.com/tboronczyk/kiwi/scope"
)

type (
	Node interface {
		Accept(Visitor)
	}

	AddNode struct {
		Left  Node
		Right Node
	}

	AndNode struct {
		Left  Node
		Right Node
	}

	AssignNode struct {
		Name string
		Expr Node
	}

	BoolNode struct {
		Value bool
	}

	CastNode struct {
		Cast string
		Term Node
	}

	DivideNode struct {
		Left  Node
		Right Node
	}

	EqualNode struct {
		Left  Node
		Right Node
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

	GreaterEqualNode struct {
		Left  Node
		Right Node
	}

	GreaterNode struct {
		Left  Node
		Right Node
	}

	IfNode struct {
		Cond Node
		Body []Node
		Else []Node
	}

	LessEqualNode struct {
		Left  Node
		Right Node
	}

	LessNode struct {
		Left  Node
		Right Node
	}

	ModuloNode struct {
		Left  Node
		Right Node
	}

	MultiplyNode struct {
		Left  Node
		Right Node
	}

	NegativeNode struct {
		Term Node
	}

	NotEqualNode struct {
		Left  Node
		Right Node
	}

	NotNode struct {
		Term Node
	}

	NumberNode struct {
		Value float64
	}

	OrNode struct {
		Left  Node
		Right Node
	}

	PositiveNode struct {
		Term Node
	}

	ProgramNode struct {
		*scope.Scope
		Stmts []Node
	}

	ReturnNode struct {
		Expr Node
	}

	StringNode struct {
		Value string
	}

	SubtractNode struct {
		Left  Node
		Right Node
	}

	VariableNode struct {
		Name string
	}

	WhileNode struct {
		Cond Node
		Body []Node
	}
)

func (n *AddNode) Accept(v Visitor) {
	v.VisitAddNode(n)
}

func (n *AndNode) Accept(v Visitor) {
	v.VisitAndNode(n)
}

func (n *AssignNode) Accept(v Visitor) {
	v.VisitAssignNode(n)
}

func (n *BoolNode) Accept(v Visitor) {
	v.VisitBoolNode(n)
}

func (n *CastNode) Accept(v Visitor) {
	v.VisitCastNode(n)
}

func (n *DivideNode) Accept(v Visitor) {
	v.VisitDivideNode(n)
}

func (n *EqualNode) Accept(v Visitor) {
	v.VisitEqualNode(n)
}

func (n *FuncCallNode) Accept(v Visitor) {
	v.VisitFuncCallNode(n)
}

func (n *FuncDefNode) Accept(v Visitor) {
	v.VisitFuncDefNode(n)
}

func (n *GreaterEqualNode) Accept(v Visitor) {
	v.VisitGreaterEqualNode(n)
}

func (n *GreaterNode) Accept(v Visitor) {
	v.VisitGreaterNode(n)
}

func (n *IfNode) Accept(v Visitor) {
	v.VisitIfNode(n)
}

func (n *LessEqualNode) Accept(v Visitor) {
	v.VisitLessEqualNode(n)
}

func (n *LessNode) Accept(v Visitor) {
	v.VisitLessNode(n)
}

func (n *ModuloNode) Accept(v Visitor) {
	v.VisitModuloNode(n)
}

func (n *MultiplyNode) Accept(v Visitor) {
	v.VisitMultiplyNode(n)
}

func (n *NegativeNode) Accept(v Visitor) {
	v.VisitNegativeNode(n)
}

func (n *NotEqualNode) Accept(v Visitor) {
	v.VisitNotEqualNode(n)
}

func (n *NotNode) Accept(v Visitor) {
	v.VisitNotNode(n)
}

func (n *NumberNode) Accept(v Visitor) {
	v.VisitNumberNode(n)
}

func (n *OrNode) Accept(v Visitor) {
	v.VisitOrNode(n)
}

func (n *PositiveNode) Accept(v Visitor) {
	v.VisitPositiveNode(n)
}

func (n *ProgramNode) Accept(v Visitor) {
	v.VisitProgramNode(n)
}

func (n *ReturnNode) Accept(v Visitor) {
	v.VisitReturnNode(n)
}

func (n *StringNode) Accept(v Visitor) {
	v.VisitStringNode(n)
}

func (n *SubtractNode) Accept(v Visitor) {
	v.VisitSubtractNode(n)
}

func (n *VariableNode) Accept(v Visitor) {
	v.VisitVariableNode(n)
}

func (n *WhileNode) Accept(v Visitor) {
	v.VisitWhileNode(n)
}
