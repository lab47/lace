package core

import (
	"os"
)

type (
	File struct {
		*os.File
	}
)

var _ Object = &File{}

func (f *File) ToString(env *Env, escape bool) (string, error) {
	return "#object[File]", nil
}

func (f *File) Equals(env *Env, other interface{}) bool {
	return f == other
}

func (f *File) GetInfo() *ObjectInfo {
	return nil
}

func (f *File) GetType() *Type {
	return TYPE.File
}

func (f *File) Hash(env *Env) (uint32, error) {
	return HashPtr(f), nil
}

func (f *File) WithInfo(info *ObjectInfo) Object {
	return f
}

func MakeFile(f *os.File) *File {
	return &File{f}
}

func ExtractFile(env *Env, args []Object, index int) (*File, error) {
	return EnsureFile(env, args, index)
}
