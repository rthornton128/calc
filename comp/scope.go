package comp

type Scope struct {
	Parent  *Scope
	Symbols map[string]int
}

func NewScope(p *Scope) *Scope {
	return &Scope{Parent: p, Symbols: make(map[string]int)}
}

func (c *compiler) openScope() {
	c.curScope = NewScope(c.curScope)
}

func (c *compiler) closeScope() {
	c.curScope = c.curScope.Parent
}

func (s *Scope) Lookup(name string) (int, bool) {
	off, ok := s.Symbols[name]
	if ok || s.Parent == nil {
		return off, ok
	}
	return s.Parent.Lookup(name)
}
