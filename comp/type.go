package comp

import (
	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

func validType(t *ast.Ident) bool {
	return t.Name == "int"
}

func typeOf(n ast.Node, s *ast.Scope) (t *ast.Ident) {
	switch e := n.(type) {
	case *ast.AssignExpr:
		t = typeOfObject(s.Lookup(e.Name.Name))
	case *ast.BasicLit:
		t = typeOfBasic(e)
	case *ast.BinaryExpr:
		t = typeOf(e.List[0], s)
	case *ast.CallExpr:
		t = typeOfObject(s.Lookup(e.Name.Name))
	case *ast.DeclExpr:
		t = typeOfObject(s.Lookup(e.Name.Name))
	case *ast.ExprList:
		// BUG: should follow chain of execution to make sure all return
		// values match return type
		t = typeOf(e.List[len(e.List)-1], s)
	case *ast.Ident:
		t = typeOfObject(s.Lookup(e.Name))
	case *ast.IfExpr:
		if e.Type != nil {
			t = e.Type
		}
	case *ast.UnaryExpr:
		t = typeOf(e.Value, s)
	case *ast.VarExpr:
		t = typeOf(e.Name, s)
	}

	if t == nil {
		t = &ast.Ident{Name: "unknown", NamePos: n.Pos()}
	}
	return t
}

func typeOfBasic(b *ast.BasicLit) *ast.Ident {
	switch b.Kind {
	case token.INTEGER:
		return &ast.Ident{Name: "int", NamePos: b.Pos()}
	default:
		return nil
	}
}

func typeOfObject(o *ast.Object) *ast.Ident {
	if o.Type != nil {
		return o.Type
	}
	return nil
}
