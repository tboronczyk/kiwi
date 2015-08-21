package ast

type Visitor interface {
	VisitAssignNode(*AssignNode)
	VisitBinOpNode(*BinOpNode)
	VisitCastNode(*CastNode)
	VisitFuncCallNode(*FuncCallNode)
	VisitFuncDefNode(*FuncDefNode)
	VisitIfNode(*IfNode)
	VisitProgramNode(*ProgramNode)
	VisitReturnNode(*ReturnNode)
	VisitUnaryOpNode(*UnaryOpNode)
	VisitValueNode(*ValueNode)
	VisitVariableNode(*VariableNode)
	VisitWhileNode(*WhileNode)
}
