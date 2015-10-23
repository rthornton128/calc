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
	case token.INTEGER:
		v, _ = makeInt(b) // TODO handle error
		//case token.BOOL:
		//v, _ := makeBool(b) // TODO handle error
	}
	return &Constant{
		object: newObject(v.String(), v.Type().String(), parent),
		value:  v,
	}
}

func (c *Constant) String() string {
	return c.value.String()
}

func makeBool(i *ast.Ident) (Value, error) {
	b, err := strconv.ParseBool(i.Name)
	return boolValue(b), err
}

func (v boolValue) String() string { return fmt.Sprintf("%v", bool(v)) }
func (v boolValue) Type() Type     { return Bool }

func makeInt(b *ast.BasicLit) (Value, error) {
	i, err := strconv.ParseInt(b.Lit, 0, 64)
	return intValue(i), err
}

func (v intValue) String() string { return strconv.FormatInt(int64(v), 10) }
func (v intValue) Type() Type     { return Int }
