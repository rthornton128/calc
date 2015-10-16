package comp

import (
	"strconv"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

type OptConstantFolder struct{}

func (o *OptConstantFolder) Visit(n ast.Node) bool {
	switch t := n.(type) {
	case *ast.AssignExpr:
		ast.Walk(t.Value, o)
		t.Value = fold(t.Value)
	case *ast.BinaryExpr:
		for i, l := range t.List {
			ast.Walk(l, o)
			t.List[i] = fold(l)
		}
	case *ast.DeclExpr:
		ast.Walk(t.Body, o)
		t.Body = fold(t.Body)
	case *ast.IfExpr:
		ast.Walk(t.Cond, o)
		ast.Walk(t.Then, o)
		ast.Walk(t.Else, o)
		t.Cond = fold(t.Cond)
		t.Then = fold(t.Then)
		t.Else = fold(t.Else)
	case *ast.UnaryExpr:
		ast.Walk(t.Value, o)
		t.Value = fold(t.Value)
	}
	return true
}

func fold(n ast.Expr) ast.Expr {
	switch t := n.(type) {
	case *ast.BasicLit:
		if t.Kind == token.INTEGER {
			i, _ := strconv.ParseInt(t.Lit, 0, 64)
			return &ast.Value{Value: int(i), Type: t.Kind}
		}
	case *ast.BinaryExpr:
		var res int
		for i, l := range t.List {
			if v, ok := l.(*ast.Value); !ok {
				return n
			} else {
				if i == 0 {
					res = v.Value.(int)
					continue
				}
				switch t.Op {
				case token.ADD:
					res += v.Value.(int)
				case token.SUB:
					res -= v.Value.(int)
				case token.MUL:
					res *= v.Value.(int)
				case token.QUO:
					res /= v.Value.(int)
				case token.REM:
					res %= v.Value.(int)
				}
			}
		}
		// TODO get type from t.RealType or...something
		return &ast.Value{Value: res, Type: token.INTEGER}
	case *ast.UnaryExpr:
		if v, ok := t.Value.(*ast.Value); ok {
			return &ast.Value{Value: v.Value.(int) * -1, Type: v.Type}
		}
	}
	return n
}
