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

	"github.com/rthornton128/calc/ir"
	"github.com/rthornton128/calc/token"
)

// This is a rudimentary, unoptimized x86 assembly code generator. It is
// highly unstable and a work in progress
// BUG functions don't create a stack frame
// BUG calls don't follow cdecl convension

type X86 struct{ io.Writer }

func (c *X86) emit(f string, args ...interface{}) {
	fmt.Fprintf(c.Writer, f, args...)
}

func (c *X86) emitln(args ...interface{}) {
	fmt.Fprintln(c.Writer, args...)
}

func (c *X86) genObject(o ir.Object) {
	switch t := o.(type) {
	case *ir.Binary:
		c.genObject(t.Lhs)
		c.emitln("pushl %eax")
		c.genObject(t.Rhs)
		c.emitln("movl %eax, %ecx")
		c.emitln("popl %eax")
		switch t.Op {
		case token.ADD:
			c.emitln("addl %ecx, %eax")
		case token.SUB:
			c.emitln("subl %ecx, %eax")
		case token.MUL:
			c.emitln("mull %ecx") // signed only right now
		case token.QUO:
			c.emitln("movl $0, %edx") // avoid sigfpe
			c.emitln("divl %ecx")
		case token.REM:
			c.emitln("movl $0, %edx") // avoid sigfpe
			c.emitln("divl %ecx")
			c.emitln("movl %edx, %eax")
		case token.EQL:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("sete %al")
			c.emitln("movzbl %al, %eax")
		case token.NEQ:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("setne %al")
			c.emitln("movzbl %al, %eax")
		case token.LST:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("setl %al")
			c.emitln("movzbl %al, %eax")
		case token.LTE:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("setle %al")
			c.emitln("movzbl %al, %eax")
		case token.GTT:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("setg %al")
			c.emitln("movzbl %al, %eax")
		case token.GTE:
			c.emitln("cmpl %ecx, %eax")
			c.emitln("setge %al")
			c.emitln("movzbl %al, %eax")
		}
	case *ir.Call:
		// BUG this is not proper calling convension for 386, must be pushed in
		// reverse order
		for _, e := range t.Args {
			c.genObject(e)
			c.emitln("pushl %eax")
		}
		c.emit("call %s\n", t.Name())
	case *ir.Constant:
		switch t.String() {
		case "true":
			c.emitln("movl $1, %eax")
		case "false":
			c.emitln("movl $0, %eax")
		default:
			c.emit("movl $%s, %%eax\n", t.String())
		}
	case *ir.For:
		c.emit("jmp L%d\n", t.ID())
		c.emit("L%db:\n", t.ID())
		for _, e := range t.Body {
			c.genObject(e)
		}
		c.emit("L%d:\n", t.ID())
		c.genObject(t.Cond)
		c.emitln("andl $1, %eax")
		c.emit("jnz L%db\n", t.ID())
	case *ir.If:
		c.genObject(t.Cond)
		c.emitln("andl $1, %eax")
		if t.Else != nil {
			c.emit("jz L%de\n", t.ID())
			c.genObject(t.Then)
			c.emit("jmp L%d\n", t.ID())
			c.emit("L%de:\n", t.ID())
			c.genObject(t.Else)
		} else {
			c.emit("jz L%d\n", t.ID())
			c.genObject(t.Then)
		}
		c.emit("L%d:\n", t.ID())
	case *ir.Function:
		for _, e := range t.Body {
			c.genObject(e)
		}
		c.emitln("ret")
	case *ir.Unary:
		c.genObject(t.Rhs)
		c.emitln("neg %eax")
	case *ir.Var:
	case *ir.Variable:
		//for _, p := range t.Args {

		//}
		for _, e := range t.Body {
			c.genObject(e)
		}
	}
}

func (c *X86) CGen(p *ir.Package) {
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
	c.emitln("push %rbp")
	c.emitln("movl %rsp, %rbp")
	c.emitln("subl $32, %rsp")
	c.emitln("call main")
	c.emitln("movl %eax, 4(%esp)")
	c.emitln("movl $fmt, 0(%esp)")
	c.emitln("call printf")
	c.emitln("movl $0, %eax")
	c.emitln("leave")
	c.emitln("ret")
}
