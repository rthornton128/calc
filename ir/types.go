package ir

type Type int

const (
	Unknown Type = iota
	Bool
	Int
)

func TypeFromString(name string) Type {
	switch name {
	case "int":
		return Int
	case "bool":
		return Bool
	default:
		return Unknown
	}
}
