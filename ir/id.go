package ir

func Tag(o Object) {
	var nextID int
	if pkg, ok := o.(*Package); ok {
		for _, v := range pkg.Scope().m {
			tag(v, &nextID)
		}
	} else {
		tag(o, &nextID)
	}
}

func tag(o Object, nextID *int) {
	switch t := o.(type) {
	case *Assignment:
		tag(t.Rhs, nextID)
	case *Binary:
		setID(t, nextID)
		tag(t.Lhs, nextID)
		tag(t.Rhs, nextID)
	case *Block:
		for _, e := range t.Exprs {
			tag(e, nextID)
		}
	case *Call:
		for _, arg := range t.Args {
			tag(arg, nextID)
		}
	case *Declaration:
		for _, p := range t.Params {
			setID(t.Scope().Lookup(p).(IDer), nextID)
		}
		tag(t.Body, nextID)
	case *If:
		setID(t, nextID)
		tag(t.Cond, nextID)
		tag(t.Then, nextID)
		if t.Else != nil {
			tag(t.Else, nextID)
		}
	case *Param:
		setID(t.Scope().Lookup(t.Name()).(IDer), nextID)
	case *Unary:
		tag(t.Rhs, nextID)
	case *Variable:
		setID(t.Scope().Lookup(t.Name()).(IDer), nextID)
	}
	return
}

func setID(o IDer, nextID *int) {
	if o.ID() == 0 {
		o.SetID(*nextID)
		*nextID++
	}
}
