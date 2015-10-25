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

func makeIf(i *ast.IfExpr, parent *Scope) *If {
	scope := newScope(parent)
	return &If{
		object: newObject("if", i.Type.Name, i.Pos(), None, scope),
		Cond:   makeExpr(i.Cond, parent),
		Then:   makeExpr(i.Then, scope),
		Else:   makeExpr(i.Else, scope),
	}
}

func (i *If) ID() int      { return i.id }
func (i *If) SetID(id int) { i.id = id }
func (i *If) String() string {
	return fmt.Sprintf("{if %s then %s else %s}", i.Cond, i.Then, i.Else)
}
