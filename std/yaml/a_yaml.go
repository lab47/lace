// This file is generated by generate-std.joke script. Do not edit manually!

package yaml

import (
	. "github.com/candid82/joker/core"
)


var __read_string__P ProcFn = __read_string_
var read_string_ Proc = Proc{Fn: __read_string__P, Name: "read_string_", Package: "std/yaml"}

func __read_string_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		_res, err := readString(s)
		return _res, err

	default:
		PanicArity(_env, _c)
	}
	return NIL, nil
}

var __write_string__P ProcFn = __write_string_
var write_string_ Proc = Proc{Fn: __write_string__P, Name: "write_string_", Package: "std/yaml"}

func __write_string_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		v, err := ExtractObject(_env, _args, 0); if err != nil { return nil, err }
		_res, err := writeString(v)
		return _res, err

	default:
		PanicArity(_env, _c)
	}
	return NIL, nil
}

func Init() {

	InternsOrThunks()
}

var yamlNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.yaml"))

func init() {
	yamlNamespace.Lazy = Init
}
