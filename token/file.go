package token

type File struct {
	base     int
	name     string
	src      string
	newlines []Pos
}

func NewFile(name, src string) *File {
	return &File{
		base:     1,
		name:     name,
		src:      src,
		newlines: make([]Pos, 0, 16),
	}
}

func (f *File) AddLine(p Pos) {
	base := Pos(1)
	if p.Valid() && p >= base && p < base+Pos(f.Size()) {
		f.newlines = append(f.newlines, p)
	}
}

func (f *File) Base() int {
	return f.base
}

func (f *File) Pos(offset int) Pos {
	if offset > len(f.src) {
		panic("illegal file offset")
	}
	return Pos(f.base + offset)
}

func (f *File) Position(p Pos) Position {
	start := Pos(0)
	col, row := int(p), 1

	for i, nl := range f.newlines {
		if p <= nl {
			col, row = int(p-start), i+1
			break
		}
		start = nl
	}

	return Position{Filename: f.name, Col: col, Row: row}
}

func (f *File) Size() int {
	return len(f.src)
}
