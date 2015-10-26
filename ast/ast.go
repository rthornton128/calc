// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ast

import "github.com/rthornton128/calc/token"

type Node interface {
	Pos() token.Pos
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
	ID    int
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
	Body   Expr
	Scope  *Scope
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
	Decls []*DeclExpr
	Scope *Scope
}

type Ident struct {
	NamePos token.Pos
	Name    string
	Type    *Ident
}

type IfExpr struct {
	Expression
	If    token.Pos
	Type  *Ident
	Cond  Expr
	Then  Expr
	Else  Expr
	Scope *Scope
}

type Object struct {
	NamePos token.Pos
	Name    string
	Kind    token.Kind
}

type Package struct {
	Scope *Scope
	Files []*File
}

type Scope struct {
	Parent *Scope
	Table  map[string]*Object
}

type UnaryExpr struct {
	OpPos token.Pos
	Op    string
	Value Expr
}

type VarExpr struct {
	Expression
	Var   token.Pos
	Name  *Ident
	Value Expr
}

func (b *BasicLit) Pos() token.Pos   { return b.LitPos }
func (e *Expression) Pos() token.Pos { return e.Opening }
func (f *File) Pos() token.Pos       { return token.NoPos }
func (i *Ident) Pos() token.Pos      { return i.NamePos }
func (p *Package) Pos() token.Pos    { return token.NoPos }
func (u *UnaryExpr) Pos() token.Pos  { return u.OpPos }

func (b *BasicLit) exprNode()   {}
func (e *Expression) exprNode() {}
func (i *Ident) exprNode()      {}
func (u *UnaryExpr) exprNode()  {}

func NewScope(parent *Scope) *Scope {
	return &Scope{Parent: parent, Table: make(map[string]*Object)}
}

func (s *Scope) Insert(ob *Object) *Object {
	if old, ok := s.Table[ob.Name]; ok {
		return old
	}
	s.Table[ob.Name] = ob
	return nil
}

func (s *Scope) Lookup(ident string) *Object {
	ob, ok := s.Table[ident]
	if ok || s.Parent == nil {
		return ob
	}
	return s.Parent.Lookup(ident)
}
