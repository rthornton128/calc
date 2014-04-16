REM Copyright (c) 2014, Rob Thornton
REM All rights reserved.
REM This source code is governed by a Simplied BSD-License. Please see the
REM LICENSE included in this distribution for a copy of the full license
REM or, if one is not included, you may also find a copy at
REM http://opensource.org/licenses/BSD-2-Clause

@ECHO OFF
SET CC=gcc.exe
SET CFLAGS=-Wall -g -Wextra -Werror -fmax-errors=10 -std=c99

calcc.exe %*
IF %ERRORLEVEL% EQU 0 %CC% %CFLAGS% %~n1.c -o %~n1.exe & DEL %~n1.c
