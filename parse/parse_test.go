// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package parse_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/parse"
	"github.com/rthornton128/calc/token"
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

type Tester struct {
	i     int
	t     *testing.T
	types []Type
}

func (t *Tester) Visit(n ast.Node) bool {
	var typ Type
	switch n.(type) {
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
		typ = IDENT
	case *ast.IfExpr:
		typ = IF
	case *ast.UnaryExpr:
		typ = UNARY
	case *ast.VarExpr:
		typ = VAR
	}
	if t.i >= len(t.types) {
		t.t.Logf("exceeded expected number of types (%d)", len(t.types))
		t.t.Fail()
		return false
	}
	if t.types[t.i] != typ {
		t.t.Log("Walk index:", t.i, "Expected:", t.types[t.i], "Got:", typ)
		t.t.Fail()
	}
	t.i++
	return true
}

func handleTests(t *testing.T, tests []Test) {
	for _, test := range tests {
		e, err := parse.ParseExpression(test.name, test.src)
		checkTest(t, test, e, err)
	}
}

func handleFileTests(t *testing.T, tests []Test) {
	for _, test := range tests {
		f, err := parse.ParseFile(token.NewFileSet(), test.name, test.src)
		checkTest(t, test, f, err)
	}
}

func checkTest(t *testing.T, test Test, n ast.Node, err error) {
	if err != nil {
		t.Logf("error: %s", err)
	}
	if n == nil && len(test.types) != 0 {
		t.Logf("%s: expr is nil, expected %d types", test.name, len(test.types))
		t.Fail()
	}

	if test.pass {
		ast.Walk(n, &Tester{types: test.types, t: t})
	}
	if t.Failed() {
		t.Logf("%s: %#v", test.name, n)
		t.Fatalf("%s: failed parsing expression: %s", test.name, test.src)
	}
}

func TestParseBasic(t *testing.T) {
	tests := []Test{
		{"integer", "24", []Type{BASIC}, true},
		{"var", "a", []Type{IDENT}, true},
	}
	handleTests(t, tests)
}

func TestParseBinary(t *testing.T) {
	tests := []Test{
		{"simple", "(+ 2 3)", []Type{BINARY, BASIC, BASIC}, true},
		{"one-var", "(+ 2 b)", []Type{BINARY, BASIC, IDENT}, true},
		{"two-vars", "(+ a b)", []Type{BINARY, IDENT, IDENT}, true},
		{"single-operand", "(- 5)", []Type{}, false},
		{"post-fix", "(3 5 +)", []Type{}, false},
		{"infix", "(3 + 4)", []Type{}, false},
		{"no-closing", "(+ 6 2", []Type{}, false},
		{"extra-open", "(d", []Type{}, false},
		{"modulus-quotient", "(% / d)", []Type{}, false},
		{"binary-and", "(& 3 5)", []Type{}, false},
		{"no-operator-nested", "((+ 3 5) 5)", []Type{}, false},
		{"multi-nested-with-empty", "(* (- 2 6) (+ 4 2)())", []Type{}, false},
	}
	handleTests(t, tests)
}

func TestParseCall(t *testing.T) {
	tests := []Test{
		{"no-args", "(nothing)", []Type{CALL}, true},
		{"two-args", "(add 1 2)", []Type{CALL, BASIC, BASIC}, true},
	}
	handleTests(t, tests)
}

func TestParseComment(t *testing.T) {
	tests := []Test{
		{"simple", "2; comment", []Type{BASIC}, true},
		{"nested-between-expr", "2; comment\na", []Type{BASIC, IDENT}, true},
		{"first-line", "; comment\na", []Type{IDENT}, true},
		{"nested-comment", "(+ 2; comment\n3)", []Type{BINARY, BASIC, BASIC}, true},
		{"comment-only", ";comment", []Type{}, false},
	}
	handleTests(t, tests)
}

func TestParseDecl(t *testing.T) {
	tests := []Test{
		{"simple", "(decl func int 0)",
			[]Type{DECL, BASIC}, true},
		{"no-param-binary", "(decl five int (+ 2 3))",
			[]Type{DECL, BINARY, BASIC, BASIC}, true},
		{"two-param-binary", "(decl add(a b int) int (+ a b))",
			[]Type{DECL, BINARY, IDENT, IDENT}, true},
		{"empty-params", "(decl main () int a)", []Type{}, false},
		{"empty-expr-list", "(decl main int ())", []Type{}, false},
		{"duplicate-param", "(decl fn (dup dup int) int 0)", []Type{}, false},
		{"no-open", "decl main int)", []Type{}, false},
		{"nested-decl", "(decl main int (decl fn int))", []Type{}, false},
	}
	handleTests(t, tests)
}

