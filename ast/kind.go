package ast

type Kind int

const (
	None Kind = iota + 1
	FuncDecl
	VarDecl
)
