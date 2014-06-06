// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package scan

import (
	"unicode"

	"github.com/rthornton128/calc/token"
)

type Scanner struct {
	ch      rune
	offset  int
	roffset int
	src     string
	file    *token.File
}

func (s *Scanner) Init(file *token.File, src string) {
	s.file = file
	s.offset, s.roffset = 0, 0
	s.src = src
	s.file.AddLine(s.offset)

	s.next()
}

func (s *Scanner) Scan() (lit string, tok token.Token, pos token.Pos) {
	s.skipWhitespace()

	if unicode.IsLetter(s.ch) {
		return s.scanIdentifier()
	}

	if unicode.IsDigit(s.ch) {
		return s.scanNumber()
	}

	ch := s.ch
	lit, pos = string(s.ch), s.file.Pos(s.offset)
	s.next()
	switch ch {
	case '(':
		tok = token.LPAREN
	case ')':
		tok = token.RPAREN
	case ',':
		tok = token.COMMA
	case '+':
		tok = token.ADD
	case '-':
		tok = token.SUB
	case '*':
		tok = token.MUL
	case '/':
		tok = token.QUO
	case '%':
		tok = token.REM
	case '<':
		if s.ch == '=' {
			tok = token.LTE
			s.next()
		} else {
			tok = token.LST
		}
	case '>':
		if s.ch == '=' {
			tok = token.GTE
			s.next()
		} else {
			tok = token.GTT
		}
	case '=':
		if s.ch == '=' {
			tok = token.EQL
			s.next()
		} else {
			tok = token.ASSIGN
		}
	case '!':
		if s.ch == '=' {
			tok = token.NEQ
			s.next()
		} else {
			tok = token.ILLEGAL
		}
	case '&':
		if s.ch == '&' {
			tok = token.AND
			s.next()
		} else {
			tok = token.ILLEGAL
		}
	case '|':
		if s.ch == '|' {
			tok = token.OR
			s.next()
		} else {
			tok = token.ILLEGAL
		}
	case ';':
		s.skipComment()
		s.next()
		return s.Scan()
	default:
		if s.offset >= len(s.src)-1 {
			tok = token.EOF
		} else {
			tok = token.ILLEGAL
		}
	}

	return
}

func (s *Scanner) next() {
	s.ch = rune(0)
	if s.roffset < len(s.src) {
		s.offset = s.roffset
		s.ch = rune(s.src[s.offset])
		if s.ch == '\n' {
			s.file.AddLine(s.offset)
		}
		s.roffset++
	}
}

func (s *Scanner) scanIdentifier() (string, token.Token, token.Pos) {
	start := s.offset

	for unicode.IsLetter(s.ch) || unicode.IsDigit(s.ch) {
		s.next()
	}
	offset := s.offset
	if s.ch == rune(0) {
		offset++
	}
	lit := s.src[start:offset]
	return lit, token.Lookup(lit), s.file.Pos(start)
}

func (s *Scanner) scanNumber() (string, token.Token, token.Pos) {
	start := s.offset

	for unicode.IsDigit(s.ch) {
		s.next()
	}
	offset := s.offset
	if s.ch == rune(0) {
		offset++
	}
	return s.src[start:offset], token.INTEGER, s.file.Pos(start)
}

func (s *Scanner) skipComment() {
	for s.ch != '\n' && s.offset < len(s.src)-1 {
		s.next()
	}
}

func (s *Scanner) skipWhitespace() {
	for unicode.IsSpace(s.ch) {
		s.next()
	}
}
