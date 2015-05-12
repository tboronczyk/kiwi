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
	"strconv"
)

type Node struct {
	token.Token
	Value    string
	Children []*Node
}

func NewNode(t interface{}, v interface{}, i uint8) *Node {
	tkn := token.UNKNOWN
	if t != nil {
		tkn = t.(token.Token)
	}

	val := ""
	if v != nil {
		val = v.(string)
	}

	return &Node{Token: tkn, Value: val, Children: make([]*Node, i)}
}

func (n Node) PrintTree() {
	n.printTree("")
}

func (n Node) printTree(s string) {
	if len(n.Children) > 0 {
		n.Children[0].printTree(s + "\t[0] ")
	}
	fmt.Printf("%s%s (%s)\n", s, n.Value, n.Token.String())
	if len(n.Children) > 1 {
		for i, node := range n.Children[1:] {
			node.printTree(s + "\t[" + strconv.Itoa(i+1) + "] ")
		}
	}
}
