// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

import (
	"fmt"
	"io"

	"github.com/rthornton128/calc/ir"
	"github.com/rthornton128/calc/token"
)

// Registers and instructions shared across 386 and AMD64
const (
	// registers
	AL  = "al"
	EAX = "eax"
	RAX = "rax"
	ECX = "ecx"
	RCX = "rcx"
	EDX = "edx"
	RDX = "rdx"
	EDI = "edi"
	RDI = "rdi"
	ESI = "esi"
	RSI = "rsi"
	EBP = "ebp"
	RBP = "rbp"
	ESP = "esp"
	RSP = "rsp"

	// jump instructions
	CALL  = "call"
	JMP   = "jmp"
	JZ    = "jz"
	JNZ   = "jnz"
	LEAVE = "leave"
	RET   = "ret"

	// inversion
	NEG = "neg"

	// set instructions
	SETE  = "sete"
	SETNE = "setne"
	SETL  = "setl"
	SETLE = "setle"
	SETG  = "setg"
	SETGE = "setge"
)

// This is a rudimentary, unoptimized x86 assembly code generator. It is
// highly unstable and a work in progress
// BUG functions don't create a stack frame
// BUG calls don't follow cdecl convension

type X86 struct{ io.Writer }

func (c *X86) emitimm(op, imm, reg0 string) {
	fmt.Fprintf(c.Writer, "%s $%s, %%%s\n", op, imm, reg0)
}

func (c *X86) emit1reg(op, reg0 string) {
	fmt.Fprintf(c.Writer, "%s %%%s\n", op, reg0)
}

func (c *X86) emit2reg(op, reg0, reg1 string) {
	fmt.Fprintf(c.Writer, "%s %%%s, %%%s\n", op, reg0, reg1)
}

func (c *X86) emitlbl(f string, arg interface{}) {
	fmt.Fprintf(c.Writer, f+":\n", arg)
}

func (c *X86) emitjmp(f, op string, arg interface{}) {
	fmt.Fprintf(c.Writer, "%s "+f+"\n", op, arg)
}

func (c *X86) emit(f string, args ...interface{}) {
	fmt.Fprintf(c.Writer, f+"\n", args...)
}

func (c *X86) genObject(o ir.Object) {
	switch t := o.(type) {
	case *ir.Binary:
		c.genObject(t.Lhs)
		c.emit1reg(PUSH, A)
		c.genObject(t.Rhs)
		c.emit2reg(MOV, A, C)
		c.emit1reg(POP, A)
		switch t.Op {
		case token.ADD:
			c.emit2reg(ADD, C, A)
		case token.SUB:
			c.emit2reg(SUB, C, A)
		case token.MUL:
			c.emit1reg(MUL, C) // signed only right now
		case token.QUO:
			c.emitimm(MOV, "0", D) // avoid sigfpe
			c.emit1reg(DIV, C)
		case token.REM:
			c.emitimm(MOV, "0", D) // avoid sigfpe
			c.emit1reg(DIV, C)
			c.emit2reg(MOV, D, A)
		case token.EQL:
			c.emit2reg(CMP, C, A)
			c.emit(SETE, AL)
			c.emit2reg(MOVZB, AL, A)
		case token.NEQ:
			c.emit2reg(CMP, C, A)
			c.emit1reg(SETNE, AL)
			c.emit2reg(MOVZB, AL, A)
		case token.LST:
			c.emit2reg(CMP, C, A)
			c.emit1reg(SETL, AL)
			c.emit2reg(MOVZB, AL, A)
		case token.LTE:
			c.emit2reg(CMP, C, A)
			c.emit1reg(SETLE, AL)
			c.emit2reg(MOVZB, AL, A)
		case token.GTT:
			c.emit2reg(CMP, C, A)
			c.emit1reg(SETG, AL)
			c.emit2reg(MOVZB, AL, A)
		case token.GTE:
			c.emit2reg(CMP, C, A)
			c.emit1reg(SETGE, AL)
			c.emit2reg(MOVZB, AL, A)
		}
	case *ir.Call:
		for _, e := range t.Args {
			c.genObject(e)
			c.emit1reg(PUSH, A)
		}
		c.emitjmp("%s", CALL, t.Name())
	case *ir.Constant:
		switch t.String() {
		case "true":
			c.emitimm(MOV, "1", A)
		case "false":
			c.emitimm(MOV, "0", A)
		default:
			c.emitimm(MOV, t.String(), A)
		}
	case *ir.For:
		c.emitjmp("L%d", JMP, t.ID())
		c.emitlbl("L%db", t.ID())
		for _, e := range t.Body {
			c.genObject(e)
		}
		c.emitlbl("L%d", t.ID())
		c.genObject(t.Cond)
		c.emit2reg(AND, "1", A)
		c.emitjmp("L%db", JNZ, t.ID())
	case *ir.If:
		c.genObject(t.Cond)
		c.emitimm(AND, "1", A)
		if t.Else != nil {
			c.emitjmp("L%de", JZ, t.ID())
			c.genObject(t.Then)
			c.emitjmp("L%d", JMP, t.ID())
			c.emitlbl("L%de", t.ID())
			c.genObject(t.Else)
		} else {
			c.emitjmp("L%d", JZ, t.ID())
			c.genObject(t.Then)
		}
		c.emitlbl("L%d", t.ID())
	case *ir.Function:
		for _, e := range t.Body {
			c.genObject(e)
		}
		c.emit(RET)
	case *ir.Unary:
		c.genObject(t.Rhs)
		c.emit1reg(NEG, A)
	case *ir.Var:
	case *ir.Variable:
		//for _, p := range t.Args {

		//}
		for _, e := range t.Body {
			c.genObject(e)
		}
	}
}
