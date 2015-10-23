package ir

import "go/token"

type Binary struct {
	ID  int
	Lhs Object
	Op  token.Token
	Rhs Object
	t   Type
}

func (b *Binary) Type() Type { return b.t }

type Unary struct {
	Op  token.Token
	Rhs Object
	t   Type
}

func (u *Unary) Type() Type { return u.t }
