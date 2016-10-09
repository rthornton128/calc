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

type Object interface {
	Copy() Object
	ID() int
	Kind() ast.Kind
	Name() string
	Package() *Package
	Loc() string
	Pos() token.Pos
	Scope() *Scope
	SetLoc(string)
	String() string
	Type() Type
}

type object struct {
	id    int
	kind  ast.Kind
	name  string
	pkg   *Package
	pos   token.Pos
	loc   string
	scope *Scope
	typ   Type
}

func (o object) copy(id int) object {
	return object{
		id:    id,
		kind:  o.kind,
		name:  o.name,
		pkg:   o.pkg,
		pos:   o.pos,
		loc:   o.loc,
		scope: o.Package().scope,
		typ:   o.typ,
	}
}

func (o object) ID() int           { return o.id }
func (o object) Kind() ast.Kind    { return o.kind }
func (o object) Name() string      { return o.name }
func (o object) Package() *Package { return o.pkg }
func (o object) Pos() token.Pos    { return o.pos }
func (o object) Loc() string       { return o.loc }
func (o object) Scope() *Scope {
	if o.scope == nil {
		return o.pkg.Scope()
	}
	return o.scope
}
func (o *object) SetLoc(loc string) { o.loc = loc }
func (o object) Type() Type         { return o.typ }
func (o object) String() string {
	if o.id != 0 {
		return fmt.Sprintf("%s%d", o.name, o.id)
	}
	return o.Name()
}
