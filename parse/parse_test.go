// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package parse_test

import (
	"testing"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/parse"
)

const (
	BASIC = iota
	BINARY
	IDENT
	FILE
	VAR
)

func nodeTest(types []int, t *testing.T) func(node ast.Node) {
	i := 0
	return func(node ast.Node) {
		switch node.(type) {
		case *ast.File:
			if types[i] != FILE {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", FILE)
			}
		case *ast.BasicLit:
			if types[i] != BASIC {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", BASIC)
			}
		case *ast.BinaryExpr:
			if types[i] != BINARY {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", BINARY)
			}
		case *ast.Ident:
			if types[i] != IDENT {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", IDENT)
			}
		case *ast.VarExpr:
			if types[i] != VAR {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", IDENT)
			}
		default:
			t.Fatal("Walk index:", i, "Expected:", types[i], "Got: Unknown")
		}
		i++
	}
}

func TestParseFileBasic(t *testing.T) {
	var types = []int{FILE, IDENT, BINARY, BASIC, BASIC}

	f := parse.ParseFile("test.calc", "(+ 3 5)")
	ast.Walk(f, nodeTest(types, t))
}

func TestParseNested(t *testing.T) {
	var types = []int{FILE, IDENT, BINARY, BINARY, BASIC, BASIC, BASIC, BINARY,
		BASIC, BASIC}
	f := parse.ParseFile("test.calc", ";comment\n(+ (/ 9 3) 5 (- 3 1))")
	ast.Walk(f, nodeTest(types, t))
}

func TestParseVar(t *testing.T) {
	var types = []int{FILE, IDENT, VAR, IDENT, IDENT, BASIC}
	f := parse.ParseFile("test.calc", "(var a 5 int)")
	ast.Walk(f, nodeTest(types, t))

	types = []int{FILE, IDENT, VAR, IDENT, IDENT}
	f = parse.ParseFile("test.calc", "(var a int)")
	ast.Walk(f, nodeTest(types, t))

	types = []int{FILE, IDENT, VAR, IDENT, BASIC}
	f = parse.ParseFile("test.calc", "(var a 5)")
	ast.Walk(f, nodeTest(types, t))

	f = parse.ParseFile("test.calc", "(var a)")
	if f != nil {
		t.Fatal("Parser returned valid object but expected nil")
	}
}

func TestExpectFail(t *testing.T) {
	var tests = []string{
		"+ 3 5)",
		"(- 5)",
		"(3 5 +)",
		"(3 + 4)",
		"(+ 6 2",
		"(- 4 5)2",
		"(d",
		"(* a 3)",
		"(/ 5 b)",
		"(% / d)",
		"(& 3 5)",
		"((+ 3 5) 5)",
		"(* (- 2 6) (+ 4 2)())",
		";comment",
	}
	for _, src := range tests {
		if f := parse.ParseFile("test", src); f != nil {
			t.Log(src, "- not nil")
			t.Fail()
		}
	}
}
