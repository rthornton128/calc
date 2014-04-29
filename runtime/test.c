#include "cmp.h"

#include <assert.h>
#include <stdio.h>

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

int main() {
	cmp_tests();

	return 0;
}
