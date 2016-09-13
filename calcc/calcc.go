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

	"github.com/rthornton128/calc/cgen"
)

var binExt = ""    // binary extension
var cgenExt = ".c" // code generation extension; C is the default backend
var cflags = []string{"-c", "-g"}
var lflags = []string{}

func cleanup(filename string) {
	os.Remove(filename + cgenExt)
	os.Remove(filename + ".o")
}

func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

func stripExt(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

func main() {
	if runtime.GOOS == "windows" {
		binExt = ".exe"
	}
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Calc Compiler Tool Version 2.1")
		fmt.Fprintln(os.Stderr, "\nUsage of:", os.Args[0])
		fmt.Fprintln(os.Stderr, os.Args[0], "[flags] <filename>")
		flag.PrintDefaults()
	}

	var (
		interOnly   = flag.Bool("s", false, "compile intermediate code")
		backEnd     = flag.String("t", "c", "target: amd64, c, x86")
		compileOnly = flag.Bool("c", false, "compile to object code, no linking")
		optimize    = flag.Bool("opt", true, "run optimization pass")
	)
	flag.Parse()

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
	// path is either file with .calc ext or directory name

	var c cgen.CodeGenerator
	switch *backEnd {
	case "amd64":
		c = &cgen.Amd64{}
		cgenExt = ".s"
	case "c":
		cflags = append(cflags, "-std=gnu99")
		c = &cgen.StdC{}
	case "x86":
		c = &cgen.X86{}
		cgenExt = ".s"
	default:
		fmt.Println("invalid code generator backend selected")
		os.Exit(1)
	}

	// writer should output .c or .s file
	w, err := os.Create(stripExt(path) + cgenExt)
	if err != nil {
		fatal(err)
	}
	fi, err := os.Stat(path) // needs raw path
	if err != nil {
		fatal(err)
	}
	if fi.IsDir() {
		err = cgen.CompileDir(w, c, path, *optimize)
		path = filepath.Join(path, filepath.Base(path))
	} else {
		err = cgen.CompileFile(w, c, path, *optimize)
	}
	if err != nil {
		cleanup(stripExt(path))
		fatal(err)
	}

	// stop processessing if only producing intermediate code
	if *interOnly {
		os.Exit(0)
	}

	path = stripExt(path)
	cflags = append(cflags, fmt.Sprintf("--output=%s%s", path, ".o"))
	cflags = append(cflags, stripExt(path)+cgenExt)
	out, err := exec.Command("gcc", cflags...).CombinedOutput()
	if err != nil {
		cleanup(path)
		fatal(string(out), err)
	}

	if *compileOnly {
		os.Exit(0)
	}

	lflags = append(lflags, fmt.Sprintf("--output=%s%s", path, binExt))
	lflags = append(lflags, stripExt(path)+".o")
	out, err = exec.Command("gcc", lflags...).CombinedOutput()
	if err != nil {
		cleanup(path)
		fatal(string(out), err)
	}
	cleanup(path)
}
