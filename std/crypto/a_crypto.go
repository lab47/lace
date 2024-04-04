// This file is generated by generate-std.joke script. Do not edit manually!

package crypto

import (
	. "github.com/candid82/joker/core"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)


var __hmac__P ProcFn = __hmac_
var hmac_ Proc = Proc{Fn: __hmac__P, Name: "hmac_", Package: "std/crypto"}

func __hmac_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 3:
		algorithm := ExtractKeyword(_args, 0)
		message := ExtractString(_args, 1)
		key := ExtractString(_args, 2)
		_res := hmacSum(algorithm, message, key)
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __md5__P ProcFn = __md5_
var md5_ Proc = Proc{Fn: __md5__P, Name: "md5_", Package: "std/crypto"}

func __md5_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		data := ExtractString(_args, 0)
		 t := md5.Sum([]byte(data))
		_res := string(t[:])
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __sha1__P ProcFn = __sha1_
var sha1_ Proc = Proc{Fn: __sha1__P, Name: "sha1_", Package: "std/crypto"}

func __sha1_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		data := ExtractString(_args, 0)
		 t := sha1.Sum([]byte(data))
		_res := string(t[:])
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __sha224__P ProcFn = __sha224_
var sha224_ Proc = Proc{Fn: __sha224__P, Name: "sha224_", Package: "std/crypto"}

func __sha224_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		data := ExtractString(_args, 0)
		 t := sha256.Sum224([]byte(data))
		_res := string(t[:])
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __sha256__P ProcFn = __sha256_
var sha256_ Proc = Proc{Fn: __sha256__P, Name: "sha256_", Package: "std/crypto"}

func __sha256_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		data := ExtractString(_args, 0)
		 t := sha256.Sum256([]byte(data))
		_res := string(t[:])
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __sha384__P ProcFn = __sha384_
var sha384_ Proc = Proc{Fn: __sha384__P, Name: "sha384_", Package: "std/crypto"}

func __sha384_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		data := ExtractString(_args, 0)
		 t := sha512.Sum384([]byte(data))
		_res := string(t[:])
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __sha512__P ProcFn = __sha512_
var sha512_ Proc = Proc{Fn: __sha512__P, Name: "sha512_", Package: "std/crypto"}

func __sha512_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		data := ExtractString(_args, 0)
		 t := sha512.Sum512([]byte(data))
		_res := string(t[:])
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __sha512_224__P ProcFn = __sha512_224_
var sha512_224_ Proc = Proc{Fn: __sha512_224__P, Name: "sha512_224_", Package: "std/crypto"}

func __sha512_224_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		data := ExtractString(_args, 0)
		 t := sha512.Sum512_224([]byte(data))
		_res := string(t[:])
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

var __sha512_256__P ProcFn = __sha512_256_
var sha512_256_ Proc = Proc{Fn: __sha512_256__P, Name: "sha512_256_", Package: "std/crypto"}

func __sha512_256_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		data := ExtractString(_args, 0)
		 t := sha512.Sum512_256([]byte(data))
		_res := string(t[:])
		return MakeString(_res)

	default:
		PanicArity(_env, _c)
	}
	return NIL
}

func Init() {

	InternsOrThunks()
}

var cryptoNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.crypto"))

func init() {
	cryptoNamespace.Lazy = Init
}
