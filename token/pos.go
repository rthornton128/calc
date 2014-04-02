package token

type Pos uint

var illegalPos = Pos(0)

func (p Pos) Valid() bool {
	return p != illegalPos
}
