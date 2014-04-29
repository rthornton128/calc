#include "registers.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MIN_STACK (size_t) 4096

char *ss = NULL;

void stack_init() {
	ss = malloc(MIN_STACK);
	if (ss == NULL) {
		fprintf(stderr, "Failed to init stack: out of memory\n");
		exit(EXIT_FAILURE);
	}
	memset(ss, 0, MIN_STACK);
	esp = &ss[0];
	ebp = &ss[0];
}

void stack_end() {
	free(ss);
	ss = NULL;
}

