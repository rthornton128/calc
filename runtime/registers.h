/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

#ifndef RT_REGISTERS_H
#define RT_REGISTERS_H

#include <stdint.h>

extern intptr_t ax; /* accumulator register */
extern intptr_t dx; /* data register */
extern uintptr_t *bp; /* base pointer */
extern uintptr_t *sp; /* stack pointer */

void reg_init();

#endif
