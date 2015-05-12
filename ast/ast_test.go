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

func capture(n *Node) string {
	// re-assign stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	n.PrintTree()

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

func TestPrintTreeUnary(t *testing.T) {
	n := NewNode(token.ADD, "+", 1)
	n.Children[0] = NewNode(token.NUMBER, "2", 0)
	expected := "\t[0] 2 (NUMBER)\n+ (+)\n"

	actual := capture(n)
	assert.Equal(t, expected, actual)
}

func TestPrintTreeBinary(t *testing.T) {
	n := NewNode(token.ADD, "+", 2)
	n.Children[0] = NewNode(token.NUMBER, "2", 0)
	n.Children[1] = NewNode(token.NUMBER, "4", 0)
	expected := "\t[0] 2 (NUMBER)\n+ (+)\n\t[1] 4 (NUMBER)\n"

	actual := capture(n)
	assert.Equal(t, expected, actual)
}

func TestPrintTreeArbitrary(t *testing.T) {
	n := NewNode(token.IF, "if", 3)
	n.Children[0] = NewNode(token.TRUE, "true", 0)
	n.Children[1] = NewNode(token.ASSIGN, ":=", 0)
	n.Children[2] = NewNode(token.ASSIGN, ":=", 0)
	expected := "\t[0] true (true)\nif (if)\n\t[1] := (:=)\n\t[2] := (:=)\n"

	actual := capture(n)
	assert.Equal(t, expected, actual)
}
