// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

import (
	"fmt"
	"io"
	"strings"

	"github.com/rthornton128/calc/ir"
)

type StdC struct{ Emitter }

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

func (c *StdC) emitHeaders() {
	c.Emit("#include <stdio.h>")
	c.Emit("#include <stdint.h>")
	c.Emit("#include <stdbool.h>")
}

func (c *StdC) emitMain() {
	c.Emit("int main(void) {")
	c.Emit("printf(\"%d\\n\", _main());")
	c.Emit("return 0;")
	c.Emit("}")
}

/* Main Compiler */

func (c *StdC) compObject(o ir.Object) string {
	switch t := o.(type) {
	case *ir.Assignment:
		return c.compAssignment(t)
	case *ir.Constant:
		return c.compConstant(t)
	case *ir.Binary:
		return c.compBinary(t)
	case *ir.Call:
		return c.compCall(t)
	case *ir.For:
		return c.compFor(t)
	case *ir.If:
		return c.compIf(t)
	case *ir.Unary:
		return c.compUnary(t)
	case *ir.Var:
		return c.compVar(t)
	case *ir.Variable:
		return c.compVariable(t)
	}
	return ""
}

func (c *StdC) compAssignment(a *ir.Assignment) string {
	c.Emitf("%s%d = %s;", a.Lhs, a.Scope().Lookup(a.Lhs).ID(),
		c.compObject(a.Rhs))
	return a.Lhs
}

func (c *StdC) compBinary(b *ir.Binary) string {
	return fmt.Sprintf("(%s %s %s)",
		c.compObject(b.Lhs), b.Op.String(), c.compObject(b.Rhs))
}

func (c *StdC) compCall(call *ir.Call) string {
	args := make([]string, len(call.Args))
	for i, a := range call.Args {
		args[i] = fmt.Sprintf("%s", c.compObject(a))
	}
	return fmt.Sprintf("_%s(%s)", call.Name(), strings.Join(args, ","))
}

func (c *StdC) compConstant(con *ir.Constant) string {
	return con.String()
}

func (c *StdC) compFor(f *ir.For) string {
	c.Emitf("%s %s%d = 0;", cType(f.Type()), f.Name(), f.ID())
	c.Emitf("while (%s) {", c.compObject(f.Cond))
	for _, e := range f.Body[:len(f.Body)-1] {
		c.compObject(e)
	}
	c.Emitf("}\n%s%d = %s;", f.Name(), f.ID(),
		c.compObject(f.Body[len(f.Body)-1]))
	return fmt.Sprintf("%s%d", f.Name(), f.ID())
}

func (c *StdC) compFunction(f *ir.Function) {
	c.Emitf("%s {", c.compSignature(f))
	for _, e := range f.Body[:len(f.Body)-1] {
		c.compObject(e)
	}
	c.Emitf("return %s;\n}", c.compObject(f.Body[len(f.Body)-1]))
}

func (c *StdC) compIdent(i *ir.Var) string {
	fmt.Printf("ident: %s%d\n", i.Name(), i.ID())
	return fmt.Sprintf("%s%d", i.Name(), i.ID())
}

func (c *StdC) compIf(i *ir.If) string {
	c.Emitf("%s if%d = 0; /* %s */", cType(i.Type()), i.ID(), i.Name())
	c.Emitf("if (%s) {", c.compObject(i.Cond))
	c.Emitf("if%d = %s;", i.ID(), c.compObject(i.Then))
	if i.Else != nil {
		c.Emit("} else {")
		c.Emitf("if%d = %s;", i.ID(), c.compObject(i.Else))
	}
	c.Emit("}")
	return fmt.Sprintf("if%d", i.ID())
}

func (c *StdC) CGen(w io.Writer, p *ir.Package) {
	c.Emitter = &Writer{w}
	c.emitHeaders()

	for _, name := range p.Scope().Names() {
		if f, ok := p.Scope().Lookup(name).(*ir.Function); ok {
			c.Emitf("%s;", c.compSignature(f))
			params := make([]string, len(f.Params))
			for i, p := range f.Params {
				params[i] = cType(f.Scope().Lookup(p.Name()).(*ir.Param).Type())
			}
			c.Emitf("%s (*_%s)(%s) = f%d;", cType(f.Type()), name,
				strings.Join(params, ","), f.ID())
			defer c.compFunction(f)
		}
	}
	c.emitMain()
}

func (c *StdC) compSignature(f *ir.Function) string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		param := f.Scope().Lookup(p.Name()).(*ir.Param)
		params[i] = fmt.Sprintf("%s %s%d", cType(param.Type()), param.Name(),
			param.ID())
	}
	return fmt.Sprintf("%s f%d(%s)", cType(f.Type()), f.ID(),
		strings.Join(params, ","))
}

func (c *StdC) compUnary(u *ir.Unary) string {
	return fmt.Sprintf("%s%s", u.Op, c.compObject(u.Rhs))
}

func (c *StdC) compVar(v *ir.Var) string {
	switch t := v.Scope().Lookup(v.Name()).(type) {
	case *ir.Param:
		fmt.Printf("param: %s%d\n", t.Name(), t.ID())
		return fmt.Sprintf("%s%d", t.Name(), t.ID())
	}
	panic("unreachable")
}

func (c *StdC) compVariable(v *ir.Variable) string {
	c.Emitf("%s %s = 0;", cType(v.Type()), v.Name())
	for _, p := range v.Params {
		param := v.Scope().Lookup(p.Name()).(*ir.Param)
		c.Emitf("%s %s%d = 0;", cType(param.Type()), param.Name(), param.ID())
	}
	for _, e := range v.Body[:len(v.Body)-1] {
		c.compObject(e)
	}
	c.Emitf("%s = %s;", v.Name(), c.compObject(v.Body[len(v.Body)-1]))
	return fmt.Sprintf("%s", v.Name())
}
