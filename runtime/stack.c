#include "registers.h"

#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#define MIN_STACK (size_t) 4096

char *ss = NULL;
uint32_t bpi = 0;
uint32_t spi = 0;

void stack_init() {
	ss = malloc(MIN_STACK);
	if (ss == NULL) {
		fprintf(stderr, "Failed to init stack: out of memory\n");
		exit(EXIT_FAILURE);
	}
	memset(ss, 0, MIN_STACK);
	bpi = 0;
	spi = 0;
	ebp = &ss[bpi];
	esp = &ss[spi];
}

void stack_end() {
	free(ss);
	ss = NULL;
}

