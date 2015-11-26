// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package token

// File represents a single source file. It is used to track the number of
// newlines in the file, it's size, name and position within a fileset.
type File struct {
	base  int
	name  string
	lines []int
	size  int
}

// NewFile returns a new file object
func NewFile(name string, base, size int) *File {
	return &File{
		base:  base,
		name:  name,
		lines: make([]int, 0, 16),
		size:  size,
	}
}

// AddLine adds the position of the start of a line in the source file at
// the given offset. Every file consists of at least one line at offset
// zero.
func (f *File) AddLine(offset int) {
	if offset >= f.base-1 && offset < f.base+f.size {
		f.lines = append(f.lines, offset)
	}
}

// Base returns the base offset of the file within a fileset
func (f *File) Base() int {
	return f.base
}

// Pos generates a Pos based on the offset. The position is the file's
// base+offset
func (f *File) Pos(offset int) Pos {
	if offset < 0 || offset > f.size {
		panic("illegal file offset")
	}
	return Pos(f.base + offset)
}

// Position returns the column and row position of a Pos within the file
func (f *File) Position(p Pos) Position {
	col, row := int(p)-f.Base()+1, 1

	for i, nl := range f.lines {
		if p > f.Pos(nl) {
			col, row = int(p-f.Pos(nl))-f.Base()+1, i+1
		}
	}

	return Position{Filename: f.name, Col: col, Row: row}
}

// Size returns the length of the source code of the file.
func (f *File) Size() int {
	return f.size
}
