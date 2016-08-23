// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

import (
	"io"

	"github.com/rthornton128/calc/ir"
)

// Registers and instructions specific to AMD64
const (
	A     = RAX
	BP    = RBP
	C     = RCX
	D     = RDX
	DI    = RDI
	SI    = RSI
	SP    = RSP
	ADD   = "addq"
	AND   = "andq"
	CMP   = "cmpq"
	DIV   = "divq"
	POP   = "popq"
	PUSH  = "pushq"
	MOV   = "movq"
	MOVZB = "movzbq"
	MUL   = "mulq"
	SUB   = "subq"
)

func (c *X86) CGen(w io.Writer, pkg *ir.Package) {
	c.Writer = w
	//c.emit(".file %s\n", "xxx.calc")
	c.emit(".global main")
	for _, name := range pkg.Scope().Names() {
		if d, ok := pkg.Scope().Lookup(name).(*ir.Define); ok {
			if f, ok := d.Body.(*ir.Function); ok {
				c.emit(".global _%s\n", name)
				defer func(name string) {
					c.emit("") // add a space for clarity
					c.emit("_%s:", name)
					c.genObject(f)
				}(name)
			}
		}
	}
	c.emit(".data")
	c.emit("fmt: .asciz \"%%d\\12\"")
	c.emit("")
	c.emit(".text")
	c.emitMain()
}
