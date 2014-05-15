/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

#ifndef RT_INSTRUCTIONS_H
#define RT_INSTRUCTIONS_H

#include <stdint.h>

void enter(const int32_t n);
void leave(void);
void popl(char *dest);
void pushl(const char *src);

void movl(const char *src, char *dest);
void setl(const int32_t n, char *dest);

void addl(const char *src, char *dest);
void divl(const char *src, char *dest);
void mull(const char *src, char *dest);
void reml(const char *src, char *dest);
void subl(const char *src, char *dest);

#endif
