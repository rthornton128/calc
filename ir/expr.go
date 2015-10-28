package ir

import "github.com/rthornton128/calc/ast"

func MakeExpr(e ast.Expr, parent *Scope) Object {
	switch t := e.(type) {
	case *ast.AssignExpr:
		return makeAssignment(t, parent)
	case *ast.BasicLit:
		return makeConstant(t, parent)
	case *ast.BinaryExpr:
		return makeBinary(t, parent)
	case *ast.CallExpr:
		return makeCall(t, parent)
	case *ast.ExprList:
		return makeBlock(t, parent)
	case *ast.Ident:
		return makeVar(t, parent)
	case *ast.IfExpr:
		return makeIf(t, parent)
	case *ast.UnaryExpr:
		return makeUnary(t, parent)
	case *ast.VarExpr:
		return makeVariable(t, parent)
	default:
		panic("unreachable")
	}
}
