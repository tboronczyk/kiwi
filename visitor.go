package main

type Visitor interface {
	VisitAddNode(*AstAddNode)
	VisitAndNode(*AstAndNode)
	VisitAssignNode(*AstAssignNode)
	VisitBoolNode(*AstBoolNode)
	VisitCastNode(*AstCastNode)
	VisitDivideNode(*AstDivideNode)
	VisitEqualNode(*AstEqualNode)
	VisitFuncCallNode(*AstFuncCallNode)
	VisitFuncDefNode(*AstFuncDefNode)
	VisitGreaterEqualNode(*AstGreaterEqualNode)
	VisitGreaterNode(*AstGreaterNode)
	VisitIfNode(*AstIfNode)
	VisitLessEqualNode(*AstLessEqualNode)
	VisitLessNode(*AstLessNode)
	VisitModuloNode(*AstModuloNode)
	VisitMultiplyNode(*AstMultiplyNode)
	VisitNegativeNode(*AstNegativeNode)
	VisitNotEqualNode(*AstNotEqualNode)
	VisitNotNode(*AstNotNode)
	VisitNumberNode(*AstNumberNode)
	VisitOrNode(*AstOrNode)
	VisitPositiveNode(*AstPositiveNode)
	VisitProgramNode(*AstProgramNode)
	VisitReturnNode(*AstReturnNode)
	VisitStringNode(*AstStringNode)
	VisitSubtractNode(*AstSubtractNode)
	VisitVariableNode(*AstVariableNode)
	VisitWhileNode(*AstWhileNode)
}