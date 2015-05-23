/*
 * Copyright (c) 2012, 2015 Timothy Boronczyk
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  1. Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *
 *  2. Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *
 *  3. The names of the authors may not be used to endorse or promote products
 *     derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED "AS IS" AND WITHOUT ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, WITHOUT LIMITATION, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE.
 */

package ast

import (
	"fmt"
	"github.com/tboronczyk/kiwi/token"
)

type Node interface {
	print(string)
}

type (
	Literal struct {
		Type  token.Token
		Value string
	}

	Operator struct {
		Op    token.Token
		Left  Node
		Right Node
	}

	List struct {
		Node
		Prev Node
	}

	FuncCall struct {
		Name Literal
		Args Node
	}

	FuncDef struct {
		Name   Literal
		Params Node
		Body   Node
	}

	If struct {
		Condition Node
		Body      Node
	}

	Return struct {
		Expr Node
	}

	While struct {
		Condition Node
		Body      Node
	}
)

func Print(n Node) {
	n.print("")
}

func (n Literal) print(s string) {
	fmt.Printf("%s%s (%s)\n", s, n.Value, n.Type.String())
}

func (n Operator) print(s string) {
	if n.Left != nil {
		n.Left.print(s + "OP.L ")
	}
	fmt.Printf("%sOP %s\n", s, n.Op.String())
	if n.Right != nil {
		n.Right.print(s + "OP.R ")
	}
}

func (n FuncCall) print(s string) {
	n.Name.print(s + "FC.N ")
	if n.Args != nil {
		n.Args.print(s + "FC.A ")
	}
}

func (n List) print(s string) {
	if n.Prev != nil {
		n.Prev.print(s + "L.P ")
	}
	if n.Node != nil {
		n.Node.print(s + "L.N ")
	}
}

func (n FuncDef) print(s string) {
	n.Name.print(s + "FD.N ")
	if n.Params != nil {
		n.Params.print(s + "FD.P ")
	}
	if n.Body != nil {
		n.Body.print(s + "FD.B ")
	}
}

func (n If) print(s string) {
	if n.Condition != nil {
		n.Condition.print(s + "IF.C ")
	}
	if n.Body != nil {
		n.Body.print(s + "IF.B ")
	}
}

func (n Return) print(s string) {
	if n.Expr == nil {
		fmt.Printf("%sRet\n", s)
	} else {
		n.Expr.print(s + "Ret.E ")
	}
}

func (n While) print(s string) {
	if n.Condition != nil {
		n.Condition.print(s + "WL.C ")
	}
	if n.Body != nil {
		n.Body.print(s + "WL.B ")
	}
}
