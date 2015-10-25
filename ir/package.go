package ir

import (
	"fmt"

	"github.com/rthornton128/calc/ast"
)

type Package struct {
	object
}

func MakePackage(pkg *ast.Package, name string) *Package {
	p := &Package{object: newObject(name, "", pkg.Pos(), None, newScope(nil))}
	for _, f := range pkg.Files {
		MakeFile(f, p.Scope())
	}
	return p
}

func (p *Package) String() string {
	return fmt.Sprintf("package: %s {%s}", p.name, p.scope)
}

func MakeFile(f *ast.File, parent *Scope) {
	for _, d := range f.Decls {
		parent.Insert(MakeDeclaration(d, parent))
	}
}
