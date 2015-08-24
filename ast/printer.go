package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tboronczyk/kiwi/util"
)

// AstPrinter implements the Visitor interface to traverse AST nodes to pretty
// print the tree.
type AstPrinter struct {
	stack util.Stack
}

func NewAstPrinter() AstPrinter {
	p := AstPrinter{stack: util.Stack{}}
	p.push("")
	return p
}

func (p *AstPrinter) push(s string) {
	p.stack.Push(s)
}

func (p AstPrinter) peek() string {
	return p.stack.Peek().(string)
}

func (p *AstPrinter) pop() string {
	return p.stack.Pop().(string)
}

func (p AstPrinter) VisitAddNode(n *AddNode) {
	fmt.Println("AddNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitAndNode(n *AndNode) {
	fmt.Println("AndNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitAssignNode(n *AssignNode) {
	fmt.Println("AssignNode")
	fmt.Println(p.peek() + "├ Name: " + n.Name)
	fmt.Print(p.peek() + "╰ Expr: ")
	p.push(p.peek() + "        ")
	n.Expr.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitBoolNode(n *BoolNode) {
	fmt.Println("BoolNode")
	value := strconv.FormatBool(n.Value)
	fmt.Println(p.peek() + "╰ Value: " + value)
}

func (p AstPrinter) VisitCastNode(n *CastNode) {
	fmt.Println("CastNode")
	fmt.Println(p.peek() + "├ Cast: " + n.Cast)
	fmt.Print(p.peek() + "╰ Term: ")
	p.push(p.peek() + "        ")
	n.Term.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitDivideNode(n *DivideNode) {
	fmt.Println("DivideNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitEqualNode(n *EqualNode) {
	fmt.Println("EqualNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitFuncCallNode(n *FuncCallNode) {
	fmt.Println("FuncCallNode")
	fmt.Println(p.peek() + "├ Name: " + n.Name)
	fmt.Print(p.peek() + "╰ Args: ")
	if n.Args == nil || len(n.Args) == 0 {
		fmt.Println("0x0")
		return
	}
	p.push(p.peek() + "        ")
	n.Args[0].Accept(p)
	for _, arg := range n.Args[1:] {
		fmt.Print(p.peek())
		arg.Accept(p)
	}
	p.pop()
}

func (p AstPrinter) VisitFuncDefNode(n *FuncDefNode) {
	fmt.Println("FuncDefNode")
	fmt.Println(p.peek() + "├ Name: " + n.Name)
	fmt.Print(p.peek() + "├ Args: ")
	if n.Args == nil || len(n.Args) == 0 {
		fmt.Println("0x0")
	} else {
		fmt.Println(n.Args[0])
		if len(n.Args) > 1 {
			for _, arg := range n.Args[1:] {
				fmt.Println(p.peek() + "│       " + arg)
			}
		}
	}
	fmt.Print(p.peek() + "╰ Body: ")
	if n.Body == nil || len(n.Body) == 0 {
		fmt.Println("0x0")
	} else {
		p.push(p.peek() + "        ")
		n.Body[0].Accept(p)
		for _, arg := range n.Body[1:] {
			fmt.Print(p.peek())
			arg.Accept(p)
		}
		p.pop()
	}
}

func (p AstPrinter) VisitGreaterEqualNode(n *GreaterEqualNode) {
	fmt.Println("GreaterEqualNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitGreaterNode(n *GreaterNode) {
	fmt.Println("GreaterNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitIfNode(n *IfNode) {
	fmt.Println("IfNode")
	fmt.Print(p.peek() + "├ Cond: ")
	p.push(p.peek() + "│       ")
	n.Cond.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "├ Body: ")
	if n.Body == nil || len(n.Body) == 0 {
		fmt.Println("0x0")
	} else {
		p.push(p.peek() + "│       ")
		n.Body[0].Accept(p)
		if len(n.Body) > 1 {
			for _, stmt := range n.Body[1:] {
				fmt.Print(p.peek())
				stmt.Accept(p)
			}
		}
		p.pop()
	}
	fmt.Print(p.peek() + "╰ Else: ")
	if n.Else == nil || len(n.Else) == 0 {
		fmt.Println("0x0")
	} else {
		p.push(p.peek() + "        ")
		n.Else[0].Accept(p)
		if len(n.Else) > 1 {
			for _, stmt := range n.Else[1:] {
				fmt.Print(p.peek())
				stmt.Accept(p)
			}
		}
		p.pop()
	}
}

func (p AstPrinter) VisitLessEqualNode(n *LessEqualNode) {
	fmt.Println("LessEqualNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitLessNode(n *LessNode) {
	fmt.Println("LessNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitModuloNode(n *ModuloNode) {
	fmt.Println("ModuloNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitMultiplyNode(n *MultiplyNode) {
	fmt.Println("MultiplyNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitNegativeNode(n *NegativeNode) {
	fmt.Println("NegativeNode")
	fmt.Print(p.peek() + "╰ Term: ")
	p.push(p.peek() + "        ")
	n.Term.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitNotEqualNode(n *NotEqualNode) {
	fmt.Println("NotEqualNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitNotNode(n *NotNode) {
	fmt.Println("NotNode")
	fmt.Print(p.peek() + "╰ Term: ")
	p.push(p.peek() + "        ")
	n.Term.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitNumberNode(n *NumberNode) {
	fmt.Println("NumberNode")
	// numbers are presented as integers if they are whole, as
	// floats if they have a decimal
	value := fmt.Sprintf("%f", n.Value)
	value = strings.TrimRight(value, "0")
	value = strings.TrimRight(value, ".")
	fmt.Println(p.peek() + "╰ Value: " + value)
}

func (p AstPrinter) VisitOrNode(n *OrNode) {
	fmt.Println("OrNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitPositiveNode(n *PositiveNode) {
	fmt.Println("PositiveNode")
	fmt.Print(p.peek() + "╰ Term: ")
	p.push(p.peek() + "        ")
	n.Term.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitProgramNode(n *ProgramNode) {
	fmt.Println("ProgramNode")
	fmt.Print(p.peek() + "╰ Stmts: ")
	if n.Stmts == nil || len(n.Stmts) == 0 {
		fmt.Println("0x0")
	} else {
		p.push(p.peek() + "         ")
		n.Stmts[0].Accept(p)
		if len(n.Stmts) > 1 {
			for _, stmt := range n.Stmts[1:] {
				fmt.Print(p.peek())
				stmt.Accept(p)
			}
		}
		p.pop()
	}
}

func (p AstPrinter) VisitReturnNode(n *ReturnNode) {
	fmt.Println("ReturnNode")
	fmt.Print(p.peek() + "╰ Expr: ")
	p.push(p.peek() + "        ")
	n.Expr.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitStringNode(n *StringNode) {
	fmt.Println("StringNode")
	// strings are presented in quotes with its special characters
	// escaped
	r := strings.NewReplacer(
		"\\\\", "\\",
		"\r", "\\r",
		"\n", "\\n",
		"\t", "\\t",
		"\"", "\\\"",
	)
	value := "\"" + r.Replace(n.Value) + "\""
	fmt.Println(p.peek() + "╰ Value: " + value)
}

func (p AstPrinter) VisitSubtractNode(n *SubtractNode) {
	fmt.Println("SubtractNode")
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitVariableNode(n *VariableNode) {
	fmt.Println("VariableNode")
	fmt.Println(p.peek() + "╰ Name: " + n.Name)
}

func (p AstPrinter) VisitWhileNode(n *WhileNode) {
	fmt.Println("WhileNode")
	fmt.Print(p.peek() + "├ Cond: ")
	p.push(p.peek() + "│       ")
	n.Cond.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Body: ")
	if n.Body == nil || len(n.Body) == 0 {
		fmt.Println("0x0")
		return
	}
	p.push(p.peek() + "        ")
	n.Body[0].Accept(p)
	if len(n.Body) > 1 {
		for _, stmt := range n.Body[1:] {
			fmt.Print(p.peek())
			stmt.Accept(p)
		}
	}
	p.pop()
}
