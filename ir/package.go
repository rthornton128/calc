package ir

import "github.com/rthornton128/calc/ast"

type Package struct {
	scope *Scope
}

func MakePackage(pkg *ast.Package) *Package {
	p := &Package{scope: newScope(nil)}
	for _, f := range pkg.Files {
		MakeFile(f, p.scope)
	}
	return p
}

func MakeFile(f *ast.File, parent *Scope) {
	for _, d := range f.Decls {
		parent.Insert(MakeDeclaration(d, parent))
	}
}
