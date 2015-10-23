package ir

type If struct {
	t     Type
	Cond  Object
	Then  Object
	Else  Object
	Scope *Scope
}

func (i *If) Type() Type {
	return i.t
}
