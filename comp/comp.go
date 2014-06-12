// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

// Package comp comprises the code generation portion of the Calc
// programming language
package comp

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/parse"
	"github.com/rthornton128/calc/token"
)

type compiler struct {
	fp       *os.File
	fset     *token.FileSet
	errors   token.ErrorList
	offset   int
	curScope *ast.Scope
	topScope *ast.Scope
}

// CompileFile generates a C source file for the corresponding file
// specified by path. The .calc extension for the filename in path is
// replaced with .c for the C source output.
func CompileFile(path string) {
	var c compiler

	c.fset = token.NewFileSet()
	f, err := parse.ParseFile(c.fset, path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	path = path[:len(path)-len(filepath.Ext(path))]
	fp, err := os.Create(path + ".c")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fp.Close()

	c.fp = fp
	c.compFile(f)

	if c.errors.Count() != 0 {
		c.errors.Print()
		os.Exit(1)
	}
}

// CompileDir generates C source code for the Calc sources found in the
// directory specified by path. The C source file uses the same name as
// directory rather than any individual file.
func CompileDir(path string) {
	fs := token.NewFileSet()
	pkg, err := parse.ParseDir(fs, path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fp, err := os.Create(filepath.Join(path, filepath.Base(path)) + ".c")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fp.Close()

	c := &compiler{fp: fp, fset: fs}
	c.compPackage(pkg)

	if c.errors.Count() != 0 {
		c.errors.Print()
		os.Exit(1)
	}
}

/* Utility */

// Error adds an error to the compiler at the given position. The remaining
// arguments are used to generate the error message.
func (c *compiler) Error(pos token.Pos, args ...interface{}) {
	c.errors.Add(c.fset.Position(pos), args...)
}

func roundUp16(n int) int {
	if r := n % 16; r != 0 {
		return n + (16 - r)
	}
	return n
}

func (c *compiler) nextOffset() (offset int) {
	offset = c.offset
	c.offset += 4
	return
}

/* Scope */

func (c *compiler) openScope(s *ast.Scope) {
	c.curScope = s
}

func (c *compiler) closeScope() {
	c.curScope = c.curScope.Parent
}

/* Main Compiler */

func (c *compiler) compNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.AssignExpr:
		c.compAssignExpr(n)
	case *ast.BasicLit:
		c.compInt(n, "eax")
	case *ast.BinaryExpr:
		c.compBinaryExpr(n)
	case *ast.CallExpr:
		c.compCallExpr(n)
	case *ast.DeclExpr:
		c.compDeclExpr(n)
	case *ast.ExprList:
		for i := range n.List {
			c.compNode(n.List[i])
		}
	case *ast.Ident:
		c.compIdent(n, "movl(ebp+%d, eax);\n")
	case *ast.IfExpr:
		c.compIfExpr(n)
	case *ast.UnaryExpr:
		c.compUnaryExpr(n)
	case *ast.VarExpr:
		c.compVarExpr(n)
	}
	return
}

func (c *compiler) compAssignExpr(a *ast.AssignExpr) {
	ob := c.curScope.Lookup(a.Name.Name)
	if ob == nil {
		c.Error(a.Name.NamePos, "undeclared variable '", a.Name.Name, "'")
		return
	}

	if ob.Type == nil {
		ob.Type = typeOf(a.Value, c.curScope)
	} else {
		c.matchTypes(a.Name, a.Value)
	}

	ob.Value = a.Value

	switch n := ob.Value.(type) {
	case *ast.BasicLit:
		c.compInt(n, fmt.Sprintf("ebp+%d", ob.Offset))
		return
	case *ast.BinaryExpr:
		c.compBinaryExpr(n)
	case *ast.CallExpr:
		c.compCallExpr(n)
	case *ast.IfExpr:
		c.compIfExpr(n)
	case *ast.Ident:
		c.compIdent(n, fmt.Sprintf("movl(ebp+%%d, ebp+%d);\n", ob.Offset))
	}
	fmt.Fprintf(c.fp, "movl(eax, ebp+%d);\n", ob.Offset)
}

