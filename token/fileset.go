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
	files = append(files, NewFile(name, base, size))
	base += size
}

func (fs *FileSet) Position(p Pos) Position {
	for _, f := range fs.files {
		if p >= Pos(f.Base()) && p < Pos(f.Base()+f.Size()) {
			return f.Position(p)
		}
	}
}
