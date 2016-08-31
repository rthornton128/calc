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
		offset := 0

		switch t := pkg.Lookup(f).(type) {
		case *Function:
			for _, p := range t.Params {
				p.off = offset
				offset += 4
			}
			defer a.walk(t)
		case *Variable:
			for _, p := range t.Params {
				p.off = offset
				offset += 4
			}
			defer a.walk(t)
		}
	}
}

// allocator is a rudimentary register allocator that mainly just spills
// everything to the stack and does little to no optimization
type allocator struct {
	off int
	sz  int
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
		a.walk(t.Lhs)
		a.walk(t.Rhs)
	case *Function:
		for _, o := range t.Body {
			a.walk(o)
		}
	case *Var:
		if t.off == 0 {
			t.off = a.nextOffset()
		}
	case *Variable:
		for _, o := range t.Body {
			a.walk(o)
		}
	}
}
