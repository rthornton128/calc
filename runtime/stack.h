/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

#ifndef RT_STACK_H
#define RT_STACK_H

void stack_init(void);
void stack_end(void);
void stack_check_overflow(register const int n);
void stack_check_underflow(register const int n);

#endif
