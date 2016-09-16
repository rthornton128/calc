// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

import "fmt"

// TODO amd64 only
var paramRegisters = []string{"%rcx", "%rdx", "%r8", "%r9"}

// TODO amd64 only
func (s *stack64) ArgumentLoc(i int) string {
	if i < len(paramRegisters) {
		return paramRegisters[i]
	}
	i = i - (len(paramRegisters) - 1)
	return fmt.Sprintf("%d(%%rsp)", 32+(i*8))
}
