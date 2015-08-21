package ast

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestPrintBoolNode(t *testing.T) {
	expected := "BoolNode\n" +
		"╰ Value: true\n"
	actual := capture(func() {
		n := &BoolNode{Value: true}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintNumberNode(t *testing.T) {
	expected := "NumberNode\n" +
		"╰ Value: 42\n"
	actual := capture(func() {
		n := &NumberNode{Value: 42.0}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintStringNode(t *testing.T) {
	expected := "StringNode\n" +
		"╰ Value: \"foo\"\n"
	actual := capture(func() {
		n := &StringNode{Value: "foo"}
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

func TestPrintAddNode(t *testing.T) {
	expected := "AddNode\n" +
		"├ Left: NumberNode\n" +
		"│       ╰ Value: 42\n" +
		"╰ Right: NumberNode\n" +
		"         ╰ Value: 73\n"
	actual := capture(func() {
		n := &AddNode{
			Left:  &NumberNode{Value: 42},
			Right: &NumberNode{Value: 73},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintAndNode(t *testing.T) {
	expected := "AndNode\n" +
		"├ Left: BoolNode\n" +
		"│       ╰ Value: true\n" +
		"╰ Right: BoolNode\n" +
		"         ╰ Value: false\n"
	actual := capture(func() {
		n := &AndNode{
			Left:  &BoolNode{Value: true},
			Right: &BoolNode{Value: false},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintAssignNode(t *testing.T) {
	expected := "AssignNode\n" +
		"├ Name: foo\n" +
		"╰ Expr: StringNode\n" +
		"        ╰ Value: \"bar\"\n"
	actual := capture(func() {
		n := &AssignNode{
			Name: "foo",
			Expr: &StringNode{Value: "bar"},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintCast(t *testing.T) {
	expected := "CastNode\n" +
		"├ Cast: number\n" +
		"╰ Term: StringNode\n" +
		"        ╰ Value: \"42\"\n"
	actual := capture(func() {
		n := &CastNode{
			Cast: "number",
			Term: &StringNode{Value: "42"},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintDivideNode(t *testing.T) {
	expected := "DivideNode\n" +
		"├ Left: NumberNode\n" +
		"│       ╰ Value: 42\n" +
		"╰ Right: NumberNode\n" +
		"         ╰ Value: 21\n"
	actual := capture(func() {
		n := &DivideNode{
			Left:  &NumberNode{Value: 42},
			Right: &NumberNode{Value: 21},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintEqualNode(t *testing.T) {
	expected := "EqualNode\n" +
		"├ Left: StringNode\n" +
		"│       ╰ Value: \"foo\"\n" +
		"╰ Right: StringNode\n" +
		"         ╰ Value: \"bar\"\n"
	actual := capture(func() {
		n := &EqualNode{
			Left:  &StringNode{Value: "foo"},
			Right: &StringNode{Value: "bar"},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintFuncCallNode(t *testing.T) {
	expected := "FuncCallNode\n" +
		"├ Name: foo\n" +
		"╰ Args: BoolNode\n" +
		"        ╰ Value: true\n" +
		"        NumberNode\n" +
		"        ╰ Value: 42\n"
	actual := capture(func() {
		n := &FuncCallNode{
			Name: "foo",
			Args: []Node{
				&BoolNode{Value: true},
				&NumberNode{Value: 42},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintFuncCallNodeNoArgs(t *testing.T) {
	expected := "FuncCallNode\n" +
		"├ Name: foo\n" +
		"╰ Args: 0x0\n"
	actual := capture(func() {
		n := &FuncCallNode{Name: "foo"}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintGreaterEqualNode(t *testing.T) {
	expected := "GreaterEqualNode\n" +
		"├ Left: NumberNode\n" +
		"│       ╰ Value: 1984\n" +
		"╰ Right: NumberNode\n" +
		"         ╰ Value: 1776\n"
	actual := capture(func() {
		n := &GreaterEqualNode{
			Left:  &NumberNode{Value: 1984},
			Right: &NumberNode{Value: 1776},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintGreaterNode(t *testing.T) {
	expected := "GreaterNode\n" +
		"├ Left: NumberNode\n" +
		"│       ╰ Value: 1984\n" +
		"╰ Right: NumberNode\n" +
		"         ╰ Value: 1776\n"
	actual := capture(func() {
		n := &GreaterNode{
			Left:  &NumberNode{Value: 1984},
			Right: &NumberNode{Value: 1776},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintIfNode(t *testing.T) {
	expected := "IfNode\n" +
		"├ Cond: BoolNode\n" +
		"│       ╰ Value: true\n" +
		"├ Body: AssignNode\n" +
		"│       ├ Name: foo\n" +
		"│       ╰ Expr: NumberNode\n" +
		"│               ╰ Value: 42\n" +
		"│       AssignNode\n" +
		"│       ├ Name: bar\n" +
		"│       ╰ Expr: StringNode\n" +
		"│               ╰ Value: \"baz\"\n" +
		"╰ Else: IfNode\n" +
		"        ├ Cond: BoolNode\n" +
		"        │       ╰ Value: true\n" +
		"        ├ Body: 0x0\n" +
		"        ╰ Else: 0x0\n"
	actual := capture(func() {
		n := &IfNode{
			Cond: &BoolNode{Value: true},
			Body: []Node{
				&AssignNode{
					Name: "foo",
					Expr: &NumberNode{Value: 42},
				},
				&AssignNode{
					Name: "bar",
					Expr: &StringNode{Value: "baz"},
				},
			},
			Else: &IfNode{
				Cond: &BoolNode{Value: true},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintIfNodeNoBody(t *testing.T) {
	expected := "IfNode\n" +
		"├ Cond: BoolNode\n" +
		"│       ╰ Value: true\n" +
		"├ Body: 0x0\n" +
		"╰ Else: 0x0\n"
	actual := capture(func() {
		n := &IfNode{Cond: &BoolNode{Value: true}}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintLessEqualNode(t *testing.T) {
	expected := "LessEqualNode\n" +
		"├ Left: NumberNode\n" +
		"│       ╰ Value: 1776\n" +
		"╰ Right: NumberNode\n" +
		"         ╰ Value: 1984\n"
	actual := capture(func() {
		n := &LessEqualNode{
			Left:  &NumberNode{Value: 1776},
			Right: &NumberNode{Value: 1984},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintLessNode(t *testing.T) {
	expected := "LessNode\n" +
		"├ Left: NumberNode\n" +
		"│       ╰ Value: 1776\n" +
		"╰ Right: NumberNode\n" +
		"         ╰ Value: 1984\n"
	actual := capture(func() {
		n := &LessNode{
			Left:  &NumberNode{Value: 1776},
			Right: &NumberNode{Value: 1984},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintModuloNode(t *testing.T) {
	expected := "ModuloNode\n" +
		"├ Left: NumberNode\n" +
		"│       ╰ Value: 11\n" +
		"╰ Right: NumberNode\n" +
		"         ╰ Value: 7\n"
	actual := capture(func() {
		n := &ModuloNode{
			Left:  &NumberNode{Value: 11},
			Right: &NumberNode{Value: 7},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintMultiplyNode(t *testing.T) {
	expected := "MultiplyNode\n" +
		"├ Left: NumberNode\n" +
		"│       ╰ Value: 21\n" +
		"╰ Right: NumberNode\n" +
		"         ╰ Value: 2\n"
	actual := capture(func() {
		n := &MultiplyNode{
			Left:  &NumberNode{Value: 21},
			Right: &NumberNode{Value: 2},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintNegativeNode(t *testing.T) {
	expected := "NegativeNode\n" +
		"╰ Term: NumberNode\n" +
		"        ╰ Value: 42\n"
	actual := capture(func() {
		n := &NegativeNode{
			Term: &NumberNode{Value: 42},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintNotEqualNode(t *testing.T) {
	expected := "NotEqualNode\n" +
		"├ Left: StringNode\n" +
		"│       ╰ Value: \"foo\"\n" +
		"╰ Right: StringNode\n" +
		"         ╰ Value: \"bar\"\n"
	actual := capture(func() {
		n := &NotEqualNode{
			Left:  &StringNode{Value: "foo"},
			Right: &StringNode{Value: "bar"},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintNotNode(t *testing.T) {
	expected := "NotNode\n" +
		"╰ Term: BoolNode\n" +
		"        ╰ Value: false\n"
	actual := capture(func() {
		n := &NotNode{
			Term: &BoolNode{Value: false},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintOrNode(t *testing.T) {
	expected := "OrNode\n" +
		"├ Left: BoolNode\n" +
		"│       ╰ Value: true\n" +
		"╰ Right: BoolNode\n" +
		"         ╰ Value: false\n"
	actual := capture(func() {
		n := &OrNode{
			Left:  &BoolNode{Value: true},
			Right: &BoolNode{Value: false},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintPositiveNode(t *testing.T) {
	expected := "PositiveNode\n" +
		"╰ Term: NumberNode\n" +
		"        ╰ Value: 42\n"
	actual := capture(func() {
		n := &PositiveNode{
			Term: &NumberNode{Value: 42},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintProgramNode(t *testing.T) {
	expected := "ProgramNode\n" +
		"╰ Stmts: AssignNode\n" +
		"         ├ Name: foo\n" +
		"         ╰ Expr: VariableNode\n" +
		"                 ╰ Name: bar\n" +
		"         AssignNode\n" +
		"         ├ Name: baz\n" +
		"         ╰ Expr: VariableNode\n" +
		"                 ╰ Name: quux\n"
	actual := capture(func() {
		n := &ProgramNode{
			Stmts: []Node{
				&AssignNode{
					Name: "foo",
					Expr: &VariableNode{Name: "bar"},
				},
				&AssignNode{
					Name: "baz",
					Expr: &VariableNode{Name: "quux"},
				},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintProgramNodeNoStmts(t *testing.T) {
	expected := "ProgramNode\n" +
		"╰ Stmts: 0x0\n"
	actual := capture(func() {
		n := &ProgramNode{}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintReturnNode(t *testing.T) {
	expected := "ReturnNode\n" +
		"╰ Expr: BoolNode\n" +
		"        ╰ Value: true\n"
	actual := capture(func() {
		n := &ReturnNode{Expr: &BoolNode{Value: true}}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintSubtractNode(t *testing.T) {
	expected := "SubtractNode\n" +
		"├ Left: NumberNode\n" +
		"│       ╰ Value: 73\n" +
		"╰ Right: NumberNode\n" +
		"         ╰ Value: 42\n"
	actual := capture(func() {
		n := &SubtractNode{
			Left:  &NumberNode{Value: 73},
			Right: &NumberNode{Value: 42},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintWhileNode(t *testing.T) {
	expected := "WhileNode\n" +
		"├ Cond: BoolNode\n" +
		"│       ╰ Value: true\n" +
		"╰ Body: AssignNode\n" +
		"        ├ Name: foo\n" +
		"        ╰ Expr: NumberNode\n" +
		"                ╰ Value: 42\n" +
		"        AssignNode\n" +
		"        ├ Name: bar\n" +
		"        ╰ Expr: StringNode\n" +
		"                ╰ Value: \"baz\"\n"
	actual := capture(func() {
		n := &WhileNode{
			Cond: &BoolNode{Value: true},
			Body: []Node{
				&AssignNode{
					Name: "foo",
					Expr: &NumberNode{Value: 42},
				},
				&AssignNode{
					Name: "bar",
					Expr: &StringNode{Value: "baz"},
				},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestPrintWhileNodeNoBody(t *testing.T) {
	expected := "WhileNode\n" +
		"├ Cond: BoolNode\n" +
		"│       ╰ Value: true\n" +
		"╰ Body: 0x0\n"
	actual := capture(func() {
		n := &WhileNode{Cond: &BoolNode{Value: true}}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}
