// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

// cdecl and stdcall code should be here

type Arch32 struct{}

var registers32 = []string{
	AX: "%eax",
	BP: "%ebp",
	SP: "%esp",
}

var instructions32 = []string{
	ADD:  "addl",
	MOV:  "movl",
	POP:  "popl",
	PUSH: "pushl",
	SUB:  "subl",
}

func (a *Arch32) Instruction(i Instruction) string { return instructions32[i] }
func (a *Arch32) Register(r Register) string       { return registers32[r] }
func (a *Arch32) Width() int                       { return 4 }
