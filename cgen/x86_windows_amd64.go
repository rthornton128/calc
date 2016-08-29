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
	c.emit("push %ebp")
	c.emit("movl %esp, %ebp")
	c.emit("subl $32, %esp")
	c.emit("call _main")
	c.emit("movl %eax, %edx")
	c.emit("movl $fmt, %ecx")
	c.emit("call _printf")
	c.emit("mov $0, %eax")
	c.emit("movl %ebp, %esp")
	c.emit("pop %ebp")
	c.emit("ret")
}
