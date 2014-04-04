package ast

type Func func(Node)

func Walk(node Node, f Func) {
	if node == nil {
		panic("Node is nil!")
	}

	if f != nil {
		f(node)
	}
	switch n := node.(type) {
	case *File:
		Walk(n.Root, f)
	case *BinaryExpr:
		for _, v := range n.List {
			Walk(v, f)
		}
	}
}