func (c *compiler) compBinaryExpr(b *ast.BinaryExpr) {
	if x, ok := c.compTryOptimizeBinaryOrInt(b); ok {
		fmt.Fprintf(c.fp, "setl(%d, eax);\n", x)
		return
	}
	c.compNode(b.List[0])

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
		case *ast.UnaryExpr:
			fmt.Fprintln(c.fp, "pushl(eax);")
			c.compUnaryExpr(n)
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
}

func (c *compiler) compCallExpr(e *ast.CallExpr) {
	offset := 4

	ob := c.curScope.Lookup(e.Name.Name)
	switch {
	case e.Name.Name == "main":
		c.Error(e.Name.NamePos, "illegal to call function 'main'")
	case ob == nil:
		c.Error(e.Name.NamePos, "call to undeclared function '", e.Name.Name, "'")
	case ob.Kind != ast.Decl:
		c.Error(e.Name.NamePos, "may not call object that is not a function")
	case len(ob.Value.(*ast.DeclExpr).Params) != len(e.Args):
		c.Error(e.Name.NamePos, "number of arguments in function call do not "+
			"match declaration, expected ", len(ob.Value.(*ast.DeclExpr).Params),
			" got ", len(e.Args))
	}

	decl := ob.Value.(*ast.DeclExpr)
	for i, v := range e.Args {
		atype, dtype := typeOf(v, c.curScope), typeOf(decl.Params[i], decl.Scope)
		if atype.Name != dtype.Name {
			c.Error(e.Name.NamePos, "type mismatch, argument ", i+1, " of ",
				e.Name.Name, " is of type ", atype.Name, " but expected ", dtype.Name)
		}
	}
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
	fmt.Fprintf(c.fp, "_%s();\n", e.Name.Name)
	return
}

func (c *compiler) compDeclExpr(d *ast.DeclExpr) {
	c.openScope(d.Scope)
	c.compScopeDecls()

	last := c.offset
	c.offset = 0
	for _, p := range d.Params {
		ob := c.curScope.Lookup(p.Name)
		ob.Offset = c.nextOffset()
	}
	x := c.countVars(d)
	fmt.Fprintf(c.fp, "void _%s(void) {\n", d.Name.Name)

	if x > 0 {
		fmt.Fprintf(c.fp, "enter(%d);\n", roundUp16(x))
		c.compNode(d.Body)
		fmt.Fprintln(c.fp, "leave();")
	} else {
		c.compNode(d.Body)
	}

	fmt.Fprintln(c.fp, "}")
	c.offset = last
	if d.Body != nil {
		c.matchTypes(d, d.Body)
	}
	c.closeScope()
	return
}

func (c *compiler) compFile(f *ast.File) {
	c.topScope = f.Scope
	c.curScope = c.topScope
	c.compTopScope()
}

func (c *compiler) compIdent(n *ast.Ident, format string) {
	ob := c.curScope.Lookup(n.Name)
	if ob == nil {
		panic("no offset for identifier")
	}
	fmt.Fprintf(c.fp, format, ob.Offset)
}

