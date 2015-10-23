package ir

type Object interface {
	Name() string
	Type() Type
	Scope() *Scope
	String() string
}

type IDer interface {
	ID() int
	SetID(int)
}

type object struct {
	name  string
	typ   Type
	scope *Scope
}

func newObject(name, t string, s *Scope) object {
	return object{
		name:  name,
		typ:   typeFromString(t),
		scope: s,
	}
}

func (o object) Name() string  { return o.name }
func (o object) Type() Type    { return o.typ }
func (o object) Scope() *Scope { return o.scope }
