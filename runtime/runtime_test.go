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
	varTable := symtable.New()
	Eval(node, varTable, nil)

	value, dtype, ok := varTable.Get("foo")
	assert.Equal(t, "bar", value)
	assert.Equal(t, symtable.STRING, dtype)
	assert.True(t, ok)
}

func TestEvalFuncDef(t *testing.T) {
	// todo...
}

func TestEvalFuncCallUserDefined(t *testing.T) {
	// todo...
}

func TestEvalFuncCallUserDefinedWithReturn(t *testing.T) {
	// todo...
}

func TestEvalFuncCallBuiltin(t *testing.T) {
	// todo...
}

func TestEvalFuncCallArityMismatch(t *testing.T) {
	// todo...
}

func TestEvalFuncCallNotExist(t *testing.T) {
	// todo...
}

func TestEvalIfStmtTrue(t *testing.T) {
	// todo...
}

func TestEvalIfStmtTrueWithReturn(t *testing.T) {
	// todo...
}

func TestEvalIfStmtFalse(t *testing.T) {
	// todo...
}

func TestEvalIfStmtNotBoolCondition(t *testing.T) {
	// todo...
}

func TestEvalReturnStmt(t *testing.T) {
	node := ast.ReturnStmt{
		Expr: ast.ValueExpr{
			Value: "foo",
			Type:  token.STRING,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, "foo", value)
	assert.Equal(t, symtable.STRING, dtype)
	assert.True(t, isReturn)
}

func TestEvalWhileStmt(t *testing.T) {
	// todo...
}

func TestEvalWhileStmtWithReturn(t *testing.T) {
	// todo...
}

func TestEvalWhileStmtNotBoolCondition(t *testing.T) {
	// todo...
}

func TestEvalVariableExpr(t *testing.T) {
	node := ast.VariableExpr{Name: "foo"}
	varTable := symtable.New()
	varTable.Set("foo", "bar", symtable.STRING)
	value, dtype, isReturn := Eval(node, varTable, nil)
	assert.Equal(t, "bar", value)
	assert.Equal(t, symtable.STRING, dtype)
	assert.False(t, isReturn)
}

func TestEvalVariableExprNotSet(t *testing.T) {
	node := ast.VariableExpr{Name: "foo"}
	varTable := symtable.New()
	assert.Panics(t, func() {
		Eval(node, varTable, nil)
	})
}

func TestEvalValueExpr(t *testing.T) {
	node := ast.ValueExpr{Value: "foo", Type: token.STRING}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, "foo", value)
	assert.Equal(t, symtable.STRING, dtype)
	assert.False(t, isReturn)
}

func TestEvalUnaryExprAdd(t *testing.T) {
	node := ast.UnaryExpr{
		Op: token.ADD,
		Right: ast.ValueExpr{
			Value: "-42",
			Type:  token.NUMBER,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, 42.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
}

func TestEvalUnaryExprSubtract(t *testing.T) {
	node := ast.UnaryExpr{
		Op: token.SUBTRACT,
		Right: ast.ValueExpr{
			Value: "42",
			Type:  token.NUMBER,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, -42.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
}

func TestEvalUnaryExprNot(t *testing.T) {
	node := ast.UnaryExpr{
		Op: token.NOT,
		Right: ast.ValueExpr{
			Value: "true",
			Type:  token.BOOL,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.False(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, 18.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, 4.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, 77.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, 11/7.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, 4.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.False(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.False(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.False(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.True(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.True(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, "foobar", value)
	assert.Equal(t, symtable.STRING, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.False(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.True(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.False(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.True(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.True(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
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
		Eval(node, nil, nil)
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
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, "123", value)
	assert.Equal(t, symtable.STRING, dtype)
	assert.False(t, isReturn)
}

func TestEvalCastExprStringToNumber(t *testing.T) {
	node := ast.CastExpr{
		Cast: "number",
		Expr: ast.ValueExpr{
			Value: "123",
			Type:  token.STRING,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, 123.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
}

func TestEvalCastExprStringToNumberBadString(t *testing.T) {
	node := ast.CastExpr{
		Cast: "number",
		Expr: ast.ValueExpr{
			Value: "",
			Type:  token.STRING,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, 0.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
}

func TestEvalCastExprStringToBool(t *testing.T) {
	node := ast.CastExpr{
		Cast: "bool",
		Expr: ast.ValueExpr{
			Value: "123",
			Type:  token.STRING,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.True(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
}

func TestEvalCastExprNumberToString(t *testing.T) {
	node := ast.CastExpr{
		Cast: "string",
		Expr: ast.ValueExpr{
			Value: "123",
			Type:  token.NUMBER,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, "123", value)
	assert.Equal(t, symtable.STRING, dtype)
	assert.False(t, isReturn)
}

func TestEvalCastExprNumberToNumber(t *testing.T) {
	node := ast.CastExpr{
		Cast: "number",
		Expr: ast.ValueExpr{
			Value: "123",
			Type:  token.NUMBER,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, 123.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
}

func TestEvalCastExprNumberToBool(t *testing.T) {
	node := ast.CastExpr{
		Cast: "bool",
		Expr: ast.ValueExpr{
			Value: "123",
			Type:  token.NUMBER,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.True(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
}

func TestEvalCastExprBoolToString(t *testing.T) {
	node := ast.CastExpr{
		Cast: "string",
		Expr: ast.ValueExpr{
			Value: "true",
			Type:  token.BOOL,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, "true", value)
	assert.Equal(t, symtable.STRING, dtype)
	assert.False(t, isReturn)
}

func TestEvalCastExprBoolToNumberTrue(t *testing.T) {
	node := ast.CastExpr{
		Cast: "number",
		Expr: ast.ValueExpr{
			Value: "true",
			Type:  token.BOOL,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, 1.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
}

func TestEvalCastExprBoolToNumberFalse(t *testing.T) {
	node := ast.CastExpr{
		Cast: "number",
		Expr: ast.ValueExpr{
			Value: "false",
			Type:  token.BOOL,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.Equal(t, 0.0, value)
	assert.Equal(t, symtable.NUMBER, dtype)
	assert.False(t, isReturn)
}

func TestEvalCastExprBoolToBool(t *testing.T) {
	node := ast.CastExpr{
		Cast: "bool",
		Expr: ast.ValueExpr{
			Value: "true",
			Type:  token.BOOL,
		},
	}
	value, dtype, isReturn := Eval(node, nil, nil)
	assert.True(t, value.(bool))
	assert.Equal(t, symtable.BOOL, dtype)
	assert.False(t, isReturn)
}

func TestEvalCastExprBadCast(t *testing.T) {
	node := ast.CastExpr{
		Cast: "foo",
		Expr: ast.ValueExpr{},
	}
	assert.Panics(t, func() {
		Eval(node, nil, nil)
	})
}
