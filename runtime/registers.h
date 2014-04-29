#ifndef RT_REGISTERS_H
#define RT_REGISTERS_H

extern char *eax; /* accumulator register */
extern char *edx; /* data register */
extern char *ebp; /* base pointer */
extern char *esp; /* stack pointer */

void reg_init();

#endif
