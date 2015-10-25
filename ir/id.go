package ir

import "fmt"

func Tag(pkg *Package) {
	for _, v := range pkg.Scope().m {
		tag(v, 1)
	}
}

func tag(o Object, nextID int) int {
	switch t := o.(type) {
	case *Assignment:
		nextID = tag(t.Rhs, nextID)
	case *Binary:
		nextID = setID(t, nextID)
		fmt.Println("set binary", t, "to id", t.ID())
		nextID = tag(t.Lhs, nextID)
		nextID = tag(t.Rhs, nextID)
	case *Block:
		for _, e := range t.Exprs {
			nextID = tag(e, nextID)
		}
	case *Call:
		for _, arg := range t.Args {
			nextID = tag(arg, nextID)
		}
	case *Declaration:
		for _, p := range t.Params {
			nextID = setID(t.Scope().Lookup(p).(IDer), nextID)
		}
		nextID = tag(t.Body, nextID)
	case *If:
		nextID = setID(t, nextID)
		nextID = tag(t.Cond, nextID)
		nextID = tag(t.Then, nextID)
		if t.Else != nil {
			nextID = tag(t.Else, nextID)
		}
	case *Param:
		nextID = setID(t.Scope().Lookup(t.Name()).(IDer), nextID)
	case *Unary:
		nextID = tag(t.Rhs, nextID)
	case *Variable:
		nextID = setID(t.Scope().Lookup(t.Name()).(IDer), nextID)
	}
	return nextID
}

func setID(o IDer, nextID int) int {
	if o.ID() == 0 {
		o.SetID(nextID)
		nextID++
	}
	return nextID
}
