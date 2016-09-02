// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

// Windows AMD64 calling convention

package cgen

func (c *X86) emitMain() {
	c.emit("main:")
	c.genEnter(24)
	c.emit("call _main")
	c.emit("cltq") // promite %eax to %rax
	c.emit("movq %rax, %rdx")
	c.emit("movq $fmt, %rcx")
	c.emit("call printf")
	c.emit("movq $0, %rax")
	c.genLeave()
	c.emit("ret")
}
