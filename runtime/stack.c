/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

#include "registers.h"

#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>

#define MIN_STACK (size_t) 1024

uintptr_t *ss = NULL; /* stack segment */
size_t scap = 0; /* stack capacity */

void
stack_init(void)
{
	ss = calloc(sizeof (intptr_t), MIN_STACK);
	if (ss == NULL) {
		fprintf(stderr, "panic: failed to init stack\n");
		exit(EXIT_FAILURE);
	}
	scap = MIN_STACK;
	bp = ss;
	sp = ss;
}

void
stack_end(void)
{
	bp = NULL;
	sp = NULL;
	free(ss);
	ss = NULL;
}

void
stack_check_overflow(register const int n)
{
	if (sp + n >= (uintptr_t *)&ss[scap-1]) {
		fprintf(stderr, "panic: stack overflow!\n");
		exit(EXIT_FAILURE);
	}
}

void
stack_check_underflow(register const int n)
{
	if (sp - n < (uintptr_t *) ss) {
		fprintf(stderr, "panic: stack underflow!\n");
		exit(EXIT_FAILURE);
	}
}
