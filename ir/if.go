package ir

import (
	"fmt"

	"github.com/rthornton128/calc/ast"
)

type If struct {
	object
	id   int
	Cond Object
	Then Object
	Else Object
}

func makeIf(ie *ast.IfExpr, parent *Scope) *If {
	scope := newScope(parent)
	i := &If{
		object: newObject("if", ie.Type.Name, ie.Pos(), ast.None, scope),
		Cond:   makeExpr(ie.Cond, parent),
		Then:   makeExpr(ie.Then, scope),
	}
	if ie.Else != nil {
		i.Else = makeExpr(ie.Else, scope)
	}
	return i
}

func (i *If) ID() int      { return i.id }
func (i *If) SetID(id int) { i.id = id }
func (i *If) String() string {
	return fmt.Sprintf("{if %s then %s else %s}", i.Cond, i.Then, i.Else)
}
