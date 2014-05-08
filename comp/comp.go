// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package comp

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/parse"
	"github.com/rthornton128/calc/token"
)

type compiler struct {
	fp *os.File
}

func CompileFile(fname, src string) {

	var c compiler
	fp, err := os.Create(fname + ".c")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fp.Close()

	f := parse.ParseFile(fname, src)
	if f == nil {
		os.Exit(1)
	}
	c.fp = fp
	c.compFile(f)
}

func (c *compiler) compNode(node ast.Node) int {
	switch n := node.(type) {
	case *ast.BasicLit:
		i, err := strconv.Atoi(n.Lit)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return i
	case *ast.BinaryExpr:
		return c.compBinaryExpr(n)
	default:
		return 0 /* can't be reached */
	}
}

func (c *compiler) compBinaryExpr(b *ast.BinaryExpr) int {
	var tmp int

	tmp = c.compNode(b.List[0])

	for _, node := range b.List[1:] {
		switch b.Op {
		case token.ADD:
			tmp += c.compNode(node)
		case token.SUB:
			tmp -= c.compNode(node)
		case token.MUL:
			tmp *= c.compNode(node)
		case token.QUO:
			tmp /= c.compNode(node)
		case token.REM:
			tmp %= c.compNode(node)
		}
	}

	return tmp
}

func (c *compiler) compDecl(ob *ast.Object) {
	fmt.Fprintf(c.fp, "%s %s(void) {\n", ob.Type.Name, ob.Name)
	fmt.Fprintf(c.fp, "printf(\"%%d\", %d);\n", c.compNode(ob.Value))
	fmt.Fprintf(c.fp, "return *(int32_t *)eax;\n")
	fmt.Fprintln(c.fp, "}")
}

func (c *compiler) compFile(f *ast.File) {
	fmt.Fprintln(c.fp, "#include <stdio.h>")
	fmt.Fprintln(c.fp, "#include <runtime.h>")
	c.compDecl(f.Scope.Lookup("main"))
}
