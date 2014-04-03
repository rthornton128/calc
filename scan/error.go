package scan

import (
	"fmt"

	"github.com/rthornton128/calc/token"
)

type Error struct {
	pos token.Position
	msg string
}

func (e Error) Error() string {
	return fmt.Sprint(e.pos, " ", e.msg)
}

type ErrorList struct {
	errors []*Error
}

func (el *ErrorList) ErrorCount() int {
	return len(el.errors)
}

func (el *ErrorList) Add(p token.Position, msg string) {
	el.errors = append(el.errors, &Error{p, msg})
}
