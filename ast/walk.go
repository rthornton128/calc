// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

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
