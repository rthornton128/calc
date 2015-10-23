package ir

type Type int

const (
	Unknown Type = iota
	Bool
	Int
)

var typeStrings = []string{
	Unknown: "unknown type",
	Bool:    "bool",
	Int:     "int",
}

func typeFromString(name string) Type {
	for i, s := range typeStrings {
		if name == s {
			return Type(i)
		}
	}
	return Unknown
}

func (t Type) String() string {
	return typeStrings[t]
}
