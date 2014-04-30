#ifndef RT_STACK_H
#define RT_STACK_H

#include <stdint.h>

extern char *ss;
extern uint32_t bpi; /* base pointer index */
extern uint32_t spi; /* stack pointer index */

void stack_init();
void stack_end();

#endif
