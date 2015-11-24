// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir

import (
	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

type typeChecker struct {
	token.ErrorList
	fset *token.FileSet
}

func TypeCheck(o Object, fs *token.FileSet) error {
	t := &typeChecker{ErrorList: make(token.ErrorList, 0), fset: fs}
	if pkg, ok := o.(*Package); ok {
		for _, decl := range pkg.Scope().m {
			t.check(decl)
		}
	} else {
		t.check(o)
	}
	if t.ErrorList.Count() != 0 {
		return t.ErrorList
	}
	return nil
}

func (tc *typeChecker) check(o Object) {
	switch t := o.(type) {
	case *Assignment:
		o := t.Scope().Lookup(t.Lhs)
		if o == nil {
			tc.error(t.Pos(), "undeclared variable ", t.Lhs)
			return
		}
		if o.Kind() != ast.VarDecl {
			tc.error(t.Pos(), "may only assign to variables but ", o.Name(), " is ",
				o.Kind())
			return
		}
		tc.check(t.Rhs)
		if o.Type() != t.Rhs.Type() {
			tc.error(t.Pos(), "variable ", t.Name(), " is of type ", t.Type(), " but ",
				" assignment of type ", t.Rhs.Type())
			return
		}
	case *Binary:
		tc.check(t.Lhs)
		tc.check(t.Rhs)
		typ := Int
		if (t.Op == token.EQL || t.Op == token.NEQ) && t.Lhs.Type() == Bool {
			typ = Bool
		}
		if t.Lhs.Type() != typ {
			tc.error(t.Pos(), "binary expected type ", typ, " but lhs is type ",
				t.Lhs.Type())
			return
		}
		if t.Rhs.Type() != typ {
			tc.error(t.Pos(), "binary expected type ", typ, " but rhs is type ",
				t.Rhs.Type())
			return
		}
	case *Call:
		o := t.Scope().Lookup(t.Name())
		if o == nil {
			tc.error(t.Pos(), "calling undeclared function ", t.Name())
			return
		}
		if o.Kind() != ast.FuncDecl {
			tc.error(t.Pos(), "call expects function got ", o.Kind())
			return
		}
		f := o.(*Define).Body.(*Function)

		if len(t.Args) != len(f.Params) {
			tc.error(t.Pos(), "function ", t.Name(), " expects ", len(f.Params),
				" arguments but received ", len(t.Args))
			return
		}

		for i, a := range t.Args {
			tc.check(a)
			p := f.Scope().Lookup(f.Params[i])
			if a.Type() != p.Type() {
				tc.error(t.Pos(), "parameter ", i, " of function ", t.Name(),
					" expects type ", p.Type(), " but argument ", i, " is of type ",
					a.Type())
			}
		}
		t.object.typ = f.Type()
	case *Define:
		tc.check(t.Body)
	case *For:
		tc.check(t.Cond)
		if t.Cond.Type() != Bool {
			tc.error(t.Pos(), "conditional must be type bool, got ", t.Cond.Type())
			return
		}
		tc.checkBody(t.Type(), t.Pos(), t.Body)
	case *Function:
		tc.checkBody(t.Type(), t.Pos(), t.Body)
	case *If:
		tc.check(t.Cond)
		if t.Cond.Type() != Bool {
			tc.error(t.Pos(), "conditional must be type bool, got ", t.Cond.Type())
			return
		}
		tc.check(t.Then)
		if t.Type() != t.Then.Type() {
			tc.error(t.Pos(), "if expects type ", t.Type(),
				" but then clause of type ", t.Then.Type())

			return
		}
		if t.Else != nil {
			tc.check(t.Else)
			if t.Type() != t.Else.Type() {
				tc.error(t.Pos(), "if expects type ", t.Type(),
					" but else clause of type ", t.Else.Type())

				return
			}
		}
	case *Var:
		o := t.Scope().Lookup(t.Name())
		if o == nil {
			tc.error(t.Pos(), "undeclared variable ", t.Name())
			return
		}
		if o.Kind() == ast.FuncDecl {
			tc.error(t.Pos(), "function ", t.Name(), "used as variable; must be used "+
				"in call form (surrounded in parentheses")
			return
		}
		t.object.typ = o.Type()
	case *Variable:
		tc.checkBody(t.Type(), t.Pos(), t.Body)
	}
}

func (tc *typeChecker) checkBody(t Type, p token.Pos, body []Object) {
	for _, e := range body {
		tc.check(e)
	}

	tail := body[len(body)-1]
	if t != tail.Type() {
		tc.error(p, "last expression of var is of type ", tail.Type(),
			" but var is of type ", t)
	}
}

func (t *typeChecker) error(p token.Pos, args ...interface{}) {
	t.Add(t.fset.Position(p), args...)
}
