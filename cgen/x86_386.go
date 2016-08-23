// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

import "github.com/rthornton128/calc/ir"

// This is a rudimentary, unoptimized x86 assembly code generator. It is
// highly unstable and a work in progress
// BUG functions don't create a stack frame
// BUG calls don't follow cdecl convension

// Registers and instructions specific to AMD64
const (
	A     = EAX
	BP    = EBP
	C     = ECX
	D     = EBP
	SP    = ESP
	ADD   = "addl"
	AND   = "andl"
	CMP   = "cmpl"
	DIV   = "divl"
	POP   = "popl"
	PUSH  = "pushl"
	MOV   = "movl"
	MOVZB = "movzbl"
	MUL   = "mull"
	SUB   = "subl"
)

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
	c.emit(".data")
	c.emit("fmt: .asciz \"%%d\\12\"")
	c.emit()
	c.emit(".text")
	c.emit("_main:")
	c.emit("push %ebp")
	c.emit("movl %esp, %ebp")
	c.emit("subl $32, %esp")
	c.emit("call main")
	c.emit("movl %eax, 4(%esp)")
	c.emit("movl $fmt, 0(%esp)")
	c.emit("call printf")
	c.emit("movl $0, %eax")
	c.emit("leave")
	c.emit("ret")
}
