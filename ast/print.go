package ast

import (
	"fmt"
)

func Print(n Node) {
	n.print("")
}

func (e ValueExpr) print(pad string) {
	fmt.Println("ValueExpr")
	fmt.Println(pad + "├ Value: " + e.Value)
	fmt.Println(pad + "╰ Type: " + e.Type.String())
}

func (e CastExpr) print(pad string) {
	fmt.Println("CastExpr")
	fmt.Println(pad + "├ Cast: " + e.Cast)
	fmt.Print(pad + "╰ Expr: ")
	e.Expr.print(pad + "        ")
}

func (e VariableExpr) print(pad string) {
	fmt.Println("VariableExpr")
	fmt.Println(pad + "╰ Name: " + e.Name)
}

func (e UnaryExpr) print(pad string) {
	fmt.Println("UnaryExpr")
	fmt.Println(pad + "├ Op: " + e.Op.String())
	fmt.Print(pad + "╰ Right: ")
	e.Right.print(pad + "         ")
}

func (e BinaryExpr) print(pad string) {
	fmt.Println("BinaryExpr")
	fmt.Println(pad + "├ Op: " + e.Op.String())
	fmt.Print(pad + "├ Left: ")
	e.Left.print(pad + "│       ")
	fmt.Print(pad + "╰ Right: ")
	e.Right.print(pad + "         ")
}

func (f FuncCall) print(pad string) {
	fmt.Println("FuncCall")
	fmt.Println(pad + "├ Name: " + f.Name)
	fmt.Print(pad + "╰ Args: ")
	if f.Args == nil {
		fmt.Println()
	} else {
		f.Args[0].print(pad + "        ")
		if len(f.Args) > 1 {
			for _, arg := range f.Args[1:] {
				fmt.Print(pad + "        ")
				arg.print(pad + "        ")
			}
		}
	}
}

func (s AssignStmt) print(pad string) {
	fmt.Println("AssignStmt")
	fmt.Println(pad + "├ Name: " + s.Name)
	fmt.Print(pad + "╰ Expr: ")
	s.Expr.print(pad + "        ")
}

func (s FuncDef) print(pad string) {
	fmt.Println("FuncDef")
	fmt.Println(pad + "├ Name: " + s.Name)
	fmt.Print(pad + "├ Args: ")
	if s.Args == nil {
		fmt.Println()
	} else {
		fmt.Println(s.Args[0])
		if len(s.Args) > 1 {
			for _, arg := range s.Args[1:] {
				fmt.Println(pad + "│       " + arg)
			}
		}
	}
	fmt.Print(pad + "╰ Body: ")
	if s.Body == nil {
		fmt.Println()
	} else {
		s.Body[0].print(pad + "        ")
		if len(s.Body) > 1 {
			for _, stmt := range s.Body[1:] {
				fmt.Print(pad + "        ")
				stmt.print(pad + "        ")
			}
		}
	}
}

func (s IfStmt) print(pad string) {
	fmt.Println("IfStmt")
	fmt.Print(pad + "├ Condition: ")
	s.Condition.print(pad + "│            ")
	fmt.Print(pad + "├ Body: ")
	if s.Body == nil {
		fmt.Println()
	} else {
		s.Body[0].print(pad + "│       ")
		if len(s.Body) > 1 {
			for _, stmt := range s.Body[1:] {
				fmt.Print(pad + "│       ")
				stmt.print(pad + "│       ")
			}
		}
	}
	fmt.Print(pad + "╰ Else: ")

	if s.Else == nil {
		fmt.Println()
	} else {
		s.Else.print(pad + "        ")
	}
}
func (s ReturnStmt) print(pad string) {
	fmt.Println("ReturnStmt")
	fmt.Print(pad + "╰ Expr: ")
	s.Expr.print(pad + "        ")
}

func (s WhileStmt) print(pad string) {
	fmt.Println("WhileStmt")
	fmt.Print(pad + "├ Condition: ")
	s.Condition.print(pad + "│            ")
	fmt.Print(pad + "╰ Body: ")
	if s.Body == nil {
		fmt.Println()
	} else {
		s.Body[0].print(pad + "        ")
		if len(s.Body) > 1 {
			for _, stmt := range s.Body[1:] {
				fmt.Print(pad + "        ")
				stmt.print(pad + "        ")
			}
		}
	}
}
