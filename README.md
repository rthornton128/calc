# Prerequisites

 * Make sure GOPATH/bin is in your PATH environmental variable
 * GCC is installed (see: compilers to use alternates)

## Windows

Any command, in the following instructions, ending in '.sh' should instead 
have the '.bat' extension for Windows batch scripts.

## Alternate Compilers

If you want to use LLVM/Clang or another compiler you will need to edit
the compiler script calc.bash.

# Install

$ git clone http://github.com/rthornton128/calc1
$ cd github.com/rthornton128/calc1/calcc

Check both install.bash and calc.bash to ensure they are setup correctly.
There will be a few options at the top of each script that can be tailored
for your particular setup.

Once you are done, install calc with the following command:

$ install.sh

*go get should work to get the compiler but will not pull in the build
scripts*

# Usage:

$ calc.sh **filename**.calc

This script with invoke the calcc compiler and the C compiler on your
system to produce an executable binary.
