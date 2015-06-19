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

func TestValueExpr(t *testing.T) {
	expr := ValueExpr{Value: "foo", Type: token.STRING}
	expected := "ValueExpr\n" +
		"├ Value: foo\n" +
		"╰ Type: STRING\n"

	actual := capture(expr)
	assert.Equal(t, expected, actual)
}

func TestCast(t *testing.T) {
	expr := CastExpr{Cast: "string", Expr: ValueExpr{Value: "foo", Type: token.STRING}}
	expected := "CastExpr\n" +
		"├ Cast: string\n" +
		"╰ Expr: ValueExpr\n" +
		"        ├ Value: foo\n" +
		"        ╰ Type: STRING\n"

	actual := capture(expr)
	assert.Equal(t, expected, actual)
}

func TestVariableExpr(t *testing.T) {
	expr := VariableExpr{Name: "foo"}
	expected := "VariableExpr\n" +
		"╰ Name: foo\n"

	actual := capture(expr)
	assert.Equal(t, expected, actual)
}

func TestUnaryExpr(t *testing.T) {
	expr := UnaryExpr{
		Op:    token.SUBTRACT,
		Right: VariableExpr{Name: "foo"},
	}
	expected := "UnaryExpr\n" +
		"├ Op: -\n" +
		"╰ Right: VariableExpr\n" +
		"         ╰ Name: foo\n"

	actual := capture(expr)
	assert.Equal(t, expected, actual)
}

func TestBinaryExpr(t *testing.T) {
	expr := BinaryExpr{
		Op:    token.ADD,
		Left:  VariableExpr{Name: "foo"},
		Right: VariableExpr{Name: "bar"},
	}
	expected := "BinaryExpr\n" +
		"├ Op: +\n" +
		"├ Left: VariableExpr\n" +
		"│       ╰ Name: foo\n" +
		"╰ Right: VariableExpr\n" +
		"         ╰ Name: bar\n"

	actual := capture(expr)
	assert.Equal(t, expected, actual)
}

func TestAssignStmt(t *testing.T) {
	stmt := AssignStmt{
		Name: "foo",
		Expr: VariableExpr{Name: "bar"},
	}
	expected := "AssignStmt\n" +
		"├ Name: foo\n" +
		"╰ Expr: VariableExpr\n" +
		"        ╰ Name: bar\n"

	actual := capture(stmt)
	assert.Equal(t, expected, actual)
}

func TestFuncDefNoArgsOrBody(t *testing.T) {
	stmt := FuncDef{Name: "foo"}
	expected := "FuncDef\n" +
		"├ Name: foo\n" +
		"├ Args: \n" +
		"╰ Body: \n"

	actual := capture(stmt)
	assert.Equal(t, expected, actual)
}

func TestFuncDef(t *testing.T) {
	stmt := FuncDef{
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

	actual := capture(stmt)
	assert.Equal(t, expected, actual)
}

func TestFuncCallNoArgs(t *testing.T) {
	expr := FuncCall{Name: "foo"}
	expected := "FuncCall\n" +
		"├ Name: foo\n" +
		"╰ Args: \n"

	actual := capture(expr)
	assert.Equal(t, expected, actual)
}

func TestFuncCall(t *testing.T) {
	expr := FuncCall{
		Name: "foo",
		Args: []Node{
			VariableExpr{Name: "foo"},
			VariableExpr{Name: "bar"},
		},
	}
	expected := "FuncCall\n" +
		"├ Name: foo\n" +
		"╰ Args: VariableExpr\n" +
		"        ╰ Name: foo\n" +
		"        VariableExpr\n" +
		"        ╰ Name: bar\n"

	actual := capture(expr)
	assert.Equal(t, expected, actual)
}

func TestIfStmtNoBody(t *testing.T) {
	stmt := IfStmt{Condition: VariableExpr{Name: "foo"}}
	expected := "IfStmt\n" +
		"├ Condition: VariableExpr\n" +
		"│            ╰ Name: foo\n" +
		"╰ Body: \n"

	actual := capture(stmt)
	assert.Equal(t, expected, actual)
}

func TestIfStmt(t *testing.T) {
	stmt := IfStmt{
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
	expected := "IfStmt\n" +
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

	actual := capture(stmt)
	assert.Equal(t, expected, actual)
}

func TestReturnStmt(t *testing.T) {
	stmt := ReturnStmt{Expr: VariableExpr{Name: "foo"}}
	expected := "ReturnStmt\n" +
		"╰ Expr: VariableExpr\n" +
		"        ╰ Name: foo\n"

	actual := capture(stmt)
	assert.Equal(t, expected, actual)
}

func TestWhileStmtNoBody(t *testing.T) {
	stmt := WhileStmt{Condition: VariableExpr{Name: "foo"}}
	expected := "WhileStmt\n" +
		"├ Condition: VariableExpr\n" +
		"│            ╰ Name: foo\n" +
		"╰ Body: \n"

	actual := capture(stmt)
	assert.Equal(t, expected, actual)
}

func TestWhileStmt(t *testing.T) {
	stmt := WhileStmt{
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

	actual := capture(stmt)
	assert.Equal(t, expected, actual)
}
