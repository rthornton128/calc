package ir

import (
	"fmt"

	"github.com/rthornton128/calc/ast"
)

type Assignment struct {
	object
	lhs string
	rhs Object
}

func makeAssignment(a *ast.AssignExpr, parent *Scope) *Assignment {
	return &Assignment{
		object: newObject(a.Name.Name, "", parent),
		lhs:    a.Name.Name,
		rhs:    makeExpr(a.Value, parent),
	}
}

func (a *Assignment) String() string {
	return fmt.Sprintf("{%s=%s}", a.lhs, a.rhs.String())
}
