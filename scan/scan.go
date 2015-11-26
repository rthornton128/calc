// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

// Package scan is the Calc source code scanner
package scan

import (
	"bufio"
	"io"
	"unicode"

	"github.com/rthornton128/calc/token"
)

// Scanner...
type Scanner struct {
	ch      rune
	offset  int
	roffset int
	src     *bufio.Reader
	file    *token.File
}

// Init initializes Scanner and makes the source code ready to Scan
func (s *Scanner) Init(file *token.File, src io.Reader) {
	s.file = file
	s.offset, s.roffset = 0, 0
	s.src = bufio.NewReader(src)
	s.file.AddLine(s.offset) // TODO no sir, don't like it

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
	case 0:
		tok = token.EOF
	case '(':
		tok = token.LPAREN
	case ')':
		tok = token.RPAREN
	case ':':
		tok = token.COLON
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
		tok = s.selectToken('=', token.LTE, token.LST)
	case '>':
		tok = s.selectToken('=', token.GTE, token.GTT)
	case '=':
		tok = s.selectToken('=', token.EQL, token.ASSIGN)
	case '!':
		tok = s.selectToken('=', token.NEQ, token.ILLEGAL)
	case '&':
		tok = s.selectToken('&', token.AND, token.ILLEGAL)
	case '|':
		tok = s.selectToken('|', token.OR, token.ILLEGAL)
	case ';':
		s.skipComment()
		s.next()
		return s.Scan()
	default:
		tok = token.ILLEGAL
	}

	return
}

func (s *Scanner) next() {
	r, w, err := s.src.ReadRune()
	s.offset = s.roffset
	s.roffset += w
	if r == '\n' {
		s.file.AddLine(s.offset)
	}
	s.ch = r
	if err != nil {
		s.ch = 0
	}
}

func (s *Scanner) scanIdentifier() (string, token.Token, token.Pos) {
	start := s.offset
	var str string

	for unicode.IsLetter(s.ch) || unicode.IsDigit(s.ch) {
		str += string(s.ch)
		s.next()
	}
	return str, token.Lookup(str), s.file.Pos(start)
}

func (s *Scanner) scanNumber() (string, token.Token, token.Pos) {
	start := s.offset
	var str string

	for unicode.IsDigit(s.ch) {
		str += string(s.ch)
		s.next()
	}
	return str, token.INTEGER, s.file.Pos(start)
}

func (s *Scanner) selectToken(r rune, a, b token.Token) token.Token {
	if s.ch == r {
		s.next()
		return a
	}
	return b
}

func (s *Scanner) skipComment() {
	for s.ch != '\n' && s.ch != 0 {
		s.next()
	}
}

func (s *Scanner) skipWhitespace() {
	for unicode.IsSpace(s.ch) {
		s.next()
	}
}
