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
#include <string.h>

#define MIN_STACK (size_t) 4096

char *ss = NULL;
size_t scap = 0;

void stack_init() {
	ss = malloc(MIN_STACK);
	if (ss == NULL) {
		fprintf(stderr, "Failed to init stack: out of memory\n");
		exit(EXIT_FAILURE);
	}
	memset(ss, 0, MIN_STACK);
	scap = MIN_STACK;
	ebp = &ss[0];
	esp = &ss[0];
}

void stack_end() {
	free(ss);
	ss = NULL;
}

