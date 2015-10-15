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
}

// CompileFile generates a C source file for the corresponding file
// specified by path. The .calc extension for the filename in path is
// replaced with .c for the C source output.
func CompileFile(path string) error {
	var c compiler

	c.fset = token.NewFileSet()
	f, err := parse.ParseFile(c.fset, path, nil)
	if err != nil {
		return err
	}

	// type checking pass
	var t ast.TypeChecker
	ast.Walk(f, &t)

	// TODO optimization pass(es)

	path = path[:len(path)-len(filepath.Ext(path))]
	fp, err := os.Create(path + ".c")
	if err != nil {
		return err
	}
	defer fp.Close()

	c.fp = fp

	c.emitHeaders()
	c.compFile(f)
	c.emitMain()

	if c.errors.Count() != 0 {
		return c.errors
	}
	return nil
}

// CompileDir generates C source code for the Calc sources found in the
// directory specified by path. The C source file uses the same name as
// directory rather than any individual file.
func CompileDir(path string) error {
	fs := token.NewFileSet()
	pkg, err := parse.ParseDir(fs, path)
	if err != nil {
		return err
	}

	fp, err := os.Create(filepath.Join(path, filepath.Base(path)) + ".c")
	if err != nil {
		return err
	}
	defer fp.Close()

	c := &compiler{fp: fp, fset: fs}

	c.emitHeaders()
	c.compPackage(pkg)
	c.emitMain()

	if c.errors.Count() != 0 {
		return c.errors
	}
	return nil
}

/* Utility */

// Error adds an error to the compiler at the given position. The remaining
// arguments are used to generate the error message.
func (c *compiler) Error(pos token.Pos, args ...interface{}) {
	c.errors.Add(c.fset.Position(pos), args...)
}

func (c *compiler) nextOffset() (offset int) {
	offset = c.offset
	c.offset += 1
	return
}

func (c *compiler) emit(s string, args ...interface{}) {
	fmt.Fprintf(c.fp, s, args...)
}

func (c *compiler) emitln(args ...interface{}) {
	fmt.Fprintln(c.fp, args...)
}

func (c *compiler) emitHeaders() {
	c.emitln("#include <stdio.h>")
	c.emitln("#include <runtime.h>")
}

func (c *compiler) emitMain() {
	c.emitln("int main(void) {")
	c.emitln("stack_init();")
	c.emitln("_main();")
	c.emitln("printf(\"%d\\n\", (int32_t) ax);")
	c.emitln("stack_end();")
	c.emitln("return 0;")
	c.emitln("}")
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
		c.compInt(n, "ax")
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
		c.compIdent(n, "ax = *(bp+%d);\n")
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

	ob.Value = a.Value

	switch n := ob.Value.(type) {
	case *ast.BasicLit:
		c.compInt(n, fmt.Sprintf("*(bp+%d)", ob.Offset))
		return
	case *ast.BinaryExpr:
		c.compBinaryExpr(n)
	case *ast.CallExpr:
		c.compCallExpr(n)
	case *ast.IfExpr:
		c.compIfExpr(n)
	case *ast.Ident:
		c.compIdent(n, fmt.Sprintf("*(bp+%d) = *(bp+%%d);\n", ob.Offset))
		return
	}
	c.emit("*(bp+%d) = ax;\n", ob.Offset)
}

func (c *compiler) compBinaryExpr(b *ast.BinaryExpr) {
	if x, ok := c.compTryOptimizeBinaryOrInt(b); ok {
		c.emit("ax = %d;\n", x)
		return
	}
	c.compNode(b.List[0])

	for _, node := range b.List[1:] {
		switch n := node.(type) {
		case *ast.BasicLit:
			c.compInt(n, "dx")
		case *ast.BinaryExpr:
			c.emitln("push(ax);")
			c.compBinaryExpr(n)
			c.emitln("dx = ax;")
			c.emitln("pop(ax);")
		case *ast.CallExpr:
			c.emitln("push(ax);")
			c.compCallExpr(n)
			c.emitln("dx = ax;")
			c.emitln("pop(ax);")
		case *ast.Ident:
			c.compIdent(n, "dx = *(bp+%d);\n")
		case *ast.UnaryExpr:
			c.emitln("push(ax);")
			c.compUnaryExpr(n)
			c.emitln("dx = ax;")
			c.emitln("pop(ax);")
		}
		switch b.Op {
		case token.ADD:
			c.emitln("ax += dx;")
		case token.SUB:
			c.emitln("ax -= dx;")
		case token.MUL:
			c.emitln("ax *= dx;")
		case token.QUO:
			c.emitln("ax /= dx;")
		case token.REM:
			c.emitln("ax %= dx;")
		case token.AND:
			c.emitln("ax &= dx;")
		case token.OR:
			c.emitln("ax |= dx;")
		case token.EQL:
			c.emitln("ax = ax == dx;")
		case token.GTE:
			c.emitln("ax = ax >= dx;")
		case token.GTT:
			c.emitln("ax = ax > dx;")
		case token.LST:
			c.emitln("ax = ax < dx;")
		case token.LTE:
			c.emitln("ax = ax <= dx;")
		case token.NEQ:
			c.emitln("ax = ax != dx;")
		}
	}
}

