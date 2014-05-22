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
	"github.com/rthornton128/calc/token"
)

type Type int

const (
	ASSIGN Type = iota
	BASIC
	BINARY
	CALL
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
	CALL:   "callexpr",
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
		case *ast.CallExpr:
			if types[i] != CALL {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", CALL)
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
	types := []Type{FILE, DECL, IDENT, IDENT, BINARY, BASIC, BASIC}
	src := "(decl main int (+ 1 2))"
	file := token.NewFile("basicdec", 1, len(src))
	f := parse.ParseFile(file, "basicdec", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseCall(t *testing.T) {
	types := []Type{FILE, DECL, IDENT, IDENT, CALL, IDENT, BASIC, BASIC,
		BASIC, BASIC}
	src := "(decl main int (add 1 2 3 4))"
	file := token.NewFile("call1", 1, len(src))
	f := parse.ParseFile(file, "call1", "(decl main int (add 1 2 3 4))")
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, DECL, IDENT, IDENT, CALL, IDENT}
	src = "(decl main int (nothing))"
	file = token.NewFile("call1", 1, len(src))
	f = parse.ParseFile(file, "call1", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseDecl(t *testing.T) {
	types := []Type{FILE, DECL, IDENT, IDENT, IDENT}
	src := "(decl func int a)"
	file := token.NewFile("decl1.calc", 1, len(src))
	f := parse.ParseFile(file, "decl1.calc", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, DECL, IDENT, IDENT, BINARY, BASIC, BASIC}
	src = "(decl five int (+ 2 3))"
	file = token.NewFile("decl2.calc", 1, len(src))
	f = parse.ParseFile(file, "decl1.calc", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, DECL, IDENT, IDENT, IDENT, IDENT, BINARY,
		IDENT, IDENT}
	src = "(decl add(a b int) int (+ a b))"
	file = token.NewFile("decl3.calc", 1, len(src))
	f = parse.ParseFile(file, "decl2.calc", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseIf(t *testing.T) {
	types := []Type{FILE, DECL, IDENT, IDENT, IF, IDENT, IDENT, BASIC}
	src := "(decl main int (if true int 3))"
	file := token.NewFile("if.calc", 1, len(src))
	f := parse.ParseFile(file, "if.calc", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, DECL, IDENT, IDENT, IF, BINARY, IDENT, IDENT, IDENT,
		IDENT, LIST, BINARY, IDENT, BASIC, IDENT}
	src = "(decl main int (if (< a b) int a ((+ b 1) b)))"
	file = token.NewFile("if2.calc", 1, len(src))
	f = parse.ParseFile(file, "if2.calc", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, DECL, IDENT, IDENT, IF, BINARY, IDENT, IDENT, LIST,
		ASSIGN, IDENT, IDENT}
	src = "(decl main int (if (< a b) ((= a b))))"
	file = token.NewFile("if3.calc", 1, len(src))
	f = parse.ParseFile(file, "if3.calc", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseNested(t *testing.T) {
	var types = []Type{FILE, DECL, IDENT, IDENT, BINARY, BINARY, BASIC, BASIC,
		BASIC, BINARY, BASIC, BASIC}
	src := ";comment\n(decl main int (+ (/ 9 3) 5 (- 3 1)))"
	file := token.NewFile("nested.calc", 1, len(src))
	f := parse.ParseFile(file, "nested.calc", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseVar(t *testing.T) {
	types := []Type{FILE, DECL, IDENT, IDENT, VAR, IDENT, IDENT}
	src := "(decl main int (var a int))"
	file := token.NewFile("var.calc", 1, len(src))
	f := parse.ParseFile(file, "var.calc", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, DECL, IDENT, IDENT, VAR, IDENT, IDENT, ASSIGN,
		IDENT, BASIC}
	src = "(decl main int (var (= a 5) int))"
	file = token.NewFile("var2.calc", 1, len(src))
	f = parse.ParseFile(file, "var2.calc", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{FILE, DECL, IDENT, IDENT, VAR, IDENT, ASSIGN, IDENT, BASIC}
	src = "(decl main int (var (= a 5)))"
	file = token.NewFile("var3.calc", 1, len(src))
	f = parse.ParseFile(file, "var3.calc", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
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
		"(var a)",
		"(decl main () int a)",
		"(decl main int ())",
		"(decl main int a)(decl main in b)",
		"(var a)",
		"(var (+ a b))",
		"(var 23)",
		"(var a int)(var a int)",
	}
	for _, src := range tests {
		file := token.NewFile("expectfail", 1, len(src))
		if f := parse.ParseFile(file, "expectfail", src); f != nil {
			t.Log(src, "- not nil")
			t.Fail()
		}
	}
}
