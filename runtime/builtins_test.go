package runtime

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/symtable"
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

func TestWrite(t *testing.T) {
	node := ast.FuncCall{
		Name: "write",
		Args: []ast.Node{
			ast.ValueExpr{
				Value: "foo",
				Type:  token.STRING,
			},
		},
	}
	funcTable := symtable.New()
	LoadBuiltins(funcTable)

	varTable := symtable.New()
	varTable.Set("value", "foo", symtable.STRING)

	expected := "foo"
	actual := capture(func() {
		Eval(node, varTable, funcTable)
	})
	assert.Equal(t, expected, actual)
}
