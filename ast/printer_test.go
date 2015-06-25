package ast

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/token"
)

func capture(f func()) string {
	// re-assign stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()
	// read output
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

func TestPrintValueNodeNumber(t *testing.T) {
	expected := "ValueNode\n" +
		"├ Value: 1\n" +
		"╰ Type: NUMBER\n"
	actual := capture(func() {
		n := &ValueNode{Value: "1.0", Type: token.NUMBER}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintValueNodeString(t *testing.T) {
	expected := "ValueNode\n" +
		"├ Value: \"foo\"\n" +
		"╰ Type: STRING\n"
	actual := capture(func() {
		n := &ValueNode{Value: "foo", Type: token.STRING}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintValueNodeBool(t *testing.T) {
	expected := "ValueNode\n" +
		"├ Value: true\n" +
		"╰ Type: BOOL\n"
	actual := capture(func() {
		n := &ValueNode{Value: "True", Type: token.BOOL}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintCast(t *testing.T) {
	expected := "CastNode\n" +
		"├ Cast: string\n" +
		"╰ Expr: ValueNode\n" +
		"        ├ Value: \"foo\"\n" +
		"        ╰ Type: STRING\n"
	actual := capture(func() {
		n := &CastNode{
			Cast: "string",
			Expr: &ValueNode{Value: "foo", Type: token.STRING},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintVariableNode(t *testing.T) {
	expected := "VariableNode\n" +
		"╰ Name: foo\n"
	actual := capture(func() {
		n := &VariableNode{Name: "foo"}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintUnaryOpNode(t *testing.T) {
	expected := "UnaryOpNode\n" +
		"├ Op: -\n" +
		"╰ Right: VariableNode\n" +
		"         ╰ Name: foo\n"
	actual := capture(func() {
		n := &UnaryOpNode{
			Op:    token.SUBTRACT,
			Right: &VariableNode{Name: "foo"},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintBinaryOpNode(t *testing.T) {
	expected := "BinaryOpNode\n" +
		"├ Op: +\n" +
		"├ Left: VariableNode\n" +
		"│       ╰ Name: foo\n" +
		"╰ Right: VariableNode\n" +
		"         ╰ Name: bar\n"
	actual := capture(func() {
		n := &BinaryOpNode{
			Op:    token.ADD,
			Left:  &VariableNode{Name: "foo"},
			Right: &VariableNode{Name: "bar"},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintAssignNode(t *testing.T) {
	expected := "AssignNode\n" +
		"├ Name: foo\n" +
		"╰ Expr: VariableNode\n" +
		"        ╰ Name: bar\n"
	actual := capture(func() {
		n := &AssignNode{
			Name: "foo",
			Expr: &VariableNode{Name: "bar"},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintFuncDefNodeNoArgsOrBody(t *testing.T) {
	expected := "FuncDefNode\n" +
		"├ Name: foo\n" +
		"├ Args: ␀\n" +
		"╰ Body: ␀\n"
	actual := capture(func() {
		n := &FuncDefNode{Name: "foo"}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintFuncDefNode(t *testing.T) {
	expected := "FuncDefNode\n" +
		"├ Name: foo\n" +
		"├ Args: bar\n" +
		"│       baz\n" +
		"╰ Body: AssignNode\n" +
		"        ├ Name: bar\n" +
		"        ╰ Expr: VariableNode\n" +
		"                ╰ Name: baz\n" +
		"        ReturnNode\n" +
		"        ╰ Expr: VariableNode\n" +
		"                ╰ Name: baz\n"
	actual := capture(func() {
		n := &FuncDefNode{
			Name: "foo",
			Args: []string{"bar", "baz"},
			Body: []Node{
				&AssignNode{
					Name: "bar",
					Expr: &VariableNode{Name: "baz"},
				},
				&ReturnNode{Expr: &VariableNode{Name: "baz"}},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintFuncCallNodeNoArgs(t *testing.T) {
	expected := "FuncCallNode\n" +
		"├ Name: foo\n" +
		"╰ Args: ␀\n"
	actual := capture(func() {
		n := &FuncCallNode{Name: "foo"}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintFuncCallNode(t *testing.T) {
	expected := "FuncCallNode\n" +
		"├ Name: foo\n" +
		"╰ Args: VariableNode\n" +
		"        ╰ Name: foo\n" +
		"        VariableNode\n" +
		"        ╰ Name: bar\n"
	actual := capture(func() {
		n := &FuncCallNode{
			Name: "foo",
			Args: []Node{
				&VariableNode{Name: "foo"},
				&VariableNode{Name: "bar"},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintIfNodeNoBody(t *testing.T) {
	expected := "IfNode\n" +
		"├ Condition: VariableNode\n" +
		"│            ╰ Name: foo\n" +
		"├ Body: ␀\n" +
		"╰ Else: ␀\n"
	actual := capture(func() {
		n := &IfNode{Condition: &VariableNode{Name: "foo"}}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintIfNode(t *testing.T) {
	expected := "IfNode\n" +
		"├ Condition: VariableNode\n" +
		"│            ╰ Name: foo\n" +
		"├ Body: AssignNode\n" +
		"│       ├ Name: bar\n" +
		"│       ╰ Expr: VariableNode\n" +
		"│               ╰ Name: baz\n" +
		"│       AssignNode\n" +
		"│       ├ Name: quux\n" +
		"│       ╰ Expr: VariableNode\n" +
		"│               ╰ Name: norf\n" +
		"╰ Else: IfNode\n" +
		"        ├ Condition: ValueNode\n" +
		"        │            ├ Value: true\n" +
		"        │            ╰ Type: BOOL\n" +
		"        ├ Body: ␀\n" +
		"        ╰ Else: ␀\n"
	actual := capture(func() {
		n := &IfNode{
			Condition: &VariableNode{Name: "foo"},
			Body: []Node{
				&AssignNode{
					Name: "bar",
					Expr: &VariableNode{Name: "baz"},
				},
				&AssignNode{
					Name: "quux",
					Expr: &VariableNode{Name: "norf"},
				},
			},
			Else: &IfNode{
				Condition: &ValueNode{
					Value: "true",
					Type:  token.BOOL,
				},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintReturnNode(t *testing.T) {
	expected := "ReturnNode\n" +
		"╰ Expr: VariableNode\n" +
		"        ╰ Name: foo\n"
	actual := capture(func() {
		n := &ReturnNode{Expr: &VariableNode{Name: "foo"}}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintWhileNodeNoBody(t *testing.T) {
	expected := "WhileNode\n" +
		"├ Condition: VariableNode\n" +
		"│            ╰ Name: foo\n" +
		"╰ Body: ␀\n"
	actual := capture(func() {
		n := &WhileNode{Condition: &VariableNode{Name: "foo"}}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintWhileNode(t *testing.T) {
	expected := "WhileNode\n" +
		"├ Condition: VariableNode\n" +
		"│            ╰ Name: foo\n" +
		"╰ Body: AssignNode\n" +
		"        ├ Name: bar\n" +
		"        ╰ Expr: VariableNode\n" +
		"                ╰ Name: baz\n" +
		"        AssignNode\n" +
		"        ├ Name: quux\n" +
		"        ╰ Expr: VariableNode\n" +
		"                ╰ Name: norf\n"
	actual := capture(func() {
		n := &WhileNode{
			Condition: &VariableNode{Name: "foo"},
			Body: []Node{
				&AssignNode{
					Name: "bar",
					Expr: &VariableNode{Name: "baz"},
				},
				&AssignNode{
					Name: "quux",
					Expr: &VariableNode{Name: "norf"},
				},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}
