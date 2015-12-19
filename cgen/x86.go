// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

import (
	"github.com/rthornton128/calc/ir"
	"github.com/rthornton128/calc/token"
)

// This is a redimentary, unoptimized x86 assembly code generator. It is
// highly unstable and a work in progress

type x86 struct{ compiler }

func (c *x86) genObject(o ir.Object) {
	switch t := o.(type) {
	case *ir.Binary:
		c.genObject(t.Lhs)
		c.emitln("push %eax")
		c.genObject(t.Rhs)
		c.emitln("movl %eax, %ecx")
		c.emitln("pop %eax")
		switch t.Op {
		case token.ADD:
			c.emitln("addl %ecx, %eax")
		case token.SUB:
			c.emitln("subl %ecx, %eax")
		case token.MUL:
			c.emitln("mull %ecx") // signed only right now
		case token.QUO:
			c.emitln("mov $0, %edx") // avoid sigfpe
			c.emitln("divl %ecx")
		case token.REM:
			c.emitln("mov $0, %edx") // avoid sigfpe
			c.emitln("divl %ecx")
			c.emitln("movl %edx, %eax")
		case token.EQL:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("je %ebx")
		case token.NEQ:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("jz %ebx")
		case token.LST:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("jl %ebx")
		case token.LTE:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("jle %ebx")
		case token.GTT:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("jg %ebx")
		case token.GTE:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("jge %ebx")
		}
	case *ir.Call:
		for _, e := range t.Args {
			c.genObject(e)
			c.emitln("push %eax") // would prefer to place into specific offset
		}
		c.emit("call %s\n", t.Name())
	case *ir.Constant:
		c.emit("movl $%s, %%eax\n", t.String())
	case *ir.If:
		c.emit("movl then%d, %%ebx\n", t.ID())
		c.genObject(t.Cond)
		if t.Else != nil {
			c.genObject(t.Else)
		}
		c.emit("then%d:\n", t.ID())
		c.genObject(t.Then)
	case *ir.Function:
		for _, e := range t.Body {
			c.genObject(e)
		}
		c.emitln("ret")
	case *ir.Unary:
		c.genObject(t.Rhs)
		c.emitln("neg %eax")
	case *ir.Var:
	}
}

func (c *x86) genPackage(pkg *ir.Package) {
	//c.emit(".file %s\n", "xxx.calc")
	c.emitln(".global _main")
	for _, name := range pkg.Scope().Names() {
		if d, ok := pkg.Scope().Lookup(name).(*ir.Define); ok {
			if f, ok := d.Body.(*ir.Function); ok {
				c.emit(".global %s\n", name)
				defer func(name string) {
					c.emit("%s:\n", name)
					c.genObject(f)
				}(name)
			}
		}
	}
	c.emitln(".data")
	c.emitln("fmt: .asciz \"%d\\12\"")
	c.emitln()
	c.emitln(".text")
	c.emitln("_main:")
	c.emitln("push %ebp")
	c.emitln("movl %esp, %ebp")
	c.emitln("subl $16, %esp")
	c.emitln("call main")
	c.emitln("movl %eax, 4(%esp)")
	c.emitln("movl $fmt, (%esp)")
	c.emitln("call _printf")
	c.emitln("movl $0, %eax")
	c.emitln("leave")
	c.emitln("ret")
}
