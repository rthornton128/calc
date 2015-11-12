// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package token

type Token int

const (
	tok_start Token = iota

	EOF
	ILLEGAL

	lit_start
	BOOL
	IDENT
	INTEGER
	lit_end

	op_start
	LPAREN
	RPAREN
	COLON

	ADD
	SUB
	MUL
	QUO
	REM

	ASSIGN

	AND
	OR

	EQL
	NEQ
	LST
	GTT
	LTE
	GTE
	op_end

	key_start
	DEFINE
	FUNC
	IF
	VAR
	key_end

	tok_end
)

var tok_strings = map[Token]string{
	EOF:     "EOF",
	ILLEGAL: "Illegal",
	BOOL:    "Boolean",
	IDENT:   "Identifier",
	INTEGER: "Integer",
	LPAREN:  "(",
	RPAREN:  ")",
	COLON:   ":",
	ADD:     "+",
	SUB:     "-",
	MUL:     "*",
	QUO:     "/",
	REM:     "%",
	ASSIGN:  "=",
	AND:     "&&",
	OR:      "||",
	EQL:     "==",
	NEQ:     "!=",
	LST:     "<",
	GTT:     ">",
	LTE:     "<=",
	GTE:     ">=",
	DEFINE:  "define",
	FUNC:    "func",
	IF:      "if",
	VAR:     "var",
}

func (t Token) IsLiteral() bool {
	return t > lit_start && t < lit_end
}

func (t Token) IsOperator() bool {
	return t > op_start && t < op_end
}

func (t Token) IsKeyword() bool {
	return t > key_start && t < key_end
}

func Lookup(str string) Token {
	if str == "true" || str == "false" {
		return BOOL
	}
	for t, s := range tok_strings {
		if s == str {
			return t
		}
	}
	return IDENT
}

func (t Token) String() string {
	return tok_strings[t]
}

func (t Token) Valid() bool {
	return t > tok_start && t < tok_end
}
