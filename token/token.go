package token

import (
	"strconv"
)

type Token uint8

const (
	UNKNOWN Token = iota
	MALFORMED
	EOF
	COMMENT

	addop_start
	ADD
	SUBTRACT
	addop_end

	mulop_start
	MULTIPLY
	DIVIDE
	MODULO
	mulop_end

	cmpop_start
	EQUAL
	NOT_EQUAL
	GREATER
	GREATER_EQ
	LESS
	LESS_EQ
	cmpop_end

	logop_start
	AND
	OR
	NOT
	logop_end

	ASSIGN
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	DOT
	COMMA

	stmtkwd_start
	IF
	WHILE
	FUNC
	RETURN
	stmtkwd_end

	lit_start
	IDENTIFIER
	TRUE
	FALSE
	NUMBER
	STRING
	lit_end
)

var tokens = [...]string{
	UNKNOWN:    "UNKNOWN",
	MALFORMED:  "MALFORMED",
	EOF:        "EOF",
	COMMENT:    "COMMENT",
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
	ASSIGN:     ":=",
	LPAREN:     "(",
	RPAREN:     ")",
	LBRACE:     "{",
	RBRACE:     "}",
	DOT:        ".",
	COMMA:      ",",
	IF:         "if",
	WHILE:      "while",
	FUNC:       "func",
	RETURN:     "return",
	TRUE:       "true",
	FALSE:      "false",
	NUMBER:     "NUMBER",
	STRING:     "STRING",
	IDENTIFIER: "IDENTIFIER",
}

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

func (t Token) IsExprOp() bool {
	return t.IsAddOp() || t == NOT
}

func (t Token) IsLiteral() bool {
	return t > lit_start && t < lit_end
}

func (t Token) IsStmtKeyword() bool {
	return t > stmtkwd_start && t < stmtkwd_end
}
