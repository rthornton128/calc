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
			id:    pkg.getID(),
			kind:  ast.VarDecl,
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

	v := &Variable{
		object: object{
			id:    pkg.getID(),
			kind:  ast.VarDecl,
			name:  "var",
			pos:   ve.Pos(),
			scope: pkg.scope,
			typ:   typeFromString(ve.Type.Name),
		},
		Params: makeParamList(pkg, ve.Params),
		Body:   MakeExprList(pkg, ve.Body),
	}

	return v
}

func (v *Variable) Copy(name string, id int) *Variable {
	params := make([]*Param, len(v.Params))
	for i, p := range v.Params {
		params[i] = p.Copy()
	}
	return &Variable{
		object: object{
			id:    id,
			kind:  ast.VarDecl,
			name:  name,
			pos:   v.Pos(),
			scope: v.Scope(),
			typ:   v.Type(),
		},
		Params: params,
		Body:   v.Body,
	}

}

func (v *Variable) Name() string {
	return fmt.Sprintf("%s%d", v.name, v.id)
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
