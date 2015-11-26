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

type FoldTest struct {
	src, expect string
}

func TestAssignmentFolding(t *testing.T) {
	test := FoldTest{src: "(= a (* 1 1))", expect: "1"}
	name := "assign"
	expr, _ := parse.ParseExpression(name, test.src)
	o := ir.FoldConstants(ir.MakeExpr(new(Package), expr))
	validate_constant(t, name, o.(*ir.Assignment).Rhs, test)
}

func TestBinaryFolding(t *testing.T) {
	tests := []FoldTest{
		{src: "(+ 21 21)", expect: "42"},
		{src: "(* 21 2)", expect: "42"},
		{src: "(/ 126 3)", expect: "42"},
		{src: "(- 0 42)", expect: "-42"},
		{src: "(% 5 2)", expect: "1"},
		{src: "(+ 1 2 3 4)", expect: "10"},
		{src: "(* 1 2 3 4)", expect: "24"},
		{src: "(/ 18 3 3)", expect: "2"},
		{src: "(- 15 5 10)", expect: "0"},
		{src: "(% 15 10 2)", expect: "1"},
		{src: "(== 42 42)", expect: "true"},
		{src: "(!= 24 24)", expect: "false"},
		{src: "(< 126 3)", expect: "false"},
		{src: "(<= 0 42)", expect: "true"},
		{src: "(> 5 2)", expect: "true"},
		{src: "(>= 3 4)", expect: "false"},
		{src: "(== true true)", expect: "true"},
		{src: "(!= true false)", expect: "true"},
	}
	for i, test := range tests {
		test_folding(t, fmt.Sprintf("binary%d", i), test)
	}
}

func TestCallFolding(t *testing.T) {
	src := "(fn (== 3 2) (+ 2 2))"
	name := "call"
	expr, _ := parse.ParseExpression(name, src)
	o := ir.FoldConstants(ir.MakeExpr(expr, ir.NewScope(nil)))
	validate_constant(t, name, o.(*ir.Call).Args[0], FoldTest{src, "false"})
	validate_constant(t, name, o.(*ir.Call).Args[1], FoldTest{src, "4"})
}

func TestIfFolding(t *testing.T) {
	src := "(if (== false (!= 3 3)):int (/ 9 3) (* 1 2 3))"
	name := "if"
	expr, _ := parse.ParseExpression(name, src)
	o := ir.FoldConstants(ir.MakeExpr(expr, ir.NewScope(nil)))
	validate_constant(t, name, o.(*ir.If).Cond, FoldTest{src, "true"})
	validate_constant(t, name, o.(*ir.If).Then, FoldTest{src, "3"})
	validate_constant(t, name, o.(*ir.If).Else, FoldTest{src, "6"})
}

func TestPackageFolding(t *testing.T) {
	fs := token.NewFileSet()
	f1, _ := parse.ParseFile(fs, "package", "(define f1 (func:int (+ 1 2)))")
	f2, _ := parse.ParseFile(fs, "package", "(define f2 (func:int (* 8 2)))")
	pkg := &ast.Package{Files: []*ast.File{f1, f2}}
	o := ir.FoldConstants(ir.MakePackage(pkg, "package"))
	o1 := o.(*ir.Package).Scope().Lookup("f1")
	o2 := o.(*ir.Package).Scope().Lookup("f2")
	validate_constant(t, "package", o1.(*ir.Define).Body.(*ir.Function).Body[0],
		FoldTest{"", "3"})
	validate_constant(t, "package", o2.(*ir.Define).Body.(*ir.Function).Body[0],
		FoldTest{"", "16"})
}

func TestUnaryFolding(t *testing.T) {
	tests := []FoldTest{
		{src: "-42)", expect: "-42"},
		{src: "+42", expect: "42"},
		{src: "+(- 2 4)", expect: "2"},
		{src: "-(+ 2 4)", expect: "-6"},
	}
	for i, test := range tests {
		test_folding(t, fmt.Sprintf("unary%d", i), test)
	}
}

func TestVarFolding(t *testing.T) {
	test := FoldTest{src: "(var (a:int):int (= a (/ 24 3)))", expect: "8"}
	name := "var"
	expr, _ := parse.ParseExpression(name, test.src)
	o := ir.FoldConstants(ir.MakeExpr(expr, ir.NewScope(nil)))
	o = o.(*ir.Variable).Body[0].(*ir.Assignment).Rhs
	validate_constant(t, name, o, test)
}

func test_folding(t *testing.T, name string, test FoldTest) {
	expr, _ := parse.ParseExpression(name, test.src)
	o := ir.FoldConstants(ir.MakeExpr(expr, ir.NewScope(nil)))
	validate_constant(t, name, o, test)
}

func validate_constant(t *testing.T, name string, o ir.Object, test FoldTest) {
	if c, ok := o.(*ir.Constant); !ok {
		t.Fatalf("%s: expected constant with value '%s' but got: %s",
			name, test.expect, o)
	} else {
		if c.String() != test.expect {
			t.Fatalf("%s: expected constant with value '%s' but got value: %s",
				name, test.expect, c.String())
		}
	}
}
