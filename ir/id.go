// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir

func Tag(o Object) {
	var nextID int
	if pkg, ok := o.(*Package); ok {
		for _, v := range pkg.Scope().m {
			tag(v, &nextID)
		}
	} else {
		tag(o, &nextID)
	}
}

func tag(o Object, nextID *int) {
	switch t := o.(type) {
	case *Assignment:
		tag(t.Rhs, nextID)
	case *Binary:
		setID(t, nextID)
		tag(t.Lhs, nextID)
		tag(t.Rhs, nextID)
	case *Call:
		for _, arg := range t.Args {
			tag(arg, nextID)
		}
	case *Define:
		tag(t.Body, nextID)
	case *Function:
		for _, p := range t.Params {
			setID(t.Scope().Lookup(p).(IDer), nextID)
		}
		for _, e := range t.Body {
			tag(e, nextID)
		}
	case *If:
		setID(t, nextID)
		tag(t.Cond, nextID)
		tag(t.Then, nextID)
		if t.Else != nil {
			tag(t.Else, nextID)
		}
	case *Param:
		setID(t.Scope().Lookup(t.Name()).(IDer), nextID)
	case *Unary:
		tag(t.Rhs, nextID)
	case *Variable:
		for _, p := range t.Params {
			setID(t.Scope().Lookup(p).(IDer), nextID)
		}
		for _, e := range t.Body {
			tag(e, nextID)
		}
	}
	return
}

func setID(o IDer, nextID *int) {
	if o.ID() == 0 {
		o.SetID(*nextID)
		*nextID++
	}
}
