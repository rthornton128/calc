package token

import "fmt"

type Pos uint

var illegalPos = Pos(0)

func (p Pos) Valid() bool {
	return p != illegalPos
}

type Position struct {
	Filename string
	Col, Row int
}

func (p Position) String() string {
	if p.Filename == "" {
		return fmt.Sprintf("%d:%d", p.Col, p.Row)
	}
	return fmt.Sprintf("%s:%d:%d", p.Filename, p.Col, p.Row)
}
