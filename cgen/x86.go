// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

import (
	"github.com/rthornton128/calc/ir"
	"github.com/rthornton128/calc/token"
)

type X86 struct {
	Emitter
	Arch
	a *allocator
}

func (c *X86) EmitPrologue(sz int) {
	c.Emit(c.Instruction(PUSH), c.Register(BP))
	c.Emitf("%s %s, %s", c.Instruction(MOV), c.Register(SP), c.Register(BP))
	c.Emitf("%s $%d, %s", c.Instruction(SUB), sz, c.Register(SP))
}

func (c *X86) EmitEpilogue(sz int) {
	c.Emitf("%s $%d, %s", c.Instruction(ADD), sz, c.Register(SP))
	c.Emit(c.Instruction(POP), c.Register(BP))
}

func (c *X86) CGen(e Emitter, pkg *ir.Package) {

	// set stack offsets and function stack sizes
	c.a = StackAlloc(pkg, c.Arch)

	for _, name := range pkg.Scope().Names() {
		//fmt.Println("x86: gen:", name)
		if f, ok := pkg.Scope().Lookup(name).(*ir.Function); ok {
			//fmt.Println("gen fun")
			//if f, ok := d.Body.(*ir.Function); ok {
			c.Emitf(".global _%s", name)
			defer func(name string) {
				c.Emitf("# %s @function, locals: %x, params: %x", name,
					f.SizeLocals, f.SizeMaxArgs)
				szStack := align16(f.SizeLocals + f.SizeMaxArgs)
				c.Emitf("_%s:", name)
				c.EmitPrologue(szStack)

				c.genObject(f, false, "%eax")

				c.EmitEpilogue(szStack)
				c.Emit("ret")
				c.Emit()
			}(name)
			//}
		}
	}
	c.Emit()
}

func (c *X86) genObject(o ir.Object, jmp bool, dest string) {
	switch t := o.(type) {
	case *ir.Assignment:
		// TODO optimize to allow constant/variable to be directly moved into
		// location
		c.genObject(t.Rhs, false, dest)
		c.Emitf("movl %s, %s", dest, t.Scope().Lookup(t.Lhs).Loc())
	case *ir.Binary:
		c.genBinary(t, jmp, dest)
	case *ir.Call:
		for i, arg := range t.Args {
			//fmt.Println("call: processing argument", i, "-", arg)
			switch arg.(type) {
			case *ir.Constant:
				c.genObject(arg, false, c.a.ArgumentLoc(i))
			default:
				c.genObject(arg, false, dest) //"%eax")
				c.Emitf("movl %s, %s", dest, c.a.ArgumentLoc(i))
			}
		}
		c.Emitf("call _%s", t.Name())
	case *ir.Constant:
		var val string
		switch t.String() {
		case "true":
			val = "1"
		case "false":
			val = "0"
		default:
			val = t.String()
		}
		c.Emitf("movl $%s, %s", val, dest)
	case *ir.For:
		c.Emitf("jmp %s", t.CondLabel())
		c.Emitf("%s:", t.BodyLabel())
		for _, e := range t.Body {
			c.genObject(e, false, "%eax")
		}
		c.Emitf("%s:", t.CondLabel())
		c.genObject(t.Cond, true, t.BodyLabel())
	case *ir.If:
		c.genIf(t)
	case *ir.Function:
		if len(t.Body) == 0 {
			c.Emit("xorl %eax, %eax")
			return
		}
		for _, e := range t.Body {
			c.genObject(e, false, "%eax")
		}
	case *ir.Unary:
		c.genObject(t.Rhs, false, dest)
		c.Emitf("neg %s", dest)
	case *ir.Var:
		//fmt.Println("looking up var:", t.Name(), t.ID(), t.Loc())
		v := t.Scope().Lookup(t.Name()) //fmt.Sprint(t.Name(), t.ID()))
		//fmt.Println("found:", v, v.ID(), "Loc:", v.Loc())
		if o, ok := v.(*ir.Variable); ok {
			c.genObject(o, false, o.Loc())
			return
		}
		v = t.Scope().Lookup(t.Name())
		c.Emitf("movl %s, %s", v.Loc(), dest)
	case *ir.Variable:
		c.Emitf("# variable %d start", t.ID())
		for _, p := range t.Params {
			c.Emitf("movl $0, %s", p.Loc())
		}
		for _, e := range t.Body {
			//fmt.Println("variable: processing body expr", i, "-", e)
			c.genObject(e, false, "%eax")
		}
		c.Emitf("movl %%eax, %s", t.Loc())
		c.Emitf("# variable %d end", t.ID())
	}
}

func (c *X86) genBinary(b *ir.Binary, jump bool, dest string) {
	c.genObject(b.Lhs, false, "%eax")

	switch b.Rhs.(type) {
	case *ir.Constant, *ir.Var:
		if b.Op == token.QUO || b.Op == token.REM {
			c.genObject(b.Rhs, false, "%ecx")
		} else {
			c.genObject(b.Rhs, false, "%edx")
		}
	default:
		c.Emitf("movl %%eax, %s", b.Rhs.Loc())
		c.genObject(b.Rhs, false, "%eax")
		if b.Op == token.QUO || b.Op == token.REM {
			c.Emit("movl %eax, %ecx")
		} else {
			c.Emit("movl %eax, %edx")
		}
		c.Emitf("movl %s, %%eax", b.Rhs.Loc())
	}
	switch b.Op {
	case token.ADD:
		c.Emit("addl %edx, %eax")
	case token.SUB:
		c.Emit("subl %edx, %eax")
	case token.MUL:
		c.Emit("imul %edx, %eax")
	case token.QUO:
		c.Emit("cdq")
		c.Emit("idiv %ecx")
	case token.REM:
		c.Emit("cdq")
		c.Emit("idiv %ecx, %eax")
		c.Emit("movl %edx, %eax")
	default:
		c.Emit("cmpl %edx, %eax")
		if jump {
			c.genJump(b, dest)
			return
		}
		switch b.Op {
		case token.EQL:
			c.Emit("sete %al")
		case token.NEQ:
			c.Emit("setne %al")
		case token.LST:
			c.Emit("setl %al")
		case token.LTE:
			c.Emit("setle %al")
		case token.GTT:
			c.Emit("setg %al")
		case token.GTE:
			c.Emit("setge %al")
		}
		c.Emit("movzbl %al, %eax")
	}
}

func (c *X86) genJump(b *ir.Binary, label string) {
	var inst string
	switch b.Op {
	case token.EQL:
		inst = "je"
	case token.NEQ:
		inst = "jne"
	case token.LST:
		inst = "jl"
	case token.LTE:
		inst = "jle"
	case token.GTT:
		inst = "jg"
	case token.GTE:
		inst = "jge"
	}
	c.Emitf("%s %s", inst, label)
}

func (c *X86) genIf(i *ir.If) {
	//fmt.Println("if:", i.ID(), i.ThenLabel, i.EndLabel)
	switch t := i.Cond.(type) {
	case *ir.Binary:
		c.genBinary(t, true, i.ThenLabel)
	default:
		c.genObject(t, false, "%eax") // s/b genConstant() or genBoolean()
		c.Emit("cmpl $0, %eax")
		c.Emitf("jnz %s", i.ThenLabel)
	}
	c.genObject(i.Else, false, "%eax")
	c.Emitf("jmp %s", i.EndLabel)
	c.Emitf("%s:", i.ThenLabel)
	c.genObject(i.Then, false, "%eax")
	c.Emitf("%s:", i.EndLabel)
}
