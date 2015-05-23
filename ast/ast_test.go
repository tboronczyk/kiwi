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
	node := Literal{Type: token.IDENTIFIER, Value: "foo"}
	expected := "foo (IDENTIFIER)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestOperator(t *testing.T) {
	node := Operator{
		Op: token.ADD,
		Left: Operator{
			Op:    token.MULTIPLY,
			Left:  Literal{Type: token.NUMBER, Value: "2"},
			Right: Literal{Type: token.NUMBER, Value: "4"}},
		Right: Literal{Type: token.NUMBER, Value: "8"}}
	expected := "OP.L OP.L 2 (NUMBER)\nOP.L OP *\nOP.L OP.R 4 (NUMBER)\nOP +\nOP.R 8 (NUMBER)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestList(t *testing.T) {
	node := List{
		Node: Operator{
			Op:    token.ADD,
			Left:  Literal{Type: token.NUMBER, Value: "2"},
			Right: Literal{Type: token.NUMBER, Value: "4"}},
		Prev: List{
			Node: Operator{
				Op:    token.SUBTRACT,
				Left:  Literal{Type: token.NUMBER, Value: "6"},
				Right: Literal{Type: token.NUMBER, Value: "8"}}}}

	expected := "L.P L.N OP.L 6 (NUMBER)\nL.P L.N OP -\nL.P L.N OP.R 8 (NUMBER)\nL.N OP.L 2 (NUMBER)\nL.N OP +\nL.N OP.R 4 (NUMBER)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestFuncDef(t *testing.T) {
	node := FuncDef{
		Name: Literal{Type: token.IDENTIFIER, Value: "foo"},
		Params: List{
			Node: Literal{Type: token.IDENTIFIER, Value: "c"},
			Prev: List{
				Node: Literal{Type: token.IDENTIFIER, Value: "b"},
				Prev: List{
					Node: Literal{Type: token.IDENTIFIER, Value: "a"}}}},
		Body: List{
			Node: Return{}}}
	expected := "FD.N foo (IDENTIFIER)\nFD.P L.P L.P L.N a (IDENTIFIER)\nFD.P L.P L.N b (IDENTIFIER)\nFD.P L.N c (IDENTIFIER)\nFD.B L.N Ret\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestFuncCall(t *testing.T) {
	node := FuncCall{
		Name: Literal{Type: token.IDENTIFIER, Value: "foo"},
		Args: Operator{
			Op:    token.ADD,
			Left:  Literal{Type: token.NUMBER, Value: "2"},
			Right: Literal{Type: token.NUMBER, Value: "4"}}}
	expected := "FC.N foo (IDENTIFIER)\nFC.A OP.L 2 (NUMBER)\nFC.A OP +\nFC.A OP.R 4 (NUMBER)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestIf(t *testing.T) {
	node := If{
		Condition: Operator{
			Op:    token.EQUAL,
			Left:  Literal{Type: token.IDENTIFIER, Value: "foo"},
			Right: Literal{Type: token.TRUE, Value: "true"}},
		Body: Operator{
			Op:    token.ASSIGN,
			Left:  Literal{Type: token.IDENTIFIER, Value: "bar"},
			Right: Literal{Type: token.FALSE, Value: "false"}}}
	expected := "IF.C OP.L foo (IDENTIFIER)\nIF.C OP =\nIF.C OP.R true (true)\nIF.B OP.L bar (IDENTIFIER)\nIF.B OP :=\nIF.B OP.R false (false)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestReturn(t *testing.T) {
	node := Return{
		Expr: Operator{
			Op:    token.EQUAL,
			Left:  Literal{Type: token.IDENTIFIER, Value: "foo"},
			Right: Literal{Type: token.TRUE, Value: "true"}}}
	expected := "Ret.E OP.L foo (IDENTIFIER)\nRet.E OP =\nRet.E OP.R true (true)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestReturnNoExpr(t *testing.T) {
	node := Return{}
	expected := "Ret\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}

func TestWhile(t *testing.T) {
	node := While{
		Condition: Operator{
			Op:    token.EQUAL,
			Left:  Literal{Type: token.IDENTIFIER, Value: "foo"},
			Right: Literal{Type: token.TRUE, Value: "true"}},
		Body: Operator{
			Op:    token.ASSIGN,
			Left:  Literal{Type: token.IDENTIFIER, Value: "foo"},
			Right: Literal{Type: token.FALSE, Value: "false"}}}
	expected := "WL.C OP.L foo (IDENTIFIER)\nWL.C OP =\nWL.C OP.R true (true)\nWL.B OP.L foo (IDENTIFIER)\nWL.B OP :=\nWL.B OP.R false (false)\n"

	actual := capture(node)
	assert.Equal(t, expected, actual)
}
