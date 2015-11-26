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
	pkg.scope = NewScope(pkg.scope)
	defer func() { pkg.scope = pkg.scope.parent }()

	return &Function{
		object: object{
			id:    pkg.getID(),
			pkg:   pkg,
			pos:   f.Pos(),
			kind:  ast.FuncDecl,
			scope: pkg.scope,
			typ:   typeFromString(f.Type.Name)},
		Params: makeParamList(pkg, f.Params),
		Body:   MakeExprList(pkg, f.Body),
	}
}

func (d *Function) String() string {
	params := make([]string, len(d.Params))
	for i, s := range d.Params {
		params[i] = d.scope.Lookup(s).String()
	}

	exprs := make([]string, len(d.Body))
	for i, s := range d.Body {
		exprs[i] = s.String()
	}

	return fmt.Sprintf("func:%s (%s) {%s}", d.typ, strings.Join(params, ","),
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
		typ:  typeFromString(p.Type.Name)}}
}

func (p *Param) String() string {
	return fmt.Sprintf("%s[%s]", p.name, p.typ)
}

func makeParamList(pkg *Package, pl []*ast.Param) []string {
	params := make([]string, len(pl))
	for i, p := range pl {
		params[i] = p.Name.Name
		pkg.scope.Insert(makeParam(pkg, p))
	}
	return params
}
