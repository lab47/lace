// This file is generated by generate-std.joke script. Do not edit manually!

package io

import (
	. "github.com/candid82/joker/core"
	"io"
)


var __close__P ProcFn = __close_
var close_ Proc = Proc{Fn: __close__P, Name: "close_", Package: "std/io"}

func __close_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		f := ExtractObject(_env, _args, 0)
		_res := close(f)
		return _res

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __copy__P ProcFn = __copy_
var copy_ Proc = Proc{Fn: __copy__P, Name: "copy_", Package: "std/io"}

func __copy_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		dst := ExtractIOWriter(_env, _args, 0)
		src := ExtractIOReader(_env, _args, 1)
		 n, err := io.Copy(dst, src)
		PanicOnErr(err)
		_res := int(n)
		return MakeInt(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __pipe__P ProcFn = __pipe_
var pipe_ Proc = Proc{Fn: __pipe__P, Name: "pipe_", Package: "std/io"}

func __pipe_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 0:
		_res := pipe()
		return _res

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

func Init() {

	InternsOrThunks()
}

var ioNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.io"))

func init() {
	ioNamespace.Lazy = Init
}
