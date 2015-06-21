package runtime

import (
	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/symtable"
	"github.com/tboronczyk/kiwi/token"
	"testing"
)

func TestEvalAssignStmt(t *testing.T) {
	node := ast.AssignStmt{
		Name: "foo",
		Expr: ast.ValueExpr{
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

func TestEvalFuncDef(t *testing.T) {
	node := ast.FuncDef{
		Name: "foo",
	}

	r := New()
	node.Accept(r)

	value, dtype, _ := r.funcGet("foo")
	assert.Equal(t, node, value)
	assert.Equal(t, symtable.USRFUNC, dtype)
}

func TestEvalFuncCallUserDefined(t *testing.T) {
	r := New()
	r.funcSet(
		"foo",
		ast.FuncDef{
			Args: []string{},
			Body: []ast.Node{
				ast.AssignStmt{
					Name: "bar",
					Expr: ast.ValueExpr{
						Value: "baz",
						Type:  token.STRING,
					},
				},
			},
		},
		symtable.USRFUNC)

	node := ast.FuncCall{
		Name: "foo",
		Args: []ast.Node{},
	}

	node.Accept(r)

	assert.Panics(t, func() {
		r.popStack()
	})
}

func TestEvalFuncCallBuiltin(t *testing.T) {
	r := New()
	r.funcSet(
		"foo",
		ast.FuncDef{Name: "foo", Args: []string{}},
		symtable.BUILTIN,
	)
	builtinFuncs["foo"] = func(r *Runtime) {
		r.pushStack(true, symtable.BOOL)
	}

	node := ast.FuncCall{
		Name: "foo",
		Args: []ast.Node{},
	}

	node.Accept(r)

	actual := r.popStack()
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, symtable.BOOL, actual.Type)
}

func TestEvalFuncCallUserDefinedWithArgs(t *testing.T) {
	r := New()
	r.funcSet(
		"foo",
		ast.FuncDef{
			Args: []string{"bar"},
			Body: []ast.Node{
				ast.AssignStmt{
					Name: "baz",
					Expr: ast.VariableExpr{Name: "bar"},
				},
			},
		},
		symtable.USRFUNC)

	node := ast.FuncCall{
		Name: "foo",
		Args: []ast.Node{
			ast.ValueExpr{
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

func TestEvalFuncCallUserDefinedWithReturn(t *testing.T) {
	r := New()
	r.funcSet(
		"foo",
		ast.FuncDef{
			Args: []string{"bar"},
			Body: []ast.Node{
				ast.ReturnStmt{
					Expr: ast.VariableExpr{Name: "bar"},
				},
			},
		},
		symtable.USRFUNC)

	node := ast.FuncCall{
		Name: "foo",
		Args: []ast.Node{
			ast.ValueExpr{
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

func TestEvalFuncCallNotExist(t *testing.T) {
	node := ast.FuncCall{Name: "foo"}

	r := New()

	assert.Panics(t, func() {
		node.Accept(r)
	})
}

func TestEvalFuncCallArityMismatch(t *testing.T) {
	r := New()
	r.funcSet("foo", ast.FuncDef{}, symtable.USRFUNC)

	node := ast.FuncCall{
		Name: "foo",
		Args: []ast.Node{
			ast.ValueExpr{
				Value: "bar",
				Type:  token.STRING,
			},
		},
	}

	assert.Panics(t, func() {
		node.Accept(r)
	})
}

func TestEvalIfStmtTrue(t *testing.T) {
	node := ast.IfStmt{
		Condition: ast.ValueExpr{
			Value: "true",
			Type:  token.BOOL,
		},
		Body: []ast.Node{
			ast.AssignStmt{
				Name: "foo",
				Expr: ast.ValueExpr{
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

func TestEvalIfStmtFalse(t *testing.T) {
	node := ast.IfStmt{
		Condition: ast.ValueExpr{
			Value: "false",
			Type:  token.BOOL,
		},
		Else: ast.IfStmt{
			Condition: ast.ValueExpr{
				Value: "true",
				Type:  token.BOOL,
			},
			Body: []ast.Node{
				ast.AssignStmt{
					Name: "foo",
					Expr: ast.ValueExpr{
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

func TestEvalIfStmtNotBoolCondition(t *testing.T) {
	node := ast.IfStmt{
		Condition: ast.ValueExpr{
			Value: "foo",
			Type:  token.STRING,
		},
	}

	assert.Panics(t, func() {
		node.Accept(New())
	})
}

func TestEvalIfStmtWithReturn(t *testing.T) {
	node := ast.IfStmt{
		Condition: ast.ValueExpr{
			Value: "true",
			Type:  token.BOOL,
		},
		Body: []ast.Node{
			ast.ReturnStmt{
				Expr: ast.ValueExpr{
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

func TestEvalReturnStmt(t *testing.T) {
	node := ast.ReturnStmt{
		Expr: ast.ValueExpr{
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

func TestEvalWhileStmt(t *testing.T) {
	node := ast.WhileStmt{
		Condition: ast.BinaryExpr{
			Op:   token.LESS,
			Left: ast.VariableExpr{Name: "foo"},
			Right: ast.ValueExpr{
				Value: "3",
				Type:  token.NUMBER,
			},
		},
		Body: []ast.Node{
			ast.AssignStmt{
				Name: "foo",
				Expr: ast.BinaryExpr{
					Op:   token.ADD,
					Left: ast.VariableExpr{Name: "foo"},
					Right: ast.ValueExpr{
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

func TestEvalWhileStmtWithReturn(t *testing.T) {
	node := ast.WhileStmt{
		Condition: ast.ValueExpr{
			Value: "true",
			Type:  token.BOOL,
		},
		Body: []ast.Node{
			ast.ReturnStmt{
				Expr: ast.ValueExpr{
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

func TestEvalWhileStmtNotBoolCondition(t *testing.T) {
	node := ast.WhileStmt{
		Condition: ast.ValueExpr{
			Value: "foo",
			Type:  token.STRING,
		},
	}

	assert.Panics(t, func() {
		node.Accept(New())
	})
}

func TestEvalVariableExpr(t *testing.T) {
	node := ast.VariableExpr{Name: "foo"}

	r := New()
	r.varSet("foo", "bar", symtable.STRING)
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "bar", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
}

func TestEvalVariableExprNotSet(t *testing.T) {
	node := ast.VariableExpr{Name: "foo"}

	assert.Panics(t, func() {
		node.Accept(New())
	})
}

func TestEvalValueExpr(t *testing.T) {
	node := ast.ValueExpr{Value: "foo", Type: token.STRING}

	r := New()
	node.Accept(r)

	actual := r.popStack()
	assert.Equal(t, "foo", actual.Value)
	assert.Equal(t, symtable.STRING, actual.Type)
}

func TestEvalUnaryExprAdd(t *testing.T) {
	node := ast.UnaryExpr{
		Op: token.ADD,
		Right: ast.ValueExpr{
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

func TestEvalUnaryExprSubtract(t *testing.T) {
	node := ast.UnaryExpr{
		Op: token.SUBTRACT,
		Right: ast.ValueExpr{
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

func TestEvalUnaryExprNot(t *testing.T) {
	node := ast.UnaryExpr{
		Op: token.NOT,
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprNumberAdd(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.ADD,
		Left: ast.ValueExpr{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprNumberSubtract(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.SUBTRACT,
		Left: ast.ValueExpr{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprNumberMultiply(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.MULTIPLY,
		Left: ast.ValueExpr{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprNumberDivide(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.DIVIDE,
		Left: ast.ValueExpr{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueExpr{
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

func TestEvalUnaryExprModulo(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.MODULO,
		Left: ast.ValueExpr{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprNumberEqual(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.EQUAL,
		Left: ast.ValueExpr{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprNumberLess(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.LESS,
		Left: ast.ValueExpr{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprNumberLessEq(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.LESS_EQ,
		Left: ast.ValueExpr{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprNumberGreater(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.GREATER,
		Left: ast.ValueExpr{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprNumberGreaterEq(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.GREATER_EQ,
		Left: ast.ValueExpr{
			Value: "11",
			Type:  token.NUMBER,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprStringAdd(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.ADD,
		Left: ast.ValueExpr{
			Value: "foo",
			Type:  token.STRING,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprStringEqual(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.EQUAL,
		Left: ast.ValueExpr{
			Value: "foo",
			Type:  token.STRING,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprBoolAnd(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.AND,
		Left: ast.ValueExpr{
			Value: "true",
			Type:  token.BOOL,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprBoolAndShortCircuit(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.AND,
		Left: ast.ValueExpr{
			Value: "false",
			Type:  token.BOOL,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprBoolOr(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.OR,
		Left: ast.ValueExpr{
			Value: "false",
			Type:  token.BOOL,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprBoolOrShortCircuit(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.OR,
		Left: ast.ValueExpr{
			Value: "true",
			Type:  token.BOOL,
		},
		Right: ast.ValueExpr{
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

func TestEvalBinaryExprTypeMismatch(t *testing.T) {
	node := ast.BinaryExpr{
		Op: token.AND,
		Left: ast.ValueExpr{
			Value: "true",
			Type:  token.BOOL,
		},
		Right: ast.ValueExpr{
			Value: "1",
			Type:  token.NUMBER,
		},
	}

	assert.Panics(t, func() {
		node.Accept(New())
	})
}

func TestEvalCastExprStringToString(t *testing.T) {
	node := ast.CastExpr{
		Cast: "string",
		Expr: ast.ValueExpr{
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

func TestEvalCastExprStringToNumber(t *testing.T) {
	node := ast.CastExpr{
		Cast: "number",
		Expr: ast.ValueExpr{
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

func TestEvalCastExprStringToNumberBadString(t *testing.T) {
	node := ast.CastExpr{
		Cast: "number",
		Expr: ast.ValueExpr{
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

func TestEvalCastExprStringToBool(t *testing.T) {
	node := ast.CastExpr{
		Cast: "bool",
		Expr: ast.ValueExpr{
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

func TestEvalCastExprNumberToString(t *testing.T) {
	node := ast.CastExpr{
		Cast: "string",
		Expr: ast.ValueExpr{
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

func TestEvalCastExprNumberToNumber(t *testing.T) {
	node := ast.CastExpr{
		Cast: "number",
		Expr: ast.ValueExpr{
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

func TestEvalCastExprNumberToBool(t *testing.T) {
	node := ast.CastExpr{
		Cast: "bool",
		Expr: ast.ValueExpr{
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

func TestEvalCastExprBoolToString(t *testing.T) {
	node := ast.CastExpr{
		Cast: "string",
		Expr: ast.ValueExpr{
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

func TestEvalCastExprBoolToNumberTrue(t *testing.T) {
	node := ast.CastExpr{
		Cast: "number",
		Expr: ast.ValueExpr{
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

func TestEvalCastExprBoolToNumberFalse(t *testing.T) {
	node := ast.CastExpr{
		Cast: "number",
		Expr: ast.ValueExpr{
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

func TestEvalCastExprBoolToBool(t *testing.T) {
	node := ast.CastExpr{
		Cast: "bool",
		Expr: ast.ValueExpr{
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

func TestEvalCastExprBadCast(t *testing.T) {
	node := ast.CastExpr{
		Cast: "foo",
		Expr: ast.ValueExpr{
			Value: "true",
			Type:  token.BOOL,
		},
	}

	assert.Panics(t, func() {
		node.Accept(New())
	})
}
