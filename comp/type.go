package comp

import (
	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

func validType(t string) bool {
	return t == "int"
}

func typeOf(n ast.Node) string {
	t := "unkown"
	switch e := n.(type) {
	case *ast.BasicLit:
		t = typeOfBasic(e)
	case *ast.BinaryExpr:
		t = "int"
	case *ast.DeclExpr:
		if e.Type != nil {
			t = typeOfIdent(e.Type)
		}
	case *ast.Ident:
		t = typeOfIdent(e)
	}
	return t
}

func typeOfBasic(b *ast.BasicLit) string {
	switch b.Kind {
	case token.INTEGER:
		return "int"
	default:
		return "unknown"
	}
}

func typeOfIdent(i *ast.Ident) string {
	if validType(i.Name) {
		return i.Name
	}
	return "unknown"
}
