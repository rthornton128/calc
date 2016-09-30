// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

import (
	"fmt"

	"github.com/rthornton128/calc/ir"
)

type Register int
type Instruction int

const (
	AX Register = iota // Register A
	BP                 // Base Pointer
	SP                 // Stack Pointer
)

const (
	ADD  Instruction = iota // Addition
	MOV                     // Move
	POP                     // Pop
	PUSH                    // Push
	SUB                     // Subtraction
)

type Arch interface {
	Instruction(Instruction) string
	Register(Register) string
	Width() int
}

// StackAlloc sets stack offsets for all expressions needing one. It does is a
// rudimentary register allocator that mainly just spills everything to the
// stack and does little to no optimization
func StackAlloc(pkg *ir.Package, arch Arch) *allocator {
	a := allocator{
		nextOffset: 0 - arch.Width(),
		Arch:       arch,
	}

	// assign offsets to parameters of all functions first
	for _, f := range pkg.Scope().Names() {
		if d, ok := pkg.Lookup(f).(*ir.Define); ok {
			switch t := d.Body.(type) {
			case *ir.Function:
				a.walk(t)
			}
		}
	}
	return &a
}

func align16(n int) int { return (n & -16) + 16 }

// allocator is a rudimentary register allocator that mainly just spills
// everything to the stack and does little to no optimization
type allocator struct {
	Arch
	nextOffset int
	locals     int
	maxArgs    int
}

func (a *allocator) ArgumentLoc(i int) string {
	return fmt.Sprintf("%d(%s)", i*a.Width(), a.Register(SP))
}

func (a *allocator) CallStackOffset(i int) string {
	return fmt.Sprintf("%d(%s)", (i*a.Width())+(a.Width()*2), a.Register(BP))
}

func (a *allocator) ParameterLoc(i int) string {
	return a.CallStackOffset(i)
}

func (a *allocator) nextLoc() string {
	s := fmt.Sprintf("%d(%s)", a.nextOffset, a.Register(BP))
	a.nextOffset -= a.Width()
	a.locals++
	return s
}

func (a *allocator) walk(o ir.Object) {
	switch t := o.(type) {
	case *ir.Assignment:
		a.walk(t.Rhs)
	case *ir.Binary:
		a.walk(t.Lhs)
		a.walk(t.Rhs)
		switch t.Rhs.(type) {
		case *ir.Constant, *ir.Var:
			//no nothing
		default:
			t.Rhs.SetLoc(a.nextLoc())
		}
	case *ir.Call:
		for _, arg := range t.Args {
			a.walk(arg)

		}
		if a.maxArgs < len(t.Args) {
			a.maxArgs = len(t.Args)
		}
	case *ir.Function:
		// set parameter offsets
		for i, p := range t.Params {
			p.SetLoc(a.CallStackOffset(i))
		}

		// locals
		for _, o := range t.Body {
			a.walk(o)
		}
		t.SizeLocals = a.locals * a.Width()
		t.SizeMaxArgs = a.maxArgs * a.Width()
		a.locals, a.maxArgs = 0, 0
	case *ir.For:
		a.walk(t.Cond)
		for _, e := range t.Body {
			a.walk(e)
		}
	case *ir.If:
		a.walk(t.Cond)
		a.walk(t.Then)
		if t.Else != nil {
			a.walk(t.Else)
		}
	case *ir.Var:
		if d, ok := t.Scope().Lookup(t.Name()).(*ir.Define); ok {
			switch b := d.Body.(type) {
			case *ir.Variable:
				x := b.Copy(t.Name(), t.ID())
				t.Scope().Insert(x)
				a.walk(x)
			default:
				t.SetLoc(a.nextLoc())
			}
		}
	case *ir.Variable:
		t.SetLoc(a.nextLoc())
		for _, p := range t.Params {
			loc := a.nextLoc()
			p.SetLoc(loc)
			t.Scope().Lookup(p.Name()).SetLoc(loc)
		}
		for _, o := range t.Body {
			a.walk(o)
		}
	}
}
