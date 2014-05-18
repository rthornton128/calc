/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

#include "cmp.h"
#include "registers.h"
#include "instructions.h"

#include <string.h>

void gel(char *a, char *b) { setl(memcmp(a, b, 4) >= 0, ecx); }
void gtl(char *a, char *b) { setl(memcmp(a, b, 4) >= 1, ecx); }
void lel(char *a, char *b) { setl(memcmp(a, b, 4) <= 0, ecx); }
void ltl(char *a, char *b) { setl(memcmp(a, b, 4) <= -1, ecx); }
void eql(char *a, char *b) { setl(memcmp(a, b, 4) == 0, ecx); }
void nel(char *a, char *b) { setl(memcmp(a, b, 4) != 0, ecx); }

void andl(char *a, char *b) {
       setl(*(int32_t *)a >= 1 && *(int32_t *)b >= 1, ecx);
}

void orl(char *a, char *b) {
       setl(*(int32_t *)a >= 1 || *(int32_t *)b >= 1, ecx);
}
