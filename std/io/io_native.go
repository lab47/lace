package io

import (
	"io"

	. "github.com/lab47/lace/core"
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

func close(env *Env, f Object) (Nil, error) {
	if c, ok := f.(io.Closer); ok {
		if err := c.Close(); err != nil {
			return NIL, env.RT.NewError(err.Error())
		}
		return NIL, nil
	}
	s, err := f.ToString(env, false)
	if err != nil {
		return NIL, err
	}

	return NIL, env.RT.NewError("Object is not closable: " + s)
}
