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
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rthornton128/calc/cgen"
)

var binExt = ""    // binary extension
var cgenExt = ".c" // code generation extension; C is the default backend
var objExt = ".o"  // object extension

func cleanup(filename string) {
	os.Remove(filename + cgenExt)
	os.Remove(filename + objExt)
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
	fmt.Fprintln(os.Stderr, "Calc Compiler Tool Version 2.1")
}

func main() {
	if runtime.GOOS == "windows" {
		binExt = ".exe"
	}
	flag.Usage = func() {
		printVersion()
		fmt.Fprintln(os.Stderr, "\nUsage of:", os.Args[0])
		fmt.Fprintln(os.Stderr, os.Args[0], "[flags] <filename>")
		flag.PrintDefaults()
	}
	var (
		asm  = flag.Bool("s", false, "output intermediate code only")
		be   = flag.String("cgen", "c", "code generator: c, x86")
		cc   = flag.String("cc", "gcc", "C compiler to use (C backend only)")
		cfl  = flag.String("cflags", "-c -g -std=gnu99", "C compiler flags")
		cout = flag.String("cout", "--output=", "C compiler output flag")
		ld   = flag.String("ld", "gcc", "linker")
		ldf  = flag.String("ldflags", "", "linker flags")
		opt  = flag.Bool("o", true, "run optimization pass")
		ver  = flag.Bool("v", false, "Print version number and exit")
	)
	flag.Parse()

	if *ver {
		printVersion()
		os.Exit(1)
	}
	var path string
	switch flag.NArg() {
	case 0:
		path, _ = filepath.Abs(".")
	case 1:
		path, _ = filepath.Abs(flag.Arg(0))
	default:
		flag.Usage()
		os.Exit(1)
	}

	fi, err := os.Stat(path)
	if err != nil {
		fatal(err)
	}

	opath := path[:len(path)-len(filepath.Ext(path))]
	w, err := os.Create(opath + cgenExt) // BUG cgenExt not changed yet...
	if err != nil {
		fatal(err)
	}
	var c cgen.CodeGenerator
	switch *be {
	case "c":
		c = &cgen.StdC{w}
	case "x86":
		c = &cgen.X86{w}
		cgenExt = ".s"
	default:
		fmt.Println("invalid code generator backend selected")
		os.Exit(1)
	}
	if fi.IsDir() {
		err = cgen.CompileDir(c, path, *opt)
		path = filepath.Join(path, filepath.Base(path))
	} else {
		err = cgen.CompileFile(c, path, *opt)
	}

	path = opath
	if err != nil {
		cleanup(path)
		fatal(err)
	}
	if !*asm {
		/* compile to object code */
		var out []byte
		args := make_args(*cfl, *cout+path+objExt, path+cgenExt)
		out, err := exec.Command(*cc+binExt,
			strings.Split(args, " ")...).CombinedOutput()
		if err != nil {
			cleanup(path)
			fatal(string(out), err)
		}

		/* link to executable */
		args = make_args(*ldf, *cout+path+binExt, path+objExt)
		out, err = exec.Command(*ld+binExt,
			strings.Split(args, " ")...).CombinedOutput()
		if err != nil {
			cleanup(path)
			fatal(string(out), err)
		}
		cleanup(path)
	}
}
