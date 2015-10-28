package ir

import (
	"fmt"
	"strings"

	"github.com/rthornton128/calc/ast"
)

type Block struct {
	object
	Exprs []Object
}

func makeBlock(l *ast.ExprList, parent *Scope) *Block {
	list := make([]Object, len(l.List))
	for i, e := range l.List {
		list[i] = MakeExpr(e, parent)
	}
	return &Block{
		object: newObject("block", "", l.Pos(), ast.None, parent),
		Exprs:  list,
	}
}

func (b *Block) String() string {
	var out []string
	for _, e := range b.Exprs {
		out = append(out, e.String())
	}
	return fmt.Sprintf("{%s}", strings.Join(out, ","))
}
