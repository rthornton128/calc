// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir

import (
	"fmt"
	"strings"

	"github.com/rthornton128/calc/ast"
)

// Call is a function call of the form (function [args...]) where args is
// zero or more optional arguments.
type Call struct {
	object
	Args []Object
}

func makeCall(pkg *Package, c *ast.CallExpr) *Call {
	args := make([]Object, len(c.Args))
	for i, a := range c.Args {
		args[i] = MakeExpr(pkg, a)
	}
	return &Call{
		object: object{name: c.Name.Name, pkg: pkg, pos: c.Pos(), scope: pkg.scope},
		Args:   args,
	}
}

// Copy makes a deep copy of the Call object
func (c *Call) Copy() Object {
	args := make([]Object, len(c.Args))
	for i, a := range c.Args {
		args[i] = a.Copy()
	}
	return &Call{
		object: c.object.copy(0),
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
