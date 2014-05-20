// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package token

import (
	"fmt"
)

type Error struct {
	pos Position
	msg string
}

func (e Error) Error() string {
	return fmt.Sprint(e.pos, " ", e.msg)
}

type ErrorList []*Error

func (el ErrorList) Count() int {
	return len(el)
}

func (el *ErrorList) Add(p Position, args ...interface{}) {
	*el = append(*el, &Error{pos: p, msg: fmt.Sprint(args...)})
}

func (el *ErrorList) Print() {
	for _, err := range *el {
		fmt.Println(err)
	}
}
