package ir

import (
	"fmt"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

type Binary struct {
	object
	id  int
	op  token.Token
	lhs Object
	rhs Object
}

func makeBinary(b *ast.BinaryExpr, parent *Scope) *Binary {
	o := newObject("binary", "", parent)
	o.typ = binaryType(b.Op)

	lhs := makeExpr(b.List[0], parent)
	for _, e := range b.List[1:] {
		rhs := makeExpr(e, parent)
		lhs = Object(&Binary{
			object: o,
			op:     b.Op,
			lhs:    lhs,
			rhs:    rhs,
		})
	}
	return lhs.(*Binary)
}

func binaryType(t token.Token) Type {
	switch t {
	case token.ADD, token.MUL, token.QUO, token.REM, token.SUB:
		return Int
	default:
		return Bool
	}
}

func (b *Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.lhs.String(), b.op, b.rhs.String())
}

type Unary struct {
	object
	op  string
	rhs Object
}

func makeUnary(u *ast.UnaryExpr, parent *Scope) *Unary {
	o := newObject("unary", "", parent)
	o.typ = Int

	return &Unary{
		object: o,
		op:     u.Op,
		rhs:    makeExpr(u.Value, parent),
	}
}

func (u *Unary) String() string {
	return fmt.Sprintf("%s(%s)", u.op, u.rhs.String())
}
