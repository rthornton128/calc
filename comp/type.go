package comp

import (
	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

func validType(t string) bool {
	return t == "int"
}

func (c *compiler) typeOf(n ast.Node) string {
	t := "unknown"
	switch e := n.(type) {
	case *ast.BasicLit:
		t = typeOfBasic(e)
	case *ast.BinaryExpr:
		t = "int"
	case *ast.CallExpr:
		ob := c.curScope.Lookup(e.Name.Name)
		t = typeOfObject(ob)
	case *ast.DeclExpr:
		ob := c.curScope.Lookup(e.Name.Name)
		t = typeOfObject(ob)
	case *ast.Ident:
		ob := c.curScope.Lookup(e.Name)
		t = typeOfObject(ob)
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
