package ir

import (
	"fmt"
	"strings"

	"github.com/rthornton128/calc/ast"
)

type Block struct {
	object
	exprs []Object
}

func makeBlock(l *ast.ExprList, parent *Scope) *Block {
	list := make([]Object, len(l.List))
	for i, e := range l.List {
		list[i] = makeExpr(e, parent)
	}
	return &Block{
		object: newObject("block", "", parent),
		exprs:  list,
	}
}

func (b *Block) String() string {
	var out []string
	for _, e := range b.exprs {
		out = append(out, e.String())
	}
	return fmt.Sprintf("{%s}", strings.Join(out, ","))
}
