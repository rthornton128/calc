#!/bin/bash

# set this to suite your C compiler
CC=gcc
CFLAGS="-g -Wall -Wextra -Werror -fmax-errors=10 -std=c99"
CALCC="$GOPATH/bin/calcc"

filename=$(basename "$1")
filename="${filename%.*}"

$CALCC $1 && $CC $CFLAGS -o $filename ""$filename".c" && rm -f ""$filename".c"
