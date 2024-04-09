package core

import (
	"io"
	"unsafe"
)

type (
	IOReader struct {
		io.Reader
		hash uint32
	}
)

var _ Object = &IOReader{}

func MakeIOReader(r io.Reader) *IOReader {
	res := &IOReader{r, 0}
	res.hash = HashPtr(uintptr(unsafe.Pointer(res)))
	return res
}

func (ior *IOReader) ToString(env *Env, escape bool) (string, error) {
	return "#object[IOReader]", nil
}

func (ior *IOReader) Equals(env *Env, other interface{}) bool {
	return ior == other
}

func (ior *IOReader) GetInfo() *ObjectInfo {
	return nil
}

func (ior *IOReader) GetType() *Type {
	return TYPE.IOReader
}

func (ior *IOReader) Hash(env *Env) (uint32, error) {
	return ior.hash, nil
}

func (ior *IOReader) WithInfo(info *ObjectInfo) Object {
	return ior
}

func (ior *IOReader) Close(env *Env) error {
	if c, ok := ior.Reader.(io.Closer); ok {
		return c.Close()
	} else {
		return env.RT.NewError("Object is not closable: #object[IOReader]")
	}
}
