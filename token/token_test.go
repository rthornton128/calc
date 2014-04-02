package token_test

import (
	"testing"

	"github.com/rthornton128/calc/token"
)

var test_expr = "(+ 2 3)\n(- 5 4)"

func TestFilePosition(t *testing.T) {
	var tests = []struct {
		col, row uint
		pos      token.Pos
	}{
		{1, 1, token.Pos(1)},
		{1, 2, token.Pos(8)},
		{7, 2, token.Pos(14)},
	}
	f, err := token.NewFile("test", test_expr, 1)
	if err != nil {
		t.Fatal(err)
	}
	f.AddLine(token.Pos(7))
	f.AddLine(token.Pos(14))
	for _, v := range tests {
		col, row := f.Position(v.pos)
		if col != v.col || row != v.row {
			t.Fatal("For:", v.pos, "Expected:", v.col, "and", v.row, "Got:",
				col, "and", row)
		}
	}
}

func TestLookup(t *testing.T) {
	var tests = []struct {
		str string
		tok token.Token
	}{
		{"+", token.ADD},
		{"%", token.REM},
		{"EOF", token.EOF},
		{"Integer", token.INTEGER},
		{"Comment", token.COMMENT},
		{"", token.ILLEGAL},
	}

	for i, v := range tests {
		if res := token.Lookup(v.str); res != v.tok {
			t.Fatal(i, "- Expected:", v.tok, "Got:", res)
		}
	}
}

func TestIsLiteral(t *testing.T) {
	var tests = []struct {
		tok token.Token
		exp bool
	}{
		{token.ADD, false},
		{token.REM, false},
		{token.EOF, false},
		{token.INTEGER, true},
		{token.COMMENT, false},
	}

	for _, v := range tests {
		if res := v.tok.IsLiteral(); res != v.exp {
			t.Fatal(v.tok, "- Expected:", v.exp, "Got:", res)
		}
	}
}

func TestIsOperator(t *testing.T) {
	var tests = []struct {
		tok token.Token
		exp bool
	}{
		{token.ADD, true},
		{token.REM, true},
		{token.EOF, false},
		{token.INTEGER, false},
		{token.COMMENT, false},
	}

	for _, v := range tests {
		if res := v.tok.IsOperator(); res != v.exp {
			t.Fatal(v.tok, "- Expected:", v.exp, "Got:", res)
		}
	}
}

func TestString(t *testing.T) {
	var tests = []struct {
		str string
		tok token.Token
	}{
		{"+", token.ADD},
		{"%", token.REM},
		{"EOF", token.EOF},
		{"Integer", token.INTEGER},
		{"Comment", token.COMMENT},
	}

	for i, v := range tests {
		if res := v.tok.String(); res != v.str {
			t.Fatal(i, "- Expected:", v.str, "Got:", res)
		}
	}
}

func TestValid(t *testing.T) {
	var tests = []struct {
		tok token.Token
		exp bool
	}{
		{token.ADD, true},
		{token.REM, true},
		{token.EOF, true},
		{token.INTEGER, true},
		{token.COMMENT, true},
		{token.Token(-1), false},
		{token.Token(999999), false},
	}

	for _, v := range tests {
		if res := v.tok.Valid(); res != v.exp {
			t.Fatal(v.tok, "- Expected:", v.exp, "Got:", res)
		}
	}
}
