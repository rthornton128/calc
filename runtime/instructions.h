/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

#ifndef RT_INSTRUCTIONS_H
#define RT_INSTRUCTIONS_H

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#define enter(n) {\
	stack_check_overflow(n);\
	*sp = (uintptr_t) bp;\
	bp = ++sp;\
	sp += n;\
}

#define leave() {\
	stack_check_underflow((sp - bp) / sizeof (uintptr_t));\
	sp = bp;\
	bp = (uintptr_t *) *--sp;\
}

#define pop(dest) {\
	stack_check_underflow(1);\
	dest = (uintptr_t) *--sp;\
}

#define push(src) {\
	stack_check_overflow(1);\
	*sp++ = (uintptr_t) src;\
}

#endif
