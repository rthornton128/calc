/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

//#include "cmp.h"
#include "instructions.h"
#include "registers.h"
#include "stack.h"

#include <assert.h>

void
stack_tests()
{
	stack_init();
	enter(2);
	*(bp+0) = 24;
	*(bp+1) = 18;

	/* simulate another function call */
	enter(2);
	*(bp+0) = 5;
	*(bp+1) = 3;
	ax = *(bp+1);
	ax += *(bp+0);
	leave();
	assert((int32_t)ax == 8);
	/* end inner function */

	ax = *(bp+1);
	ax += *(bp+0);
	leave();
	stack_end();
	assert((int32_t)ax == 42);
}

int main() {
	stack_tests();

	return 0;
}
