package ir

import (
	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

type typeChecker struct {
	token.ErrorList
	fset *token.FileSet
}

func TypeCheck(pkg *Package, fs *token.FileSet) error {
	t := &typeChecker{ErrorList: make(token.ErrorList, 0), fset: fs}
	for _, decl := range pkg.Scope().m {
		t.check(decl)
	}
	return t.ErrorList
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
		if t.Lhs.Type() != Int {
			tc.error(t.Pos(), "expected type ", Int, " but argument is type ",
				t.Lhs.Type())
			return
		}
		if t.Rhs.Type() != Int {
			tc.error(t.Pos(), "expected type ", Int, " but argument is type ",
				t.Rhs.Type())
			return
		}
	case *Block:
		for _, e := range t.Exprs {
			tc.check(e)
		}
		t.object.typ = t.Exprs[len(t.Exprs)-1].Type()
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
		decl := o.(*Declaration)
		if len(t.Args) != len(decl.Params) {
			tc.error(t.Pos(), "function ", decl.Name(), " expects ", len(decl.Params),
				" arguments but received ", len(t.Args))
			return
		}
		t.object.typ = decl.Type()
	case *Declaration:
		tc.check(t.Body)
		if t.Type() != t.Body.Type() {
			tc.error(t.Pos(), "declaration of type ", t.Type(),
				" but body returns type ", t.Body.Type())

			return
		}
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
		if o.Kind() != ast.VarDecl {
			tc.error(t.Pos(), "may non reference non-variable ", o.Kind(), " named ",
				t.Name())
			return
		}
		t.object.typ = o.Type()
	case *Variable:
		if t.Assign != nil {
			assign := t.Assign.(*Assignment)
			tc.check(assign.Rhs)
			if t.Type() == Unknown {
				t.object.typ = assign.Rhs.Type()
			}
			if t.Type() != assign.Rhs.Type() {
				tc.error(t.Pos(), "variable ", t.Name(), " expects type ", t.Type(),
					" but initializer is of type ", assign.Rhs.Type())
				return
			}
		}
	}
}

func (t *typeChecker) error(p token.Pos, args ...interface{}) {
	t.Add(t.fset.Position(p), args...)
}
