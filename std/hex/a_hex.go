// This file is generated by generate-std.joke script. Do not edit manually!

package hex

import (
	"encoding/hex"

	. "github.com/lab47/lace/core"
)

var __decode_string__P ProcFn = __decode_string_
var decode_string_ Proc = Proc{Fn: __decode_string__P, Name: "decode_string_", Package: "std/hex"}

func __decode_string_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0)
		if err != nil {
			return nil, err
		}
		t, err := hex.DecodeString(s)
		_res := string(t)
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __encode_string__P ProcFn = __encode_string_
var encode_string_ Proc = Proc{Fn: __encode_string__P, Name: "encode_string_", Package: "std/hex"}

func __encode_string_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0)
		if err != nil {
			return nil, err
		}
		_res, err := hex.EncodeToString([]byte(s)), nil
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

func Init() {

	InternsOrThunks()
}

var hexNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("lace.hex"))

func init() {
	hexNamespace.Lazy = Init
}
