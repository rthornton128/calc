#ifndef RT_STACK_H
#define RT_STACK_H

#include <stddef.h>

extern char *ss;
extern size_t scap;

void stack_init();
void stack_end();

#endif
