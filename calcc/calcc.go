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

func findRuntime() string {
	var paths []string
	rpath := "/src/github.com/rthornton128/calc/runtime"
	if runtime.GOOS == "Windows" {
		paths = strings.Split(os.Getenv("GOPATH"), ";")
	} else {
		paths = strings.Split(os.Getenv("GOPATH"), ":")
	}
	for _, path := range paths {
		path = filepath.Join(path, rpath)
		_, err := os.Stat(filepath.Join(path, "runtime.a"))
		if err == nil || os.IsExist(err) {
			return path
		}
	}
	return ""
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
	fmt.Fprintln(os.Stderr, "Calc Compiler Tool Version 2.0")
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
		asm  = flag.Bool("s", false, "generate C code but do not compile")
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

	/* do a preemptive search to see if runtime can be found. Does not
	 * guarantee it will be there at link time */
	rpath := findRuntime()
	if rpath == "" {
		fatal("Unable to find runtime in GOPATH. Make sure 'make' command was " +
			"run in source directory")
	}

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fatal(err)
	}
	fmt.Println("Compiling:", filename)

	filename = filename[:len(filename)-len(calcExt)]
	comp.CompileFile(filename, string(src))

	if !*asm {
		/* compile to object code */
		var out []byte
		args := make_args(*cfl, "-I "+rpath, *cout+filename+".o", filename+".c")
		out, err = exec.Command(*cc+ext, strings.Split(args, " ")...).CombinedOutput()
		if err != nil {
			cleanup(filename)
			fatal(string(out), err)
		}

		/* link to executable */
		args = make_args(*ldf, *cout+filename+ext, filename+".o", rpath+"/runtime.a")
		out, err = exec.Command(*ld+ext, strings.Split(args, " ")...).CombinedOutput()
		if err != nil {
			cleanup(filename)
			fatal(string(out), err)
		}
		cleanup(filename)
	}
}
