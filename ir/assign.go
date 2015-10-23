package ir

type Assignment struct {
	name string
	Rhs  Object
}

func (a *Assignment) Name() string { return a.name }
func (a *Assignment) Type() Type   { return a.Rhs.Type() } /* no? */
