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

// Constant represents a basic literal's actual value. It is not a
// syntactical construct and will not be found in the original source.
// As such, it's position will always be NoPos as it should not be used
// in error reporting. It will be created only by the optimizer for
// constant folding.
type Constant struct {
	Kind  token.Token // TODO types.Kind?
	Value interface{} // TODO placeholder, needs proper interface/type
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
	Opening  token.Pos
	Closing  token.Pos
	RealType Type
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
	Object  *Object // may be nil (ie. Name is a type keyword)
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
	NamePos  token.Pos
	Name     string
	Kind     ObKind
	Offset   int
	RealType Type
	Type     *Ident
	Value    Expr
}

type ObKind int

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
	Var    token.Pos
	Name   *Ident
	Object *Object
}

func (b *BasicLit) Pos() token.Pos   { return b.LitPos }
func (c *Constant) Pos() token.Pos   { return token.NoPos }
func (e *Expression) Pos() token.Pos { return e.Opening }
func (f *File) Pos() token.Pos       { return token.NoPos }
func (i *Ident) Pos() token.Pos      { return i.NamePos }
func (p *Package) Pos() token.Pos    { return token.NoPos }
func (u *UnaryExpr) Pos() token.Pos  { return u.OpPos }

func (b *BasicLit) End() token.Pos   { return b.LitPos + token.Pos(len(b.Lit)) }
func (c *Constant) End() token.Pos   { return token.NoPos }
func (e *Expression) End() token.Pos { return e.Closing }
func (f *File) End() token.Pos       { return token.NoPos }
func (i *Ident) End() token.Pos      { return i.NamePos + token.Pos(len(i.Name)) }
func (p *Package) End() token.Pos    { return token.NoPos }
func (u *UnaryExpr) End() token.Pos  { return u.Value.End() }

func (b *BasicLit) exprNode()   {}
func (b *Constant) exprNode()   {} // not an expression
func (e *Expression) exprNode() {}
func (i *Ident) exprNode()      {}
func (u *UnaryExpr) exprNode()  {}

const (
	Decl ObKind = iota
	Var
)

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
