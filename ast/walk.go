// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ast

type Visitor interface {
	Visit(node Node) bool
}

func Walk(node Node, v Visitor) {
	if !v.Visit(node) {
		return
	}

	switch n := node.(type) {
	case *AssignExpr:
		Walk(n.Name, v)
		Walk(n.Value, v)
	case *BinaryExpr:
		for _, x := range n.List {
			Walk(x, v)
		}
	case *CallExpr:
		Walk(n.Name, v)
		for _, x := range n.Args {
			Walk(x, v)
		}
	case *DeclExpr:
		Walk(n.Name, v)
		for _, x := range n.Params {
			Walk(x, v)
		}
		Walk(n.Type, v)
		Walk(n.Body, v)
	case *ExprList:
		for _, x := range n.List {
			Walk(x, v)
		}
	case *File:
		for _, x := range n.Scope.Table {
			Walk(x.Value, v)
		}
	case *IfExpr:
		Walk(n.Cond, v)
		Walk(n.Type, v)
		Walk(n.Then, v)
		Walk(n.Else, v)
	case *Package:
		for _, x := range n.Scope.Table {
			Walk(x.Value, v)
		}
	case *UnaryExpr:
		Walk(n.Value, v)
	case *VarExpr:
		Walk(n.Name, v)
	}
}
