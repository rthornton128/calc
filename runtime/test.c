/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

#include "cmp.h"
#include "instructions.h"
#include "registers.h"
#include "stack.h"

#include <assert.h>
#include <stdio.h>

void cmp_tests() {
	int a = 1, b = 2;
	
	/* greater tests */
	gtl((char *)&a, (char *)&b); assert(*(int32_t *)eax == 0);
	gtl((char *)&b, (char *)&a); assert(*(int32_t *)eax == 1);
	gel((char *)&a, (char *)&b); assert(*(int32_t *)eax == 0);
	gel((char *)&a, (char *)&a); assert(*(int32_t *)eax == 1);
	gel((char *)&b, (char *)&a); assert(*(int32_t *)eax == 1);
	a = 1;
	/* less tests */
	ltl((char *)&a, (char *)&b); assert(*(int32_t *)eax == 1);
	ltl((char *)&b, (char *)&a); assert(*(int32_t *)eax == 0);
	lel((char *)&a, (char *)&b); assert(*(int32_t *)eax == 1);
	lel((char *)&a, (char *)&a); assert(*(int32_t *)eax == 1);
	lel((char *)&b, (char *)&a); assert(*(int32_t *)eax == 0);

	/* (not) equal tests */
	eql((char *)&a, (char *)&b); assert(*(int32_t *)eax == 0);
	eql((char *)&b, (char *)&b); assert(*(int32_t *)eax == 1);
	nel((char *)&a, (char *)&b); assert(*(int32_t *)eax == 1);
	nel((char *)&a, (char *)&a); assert(*(int32_t *)eax == 0);

	/* and/or tests */
	a = 0, b = 1;
	andl((char *)&a, (char *)&a); assert(*(int32_t *)eax == 0);
	andl((char *)&a, (char *)&b); assert(*(int32_t *)eax == 0);
	andl((char *)&b, (char *)&b); assert(*(int32_t *)eax == 1);
	orl((char *)&a, (char *)&a); assert(*(int32_t *)eax == 0);
	orl((char *)&a, (char *)&b); assert(*(int32_t *)eax == 1);
	orl((char *)&b, (char *)&b); assert(*(int32_t *)eax == 1);
}

void instructions_tests() {
	/* 32 bit copy */
	setl(42, eax);
	assert(*(int32_t *)eax == 42);
	setl(1000000, eax);
	assert(*(int32_t *)eax == 1000000);
	movl(eax, edx);
	assert(*(int32_t *)edx == 1000000);

	/* addition */
	setl(3, eax);
	setl(5, edx);
	addl(edx, eax);
	assert(*(int32_t *)eax == 8);

	/* division */
	setl(15, eax);
	setl(5, edx);
	divl(edx, eax);
	assert(*(int32_t *)eax == 3);

	/* multiplication */
	setl(3, eax);
	setl(5, edx);
	mull(edx, eax);
	assert(*(int32_t *)eax == 15);

	/* subtraction */
	setl(3, eax);
	setl(5, edx);
	subl(edx, eax);
	assert(*(int32_t *)eax == -2);
}

void stack_tests() {
	stack_init();
	//printf("%p, %p\n", ebp, esp);
	enter(16);
	//printf("%p, %p\n", ebp, esp);
	setl(24, ebp+0);
	setl(18, ebp+4);

	/* simulate another function call */
	enter(16);
	//printf("%p, %p\n", ebp, esp);
	setl(5, ebp+0);
	setl(3, ebp+4);
	movl(ebp+4, eax);
	addl(ebp+0, eax);
	leave();
	//printf("%p, %p\n", ebp, esp);
	assert(*(int32_t *)eax == 8);
	/* end inner function */

	movl(ebp+4, eax);
	addl(ebp+0, eax);
	leave();
	//printf("%p, %p\n", ebp, esp);
	stack_end();
	assert(*(int32_t *)eax == 42);
}

int main() {
	cmp_tests();
	instructions_tests();
	stack_tests();

	return 0;
}
