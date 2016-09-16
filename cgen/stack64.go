// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

import "fmt"

type stack64 struct{}

func (s *stack64) Size() int { return 8 }

func (s *stack64) ParameterLoc(i int) string {
	if i < len(paramRegisters) {
		return paramRegisters[i]
	}
	return s.CallStackOffset(i - (len(paramRegisters) - 1))
}

func (s *stack64) CallStackOffset(i int) string {
	return fmt.Sprintf("%d(%%rbp)", (i*8)+16)
}
