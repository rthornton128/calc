// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

// Windows AMD64 calling convention

package cgen

func (c *Amd64) emitMain() {
	c.emit("main:")
	c.emit("call _main")
	c.emit("cltq") // promote %eax to %rax
	c.emit("movq %rax, %rdx")
	c.emit("movq $fmt, %rcx")
	c.emit("callq printf")
	c.emit("xorq %rcx, %rcx")
	c.emit("callq ExitProcess")
	c.emit()
}
