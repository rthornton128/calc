package comp

import (
	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

func validType(t string) bool {
	return t == "int"
}

func typeOf(n ast.Node, s *ast.Scope) string {
	t := "unknown"
	switch e := n.(type) {
	case *ast.BasicLit:
		t = typeOfBasic(e)
	case *ast.BinaryExpr, *ast.UnaryExpr:
		t = "int"
	case *ast.CallExpr:
		ob := s.Lookup(e.Name.Name)
		t = typeOfObject(ob)
	case *ast.DeclExpr:
		ob := s.Lookup(e.Name.Name)
		t = typeOfObject(ob)
	case *ast.Ident:
		ob := s.Lookup(e.Name)
		t = typeOfObject(ob)
	case *ast.IfExpr:
		if e.Type != nil {
			t = e.Type.Name
		}
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

func typeOfObject(o *ast.Object) string {
	t := "unknown"
	if o.Type != nil {
		if validType(o.Type.Name) {
			t = o.Type.Name
		}
	}
	return t
}
