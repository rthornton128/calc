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

#define check_overflow(n) {\
	if (sp + n >= (uintptr_t *)&ss[scap-1]) {\
		fprintf(stderr, "panic: stack underflow!\n");\
		exit(EXIT_FAILURE);\
	}\
}

#define check_underflow(n) {\
	if (sp - n < (uintptr_t *) ss) {\
		fprintf(stderr, "panic: stack underflow!\n");\
		exit(EXIT_FAILURE);\
	}\
}

#define enter(n) {\
	check_overflow(n);\
	*sp = (uintptr_t) bp;\
	bp = ++sp;\
	sp += n;\
}

#define leave() {\
	check_underflow((sp - bp) / sizeof (uintptr_t));\
	sp = bp;\
	bp = (uintptr_t *) *--sp;\
}

#define pop(dest) {\
	check_underflow(1);\
	dest = (uintptr_t) *--sp;\
}

#define push(src) {\
	check_overflow(1);\
	*sp++ = (uintptr_t) src;\
}

#endif
