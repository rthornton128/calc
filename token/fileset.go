// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package token

// FileSet holds all the files for the source code
type FileSet struct {
	base  int
	files []*File
}

// NewFileSet creates a new FileSet object
func NewFileSet() *FileSet {
	return &FileSet{base: 1}
}

// Add appends a new file to the fileset
func (fs *FileSet) Add(name string, sz int) *File {
	f := NewFile(name, fs.base, sz)
	fs.files = append(fs.files, f)
	fs.base += sz
	return f
}

// Position returns the row and column position of the given Pos p
func (fs *FileSet) Position(p Pos) Position {
	var pos Position
	if !p.Valid() {
		panic("invalid position")
	}
	for _, f := range fs.files {
		if p >= Pos(f.Base()) && p < Pos(f.Base()+f.Size()) {
			pos = f.Position(p)
		}
	}
	return pos
}
