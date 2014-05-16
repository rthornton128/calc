/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

#include "cmp.h"
#include "instructions.h"

#include <string.h>

void gel(char *a, char *b) { setl(memcmp(a, b, 4) != -1, b); }
void gtl(char *a, char *b) { setl(memcmp(a, b, 4) == 1, b); }
void lel(char *a, char *b) { setl(memcmp(a, b, 4) != 1, b); }
void ltl(char *a, char *b) { setl(memcmp(a, b, 4) == -1, b); }
void eql(char *a, char *b) { setl(memcmp(a, b, 4) == 0, b); }
void nel(char *a, char *b) { setl(memcmp(a, b, 4) != 0, b); }

void andl(char *a, char *b) {}
void orl(char *a, char *b) {}
