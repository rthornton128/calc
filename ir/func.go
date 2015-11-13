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

func makeFunc(f *ast.FuncExpr, parent *Scope) *Function {
	scope := NewScope(parent)

	return &Function{
		object: newObject("", f.Type.Name, f.Pos(), ast.FuncDecl, scope),
		Params: makeParamList(f.Params, scope),
		Body:   MakeExprList(f.Body, scope),
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
	id int
}

func makeParam(p *ast.Param, parent *Scope) *Param {
	return &Param{object: newObject(p.Name.Name, p.Type.Name, p.Pos(),
		ast.VarDecl, parent)}
}

func (p *Param) ID() int      { return p.id }
func (p *Param) SetID(id int) { p.id = id }
func (p *Param) String() string {
	return fmt.Sprintf("%s:%s", p.name, p.typ)
}

func makeParamList(pl []*ast.Param, parent *Scope) []string {
	params := make([]string, len(pl))
	for i, p := range pl {
		params[i] = p.Name.Name
		parent.Insert(makeParam(p, parent))
	}
	return params
}
