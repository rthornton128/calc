package ir

import (
	"fmt"

	"github.com/rthornton128/calc/ast"
)

type Assignment struct {
	object
	Lhs string
	Rhs Object
}

func makeAssignment(a *ast.AssignExpr, parent *Scope) *Assignment {
	return &Assignment{
		object: newObject(a.Name.Name, "", a.Pos(), ast.None, parent),
		Lhs:    a.Name.Name,
		Rhs:    MakeExpr(a.Value, parent),
	}
}

func (a *Assignment) String() string {
	return fmt.Sprintf("{%s=%s}", a.Lhs, a.Rhs.String())
}
