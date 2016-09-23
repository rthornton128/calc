// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

// Package win32 implements the C calling convention for 32 bit Windows
// platforms
package win32

import (
	"io"

	"github.com/rthornton128/calc/cgen"
	"github.com/rthornton128/calc/ir"
)

type Win32Gen struct{}

func (gen *Win32Gen) CGen(w io.Writer, pkg *ir.Package) {
	e := &cgen.Writer{w}

	e.Emit(".data # .section rodata")
	e.Emit("fmt: .asciz \"%%d\\12\"")
	e.Emit()

	e.Emit(".text")
	e.Emit("global main")

	// generate sources
	c := cgen.X86{Emitter: e, Arch: new(cgen.Arch32)}
	c.CGen(e, pkg)

	e.Emit("main:")
	// prologue
	e.Emit("pushl %ebp")
	e.Emit("movl %esp, %ebp")
	e.Emit("subl $8, %esp")

	// call main
	e.Emit("call _main")

	// display result
	e.Emit("movl %eax, 4(%esp)")
	e.Emit("movl $fmt, (%esp)")
	e.Emit("call _printf")

	// exit
	e.Emit("movl $0, (%esp)")
	e.Emit("call _exit")

	// epilogue
	e.Emit("leave")
	e.Emit("ret")
	e.Emit()
}
