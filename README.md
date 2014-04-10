Prerequisites:
--------------

 * Make sure GOPATH/bin is in your PATH environmental variable
 * GCC is installed (see: compilers to use alternates)

Alternate Compilers:
====================

 If you want to use LLVM/Clang or another compiler you will need to edit
 the compiler scripts calc.sh (Linux) or calc.bat (Windows)

Install:
--------

$ git clone http://github.com/rthornton128/calc1
$ cd github.com/rthornton128/calc1/calcc

Linux:
======

$ install.sh

Windows:
=======

$ install.bat

*go get should work to get the compiler but will not pull in the build
scripts*

Usage:
------

$ calc 'filename'.calc
