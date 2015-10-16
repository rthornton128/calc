package ast

import "github.com/rthornton128/calc/token"

type Type int

const (
	Invalid Type = iota
	Int
)

var types = []string{
	Invalid: "invalid",
	Int:     "int",
}

func (t Type) String() string {
	return types[t]
}

func typeLookup(i *Ident) Type {
	for x := range types {
		if types[x] == i.Name {
			return Type(x)
		}
	}
	return Invalid
}

type TypeChecker struct {
	ErrorHandler token.ErrorHandler
	i            Type
	scope        *Scope
}

func (t *TypeChecker) Visit(n Node) bool {
	switch x := n.(type) {
	case *Package:
		/* TODO should switch to package scope? */
		/* do nothing */
	case *File:
		t.scope = x.Scope
		for _, d := range x.Decls {
			o := x.Scope.Lookup(d.Name.Name)
			o.RealType = typeLookup(d.Type)
		}
	case *BasicLit:
		switch x.Kind {
		case token.INTEGER:
			t.i = Int
		}
	case *BinaryExpr:
		if x.RealType != Invalid {
			t.i = x.RealType
			break
		}

		for _, e := range x.List {
			Walk(e, t)
			if t.i != Int {
				t.ErrorHandler(x.Pos(),
					"Binary expr expects operands of type Int, Got:", t.i)
			}
		}
		x.RealType = Int
		t.i = Int
	case *CallExpr:
		// already checked this call, no need to do it again
		if x.RealType != Invalid {
			t.i = x.RealType
			break
		}

		o := t.scope.Lookup(x.Name.Name)
		if o.Kind != Decl {
			t.ErrorHandler(x.Name.Pos(), "call of non-function", x.Name.Name)
		}
		d := o.Value.(*DeclExpr)
		for i, a := range x.Args {
			Walk(a, t)
			if t.i != d.Params[i].Object.RealType {
				t.ErrorHandler(a.Pos(),
					"argument ", i, " of ", x.Name.Name, " wrong type; expected ",
					t.i, " got ", d.Params[i].Object.RealType)
			}
		}
		x.RealType = d.RealType
		t.i = x.RealType
	case *DeclExpr:
		t.scope = x.Scope
		x.RealType = typeLookup(x.Type)
		for _, p := range x.Params {
			p.Object.RealType = typeLookup(p.Object.Type)
		}
	case *Ident:
		if x.Object == nil {
			x.Object = t.scope.Lookup(x.Name)
		}
		if x.Object.RealType == Invalid {
			x.Object.RealType = typeLookup(x.Object.Type)
		}
		t.i = x.Object.RealType
	case *IfExpr:
		t.scope = x.Scope
		x.RealType = typeLookup(x.Type)
	case *UnaryExpr:
		Walk(x.Value, t)
		if t.i != Int {
			t.ErrorHandler(x.Pos(), " unary expression expects integer value")
		}
	case *VarExpr:
		// TODO bug: does not infer type when type is not explicitly set
		x.Object.RealType = typeLookup(x.Object.Type)
	}

	return true
}
