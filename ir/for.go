package ir

import (
	"fmt"
	"strings"

	"github.com/rthornton128/calc/ast"
)

type For struct {
	object
	id   int
	Cond Object
	Body []Object
}

func makeFor(f *ast.ForExpr, parent *Scope) *For {
	body := make([]Object, len(f.Body))
	for i, e := range f.Body {
		body[i] = MakeExpr(e, parent)
	}
	return &For{
		object: newObject("for", f.Type.Name, f.Pos(), ast.None, parent),
		Cond:   MakeExpr(f.Cond, parent),
		Body:   body,
	}
}

func (f *For) ID() int      { return f.id }
func (f *For) SetID(id int) { f.id = id }
func (f *For) String() string {
	body := make([]string, len(f.Body))
	for i, o := range f.Body {
		body[i] = o.String()
	}
	return fmt.Sprintf("{for[%s] %s {%s}}", f.typ, f.Cond, strings.Join(body, ","))
}
