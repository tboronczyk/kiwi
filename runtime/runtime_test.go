package runtime

import (
	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/symtable"
	"github.com/tboronczyk/kiwi/token"
	"testing"
)

func TestEvalAssignNode(t *testing.T) {
	node := ast.AssignNode{
		Name: "foo",
		Expr: ast.ValueNode{
			Value: "bar",
			Type:  token.STRING,
		},
	}

	r := New()
	node.Accept(r)

	value, dtype, _ := r.varGet("foo")
	assert.Equal(t, "bar", value)
	assert.Equal(t, symtable.STRING, dtype)
}

func TestEvalFuncDefNode(t *testing.T) {
	node := ast.FuncDefNode{
		Name: "foo",
	}

	r := New()
	node.Accept(r)

	value, dtype, _ := r.funcGet("foo")
	assert.Equal(t, node, value)
	assert.Equal(t, symtable.USRFUNC, dtype)
}

func TestEvalFuncCallNodeUserDefined(t *testing.T) {
	r := New()
	r.funcSet(
		"foo",
		ast.FuncDefNode{
			Args: []string{},
			Body: []ast.Node{
				ast.AssignNode{
					Name: "bar",
					Expr: ast.ValueNode{
						Value: "baz",
						Type:  token.STRING,
					},
				},
			},
		},
		symtable.USRFUNC)

	node := ast.FuncCallNode{
		Name: "foo",
		Args: []ast.Node{},
	}

	node.Accept(r)

	assert.Panics(t, func() {
		r.popStack()
	})
}

