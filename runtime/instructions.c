#include "instructions.h"
#include "registers.h"
#include "stack.h"

#include <stdio.h> /* printf */
#include <string.h>

/* stack */
void enter(const int32_t n) {
	setl(ebp-&ss[0], esp);
	esp += 4;
	ebp = esp;
	esp += n;
}

void leave() {
	esp = ebp;
	esp -= 4;
	ebp = &ss[0]+(int)*esp;
}

/* memory */
void movl(const char *src, char *dest) { memmove(dest, src, sizeof (int32_t)); }
void setl(const int32_t n, char *dest) { movl((char *)&n, dest); }

/* arithmatic */
void addl(const char *src, char *dest) { *(int32_t *)dest += *(int32_t *)src; }
void divl(const char *src, char *dest) { *(int32_t *)dest /= *(int32_t *)src; }
void mull(const char *src, char *dest) { *(int32_t *)dest *= *(int32_t *)src; }
void reml(const char *src, char *dest) { *(int32_t *)dest %= *(int32_t *)src; }
void subl(const char *src, char *dest) { *(int32_t *)dest -= *(int32_t *)src; }
