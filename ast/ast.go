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

type AssignExpr struct {
	Expression
	Equal token.Pos
	Name  *Ident
	Value Expr
}

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

type CallExpr struct {
	Expression
	Name *Ident
	Args []Expr
}

type DeclExpr struct {
	Expression
	Decl   token.Pos
	Name   *Ident
	Type   *Ident
	Params []*Ident
	Body   ExprList
}

type Expression struct {
	Opening token.Pos
	Closing token.Pos
}

type ExprList struct {
	Expression
	List []Expr
}

type File struct {
	/* Empty */
}

type Ident struct {
	NamePos token.Pos
	Name    string
	Value   Expr
}

type IfExpr struct {
	Expression
	If   token.Pos
	Type *Ident
	Cond Expr
	Then ExprList
	Else ExprList
}

type Package struct {
	Scope *Scope
	//Files map[string]*File
}

type Scope struct {
	parent *Scope
	table  map[string]Expr
}

type VarExpr struct {
	Expression
	Var   token.Pos
	Name  *Ident
	Type  *Ident
	Value Expr
}

func (b *BasicLit) Pos() token.Pos   { return b.LitPos }
func (e *Expression) Pos() token.Pos { return e.Opening }
func (f *File) Pos() token.Pos       { return token.NoPos }
func (p *Package) Pos() token.Pos    { return token.NoPos }

func (b *BasicLit) End() token.Pos   { return b.LitPos + token.Pos(len(b.Lit)) }
func (e *Expression) End() token.Pos { return e.Closing }
func (f *File) End() token.Pos       { return token.NoPos }
func (p *Package) End() token.Pos    { return token.NoPos }

func (b *BasicLit) exprNode()   {}
func (e *Expression) exprNode() {}
func (e *ExprList) exprNode()   {}

func NewScope(parent *Scope) *Scope {
	return &Scope{parent: parent, table: make(map[string]Expr)}
}

func (s *Scope) Insert(ident *Ident, exp Expr) {
	s.table[ident.Name] = exp
}

func (s *Scope) Lookup(ident string) Expr {
	var exp Expr
	var ok bool
	if exp, ok = s.table[ident]; !ok {
		if s.parent == nil {
			return nil
		}
		return s.parent.Lookup(ident)
	}
	return exp
}
