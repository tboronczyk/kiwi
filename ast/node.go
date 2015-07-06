package ast

import (
	"github.com/tboronczyk/kiwi/token"
)

// Node is the basic type for all AST nodes. The visitor pattern is used for
// node access.
type Node interface {
	Accept(NodeVisitor)
}

// CastNode represents a cast expression as an AST node.
type CastNode struct {
	Cast string
	Expr Node
}

// Accept visits the CastNode node using v.
func (n *CastNode) Accept(v NodeVisitor) {
	v.VisitCastNode(n)
}

// ValueNode represents a value expression as an AST node.
type ValueNode struct {
	Value string
	Type  token.Token
}

// Accept visits the value expression node using v.
func (n *ValueNode) Accept(v NodeVisitor) {
	v.VisitValueNode(n)
}

// VariableNode represents a variable expression as an AST node.
type VariableNode struct {
	Name string
}

// Accept visits the variable expression node using v.
func (n *VariableNode) Accept(v NodeVisitor) {
	v.VisitVariableNode(n)
}

// UnaryOpNode represents a unary operator expression as an AST node.
type UnaryOpNode struct {
	Op   token.Token
	Expr Node
}

// Accept visits the unary operator expression node using v.
func (n *UnaryOpNode) Accept(v NodeVisitor) {
	v.VisitUnaryOpNode(n)
}

// BinaryOpNode represents a binary operator expression as an AST node.
type BinaryOpNode struct {
	Op    token.Token
	Left  Node
	Right Node
}

// Accept visits the binary operator expression node using v.
func (n *BinaryOpNode) Accept(v NodeVisitor) {
	v.VisitBinaryOpNode(n)
}

// FuncCallNode represents a function call as an AST node.
type FuncCallNode struct {
	Name string
	Args []Node
}

// Accept visits the function call node using v.
func (n *FuncCallNode) Accept(v NodeVisitor) {
	v.VisitFuncCallNode(n)
}

// AssignNode represents an assignment operation as an AST.
type AssignNode struct {
	Name string
	Expr Node
}

// Accept visits the assignment node using v.
func (n *AssignNode) Accept(v NodeVisitor) {
	v.VisitAssignNode(n)
}

// FuncDefNode represents the defining of a function as an AST node.
type FuncDefNode struct {
	Name string
	Args []string
	Body []Node
}

// Accept visits the function definition node using v.
func (n *FuncDefNode) Accept(v NodeVisitor) {
	v.VisitFuncDefNode(n)
}

// IfNode represents an if construct as an AST node.
type IfNode struct {
	Condition Node
	Body      []Node
	Else      Node
}

// Accept visits the if construct node using v.
func (n *IfNode) Accept(v NodeVisitor) {
	v.VisitIfNode(n)
}

// ReturnNode represents a return statement as an AST node.
type ReturnNode struct {
	Expr Node
}

// Accept visits the return statement node using v.
func (n *ReturnNode) Accept(v NodeVisitor) {
	v.VisitReturnNode(n)
}

// WhileNode represents a while construct as an AST node.
type WhileNode struct {
	Condition Node
	Body      []Node
}

// Accept visits the while construct node using v.
func (n *WhileNode) Accept(v NodeVisitor) {
	v.VisitWhileNode(n)
}
