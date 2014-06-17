/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

#ifndef RT_REGISTERS_H
#define RT_REGISTERS_H

extern char *eax; /* accumulator register */
extern char *edx; /* data register */
extern char *ebp; /* base pointer */
extern char *esp; /* stack pointer */

void reg_init();

#endif
