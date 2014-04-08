package parse_test

import (
	"testing"

	"github.com/rthornton128/calc1/ast"
	"github.com/rthornton128/calc1/parse"
)

const (
	FILE = iota
	BASIC
	BINARY
)

func TestParseFileBasic(t *testing.T) {
	f := parse.ParseFile("test.calc", "(+ 3 5)")
	i := 0

	var types = []int{FILE, BINARY, BASIC, BASIC}

	ast.Walk(f, func(node ast.Node) {
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
		}
		i++
	})
}

func TestParseNested(t *testing.T) {
	f := parse.ParseFile("test.calc", "(+ (/ 9 3) 5 (- 3 1))")
	i := 0

	var types = []int{FILE, BINARY, BINARY, BASIC, BASIC, BASIC, BINARY,
		BASIC, BASIC}

	ast.Walk(f, func(node ast.Node) {
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
		}
		i++
	})
}

func TestTemp(t *testing.T) {
	var tests = []string{
		"+ 3 5)",
		"(3 5 +)",
		"(3 + 4)",
		"(+ 6 2",
		"(d",
		"(* a 3)",
		"(/ 5 b)",
		"(% / d)",
		"(& 3 5)",
		"((+ 3 5) 5)",
		"(* (- 2 6) (+ 4 2)())",
	}
	for _, src := range tests {
		if f := parse.ParseFile("test", src); f != nil {
			t.Log(src, "- not nil")
			t.Fail()
		}
	}
}
