package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

// convenient representation of EOF
const eof = rune(0)

// Scanner lexes a stream of characters (runes) into tokens and lexemes.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new scanner that reads from r.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{bufio.NewReader(r)}
}

// read manages the scanner's buffer and returns runes from it.
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		ch = eof
	}
	return ch
}

// unread pushes the most recently read rune back to the stream.
func (s *Scanner) unread() {
	s.r.UnreadRune()
}

// Scan consumes a lexeme from the reader's stream and returns its Token and
// string values.
func (s *Scanner) Scan() (Token, string) {
	s.skipWhitespace()
	ch := s.read()

	switch ch {
	case eof:
		return TkEOF, ""
	case '+':
		return TkAdd, "+"
	case '-':
		return TkSubtract, "-"
	case '*':
		return TkMultiply, "*"
	case '/':
		ch = s.read()
		if ch == '/' {
			return s.scanLineComment()
		}
		if ch == '*' {
			return s.scanMultiComment()
		}
		s.unread()
		return TkDivide, "/"
	case '%':
		return TkModulo, "%"
	case ':':
		ch = s.read()
		if ch == '=' {
			return TkAssign, ":="
		}
		s.unread()
		return TkColon, ":"
	case '=':
		return TkEqual, "="
	case '<':
		ch = s.read()
		if ch == '=' {
			return TkLessEq, "<="
		}
		s.unread()
		return TkLess, "<"
	case '>':
		ch = s.read()
		if ch == '=' {
			return TkGreaterEq, ">="
		}
		s.unread()
		return TkGreater, ">"
	case '&':
		ch = s.read()
		if ch == '&' {
			return TkAnd, "&&"
		}
		s.unread()
		return TkUnknown, "&"
	case '|':
		ch = s.read()
		if ch == '|' {
			return TkOr, "||"
		}
		s.unread()
		return TkUnknown, "|"
	case '~':
		ch = s.read()
		if ch == '=' {
			return TkNotEqual, "~="
		}
		s.unread()
		return TkIf, "~"
	case '(':
		return TkLParen, "("
	case ')':
		return TkRParen, ")"
	case '{':
		return TkLBrace, "{"
	case '}':
		return TkRBrace, "}"
	case ',':
		return TkComma, ","
	case '"':
		return s.scanString()
	case '`':
		s.unread()
		return s.scanIdent()
	}

	if unicode.IsLetter(ch) {
		s.unread()
		return s.scanIdent()
	}
	if unicode.IsDigit(ch) {
		s.unread()
		return s.scanNumber()
	}

	return TkUnknown, string(ch)
}

// skipWhitespace consumes whitespace by reading up to the first
// non-whitespace rune it encounters.
func (s *Scanner) skipWhitespace() {
	for {
		if ch := s.read(); !unicode.IsSpace(ch) {
			s.unread()
			break
		}
	}
}

// scanString consumes a string lexeme and returns its token and value. Escape
// sequences in the string are evaluated and replaced.
func (s *Scanner) scanString() (Token, string) {
	var buf bytes.Buffer
	for {
		if ch := s.read(); ch != '"' {
			// must have a closing quote
			if ch == eof {
				return TkUnknown, buf.String()
			}
			if ch == '\\' {
				switch s.read() {
				case '\\':
					ch = '\\'
					break
				case 'r':
					ch = '\r'
					break
				case 'n':
					ch = '\n'
					break
				case 't':
					ch = '\t'
					break
				case '"':
					ch = '"'
					break
				default:
					s.unread()
				}
			}
			buf.WriteRune(ch)
		} else {
			break
		}
	}
	return TkString, buf.String()
}

// scanIdent consumes an identifier lexeme and returns its token and value. An
// identifier will be recognized as a keyword if it matches the list of Kiwi
// keywords and is not escaped.
func (s *Scanner) scanIdent() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); unicode.IsLetter(ch) ||
			unicode.IsDigit(ch) || ch == '_' {
			buf.WriteRune(ch)
		} else {
			s.unread()
			break
		}
	}

	str := buf.String()
	if strings.IndexRune(str, '`') == 0 {
		str = str[1:]
	} else {
		switch strings.ToUpper(str) {
		case "ELSE":
			return TkElse, str
		case "FALSE":
			return TkBool, strings.ToUpper(str)
		case "FUNC":
			return TkFunc, str
		case "IF":
			return TkIf, str
		case "RETURN":
			return TkReturn, str
		case "TRUE":
			return TkBool, strings.ToUpper(str)
		case "WHILE":
			return TkWhile, str
		}
	}
	return TkIdentifier, str
}

// scanNumber consumes a numeric lexeme and returns its token and value. The
// numeric value may be an integer or real number.
func (s *Scanner) scanNumber() (Token, string) {
	var ch rune
	var buf bytes.Buffer

	buf.WriteRune(s.read())
	for {
		ch = s.read()
		if !unicode.IsDigit(ch) {
			break
		}
		buf.WriteRune(ch)
	}

	if ch == '.' {
		buf.WriteRune(ch)
		for {
			ch = s.read()
			if !unicode.IsDigit(ch) {
				break
			}
			buf.WriteRune(ch)
		}
	}
	s.unread()

	return TkNumber, buf.String()
}

// scanLineComment consumes a full-line comment and returns its token and
// lexeme value. The line comment ends when either a newline character or EOF
// is read.
func (s *Scanner) scanLineComment() (Token, string) {
	var buf bytes.Buffer
	buf.WriteString("//")
	for {
		ch := s.read()
		if ch == '\n' || ch == eof {
			break
		}
		buf.WriteRune(ch)
	}
	return TkComment, buf.String()
}

// scanMultiComment consumes a muti-line comment and returns its token and
// lexeme value. Nested multi-line comments are accomodated.
func (s *Scanner) scanMultiComment() (Token, string) {
	var buf bytes.Buffer
	buf.WriteString("/*")

	ch1 := s.read()
	ch2 := s.read()
	for {
		// must have a proper closing
		if ch1 == eof {
			return TkUnknown, buf.String()
		}
		if ch1 == '*' && ch2 == '/' {
			buf.WriteString("*/")
			break
		}
		// found a nested comment
		if ch1 == '/' && ch2 == '*' {
			_, str := s.scanMultiComment()
			buf.WriteString(str)
			ch2 = s.read()
		} else {
			buf.WriteRune(ch1)
		}
		ch1 = ch2
		ch2 = s.read()
	}
	return TkComment, buf.String()
}
