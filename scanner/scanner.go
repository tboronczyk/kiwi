package scanner

import (
	"bufio"
	"bytes"
	"github.com/tboronczyk/kiwi/token"
	"io"
	"strings"
	"unicode"
)

type Scanner interface {
	Scan() (token.Token, string)
}

const (
	bufSize = 3
	eof     = rune(0)
)

type scanner struct {
	r          *bufio.Reader
	chars      [bufSize]rune
	rPos, wPos uint8
}

func New(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r)}
}

func (s *scanner) read() rune {
	if s.rPos == s.wPos {
		if s.wPos++; s.wPos == bufSize {
			s.wPos = 0
		}
		ch, _, err := s.r.ReadRune()
		if err != nil {
			ch = eof
		}
		s.chars[s.wPos] = ch
	}

	if s.rPos++; s.rPos == bufSize {
		s.rPos = 0
	}
	ch := s.chars[s.rPos]
	return ch
}

func (s *scanner) unread() {
	if s.rPos == 0 {
		s.rPos = bufSize
	}
	s.rPos--
}

func (s *scanner) Scan() (token.Token, string) {
	s.skipWhitespace()
	ch := s.read()

	switch ch {
	case eof:
		return token.EOF, ""
	case '+':
		return token.ADD, "+"
	case '-':
		return token.SUBTRACT, "-"
	case '*':
		return token.MULTIPLY, "*"
	case '/':
		ch = s.read()
		if ch == '/' {
			return s.scanLineComment()
		}
		if ch == '*' {
			return s.scanMultiComment()
		}
		s.unread()
		return token.DIVIDE, "/"
	case '%':
		return token.MODULO, "%"
	case ':':
		ch = s.read()
		if ch == '=' {
			return token.ASSIGN, ":="
		}
		s.unread()
		return token.MALFORMED, ":"
	case '=':
		return token.EQUAL, "="
	case '<':
		ch = s.read()
		if ch == '=' {
			return token.LESS_EQ, "<="
		}
		s.unread()
		return token.LESS, "<"
	case '>':
		ch = s.read()
		if ch == '=' {
			return token.GREATER_EQ, ">="
		}
		s.unread()
		return token.GREATER, ">"
	case '&':
		ch = s.read()
		if ch == '&' {
			return token.AND, "&&"
		}
		s.unread()
		return token.MALFORMED, "&"
	case '|':
		ch = s.read()
		if ch == '|' {
			return token.OR, "||"
		}
		s.unread()
		return token.MALFORMED, "|"
	case '~':
		ch = s.read()
		if ch == '=' {
			return token.NOT_EQUAL, "~="
		}
		s.unread()
		return token.NOT, "~"
	case '(':
		return token.LPAREN, "("
	case ')':
		return token.RPAREN, ")"
	case '{':
		return token.LBRACE, "{"
	case '}':
		return token.RBRACE, "}"
	case '.':
		return token.DOT, "."
	case ',':
		return token.COMMA, ","
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

	return token.UNKNOWN, string(ch)
}

func (s *scanner) skipWhitespace() {
	for {
		if ch := s.read(); !unicode.IsSpace(ch) {
			s.unread()
			break
		}
	}
}

func (s *scanner) scanString() (token.Token, string) {
	var buf bytes.Buffer
	for {
		if ch := s.read(); ch != '"' {
			if ch == eof {
				return token.MALFORMED, buf.String()
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
	return token.STRING, buf.String()
}

func (s *scanner) scanIdent() (token.Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
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
		case "FUNC":
			return token.FUNC, str
		case "IF":
			return token.IF, str
		case "RETURN":
			return token.RETURN, str
		case "WHILE":
			return token.WHILE, str
		case "TRUE":
			return token.BOOL, strings.ToUpper(str)
		case "FALSE":
			return token.BOOL, strings.ToUpper(str)
		}
	}
	return token.IDENTIFIER, str
}

func (s *scanner) scanNumber() (token.Token, string) {
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
		ch = s.read()
		if !unicode.IsDigit(ch) {
			s.unread()
			s.unread()
			return token.NUMBER, buf.String()
		}
		buf.WriteRune('.')
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

	return token.NUMBER, buf.String()
}

func (s *scanner) scanLineComment() (token.Token, string) {
	var buf bytes.Buffer
	buf.WriteString("//")
	for {
		ch := s.read()
		if ch == '\n' || ch == eof {
			break
		}
		buf.WriteRune(ch)
	}
	return token.COMMENT, buf.String()
}

func (s *scanner) scanMultiComment() (token.Token, string) {
	var buf bytes.Buffer
	buf.WriteString("/*")

	ch1 := s.read()
	ch2 := s.read()
	for {
		if ch1 == eof {
			return token.MALFORMED, buf.String()
		}
		if ch1 == '*' && ch2 == '/' {
			buf.WriteString("*/")
			break
		}
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
	return token.COMMENT, buf.String()
}
