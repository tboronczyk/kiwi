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
	MALFORMED:  "MALFORMED",
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

// Precedence returns whether the operator represented by t1 has a higher
// precedence than the one represented by t2. If either t1 or t2 is not an
// operator token then err is true.
func Precedence(t1, t2 Token) (bool, err bool) {
	p1 := precedence(t1)
	p2 := precedence(t2)
	if p1 == 0 || p2 == 0 {
		return false, true
	}
	return p1 > p2, false
}

// precedence returns an operator's precedence. A higher value is a higher
// precedence.
func precedence(t Token) uint8 {
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
// that may form an expression (left-binding binary operators).
func (t Token) IsExprOp() bool {
	return (t.IsAddOp() || t.IsMulOp() || t.IsCmpOp() || t.IsLogOp()) &&
		t != NOT
}

// IsTermOp returns bool indicating whether the token represents an operator
// that may lead a term (right-binding unary operators).
func (t Token) IsTermOp() bool {
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
