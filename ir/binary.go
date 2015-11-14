// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

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
	o := newObject("binary", "", b.Pos(), ast.None, parent)
	o.typ = binaryType(b.Op)

	lhs := MakeExpr(b.List[0], parent)
	for _, e := range b.List[1:] {
		rhs := MakeExpr(e, parent)
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
	o := newObject("unary", "", u.Pos(), ast.None, parent)
	o.typ = Int

	return &Unary{
		object: o,
		Op:     u.Op,
		Rhs:    MakeExpr(u.Value, parent),
	}
}

func (u *Unary) String() string {
	return fmt.Sprintf("%s(%s)", u.Op, u.Rhs.String())
}
