// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

// Package pos64 is the POSIX interface for x86-64 on Linux, BSD and MacOS
package pos64

import (
	"io"

	"github.com/rthornton128/calc/cgen"
	"github.com/rthornton128/calc/ir"
)

// Pos64 is the Posix (SysV) 64 bit wrapper for x86-64
type Pos64 struct{}

func (gen *Pos64) CGen(w io.Writer, pkg *ir.Package) {
	e := &cgen.Writer{w}

	e.Emit(".data # .section rodata")
	e.Emit("fmt: .asciz \"%%d\\12\"")
	e.Emit()

	e.Emit(".text")
	e.Emit("global main")

	// generate sources
	c := cgen.X86{Emitter: e, Arch: new(cgen.Arch64)}
	c.CGen(e, pkg)

	e.Emit("main:")
	// prologue
	e.Emit("pushq %rbp")
	e.Emit("movq %rsp, %rbp")

	// call main
	e.Emit("callq _main")

	// display result
	e.Emit("movq %eax, %rsi")
	e.Emit("movq $fmt, %rdi)")
	e.Emit("callq printf")

	// exit
	e.Emit("movq $0, %rdi")
	e.Emit("callq exit")

	// epilogue
	e.Emit("leave")
	e.Emit("retq")
	e.Emit()
}
