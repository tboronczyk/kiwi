// Package token defines the token types used to represent Kiwi lexemes.
package token

import (
	"strconv"
)

// Token is a member in the lexical types set.
type Token uint

const (
	UNKNOWN Token = iota
	EOF
	NEWLINE

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
	ELSE
	LPAREN
	RPAREN
)

var tokens = []string{
	UNKNOWN:    "UNKNOWN",
	EOF:        "EOF",
	NEWLINE:    "NEWLINE",
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

// Precedence returns the relative precedence of an operator. A higher value
// is a higher precedence.
func Precedence(t Token) int {
	if t.IsLogOp() {
		return 1
	} else if t.IsCmpOp() {
		return 2
	} else if t.IsAddOp() {
		return 3
	} else if t.IsMulOp() {
		return 4
	} else {
		return 0
	}
}

// IsAddOp returns bool to indicate whether the token represents an
// addition-level operator.
func (t Token) IsAddOp() bool {
	return t > addop_start && t < addop_end
}

// IsMulOp returns bool to indicate whether the token represents a
// multiplication-level operator.
func (t Token) IsMulOp() bool {
	return t > mulop_start && t < mulop_end
}

// IsCmpOp returns bool to indicate whether the token represents a comparision
// operator.
func (t Token) IsCmpOp() bool {
	return t > cmpop_start && t < cmpop_end
}

// IsLogOp returns bool to indicate whether the token represents a logic
// operator.
func (t Token) IsLogOp() bool {
	return t > logop_start && t < logop_end
}

// IsBinOp returns bool to indicate whether the token represents a left-binding
// binary operator.
func (t Token) IsBinOp() bool {
	return (t.IsAddOp() || t.IsMulOp() || t.IsCmpOp() || t.IsLogOp()) &&
		t != NOT
}

// IsUnaryOp returns bool to indicate whether the token represents a
// right-binding operator.
func (t Token) IsUnaryOp() bool {
	return t.IsAddOp() || t == NOT
}

// IsStmtKeyword returns bool to indicate whether the token represents a keyword
// that may begin a statement.
func (t Token) IsStmtKeyword() bool {
	return t > stmtkwd_start && t < stmtkwd_end
}

// IsLiteral returns bool to indicate whether the token represents a literal
// value.
func (t Token) IsLiteral() bool {
	return t > lit_start && t < lit_end
}
