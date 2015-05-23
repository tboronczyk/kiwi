/*
 * Copyright (c) 2012, 2015 Timothy Boronczyk
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  1. Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *
 *  2. Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *
 *  3. The names of the authors may not be used to endorse or promote products
 *     derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED "AS IS" AND WITHOUT ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, WITHOUT LIMITATION, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE.
 */

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

var eof = rune(0)

type scanner struct {
	r     *bufio.Reader
	atEOF bool
}

func NewScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r)}
}

func (s *scanner) read() rune {
	if s.atEOF {
		return eof
	}

	ch, _, err := s.r.ReadRune()
	if err != nil {
		s.atEOF = true
		ch = eof
	}
	return ch
}

func (s *scanner) unread() {
	s.r.UnreadRune()
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
	case ';':
		return token.SEMICOLON, ";"
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
			buf.WriteRune(ch)
			if ch == '\\' {
				ch = s.read()
				buf.WriteRune(ch)
			}
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
			return token.TRUE, str
		case "FALSE":
			return token.FALSE, str
		}
	}
	return token.IDENTIFIER, str
}

func (s *scanner) scanNumber() (token.Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); unicode.IsDigit(ch) {
			buf.WriteRune(ch)
		} else {
			s.unread()
			break
		}
	}
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
