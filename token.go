package main

import "strconv"

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
	TkUnknown:    "TkUnknown",
	TkEof:        "TkEof",
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
	TkBool:       "TkBool",
	TkIdentifier: "TkIdentifier",
	TkNumber:     "TkNumber",
	TkString:     "TkString",
	TkAssign:     ":=",
	TkLBrace:     "{",
	TkRBrace:     "}",
	TkColon:      ":",
	TkComma:      ",",
	TkComment:    "TkComment",
	TkElse:       "else",
	TkLParen:     "(",
	TkRParent:    ")",
}

func (t Token) String() string {
	if t < Token(len(tokens)) {
		return tokens[t]
	}
	return "Token(" + strconv.Itoa(int(t)) + ")"
}

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
	panic("token is not an operator")
}

func (t Token) IsAddOp() bool {
	return t > addop_start && t < addop_end
}

func (t Token) IsMulOp() bool {
	return t > mulop_start && t < mulop_end
}

func (t Token) IsCmpOp() bool {
	return t > cmpop_start && t < cmpop_end
}

func (t Token) IsLogOp() bool {
	return t > logop_start && t < logop_end
}

func (t Token) IsBinOp() bool {
	return (t.IsAddOp() || t.IsMulOp() || t.IsCmpOp() || t.IsLogOp()) && t != TkNot
}

func (t Token) IsUnaryOp() bool {
	return t.IsAddOp() || t == TkNot
}

func (t Token) IsStmtKeyword() bool {
	return t > stmtkwd_start && t < stmtkwd_end
}

func (t Token) IsLiteral() bool {
	return t > lit_start && t < lit_end
}
