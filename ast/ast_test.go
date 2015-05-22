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
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/token"
	"io"
	"os"
	"testing"
)

func capture(n Node) string {
	// re-assign stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print(n)

	// capture output
	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		out <- buf.String()
	}()

	// restore stdout
	w.Close()
	os.Stdout = old

	return <-out
}

func TestLiteral(t *testing.T) {
	node := NewLiteral(token.IDENTIFIER, "foo")
	expected := "foo (IDENTIFIER)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestOperator(t *testing.T) {
	node := NewOperator(token.ADD)
	n := NewOperator(token.MULTIPLY)
	n.Left = NewLiteral(token.NUMBER, "2")
	n.Right = NewLiteral(token.NUMBER, "4")
	node.Left = n
	node.Right = NewLiteral(token.NUMBER, "8")
	expected := "OP.L OP.L 2 (NUMBER)\nOP.L OP *\nOP.L OP.R 4 (NUMBER)\nOP +\nOP.R 8 (NUMBER)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestList(t *testing.T) {
	node := NewList()
	n := NewOperator(token.ADD)
	n.Left = NewLiteral(token.NUMBER, "2")
	n.Right = NewLiteral(token.NUMBER, "4")
	node.Node = n
	n = NewOperator(token.SUBTRACT)
	n.Left = NewLiteral(token.NUMBER, "6")
	n.Right = NewLiteral(token.NUMBER, "8")
	node.Next = n
	expected := "L.N OP.L 2 (NUMBER)\nL.N OP +\nL.N OP.R 4 (NUMBER)\nL.n OP.L 6 (NUMBER)\nL.n OP -\nL.n OP.R 8 (NUMBER)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestFuncCall(t *testing.T) {
	node := NewFuncCall("foo")
	n := NewOperator(token.ADD)
	n.Left = NewLiteral(token.NUMBER, "2")
	n.Right = NewLiteral(token.NUMBER, "4")
	node.Body = n
	expected := "F.N foo\nF.B OP.L 2 (NUMBER)\nF.B OP +\nF.B OP.R 4 (NUMBER)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestIf(t *testing.T) {
	node := NewIf()
	n := NewOperator(token.EQUAL)
	n.Left = NewLiteral(token.IDENTIFIER, "foo")
	n.Right = NewLiteral(token.NUMBER, "true")
	node.Condition = n
	n = NewOperator(token.ASSIGN)
	n.Left = NewLiteral(token.IDENTIFIER, "bar")
	n.Right = NewLiteral(token.NUMBER, "false")
	node.Body = n
	expected := "IF.C OP.L foo (IDENTIFIER)\nIF.C OP =\nIF.C OP.R true (NUMBER)\nIF.B OP.L bar (IDENTIFIER)\nIF.B OP :=\nIF.B OP.R false (NUMBER)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestWhile(t *testing.T) {
	node := NewWhile()
	n := NewOperator(token.EQUAL)
	n.Left = NewLiteral(token.IDENTIFIER, "foo")
	n.Right = NewLiteral(token.NUMBER, "true")
	node.Condition = n
	n = NewOperator(token.ASSIGN)
	n.Left = NewLiteral(token.IDENTIFIER, "foo")
	n.Right = NewLiteral(token.NUMBER, "false")
	node.Body = n
	expected := "WL.C OP.L foo (IDENTIFIER)\nWL.C OP =\nWL.C OP.R true (NUMBER)\nWL.B OP.L foo (IDENTIFIER)\nWL.B OP :=\nWL.B OP.R false (NUMBER)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}


