package core

import (
	"io"
)

type (
	IOReader struct {
		io.Reader
		hash uint32
	}
)

var _ any = &IOReader{}

func MakeIOReader(r io.Reader) *IOReader {
	res := &IOReader{r, 0}
	res.hash = HashPtr(res)
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

func (ior *IOReader) WithInfo(info *ObjectInfo) any {
	return ior
}

func (ior *IOReader) Close(env *Env) error {
	if c, ok := ior.Reader.(io.Closer); ok {
		return c.Close()
	} else {
		return env.NewError("Object is not closable: #object[IOReader]")
	}
}