func TestParseDeclFile(t *testing.T) {
	tests := []Test{
		{"simple", "(decl main int 0)", []Type{FILE, DECL, BASIC}, true},
		{"no-source-no-file", "", []Type{}, false},
		{"no-decls", "42", []Type{}, false},
		{"duplicate-decl", "(decl fn int 1)(decl fn int 1)",
			[]Type{FILE, DECL, BINARY, DECL, BINARY}, false},
		{"redeclared-var-decl", "(var a int)(decl a int 1)",
			[]Type{FILE, VAR, DECL, BINARY}, false},
	}
	handleFileTests(t, tests)
}

func TestParseIf(t *testing.T) {
	tests := []Test{
		{"then-only", "(if false int 3)", []Type{IF, BASIC, BASIC}, true},
		{"then-else", "(if false int 3 4)", []Type{IF, BASIC, BASIC, BASIC}, true},
		{"no-type", "(if false 0 1)", []Type{}, false},
		{"integer-cond", "(if 1 int 3)", []Type{IF, BASIC, BASIC}, true},
		{"var-cond", "(if asdf int 3)", []Type{IF, IDENT, BASIC}, true},
		{"var-keyword", "(if var int 3)", []Type{}, false},
		{"logical-cond-nested-binary-then", "(if (< a b) int a ((+ b 1) b))",
			[]Type{IF, BINARY, IDENT, IDENT, IDENT, LIST, BINARY, IDENT,
				BASIC, IDENT}, true},
		{"logical-cond-assign-then", "(if (< a b) int ((= a b)))",
			[]Type{IF, BINARY, IDENT, IDENT, LIST, ASSIGN, IDENT, IDENT}, true},
	}
	handleTests(t, tests)
}

func TestParseUnary(t *testing.T) {
	var tests = []Test{
		{"negate-integer", "-24", []Type{UNARY, BASIC}, true},
		{"negate-var", "-a", []Type{UNARY, IDENT}, true},
		{"negate-call", "-(foo)", []Type{UNARY, CALL, IDENT}, true},
		{"positive-binary", "+(+ 2 3)", []Type{UNARY, BINARY, BASIC, BASIC}, true},
		{"positive-decl", "+(decl foo int)", []Type{}, false},
	}
	handleTests(t, tests)
}

func TestParseVar(t *testing.T) {
	tests := []Test{
		{"simple", "(var a int)", []Type{VAR}, true},
		{"with-assign", "(var (= a 5) int)",
			[]Type{VAR, ASSIGN, BASIC}, true},
		{"infer-type", "(var (= a 5))", []Type{VAR, ASSIGN, BASIC}, true},
		{"no-type", "(var a)", []Type{}, false},
		{"binary", "(var (+ a b))", []Type{}, false},
		{"integer", "(var 23)", []Type{}, false},
		{"redeclare", "(decl fn int ((var a int)(var a bool)0))", []Type{}, false},
	}
	handleTests(t, tests)
}

func TestParseFile(t *testing.T) {
	test := Test{"bad-ext", "", []Type{}, false}

	// test for file with bad extension
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Log(err)
	}
	defer func() {
		f.Close()
		err := os.Remove(f.Name())
		t.Log(err)
	}()
	n, err := parse.ParseFile(token.NewFileSet(), f.Name(), "")
	checkTest(t, test, n, err)

	test = Test{"simple.calc", "(decl main int 0)", []Type{FILE, DECL, BASIC},
		true}
	f, err = os.Create(test.name)
	if err != nil {
		t.Log(err)
	}
	defer func() {
		f.Close()
		if err := os.Remove(f.Name()); err != nil {
			t.Log(err)
		}
	}()

	_, err = f.WriteString(test.src)
	if err != nil {
		t.Fatal(err)
	}

	n, err = parse.ParseFile(token.NewFileSet(), f.Name(), "")
	checkTest(t, test, n, err)
}

func TestDirectory(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	var path string
	for _, p := range strings.Split(gopath, ";") {
		tmp := filepath.Join(p, "src", "github.com", "rthornton128", "calc",
			"examples", "package")
		t.Log("testing path:", tmp)
		if _, err := os.Stat(tmp); err == nil {
			path = tmp
			break
		}
	}

	t.Log("using path:", path)
	_, err := parse.ParseDir(token.NewFileSet(), path)
	if err != nil {
		t.Fatal(err)
	}
}
