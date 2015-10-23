package ir

type Call struct {
	Name string
	t    Type
	Args []Object
}

func (c *Call) Type() Type { return c.t }
