package token

type Token int

const (
	tok_start Token = iota

	EOF
	ILLEGAL
	COMMENT

	lit_start
	INTEGER
	lit_end

	op_start
	LPAREN
	RPAREN
	ADD
	SUB
	MUL
	QUO
	REM
	op_end

	tok_end
)

var tok_strings = map[Token]string{
	EOF:     "EOF",
	ILLEGAL: "Illegal",
	COMMENT: "Comment",
	INTEGER: "Integer",
	LPAREN:  "(",
	RPAREN:  ")",
	ADD:     "+",
	SUB:     "-",
	MUL:     "*",
	QUO:     "/",
	REM:     "%",
}

func (t Token) IsLiteral() bool {
	return t > lit_start && t < lit_end
}

func (t Token) IsOperator() bool {
	return t > op_start && t < op_end
}

func Lookup(str string) Token {
	for t, s := range tok_strings {
		if s == str {
			return t
		}
	}
	return ILLEGAL
}

func (t Token) String() string {
	return tok_strings[t]
}

func (t Token) Valid() bool {
	return t > tok_start && t < tok_end
}
