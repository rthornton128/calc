#include "cmp.h"
#include "instructions.h"
#include "registers.h"
#include "stack.h"

#include <assert.h>
//#include <stdio.h>

void cmp_tests() {
	int a = 1, b = 2;
	
	/* greater tests */
	assert(gtl((char *)&a, (char *)&b) == 0);
	assert(gtl((char *)&b, (char *)&a) == 1);
	assert(gel((char *)&a, (char *)&b) == 0);
	assert(gel((char *)&a, (char *)&a) == 1);
	assert(gel((char *)&b, (char *)&a) == 1);

	/* less tests */
	assert(ltl((char *)&a, (char *)&b) == 1);
	assert(ltl((char *)&b, (char *)&a) == 0);
	assert(lel((char *)&a, (char *)&b) == 1);
	assert(lel((char *)&a, (char *)&a) == 1);
	assert(lel((char *)&b, (char *)&a) == 0);

	/* not equal tests */
	assert(nel((char *)&a, (char *)&b) == 1);
	assert(nel((char *)&a, (char *)&a) == 0);
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
