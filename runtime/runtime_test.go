package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tboronczyk/kiwi/analyzer"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/symtable"
	"github.com/tboronczyk/kiwi/token"
)

func TestEvalValueNodeNumber(t *testing.T) {
	n := &ast.ValueNode{Value: "42", Type: token.NUMBER}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, 42.0, actual.Value)
	assert.Equal(t, analyzer.NUMBER, actual.Type)
}

func TestEvalValueNodeString(t *testing.T) {
	n := &ast.ValueNode{Value: "foo", Type: token.STRING}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, "foo", actual.Value)
	assert.Equal(t, analyzer.STRING, actual.Type)
}

func TestEvalValueNodeBool(t *testing.T) {
	n := &ast.ValueNode{Value: "true", Type: token.BOOL}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalAssignNode(t *testing.T) {
	n := &ast.AssignNode{
		Name: "foo",
		Expr: &ast.ValueNode{Value: "bar", Type: token.STRING},
	}
	n.Accept(analyzer.New())
	n.Accept(New())

	e, _ := n.SymTable.Get("foo", symtable.VARIABLE)
	assert.Equal(t, "bar", e.(stackEntry).Value)
	assert.Equal(t, analyzer.STRING, e.(stackEntry).Type)
}

func TestEvalVariableNode(t *testing.T) {
	n := &ast.VariableNode{Name: "foo"}
	n.SymTable = symtable.New()
	n.SymTable.Set("foo", symtable.VARIABLE,
		stackEntry{Value: 42, Type: analyzer.NUMBER},
	)
	n.Accept(New())

	e, _ := n.SymTable.Get("foo", symtable.VARIABLE)
	assert.Equal(t, 42, e.(stackEntry).Value)
	assert.Equal(t, analyzer.NUMBER, e.(stackEntry).Type)
}

func TestEvalUnaryOpNodeNot(t *testing.T) {
	n := &ast.UnaryOpNode{
		Op:   token.NOT,
		Expr: &ast.ValueNode{Value: "false", Type: token.BOOL},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalUnaryOpNodePositive(t *testing.T) {
	n := &ast.UnaryOpNode{
		Op:   token.ADD,
		Expr: &ast.ValueNode{Value: "-42", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, 42.0, actual.Value)
	assert.Equal(t, analyzer.NUMBER, actual.Type)
}

func TestEvalUnaryOpNodeNegative(t *testing.T) {
	n := &ast.UnaryOpNode{
		Op:   token.SUBTRACT,
		Expr: &ast.ValueNode{Value: "42", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, -42.0, actual.Value)
	assert.Equal(t, analyzer.NUMBER, actual.Type)
}

func TestEvalCastNodeStringToString(t *testing.T) {
	n := &ast.CastNode{
		Cast: "string",
		Expr: &ast.ValueNode{Value: "foo", Type: token.STRING},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, "foo", actual.Value)
	assert.Equal(t, analyzer.STRING, actual.Type)
}

func TestEvalCastNodeNumberToString(t *testing.T) {
	n := &ast.CastNode{
		Cast: "string",
		Expr: &ast.ValueNode{Value: "42", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, "42", actual.Value)
	assert.Equal(t, analyzer.STRING, actual.Type)
}

func TestEvalCastNodeBoolToString(t *testing.T) {
	n := &ast.CastNode{
		Cast: "string",
		Expr: &ast.ValueNode{Value: "TRUE", Type: token.BOOL},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, "true", actual.Value)
	assert.Equal(t, analyzer.STRING, actual.Type)
}

func TestEvalCastNodeStringToNumber(t *testing.T) {
	n := &ast.CastNode{
		Cast: "number",
		Expr: &ast.ValueNode{Value: "42", Type: token.STRING},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, 42.0, actual.Value)
	assert.Equal(t, analyzer.NUMBER, actual.Type)
}

func TestEvalCastNodeNumberToNumber(t *testing.T) {
	n := &ast.CastNode{
		Cast: "number",
		Expr: &ast.ValueNode{Value: "42", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, 42.0, actual.Value)
	assert.Equal(t, analyzer.NUMBER, actual.Type)
}

func TestEvalCastNodeBoolToNumber(t *testing.T) {
	n := &ast.CastNode{
		Cast: "number",
		Expr: &ast.ValueNode{Value: "TRUE", Type: token.BOOL},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, 1.0, actual.Value)
	assert.Equal(t, analyzer.NUMBER, actual.Type)
}

func TestEvalCastNodeStringToBoolTrue(t *testing.T) {
	n := &ast.CastNode{
		Cast: "bool",
		Expr: &ast.ValueNode{Value: "true", Type: token.STRING},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalCastNodeStringToBoolFalse(t *testing.T) {
	n := &ast.CastNode{
		Cast: "bool",
		Expr: &ast.ValueNode{Value: "false", Type: token.STRING},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalCastNodeNumberToBoolTrue(t *testing.T) {
	n := &ast.CastNode{
		Cast: "bool",
		Expr: &ast.ValueNode{Value: "42", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalCastNodeNumberToBoolFalse(t *testing.T) {
	n := &ast.CastNode{
		Cast: "bool",
		Expr: &ast.ValueNode{Value: "0", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalCastNodeBoolToBool(t *testing.T) {
	n := &ast.CastNode{
		Cast: "bool",
		Expr: &ast.ValueNode{Value: "TRUE", Type: token.BOOL},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalBinaryOpNumberAdd(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.ADD,
		Left:  &ast.ValueNode{Value: "11", Type: token.NUMBER},
		Right: &ast.ValueNode{Value: "7", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, 18.0, actual.Value.(float64))
	assert.Equal(t, analyzer.NUMBER, actual.Type)
}

func TestEvalBinaryOpNumberSubtract(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.SUBTRACT,
		Left:  &ast.ValueNode{Value: "11", Type: token.NUMBER},
		Right: &ast.ValueNode{Value: "7", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, 4.0, actual.Value.(float64))
	assert.Equal(t, analyzer.NUMBER, actual.Type)
}

func TestEvalBinaryOpNumberMultiply(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.MULTIPLY,
		Left:  &ast.ValueNode{Value: "11", Type: token.NUMBER},
		Right: &ast.ValueNode{Value: "7", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, 77.0, actual.Value.(float64))
	assert.Equal(t, analyzer.NUMBER, actual.Type)
}

func TestEvalBinaryOpNumberDivide(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.DIVIDE,
		Left:  &ast.ValueNode{Value: "11", Type: token.NUMBER},
		Right: &ast.ValueNode{Value: "7", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, 11/7.0, actual.Value.(float64))
	assert.Equal(t, analyzer.NUMBER, actual.Type)
}

func TestEvalBinaryOpNumberModulo(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.MODULO,
		Left:  &ast.ValueNode{Value: "11", Type: token.NUMBER},
		Right: &ast.ValueNode{Value: "7", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, 4.0, actual.Value.(float64))
	assert.Equal(t, analyzer.NUMBER, actual.Type)
}

func TestEvalBinaryOpNumberEqual(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.EQUAL,
		Left:  &ast.ValueNode{Value: "11", Type: token.NUMBER},
		Right: &ast.ValueNode{Value: "7", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}
func TestEvalBinaryOpNumberNotEqual(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.NOT_EQUAL,
		Left:  &ast.ValueNode{Value: "11", Type: token.NUMBER},
		Right: &ast.ValueNode{Value: "7", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}
func TestEvalBinaryOpNumberLess(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.LESS,
		Left:  &ast.ValueNode{Value: "11", Type: token.NUMBER},
		Right: &ast.ValueNode{Value: "7", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}
func TestEvalBinaryOpNumberLessEq(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.LESS_EQ,
		Left:  &ast.ValueNode{Value: "11", Type: token.NUMBER},
		Right: &ast.ValueNode{Value: "7", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}
func TestEvalBinaryOpNumberGreater(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.GREATER,
		Left:  &ast.ValueNode{Value: "11", Type: token.NUMBER},
		Right: &ast.ValueNode{Value: "7", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalBinaryOpNumberGreaterEq(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.GREATER_EQ,
		Left:  &ast.ValueNode{Value: "11", Type: token.NUMBER},
		Right: &ast.ValueNode{Value: "7", Type: token.NUMBER},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalBinaryOpStringAdd(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.ADD,
		Left:  &ast.ValueNode{Value: "foo", Type: token.STRING},
		Right: &ast.ValueNode{Value: "bar", Type: token.STRING},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.Equal(t, "foobar", actual.Value.(string))
	assert.Equal(t, analyzer.STRING, actual.Type)
}

func TestEvalBinaryOpStringEqual(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.EQUAL,
		Left:  &ast.ValueNode{Value: "foo", Type: token.STRING},
		Right: &ast.ValueNode{Value: "bar", Type: token.STRING},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalBinaryOpStringNotEqual(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.NOT_EQUAL,
		Left:  &ast.ValueNode{Value: "foo", Type: token.STRING},
		Right: &ast.ValueNode{Value: "bar", Type: token.STRING},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalBinaryOpBoolAnd(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.AND,
		Left:  &ast.ValueNode{Value: "true", Type: token.BOOL},
		Right: &ast.ValueNode{Value: "false", Type: token.BOOL},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalBinaryOpBoolOr(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.OR,
		Left:  &ast.ValueNode{Value: "false", Type: token.BOOL},
		Right: &ast.ValueNode{Value: "true", Type: token.BOOL},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalBinaryOpBoolEqual(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.EQUAL,
		Left:  &ast.ValueNode{Value: "true", Type: token.BOOL},
		Right: &ast.ValueNode{Value: "false", Type: token.BOOL},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalBinaryOpBoolNotEqual(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.NOT_EQUAL,
		Left:  &ast.ValueNode{Value: "true", Type: token.BOOL},
		Right: &ast.ValueNode{Value: "false", Type: token.BOOL},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalBinaryOpBoolShortCircuitAnd(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.AND,
		Left:  &ast.ValueNode{Value: "false", Type: token.BOOL},
		Right: &ast.ValueNode{Value: "true", Type: token.BOOL},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.False(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalBinaryOpBoolShortCircuitOr(t *testing.T) {
	n := &ast.BinaryOpNode{
		Op:    token.OR,
		Left:  &ast.ValueNode{Value: "true", Type: token.BOOL},
		Right: &ast.ValueNode{Value: "false", Type: token.BOOL},
	}
	r := New()
	n.Accept(r)

	actual := r.stack.Pop().(stackEntry)
	assert.True(t, actual.Value.(bool))
	assert.Equal(t, analyzer.BOOL, actual.Type)
}

func TestEvalWhileNode(t *testing.T) {
	n := &ast.WhileNode{
		Condition: &ast.BinaryOpNode{
			Op:    token.LESS,
			Left:  &ast.VariableNode{Name: "foo"},
			Right: &ast.ValueNode{Value: "10", Type: token.NUMBER},
		},
		Body: []ast.Node{
			&ast.AssignNode{
				Name: "foo",
				Expr: &ast.BinaryOpNode{
					Op:   token.ADD,
					Left: &ast.VariableNode{Name: "foo"},
					Right: &ast.ValueNode{
						Value: "1",
						Type:  token.NUMBER,
					},
				},
			},
		},
	}
	s := symtable.New()
	s.Set("foo", symtable.VARIABLE,
		stackEntry{Value: 0.0, Type: analyzer.NUMBER},
	)
	n.Condition.(*ast.BinaryOpNode).Left.(*ast.VariableNode).SymTable = s
	n.Body[0].(*ast.AssignNode).SymTable = s
	n.Body[0].(*ast.AssignNode).Expr.(*ast.BinaryOpNode).Left.(*ast.VariableNode).SymTable = s
	r := New()
	n.Accept(r)

	e, _ := s.Get("foo", symtable.VARIABLE)
	assert.Equal(t, 10.0, e.(stackEntry).Value)
}

func TestEvalIfNodeTrue(t *testing.T) {
	n := &ast.IfNode{
		Condition: &ast.BinaryOpNode{
			Op:    token.LESS,
			Left:  &ast.VariableNode{Name: "foo"},
			Right: &ast.ValueNode{Value: "10", Type: token.NUMBER},
		},
		Body: []ast.Node{
			&ast.AssignNode{
				Name: "foo",
				Expr: &ast.BinaryOpNode{
					Op:   token.ADD,
					Left: &ast.VariableNode{Name: "foo"},
					Right: &ast.ValueNode{
						Value: "1",
						Type:  token.NUMBER,
					},
				},
			},
		},
	}
	s := symtable.New()
	s.Set("foo", symtable.VARIABLE,
		stackEntry{Value: 0.0, Type: analyzer.NUMBER},
	)
	n.Condition.(*ast.BinaryOpNode).Left.(*ast.VariableNode).SymTable = s
	n.Body[0].(*ast.AssignNode).SymTable = s
	n.Body[0].(*ast.AssignNode).Expr.(*ast.BinaryOpNode).Left.(*ast.VariableNode).SymTable = s
	r := New()
	n.Accept(r)

	e, _ := s.Get("foo", symtable.VARIABLE)
	assert.Equal(t, 1.0, e.(stackEntry).Value)
}

func TestEvalIfNodeFalse(t *testing.T) {
	n := &ast.IfNode{
		Condition: &ast.BinaryOpNode{
			Op:    token.LESS,
			Left:  &ast.VariableNode{Name: "foo"},
			Right: &ast.ValueNode{Value: "10", Type: token.NUMBER},
		},
		Else: &ast.IfNode{
			Condition: &ast.ValueNode{Value: "true", Type: token.BOOL},
			Body: []ast.Node{
				&ast.AssignNode{
					Name: "foo",
					Expr: &ast.BinaryOpNode{
						Op:   token.ADD,
						Left: &ast.VariableNode{Name: "foo"},
						Right: &ast.ValueNode{
							Value: "1",
							Type:  token.NUMBER,
						},
					},
				},
			},
		},
	}
	s := symtable.New()
	s.Set("foo", symtable.VARIABLE,
		stackEntry{Value: 20.0, Type: analyzer.NUMBER},
	)
	n.Condition.(*ast.BinaryOpNode).Left.(*ast.VariableNode).SymTable = s
	n.Else.(*ast.IfNode).Body[0].(*ast.AssignNode).SymTable = s
	n.Else.(*ast.IfNode).Body[0].(*ast.AssignNode).Expr.(*ast.BinaryOpNode).Left.(*ast.VariableNode).SymTable = s
	r := New()
	n.Accept(r)

	e, _ := s.Get("foo", symtable.VARIABLE)
	assert.Equal(t, 21.0, e.(stackEntry).Value)
}
