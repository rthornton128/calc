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
	c.emit1reg(PUSH, BP)
	c.emit2reg(MOV, SP, BP)
	c.emitimm(SUB, "32", SP)
	c.emit("call _main")
	c.emit2reg(MOV, A, D)
	c.emitimm(MOV, "fmt", C)
	c.emit("call printf")
	c.emitimm(MOV, "0", A)
	c.emitimm(ADD, "32", SP)
	c.emit1reg(POP, BP)
	c.emit(RET)
}
