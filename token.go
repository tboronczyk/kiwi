package main

type Token uint

const (
	TkUnknown Token = iota
	TkEOF

	addopStart
	// addition-level operators
	TkAdd
	TkSubtract
	addopEnd

	mulopStart
	// multiplication-level operators
	TkMultiply
	TkDivide
	TkModulo
	mulopEnd

	cmpopStart
	// comparision operators
	TkEqual
	TkNotEqual
	TkGreater
	TkGreaterEq
	TkLess
	TkLessEq
	cmpopEnd

	logopStart
	// logic operators
	TkAnd
	TkOr
	TkNot
	logopEnd

	stmtkwdStart
	// statement keywords
	TkIf
	TkFunc
	TkReturn
	TkWhile
	stmtkwdEnd

	litStart
	// literal values
	TkBool
	TkIdentifier
	TkNumber
	TkString
	litEnd

	TkAssign
	TkLBrace
	TkRBrace
	TkColon
	TkComma
	TkComment
	TkElse
	TkLParen
	TkRParent

	// end of tokens
	endTokens
)

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
	return t > addopStart && t < addopEnd
}

func (t Token) IsMulOp() bool {
	return t > mulopStart && t < mulopEnd
}

func (t Token) IsCmpOp() bool {
	return t > cmpopStart && t < cmpopEnd
}

func (t Token) IsLogOp() bool {
	return t > logopStart && t < logopEnd
}

func (t Token) IsBinOp() bool {
	return (t.IsAddOp() || t.IsMulOp() || t.IsCmpOp() || t.IsLogOp()) && t != TkNot
}

func (t Token) IsUnaryOp() bool {
	return t.IsAddOp() || t == TkNot
}

func (t Token) IsStmtKeyword() bool {
	return t > stmtkwdStart && t < stmtkwdEnd
}

func (t Token) IsLiteral() bool {
	return t > litStart && t < litEnd
}