func (c *compiler) compCallExpr(e *ast.CallExpr) {
	offset := 1

	ob := c.curScope.Lookup(e.Name.Name)
	switch {
	case e.Name.Name == "main":
		c.Error(e.Name.NamePos, "illegal to call function 'main'")
		return
	case ob == nil:
		c.Error(e.Name.NamePos, "call to undeclared function '", e.Name.Name, "'")
		return
	case ob.Kind != ast.Decl:
		c.Error(e.Name.NamePos, "may not call object that is not a function")
		return
	case len(ob.Value.(*ast.DeclExpr).Params) != len(e.Args):
		c.Error(e.Name.NamePos, "number of arguments in function call do not "+
			"match declaration, expected ", len(ob.Value.(*ast.DeclExpr).Params),
			" got ", len(e.Args))
		return
	}

	for _, v := range e.Args {
		switch n := v.(type) {
		case *ast.BasicLit:
			c.compInt(n, fmt.Sprintf("*(sp+%d)", offset))
		default:
			c.compNode(n)
			c.emit("*(sp+%d) = ax;\n", offset)
		}
		offset += 1
	}
	c.emit("_%s();\n", e.Name.Name)
	return
}

func (c *compiler) compDeclExpr(d *ast.DeclExpr) {
	c.openScope(d.Scope)

	c.offset = 0
	for _, p := range d.Params {
		ob := c.curScope.Lookup(p.Name)
		ob.Offset = c.nextOffset()
	}

	c.emit("void _%s(void) {\n", d.Name.Name)
	if x := c.countVars(d); x > 0 {
		c.emit("enter(%d);\n", x)
		c.compNode(d.Body)
		c.emitln("leave();")
	} else {
		c.compNode(d.Body)
	}
	c.emitln("}")

	c.closeScope()
	return
}

func (c *compiler) compFile(f *ast.File) {
	c.curScope = f.Scope
	c.compDeclProto(f)
}

func (c *compiler) compIdent(n *ast.Ident, format string) {
	ob := c.curScope.Lookup(n.Name)
	if ob == nil {
		panic("no offset for identifier")
	}
	fmt.Fprintf(c.fp, format, ob.Offset)
}

func (c *compiler) compIfExpr(n *ast.IfExpr) {
	c.compNode(n.Cond)

	c.emitln("if ((int32_t)ax == 1) {")
	c.openScope(n.Scope)
	c.compNode(n.Then)
	if n.Else != nil && !reflect.ValueOf(n.Else).IsNil() {
		c.emitln("} else {")
		c.compNode(n.Else)
	}
	c.closeScope()
	c.emitln("}")
}

func (c *compiler) compInt(n *ast.BasicLit, reg string) {
	i, err := strconv.Atoi(n.Lit)
	if err != nil {
		c.Error(n.Pos(), "bad conversion:", err)
	}
	c.emit("%s = %d;\n", reg, i)
}

func (c *compiler) compPackage(p *ast.Package) {
	c.curScope = p.Scope

	for _, f := range p.Files {
		c.compFile(f)
	}
}

func (c *compiler) compDeclProto(f *ast.File) {
	for _, decl := range f.Decls {
		c.emit("void _%s(void);\n", decl.Name.Name)
		defer c.compNode(decl)
	}
}

func (c *compiler) compUnaryExpr(u *ast.UnaryExpr) {
	c.compNode(u.Value)
	c.emitln("ax *= -1;")
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
	// TODO: implement proper zero value code for additional types
	c.emit("*(bp+%d) = 0;\n", ob.Offset)
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

// TODO remove entirely and place into optimization pass
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
