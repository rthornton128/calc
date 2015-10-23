package ir

import (
	"fmt"
	"strings"

	"github.com/rthornton128/calc/ast"
)

type Call struct {
	object
	args []Object
}

func makeCall(c *ast.CallExpr, parent *Scope) *Call {
	args := make([]Object, len(c.Args))
	for i, a := range c.Args {
		args[i] = makeExpr(a, parent)
	}
	return &Call{
		object: newObject(c.Name.Name, "", parent),
		args:   args,
	}
}

func (c *Call) String() string {
	var out []string
	for _, a := range c.args {
		out = append(out, a.String())
	}
	return fmt.Sprintf("{call: %s (%s)}", c.Name(), strings.Join(out, ","))
}