func TestEvalFuncCallNodeBuiltin(t *testing.T) {
	r := New()
	r.funcSet(
		"foo",
		ast.FuncDefNode{Name: "foo", Args: []string{}},
		symtable.BUILTIN,
	)
	builtinFuncs["foo"] = func(r *Runtime) {
		r.pushStack(true, symtable.BOOL)
	}

	node := ast.FuncCallNode{
		Name: "foo",
		Args: []ast.Node{},
	}

	node.Accept(r)

	actual := r.popStack()
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalFuncCallNodeUserDefinedWithArgs(t *testing.T) {
	r := New()
	r.funcSet(
		"foo",
		ast.FuncDefNode{
			Args: []string{"bar"},
			Body: []ast.Node{
				ast.AssignNode{
					Name: "baz",
					Expr: ast.VariableNode{Name: "bar"},
				},
			},
		},
		symtable.USRFUNC)

	node := ast.FuncCallNode{
		Name: "foo",
		Args: []ast.Node{
			ast.ValueNode{
				Value: "bar",
				Type:  token.STRING,
			},
		},
	}

	node.Accept(r)

	assert.Panics(t, func() {
		r.popStack()
	})
}

func TestEvalFuncCallNodeUserDefinedWithReturn(t *testing.T) {
	r := New()
	r.funcSet(
		"foo",
		ast.FuncDefNode{
			Args: []string{"bar"},
			Body: []ast.Node{
				ast.ReturnNode{
					Expr: ast.VariableNode{Name: "bar"},
				},
			},
		},
		symtable.USRFUNC)

	node := ast.FuncCallNode{
		Name: "foo",
		Args: []ast.Node{
			ast.ValueNode{
				Value: "bar",
				Type:  token.STRING,
			},
		},
	}

	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "bar", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
}

func TestEvalFuncCallNodeFuncNoExist(t *testing.T) {
	node := ast.FuncCallNode{Name: "foo"}

	r := New()

	assert.Panics(t, func() {
		node.Accept(r)
	})
}

func TestEvalFuncCallNodeArityMismatch(t *testing.T) {
	r := New()
	r.funcSet("foo", ast.FuncDefNode{}, symtable.USRFUNC)

	node := ast.FuncCallNode{
		Name: "foo",
		Args: []ast.Node{
			ast.ValueNode{
				Value: "bar",
				Type:  token.STRING,
			},
		},
	}

	assert.Panics(t, func() {
		node.Accept(r)
	})
}

func TestEvalIfNodeTrue(t *testing.T) {
	node := ast.IfNode{
		Condition: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
		Body: []ast.Node{
			ast.AssignNode{
				Name: "foo",
				Expr: ast.ValueNode{
					Value: "bar",
					Type:  token.STRING,
				},
			},
		},
	}

	r := New()
	node.Accept(r)

	value, dtype, _ := r.varGet("foo")
	assert.Equal(t, "bar", value)
	assert.Equal(t, symtable.STRING, dtype)
}

func TestEvalIfNodeFalse(t *testing.T) {
	node := ast.IfNode{
		Condition: ast.ValueNode{
			Value: "false",
			Type:  token.BOOL,
		},
		Else: ast.IfNode{
			Condition: ast.ValueNode{
				Value: "true",
				Type:  token.BOOL,
			},
			Body: []ast.Node{
				ast.AssignNode{
					Name: "foo",
					Expr: ast.ValueNode{
						Value: "bar",
						Type:  token.STRING,
					},
				},
			},
		},
	}

	r := New()
	node.Accept(r)

	value, dtype, _ := r.varGet("foo")
	assert.Equal(t, "bar", value)
	assert.Equal(t, symtable.STRING, dtype)
}

func TestEvalIfNodeNotBoolCondition(t *testing.T) {
	node := ast.IfNode{
		Condition: ast.ValueNode{
			Value: "foo",
			Type:  token.STRING,
		},
	}

	assert.Panics(t, func() {
		node.Accept(New())
	})
}

func TestEvalIfNodeWithReturn(t *testing.T) {
	node := ast.IfNode{
		Condition: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
		Body: []ast.Node{
			ast.ReturnNode{
				Expr: ast.ValueNode{
					Value: "foo",
					Type:  token.STRING,
				},
			},
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "foo", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
}

func TestEvalReturnNode(t *testing.T) {
	node := ast.ReturnNode{
		Expr: ast.ValueNode{
			Value: "foo",
			Type:  token.STRING,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "foo", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
	assert.True(t, r.Return)
}

func TestEvalWhileNode(t *testing.T) {
	node := ast.WhileNode{
		Condition: ast.BinaryOpNode{
			Op:   token.LESS,
			Left: ast.VariableNode{Name: "foo"},
			Right: ast.ValueNode{
				Value: "3",
				Type:  token.NUMBER,
			},
		},
		Body: []ast.Node{
			ast.AssignNode{
				Name: "foo",
				Expr: ast.BinaryOpNode{
					Op:   token.ADD,
					Left: ast.VariableNode{Name: "foo"},
					Right: ast.ValueNode{
						Value: "1",
						Type:  token.NUMBER,
					},
				},
			},
		},
	}

	r := New()
	r.varSet("foo", 0.0, symtable.NUMBER)
	node.Accept(r)

	value, dtype, _ := r.varGet("foo")
	assert.Equal(t, 3.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
}

func TestEvalWhileNodeWithReturn(t *testing.T) {
	node := ast.WhileNode{
		Condition: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
		Body: []ast.Node{
			ast.ReturnNode{
				Expr: ast.ValueNode{
					Value: "foo",
					Type:  token.STRING,
				},
			},
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "foo", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
}

func TestEvalWhileNodeNotBoolCondition(t *testing.T) {
	node := ast.WhileNode{
		Condition: ast.ValueNode{
			Value: "foo",
			Type:  token.STRING,
		},
	}

	assert.Panics(t, func() {
		node.Accept(New())
	})
}

func TestEvalVariableNode(t *testing.T) {
	node := ast.VariableNode{Name: "foo"}

	r := New()
	r.varSet("foo", "bar", symtable.STRING)
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "bar", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
}

func TestEvalVariableNodeExprNotSet(t *testing.T) {
	node := ast.VariableNode{Name: "foo"}

	assert.Panics(t, func() {
		node.Accept(New())
	})
}

func TestEvalValueNode(t *testing.T) {
	node := ast.ValueNode{Value: "foo", Type: token.STRING}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "foo", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
}

func TestEvalUnaryOpNodeAdd(t *testing.T) {
	node := ast.UnaryOpNode{
		Op: token.ADD,
		Right: ast.ValueNode{
			Value: "-42",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, 42.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalUnaryOpNodeSubtract(t *testing.T) {
	node := ast.UnaryOpNode{
		Op: token.SUBTRACT,
		Right: ast.ValueNode{
			Value: "42",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, -42.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalUnaryOpNodeNot(t *testing.T) {
	node := ast.UnaryOpNode{
		Op: token.NOT,
		Right: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalBinaryOpNodeNumberAdd(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.ADD,
		Left: ast.ValueNode{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueNode{
			Value: "7",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, 18.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalBinaryOpNodeNumberSubtract(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.SUBTRACT,
		Left: ast.ValueNode{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueNode{
			Value: "7",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, 4.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalBinaryOpNodeNumberMultiply(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.MULTIPLY,
		Left: ast.ValueNode{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueNode{
			Value: "7",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, 77.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalBinaryOpNodeNumberDivide(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.DIVIDE,
		Left: ast.ValueNode{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueNode{
			Value: "7",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, 11/7.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalUnaryOpNodeModulo(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.MODULO,
		Left: ast.ValueNode{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueNode{
			Value: "7",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, 4.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalBinaryOpNodeNumberEqual(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.EQUAL,
		Left: ast.ValueNode{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueNode{
			Value: "7",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalBinaryOpNodeNumberLess(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.LESS,
		Left: ast.ValueNode{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueNode{
			Value: "7",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalBinaryOpNodeNumberLessEq(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.LESS_EQ,
		Left: ast.ValueNode{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueNode{
			Value: "7",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalBinaryOpNodeNumberGreater(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.GREATER,
		Left: ast.ValueNode{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueNode{
			Value: "7",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalBinaryOpNodeNumberGreaterEq(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.GREATER_EQ,
		Left: ast.ValueNode{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueNode{
			Value: "7",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalBinaryOpNodeStringAdd(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.ADD,
		Left: ast.ValueNode{
			Value: "foo",
			Type:  token.STRING,
		},
		Right: ast.ValueNode{
			Value: "bar",
			Type:  token.STRING,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "foobar", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
}

func TestEvalBinaryOpNodeStringEqual(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.EQUAL,
		Left: ast.ValueNode{
			Value: "foo",
			Type:  token.STRING,
		},
		Right: ast.ValueNode{
			Value: "bar",
			Type:  token.STRING,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalBinaryOpNodeBoolAnd(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.AND,
		Left: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
		Right: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalBinaryOpNodeBoolAndShortCircuit(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.AND,
		Left: ast.ValueNode{
			Value: "false",
			Type:  token.BOOL,
		},
		Right: ast.ValueNode{
			Value: "0",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalBinaryOpNodeBoolOr(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.OR,
		Left: ast.ValueNode{
			Value: "false",
			Type:  token.BOOL,
		},
		Right: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalBinaryOpNodeBoolOrShortCircuit(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.OR,
		Left: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
		Right: ast.ValueNode{
			Value: "0",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalBinaryOpNodeTypeMismatch(t *testing.T) {
	node := ast.BinaryOpNode{
		Op: token.AND,
		Left: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
		Right: ast.ValueNode{
			Value: "1",
			Type:  token.NUMBER,
		},
	}

	assert.Panics(t, func() {
		node.Accept(New())
	})
}

func TestEvalCastNodeStringToString(t *testing.T) {
	node := ast.CastNode{
		Cast: "string",
		Expr: ast.ValueNode{
			Value: "123",
			Type:  token.STRING,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "123", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
}

func TestEvalCastNodeStringToNumber(t *testing.T) {
	node := ast.CastNode{
		Cast: "number",
		Expr: ast.ValueNode{
			Value: "123",
			Type:  token.STRING,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, 123.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalCastNodeStringToNumberBadString(t *testing.T) {
	node := ast.CastNode{
		Cast: "number",
		Expr: ast.ValueNode{
			Value: "",
			Type:  token.STRING,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, 0.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalCastNodeStringToBool(t *testing.T) {
	node := ast.CastNode{
		Cast: "bool",
		Expr: ast.ValueNode{
			Value: "123",
			Type:  token.STRING,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalCastNodeNumberToString(t *testing.T) {
	node := ast.CastNode{
		Cast: "string",
		Expr: ast.ValueNode{
			Value: "123",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "123", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
}

func TestEvalCastNodeNumberToNumber(t *testing.T) {
	node := ast.CastNode{
		Cast: "number",
		Expr: ast.ValueNode{
			Value: "123",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, 123.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalCastNodeNumberToBool(t *testing.T) {
	node := ast.CastNode{
		Cast: "bool",
		Expr: ast.ValueNode{
			Value: "123",
			Type:  token.NUMBER,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalCastNodeBoolToString(t *testing.T) {
	node := ast.CastNode{
		Cast: "string",
		Expr: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "true", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
}

func TestEvalCastNodeBoolToNumberTrue(t *testing.T) {
	node := ast.CastNode{
		Cast: "number",
		Expr: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, 1.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalCastNodeBoolToNumberFalse(t *testing.T) {
	node := ast.CastNode{
		Cast: "number",
		Expr: ast.ValueNode{
			Value: "false",
			Type:  token.BOOL,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, 0.0, actual.Value)
	assert.Equal(t, symtable.NUMBER, actual.Type)
}

func TestEvalCastNodeBoolToBool(t *testing.T) {
	node := ast.CastNode{
		Cast: "bool",
		Expr: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
	}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalCastNodeBadCast(t *testing.T) {
	node := ast.CastNode{
		Cast: "foo",
		Expr: ast.ValueNode{
			Value: "true",
			Type:  token.BOOL,
		},
	}

	assert.Panics(t, func() {
		node.Accept(New())
	})
}
