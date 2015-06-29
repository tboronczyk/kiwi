// Package token defines the token types used to represent Kiwi lexemes.
package token

import (
	"strconv"
)

// Token is a member in the lexical types set.
type Token uint8

const (
	UNKNOWN Token = iota
	MALFORMED
	EOF

	addop_start
	// addition-level operators
	ADD
	SUBTRACT
	addop_end

	mulop_start
	// multiplication-level operators
	MULTIPLY
	DIVIDE
	MODULO
	mulop_end

	cmpop_start
	// comparision operators
	EQUAL
	NOT_EQUAL
	GREATER
	GREATER_EQ
	LESS
	LESS_EQ
	cmpop_end

	logop_start
	// logic operators
	AND
	OR
	NOT
	logop_end

	stmtkwd_start
	// statement keywords
	IF
	FUNC
	RETURN
	WHILE
	stmtkwd_end

	lit_start
	// literal values
	BOOL
	IDENTIFIER
	NUMBER
	STRING
	lit_end

	ASSIGN
	LBRACE
	RBRACE
	COLON
	COMMA
	COMMENT
	DOT
	ELSE
	LPAREN
	RPAREN
)

var tokens = []string{
	UNKNOWN:    "UNKNOWN",
	MALFORMED:  "MALFORMED",
	EOF:        "EOF",
	ADD:        "+",
	SUBTRACT:   "-",
	MULTIPLY:   "*",
	DIVIDE:     "/",
	MODULO:     "%",
	EQUAL:      "=",
	NOT_EQUAL:  "~=",
	GREATER:    ">",
	GREATER_EQ: ">=",
	LESS:       "<",
	LESS_EQ:    "<=",
	AND:        "&&",
	OR:         "||",
	NOT:        "~",
	IF:         "if",
	FUNC:       "func",
	RETURN:     "return",
	WHILE:      "while",
	BOOL:       "BOOL",
	IDENTIFIER: "IDENTIFIER",
	NUMBER:     "NUMBER",
	STRING:     "STRING",
	ASSIGN:     ":=",
	LBRACE:     "{",
	RBRACE:     "}",
	COLON:      ":",
	COMMA:      ",",
	COMMENT:    "COMMENT",
	DOT:        ".",
	ELSE:       "else",
	LPAREN:     "(",
	RPAREN:     ")",
}

// String returns the string representation of a token.
func (t Token) String() string {
	str := ""
	if t >= 0 && t < Token(len(tokens)) {
		str = tokens[t]
	}
	if str == "" {
		str = "Token(" + strconv.Itoa(int(t)) + ")"
	}
	return str
}

// IsAddOp returns bool indicating whether the token represents an
// addition-level operator.
func (t Token) IsAddOp() bool {
	return t > addop_start && t < addop_end
}

// IsMulOp returns bool indicating whether the token represents a
// multiplication-level operator.
func (t Token) IsMulOp() bool {
	return t > mulop_start && t < mulop_end
}

// IsCmpOp returns bool indicating whether the token represents a comparision
// operator.
func (t Token) IsCmpOp() bool {
	return t > cmpop_start && t < cmpop_end
}

// IsLogOp returns bool indicating whether the token represents a logic
// operator.
func (t Token) IsLogOp() bool {
	return t > logop_start && t < logop_end
}

// IsExprOp returns bool indicating whether the token represents an operator
// that may lead an expression (right-binding unary operators).
func (t Token) IsExprOp() bool {
	return t.IsAddOp() || t == NOT
}

// IsStmtKeyword returns bool indicating whether the token represents a keyword
// that may begin a statement.
func (t Token) IsStmtKeyword() bool {
	return t > stmtkwd_start && t < stmtkwd_end
}

// IsLiteral returns bool indicating whether the token represents a literal
// value.
func (t Token) IsLiteral() bool {
	return t > lit_start && t < lit_end
}
