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

type Type interface {
	Base() BasicType
	String() string
}

func GetType(e ast.Expr) Type {
	if e == nil {
		return TypeList[Unknown]
	}
	switch t := e.(type) {
	case *ast.Ident:
		for _, typ := range TypeList {
			if t.Name == typ.name {
				return typ
			}
		}
		return TypeList[Unknown]
	case *ast.FuncType:
		params := make([]Type, len(t.Params))
		for i, p := range t.Params {
			params[i] = GetType(p.Type)
		}
		return FuncType{Params: params, Return: GetType(t.Type)}
	default:
		return TypeList[Unknown]
	}
}

type BasicKind int

const (
	Unknown BasicKind = iota
	Bool
	Int
)

type BasicType struct {
	kind  BasicKind
	name  string
	cName string
}

func (b BasicType) Base() BasicType { return b }
func (b BasicType) CType() string   { return b.cName }
func (b BasicType) Kind() BasicKind { return b.kind }
func (b BasicType) String() string  { return b.name }

var TypeList = []BasicType{
	Unknown: {Unknown, "unknown", "void"},
	Bool:    {Bool, "bool", "bool"},
	Int:     {Int, "int", "int32_t"},
}

type FuncType struct {
	Params []Type
	Return Type
}

func (f FuncType) Base() BasicType { return f.Return.Base() }
func (f FuncType) String() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = p.String()
	}
	return fmt.Sprintf("func(%s)%s", strings.Join(params, ","), f.Return)
}
