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
}

func MakePackage(pkg *ast.Package, name string) *Package {
	p := &Package{object: newObject(name, "", pkg.Pos(), ast.None, NewScope(nil))}
	for _, f := range pkg.Files {
		MakeFile(f, p.Scope())
	}
	return p
}

func (p *Package) String() string {
	return fmt.Sprintf("package: %s {%s}", p.name, p.scope)
}

func MakeFile(f *ast.File, parent *Scope) {
	for _, d := range f.Defs {
		parent.Insert(MakeDefine(d, parent))
	}
}
