package io

import (
	"io"

	. "github.com/candid82/joker/core"
)

func pipe() (Object, error) {
	r, w := io.Pipe()
	res := EmptyVector()
	res, err := res.Conjoin(MakeIOReader(r))
	if err != nil {
		return nil, err
	}
	return res.Conjoin(MakeIOWriter(w))
}

func close(f Object) (Nil, error) {
	if c, ok := f.(io.Closer); ok {
		if err := c.Close(); err != nil {
			return NIL, StubNewError(err.Error())
		}
		return NIL, nil
	}
	return NIL, StubNewError("Object is not closable: " + f.ToString(false))
}
