package ir

import (
	"fmt"
	"reflect"

	"github.com/rthornton128/calc/ast"
)

type Var struct{ object }

func makeVar(i *ast.Ident, parent *Scope) *Var {
	return &Var{object: newObject(i.Name, "", parent)}
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
	var t string
	if ve.Object.Type != nil {
		t = ve.Object.Type.Name
	}

	var assign Object = nil
	if ve.Object.Value != nil && !reflect.ValueOf(ve.Object.Value).IsNil() {
		assign = makeExpr(ve.Object.Value, parent)
	}

	v := &Variable{
		object: newObject(ve.Name.Name, t, parent),
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
