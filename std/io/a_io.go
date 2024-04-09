// This file is generated by generate-std.clj script. Do not edit manually!

package io

import (
	. "github.com/lab47/lace/core"
	"io"
)


var __close__P ProcFn = __close_
var close_ Proc = Proc{Fn: __close__P, Name: "close_", Package: "std/io"}

func __close_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		f, err := ExtractObject(_env, _args, 0); if err != nil { return nil, err }
		_res, err := close(_env, f)
		return _res, err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __copy__P ProcFn = __copy_
var copy_ Proc = Proc{Fn: __copy__P, Name: "copy_", Package: "std/io"}

func __copy_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 2:
		var err error
		dst, err := ExtractIOWriter(_env, _args, 0); if err != nil { return nil, err }
		src, err := ExtractIOReader(_env, _args, 1); if err != nil { return nil, err }
		 n, err := io.Copy(dst, src)
		_res := int(n)
		return MakeInt(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __pipe__P ProcFn = __pipe_
var pipe_ Proc = Proc{Fn: __pipe__P, Name: "pipe_", Package: "std/io"}

func __pipe_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 0:
		var err error
		_res, err := pipe()
		return _res, err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

func Init(env *Env, ns *Namespace) {

	InternsOrThunks(env, ns)
}

func init() {
	AddNativeNamespace("lace.io", func(env *Env) error {
		ns := env.EnsureNamespace(MakeSymbol("lace.io"))
		Init(env, ns)
		return nil
	})
}
