// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package parse_test

import (
	"fmt"
	"testing"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/parse"
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
	UNARY
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
	UNARY:  "unaryexpr",
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
		case *ast.UnaryExpr:
			if types[i] != UNARY {
				t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", UNARY)
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
	types := []Type{DECL, IDENT, IDENT, BINARY, BASIC, BASIC}
	src := "(decl main int (+ 1 2))"
	f := parse.ParseExpression("basicdec", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseCall(t *testing.T) {
	types := []Type{DECL, IDENT, IDENT, CALL, IDENT, BASIC, BASIC,
		BASIC, BASIC}
	src := "(decl main int (add 1 2 3 4))"
	f := parse.ParseExpression("call1", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{DECL, IDENT, IDENT, CALL, IDENT}
	src = "(decl main int (nothing))"
	f = parse.ParseExpression("call2", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseDecl(t *testing.T) {
	types := []Type{DECL, IDENT, IDENT}
	src := "(decl func int)"
	f := parse.ParseExpression("decl1", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{DECL, IDENT, IDENT, BINARY, BASIC, BASIC}
	src = "(decl five int (+ 2 3))"
	f = parse.ParseExpression("decl2", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{DECL, IDENT, IDENT, IDENT, IDENT, BINARY,
		IDENT, IDENT}
	src = "(decl add(a b int) int (+ a b))"
	f = parse.ParseExpression("decl3", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseIf(t *testing.T) {
	types := []Type{IF, BASIC, IDENT, BASIC}
	src := "(if 1 int 3)"
	f := parse.ParseExpression("if1", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{IF, BINARY, IDENT, IDENT, IDENT,
		IDENT, LIST, BINARY, IDENT, BASIC, IDENT}
	src = "(if (< a b) int a ((+ b 1) b))"
	f = parse.ParseExpression("if2", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{IF, BINARY, IDENT, IDENT, LIST, ASSIGN, IDENT, IDENT}
	src = "(if (< a b) ((= a b)))"
	f = parse.ParseExpression("if3", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseNested(t *testing.T) {
	var types = []Type{BINARY, BINARY, BASIC, BASIC, BASIC, BINARY, BASIC,
		BASIC}
	src := ";comment\n(+ (/ 9 3) 5 (- 3 1))"
	f := parse.ParseExpression("nested", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))
}

func TestParseUnary(t *testing.T) {
	var tests = []struct {
		name  string
		src   string
		types []Type
		pass  bool
	}{
		{"unary1", "-24", []Type{UNARY, BASIC}, true},
		{"unary2", "-a", []Type{UNARY, IDENT}, true},
		{"unary3", "-(foo)", []Type{UNARY, CALL, IDENT}, true},
		{"unary4", "-(+ 2 3)", []Type{UNARY, BINARY, BASIC, BASIC}, true},
		{"unary5", "-(decl foo int)", []Type{UNARY, DECL, IDENT, IDENT}, false},
	}
	for _, test := range tests {
		f := parse.ParseExpression(test.name, test.src)
		if f == nil && test.pass {
			t.Log(f == nil)
			t.Log(!test.pass)
			t.Fatal("Failed to parse")
		}
		ast.Walk(f, nodeTest(test.types, t))
	}
}

func TestParseVar(t *testing.T) {
	types := []Type{VAR, IDENT, IDENT}
	src := "(var a int)"
	f := parse.ParseExpression("var", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{VAR, IDENT, IDENT, ASSIGN, IDENT, BASIC}
	src = "(var (= a 5) int)"
	f = parse.ParseExpression("var2", src)
	if f == nil {
		t.Fatal("Failed to parse")
	}
	ast.Walk(f, nodeTest(types, t))

	types = []Type{VAR, IDENT, ASSIGN, IDENT, BASIC}
	src = "(var (= a 5))"
	f = parse.ParseExpression("var3", src)
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
	}
	for i, src := range tests {
		name := fmt.Sprint("expectfail", i)
		if f := parse.ParseExpression(name, src); f != nil {
			t.Log(name, ":", src, "- not nil")
			t.Fail()
		}
	}
}
