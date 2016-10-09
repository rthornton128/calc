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

// Assignment is an object representing the (= ident expr) expression
type Assignment struct {
	object
	Lhs string
	Rhs Object
}

func makeAssignment(pkg *Package, a *ast.AssignExpr) *Assignment {
	return &Assignment{
		object: object{name: "assign", pkg: pkg, pos: a.Pos(), scope: pkg.scope},
		Lhs:    a.Name.Name,
		Rhs:    MakeExpr(pkg, a.Value),
	}
}

// Copy makes a deep copy of the Assignment object
func (a *Assignment) Copy() Object {
	return &Assignment{
		object: a.object.copy(0),
		Lhs:    a.Lhs,
		Rhs:    a.Rhs.Copy(),
	}
}

func (a *Assignment) String() string {
	return fmt.Sprintf("{%s=%s}", a.Lhs, a.Rhs.String())
}
