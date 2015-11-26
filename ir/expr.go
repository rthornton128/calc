// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir

import "github.com/rthornton128/calc/ast"

func MakeExpr(pkg *Package, e ast.Expr) Object {
	switch t := e.(type) {
	case *ast.AssignExpr:
		return makeAssignment(pkg, t)
	case *ast.BasicLit:
		return makeConstant(t)
	case *ast.BinaryExpr:
		return makeBinary(pkg, t)
	case *ast.CallExpr:
		return makeCall(pkg, t)
	case *ast.ForExpr:
		return makeFor(pkg, t)
	case *ast.FuncExpr:
		return makeFunc(pkg, t)
	case *ast.Ident:
		return makeVar(pkg, t)
	case *ast.IfExpr:
		return makeIf(pkg, t)
	case *ast.UnaryExpr:
		return makeUnary(pkg, t)
	case *ast.VarExpr:
		return makeVariable(pkg, t)
	default:
		panic("unreachable")
	}
}

func MakeExprList(pkg *Package, el []ast.Expr) []Object {
	ol := make([]Object, 0)
	for _, e := range el {
		ol = append(ol, MakeExpr(pkg, e))
	}
	return ol
}
