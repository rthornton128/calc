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

Depending on your needs, you probably want to clone a specific version of
Calc. If you clone Calc without specifying a branch you will clone the master
branch, which is unstable and likely not what you want.

Choose the version of Calc you want to install and insert the branch name of
the corresponding version you want. In the example below, Calc 1 has been used:

 # $mkdir -p github.com/rthornton128
 # #cd github.com/rthornton128
 # $git clone -b calc1 http://github.com/rthornton128/calc

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
