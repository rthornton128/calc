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
	case *ast.FuncExpr:
		return makeFunc(t, parent)
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

func MakeExprList(el []ast.Expr, parent *Scope) []Object {
	ol := make([]Object, 0)
	for _, e := range el {
		ol = append(ol, MakeExpr(e, parent))
	}
	return ol
}
