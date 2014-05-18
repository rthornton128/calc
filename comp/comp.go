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
	fp       *os.File
	curScope *Scope
	topScope *Scope
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
	c.topScope = NewScope(nil)
	c.curScope = c.topScope
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
	case *ast.BasicLit:
		c.compInt(n, "eax")
	case *ast.BinaryExpr:
		return c.compBinaryExpr(n)
	case *ast.CallExpr:
		return c.compCallExpr(n)
	case *ast.DeclExpr:
		return c.compDeclExpr(n)
	case *ast.ExprList:
		for i := range n.List {
			c.compNode(n.List[i])
		}
	case *ast.Ident:
		c.compIdent(n, "movl(ebp+%d, eax);\n")
	case *ast.IfExpr:
		c.compIfExpr(n)
	}
	return 0
}

func (c *compiler) compBinaryExpr(b *ast.BinaryExpr) int {
	switch n := b.List[0].(type) {
	case *ast.BasicLit:
		c.compInt(n, "eax")
	case *ast.BinaryExpr:
		c.compBinaryExpr(n)
	case *ast.CallExpr:
		c.compCallExpr(n)
	case *ast.Ident:
		c.compIdent(n, "movl(ebp+%d, eax);\n")
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
			c.compIdent(n, "movl(ebp+%d, edx);\n")
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
		case token.AND:
			fmt.Fprintln(c.fp, "andl(eax, edx);")
		case token.EQL:
			fmt.Fprintln(c.fp, "eql(eax, edx);")
		case token.GTE:
			fmt.Fprintln(c.fp, "gel(eax, edx);")
		case token.GTT:
			fmt.Fprintln(c.fp, "gtl(eax, edx);")
		case token.LST:
			fmt.Fprintln(c.fp, "ltl(eax, edx);")
		case token.LTE:
			fmt.Fprintln(c.fp, "lel(eax, edx);")
		case token.NEQ:
			fmt.Fprintln(c.fp, "nel(eax, edx);")
		case token.OR:
			fmt.Fprintln(c.fp, "orl(eax, edx);")
		}
	}

	return 0
}

func (c *compiler) compCallExpr(e *ast.CallExpr) int {
	offset := 4

	for _, v := range e.Args {
		switch n := v.(type) {
		case *ast.BasicLit:
			c.compInt(n, fmt.Sprintf("esp+%d", offset))
		default:
			c.compNode(n)
			fmt.Fprintf(c.fp, "movl(eax, esp+%d);\n", offset)
		}
		offset += 4
	}
	fmt.Fprintf(c.fp, "%s();\n", e.Name.Name)
	return 0
}

func (c *compiler) compDeclExpr(d *ast.DeclExpr) int {
	c.openScope()
	offset := 0
	for _, p := range d.Params {
		c.curScope.Symbols[p.Name] = offset
		offset += 4
	}
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
	c.closeScope()
	return 0
}

func (c *compiler) compFile(f *ast.File) {
	fmt.Fprintln(c.fp, "#include <stdio.h>")
	fmt.Fprintln(c.fp, "#include <runtime.h>")
	c.compScopeDecls(f.Scope)
}

func (c *compiler) compIdent(n *ast.Ident, format string) {
	offset, ok := c.curScope.Lookup(n.Name)
	if !ok {
		panic("no offset for identifier")
	}
	fmt.Fprintf(c.fp, format, offset)
}

func (c *compiler) compIfExpr(n *ast.IfExpr) {
	switch e := n.Cond.(type) {
	case *ast.BasicLit:
		c.compInt(e, "eax")
	case *ast.BinaryExpr:
		c.compBinaryExpr(e)
	}
	fmt.Fprintln(c.fp, "if (*(int32_t *)ecx == 1) {")
	c.openScope()
	c.compNode(n.Then)
	if n.Type != nil {
		fmt.Fprintln(c.fp, "leave();")
		fmt.Fprintln(c.fp, "return;")
	}
	c.closeScope()
	if n.Else != nil && !reflect.ValueOf(n.Else).IsNil() {
		fmt.Fprintln(c.fp, "} else {")
		c.openScope()
		c.compNode(n.Else)
		c.closeScope()
		if n.Type != nil {
			fmt.Fprintln(c.fp, "leave();")
			fmt.Fprintln(c.fp, "return;")
		}
	}
	fmt.Fprintln(c.fp, "}")
}

func (c *compiler) compInt(n *ast.BasicLit, reg string) {
	i, err := strconv.Atoi(n.Lit)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Fprintf(c.fp, "setl(%d, %s);\n", i, reg)
}

func (c *compiler) compScopeDecls(s *ast.Scope) {
	for k, v := range s.Table {
		if v.Kind == ast.Decl {
			if v.Name != "main" {
				fmt.Fprintf(c.fp, "void %s(void);\n", k)
			}
			defer c.compNode(v.Value)
		}
	}
}

func (c *compiler) compVar(n *ast.VarExpr) {
	ob := n.Object
	if ob.Value != nil && !reflect.ValueOf(ob.Value).IsNil() {
		if ob.Type.Name != c.typeExpr(ob.Value) {
			fmt.Println("danger will robinson!")
		}
	}
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

func (c *compiler) typeExpr(e ast.Expr) string {
	switch n := e.(type) {
	case *ast.BasicLit:
		switch n.Kind {
		case token.INTEGER:
			return "int"
		}
	case *ast.BinaryExpr:
		return "int"
	case *ast.CallExpr:
		//ob := c.curScope.Lookup(n.Name.Name)
		//return c.typeExpr(ob.Value)
	case *ast.DeclExpr:
		return n.Type.Name
	case *ast.ExprList:
		return c.typeExpr(n.List[len(n.List)-1])
	case *ast.IfExpr:
		return n.Type.Name
	case *ast.VarExpr:
		return c.typeExpr(n.Object.Value)
	}
	return ""
}
