// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package ir

import (
	"fmt"

	"github.com/rthornton128/calc/ast"
)

type Package struct {
	object
	top *Scope
}

func MakePackage(pkg *ast.Package, name string) *Package {
	scope := NewScope(nil)
	p := &Package{
		object: object{name: name, pos: pkg.Pos(), scope: scope},
		top:    scope,
	}
	for _, f := range pkg.Files {
		MakeFile(p, f)
	}
	return p
}

func (p *Package) getID() int {
	p.id++
	return p.id
}

func (p *Package) String() string {
	return fmt.Sprintf("package %s {%s}", p.name, p.scope)
}

func MakeFile(pkg *Package, f *ast.File) {
	for _, d := range f.Defs {
		pkg.scope.Insert(MakeDefine(pkg, d))
	}
}
