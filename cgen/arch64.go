// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

type Arch64 struct{}

var registers64 = []string{
	AX: "%rax",
	BP: "%rbp",
	SP: "%rsp",
}

var instructions64 = []string{
	ADD:  "addq",
	MOV:  "movq",
	POP:  "popq",
	PUSH: "pushq",
	SUB:  "subq",
}

func (a *Arch64) Instruction(i Instruction) string { return instructions64[i] }
func (a *Arch64) Register(r Register) string       { return registers64[r] }
func (s *Arch64) Width() int                       { return 8 }
