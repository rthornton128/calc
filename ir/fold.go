package ir

import "github.com/rthornton128/calc/token"

type folder struct {
	nextID int
	eh     token.ErrorHandler
}

func FoldConstants(pkg *Package, eh token.ErrorHandler) *Package {
	f := &folder{nextID: 1, eh: eh}

	for k, v := range pkg.scope.m {
		pkg.scope.m[k] = f.Fold(v, eh)
	}

	return pkg
}

func (f *folder) Fold(o Object, eh token.ErrorHandler) Object {
	switch t := o.(type) {
	case *Assignment:
		t.rhs = f.Fold(t.rhs, eh)
	case *Binary:
	case *Declaration:
		t.body = f.Fold(t.body, eh)
		for _, p := range t.params {
			f.setID(t.Scope().Lookup(p).(IDer))
		}
	case *Unary:
	}
	return o
}

func (f *folder) setID(i IDer) {
	if i.ID() == 0 {
		i.SetID(f.nextID)
		f.nextID++
	}
}
