package ast

import (
	"fmt"
)

func Print(node Node) {
	Walk(node, print)
}

func print(node Node) {
	switch n := node.(type) {
	case *BasicLit:
		fmt.Println("BasicType:", n.LitPos, n.Lit)
	case *BinaryExpr:
		fmt.Println("BinaryExpr:", n.OpPos, n.Op)
	case *File:
		fmt.Println("File:")
	default:
		fmt.Println("dunno what I got...that can't be good")
		fmt.Println(n)
	}
}
