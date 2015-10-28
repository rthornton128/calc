// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ast_test

import (
	"testing"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

func TestBasicLit(t *testing.T) {
	b := &ast.BasicLit{
		LitPos: token.Pos(1),
		Kind:   token.INTEGER,
		Lit:    "24",
	}
	pos := token.Pos(1)
	if b.Pos() != pos {
		t.Fatal("Expected:", pos, "Got:", b.Pos())
	}
}

func TestBinaryExpr(t *testing.T) {
	// (+ 3 5)
	x := &ast.BasicLit{
		LitPos: token.Pos(4),
		Kind:   token.INTEGER,
		Lit:    "3",
	}
	y := &ast.BasicLit{
		LitPos: token.Pos(6),
		Kind:   token.INTEGER,
		Lit:    "5",
	}
	b := &ast.BinaryExpr{
		Expression: ast.Expression{
			Opening: token.Pos(1),
			Closing: token.Pos(7),
		},
		Op:    token.ADD,
		OpPos: token.Pos(2),
		List:  []ast.Expr{x, y},
	}

	if b.Pos() != token.Pos(1) {
		t.Fatal("BinaryExpr: Expected: 1 Got:", b.Pos())
	}
}
