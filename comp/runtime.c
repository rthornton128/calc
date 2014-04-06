#include "runtime.h"

int eax = 0;
int edx = 0;

int stack[1024]; /* 4k stack (assuming int is 32 bits) */
int sindex = 0;

void push(int n) {
	stack[sindex++] = n;
}

int pop() {
	return stack[--sindex];
}
