package comp

import (
	"fmt"
	"os"

	"github.com/rthornton128/calc1/ast"
	"github.com/rthornton128/calc1/parse"
	//"github.com/rthornton128/calc1/token"
)

type compiler struct {
	fp *os.File
}

func CompileFile(fname, src string) {
	f := parse.ParseFile(fname, src)

	var c compiler
	fp, err := os.Create(fname + ".c")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.fp = fp
  c.compile(f)
}

func (c *compiler) compile(node ast.Node) {
	switch n := node.(type) {
	case *ast.File:
		c.file(n)
	case *ast.BasicLit:
		c.basicLit(n)
	case *ast.BinaryExpr:
		c.binaryExpr(n)
	}
}

func (c *compiler) basicLit(b *ast.BasicLit) {
	fmt.Fprint(c.fp, "push(", b.Lit, ");\n")
}

func (c *compiler) binaryExpr(b *ast.BinaryExpr) {
	for _, node := range b.List {
		switch n := node.(type) {
		case *ast.BasicLit:
			c.basicLit(n)
		case *ast.BinaryExpr:
			c.binaryExpr(n)
		}
	}

	fmt.Fprintln(c.fp, "edx = pop();")
  for i := 0; i < len(b.List[1:]); i++ {
		fmt.Fprintln(c.fp, "eax = pop();")
		fmt.Fprint(c.fp, "edx ", b.Op, "= eax;\n")
	}
  fmt.Fprintln(c.fp, "push(edx);")
}

func (c *compiler) file(f *ast.File) {
	fmt.Fprintln(c.fp, "#include \"runtime.h\"")
	fmt.Fprintln(c.fp, "int main(void) {")
  c.compile(f.Root)
  fmt.Fprintln(c.fp, "eax = pop();")
	fmt.Fprintln(c.fp, "return eax;")
	fmt.Fprintln(c.fp, "}")
}
