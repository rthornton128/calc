// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package token

type FileSet struct {
	base  int
	files []*File
}

func NewFileSet() *FileSet {
	return &FileSet{base: 1}
}

func (fs *FileSet) AddFile(name, src string) {
	size := len(src)
	fs.files = append(fs.files, NewFile(name, fs.base, size))
	fs.base += size
}

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
