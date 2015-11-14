// Copyright (c) 2015, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

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
