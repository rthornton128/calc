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
	"reflect"
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

/* Utility */

func roundUp16(n int) int {
	if r := n % 16; r != 0 {
		return n + (16 - r)
	}
	return n
}

/* Main Compiler */

func (c *compiler) compNode(node ast.Node) int {
	switch n := node.(type) {
	case *ast.BasicLit: /* nop */
	case *ast.BinaryExpr:
		return c.compBinaryExpr(n)
	case *ast.DeclExpr:
		return c.compDeclExpr(n)
	}
	return 0
}

func (c *compiler) compBinaryExpr(b *ast.BinaryExpr) int {
	vars := map[string]int{}
	noffset := 0

	switch n := b.List[0].(type) {
	case *ast.BasicLit:
		c.compInt(n, "eax")
	case *ast.BinaryExpr:
		c.compBinaryExpr(n)
	case *ast.CallExpr:
		c.compCallExpr(n)
	case *ast.Ident:
		noffset = c.compIdent(n, noffset, vars, "movl(ebp+%d, eax);\n")
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
		case *ast.CallExpr:
			fmt.Fprintln(c.fp, "pushl(eax);")
			c.compCallExpr(n)
			fmt.Fprintln(c.fp, "movl(eax, edx);")
			fmt.Fprintln(c.fp, "popl(eax);")
		case *ast.Ident:
			noffset = c.compIdent(n, noffset, vars, "movl(ebp+%d, edx);\n")
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

	return 0
}

func (c *compiler) compCallExpr(e *ast.CallExpr) int {
	//vars := map[string]int{}
	offset := 4

	for _, v := range e.Args {
		switch n := v.(type) {
		case *ast.BasicLit:
			c.compInt(n, fmt.Sprintf("esp+%d", offset))
			offset += 4
		}
	}
	fmt.Fprintf(c.fp, "%s();\n", e.Name.Name)
	return 0
}

func (c *compiler) compDeclExpr(d *ast.DeclExpr) int {
	if d.Name.Name == "main" {
		fmt.Fprintf(c.fp, "int %s(void) {\n", d.Name.Name)
		fmt.Fprintln(c.fp, "stack_init();")
	} else {
		fmt.Fprintf(c.fp, "void %s(void) {\n", d.Name.Name)
	}
	x := c.countVars(d)
	if x > 0 {
		fmt.Fprintf(c.fp, "enter(%d);\n", roundUp16(x))
		c.compNode(d.Body)
		fmt.Fprintln(c.fp, "leave();")
	} else {
		c.compNode(d.Body)
	}
	if d.Name.Name == "main" {
		fmt.Fprintln(c.fp, "printf(\"%d\\n\", *(int32_t *)eax);")
		fmt.Fprintln(c.fp, "stack_end();")
		fmt.Fprintf(c.fp, "return *(int32_t *)eax;\n")
	} else {
		fmt.Fprintf(c.fp, "return;\n")
	}
	fmt.Fprintln(c.fp, "}")
	return 0
}

func (c *compiler) compIdent(
	n *ast.Ident,
	noffset int,
	vars map[string]int,
	format string,
) int {
	offset, ok := vars[n.Name]
	if !ok {
		offset = noffset
		vars[n.Name] = offset
		noffset += 4
	}
	fmt.Fprintf(c.fp, format, offset)
	return noffset
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
			if v.Name != "main" {
				fmt.Fprintf(c.fp, "void %s(void);\n", k)
			}
			defer c.compNode(v.Value)
		}
	}
}

func (c *compiler) compFile(f *ast.File) {
	fmt.Fprintln(c.fp, "#include <stdio.h>")
	fmt.Fprintln(c.fp, "#include <runtime.h>")
	c.compScopeDecls(f.Scope)
}

func (c *compiler) countVars(n ast.Node) (x int) {
	if n != nil && !reflect.ValueOf(n).IsNil() {
		switch e := n.(type) {
		case *ast.DeclExpr:
			x = len(e.Params)
			x += c.countVars(e.Body)
		case *ast.IfExpr:
			x = c.countVars(e.Then)
			x = c.countVars(e.Else)
		case *ast.ExprList:
			for _, v := range e.List {
				x += c.countVars(v)
			}
		case *ast.VarExpr:
			x = 1
		}
	}
	return
}
