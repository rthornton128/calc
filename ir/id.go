package ir

func Tag(pkg *Package) {
	for _, v := range pkg.scope.m {
		tag(v, 1)
	}
}

func tag(o Object, nextID int) int {
	//fmt.Println("loop")
	switch t := o.(type) {
	case *Assignment:
		nextID = tag(t.rhs, nextID)
	case *Binary:
		//fmt.Println("tag binary")
		nextID = tag(t.lhs, nextID)
		nextID = tag(t.rhs, nextID)
		nextID = setID(t, nextID)
		//fmt.Println("binary", t, "=", t.ID())
	case *Block:
		for _, e := range t.exprs {
			nextID = tag(e, nextID)
		}
	case *Declaration:
		//fmt.Println("declaration")
		for _, p := range t.params {
			//fmt.Println("tag param")
			nextID = setID(t.Scope().Lookup(p).(IDer), nextID)
			//fmt.Println(p, "=", t.Scope().Lookup(p).(IDer).ID())
		}
		nextID = tag(t.body, nextID)
	case *If:
		nextID = tag(t.Cond, nextID)
		nextID = tag(t.Then, nextID)
		if t.Else != nil {
			nextID = tag(t.Else, nextID)
		}
	case *Param:
		//fmt.Println("tag param", t.Name())
		nextID = setID(t.Scope().Lookup(t.Name()).(IDer), nextID)
	case *Unary:
		nextID = tag(t.rhs, nextID)
	case *Variable:
		//fmt.Println("tag variable", t.Name())
		nextID = setID(t.Scope().Lookup(t.Name()).(IDer), nextID)
		//fmt.Println(t.Name(), "=", t.Scope().Lookup(t.Name()).(IDer).ID())
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
