// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rthornton128/calc1/comp"
)

var calcExt = ".calc"

func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	filename := flag.Arg(0)

	if filepath.Ext(filename) != calcExt {
		fatal("Calc source files should have the '.calc' extension")
	}

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fatal(err)
	}
	fmt.Println("Compiling:", filename)

	comp.CompileFile(filename[:len(filename)-len(calcExt)], string(src))
}
