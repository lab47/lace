package core

import (
	"io"
)

type (
	Reader struct {
		runeReader     io.RuneReader
		rw             *RuneWindow
		line           int
		prevLineLength int
		column         int
		isEof          bool
		rewind         int
		filename       *string
		args           map[int]Symbol
		posStack       []pos
	}
)

func NewReader(runeReader io.RuneReader, filename string) *Reader {
	return &Reader{
		line:       1,
		runeReader: runeReader,
		rw:         &RuneWindow{},
		rewind:     -1,
		filename:   STRINGS.Intern(filename),
		posStack:   make([]pos, 0, 8),
	}
}

func (reader *Reader) Get() (rune, error) {
	if reader.isEof {
		return EOF, nil
	}
	if reader.rewind > -1 {
		r, err := top(reader.rw, reader.rewind)
		if err != nil {
			return 0, err
		}
		reader.rewind--
		if r == '\n' {
			reader.line++
			reader.prevLineLength = reader.column
			reader.column = 0
		} else {
			reader.column++
		}
		return r, nil
	}
	r, _, err := reader.runeReader.ReadRune()
	switch {
	case err == io.EOF:
		reader.isEof = true
		return EOF, nil
	case err != nil:
		return 0, err
	case r == '\n':
		reader.line++
		reader.prevLineLength = reader.column
		reader.column = 0
		add(reader.rw, r)
		return r, nil
	default:
		reader.column++
		add(reader.rw, r)
		return r, nil
	}
}

func (reader *Reader) Unget() {
	if reader.isEof {
		return
	}
	reader.rewind++
	if reader.column == 0 {
		reader.line--
		reader.column = reader.prevLineLength
	} else {
		reader.column--
	}
}

func (reader *Reader) Peek() rune {
	if reader.isEof {
		return EOF
	}
	r, err := reader.Get()
	if err != nil {
		return 0
	}
	reader.Unget()
	return r
}
