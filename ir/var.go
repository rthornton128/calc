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

// Var should really be Ident but isn't for some reason
type Var struct{ object }

func makeVar(pkg *Package, i *ast.Ident) *Var {
	return &Var{
		object: object{
			id:    pkg.getID(),
			kind:  ast.VarDecl,
			name:  i.Name,
			pkg:   pkg,
			pos:   i.Pos(),
			scope: pkg.scope,
		},
	}
}

// Copy makes a deep copy of the Var object
func (v *Var) Copy() Object {
	return &Var{object: v.object.copy(v.Package().getID())}
}

/*
func (v *Var) Name() string {
	return fmt.Sprintf("%s%d", v.name, v.id)
}*/

func (i *Var) String() string {
	return fmt.Sprintf("%s[%d]", i.Name(), i.ID())
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
			pkg:   pkg,
			pos:   ve.Pos(),
			scope: pkg.scope,
			typ:   typeFromString(ve.Type.Name),
		},
		Params: makeParamList(pkg, ve.Params),
		Body:   MakeExprList(pkg, ve.Body),
	}

	return v
}

// Copy is a deep copy of (var [params] expr...)
func (v *Variable) Copy() Object {
	v.Package().newScope()
	defer v.Package().closeScope()

	nVariable := &Variable{
		object: v.object.copy(v.Package().getID()),
		//Params: params,
		//Body: v.Body.Copy(),
	}
	nVariable.scope = v.Package().scope //NewScope(v.Package().Scope().parent)
	nVariable.Params = make([]*Param, len(v.Params))
	for i, p := range v.Params {
		nVariable.Params[i] = p.Copy().(*Param)
		nVariable.Scope().Insert(nVariable.Params[i], p.Name())
	}
	for _, e := range v.Body {
		//fmt.Println("copying", e)
		nVariable.Body = append(nVariable.Body, e.Copy())
	}
	return nVariable
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

	return fmt.Sprintf("var[%d]:%s (%s) {%s}", v.ID(), v.Type(), strings.Join(params, ","),
		strings.Join(body, ","))
}
