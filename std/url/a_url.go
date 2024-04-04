// This file is generated by generate-std.joke script. Do not edit manually!

package url

import (
	. "github.com/candid82/joker/core"
	"net/url"
)


var __path_escape__P ProcFn = __path_escape_
var path_escape_ Proc = Proc{Fn: __path_escape__P, Name: "path_escape_", Package: "std/url"}

func __path_escape_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_env, _args, 0)
		_res := url.PathEscape(s)
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __path_unescape__P ProcFn = __path_unescape_
var path_unescape_ Proc = Proc{Fn: __path_unescape__P, Name: "path_unescape_", Package: "std/url"}

func __path_unescape_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_env, _args, 0)
		_res := pathUnescape(s)
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __query_escape__P ProcFn = __query_escape_
var query_escape_ Proc = Proc{Fn: __query_escape__P, Name: "query_escape_", Package: "std/url"}

func __query_escape_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_env, _args, 0)
		_res := url.QueryEscape(s)
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __query_unescape__P ProcFn = __query_unescape_
var query_unescape_ Proc = Proc{Fn: __query_unescape__P, Name: "query_unescape_", Package: "std/url"}

func __query_unescape_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_env, _args, 0)
		_res := queryUnescape(s)
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

func Init() {

	InternsOrThunks()
}

var urlNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.url"))

func init() {
	urlNamespace.Lazy = Init
}
