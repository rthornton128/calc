package ir

import (
	"fmt"
	"strings"

	"github.com/rthornton128/calc/ast"
)

type For struct {
	object
	Cond Object
	Body []Object
}

func makeFor(pkg *Package, f *ast.ForExpr) *For {
	body := make([]Object, len(f.Body))
	for i, e := range f.Body {
		body[i] = MakeExpr(pkg, e)
	}
	return &For{
		object: object{id: pkg.getID(), name: "for", pos: f.Pos(), scope: pkg.scope,
			typ: GetType(f.Type)},
		Cond: MakeExpr(pkg, f.Cond),
		Body: body,
	}
}

func (f *For) String() string {
	body := make([]string, len(f.Body))
	for i, o := range f.Body {
		body[i] = o.String()
	}
	return fmt.Sprintf("{for[%s] %s {%s}}", f.Type(), f.Cond, strings.Join(body, ","))
}
