package core

import (
	"io"
	"unsafe"
)

type (
	IOWriter struct {
		io.Writer
		hash uint32
	}
)

var _ Object = &IOWriter{}

func MakeIOWriter(w io.Writer) *IOWriter {
	res := &IOWriter{w, 0}
	res.hash = HashPtr(uintptr(unsafe.Pointer(res)))
	return res
}

func (iow *IOWriter) ToString(env *Env, escape bool) (string, error) {
	return "#object[IOWriter]", nil
}

func (iow *IOWriter) Equals(env *Env, other interface{}) bool {
	return iow == other
}

func (iow *IOWriter) GetInfo() *ObjectInfo {
	return nil
}

func (iow *IOWriter) GetType() *Type {
	return TYPE.IOWriter
}

func (iow *IOWriter) Hash(env *Env) (uint32, error) {
	return iow.hash, nil
}

func (iow *IOWriter) WithInfo(info *ObjectInfo) Object {
	return iow
}

func (iow *IOWriter) Close(env *Env) error {
	if c, ok := iow.Writer.(io.Closer); ok {
		return c.Close()
	} else {
		return env.NewError("Object is not closable: #object[IOWriter]")
	}
}
