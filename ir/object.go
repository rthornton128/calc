package ir

type Object interface {
	Name() string
	Type() Type
}

type object struct {
	name   string
	t      Type
	parent *Scope
}

func (o *object) Name() string  { return o.name }
func (o *object) Type() Type    { return o.t }
func (o *object) Scope() *Scope { return o.parent }
