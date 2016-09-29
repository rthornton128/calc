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
		top:        make(map[string]*regAllocs),
		nextOffset: 0 - arch.Width(),
		Arch:       arch,
	}

	// assign offsets to parameters of all functions first
	for _, f := range pkg.Scope().Names() {
		if d, ok := pkg.Lookup(f).(*ir.Define); ok {
			switch t := d.Body.(type) {
			case *ir.Function:
				a.openScope(f)
				a.alloc(t)
				a.closeScope()
				//case *ir.Variable:
				//a.alloc(t)
			}
		}
	}
	return &a
}

func align16(n int) int { return (n & -16) + 16 }

type regAllocs struct {
	locs     map[string]string
	szParams int
	szLocals int
}

// allocator is a rudimentary register allocator that mainly just spills
// everything to the stack and does little to no optimization
type allocator struct {
	Arch
	current    *regAllocs
	top        map[string]*regAllocs
	fn         string
	nextOffset int
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

func (a *allocator) closeScope() {
	fmt.Println("close:", a.fn, a.current.szLocals, a.current.szParams)
	a.top[a.fn] = a.current
}

func (a *allocator) openScope(fn string) {
	if fn == a.fn {
		fmt.Println("same scope:", a.current)
		return
	}

	if s, ok := a.top[fn]; ok {
		fmt.Println("existing scope:", s)
		a.current = s
	} else {
		fmt.Println("new scope:", fn)
		a.current = &regAllocs{locs: make(map[string]string)}
	}
	a.fn = fn
}

func (a *allocator) getByID(id int) string {
	return a.getByName(fmt.Sprintf("%d", id))
}

func (a *allocator) getByName(name string) string {
	fmt.Printf("get: %s = %s\n", name, a.current.locs[name])
	return a.current.locs[name]
}

func (a *allocator) insertByID(id int, loc string) {
	a.insertByName(fmt.Sprintf("%d", id), loc)
}

func (a *allocator) insertByName(name string, loc string) {
	fmt.Printf("insert: %s = %s\n", name, loc)
	a.current.locs[name] = loc
}

func (a *allocator) nextLoc() string {
	s := fmt.Sprintf("%d(%s)", a.nextOffset, a.Register(BP))
	a.nextOffset -= a.Width()
	a.current.szLocals += a.Width()
	fmt.Println("szLocals:", a.current.szLocals)
	return s
}

func (a *allocator) stackSize() int {
	return align16(a.current.szParams) + align16(a.current.szLocals)
}

func (a *allocator) alloc(o ir.Object) {
	switch t := o.(type) {
	case *ir.Function:
		// set parameter registers and offsets
		for i, p := range t.Params {
			a.insertByName(p.Name(), a.CallStackOffset(i))
			a.current.szParams += a.Width()
		}

		// locals
		for _, o := range t.Body {
			a.walk(o)
		}
	case *ir.Variable:
		a.walk(t)
	}
	fmt.Println("current:", a.current)
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
			a.insertByID(t.Rhs.ID(), a.nextLoc())
		}
	case *ir.Call:
		for i, arg := range t.Args {
			a.walk(arg)

			// ensure enough stack space available for params stored on stack
			if i >= len(t.Args) {
				a.current.szParams += a.Width()
			}
		}
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
		v := t.Scope().Lookup(t.Name())
		if d, ok := v.(*ir.Define); ok {
			fmt.Println("alloc define")
			switch t := d.Body.(type) {
			case *ir.Variable:
				fmt.Println("define variable, id:", t.ID())
				loc := a.nextLoc()
				a.insertByName(d.Name(), loc)
				a.insertByID(t.ID(), loc)
				a.walk(t)
			default:
				fmt.Println("define other:", t)
				a.insertByName(d.Name(), a.nextLoc())
			}
		}
	case *ir.Variable:
		fmt.Println("var ID:", t.ID(), "name:", t.Name())
		a.insertByID(t.ID(), a.nextLoc())
		for _, p := range t.Params {
			a.insertByName(p.Name(), a.nextLoc())
		}
		for _, o := range t.Body {
			a.walk(o)
		}
	}
}
