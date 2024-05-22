package core

import (
	"io"
)

type (
	IOWriter struct {
		io.Writer
		hash uint32
	}
)

var _ any = &IOWriter{}

func MakeIOWriter(w io.Writer) *IOWriter {
	res := &IOWriter{w, 0}
	res.hash = HashPtr(res)
	return res
}

func (iow *IOWriter) ToString(env *Env, escape bool) (string, error) {
	return "#object[IOWriter]", nil
}

func (iow *IOWriter) Close(env *Env) error {
	if c, ok := iow.Writer.(io.Closer); ok {
		return c.Close()
	} else {
		return env.NewError("Object is not closable: #object[IOWriter]")
	}
}
