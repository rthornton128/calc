// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir

import (
	"fmt"
	"strings"

	"github.com/rthornton128/calc/ast"
)

type Var struct{ object }

func makeVar(pkg *Package, i *ast.Ident) *Var {
	return &Var{
		object: object{
			name:  i.Name,
			pos:   i.Pos(),
			scope: pkg.scope,
		},
	}
}

func (i *Var) String() string {
	return i.Name()
}

type Variable struct {
	object
	Params []*Param
	Body   []Object
}

func makeVariable(pkg *Package, ve *ast.VarExpr) *Variable {
	pkg.newScope()
	defer pkg.closeScope()

	return &Variable{
		object: object{
			id:    pkg.getID(),
			name:  "var",
			pos:   ve.Pos(),
			scope: pkg.scope,
			typ:   GetType(ve.Type),
		},
		Params: makeParamList(pkg, ve.Params),
		Body:   MakeExprList(pkg, ve.Body),
	}
}

func (v *Variable) String() string {
	params := make([]string, len(v.Params))
	for i, p := range v.Params {
		params[i] = v.Scope().Lookup(p.Name()).String()
	}

	body := make([]string, len(v.Body))
	for i, e := range v.Body {
		body[i] = e.String()
	}

	return fmt.Sprintf("var:%s (%s) {%s}", v.Type(), strings.Join(params, ","),
		strings.Join(body, ","))
}
