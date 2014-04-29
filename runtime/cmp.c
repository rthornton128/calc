#include "cmp.h"

#include <string.h>

int gel(char *a, char *b) { return gtl(a, b) || eql(a, b); }
int gtl(char *a, char *b) { return memcmp(a, b, 4) == 1; }
int lel(char *a, char *b) { return ltl(a, b) || eql(a, b); }
int ltl(char *a, char *b) { return memcmp(a, b, 4) == -1; }
int eql(char *a, char *b) { return memcmp(a, b, 4) == 0; }
int nel(char *a, char *b) { return memcmp(a, b, 4) != 0; }
