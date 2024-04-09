package core

import (
	"bytes"
	"unsafe"
)

type (
	Buffer struct {
		*bytes.Buffer
		hash uint32
	}
)

var _ Object = &Buffer{}

func MakeBuffer(b *bytes.Buffer) *Buffer {
	res := &Buffer{b, 0}
	res.hash = HashPtr(uintptr(unsafe.Pointer(res)))
	return res
}

func (b *Buffer) ToString(env *Env, escape bool) (string, error) {
	return b.String(), nil
}

func (b *Buffer) Equals(env *Env, other interface{}) bool {
	return b == other
}

func (b *Buffer) GetInfo() *ObjectInfo {
	return nil
}

func (b *Buffer) GetType() *Type {
	return TYPE.Buffer
}

func (b *Buffer) Hash(env *Env) (uint32, error) {
	return b.hash, nil
}

func (b *Buffer) WithInfo(info *ObjectInfo) Object {
	return b
}
