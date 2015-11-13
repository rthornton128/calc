package ir

import (
	"fmt"
	"strings"

	"github.com/rthornton128/calc/ast"
)

type Var struct{ object }

func makeVar(i *ast.Ident, parent *Scope) *Var {
	return &Var{object: newObject(i.Name, "", i.Pos(), ast.VarDecl, parent)}
}

func (i *Var) String() string {
	return i.Name()
}

type Variable struct {
	object
	id     int
	Params []string
	Body   []Object
}

func makeVariable(ve *ast.VarExpr, parent *Scope) *Variable {
	v := &Variable{
		object: newObject("", ve.Type.Name, ve.Pos(), ast.VarDecl, parent),
		Params: makeParamList(ve.Params, parent),
		Body:   MakeExprList(ve.Body, parent),
	}

	parent.Insert(v)

	return v
}

func (v *Variable) String() string {
	params := make([]string, len(v.Params))
	for i, p := range v.Params {
		params[i] = v.Scope().Lookup(p).String()
	}

	body := make([]string, len(v.Body))
	for i, e := range v.Body {
		body[i] = e.String()
	}

	return fmt.Sprintf("var:%s (%s) {%s}", v.Type(), strings.Join(params, ","),
		strings.Join(body, ","))
}

func (v *Variable) ID() int      { return v.id }
func (v *Variable) SetID(id int) { v.id = id }
