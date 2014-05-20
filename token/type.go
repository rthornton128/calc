package token

type Type int

var types = [...]string{
	"int",
}

func ValidType(typename string) bool {
	for _, t := range types {
		if typename == t {
			return true
		}
	}
	return false
}
