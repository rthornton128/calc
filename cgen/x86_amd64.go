// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

var byteOffset = 8

const (
	BP string = "%rbp"
	SP        = "%rsp"
)

func (c *X86) genEnter(sz int) {
	c.emit("pushq %rbp")
	c.emit("movq %rsp, %rbp")
	c.emitf("subq $%d, %%rsp", sz+8)
}

func (c *X86) genLeave() {
	c.emit("movq %rbp, %rsp")
	c.emit("popq %rbp")
}
