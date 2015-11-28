// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir

import (
	"fmt"

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
			tc.error(t.Pos(), "undeclared variable '%s'", t.Lhs)
			return
		}
		tc.check(t.Rhs)

		if !typeMatch(o.Type(), t.Rhs.Type()) {
			tc.error(t.Pos(), "type mismatch; '%s' is of type '%s' but expression "+
				"is of type '%s'", o.Name(), o.Type(), t.Rhs.Type())
		}

		t.object.typ = o.Type()
	case *Binary:
		tc.check(t.Lhs)
		tc.check(t.Rhs)
		switch t.Op {
		case token.GTT, token.GTE, token.LST, token.LTE:
			if t.Lhs.Type().Base().Kind() == Bool ||
				t.Rhs.Type().Base().Kind() == Bool {
				tc.error(t.Pos(), "boolean expressions can only tested for equality")
			}
			fallthrough
		case token.EQL, token.NEQ:
			if !typeMatch(t.Lhs.Type(), t.Rhs.Type()) {
				tc.error(t.Pos(), "can not compare expression of type '%s' with "+
					"expression of type '%s'", t.Lhs.Type(), t.Rhs.Type())
			}
			t.typ = TypeList[Bool]
		default: // arithmetic
			if !typeMatch(t.Lhs.Type(), t.Rhs.Type()) {
				tc.error(t.Pos(), "can not compare expression of type '%s' with "+
					"expression of type '%s'", t.Lhs.Type(), t.Rhs.Type())
			}
			t.typ = TypeList[Int]
		}
	case *Call:
		o := t.Scope().Lookup(t.Name())
		if o == nil {
			tc.error(t.Pos(), "calling undeclared function '%s'", t.Name())
			return
		}
		f, ok := o.Type().(FuncType)
		if !ok {
			tc.error(t.Pos(), "call expects function got '%s'", o.Type())
			return
		}

		if len(t.Args) != len(f.Params) {
			tc.error(t.Pos(), "function '%s' expects '%d' arguments but received %d",
				t.Name(), len(f.Params), len(t.Args))
			return
		}

		for i, a := range t.Args {
			tc.check(a)
			//p := f.Scope().Lookup(f.Params[i].Name())
			if !typeMatch(a.Type(), f.Params[i]) {
				tc.error(t.Pos(), "parameter %d of function '%s' expects type '%s' "+
					"but argument is of type '%s'", i, t.Name(), f.Params[i], a.Type())
			}
		}
		t.object.typ = f.Return
	case *Define:
		tc.check(t.Body)
		o := t.Scope().Lookup(t.Name())
		if !typeMatch(o.Type(), t.Body.Type()) {
			tc.error(t.Pos(), "type mismatch; '%s' expects type '%s' but got '%s'",
				o.Name(), o.Type(), t.Body.Type())
		}
	case *For:
		tc.check(t.Cond)
		if t.Cond.Type().Base().Kind() != Bool {
			tc.error(t.Pos(), "conditional must be type 'bool', got '%s'",
				t.Cond.Type())
			return
		}
		tc.checkBody(t, t.Body)
	case *Function:
		//fmt.Println("function type:", t.Type())
		tc.checkBody(t, t.Body)
	case *If:
		tc.check(t.Cond)
		if t.Cond.Type().Base().Kind() != Bool {
			tc.error(t.Pos(), "conditional must be type 'bool', got '%s'",
				t.Cond.Type())
			return
		}
		tc.check(t.Then)
		if !typeMatch(t.Type(), t.Then.Type()) {
			tc.error(t.Pos(), "if expects type '%s' but then clause is type '%s'",
				t.Type(), t.Then.Type())

			return
		}
		if t.Else != nil {
			tc.check(t.Else)
			if !typeMatch(t.Type(), t.Else.Type()) {
				tc.error(t.Pos(), "if expects type '%s' but else clause is type '%s'",
					t.Type(), t.Else.Type())

				return
			}
		}
	case *Var:
		o := t.Scope().Lookup(t.Name())
		if o == nil {
			tc.error(t.Pos(), "undeclared variable '%s'", t.Name())
			return
		}
		if _, ok := o.Type().(FuncType); ok {
			tc.error(t.Pos(), "function '%s' used as variable; must be used "+
				"in call form (surrounded in parentheses)", t.Name())
			return
		}
		t.object.typ = o.Type()
	case *Variable:
		tc.checkBody(t, t.Body)
	}
}

func typeMatch(a Type, b Type) bool {
	switch ta := a.(type) {
	case BasicType:
		if tb, ok := b.(BasicType); ok {
			return ta == tb
		}
		return false
	case FuncType:
		if tb, ok := b.(FuncType); ok {
			if len(ta.Params) != len(tb.Params) {
				return false
			}
			if typeMatch(ta.Return, tb.Return) == false {
				return false
			}
			for i := range ta.Params {
				if !typeMatch(ta.Params[i], tb.Params[i]) {
					return false
				}
			}
			return true
		}

	}
	return false
}

func (tc *typeChecker) checkBody(o Object, body []Object) {
	for _, e := range body {
		tc.check(e)
	}

	tail := body[len(body)-1]
	if !typeMatch(o.Type().Base(), tail.Type()) {
		tc.error(o.Pos(), "last expression of %s is of type '%s' but expects "+
			"type '%s'", o.Name(), tail.Type(), o.Type())
	}
}

func (t *typeChecker) error(p token.Pos, format string, args ...interface{}) {
	t.Add(t.fset.Position(p), fmt.Sprintf(format, args...))
}
