package io

import (
	"io"

	. "github.com/candid82/joker/core"
)

func pipe() Object {
	r, w := io.Pipe()
	res := EmptyVector()
	res = res.Conjoin(MakeIOReader(r))
	res = res.Conjoin(MakeIOWriter(w))
	return res
}

func close(f Object) Nil {
	if c, ok := f.(io.Closer); ok {
		if err := c.Close(); err != nil {
			panic(StubNewError(err.Error()))
		}
		return NIL
	}
	panic(StubNewError("Object is not closable: " + f.ToString(false)))
}
