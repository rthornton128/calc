// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

// Windows AMD64 calling convention

package cgen

func (c *x86) emitMain() {
	c.emitln(".text")
	c.emitln("main:")
	c.emitln("push %rbp")
	c.emitln("movq %rsp, %rbp")
	c.emitln("subq $32, %rsp")
	c.emitln("call _main")
	c.emitln("movq %rax, %rdx")
	c.emitln("movq $fmt, %rcx")
	c.emitln("call printf")
	c.emitln("movq $0, %rax")
	c.emitln("addq $32, %rsp")
	c.emitln("popq %rbp")
	c.emitln("ret")
}
