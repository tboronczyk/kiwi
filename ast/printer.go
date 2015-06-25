package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tboronczyk/kiwi/token"
	"github.com/tboronczyk/kiwi/util"
)

// AstPrinter implements the NodeVisitor interface to traverse nodes in an
// abstract syntax tree and pretty-print the tree.
type AstPrinter struct {
	// stack is used to manage string padding for proper indenting.
	stack util.Stack
}

// NewAstPrinter returns a new AstPrinter.
func NewAstPrinter() AstPrinter {
	p := AstPrinter{stack: util.Stack{}}
	p.push("")
	return p
}

// push pushes indent padding string s onto the stack.
func (p *AstPrinter) push(s string) {
	p.stack.Push(s)
}

// peek returns the current padding string on the stack.
func (p AstPrinter) peek() string {
	return p.stack.Peek().(string)
}

// pop removes and returns the current padding string from the stack.
func (p *AstPrinter) pop() string {
	return p.stack.Pop().(string)
}

// VisitValueNode prints the value expression node n.
func (p AstPrinter) VisitValueNode(n *ValueNode) {
	value := ""
	switch n.Type {
	case token.STRING:
		// strings are presented in quotes with its special characters
		// escaped
		r := strings.NewReplacer(
			"\\\\", "\\",
			"\r", "\\r",
			"\n", "\\n",
			"\t", "\\t",
			"\"", "\\\"",
		)
		value = "\"" + r.Replace(n.Value) + "\""
		break
	case token.NUMBER:
		// numbers are presented as integers if they are whole, as
		// floats if they have a decimal.
		var f, _ = strconv.ParseFloat(n.Value, 64)
		value = fmt.Sprintf("%f", f)
		value = strings.TrimRight(value, "0")
		value = strings.TrimRight(value, ".")
		break
	case token.BOOL:
		value = strconv.FormatBool(strings.ToUpper(n.Value) == "TRUE")
	}
	fmt.Println("ValueNode")
	fmt.Println(p.peek() + "├ Value: " + value)
	fmt.Println(p.peek() + "╰ Type: " + n.Type.String())
}

// VisitCastNode prints the cast expression node n.
func (p AstPrinter) VisitCastNode(n *CastNode) {
	fmt.Println("CastNode")
	fmt.Println(p.peek() + "├ Cast: " + n.Cast)
	fmt.Print(p.peek() + "╰ Expr: ")
	p.push(p.peek() + "        ")
	n.Expr.Accept(p)
	p.pop()
}

// VisitVariableNode prints the variable expression node n.
func (p AstPrinter) VisitVariableNode(n *VariableNode) {
	fmt.Println("VariableNode")
	fmt.Println(p.peek() + "╰ Name: " + n.Name)
}

// VisitUnaryOpNode prints the unary operator expression node n.
func (p AstPrinter) VisitUnaryOpNode(n *UnaryOpNode) {
	fmt.Println("UnaryOpNode")
	fmt.Println(p.peek() + "├ Op: " + n.Op.String())
	fmt.Print(p.peek() + "╰ Expr: ")
	p.push(p.peek() + "        ")
	n.Expr.Accept(p)
	p.pop()
}

// VisitBinaryOpNode prints the binary operator expression node n.
func (p AstPrinter) VisitBinaryOpNode(n *BinaryOpNode) {
	fmt.Println("BinaryOpNode")
	fmt.Println(p.peek() + "├ Op: " + n.Op.String())
	fmt.Print(p.peek() + "├ Left: ")
	p.push(p.peek() + "│       ")
	n.Left.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

// VisitFuncCallNode prints the function call node n.
func (p AstPrinter) VisitFuncCallNode(n *FuncCallNode) {
	fmt.Println("FuncCallNode")
	fmt.Println(p.peek() + "├ Name: " + n.Name)
	fmt.Print(p.peek() + "╰ Args: ")
	if n.Args == nil {
		fmt.Println("␀")
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

// VisitAssignNode prints the assignment node n.
func (p AstPrinter) VisitAssignNode(n *AssignNode) {
	fmt.Println("AssignNode")
	fmt.Println(p.peek() + "├ Name: " + n.Name)
	fmt.Print(p.peek() + "╰ Expr: ")
	p.push(p.peek() + "        ")
	n.Expr.Accept(p)
	p.pop()
}

// VisitFuncDefNode prints the function definition node n.
func (p AstPrinter) VisitFuncDefNode(n *FuncDefNode) {
	fmt.Println("FuncDefNode")
	fmt.Println(p.peek() + "├ Name: " + n.Name)
	fmt.Print(p.peek() + "├ Args: ")
	if n.Args == nil {
		fmt.Println("␀")
	} else {
		fmt.Println(n.Args[0])
		if len(n.Args) > 1 {
			for _, arg := range n.Args[1:] {
				fmt.Println(p.peek() + "│       " + arg)
			}
		}
	}
	fmt.Print(p.peek() + "╰ Body: ")
	if n.Body == nil {
		fmt.Println("␀")
		return
	}
	p.push(p.peek() + "        ")
	n.Body[0].Accept(p)
	for _, arg := range n.Body[1:] {
		fmt.Print(p.peek())
		arg.Accept(p)
	}
	p.pop()
}

// VisitIfNode prints the if construct node n.
func (p AstPrinter) VisitIfNode(n *IfNode) {
	fmt.Println("IfNode")
	fmt.Print(p.peek() + "├ Condition: ")
	p.push(p.peek() + "│            ")
	n.Condition.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "├ Body: ")
	if n.Body == nil {
		fmt.Println("␀")
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
	if n.Else == nil {
		fmt.Println("␀")
		return
	}
	p.push(p.peek() + "        ")
	n.Else.Accept(p)
	p.pop()
}

// VisitReturnNode prints the return statement node n.
func (p AstPrinter) VisitReturnNode(n *ReturnNode) {
	fmt.Println("ReturnNode")
	fmt.Print(p.peek() + "╰ Expr: ")
	p.push(p.peek() + "        ")
	n.Expr.Accept(p)
	p.pop()
}

// VisitWhileNode prints the while construct node n.
func (p AstPrinter) VisitWhileNode(n *WhileNode) {
	fmt.Println("WhileNode")
	fmt.Print(p.peek() + "├ Condition: ")
	p.push(p.peek() + "│            ")
	n.Condition.Accept(p)
	p.pop()
	fmt.Print(p.peek() + "╰ Body: ")
	if n.Body == nil {
		fmt.Println("␀")
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
