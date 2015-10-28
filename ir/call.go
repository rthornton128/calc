package ir

import (
	"fmt"
	"strings"

	"github.com/rthornton128/calc/ast"
)

type Call struct {
	object
	Args []Object
}

func makeCall(c *ast.CallExpr, parent *Scope) *Call {
	args := make([]Object, len(c.Args))
	for i, a := range c.Args {
		args[i] = MakeExpr(a, parent)
	}
	return &Call{
		object: newObject(c.Name.Name, "", c.Pos(), ast.None, parent),
		Args:   args,
	}
}

func (c *Call) String() string {
	var out []string
	for _, a := range c.Args {
		out = append(out, a.String())
	}
	return fmt.Sprintf("{call: %s (%s)}", c.Name(), strings.Join(out, ","))
}
