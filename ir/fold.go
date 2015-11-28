// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir

import "github.com/rthornton128/calc/token"

func FoldConstants(o Object) Object {
	if pkg, ok := o.(*Package); ok {
		for k, v := range pkg.scope.m {
			pkg.scope.m[k] = fold(v)
		}
		return pkg
	}
	return fold(o)
}

func fold(o Object) Object {
	switch t := o.(type) {
	case *Assignment:
		t.Rhs = fold(t.Rhs)
	case *Binary:
		t.Lhs = fold(t.Lhs)
		t.Rhs = fold(t.Rhs)
		return foldBinary(t)
	case *Call:
		for i, e := range t.Args {
			t.Args[i] = fold(e)
		}
	case *Define:
		t.Body = fold(t.Body)
	case *For:
		t.Cond = fold(t.Cond)
		for i, e := range t.Body {
			t.Body[i] = fold(e)
		}
	case *Function:
		for i, e := range t.Body {
			t.Body[i] = fold(e)
		}
	case *If:
		t.Cond = fold(t.Cond)
		t.Then = fold(t.Then)
		if t.Else != nil {
			t.Else = fold(t.Else)
		}
	case *Unary:
		t.Rhs = fold(t.Rhs)
		return foldUnary(t)
	case *Variable:
		for i, e := range t.Body {
			t.Body[i] = fold(e)
		}
	}
	return o
}

func foldBinary(b *Binary) Object {
	lhs, lhsOk := b.Lhs.(*Constant)
	rhs, rhsOk := b.Rhs.(*Constant)

	if lhsOk && rhsOk {
		switch b.Type().Base().Kind() {
		case Int:
			l, r := int64(lhs.value.(intValue)), int64(rhs.value.(intValue))
			switch b.Op {
			case token.ADD:
				lhs.value = intValue(l + r)
			case token.MUL:
				lhs.value = intValue(l * r)
			case token.QUO:
				// TODO div by zero
				lhs.value = intValue(l / r)
			case token.REM:
				lhs.value = intValue(l % r)
			case token.SUB:
				lhs.value = intValue(l - r)
			}
			return lhs
		case Bool:
			switch lhs.Type().Base().Kind() {
			case Bool:
				l, r := bool(lhs.value.(boolValue)), bool(rhs.value.(boolValue))
				switch b.Op {
				case token.EQL:
					lhs.value = boolValue(l == r)
				case token.NEQ:
					lhs.value = boolValue(l != r)
				}
			case Int:
				l, r := int64(lhs.value.(intValue)), int64(rhs.value.(intValue))
				switch b.Op {
				case token.EQL:
					lhs.value = boolValue(l == r)
				case token.NEQ:
					lhs.value = boolValue(l != r)
				case token.GTT:
					lhs.value = boolValue(l > r)
				case token.GTE:
					lhs.value = boolValue(l >= r)
				case token.LST:
					lhs.value = boolValue(l < r)
				case token.LTE:
					lhs.value = boolValue(l <= r)
				}
			}
			return lhs
		}
	}
	return b
}

func foldUnary(u *Unary) Object {
	if c, ok := u.Rhs.(*Constant); ok {
		switch u.Op {
		case "+":
			val := intValue(+int64(c.value.(intValue)))
			if val < 0 {
				c.value = val * -1
			}
		case "-":
			c.value = intValue(-int64(c.value.(intValue)))
		}
		return c
	}
	return u
}
