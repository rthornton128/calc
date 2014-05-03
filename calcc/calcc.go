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
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rthornton128/calc/comp"
)

var calcExt = ".calc"

func cleanup(filename string) {
	os.Remove(filename + ".c")
	os.Remove(filename + ".o")
}

func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

func make_args(options ...string) string {
	var args string
	for i, opt := range options {
		if len(opt) > 0 {
			args += opt
			if i < len(options)-1 {
				args += " "
			}
		}
	}
	return args
}

func printVersion() {
	fmt.Fprintln(os.Stderr, "Calc 1 Compiler Tool Version 1.1")
}

func main() {
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	flag.Usage = func() {
		printVersion()
		fmt.Fprintln(os.Stderr, "\nUsage of:", os.Args[0])
		fmt.Fprintln(os.Stderr, os.Args[0], "[flags] <filename>")
		flag.PrintDefaults()
	}
	var (
		cc   = flag.String("cc", "gcc", "C compiler to use")
		cfl  = flag.String("cflags", "-c -std=gnu99", "C compiler flags")
		cout = flag.String("cout", "--output=", "C compiler output flag")
		ld   = flag.String("ld", "gcc", "linker")
		ldf  = flag.String("ldflags", "", "linker flags")
		ver  = flag.Bool("version", false, "Print version number and exit")
	)
	flag.Parse()

	if *ver {
		printVersion()
		os.Exit(1)
	}
	if flag.NArg() != 1 {
		flag.Usage()
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

	filename = filename[:len(filename)-len(calcExt)]
	comp.CompileFile(filename, string(src))

	defer cleanup(filename)

	/* compile to object code */
	var out []byte
	args := make_args(*cfl, *cout+filename+".o", filename+".c")
	out, err = exec.Command(*cc+ext, strings.Split(args, " ")...).CombinedOutput()
	if err != nil {
		fatal(string(out))
	}

	/* link to executable */
	args = make_args(*ldf, *cout+filename+ext, filename+".o")
	out, err = exec.Command(*ld+ext, strings.Split(args, " ")...).CombinedOutput()
	if err != nil {
		fatal(string(out))
	}
}
