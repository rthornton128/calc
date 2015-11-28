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
	Cond Object
	Then Object
	Else Object
}

func makeIf(pkg *Package, ie *ast.IfExpr) *If {
	i := &If{
		object: object{
			id:    pkg.getID(),
			name:  "if-then",
			pkg:   pkg,
			pos:   ie.Pos(),
			scope: pkg.scope,
			typ:   GetType(ie.Type),
		},
		Cond: MakeExpr(pkg, ie.Cond),
		Then: MakeExpr(pkg, ie.Then),
	}
	if ie.Else != nil {
		i.object.name += "-else"
		i.Else = MakeExpr(pkg, ie.Else)
	}
	return i
}

func (i *If) String() string {
	return fmt.Sprintf("{if[%s] %s then %s else %s}", i.typ, i.Cond, i.Then,
		i.Else)
}
