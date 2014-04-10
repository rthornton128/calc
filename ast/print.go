// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

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
