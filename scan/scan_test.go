// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package scan_test

import (
	"testing"

	"github.com/rthornton128/calc/scan"
	"github.com/rthornton128/calc/token"
)

func test_handler(t *testing.T, src string, expected []token.Token) {
	var s scan.Scanner
	s.Init(token.NewFile("", src), src)
	lit, tok, pos := s.Scan()
	for i := 0; tok != token.EOF; i++ {
		if tok != expected[i] {
			t.Fatal(pos, "Expected:", expected[i], "Got:", tok, lit)
		}
		lit, tok, pos = s.Scan()
	}
}

func TestNumber(t *testing.T) {
	src := "9"
	expected := []token.Token{
		token.INTEGER,
		token.EOF,
	}

	test_handler(t, src, expected)
}

func TestScan(t *testing.T) {
	src := "(+ 2 (- 4 1) (* 6 5) (% 10 2) (/ 9 3)); comment"
	expected := []token.Token{
		token.LPAREN,
		token.ADD,
		token.INTEGER,
		token.LPAREN,
		token.SUB,
		token.INTEGER,
		token.INTEGER,
		token.RPAREN,
		token.LPAREN,
		token.MUL,
		token.INTEGER,
		token.INTEGER,
		token.RPAREN,
		token.LPAREN,
		token.REM,
		token.INTEGER,
		token.INTEGER,
		token.RPAREN,
		token.LPAREN,
		token.QUO,
		token.INTEGER,
		token.INTEGER,
		token.RPAREN,
		token.RPAREN,
		token.EOF,
	}
	test_handler(t, src, expected)
}

func TestScanAllTokens(t *testing.T) {
	src := "()+-*/% 1 12\t 12345 123456789 | a as ! \\ \r ;"
	expected := []token.Token{
		token.LPAREN,
		token.RPAREN,
		token.ADD,
		token.SUB,
		token.MUL,
		token.QUO,
		token.REM,
		token.INTEGER,
		token.INTEGER,
		token.INTEGER,
		token.INTEGER,
		token.ILLEGAL,
		token.ILLEGAL,
		token.ILLEGAL,
		token.ILLEGAL,
		token.ILLEGAL,
		token.ILLEGAL,
		token.EOF,
	}
	test_handler(t, src, expected)
}
