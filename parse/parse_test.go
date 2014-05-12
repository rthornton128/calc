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

type Type int

const (
	ASSIGN Type = iota
	BASIC
	BINARY
	DECL
	IDENT
	IF
	FILE
	LIST
	VAR
)

var typeStrings = []string{
	ASSIGN: "assignexpr",
	BASIC:  "basiclit",
	BINARY: "binaryexpr",
	DECL:   "declexpr",
	IDENT:  "ident",
	IF:     "if",
	FILE:   "file",
	LIST:   "exprlist",
	VAR:    "var",
}

func (t Type) String() string { return typeStrings[int(t)] }

func nodeTest(types []Type, t *testing.T) func(node ast.Node) {
	i := 0
	return func(node ast.Node) {
		switch node.(type) {
		case *ast.AssignExpr:
			if types[i] != ASSIGN {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", ASSIGN)
			}
		case *ast.BasicLit:
			if types[i] != BASIC {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", BASIC)
			}
		case *ast.BinaryExpr:
			if types[i] != BINARY {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", BINARY)
			}
		case *ast.DeclExpr:
			if types[i] != DECL {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", DECL)
			}
		case *ast.ExprList:
			if types[i] != LIST {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", LIST)
			}
		case *ast.File:
			if types[i] != FILE {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", FILE)
			}
		case *ast.Ident:
			t.Log("ident:", node.(*ast.Ident).Name)
			if types[i] != IDENT {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", IDENT)
			}
		case *ast.IfExpr:
			if types[i] != IF {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", IF)
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
	var types = []Type{FILE, IDENT, BINARY, BASIC, BASIC}

	f := parse.ParseFile("basic.calc", "(+ 3 5)")
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseDecl(t *testing.T) {
	var types = []Type{FILE, IDENT, DECL, IDENT, IDENT, BINARY, BASIC, BASIC}
	f := parse.ParseFile("decl1.calc", "(decl five int (+ 2 3))")
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, IDENT, DECL, IDENT, IDENT, IDENT, IDENT, BINARY,
		IDENT, IDENT}
	f = parse.ParseFile("decl2.calc", "(decl add(a b int) int (+ a b))")
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseIf(t *testing.T) {
	types := []Type{FILE, IDENT, IF, IDENT, IDENT, BASIC}
	f := parse.ParseFile("if.calc", "(if true int 3)")
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, IDENT, IF, BINARY, IDENT, IDENT, IDENT, IDENT, LIST,
		BINARY, IDENT, BASIC, IDENT}
	f = parse.ParseFile("if2.calc", "(if (< a b) int a ((+ b 1) b))")
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, IDENT, IF, BINARY, IDENT, IDENT, LIST, ASSIGN, IDENT,
		IDENT}
	f = parse.ParseFile("if3.calc", "(if (< a b) ((= a b)))")
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseNested(t *testing.T) {
	var types = []Type{FILE, IDENT, BINARY, BINARY, BASIC, BASIC, BASIC, BINARY,
		BASIC, BASIC}
	f := parse.ParseFile("nested.calc", ";comment\n(+ (/ 9 3) 5 (- 3 1))")
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseVar(t *testing.T) {
	var types = []Type{FILE, IDENT, VAR, IDENT, IDENT, BASIC}
	f := parse.ParseFile("var.calc", "(var a 5 int)")
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, IDENT, VAR, IDENT, IDENT}
	f = parse.ParseFile("var2.calc", "(var a int)")
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, IDENT, VAR, IDENT, BASIC}
	f = parse.ParseFile("var3.calc", "(var a 5)")
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	f = parse.ParseFile("var4.calc", "(var a)")
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
