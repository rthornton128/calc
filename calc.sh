#!/usr/bin/env bash
# Copyright (c) 2014, Rob Thornton
# All rights reserved.
# This source code is governed by a Simplied BSD-License. Please see the
# LICENSE included in this distribution for a copy of the full license
# or, if one is not included, you may also find a copy at
# http://opensource.org/licenses/BSD-2-Clause

# set this to suite your C compiler
CC=gcc
CFLAGS="-g -Wall -Wextra -Werror -fmax-errors=10 -std=c99"
CALCC="$GOPATH/bin/calcc"

filename=$(basename "$1")
filename="${filename%.*}"

$CALCC $1 && $CC $CFLAGS -o $filename ""$filename".c" && rm -f ""$filename".c"
