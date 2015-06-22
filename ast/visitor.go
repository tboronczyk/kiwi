package ast

// NodeVisitor is the basic type capable of visiting all node types.
type NodeVisitor interface {
	VisitAssignNode(AssignNode)
	VisitBinaryOpNode(BinaryOpNode)
	VisitCastNode(CastNode)
	VisitFuncCallNode(FuncCallNode)
	VisitFuncDefNode(FuncDefNode)
	VisitIfNode(IfNode)
	VisitReturnNode(ReturnNode)
	VisitUnaryOpNode(UnaryOpNode)
	VisitValueNode(ValueNode)
	VisitVariableNode(VariableNode)
	VisitWhileNode(WhileNode)
}
