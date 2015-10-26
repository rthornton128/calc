package ir

import (
	"fmt"
	"strconv"

	"github.com/rthornton128/calc/ast"
	"github.com/rthornton128/calc/token"
)

type Value interface {
	String() string
	Type() Type
}

type (
	boolValue bool
	intValue  int64
)

type Constant struct {
	object
	value Value
}

func makeConstant(b *ast.BasicLit, parent *Scope) *Constant {
	var v Value
	switch b.Kind {
	case token.BOOL:
		v, _ = makeBool(b.Lit) // TODO handle error
	case token.INTEGER:
		v, _ = makeInt(b.Lit) // TODO handle error
	}
	return &Constant{
		object: newObject(v.String(), v.Type().String(), b.Pos(), ast.None, parent),
		value:  v,
	}
}

func (c *Constant) String() string {
	return c.value.String()
}

func makeBool(lit string) (Value, error) {
	b, err := strconv.ParseBool(lit)
	return boolValue(b), err
}

func (v boolValue) String() string { return fmt.Sprintf("%v", bool(v)) }
func (v boolValue) Type() Type     { return Bool }

func makeInt(lit string) (Value, error) {
	i, err := strconv.ParseInt(lit, 0, 64)
	return intValue(i), err
}

func (v intValue) String() string { return strconv.FormatInt(int64(v), 10) }
func (v intValue) Type() Type     { return Int }
