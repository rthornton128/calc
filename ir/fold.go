package ir

import "github.com/rthornton128/calc/token"

func FoldConstants(pkg *Package) *Package {
	for k, v := range pkg.scope.m {
		pkg.scope.m[k] = fold(v)
	}

	return pkg
}

func fold(o Object) Object {
	switch t := o.(type) {
	case *Assignment:
		t.Rhs = fold(t.Rhs)
	case *Binary:
		t.Lhs = fold(t.Lhs)
		t.Rhs = fold(t.Rhs)
		return foldBinary(t)
	case *Block:
		for i, e := range t.Exprs {
			t.Exprs[i] = fold(e)
		}
	case *Declaration:
		t.Body = fold(t.Body)
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
		t.Assign = fold(t.Assign)
	}
	return o
}

func foldBinary(b *Binary) Object {
	lhs, lhsOk := b.Lhs.(*Constant)
	rhs, rhsOk := b.Rhs.(*Constant)

	if lhsOk && rhsOk {
		switch b.Type() {
		case Int:
			l, r := int64(lhs.value.(intValue)), int64(rhs.value.(intValue))
			switch b.Op {
			case token.ADD:
				lhs.value = intValue(l + r)
			case token.MUL:
				lhs.value = intValue(l * r)
			case token.QUO:
				lhs.value = intValue(l / r)
			case token.REM:
				lhs.value = intValue(l % r)
			case token.SUB:
				lhs.value = intValue(l - r)
			}
			return lhs
		case Bool:
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
			return lhs
		}
	}
	return b
}

func foldUnary(u *Unary) Object {
	if c, ok := u.Rhs.(*Constant); ok {
		switch u.Op {
		case "+":
			c.value = intValue(+int64(c.value.(intValue)))
		case "-":
			c.value = intValue(-int64(c.value.(intValue)))
		}
		return c
	}
	return u
}
