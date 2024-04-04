// This file is generated by generate-std.joke script. Do not edit manually!

package base64

import (
	. "github.com/candid82/joker/core"
)


var __decode_string__P ProcFn = __decode_string_
var decode_string_ Proc = Proc{Fn: __decode_string__P, Name: "decode_string_", Package: "std/base64"}

func __decode_string_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_env, _args, 0)
		_res := decodeString(s)
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __encode_string__P ProcFn = __encode_string_
var encode_string_ Proc = Proc{Fn: __encode_string__P, Name: "encode_string_", Package: "std/base64"}

func __encode_string_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_env, _args, 0)
		_res := encodeString(s)
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

func Init() {

	InternsOrThunks()
}

var base64Namespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.base64"))

func init() {
	base64Namespace.Lazy = Init
}
