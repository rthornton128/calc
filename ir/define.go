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

func MakeDefine(pkg *Package, d *ast.DefineStmt) *Define {
	// type may be unknown, in which case it will be inferred later
	body := MakeExpr(pkg, d.Body)
	typ := body.Type()
	if d.Type != nil {
		typ = GetType(d.Type)
	}
	return &Define{
		object: object{
			pkg:  pkg,
			name: d.Name.Name,
			pos:  d.Pos(),
			typ:  typ,
		},
		Body: body,
	}
}

func (d *Define) String() string {
	return fmt.Sprintf("define %s[%s] {%s}", d.Name(), d.Type(), d.Body)
}
