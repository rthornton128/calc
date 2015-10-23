package ir

import "github.com/rthornton128/calc/ast"

type Declaration struct {
	object
	Params []*Param
}

func MakeDeclaration(d *ast.DeclExpr, parent *Scope) *Declaration {
	decl := &Declaration{
		object: object{
			name:   d.Name.Name,
			t:      TypeFromString(d.Type.Name),
			parent: newScope(parent),
		},
		Params: make([]*Param, len(d.Params)),
	}

	for i, p := range d.Params {
		decl.Params[i] = MakeParam(p, parent)
	}

	return decl
}

type Param struct {
	object
	ID int
}

func MakeParam(p *ast.Ident, parent *Scope) *Param {
	return &Param{
		object: object{
			name:   p.Name,
			t:      TypeFromString(p.Object.Type.Name), // TODO shudder
			parent: parent,
		},
	}
}
