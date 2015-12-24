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
	case *ir.Assignment:
		c.genObject(t.Rhs)
		c.emitln("movq %rax, %rbx") // BUG temporary
	case *ir.Binary:
		c.genObject(t.Lhs)
		c.emitln("pushq %rax")
		c.genObject(t.Rhs)
		c.emitln("movq %rax, %rcx")
		c.emitln("popq %rax")
		switch t.Op {
		case token.ADD:
			c.emitln("addq %rcx, %rax")
		case token.SUB:
			c.emitln("subq %rcx, %rax")
		case token.MUL:
			c.emitln("mulq %rcx") // signed only right now
		case token.QUO:
			c.emitln("movq $0, %rdx") // avoid sigfpe
			c.emitln("divq %rcx")
		case token.REM:
			c.emitln("movq $0, %rdx") // avoid sigfpe
			c.emitln("divq %rcx")
			c.emitln("movq %rdx, %rax")
		case token.EQL:
			c.emitln("cmpq %rcx, %rax")
			c.emitln("sete %al")
			c.emitln("movzbq %al, %rax")
		case token.NEQ:
			c.emitln("cmpq %rcx, %rax")
			c.emitln("setne %al")
			c.emitln("movzbq %al, %rax")
		case token.LST:
			c.emitln("cmpq %rcx, %rax")
			c.emitln("setl %al")
			c.emitln("movzbq %al, %rax")
		case token.LTE:
			c.emitln("cmpq %rcx, %rax")
			c.emitln("setle %al")
			c.emitln("movzbq %al, %rax")
		case token.GTT:
			c.emitln("cmpq %rcx, %rax")
			c.emitln("setg %al")
			c.emitln("movzbq %al, %rax")
		case token.GTE:
			c.emitln("cmpq %rcx, %rax")
			c.emitln("setge %al")
			c.emitln("movzbq %al, %rax")
		}
	case *ir.Call:
		for _, e := range t.Args {
			c.genObject(e)
			c.emitln("pushq %rax") // would prefer to place into specific offset
		}
		c.emit("call %s\n", t.Name())
	case *ir.Constant:
		switch t.String() {
		case "true":
			c.emitln("movq $1, %rax")
		case "false":
			c.emitln("movq $0, %rax")
		default:
			c.emit("movq $%s, %%rax\n", t.String())
		}
	case *ir.For:
		c.emit("jmp L%d\n", t.ID())
		c.emit("L%db:\n", t.ID())
		for _, e := range t.Body {
			c.genObject(e)
		}
		c.emit("L%d:\n", t.ID())
		c.genObject(t.Cond)
		c.emitln("cmpq $1, %rax")
		c.emit("jnz L%db\n", t.ID())
	case *ir.If:
		c.genObject(t.Cond)
		c.emitln("cmpq $1, %rax")
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
		c.emitln("neg %rax")
	case *ir.Var:
		c.emitln("movq %rbx, %rax") // BUG temporary
	case *ir.Variable:
		for _, e := range t.Body {
			c.genObject(e)
		}
	}
}

func (c *x86) genPackage(pkg *ir.Package) {
	//c.emit(".file %s\n", "xxx.calc")
	c.emitln(".global main")
	for _, name := range pkg.Scope().Names() {
		if d, ok := pkg.Scope().Lookup(name).(*ir.Define); ok {
			if f, ok := d.Body.(*ir.Function); ok {
				c.emit(".global _%s\n", name)
				defer func(name string) {
					c.emit("_%s:\n", name)
					c.genObject(f)
				}(name)
			}
		}
	}
	c.emitln(".data")
	c.emitln("fmt: .asciz \"%d\\12\"")
	c.emitln()
	c.emitln(".text")
	c.emitMain()
}
