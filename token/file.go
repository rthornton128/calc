package token

import (
	"fmt"
)

type File struct {
	base     Pos
	newlines []Pos
	name     string
	src      string
	errs     []error
}

func NewFile(name, src string, base Pos) (*File, error) {
	if !base.Valid() {
		return nil, fmt.Errorf("Invalid base position!")
	}

	return &File{
		base:     base,
		newlines: make([]Pos, 0, 16),
		name:     name,
		src:      src,
		errs:     make([]error, 0, 10),
	}, nil
}

func (f *File) Base() Pos {
	return f.base
}

func (f *File) AddLine(p Pos) {
	f.newlines = append(f.newlines, p)
}

func (f *File) AddError(p Pos, msg string) {
	col, row := f.Position(p)
	f.errs = append(f.errs, fmt.Errorf("%s: %d,%d: %s", f.name, col, row, msg))
}

func (f *File) Position(p Pos) (col, row uint) {
	start := Pos(0)
	col, row = uint(p), 1

	for i, nl := range f.newlines {
		if p <= nl {
			col, row = uint(p-start), uint(i+1)
			break
		}
		start = nl
	}

	return
}
