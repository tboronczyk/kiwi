package ast

import (
	"github.com/tboronczyk/kiwi/token"
)

type Node interface {
	Accept(Visitor)
}

type Visitor interface {
	VisitCastExpr(CastExpr)
	VisitValueExpr(ValueExpr)
	VisitVariableExpr(VariableExpr)
	VisitUnaryExpr(UnaryExpr)
	VisitBinaryExpr(BinaryExpr)
	VisitFuncCall(FuncCall)
	VisitAssignStmt(AssignStmt)
	VisitFuncDef(FuncDef)
	VisitIfStmt(IfStmt)
	VisitReturnStmt(ReturnStmt)
	VisitWhileStmt(WhileStmt)
}

type CastExpr struct {
	Cast string
	Expr Node
}

func (n CastExpr) Accept(v Visitor) {
	v.VisitCastExpr(n)
}

type ValueExpr struct {
	Value string
	Type  token.Token
}

func (n ValueExpr) Accept(v Visitor) {
	v.VisitValueExpr(n)
}

type VariableExpr struct {
	Name string
}

func (n VariableExpr) Accept(v Visitor) {
	v.VisitVariableExpr(n)
}

type UnaryExpr struct {
	Op    token.Token
	Right Node
}

func (n UnaryExpr) Accept(v Visitor) {
	v.VisitUnaryExpr(n)
}

type BinaryExpr struct {
	Op    token.Token
	Left  Node
	Right Node
}

func (n BinaryExpr) Accept(v Visitor) {
	v.VisitBinaryExpr(n)
}

type FuncCall struct {
	Name string
	Args []Node
}

func (n FuncCall) Accept(v Visitor) {
	v.VisitFuncCall(n)
}

type AssignStmt struct {
	Name string
	Expr Node
}

func (n AssignStmt) Accept(v Visitor) {
	v.VisitAssignStmt(n)
}

type FuncDef struct {
	Name string
	Args []string
	Body []Node
}

func (n FuncDef) Accept(v Visitor) {
	v.VisitFuncDef(n)
}

type IfStmt struct {
	Condition Node
	Body      []Node
	Else      Node
}

func (n IfStmt) Accept(v Visitor) {
	v.VisitIfStmt(n)
}

type ReturnStmt struct {
	Expr Node
}

func (n ReturnStmt) Accept(v Visitor) {
	v.VisitReturnStmt(n)
}

type WhileStmt struct {
	Condition Node
	Body      []Node
}

func (n WhileStmt) Accept(v Visitor) {
	v.VisitWhileStmt(n)
}
