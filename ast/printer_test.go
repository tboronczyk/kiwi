package ast

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/token"
	"io"
	"os"
	"testing"
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

func TestValueExprNumber(t *testing.T) {
	expected := "ValueExpr\n" +
		"├ Value: 1\n" +
		"╰ Type: NUMBER\n"
	actual := capture(func() {
		n := ValueExpr{Value: "1.0", Type: token.NUMBER}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestValueExprString(t *testing.T) {
	expected := "ValueExpr\n" +
		"├ Value: \"foo\"\n" +
		"╰ Type: STRING\n"
	actual := capture(func() {
		n := ValueExpr{Value: "foo", Type: token.STRING}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestValueExprBool(t *testing.T) {
	expected := "ValueExpr\n" +
		"├ Value: true\n" +
		"╰ Type: BOOL\n"
	actual := capture(func() {
		n := ValueExpr{Value: "True", Type: token.BOOL}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestCast(t *testing.T) {
	expected := "CastExpr\n" +
		"├ Cast: string\n" +
		"╰ Expr: ValueExpr\n" +
		"        ├ Value: \"foo\"\n" +
		"        ╰ Type: STRING\n"
	actual := capture(func() {
		n := CastExpr{Cast: "string", Expr: ValueExpr{Value: "foo", Type: token.STRING}}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestVariableExpr(t *testing.T) {
	expected := "VariableExpr\n" +
		"╰ Name: foo\n"
	actual := capture(func() {
		n := VariableExpr{Name: "foo"}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestUnaryExpr(t *testing.T) {
	expected := "UnaryExpr\n" +
		"├ Op: -\n" +
		"╰ Right: VariableExpr\n" +
		"         ╰ Name: foo\n"
	actual := capture(func() {
		n := UnaryExpr{
			Op:    token.SUBTRACT,
			Right: VariableExpr{Name: "foo"},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestBinaryExpr(t *testing.T) {
	expected := "BinaryExpr\n" +
		"├ Op: +\n" +
		"├ Left: VariableExpr\n" +
		"│       ╰ Name: foo\n" +
		"╰ Right: VariableExpr\n" +
		"         ╰ Name: bar\n"
	actual := capture(func() {
		n := BinaryExpr{
			Op:    token.ADD,
			Left:  VariableExpr{Name: "foo"},
			Right: VariableExpr{Name: "bar"},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestAssignStmt(t *testing.T) {
	expected := "AssignStmt\n" +
		"├ Name: foo\n" +
		"╰ Expr: VariableExpr\n" +
		"        ╰ Name: bar\n"
	actual := capture(func() {
		n := AssignStmt{
			Name: "foo",
			Expr: VariableExpr{Name: "bar"},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestFuncDefNoArgsOrBody(t *testing.T) {
	expected := "FuncDef\n" +
		"├ Name: foo\n" +
		"├ Args: ␀\n" +
		"╰ Body: ␀\n"
	actual := capture(func() {
		n := FuncDef{Name: "foo"}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestFuncDef(t *testing.T) {
	expected := "FuncDef\n" +
		"├ Name: foo\n" +
		"├ Args: bar\n" +
		"│       baz\n" +
		"╰ Body: AssignStmt\n" +
		"        ├ Name: bar\n" +
		"        ╰ Expr: VariableExpr\n" +
		"                ╰ Name: baz\n" +
		"        ReturnStmt\n" +
		"        ╰ Expr: VariableExpr\n" +
		"                ╰ Name: baz\n"
	actual := capture(func() {
		n := FuncDef{
			Name: "foo",
			Args: []string{"bar", "baz"},
			Body: []Node{
				AssignStmt{
					Name: "bar",
					Expr: VariableExpr{Name: "baz"},
				},
				ReturnStmt{Expr: VariableExpr{Name: "baz"}},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestFuncCallNoArgs(t *testing.T) {
	expected := "FuncCall\n" +
		"├ Name: foo\n" +
		"╰ Args: ␀\n"
	actual := capture(func() {
		n := FuncCall{Name: "foo"}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestFuncCall(t *testing.T) {
	expected := "FuncCall\n" +
		"├ Name: foo\n" +
		"╰ Args: VariableExpr\n" +
		"        ╰ Name: foo\n" +
		"        VariableExpr\n" +
		"        ╰ Name: bar\n"
	actual := capture(func() {
		n := FuncCall{
			Name: "foo",
			Args: []Node{
				VariableExpr{Name: "foo"},
				VariableExpr{Name: "bar"},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestIfStmtNoBody(t *testing.T) {
	expected := "IfStmt\n" +
		"├ Condition: VariableExpr\n" +
		"│            ╰ Name: foo\n" +
		"├ Body: ␀\n" +
		"╰ Else: ␀\n"
	actual := capture(func() {
		n := IfStmt{Condition: VariableExpr{Name: "foo"}}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestIfStmt(t *testing.T) {
	expected := "IfStmt\n" +
		"├ Condition: VariableExpr\n" +
		"│            ╰ Name: foo\n" +
		"├ Body: AssignStmt\n" +
		"│       ├ Name: bar\n" +
		"│       ╰ Expr: VariableExpr\n" +
		"│               ╰ Name: baz\n" +
		"│       AssignStmt\n" +
		"│       ├ Name: quux\n" +
		"│       ╰ Expr: VariableExpr\n" +
		"│               ╰ Name: norf\n" +
		"╰ Else: IfStmt\n" +
		"        ├ Condition: ValueExpr\n" +
		"        │            ├ Value: true\n" +
		"        │            ╰ Type: BOOL\n" +
		"        ├ Body: ␀\n" +
		"        ╰ Else: ␀\n"
	actual := capture(func() {
		n := IfStmt{
			Condition: VariableExpr{Name: "foo"},
			Body: []Node{
				AssignStmt{
					Name: "bar",
					Expr: VariableExpr{Name: "baz"},
				},
				AssignStmt{
					Name: "quux",
					Expr: VariableExpr{Name: "norf"},
				},
			},
			Else: IfStmt{
				Condition: ValueExpr{
					Value: "true",
					Type:  token.BOOL,
				},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestReturnStmt(t *testing.T) {
	expected := "ReturnStmt\n" +
		"╰ Expr: VariableExpr\n" +
		"        ╰ Name: foo\n"
	actual := capture(func() {
		n := ReturnStmt{Expr: VariableExpr{Name: "foo"}}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestWhileStmtNoBody(t *testing.T) {
	expected := "WhileStmt\n" +
		"├ Condition: VariableExpr\n" +
		"│            ╰ Name: foo\n" +
		"╰ Body: ␀\n"
	actual := capture(func() {
		n := WhileStmt{Condition: VariableExpr{Name: "foo"}}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}

func TestWhileStmt(t *testing.T) {
	expected := "WhileStmt\n" +
		"├ Condition: VariableExpr\n" +
		"│            ╰ Name: foo\n" +
		"╰ Body: AssignStmt\n" +
		"        ├ Name: bar\n" +
		"        ╰ Expr: VariableExpr\n" +
		"                ╰ Name: baz\n" +
		"        AssignStmt\n" +
		"        ├ Name: quux\n" +
		"        ╰ Expr: VariableExpr\n" +
		"                ╰ Name: norf\n"
	actual := capture(func() {
		n := WhileStmt{
			Condition: VariableExpr{Name: "foo"},
			Body: []Node{
				AssignStmt{
					Name: "bar",
					Expr: VariableExpr{Name: "baz"},
				},
				AssignStmt{
					Name: "quux",
					Expr: VariableExpr{Name: "norf"},
				},
			},
		}
		n.Accept(NewAstPrinter())
	})
	assert.Equal(t, expected, actual)
}
