package main

import "strconv"

// Token represents the token types used to represent Kiwi lexemes.
type Token uint

const (
	TkUnknown Token = iota
	TkEof

	addop_start
	// addition-level operators
	TkAdd
	TkSubtract
	addop_end

	mulop_start
	// multiplication-level operators
	TkMultiply
	TkDivide
	TkModulo
	mulop_end

	cmpop_start
	// comparision operators
	TkEqual
	TkNotEqual
	TkGreater
	TkGreaterEq
	TkLess
	TkLessEq
	cmpop_end

	logop_start
	// logic operators
	TkAnd
	TkOr
	TkNot
	logop_end

	stmtkwd_start
	// statement keywords
	TkIf
	TkFunc
	TkReturn
	TkWhile
	stmtkwd_end

	lit_start
	// literal values
	TkBool
	TkIdentifier
	TkNumber
	TkString
	lit_end

	TkAssign
	TkLBrace
	TkRBrace
	TkColon
	TkComma
	TkComment
	TkElse
	TkLParen
	TkRParent
)

var tokens = []string{
	TkUnknown:    "Unknown",
	TkEof:        "eof",
	TkAdd:        "+",
	TkSubtract:   "-",
	TkMultiply:   "*",
	TkDivide:     "/",
	TkModulo:     "%",
	TkEqual:      "=",
	TkNotEqual:   "~=",
	TkGreater:    ">",
	TkGreaterEq:  ">=",
	TkLess:       "<",
	TkLessEq:     "<=",
	TkAnd:        "&&",
	TkOr:         "||",
	TkNot:        "~",
	TkIf:         "if",
	TkFunc:       "func",
	TkReturn:     "return",
	TkWhile:      "while",
	TkBool:       "Bool",
	TkIdentifier: "Identifier",
	TkNumber:     "Number",
	TkString:     "String",
	TkAssign:     ":=",
	TkLBrace:     "{",
	TkRBrace:     "}",
	TkColon:      ":",
	TkComma:      ",",
	TkComment:    "Comment",
	TkElse:       "else",
	TkLParen:     "(",
	TkRParent:    ")",
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
	}
	if t.IsCmpOp() {
		return 2
	}
	if t.IsAddOp() {
		return 3
	}
	if t.IsMulOp() {
		return 4
	}
	return 0
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
	return (t.IsAddOp() || t.IsMulOp() || t.IsCmpOp() || t.IsLogOp()) && t != TkNot
}

// IsUnaryOp returns bool to indicate whether the token represents a
// right-binding operator.
func (t Token) IsUnaryOp() bool {
	return t.IsAddOp() || t == TkNot
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
