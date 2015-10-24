package ir

import (
	"fmt"

	"github.com/rthornton128/calc/ast"
)

type Package struct {
	object
	name  string
	scope *Scope
}

func MakePackage(pkg *ast.Package, name string) *Package {
	p := &Package{name: name, scope: newScope(nil)}
	for _, f := range pkg.Files {
		MakeFile(f, p.scope)
	}
	return p
}

func (p *Package) Name() string  { return p.name }
func (p *Package) Scope() *Scope { return p.scope }
func (p *Package) String() string {
	return fmt.Sprintf("package: %s { %s }", p.name, p.scope)
}

func MakeFile(f *ast.File, parent *Scope) {
	for _, d := range f.Decls {
		parent.Insert(MakeDeclaration(d, parent))
	}
}
