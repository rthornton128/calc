// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package cgen

import "fmt"

type stack32 struct{}

func (s *stack32) Size() int { return 4 }

func (s *stack32) ArgumentLoc(i int) string {
	return fmt.Sprintf("%d(%%esp)", (i+1)*s.Size())
}

func (s *stack32) CallStackOffset(i int) string {
	return fmt.Sprintf("%d(%%ebp)", (i*4)+8)
}

func (s *stack32) ParameterLoc(i int) string {
	return s.CallStackOffset(i)
}
