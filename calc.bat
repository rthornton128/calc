@ECHO OFF
SET CC=gcc.exe
SET CFLAGS=-Wall -g -Wextra -Werror -fmax-errors=10 -std=c99

calcc.exe %*
IF %ERRORLEVEL% EQU 0 %CC% %CFLAGS% %~n1.c -o %~n1.exe & DEL %~n1.c
