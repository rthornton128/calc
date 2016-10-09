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

// Binary is the expression (operator arg0, args1 [, args...])
type Binary struct {
	object
	Op  token.Token
	Lhs Object
	Rhs Object
}

func makeBinary(pkg *Package, b *ast.BinaryExpr) *Binary {
	lhs := MakeExpr(pkg, b.List[0])
	for _, e := range b.List[1:] {
		rhs := MakeExpr(pkg, e)
		lhs = Object(&Binary{
			object: object{
				id:  pkg.getID(),
				pkg: pkg,
				pos: b.Pos(),
				typ: binaryType(b.Op)},
			Op:  b.Op,
			Lhs: lhs,
			Rhs: rhs,
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

// Copy makes a deep copy of the Binary object
func (b *Binary) Copy() Object {
	return &Binary{
		object: b.object.copy(b.Package().getID()),
		Op:     b.Op,
		Lhs:    b.Lhs.Copy(),
		Rhs:    b.Rhs.Copy(),
	}
}

func (b *Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Lhs.String(), b.Op, b.Rhs.String())
}

// Unary is the expreesion -expr
type Unary struct {
	object
	Op  string
	Rhs Object
}

func makeUnary(pkg *Package, u *ast.UnaryExpr) *Unary {
	return &Unary{
		object: object{pkg: pkg, pos: u.Pos(), scope: pkg.scope, typ: Int},
		Op:     u.Op,
		Rhs:    MakeExpr(pkg, u.Value),
	}
}

// Copy makes a deep copy of the Unary object
func (u *Unary) Copy() Object {
	return &Unary{
		object: u.object.copy(0),
		Op:     u.Op,
		Rhs:    u.Rhs.Copy(),
	}
}

func (u *Unary) String() string {
	return fmt.Sprintf("%s(%s)", u.Op, u.Rhs.String())
}
