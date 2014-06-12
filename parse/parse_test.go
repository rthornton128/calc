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

type Test struct {
	name  string
	src   string
	types []Type
	pass  bool
}

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
	UNKNOWN
	VAR
)

var typeStrings = []string{
	ASSIGN:  "assignexpr",
	BASIC:   "basiclit",
	BINARY:  "binaryexpr",
	CALL:    "callexpr",
	DECL:    "declexpr",
	IDENT:   "ident",
	IF:      "if",
	FILE:    "file",
	LIST:    "exprlist",
	UNARY:   "unaryexpr",
	UNKNOWN: "unknown",
	VAR:     "var",
}

func (t Type) String() string { return typeStrings[int(t)] }

func nodeTest(types []Type, t *testing.T) func(node ast.Node) {
	typ := UNKNOWN
	i := 0
	return func(node ast.Node) {
		switch node.(type) {
		case *ast.AssignExpr:
			typ = ASSIGN
		case *ast.BasicLit:
			typ = BASIC
		case *ast.BinaryExpr:
			typ = BINARY
		case *ast.CallExpr:
			typ = CALL
		case *ast.DeclExpr:
			typ = DECL
		case *ast.ExprList:
			typ = LIST
		case *ast.File:
			typ = FILE
		case *ast.Ident:
			t.Log("ident:", node.(*ast.Ident).Name)
			typ = IDENT
		case *ast.IfExpr:
			typ = IF
		case *ast.UnaryExpr:
			typ = UNARY
		case *ast.VarExpr:
			typ = VAR
		}
		if types[i] != typ {
			t.Fatal("Walk index:", i, "Expected:", types[i], "Got:", typ)
		}
		i++
	}
}

func handleTests(t *testing.T, tests []Test) {
	for _, test := range tests {
		f, _ := parse.ParseExpression(test.name, test.src)
		if f == nil && test.pass {
			t.Log(f == nil)
			t.Log(!test.pass)
			t.Fatal("Failed to parse")
		}
		ast.Walk(f, nodeTest(test.types, t))
	}
}

func TestParseBasic(t *testing.T) {
	tests := []Test{
		{"basic1", "24", []Type{BASIC}, true},
		{"basic2", "a", []Type{IDENT}, true},
	}
	handleTests(t, tests)
}

func TestParseBinary(t *testing.T) {
	tests := []Test{
		{"basic1", "(+ 2 3)", []Type{BINARY, BASIC, BASIC}, true},
		{"basic2", "(+ 2 b)", []Type{BINARY, BASIC, IDENT}, true},
		{"basic3", "(+ a b)", []Type{BINARY, IDENT, IDENT}, true},
		{"basic4", "+ 3 5)", []Type{}, false},
		{"basic5", "(- 5)", []Type{}, false},
		{"basic6", "(3 5 +)", []Type{}, false},
		{"basic7", "(3 + 4)", []Type{}, false},
		{"basic8", "(+ 6 2", []Type{}, false},
		{"basic9", "(d", []Type{}, false},
		{"basic10", "(% / d)", []Type{}, false},
		{"basic11", "(& 3 5)", []Type{}, false},
		{"basic12", "((+ 3 5) 5)", []Type{}, false},
		{"basic13", "(* (- 2 6) (+ 4 2)())", []Type{}, false},
	}
	handleTests(t, tests)
}

func TestParseCall(t *testing.T) {
	tests := []Test{
		{"call1", "(add 1 2)", []Type{CALL, IDENT, BASIC, BASIC}, true},
		{"call2", "(nothing)", []Type{CALL, IDENT}, true},
	}
	handleTests(t, tests)
}

func TestParseComment(t *testing.T) {
	tests := []Test{
		{"comment1", "2; comment", []Type{BASIC}, true},
		{"comment2", "2; comment\na", []Type{BASIC, IDENT}, true},
		{"comment3", "; comment\na", []Type{IDENT}, true},
		{"comment4", "(+ 2; comment\n3)", []Type{BINARY, BASIC, BASIC}, true},
		{"comment5", ";comment", []Type{}, false},
	}
	handleTests(t, tests)
}

func TestParseDecl(t *testing.T) {
	tests := []Test{
		{"decl1", "(decl func int 0)",
			[]Type{DECL, IDENT, IDENT, BASIC}, true},
		{"decl2", "(decl five int (+ 2 3))",
			[]Type{DECL, IDENT, IDENT, BINARY, BASIC, BASIC}, true},
		{"decl3", "(decl add(a b int) int (+ a b))",
			[]Type{DECL, IDENT, IDENT, IDENT, IDENT, BINARY, IDENT, IDENT}, true},
		{"decl3", "(decl main () int a)", []Type{}, false},
		{"decl3", "(decl main int ())", []Type{}, false},
	}
	handleTests(t, tests)
}

func TestParseIf(t *testing.T) {
	tests := []Test{
		{"if1", "(if 1 int 3)", []Type{IF, BASIC, IDENT, BASIC}, true},
		{"if2", "(if (< a b) int a ((+ b 1) b))",
			[]Type{IF, BINARY, IDENT, IDENT, IDENT, IDENT, LIST, BINARY, IDENT,
				BASIC, IDENT}, true},
		{"if3", "(if (< a b) ((= a b)))",
			[]Type{IF, BINARY, IDENT, IDENT, LIST, ASSIGN, IDENT, IDENT}, true},
	}
	handleTests(t, tests)
}

func TestParseNested(t *testing.T) {
	tests := []Test{
		{"nested1", "(+ (/ 9 3) 5 (- 3 1))",
			[]Type{BINARY, BINARY, BASIC, BASIC, BASIC, BINARY, BASIC, BASIC}, true},
	}
	handleTests(t, tests)
}

func TestParseUnary(t *testing.T) {
	var tests = []Test{
		{"unary1", "-24", []Type{UNARY, BASIC}, true},
		{"unary2", "-a", []Type{UNARY, IDENT}, true},
		{"unary3", "-(foo)", []Type{UNARY, CALL, IDENT}, true},
		{"unary4", "-(+ 2 3)", []Type{UNARY, BINARY, BASIC, BASIC}, true},
		{"unary5", "-(decl foo int)", []Type{}, false},
	}
	handleTests(t, tests)
}

func TestParseVar(t *testing.T) {
	tests := []Test{
		{"var1", "(var a int)", []Type{VAR, IDENT, IDENT}, true},
		{"var2", "(var (= a 5) int)",
			[]Type{VAR, IDENT, IDENT, ASSIGN, IDENT, BASIC}, true},
		{"var3", "(var (= a 5))", []Type{VAR, IDENT, ASSIGN, IDENT, BASIC}, true},
		{"var4", "(var a)", []Type{}, false},
		{"var5", "(var (+ a b))", []Type{}, false},
		{"var6", "(var 23)", []Type{}, false},
	}
	handleTests(t, tests)
}
