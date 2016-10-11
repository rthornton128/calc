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
)

type If struct {
	object
	Cond      Object
	Then      Object
	Else      Object
	ThenLabel string
	EndLabel  string
}

func makeIf(pkg *Package, ie *ast.IfExpr) *If {
	i := &If{
		object: object{
			id:    pkg.getID(),
			name:  "if-then",
			pkg:   pkg,
			pos:   ie.Pos(),
			scope: pkg.scope,
			typ:   typeFromString(ie.Type.Name),
		},
		Cond: MakeExpr(pkg, ie.Cond),
		Then: MakeExpr(pkg, ie.Then),
	}
	if ie.Else != nil {
		i.object.name += "-else"
		i.Else = MakeExpr(pkg, ie.Else)
	}
	i.ThenLabel = fmt.Sprintf("L%d", pkg.getID())
	i.EndLabel = fmt.Sprintf("L%d", i.ID())
	return i
}

// Copy makes a deep copy of the Unary object
func (i *If) Copy() Object {
	var e Object
	if i.Else != nil {
		e = i.Else.Copy()
	}
	id := i.Package().getID()
	//fmt.Println("if with new id:", id)
	return &If{
		object:    i.object.copy(id),
		Cond:      i.Cond.Copy(),
		Then:      i.Then.Copy(),
		Else:      e,
		ThenLabel: fmt.Sprintf("L%d", i.Package().getID()),
		EndLabel:  fmt.Sprintf("L%d", id),
	}
}

func (i *If) String() string {
	if i.Else == nil {
		return fmt.Sprintf("{if[%d]:%s (%s) then %s}", i.id, i.typ, i.Cond,
			i.Then)

	}
	return fmt.Sprintf("{if[%d]:%s (%s) then %s else %s}", i.id, i.typ, i.Cond, i.Then,
		i.Else)
}
