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

type Define struct {
	object
	Body Object
}

func MakeDefine(d *ast.DefineStmt, parent *Scope) *Define {
	scope := NewScope(parent)
	body := MakeExpr(d.Body, scope)
	t := body.Type().String()
	if d.Type != nil {
		t = d.Type.Name
	}

	return &Define{
		object: newObject(d.Name.Name, t, d.Pos(), body.Kind(), scope),
		Body:   body,
	}
}

func (d *Define) String() string {
	return fmt.Sprintf("define{%s:%s = %s}", d.Name(), d.Type(), d.Body)
}
