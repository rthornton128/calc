package ir

import (
	"fmt"

	"github.com/rthornton128/calc/ast"
)

type Define struct {
	object
	Body Object
}

func MakeDefine(d *ast.DefineStmt, parent *Scope) *Define {
	scope := NewScope(parent)
	body := MakeExpr(d.Body, scope)
	t := body.Type().String()
	if d.Type != nil {
		t = d.Type.Name
	}

	return &Define{
		object: newObject(d.Name.Name, t, d.Pos(), body.Kind(), scope),
		Body:   body,
	}
}

func (d *Define) String() string {
	return fmt.Sprintf("define{%s:%s = %s}", d.Name(), d.Type(), d.Body)
}
