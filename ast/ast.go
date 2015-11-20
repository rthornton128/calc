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
	Op    token.Token
	OpPos token.Pos
	ID    int
	List  []Expr
}

type CallExpr struct {
	Name *Ident
	Args []Expr
}

type DefineStmt struct {
	Define token.Pos
	Name   *Ident
	Type   *Ident
	Kind   Kind
	Body   Expr
}

type File struct {
	Defs []*DefineStmt
}

type ForExpr struct {
	For  token.Pos
	Type *Ident
	Cond Expr
	Body []Expr
}

type FuncExpr struct {
	Func   token.Pos
	Type   *Ident
	Params []*Param
	Body   []Expr
}

type Ident struct {
	NamePos token.Pos
	Name    string
	Type    *Ident
}

type IfExpr struct {
	If   token.Pos
	Type *Ident
	Cond Expr
	Then Expr
	Else Expr
}

type Object struct {
	NamePos token.Pos
	Name    string
	Kind    Kind
}

type Package struct {
	Files []*File
}

type Param struct {
	Name *Ident
	Type *Ident
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
	Var    token.Pos
	Type   *Ident
	Params []*Param
	Body   []Expr
}

func (a *AssignExpr) Pos() token.Pos { return a.Equal }
func (b *BasicLit) Pos() token.Pos   { return b.LitPos }
func (b *BinaryExpr) Pos() token.Pos { return b.OpPos }
func (c *CallExpr) Pos() token.Pos   { return c.Name.Pos() }
func (d *DefineStmt) Pos() token.Pos { return d.Define }
func (f *File) Pos() token.Pos       { return token.NoPos }
func (f *ForExpr) Pos() token.Pos    { return f.For }
func (f *FuncExpr) Pos() token.Pos   { return f.Func }
func (i *Ident) Pos() token.Pos      { return i.NamePos }
func (i *IfExpr) Pos() token.Pos     { return i.If }
func (o *Object) Pos() token.Pos     { return o.NamePos }
func (p *Package) Pos() token.Pos    { return token.NoPos }
func (p *Param) Pos() token.Pos      { return p.Name.Pos() }
func (u *UnaryExpr) Pos() token.Pos  { return u.OpPos }
func (v *VarExpr) Pos() token.Pos    { return v.Var }

func (a *AssignExpr) exprNode() {}
func (b *BasicLit) exprNode()   {}
func (b *BinaryExpr) exprNode() {}
func (c *CallExpr) exprNode()   {}
func (f *ForExpr) exprNode()    {}
func (f *FuncExpr) exprNode()   {}
func (i *IfExpr) exprNode()     {}
func (i *Ident) exprNode()      {}
func (u *UnaryExpr) exprNode()  {}
func (v *VarExpr) exprNode()    {}

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
