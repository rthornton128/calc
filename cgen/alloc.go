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

type funcAllocs struct {
	stackSz int
	offsets map[string]int
}

var fnStackAllocs = make(map[string]funcAllocs)

// StackAlloc sets stack offsets for all expressions needing one. It does is a
// rudimentary register allocator that mainly just spills everything to the
// stack and does little to no optimization
func StackAlloc(pkg *ir.Package, sz int) {
	a := allocator{ptrSz: sz}

	// assign offsets to parameters of all functions first
	for _, f := range pkg.Scope().Names() {
		if t, ok := pkg.Lookup(f).(*ir.Define); ok {
			a.curFn = f
			fnStackAllocs[f] = funcAllocs{}
			a.alloc(t.Body)
		}
	}
}

func align(n int) int { return (n & -16) + 16 }

func typeSize(t ir.Type) int {
	switch t {
	case ir.Int:
		return 4
	case ir.Bool:
		return 1
	}
	// not reachable
	return 0
}

// allocator is a rudimentary register allocator that mainly just spills
// everything to the stack and does little to no optimization
type allocator struct {
	ptrSz int
	curFn string
	off   int
	sz    int
}

func (a *allocator) alloc(o ir.Object) {
	if t, ok := o.(*ir.Function); ok {
		fnStackAllocs[a.curFn] = funcAllocs{0, make(map[string]int)}

		// set parameter offsets (bp+offset)
		offset := a.ptrSz * 3 // starting offset = size(BP) + size(IP) + size(ptr)
		tmp := fnStackAllocs[a.curFn]
		for _, p := range t.Params {
			tmp.offsets[p.Name()] = offset
			offset += a.ptrSz
		}

		// locals
		a.off = 0
		for _, o := range t.Body {
			a.walk(o)
		}
		tmp.stackSz = align(a.sz + (len(t.Params) * a.ptrSz))
		fnStackAllocs[a.curFn] = tmp

		// reset
		a.sz = 0
	}
}

func (a *allocator) nextOffset() int { //t ir.Type) int {
	o := a.off
	a.off -= a.ptrSz //typeSize(t)
	a.sz += a.ptrSz
	return o
}

func (a *allocator) walk(o ir.Object) {
	switch t := o.(type) {
	case *ir.Assignment:
		a.walk(t.Rhs)
	case *ir.Binary:
		a.walk(t.Lhs)
		a.walk(t.Rhs)
		// TODO: requests more stack space than is strictly necessary
		switch t.Rhs.(type) {
		case *ir.Constant, *ir.Var:
			//no nothing
		default:
			tmp := fnStackAllocs[a.curFn]
			tmp.offsets[fmt.Sprintf("%d", t.Rhs.ID())] = a.nextOffset()
			fnStackAllocs[a.curFn] = tmp
		}
		//t.off = a.nextOffset(t.Rhs.Type())
		//a.off += typeSize(t.Rhs.Type())
	case *ir.Call:
		// arguments for functions are placed in reverse order (right to left)
		// and are relative to the bp-(a.sz+(i*a.ptrSz)) OR
		// sp+((len(args)*a.ptrSz)-(i*a.ptrSz)). Size a.sz is unknown during
		// an initial pass, the latter method is used
		a.sz = a.ptrSz * len(t.Args)
		//a.sz += sz
		//tmp := fnStackAllocs[a.curFn]
		for _, arg := range t.Args {
			a.walk(arg)
			//tmp.offsets[arg.ID()] = sz - (i * a.ptrSz)
		}
		//fnStackAllocs[a.curFn] = tmp
	case *ir.Function:
		// unreachable
	case *ir.If:
		a.walk(t.Cond)
		a.walk(t.Then)
		if t.Else != nil {
			a.walk(t.Else)
		}
	case *ir.Var:
	case *ir.Variable:
		for _, p := range t.Params {
			tmp := fnStackAllocs[a.curFn]
			tmp.offsets[p.Name()] = a.nextOffset()
			fnStackAllocs[a.curFn] = tmp
			//a.sz += typeSize(p.Type())
		}
		//a.sz += len(t.Params) * 4
		for _, o := range t.Body {
			a.walk(o)
		}
	case *ir.Unary:
		// operates directly on .ax, nothing to do
	}
}
