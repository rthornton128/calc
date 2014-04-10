// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package scan

import (
	"fmt"

	"github.com/rthornton128/calc1/token"
)

type Error struct {
	pos token.Position
	msg string
}

func (e Error) Error() string {
	return fmt.Sprint(e.pos, " ", e.msg)
}

type ErrorList []*Error

func (el ErrorList) Count() int {
	return len(el)
}

func (el *ErrorList) Add(p token.Position, msg string) {
	*el = append(*el, &Error{p, msg})
}

func (el *ErrorList) Print() {
	for _, err := range *el {
		fmt.Println(err)
	}
}
