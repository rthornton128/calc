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

type Assignment struct {
	object
	Lhs string
	Rhs Object
}

func makeAssignment(a *ast.AssignExpr, parent *Scope) *Assignment {
	return &Assignment{
		object: newObject(a.Name.Name, "", a.Pos(), ast.None, parent),
		Lhs:    a.Name.Name,
		Rhs:    MakeExpr(a.Value, parent),
	}
}

func (a *Assignment) String() string {
	return fmt.Sprintf("{%s=%s}", a.Lhs, a.Rhs.String())
}
