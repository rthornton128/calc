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
	"strings"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/ir"
	"github.com/rthornton128/calc/parse"
	"github.com/rthornton128/calc/token"
)

type compiler struct {
	fp     *os.File
	fset   *token.FileSet
	errors token.ErrorList
}

// CompileFile generates a C source file for the corresponding file
// specified by path. The .calc extension for the filename in path is
// replaced with .c for the C source output.
func CompileFile(path string, opt bool) error {
	fset := token.NewFileSet()
	f, err := parse.ParseFile(fset, path, "")
	if err != nil {
		return err
	}

	pkg := ir.MakePackage(&ast.Package{
		Files: []*ast.File{f},
	}, filepath.Base(path))

	if err := ir.TypeCheck(pkg, fset); err != nil {
		return err
	}
	if opt {
		pkg = ir.FoldConstants(pkg).(*ir.Package)
	}
	ir.Tag(pkg)

	path = path[:len(path)-len(filepath.Ext(path))]
	fp, err := os.Create(path + ".c")
	if err != nil {
		return err
	}
	defer fp.Close()

	c := &compiler{fp: fp}

	c.emitHeaders()
	c.compPackage(pkg)
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
	fset := token.NewFileSet()
	p, err := parse.ParseDir(fset, path)
	if err != nil {
		return err
	}

	pkg := ir.MakePackage(p, filepath.Base(path))
	if err := ir.TypeCheck(pkg, fset); err != nil {
		return err
	}
	if opt {
		pkg = ir.FoldConstants(pkg).(*ir.Package)
	}
	ir.Tag(pkg)

	fp, err := os.Create(filepath.Join(path, filepath.Base(path)) + ".c")
	if err != nil {
		return err
	}
	defer fp.Close()

	c := &compiler{fp: fp, fset: fset}

	c.emitHeaders()
	c.compPackage(pkg)
	c.emitMain()

	if c.errors.Count() != 0 {
		return c.errors
	}
	return nil
}

/* Utility */

func cType(t ir.Type) string {
	switch t {
	case ir.Int:
		return "int32_t"
	case ir.Bool:
		return "bool"
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
	c.emitln("#include <stdbool.h>")
}

func (c *compiler) emitMain() {
	c.emitln("int main(void) {")
	c.emitln("init();")
	c.emitln("printf(\"%d\\n\", _main());")
	c.emitln("return 0;")
	c.emitln("}")
}

/* Main Compiler */

func (c *compiler) compObject(o ir.Object) string {
	var str string
	switch t := o.(type) {
	case *ir.Assignment:
		c.compAssignment(t)
	case *ir.Constant:
		str = c.compConstant(t)
	case *ir.Binary:
		str = c.compBinary(t)
	case *ir.Call:
		str = c.compCall(t)
	//case *ir.Define:
	//c.compDefine(t)
	case *ir.Function:
		c.compFunction(t)
	case *ir.If:
		str = c.compIf(t)
	case *ir.Unary:
		str = c.compUnary(t)
	case *ir.Var:
		str = c.compVar(t)
	case *ir.Variable:
		c.emit("%s _v%d = 0;\n", cType(t.Type()), t.ID())
		str = c.compVariable(t)
	}
	return str
}

func (c *compiler) compAssignment(a *ir.Assignment) {
	o := a.Scope().Lookup(a.Lhs)
	c.emit("_v%d = %s;\n", o.(ir.IDer).ID(), c.compObject(a.Rhs))
}

func (c *compiler) compBinary(b *ir.Binary) string {
	c.emit("%s _v%d = %s %s %s;\n", cType(b.Type()), b.ID(),
		c.compObject(b.Lhs), b.Op.String(), c.compObject(b.Rhs))
	return fmt.Sprintf("_v%d", b.ID())
}

func (c *compiler) compCall(call *ir.Call) string {
	args := make([]string, len(call.Args))
	for i, a := range call.Args {
		args[i] = fmt.Sprintf("%s", c.compObject(a))
	}
	return fmt.Sprintf("_%s(%s)", call.Name(), strings.Join(args, ","))
}

func (c *compiler) compConstant(con *ir.Constant) string {
	return con.String()
}

func (c *compiler) compDefine(d *ir.Define) {
	c.emit("%s {\n", c.compSignature(d))
	c.compFunction(d.Body.(*ir.Function))
}

func (c *compiler) compFunction(f *ir.Function) {
	for _, e := range f.Body[:len(f.Body)-1] {
		c.compObject(e)
	}
	c.emit("return %s;\n}\n", c.compObject(f.Body[len(f.Body)-1]))
}

func (c *compiler) compIdent(i *ir.Var) string {
	return fmt.Sprintf("_v%d", i.Scope().Lookup(i.Name()).(ir.IDer).ID())
}

func (c *compiler) compIf(i *ir.If) string {
	c.emit("%s _v%d = 0;\n", cType(i.Type()), i.ID())
	c.emit("if (%s) {\n", c.compObject(i.Cond))
	c.emit("_v%d = %s;\n", i.ID(), c.compObject(i.Then))
	if i.Else != nil {
		c.emitln("} else {")
		c.emit("_v%d = %s;\n", i.ID(), c.compObject(i.Else))
	}
	c.emitln("}")
	return fmt.Sprintf("_v%d", i.ID())
}

func (c *compiler) compInits(inits []*ir.Define) {
	c.emitln("void init(void) {")
	for _, d := range inits {
		c.compVariable(d.Body.(*ir.Variable))
	}
	c.emitln("}")
}

func (c *compiler) compPackage(p *ir.Package) {
	names := p.Scope().Names()
	inits := make([]*ir.Define, 0)
	for _, name := range names {
		if d, ok := p.Scope().Lookup(name).(*ir.Define); ok {
			switch t := d.Body.(type) {
			case *ir.Function:
				c.emit("%s;\n", c.compSignature(d))
				defer c.compDefine(d)
			case *ir.Variable:
				c.emit("%s _v%d = 0; // Name: %s\n", cType(t.Type()), t.ID(), d.Name())
				inits = append(inits, d)
			}
		}
	}

	c.compInits(inits)
}

func (c *compiler) compSignature(d *ir.Define) string {
	f := d.Body.(*ir.Function)
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		param := f.Scope().Lookup(p).(*ir.Param)
		params[i] = fmt.Sprintf("%s _v%d", cType(param.Type()), param.ID())
	}
	return fmt.Sprintf("%s _%s(%s)", cType(f.Type()), d.Name(),
		strings.Join(params, ","))
}

func (c *compiler) compUnary(u *ir.Unary) string {
	return fmt.Sprintf("%s%s", u.Op, c.compObject(u.Rhs))
}

func (c *compiler) compVar(v *ir.Var) string {
	var o ir.Object
	switch t := v.Scope().Lookup(v.Name()).(type) {
	case *ir.Define:
		o = t.Body
	case *ir.Param:
		o = t
	}
	return fmt.Sprintf("_v%d", o.(ir.IDer).ID())
}

func (c *compiler) compVariable(v *ir.Variable) string {
	for _, p := range v.Params {
		param := v.Scope().Lookup(p).(*ir.Param)
		c.emit("%s _v%d = 0; // name: %s\n", cType(param.Type()), param.ID(),
			param.Name())
	}
	for _, e := range v.Body[:len(v.Body)-1] {
		c.compObject(e)
	}
	c.emit("_v%d = %s;\n", v.ID(), c.compObject(v.Body[len(v.Body)-1]))
	return fmt.Sprintf("_v%d", v.ID())
}
