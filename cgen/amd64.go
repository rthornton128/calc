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

type Amd64 struct {
	io.Writer
	offsets map[string]int
}

func (c *Amd64) genEnter(sz int) {
}
func (c *Amd64) emit(args ...interface{}) {
	fmt.Fprintln(c.Writer, args...)
}

func (c *Amd64) emitf(f string, args ...interface{}) {
	fmt.Fprintf(c.Writer, f+"\n", args...)
}
func (c *Amd64) CGen(w io.Writer, pkg *ir.Package) {
	c.Writer = w

	// set stack offsets and function stack sizes
	StackAlloc(pkg, 8)

	//c.emit(".file %s\n", "xxx.calc")
	c.emit(".data")
	c.emitf("fmt: .asciz \"%%d\\12\"")
	c.emit("")
	c.emit(".text")
	c.emit(".global main")
	for _, name := range pkg.Scope().Names() {
		if d, ok := pkg.Scope().Lookup(name).(*ir.Define); ok {
			if f, ok := d.Body.(*ir.Function); ok {
				c.emitf(".global _%s", name)
				defer func(name string) {
					// label
					c.emitf("_%s:", name)

					// pre
					c.emit("pushq %rbp")
					c.emit("movq %rsp, %rbp")
					c.emitf("subq $%d, %%rsp", fnStackAllocs[name].stackSz)

					// body
					c.offsets = fnStackAllocs[name].offsets
					c.genObject(f, "%eax")

					// post
					c.emit("movq %rbp, %rsp")
					c.emit("popq %rbp")
					c.emit("ret")
					c.emit()
				}(name)
			}
		}
	}
	c.emitMain()
}

func (c *Amd64) genObject(o ir.Object, dest string) {
	switch t := o.(type) {
	case *ir.Assignment:
		c.genObject(t.Rhs, "%eax")
		c.emitf("movl %%eax, %d(%%rbp)", c.offsets[t.Name()])
	case *ir.Binary:
		c.genBinary(t, "")
	case *ir.Call:
		sz := 8 * len(t.Args)
		for i, arg := range t.Args {
			c.genObject(arg, "%eax")
			c.emitf("movl %%eax, %d(%%rsp)", sz-(i*8))
		}
		c.emitf("call _%s", t.Name())
	case *ir.Constant:
		var val string
		switch t.String() {
		case "true":
			val = "1"
		case "false":
			val = "0"
		default:
			val = t.String()
		}
		c.emitf("movl $%s, %s", val, dest)
	case *ir.For:
	case *ir.If:
		c.genIf(t)
	case *ir.Function:
		// enter
		for _, e := range t.Body {
			c.genObject(e, "%eax")
		}
	case *ir.Unary:
		c.genObject(t.Rhs, "%eax")
		c.emit("neg %eax")
	case *ir.Var:
		//o := t.Scope().Lookup(t.Name())

		c.emitf("movl %d(%%rbp), %s", c.offsets[t.Name()], dest)
	case *ir.Variable:
		for _, e := range t.Body {
			c.genObject(e, "%eax")
		}
	}
}

func (c *Amd64) genBinary(b *ir.Binary, jump string) {
	c.genObject(b.Lhs, "%eax")
	// TODO an improvement (optimization) is to reduce
	// the Rhs move and instead
	switch b.Rhs.(type) {
	case *ir.Constant, *ir.Var:
		if b.Op == token.QUO || b.Op == token.REM {
			c.genObject(b.Rhs, "%ecx")
		} else {
			c.genObject(b.Rhs, "%edx")
		}
	default:
		c.emitf("movl %%eax, %d(%%rsp)", c.offsets[b.Name()])
		c.genObject(b.Rhs, "%eax")
		if b.Op == token.QUO || b.Op == token.REM {
			c.emit("movl %eax, %ecx")
		} else {
			c.emit("movl %eax, %edx")
		}
		c.emitf("movl %d(%%rsp), %%eax", c.offsets[b.Name()])
	}
	switch b.Op {
	case token.ADD:
		c.emit("addl %edx, %eax")
	case token.SUB:
		c.emit("subl %edx, %eax")
	case token.MUL:
		c.emit("imul %edx, %eax")
	case token.QUO:
		c.emit("cdq")
		c.emit("idiv %ecx")
	case token.REM:
		c.emit("cdq")
		c.emit("idiv %ecx, %eax")
		c.emit("movl %edx, %eax")
	default:
		c.emit("cmpl %edx, %eax")
		switch b.Op {
		case token.EQL:
			if len(jump) > 0 {
				c.emitf("jne %s", jump)
				return
			}
			c.emit("sete %al")
			c.emit("movzbl %al, %eax")
		case token.NEQ:
			if len(jump) > 0 {
				c.emitf("je %s", jump)
			}
			c.emit("setne %al")
			c.emit("movzbl %al, %eax")
		case token.LST:
			if len(jump) > 0 {
				c.emitf("jge %s", jump)
				return
			}
			c.emit("setl %al")
			c.emit("movzbl %al, %eax")
		case token.LTE:
			if len(jump) > 0 {
				c.emitf("jg %s", jump)
				return
			}
			c.emit("setle %al")
			c.emit("movzbl %al, %eax")
		case token.GTT:
			if len(jump) > 0 {
				c.emitf("jle %s", jump)
				return
			}
			c.emit("setg %al")
			c.emit("movzbl %al, %eax")
		case token.GTE:
			if len(jump) > 0 {
				c.emitf("jl %s", jump)
				return
			}
			c.emit("setge %al")
			c.emit("movzbl %al, %eax")
		}
	}
}

func (c *Amd64) genIf(i *ir.If) {
	switch t := i.Cond.(type) {
	case *ir.Binary:
		c.genBinary(t, i.ElseLabel)
	default:
		c.genObject(t, "%eax") // s/b genConstant() or genBoolean()
		c.emit("cmpl $0, %eax")
		c.emitf("jz %s", i.ElseLabel)
	}
	c.genObject(i.Then, "%eax")
	c.emitf("jmp %s", i.EndLabel)
	c.emitf("%s:", i.ElseLabel)
	c.genObject(i.Else, "%eax")
	c.emitf("%s:", i.EndLabel)
}
