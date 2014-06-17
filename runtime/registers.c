/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

#include <stddef.h>

char regs[8];
char *eax = &regs[0];
char *edx = &regs[4];
char *ebp = NULL;
char *esp = NULL;
