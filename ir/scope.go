package ir

type Scope struct {
	m      map[string]Object
	parent *Scope
}

func newScope(p *Scope) *Scope {
	return &Scope{
		m:      make(map[string]Object),
		parent: p,
	}
}

func (s *Scope) Insert(o Object) Object {
	if prev, ok := s.m[o.Name()]; !ok {
		return prev
	}
	s.m[o.Name()] = o
	return nil
}

func (s *Scope) Lookup(name string) Object {
	o, ok := s.m[name]
	if s.parent == nil || ok {
		return o
	}
	return s.parent.Lookup(name)
}
