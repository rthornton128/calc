// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

// Package cgen comprises the code generation portion of the Calc
// programming language
package cgen

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/ir"
	"github.com/rthornton128/calc/parse"
	"github.com/rthornton128/calc/token"
)

type CodeGenerator interface {
	CGen(io.Writer, *ir.Package)
}

type Emitter interface {
	Emit(...interface{})
	Emitf(string, ...interface{})
}

type Writer struct{ io.Writer }

func (w *Writer) Emit(a ...interface{}) {
	fmt.Fprintln(w, a...)
}

func (w *Writer) Emitf(format string, a ...interface{}) {
	fmt.Fprintf(w, format+"\n", a...)
}

// CompileFile generates a C source file for the corresponding file
// specified by path. The .calc extension for the filename in path is
// replaced with .c for the C source output.
func CompileFile(w io.Writer, c CodeGenerator, path string, opt bool) error {
	fset := token.NewFileSet()
	f, err := parse.ParseFile(fset, path, "")
	if err != nil {
		return err
	}

	pkg := ir.MakePackage(&ast.Package{
		Files: []*ast.File{f},
	}, filepath.Base(path))

	if err := ir.TypeCheck(pkg, fset); err != nil {
		return err
	}
	if opt {
		pkg = ir.FoldConstants(pkg).(*ir.Package)
	}

	pkg.ReplaceMacros(pkg)
	//fmt.Println("package after:", pkg)

	c.CGen(w, pkg)

	return nil
}

// CompileDir generates C source code for the Calc sources found in the
// directory specified by path. The C source file uses the same name as
// directory rather than any individual file.
func CompileDir(w io.Writer, c CodeGenerator, path string, opt bool) error {
	fset := token.NewFileSet()
	p, err := parse.ParseDir(fset, path)
	if err != nil {
		return err
	}

	pkg := ir.MakePackage(p, filepath.Base(path))
	if err := ir.TypeCheck(pkg, fset); err != nil {
		return err
	}
	if opt {
		pkg = ir.FoldConstants(pkg).(*ir.Package)
	}

	c.CGen(w, pkg)

	return nil
}
