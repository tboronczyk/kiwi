package ast

import (
	"fmt"
	"github.com/tboronczyk/kiwi/token"
	"github.com/tboronczyk/kiwi/util"
	"strconv"
	"strings"
)

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

func (p AstPrinter) VisitValueExpr(n ValueExpr) {
	value := ""
	switch n.Type {
	case token.STRING:
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
		var f, _ = strconv.ParseFloat(n.Value, 64)
		value = fmt.Sprintf("%f", f)
		value = strings.TrimRight(value, "0")
		value = strings.TrimRight(value, ".")
		break
	case token.BOOL:
		value = strconv.FormatBool(strings.ToUpper(n.Value) == "TRUE")
	}
	fmt.Println("ValueExpr")
	fmt.Println(p.peek() + "├ Value: " + value)
	fmt.Println(p.peek() + "╰ Type: " + n.Type.String())
}

func (p AstPrinter) VisitCastExpr(n CastExpr) {
	fmt.Println("CastExpr")
	fmt.Println(p.peek() + "├ Cast: " + n.Cast)
	fmt.Print(p.peek() + "╰ Expr: ")
	p.push(p.peek() + "        ")
	n.Expr.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitVariableExpr(n VariableExpr) {
	fmt.Println("VariableExpr")
	fmt.Println(p.peek() + "╰ Name: " + n.Name)
}

func (p AstPrinter) VisitUnaryExpr(n UnaryExpr) {
	fmt.Println("UnaryExpr")
	fmt.Println(p.peek() + "├ Op: " + n.Op.String())
	fmt.Print(p.peek() + "╰ Right: ")
	p.push(p.peek() + "         ")
	n.Right.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitBinaryExpr(n BinaryExpr) {
	fmt.Println("BinaryExpr")
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

func (p AstPrinter) VisitFuncCall(n FuncCall) {
	fmt.Println("FuncCall")
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

func (p AstPrinter) VisitAssignStmt(n AssignStmt) {
	fmt.Println("AssignStmt")
	fmt.Println(p.peek() + "├ Name: " + n.Name)
	fmt.Print(p.peek() + "╰ Expr: ")
	p.push(p.peek() + "        ")
	n.Expr.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitFuncDef(n FuncDef) {
	fmt.Println("FuncDef")
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

func (p AstPrinter) VisitIfStmt(n IfStmt) {
	fmt.Println("IfStmt")
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

func (p AstPrinter) VisitReturnStmt(n ReturnStmt) {
	fmt.Println("ReturnStmt")
	fmt.Print(p.peek() + "╰ Expr: ")
	p.push(p.peek() + "        ")
	n.Expr.Accept(p)
	p.pop()
}

func (p AstPrinter) VisitWhileStmt(n WhileStmt) {
	fmt.Println("WhileStmt")
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
