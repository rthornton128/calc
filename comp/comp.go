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
	"strings"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/parse"
	"github.com/rthornton128/calc/token"
)

type compiler struct {
	fp       *os.File
	fset     *token.FileSet
	errors   token.ErrorList
	nextID   int
	curScope *ast.Scope
}

// CompileFile generates a C source file for the corresponding file
// specified by path. The .calc extension for the filename in path is
// replaced with .c for the C source output.
func CompileFile(path string, opt bool) error {
	var c compiler

	c.fset = token.NewFileSet()
	f, err := parse.ParseFile(c.fset, path, nil)
	if err != nil {
		return err
	}

	// type checking pass
	ast.Walk(f, &ast.TypeChecker{ErrorHandler: c.Error})

	// optimization pass(es)
	if opt {
		ast.Walk(f, &OptConstantFolder{})
	}

	path = path[:len(path)-len(filepath.Ext(path))]
	fp, err := os.Create(path + ".c")
	if err != nil {
		return err
	}
	defer fp.Close()

	c.fp = fp
	c.nextID = 1

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
func CompileDir(path string, opt bool) error {
	fs := token.NewFileSet()
	pkg, err := parse.ParseDir(fs, path)
	if err != nil {
		return err
	}

	// type checking pass
	ast.Walk(pkg, &ast.TypeChecker{})

	// optimization pass(es)
	if opt {
		ast.Walk(pkg, &OptConstantFolder{})
	}

	fp, err := os.Create(filepath.Join(path, filepath.Base(path)) + ".c")
	if err != nil {
		return err
	}
	defer fp.Close()

	c := &compiler{fp: fp, fset: fs, nextID: 1}

	c.emitHeaders()
	c.compPackage(pkg)
	c.emitMain()

	if c.errors.Count() != 0 {
		return c.errors
	}
	return nil
}

/* Utility */

func cType(t ast.Type) string {
	switch t {
	case ast.Int:
		return "int32_t"
	default:
		return "int"
	}
}

// Error adds an error to the compiler at the given position. The remaining
// arguments are used to generate the error message.
func (c *compiler) Error(pos token.Pos, args ...interface{}) {
	c.errors.Add(c.fset.Position(pos), args...)
}

func (c *compiler) emit(s string, args ...interface{}) {
	fmt.Fprintf(c.fp, s, args...)
}

func (c *compiler) emitln(args ...interface{}) {
	fmt.Fprintln(c.fp, args...)
}

func (c *compiler) emitHeaders() {
	c.emitln("#include <stdio.h>")
	c.emitln("#include <stdint.h>")
}

func (c *compiler) emitMain() {
	c.emitln("int main(void) {")
	c.emitln("printf(\"%d\\n\", _main());")
	c.emitln("return 0;")
	c.emitln("}")
}

func (c *compiler) getID(id int) int {
	if id == 0 {
		id = c.nextID
		c.nextID++
	}
	return id
}

/* Main Compiler */

type temp struct {
	ID int
}

func (t *temp) Pos() token.Pos { return token.NoPos }
func (t *temp) End() token.Pos { return token.NoPos }

func (c *compiler) compNode(node ast.Node) string {
	var str string
	switch n := node.(type) {
	case *ast.AssignExpr:
		c.compAssignExpr(n)
	case *ast.BasicLit:
		str = c.compBasicLit(n)
	case *ast.BinaryExpr:
		str = c.compBinaryExpr(n)
	case *ast.CallExpr:
		str = c.compCallExpr(n)
	case *ast.DeclExpr:
		c.compDeclExpr(n)
	case *ast.ExprList:
		for _, e := range n.List {
			str = c.compNode(e)
		}
	case *ast.Ident:
		str = c.compIdent(n)
	case *ast.IfExpr:
		str = c.compIfExpr(n)
	case *ast.UnaryExpr:
		str = c.compUnaryExpr(n)
	case *ast.Value:
		str = fmt.Sprintf("%d", n.Value)
	case *ast.VarExpr:
		c.compVarExpr(n)
	}
	return str
}

func (c *compiler) compAssignExpr(a *ast.AssignExpr) {
	c.emit("%s = %s;\n", c.compNode(a.Name), c.compNode(a.Value))
}

func (c *compiler) compBasicLit(b *ast.BasicLit) string {
	// use switch b.Kind {} when future types are added
	i, err := strconv.Atoi(b.Lit)
	if err != nil {
		c.Error(b.Pos(), "bad conversion:", err)
	}
	return fmt.Sprintf("%d", i)
}

func (c *compiler) compBinaryExpr(b *ast.BinaryExpr) string {
	lhs := ast.Node(b.List[0])

	for _, rhs := range b.List[1:] {
		c.emit("%s _v%d = %s %s %s;\n",
			cType(b.RealType),
			c.nextID,
			c.compNode(lhs),
			b.Op,
			c.compNode(rhs))
		lhs = &temp{ID: c.getID(b.ID)}
	}
	b.ID = c.getID(lhs.(*temp).ID)
	return fmt.Sprintf("_v%d", b.ID)
}

func (c *compiler) compCallExpr(e *ast.CallExpr) string {
	args := make([]string, len(e.Args))
	for i, a := range e.Args {
		args[i] = c.compNode(a)
	}
	return fmt.Sprintf("_%s(%s)", e.Name.Name, strings.Join(args, ","))
}

func (c *compiler) compDeclExpr(d *ast.DeclExpr) {
	params := make([]string, len(d.Params))
	for i, p := range d.Params {
		params[i] = cType(p.Object.RealType) + " " + c.compNode(p)
	}
	c.emit("%s _%s(%s) {\n",
		cType(d.RealType),
		d.Name.Name,
		strings.Join(params, ","))
	c.emit("return %s;\n}\n", c.compNode(d.Body))
}

func (c *compiler) compFile(f *ast.File) {
	c.compDeclProto(f)
}

func (c *compiler) compIdent(i *ast.Ident) string {
	i.Object.ID = c.getID(i.Object.ID)
	return fmt.Sprintf("_v%d", i.Object.ID)
}

func (c *compiler) compIfExpr(n *ast.IfExpr) string {
	t := &temp{ID: c.getID(0)}
	c.emit("%s _v%d = 0;\n", cType(n.RealType), t.ID)
	c.emit("if (%s == 1) {\n", c.compNode(n.Cond))
	c.emit("_v%d = %s;\n", t.ID, c.compNode(n.Then))
	if n.Else != nil && !reflect.ValueOf(n.Else).IsNil() {
		c.emitln("} else {")
		c.emit("_v%d = %s;\n", t.ID, c.compNode(n.Else))
	}
	c.emitln("}")
	return fmt.Sprintf("_v%d", t.ID)
}

func (c *compiler) compPackage(p *ast.Package) {
	for _, f := range p.Files {
		c.compFile(f)
	}
}

func (c *compiler) compDeclProto(f *ast.File) {
	for _, decl := range f.Decls {
		params := make([]string, len(decl.Params))
		for i, p := range decl.Params {
			params[i] = cType(p.Object.RealType) + " " + c.compNode(p)
		}
		c.emit("%s _%s(%s);\n",
			cType(decl.RealType),
			decl.Name.Name,
			strings.Join(params, ","))
		defer c.compNode(decl)
	}
}

func (c *compiler) compUnaryExpr(u *ast.UnaryExpr) string {
	//fmt.Println(u.Value, c.compNode(u.Value))
	return fmt.Sprintf("-%s", c.compNode(u.Value))
}

func (c *compiler) compVarExpr(v *ast.VarExpr) {
	c.emit("%s %s = 0;\n", cType(v.RealType), c.compNode(v.Name))
	if v.Object.Value != nil && !reflect.ValueOf(v.Object.Value).IsNil() {
		if val, ok := v.Object.Value.(*ast.AssignExpr); ok {
			c.compAssignExpr(val)
			return
		}
		panic("parsing error occured, object's Value is not an assignment")
	}
}
