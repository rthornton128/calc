package ir

import (
	"fmt"

	"github.com/rthornton128/calc/ast"
)

type Declaration struct {
	object
	Body   Object
	Params []string
}

func MakeDeclaration(d *ast.DeclExpr, parent *Scope) *Declaration {
	scope := NewScope(parent)
	params := make([]string, len(d.Params))
	for i, p := range d.Params {
		params[i] = p.Name
		scope.Insert(MakeParam(p, scope))
	}

	return &Declaration{
		object: newObject(d.Name.Name, d.Type.Name, d.Pos(), ast.FuncDecl, scope),
		Params: params,
		Body:   MakeExpr(d.Body, scope),
	}
}

func (d *Declaration) String() string {
	var out string
	for _, s := range d.Params {
		out += d.scope.Lookup(s).String()
	}
	return fmt.Sprintf("decl {%s %s (%s) %s}", d.name, d.typ, out, d.Body)
}

type Param struct {
	object
	id int
}

func MakeParam(p *ast.Ident, parent *Scope) *Param {
	return &Param{object: newObject(p.Name, p.Type.Name, p.Pos(), ast.VarDecl,
		parent)}
}

func (p *Param) ID() int      { return p.id }
func (p *Param) SetID(id int) { p.id = id }
func (p *Param) String() string {
	return fmt.Sprintf("{%s %s}", p.name, p.typ)
}
