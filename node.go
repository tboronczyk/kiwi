package main

type (
	AstNode interface {
		Accept(Visitor)
	}

	AstAddNode struct {
		Left  AstNode
		Right AstNode
	}

	AstAndNode struct {
		Left  AstNode
		Right AstNode
	}

	AstAssignNode struct {
		Name string
		Expr AstNode
	}

	AstBoolNode struct {
		Value bool
	}

	AstCastNode struct {
		Cast string
		Term AstNode
	}

	AstDivideNode struct {
		Left  AstNode
		Right AstNode
	}

	AstEqualNode struct {
		Left  AstNode
		Right AstNode
	}

	AstFuncCallNode struct {
		Name string
		Args []AstNode
	}

	AstFuncDefNode struct {
		*Scope
		Name string
		Args []string
		Body []AstNode
	}

	AstGreaterEqualNode struct {
		Left  AstNode
		Right AstNode
	}

	AstGreaterNode struct {
		Left  AstNode
		Right AstNode
	}

	AstIfNode struct {
		Cond AstNode
		Body []AstNode
		Else []AstNode
	}

	AstLessEqualNode struct {
		Left  AstNode
		Right AstNode
	}

	AstLessNode struct {
		Left  AstNode
		Right AstNode
	}

	AstModuloNode struct {
		Left  AstNode
		Right AstNode
	}

	AstMultiplyNode struct {
		Left  AstNode
		Right AstNode
	}

	AstNegativeNode struct {
		Term AstNode
	}

	AstNotEqualNode struct {
		Left  AstNode
		Right AstNode
	}

	AstNotNode struct {
		Term AstNode
	}

	AstNumberNode struct {
		Value float64
	}

	AstOrNode struct {
		Left  AstNode
		Right AstNode
	}

	AstPositiveNode struct {
		Term AstNode
	}

	AstProgramNode struct {
		*Scope
		Stmts []AstNode
	}

	AstReturnNode struct {
		Expr AstNode
	}

	AstStringNode struct {
		Value string
	}

	AstSubtractNode struct {
		Left  AstNode
		Right AstNode
	}

	AstVariableNode struct {
		Name string
	}

	AstWhileNode struct {
		Cond AstNode
		Body []AstNode
	}
)

func (n *AstAddNode) Accept(v Visitor) {
	v.VisitAddNode(n)
}

func (n *AstAndNode) Accept(v Visitor) {
	v.VisitAndNode(n)
}

func (n *AstAssignNode) Accept(v Visitor) {
	v.VisitAssignNode(n)
}

func (n *AstBoolNode) Accept(v Visitor) {
	v.VisitBoolNode(n)
}

func (n *AstCastNode) Accept(v Visitor) {
	v.VisitCastNode(n)
}

func (n *AstDivideNode) Accept(v Visitor) {
	v.VisitDivideNode(n)
}

func (n *AstEqualNode) Accept(v Visitor) {
	v.VisitEqualNode(n)
}

func (n *AstFuncCallNode) Accept(v Visitor) {
	v.VisitFuncCallNode(n)
}

func (n *AstFuncDefNode) Accept(v Visitor) {
	v.VisitFuncDefNode(n)
}

func (n *AstGreaterEqualNode) Accept(v Visitor) {
	v.VisitGreaterEqualNode(n)
}

func (n *AstGreaterNode) Accept(v Visitor) {
	v.VisitGreaterNode(n)
}

func (n *AstIfNode) Accept(v Visitor) {
	v.VisitIfNode(n)
}

func (n *AstLessEqualNode) Accept(v Visitor) {
	v.VisitLessEqualNode(n)
}

func (n *AstLessNode) Accept(v Visitor) {
	v.VisitLessNode(n)
}

func (n *AstModuloNode) Accept(v Visitor) {
	v.VisitModuloNode(n)
}

func (n *AstMultiplyNode) Accept(v Visitor) {
	v.VisitMultiplyNode(n)
}

func (n *AstNegativeNode) Accept(v Visitor) {
	v.VisitNegativeNode(n)
}

func (n *AstNotEqualNode) Accept(v Visitor) {
	v.VisitNotEqualNode(n)
}

func (n *AstNotNode) Accept(v Visitor) {
	v.VisitNotNode(n)
}

func (n *AstNumberNode) Accept(v Visitor) {
	v.VisitNumberNode(n)
}

func (n *AstOrNode) Accept(v Visitor) {
	v.VisitOrNode(n)
}

func (n *AstPositiveNode) Accept(v Visitor) {
	v.VisitPositiveNode(n)
}

func (n *AstProgramNode) Accept(v Visitor) {
	v.VisitProgramNode(n)
}

func (n *AstReturnNode) Accept(v Visitor) {
	v.VisitReturnNode(n)
}

func (n *AstStringNode) Accept(v Visitor) {
	v.VisitStringNode(n)
}

func (n *AstSubtractNode) Accept(v Visitor) {
	v.VisitSubtractNode(n)
}

func (n *AstVariableNode) Accept(v Visitor) {
	v.VisitVariableNode(n)
}

func (n *AstWhileNode) Accept(v Visitor) {
	v.VisitWhileNode(n)
}
