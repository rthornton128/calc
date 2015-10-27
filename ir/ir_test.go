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

func TestSimple(t *testing.T) {
	tests := []Test{
		{src: "(decl main int 42)", pass: true},
		{src: "(decl main bool true)", pass: true},
	}
	for i, test := range tests {
		test_handler(t, fmt.Sprintf("example%d", i), test)
	}
}

func TestDeclaration(t *testing.T) {
	tests := []Test{
		{src: "(decl fn(a b int) int true)", pass: false},
		{src: "(decl fn(a b int) bool a)", pass: false},
		{src: "(decl fn bool 24)", pass: false},
		{src: "(decl fn(a b int) int (+ a b))(decl main int (fn 2 3))",
			pass: true},
	}
	for i, test := range tests {
		test_handler(t, fmt.Sprintf("example%d", i), test)
	}
}

func TestBinary(t *testing.T) {
	tests := []Test{
		{src: "(decl main int (+ 2 3))", pass: true},
		{src: "(decl main int (* 2 3 4 5 6))", pass: true},
		{src: "(decl main int (/ (* 2 3) (% 4 5) (- 8 6)))", pass: true},
		{src: "(decl main bool (!= 2 3))", pass: true},
		{src: "(decl main int (+ main 1))", pass: false},
		{src: "(decl main(a b int) bool (< a b))", pass: true},
		{src: "(decl main(a bool) bool (== a true))", pass: true},
		{src: "(decl main(a bool) int (== a true))", pass: false},
		{src: "(decl main int (+ main 1))", pass: false},
	}
	for i, test := range tests {
		test_handler(t, fmt.Sprintf("example%d", i), test)
	}
}

func TestIf(t *testing.T) {
	tests := []Test{
		{src: "(decl main int (if (== 1 1) int 1 0))", pass: true},
		{src: "(decl main int (if (!= 1 1) int 1 true))", pass: false},
		{src: "(decl main int (if (< 1 1) int false 1))", pass: false},
		{src: "(decl main (a b int) int (if (<= a b) int 0 1))", pass: true},
		{src: "(decl main (a b int) int (if (> a b) int 0 1))", pass: true},
		{src: "(decl main (a b int) int (if (>= a b) int 0 1))", pass: true},
		{src: "(decl main (a bool) int (if (== a false) int 0 1))", pass: true},
		{src: "(decl main (a int) int (if (== a false) int 0 1))", pass: false},
	}
	for i, test := range tests {
		test_handler(t, fmt.Sprintf("example%d", i), test)
	}
}

func TestUnary(t *testing.T) {
	tests := []Test{
		{src: "(decl main int -24)", pass: true},
		{src: "(decl main int +(- 3 5))", pass: true},
	}
	for i, test := range tests {
		test_handler(t, fmt.Sprintf("example%d", i), test)
	}
}

func TestVar(t *testing.T) {
	tests := []Test{
		{src: "(decl main int ((var a int) a))", pass: true},
		{src: "(decl main int ((var a bool) a))", pass: false},
		{src: "(decl main int ((var (= a 42)) a))", pass: true},
		{src: "(decl main int ((var (= main 3)) a))", pass: false},
		{src: "(decl main bool ((var (= a true)) a))", pass: true},
		{src: "(decl main int ((var (= a true)) a))", pass: false},
		{src: "(decl main int ((var (= a 24) bool) a))", pass: false},
		{src: "(decl main int ((var a int)(= a 42) a))", pass: true},
	}
	for i, test := range tests {
		test_handler(t, fmt.Sprintf("example%d", i), test)
	}
}

func test_handler(t *testing.T, name string, test Test) {
	n, err := parse.ParseExpression(name, test.src)
	if err != nil {
		t.Fatal(err)
	}
	pkg := &ast.Package{
		Files: []*ast.File{n.(*ast.File)},
	}

	p := ir.MakePackage(pkg, name)
	t.Log(p)
	fset := token.NewFileSet()
	fset.Add(name, test.src)
	if err := ir.TypeCheck(p, fset); (err == nil) != test.pass {
		t.Logf("expected %v got %v", test.pass, err == nil)
		if err != nil {
			t.Log(err)
		}
		t.Fail()
	}
	t.Log(ir.FoldConstants(p))
	ir.Tag(p)

}
