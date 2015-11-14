// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir_test

import (
	"fmt"
	"testing"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/ir"
	"github.com/rthornton128/calc/parse"
	"github.com/rthornton128/calc/token"
)

type Test struct {
	src  string
	pass bool
}

func TestAssignment(t *testing.T) {
	tests := []Test{
		{src: "(= a 3)", pass: false},
		{src: "(func:int (a:int) (= a 0))", pass: false},
		{src: "(var:int (a:int) (= a true) 0)", pass: false},
	}
	for i, test := range tests {
		test_expression(t, fmt.Sprintf("assign%d", i), test)
	}
}

func TestBinary(t *testing.T) {
	tests := []Test{
		{src: "(+ 2 3)", pass: true},
		{src: "(* 2 3 4 5 6)", pass: true},
		{src: "(/ (* 2 3) (% 4 5) (- 8 6))", pass: true},
		{src: "(func:bool () (!= 2 3))", pass: true},
		{src: "(func:int () (+ main 1))", pass: false},
		{src: "(func:bool (a:int b:int) (< a b))", pass: true},
		{src: "(func:bool (a:bool) (== a true))", pass: true},
		{src: "(func:int (a:bool) (== a true))", pass: false},
		{src: "(var:int (main:bool) (+ main 1))", pass: false},
	}
	for i, test := range tests {
		test_expression(t, fmt.Sprintf("example%d", i), test)
	}
}

func TestCall(t *testing.T) {
	tests := []Test{
		{src: "(fn)", pass: false},
		//{src: "(decl fn (a int) int (a))", pass: false},
		//{src: "(decl fn int ((var a int) (a)))", pass: false},
	}
	for i, test := range tests {
		test_expression(t, fmt.Sprintf("call%d", i), test)
	}
}

func TestConstant(t *testing.T) {
	tests := []Test{
		{src: "42", pass: true},
		{src: "true", pass: true},
	}
	for i, test := range tests {
		test_expression(t, fmt.Sprintf("constant%d", i), test)
	}
}

func TestFile(t *testing.T) {
	tests := []Test{
		{src: "(define add:int (func:int (a:int b:int) (+ a b)))" +
			"(define main:int (func:int () (add 2 3)))",
			pass: true},
		{src: "(define equal (func:bool (a:int b:int) (== a b)))" +
			"(define main (func:int () (if:int (equal(+ 2 3) (*4 2)) 0 1)))",
			pass: true},
		{src: "(define equal (func:bool (a:int b:int) (== a b)))" +
			"(define main (func:int () (equal 2 3)))", pass: false},
		{src: "(define fn (func:int () 0))(define main (func:int () fn))",
			pass: false},
		{src: "(define fn (func:int (a:int) 0))" +
			"(define main (func:int () (fn 1 2)))", pass: false},
		{src: "(define fn (func:int (a:int b:int) 0))" +
			"(define main (func:int () (fn 1)))", pass: false},
		{src: "(define fn (func:int () 0))" +
			"(define main (func:int () (= fn 3) 0))", pass: false},
		{src: "(define fn (func:int (i:int b:bool) 0))" +
			"(define main:int (func:int () (fn 42 true)))", pass: true},
		{src: "(define fn (func:int (i:int b:bool) 0))" +
			"(define main (func:int () (fn 4 2)))", pass: false},
	}
	for i, test := range tests {
		test_file(t, fmt.Sprintf("file%d", i), test)
	}
}

func TestIf(t *testing.T) {
	tests := []Test{
		{src: "(if:int (== 1 1) 1 0)", pass: true},
		{src: "(if:int (!= 1 1) 1 true)", pass: false},
		{src: "(if:int (< 1 1) false 1)", pass: false},
		{src: "(if:int 1 0 1)", pass: false},
		{src: "(func:int (a:int b:int) (if:int  (<= a b) 0 1))", pass: true},
		{src: "(func:int (a:int b:int) (if:int (> a b) 0 1))", pass: true},
		{src: "(func:int (a:int b:int) (if:int (>= a b) 0 1))", pass: true},
		{src: "(func:int (a:bool) (if:int (== a false) 0 1))", pass: true},
		{src: "(func:int (a:int) (if:int (== a false) 0 1))", pass: false},
	}
	for i, test := range tests {
		test_expression(t, fmt.Sprintf("if%d", i), test)
	}
}

func TestUnary(t *testing.T) {
	tests := []Test{
		{src: "-24", pass: true},
		{src: "+(- 3 5)", pass: true},
	}
	for i, test := range tests {
		test_expression(t, fmt.Sprintf("unary%d", i), test)
	}
}

func TestVar(t *testing.T) {
	tests := []Test{
		{src: "(var:int (a:int) a)", pass: true},
		{src: "(var:int (a:bool) a)", pass: false},
		{src: "(var:int () (= a 42) a)", pass: false},
		{src: "(var:bool (a:bool) (= a true) a)", pass: true},
		{src: "(var:int (a:bool) (= a true) a)", pass: false},
		{src: "(var:int (a:bool) (= a true) a)", pass: false},
		{src: "(var:bool (a:int) (= a 24))", pass: false},
		{src: "(var:int (a:int) (= a 42) a)", pass: true},
	}
	for i, test := range tests {
		test_expression(t, fmt.Sprintf("var%d", i), test)
	}
}

func test_expression(t *testing.T, name string, test Test) {
	expr, err := parse.ParseExpression(name, test.src)
	if err != nil {
		t.Fatal(err)
	}
	test_handler(t, test, name, expr)
}

func test_file(t *testing.T, name string, test Test) {
	f, err := parse.ParseFile(token.NewFileSet(), name, test.src)
	if err != nil {
		t.Fatal(err)
	}
	test_handler(t, test, name, &ast.Package{Files: []*ast.File{f}})
}

func test_handler(t *testing.T, test Test, name string, n ast.Node) {
	var o ir.Object
	switch t := n.(type) {
	case *ast.DefineStmt:
		o = ir.MakeDefine(t, ir.NewScope(nil))
	case *ast.Package:
		o = ir.MakePackage(t, name)
	case ast.Expr:
		o = ir.MakeExpr(t, ir.NewScope(nil))
	}
	t.Log(o)
	fset := token.NewFileSet()
	fset.Add(name, test.src)
	if err := ir.TypeCheck(o, fset); (err == nil) != test.pass {
		t.Logf("expected %v got %v", test.pass, err == nil)
		if err != nil {
			t.Log(err)
		}
		t.Fail()
	} /*
		ir.Tag(o)*/
}
