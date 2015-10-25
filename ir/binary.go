package ir

import (
	"fmt"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

type Binary struct {
	object
	id  int
	Op  token.Token
	Lhs Object
	Rhs Object
}

func makeBinary(b *ast.BinaryExpr, parent *Scope) *Binary {
	o := newObject("binary", "", b.Pos(), None, parent)
	o.typ = binaryType(b.Op)

	lhs := makeExpr(b.List[0], parent)
	for _, e := range b.List[1:] {
		rhs := makeExpr(e, parent)
		lhs = Object(&Binary{
			object: o,
			Op:     b.Op,
			Lhs:    lhs,
			Rhs:    rhs,
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

func (b *Binary) ID() int      { return b.id }
func (b *Binary) SetID(id int) { b.id = id }
func (b *Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Lhs.String(), b.Op, b.Rhs.String())
}

type Unary struct {
	object
	Op  string
	Rhs Object
}

func makeUnary(u *ast.UnaryExpr, parent *Scope) *Unary {
	o := newObject("unary", "", u.Pos(), None, parent)
	o.typ = Int

	return &Unary{
		object: o,
		Op:     u.Op,
		Rhs:    makeExpr(u.Value, parent),
	}
}

func (u *Unary) String() string {
	return fmt.Sprintf("%s(%s)", u.Op, u.Rhs.String())
}
