package core

import (
	"os"
)

type (
	File struct {
		*os.File
	}
)

var _ any = &File{}

func (f *File) ToString(env *Env, escape bool) (string, error) {
	return "#object[File]", nil
}

func MakeFile(f *os.File) *File {
	return &File{f}
}

func ExtractFile(env *Env, args []any, index int) (*File, error) {
	return EnsureFile(env, args, index)
}
