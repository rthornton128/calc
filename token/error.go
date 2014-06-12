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

// Error represents an error in the source code. It consists of a position
// within the source files and message text describing the error
type Error struct {
	pos Position
	msg string
}

// Error generates an error string to satisfy the error interface
func (e Error) Error() string {
	return fmt.Sprint(e.pos, " ", e.msg)
}

// ErrorList is a slice of Error pointers
type ErrorList []*Error

// Count returns the number of errors within the list
func (el ErrorList) Count() int {
	return len(el)
}

// Add a new error the list at the given position p.
func (el *ErrorList) Add(p Position, args ...interface{}) {
	*el = append(*el, &Error{pos: p, msg: fmt.Sprint(args...)})
}

// Error returns a string containing all the errors in the error list
func (el ErrorList) Error() string {
	var msg string
	for _, err := range el {
		msg += fmt.Sprintln(err)
	}
	return msg
}

// Print outputs a message containing all the errors in the list
func (el ErrorList) Print() {
	for _, err := range el {
		fmt.Println(err)
	}
}
