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

/* Main Compiler */
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
	case *ast.DeclExpr:
		return c.compDeclExpr(n)
	default:
		return 0 /* can't be reached */
	}
}

func (c *compiler) compBinaryExpr(b *ast.BinaryExpr) int {
	var tmp int

	switch n := b.List[0].(type) {
	case *ast.BasicLit:
		c.compInt(n, "eax")
	case *ast.BinaryExpr:
		c.compBinaryExpr(n)
	}

	for _, node := range b.List[1:] {
		switch n := node.(type) {
		case *ast.BasicLit:
			c.compInt(n, "edx")
		case *ast.BinaryExpr:
			fmt.Fprintln(c.fp, "pushl(eax);")
			c.compBinaryExpr(n)
			fmt.Fprintln(c.fp, "movl(eax, edx);")
			fmt.Fprintln(c.fp, "popl(eax);")
		}
		switch b.Op {
		case token.ADD:
			fmt.Fprintln(c.fp, "addl(edx, eax);")
		case token.SUB:
			fmt.Fprintln(c.fp, "subl(edx, eax);")
		case token.MUL:
			fmt.Fprintln(c.fp, "mull(edx, eax);")
		case token.QUO:
			fmt.Fprintln(c.fp, "divl(edx, eax);")
		case token.REM:
			fmt.Fprintln(c.fp, "reml(edx, eax);")
		}
	}

	return tmp
}

func (c *compiler) compDeclExpr(d *ast.DeclExpr) int {
	fmt.Fprintf(c.fp, "int %s(void) {\n", d.Name.Name)
	if d.Name.Name == "main" {
		fmt.Fprintln(c.fp, "stack_init();")
	}
	x := len(d.Params) // TODO: need to count all locals, calls
	if x > 0 {
		fmt.Fprintf(c.fp, "enter(%d);\n", x)
		c.compNode(d.Body)
		fmt.Fprintln(c.fp, "leave();")
	} else {
		c.compNode(d.Body)
	}
	if d.Name.Name == "main" {
		fmt.Fprintln(c.fp, "printf(\"%d\\n\", *(int32_t *)eax);")
		fmt.Fprintln(c.fp, "stack_end();")
	}
	fmt.Fprintf(c.fp, "return *(int32_t *)eax;\n")
	fmt.Fprintln(c.fp, "}")
	return 0
}

func (c *compiler) compInt(n *ast.BasicLit, reg string) {
	i, err := strconv.Atoi(n.Lit)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Fprintf(c.fp, "setl(%d, %s);\n", i, reg)
}

func (c *compiler) compScopeDecls(scope *ast.Scope) {
	for k, v := range scope.Table {
		if v.Kind == ast.Decl {
			fmt.Fprintf(c.fp, "int %s(void);\n", k)
			defer c.compNode(v.Value)
		}
	}
}

func (c *compiler) compFile(f *ast.File) {
	fmt.Fprintln(c.fp, "#include <stdio.h>")
	fmt.Fprintln(c.fp, "#include <runtime.h>")
	c.compScopeDecls(f.Scope)
}
