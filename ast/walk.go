// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ast

type Func func(Node)

func Walk(node Node, f Func) {
	if node == nil {
		panic("Node is nil!")
	}

	if f != nil {
		f(node)
	}
	switch n := node.(type) {
	case *AssignExpr:
		Walk(n.Name, f)
		if n.Object.Type != nil {
			Walk(n.Object.Type, f)
		}
		Walk(n.Object.Value, f)
	case *BinaryExpr:
		for _, v := range n.List {
			Walk(v, f)
		}
	case *CallExpr:
		Walk(n.Name, f)
		for _, v := range n.Args {
			Walk(v, f)
		}
	case *DeclExpr:
		Walk(n.Name, f)
		Walk(n.Type, f)
		for _, v := range n.Params {
			Walk(v, f)
		}
		Walk(n.Body, f)
	case *ExprList:
		for _, v := range n.List {
			Walk(v, f)
		}
	case *File:
		for _, v := range n.Scope.table {
			Walk(v.Type, f)
			Walk(v.Value, f)
		}
	case *IfExpr:
		Walk(n.Cond, f)
		Walk(n.Then, f)
		Walk(n.Else, f)
	case *Package:
		for _, v := range n.Scope.table {
			Walk(v.Type, f)
			Walk(v.Value, f)
		}
	case *VarExpr:
		Walk(n.Name, f)
		if n.Object.Type != nil {
			Walk(n.Object.Type, f)
		}
		if n.Object.Value != nil {
			Walk(n.Object.Value, f)
		}
	}
}
