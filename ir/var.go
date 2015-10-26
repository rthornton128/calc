package ir

import (
	"fmt"
	"reflect"

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
	Assign Object
}

func makeVariable(ve *ast.VarExpr, parent *Scope) *Variable {
	var assign Object = nil
	if ve.Value != nil && !reflect.ValueOf(ve.Value).IsNil() {
		assign = makeExpr(ve.Value, parent)
	}

	var typ string
	if ve.Name.Type != nil {
		typ = ve.Name.Type.Name
	}

	v := &Variable{
		object: newObject(ve.Name.Name, typ, ve.Pos(), ast.VarDecl, parent),
		Assign: assign,
	}

	if prev := parent.Insert(v); prev != nil {
	}

	return v
}

func (v *Variable) String() string {
	if v.Assign != nil {
		return fmt.Sprintf("var %s = %s", v.Name(), v.Assign.String())
	}
	return fmt.Sprintf("var %s", v.Name())
}

func (v *Variable) ID() int      { return v.id }
func (v *Variable) SetID(id int) { v.id = id }
