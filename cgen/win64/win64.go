package win64

import (
	"io"

	"github.com/rthornton128/calc/cgen"
	"github.com/rthornton128/calc/ir"
)

type Win64Gen struct{}

func (gen *Win64Gen) CGen(w io.Writer, pkg *ir.Package) {
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
	e.Emit("subq $32, %rsp") // 32 bytes for shadow storage

	// call main
	e.Emit("callq _main")

	// display result
	e.Emit("movq %eax, %rdx")
	e.Emit("movq $fmt, %rcx)")
	e.Emit("callq _printf")

	// exit
	e.Emit("movq $0, %rcx")
	e.Emit("callq _exit")

	// epilogue
	e.Emit("leave")
	e.Emit("retq")
	e.Emit()
}
