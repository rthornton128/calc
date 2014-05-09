/* Copyright (c) 2014, Rob Thornton
 * All rights reserved.
 * This source code is governed by a Simplied BSD-License. Please see the
 * LICENSE included in this distribution for a copy of the full license
 * or, if one is not included, you may also find a copy at
 * http://opensource.org/licenses/BSD-2-Clause */

#include "cmp.h"

#include <string.h>

int gel(char *a, char *b) { return gtl(a, b) || eql(a, b); }
int gtl(char *a, char *b) { return memcmp(a, b, 4) == 1; }
int lel(char *a, char *b) { return ltl(a, b) || eql(a, b); }
int ltl(char *a, char *b) { return memcmp(a, b, 4) == -1; }
int eql(char *a, char *b) { return memcmp(a, b, 4) == 0; }
int nel(char *a, char *b) { return memcmp(a, b, 4) != 0; }
