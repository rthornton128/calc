// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

// 386 only code

var byteOffset = 4

const (
	BP string = "%ebp"
	SP        = "%esp"
)

func (c *X86) genEnter(sz int) {
	c.emit("pushl %ebp")
	c.emit("movl %esp, %ebp")
	c.emit("andl $-16, %esp")
	c.emitf("subl $%d, %%esp", sz+4)
}

func (c *X86) genLeave() {
	c.emit("movl %ebp, %esp")
	c.emit("popl %ebp")
}

func (c *X86) emitMain() {
	c.emit("_main:")
	c.genEnter(16)
	c.emit("call main")
	c.emit("movl %eax, 4(%esp)")
	c.emit("movl $fmt, 0(%esp)")
	c.emit("call printf")
	c.emit("movl $0, %eax")
	c.genLeave()
	c.emit("ret")
}
