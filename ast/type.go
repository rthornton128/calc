package ast

import (
	"fmt"

	"github.com/rthornton128/calc/token"
)

type Type int

const (
	Invalid Type = iota
	Int
)

var types = []string{
	Invalid: "invalid",
	Int:     "int",
}

func typeLookup(i *Ident) Type {
	for x := range types {
		if types[x] == i.Name {
			return Type(x)
		}
	}
	return Invalid
}

type TypeChecker struct{ i Type }

func (t *TypeChecker) Visit(n Node) bool {
	//fmt.Printf("%d-%#v\n", t.i, n)
	//t.i++
	switch x := n.(type) {
	case *Package, *File:
		/* do nothing */
	case *BasicLit:
		switch x.Kind {
		case token.INTEGER:
			t.i = Int
		}
	case *BinaryExpr:
	case *UnaryExpr:
		Walk(x.Value, t)
		// TODO should emit proper warning
		if t.i != Int {
			fmt.Println("unary expression expects integer value")
		}
	case *VarExpr:
		x.Object.RealType = typeLookup(x.Object.Type)
	}

	return true
}
