# Prerequisites

 * Make sure GOPATH/bin is in your PATH environmental variable
 * C compiler; GCC is used by default (see Usage: Alternate C Compilers, below
   if you wish to use something other than GCC)
 * *(Windows Only)* Mingw or makefile program compatible with GNU
   Make. An Msys install is __not__ required.

# Install

Depending on your needs, you probably want to clone a specific version of
Calc. If you clone Calc without specifying a branch you will clone the master
branch, which is unstable and likely not what you want.

Choose the version of Calc you want to install and insert the branch name of
the corresponding version you want. In the example below, Calc 1 has been used:

 1. `mkdir -p github.com/rthornton128`
 2. `cd github.com/rthornton128`
 3. `git clone -b calc1 http://github.com/rthornton128/calc`

To install and use the compiler, change into the calc directory and run the
following command (see below if using master):

	make install

This will call 'go build' to install calcc and attempt to build the C
runtime. Edit the Makefile in the root directory if you need to change the
C compiler or tune any C compiler/linker flags.

*Note* The current tip of Calc (what will be Calc 2.1) does away with the Makefile for now. You can now simply run 'go install' on the calcc directory to install the compiler.

# Usage:

	calcc [flags] **filename**.calc

Provided no errors were reported, you should be able to run the resulting
binary.

Use the -h flag to view usage and optional flags information.

## Alternate C Compilers

If you want to use LLVM/Clang or another C compiler you will need to pass
additional flags to calcc.

 * -cc=*name or path to compiler*
 * -cflags=*C flags to compile but not link, including warning flags*
 * -cout=*flag(s) to output an object by name*
 * -ld=*name or path to linker*
 * -ldflags=*linker flags*

Check the usage for examples of the defaults to properly format the flags.
Pay special attention to the -cout flag. 

*Note* This feature has not been well tested and may exhibit bad behaviour.
