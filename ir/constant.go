package ir

import (
	"fmt"
	"strconv"

	"github.com/rthornton128/calc/ast"
)

type Value interface {
	Type() Type
	String() string
}

type (
	boolValue bool
	intValue  int64
)

func (v boolValue) String() string { return fmt.Sprintf("%v", bool(v)) }
func (v boolValue) Type() Type     { return Bool }

func makeBool(i *ast.Ident) (Value, error) {
	b, err := strconv.ParseBool(i.Name)
	return boolValue(b), err
}

func (v intValue) String() string { return strconv.FormatInt(int64(v), 10) }
func (v intValue) Type() Type     { return Int }

func makeInt(b *ast.BasicLit) (Value, error) {
	i, err := strconv.ParseInt(b.Lit, 0, 64)
	return intValue(i), err
}
