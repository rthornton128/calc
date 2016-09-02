// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir

// RegAlloc sets stack offsets for all expressions needing one. It does is a
// rudimentary register allocator that mainly just spills everything to the
// stack and does little to no optimization
func RegAlloc(pkg *Package) {
	a := allocator{}

	// assign offsets to parameters of all functions first
	for _, f := range pkg.Scope().Names() {
		if t, ok := pkg.Lookup(f).(*Define); ok {
			a.process(t.Body)
		}
	}
}

// allocator is a rudimentary register allocator that mainly just spills
// everything to the stack and does little to no optimization
type allocator struct {
	off int
	sz  int
}

func (a *allocator) process(o Object) {
	offset := 0

	switch t := o.(type) {
	case *Function:
		offset = byteOffset
		for _, p := range t.Params {
			p.reg = BP
			p.off = offset
			offset += 4
		}
		defer a.walk(t)
	case *Variable:
		for _, p := range t.Params {
			p.reg = SP
			p.off = offset
			offset += 4
		}
		defer a.walk(t)
	}
}

func (a *allocator) nextOffset() int {
	o := a.off
	a.off += 4
	return o
}

func (a *allocator) walk(o Object) {
	switch t := o.(type) {
	case *Assignment:
		a.walk(t.Rhs)
	case *Binary:
		// TODO: requests more stack space than is necessary since only
		// rhs expressions would potentially be placed on the stack
		a.walk(t.Lhs)
		a.walk(t.Rhs)
		t.off = a.nextOffset()
		a.sz += 4
	case *Call:
		a.sz += len(t.Args) * 4
		for _, arg := range t.Args {
			a.walk(arg)
		}
	case *Function:
		for _, o := range t.Body {
			a.walk(o)
		}
		t.off = a.sz
		a.sz = 0
	case *If:
		a.walk(t.Cond)
		a.walk(t.Then)
		if t.Else != nil {
			a.walk(t.Else)
		}
	case *Var:
	case *Variable:
		for _, p := range t.Params {
			p.off = a.nextOffset()
			p.reg = SP
		}
		a.sz += len(t.Params) * 4
		for _, o := range t.Body {
			a.walk(o)
		}
	}
}
