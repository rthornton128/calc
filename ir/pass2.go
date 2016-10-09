// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir

// ReplaceMacros replaces defines with copies of expressions
func (pkg *Package) ReplaceMacros(o Object) Object {
	switch t := o.(type) {
	case *Assignment:
		t.Rhs = pkg.ReplaceMacros(t.Rhs)
	case *Binary:
		//fmt.Println("lhs:", t.Lhs)
		t.Lhs = pkg.ReplaceMacros(t.Lhs)
		//fmt.Println("lhs after:", t.Lhs)
		//fmt.Println("lhs:", t.Rhs)
		t.Rhs = pkg.ReplaceMacros(t.Rhs)
		//fmt.Println("rhs after:", t.Rhs)
	case *Call:
		for k, v := range t.Args {
			t.Args[k] = pkg.ReplaceMacros(v)
		}
	case *Constant:
		// nothing to do but return value
	case *Function:
		for k, v := range t.Body {
			t.Body[k] = pkg.ReplaceMacros(v)
		}
	case *For:
		for k, v := range t.Body {
			t.Body[k] = pkg.ReplaceMacros(v)
		}
	case *If:
		t.Then = pkg.ReplaceMacros(t.Then)
		t.Else = pkg.ReplaceMacros(t.Else) // will this panic?
	case *Package:
		for _, v := range pkg.Scope().Names() {
			pkg.ReplaceMacros(pkg.Scope().Lookup(v))
		}
	case *Unary:
		t.Rhs = pkg.ReplaceMacros(t.Rhs)
	case *Var:
		x := pkg.top.Lookup(t.Name())
		if x != nil {
			//fmt.Println("var", t.Name(), "not nil so making copy of:", x)
			//foo :=
			return x.Copy()
			//fmt.Println("after copy:", foo)
			//return foo
		}
	case *Variable:
		for k, v := range t.Body {
			//fmt.Println("body", k, "-", v)
			t.Body[k] = pkg.ReplaceMacros(v)
			//fmt.Println("body", k, "after -", t.Body[k])
		}
	}
	return o
}
