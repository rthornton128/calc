// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

// Package comp comprises the code generation portion of the Calc
// programming language
package cgen

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/ir"
	"github.com/rthornton128/calc/parse"
	"github.com/rthornton128/calc/token"
)

type compiler struct {
	writer io.WriteCloser
	fset   *token.FileSet
	errors token.ErrorList
}

// CompileFile generates a C source file for the corresponding file
// specified by path. The .calc extension for the filename in path is
// replaced with .c for the C source output.
func CompileFile(path string, opt bool) error {
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

	c := compiler{}
	path = path[:len(path)-len(filepath.Ext(path))]
	switch {
	default:
		var err error
		c.writer, err = os.Create(path + ".c")
		if err != nil {
			return err
		}
		defer c.writer.Close()

		cc := &cCompiler{c}
		cc.compPackage(pkg)
	}

	if c.errors.Count() != 0 {
		return c.errors
	}
	return nil
}

// CompileDir generates C source code for the Calc sources found in the
// directory specified by path. The C source file uses the same name as
// directory rather than any individual file.
func CompileDir(path string, opt bool) error {
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

	c := compiler{fset: fset}
	path = path[:len(path)-len(filepath.Ext(path))]
	switch {
	default:
		var err error
		c.writer, err = os.Create(path + ".c")
		if err != nil {
			return err
		}
		defer c.writer.Close()

		cc := &cCompiler{c}
		cc.compPackage(pkg)
	}

	if c.errors.Count() != 0 {
		return c.errors
	}
	return nil
}

// Error adds an error to the compiler at the given position. The remaining
// arguments are used to generate the error message.
func (c *compiler) Error(pos token.Pos, args ...interface{}) {
	c.errors.Add(c.fset.Position(pos), args...)
}

func (c *compiler) emit(s string, args ...interface{}) {
	fmt.Fprintf(c.writer, s, args...)
}

func (c *compiler) emitln(args ...interface{}) {
	fmt.Fprintln(c.writer, args...)
}
