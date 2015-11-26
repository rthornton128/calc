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

type Function struct {
	object
	Body   []Object
	Params []string
}

func makeFunc(pkg *Package, f *ast.FuncExpr) *Function {
	pkg.newScope()
	defer pkg.closeScope()

	fn := &Function{
		object: object{
			//id:    pkg.getID(),
			kind:  ast.FuncDecl,
			name:  "f",
			pkg:   pkg,
			pos:   f.Pos(),
			scope: pkg.scope,
			typ:   typeFromString(f.Type.Name)},
		Params: makeParamList(pkg, f.Params),
		Body:   MakeExprList(pkg, f.Body),
	}
	return fn
}

func (f *Function) String() string {
	params := make([]string, len(f.Params))
	for i, s := range f.Params {
		// TODO there is a case to be made to prevent lookup in parent scopes here
		params[i] = f.Scope().Lookup(s).String()
	}

	exprs := make([]string, len(f.Body))
	for i, s := range f.Body {
		exprs[i] = s.String()
	}

	return fmt.Sprintf("func:%s (%s) {%s}", f.typ, strings.Join(params, ","),
		strings.Join(exprs, ","))
}

type Param struct {
	object
}

func makeParam(pkg *Package, p *ast.Param) *Param {
	return &Param{object{
		id:   pkg.getID(),
		kind: ast.VarDecl,
		name: p.Name.Name,
		pos:  p.Pos(),
		typ:  typeFromString(p.Type.Name),
	}}
}

func (p *Param) String() string {
	return fmt.Sprintf("%s[%s]", p.name, p.typ)
}

func makeParamList(pkg *Package, pl []*ast.Param) []string {
	params := make([]string, len(pl))
	for i, p := range pl {
		params[i] = p.Name.Name
		pkg.Insert(makeParam(pkg, p))
	}
	return params
}
