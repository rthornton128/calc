package token

type File struct {
	newlines []Pos
	name     string
	src      string
}

func NewFile(name, src string) *File {
	return &File{
		newlines: make([]Pos, 0, 16),
		name:     name,
		src:      src,
	}
}

func (f *File) AddLine(p Pos) {
	base := Pos(1)
	if p.Valid() && p >= base && p < base+Pos(f.Size()) {
		f.newlines = append(f.newlines, p)
	}
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
