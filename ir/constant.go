// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir

import (
	"fmt"
	"strconv"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

type Value interface {
	String() string
	Type() Type
}

type (
	boolValue bool
	intValue  int64
)

type Constant struct {
	object
	value Value
}

func makeConstant(b *ast.BasicLit) *Constant {
	var v Value
	switch b.Kind {
	case token.BOOL:
		v, _ = makeBool(b.Lit) // TODO handle error
	case token.INTEGER:
		v, _ = makeInt(b.Lit) // TODO handle error
	}
	return &Constant{
		object: object{name: v.String(), pos: b.Pos(), typ: v.Type()},
		value:  v,
	}
}

func (c *Constant) String() string {
	return c.value.String()
}

func makeBool(lit string) (Value, error) {
	b, err := strconv.ParseBool(lit)
	return boolValue(b), err
}

func (v boolValue) String() string { return fmt.Sprintf("%v", bool(v)) }
func (v boolValue) Type() Type     { return Bool }

func makeInt(lit string) (Value, error) {
	i, err := strconv.ParseInt(lit, 0, 64)
	return intValue(i), err
}

func (v intValue) String() string { return strconv.FormatInt(int64(v), 10) }
func (v intValue) Type() Type     { return Int }
