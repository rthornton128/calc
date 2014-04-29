#ifndef RT_INSTRUCTIONS_H
#define RT_INSTRUCTIONS_H

#include <stdint.h>

void enter(const int32_t n);
void leave(void);

void movl(const char *src, char *dest);
void setl(const int32_t n, char *dest);

void addl(const char *src, char *dest);
void divl(const char *src, char *dest);
void mull(const char *src, char *dest);
void reml(const char *src, char *dest);
void subl(const char *src, char *dest);

#endif
