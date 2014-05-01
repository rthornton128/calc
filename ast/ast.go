// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ast

import (
	"github.com/rthornton128/calc/token"
)

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

type Expr interface {
	Node
	exprNode()
}

type AssignExpr struct{}

type BasicLit struct {
	LitPos token.Pos
	Kind   token.Token
	Lit    string
}

type BinaryExpr struct {
	Expression
	Op    token.Token
	OpPos token.Pos
	List  []Expr
}

type CallExpr struct{}

type DeclExpr struct {
	Name *Ident
	/* Type, Params, decl keyword */
}

type Expression struct {
	Opening token.Pos
	Closing token.Pos
}

type File struct {
	Root Expr
}

type Ident struct {
	NamePos token.Pos
	Name    string
	Value   Expr
}

type IfExpr struct {
	Expression
	If   token.Pos
	Cond Expr
	Then Expr
	Else Expr
	// needs Type field
}

type Scope struct{}

type VarExpr struct{}

func (b *BasicLit) Pos() token.Pos   { return b.LitPos }
func (e *Expression) Pos() token.Pos { return e.Opening }
func (f *File) Pos() token.Pos       { return f.Root.Pos() }

func (b *BasicLit) End() token.Pos   { return b.LitPos + token.Pos(len(b.Lit)) }
func (e *Expression) End() token.Pos { return e.Closing }
func (f *File) End() token.Pos       { return f.Root.End() }

func (b *BasicLit) exprNode()   {}
func (e *Expression) exprNode() {}