func (c *compiler) compIfExpr(n *ast.IfExpr) {
	if t := typeOf(n.Cond, c.curScope); t.Name != "int" {
		c.Error(n.Cond.Pos(), "Expression must be of type int, got ", t.Name)
	}
	c.compNode(n.Cond)
	fmt.Fprintln(c.fp, "if (*(int32_t *)ecx == 1) {")
	c.openScope(n.Scope)
	c.matchTypes(n, n.Then)
	c.compNode(n.Then)
	if n.Else != nil && !reflect.ValueOf(n.Else).IsNil() {
		c.matchTypes(n, n.Then)
		fmt.Fprintln(c.fp, "} else {")
		c.compNode(n.Else)
	}
	c.closeScope()
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

func (c *compiler) compPackage(p *ast.Package) {
	c.topScope = p.Scope
	c.curScope = c.topScope
	c.compTopScope()
}

func (c *compiler) compScopeDecls() {
	for k, v := range c.curScope.Table {
		if v.Kind == ast.Decl {
			fmt.Fprintf(c.fp, "void _%s(void);\n", k)
			defer c.compNode(v.Value)
		}
	}
}

func (c *compiler) compTopScope() {
	ob := c.curScope.Lookup("main")
	switch {
	case ob == nil:
		// BUG: token.NoPos will cause panic
		c.Error(token.NoPos, "no entry point, function 'main' not found")
	case ob.Kind != ast.Decl:
		c.Error(ob.NamePos, "no entry point, 'main' is not a function")
	case ob.Type == nil:
		c.Error(ob.NamePos, "'main' must be of type int but was declared as "+
			"void")
	case ob.Type.Name != "int":
		c.Error(ob.Type.NamePos, "'main' must be of type but declared as ",
			ob.Type.Name)
	}
	fmt.Fprintln(c.fp, "#include <stdio.h>")
	fmt.Fprintln(c.fp, "#include <runtime.h>")
	c.compScopeDecls()
	fmt.Fprintln(c.fp, "int main(void) {")
	fmt.Fprintln(c.fp, "stack_init();")
	fmt.Fprintln(c.fp, "_main();")
	fmt.Fprintln(c.fp, "printf(\"%d\\n\", *(int32_t *)eax);")
	fmt.Fprintln(c.fp, "stack_end();")
	fmt.Fprintln(c.fp, "return 0;")
	fmt.Fprintln(c.fp, "}")
}

func (c *compiler) compUnaryExpr(u *ast.UnaryExpr) {
	c.compNode(u.Value)
	fmt.Fprintln(c.fp, "setl(-1, edx);")
	fmt.Fprintln(c.fp, "mull(edx, eax);")
}

func (c *compiler) compVarExpr(v *ast.VarExpr) {
	ob := c.curScope.Lookup(v.Name.Name)
	ob.Offset = c.nextOffset()
	if ob.Value != nil && !reflect.ValueOf(ob.Value).IsNil() {
		if val, ok := ob.Value.(*ast.AssignExpr); ok {
			c.compAssignExpr(val)
			return
		}
		panic("parsing error occured, object's Value is not an assignment")
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

func (c *compiler) compTryOptimizeBinaryOrInt(e ast.Expr) (int, bool) {
	var ret int
	var ok bool
	switch t := e.(type) {
	case *ast.BasicLit:
		if t.Kind == token.INTEGER {
			i, err := strconv.Atoi(t.Lit)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			ret, ok = i, true
		}
	case *ast.BinaryExpr:
		for i, v := range t.List {
			var x int
			x, ok = c.compTryOptimizeBinaryOrInt(v)
			if !ok {
				break
			}
			if i == 0 {
				ret = x
				continue
			}
			switch t.Op {
			case token.ADD:
				ret += x
			case token.SUB:
				ret -= x
			case token.MUL:
				ret *= x
			case token.QUO:
				ret /= x
			case token.REM:
				ret %= x
			default:
				return 0, false
			}
		}
	}
	return ret, ok
}

func (c *compiler) matchTypes(a, b ast.Node) {
	atype, btype := typeOf(a, c.curScope), typeOf(b, c.curScope)

	switch {
	//case atype == "unknown":
	//c.Error(a.Pos(), "unknown type")
	case btype.Name == "unknown":
		c.Error(btype.Pos(), "unknown type")
	//case !validType(atype):
	//c.Error(a.Pos(), "invalid type: ", atype)
	case !validType(btype):
		c.Error(btype.Pos(), "invalid type: ", btype.Name)
	case atype.Name != btype.Name:
		c.Error(btype.Pos(), "type mismatch: ", btype.Name, " vs ", atype.Name)
	}
}
