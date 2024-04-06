// This file is generated by generate-std.joke script. Do not edit manually!

package csv

import (
	. "github.com/lab47/lace/core"
)

var __csv_seq__P ProcFn = __csv_seq_
var csv_seq_ Proc = Proc{Fn: __csv_seq__P, Name: "csv_seq_", Package: "std/csv"}

func __csv_seq_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		rdr, err := ExtractObject(_env, _args, 0)
		if err != nil {
			return nil, err
		}
		_res, err := csvSeqOpts(_env, rdr, EmptyArrayMap())
		return _res, err

	case _c == 2:
		var err error
		rdr, err := ExtractObject(_env, _args, 0)
		if err != nil {
			return nil, err
		}
		opts, err := ExtractMap(_env, _args, 1)
		if err != nil {
			return nil, err
		}
		_res, err := csvSeqOpts(_env, rdr, opts)
		return _res, err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __write__P ProcFn = __write_
var write_ Proc = Proc{Fn: __write__P, Name: "write_", Package: "std/csv"}

func __write_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 2:
		var err error
		f, err := ExtractIOWriter(_env, _args, 0)
		if err != nil {
			return nil, err
		}
		data, err := ExtractSeqable(_env, _args, 1)
		if err != nil {
			return nil, err
		}
		_res, err := write(_env, f, data, EmptyArrayMap())
		return _res, err

	case _c == 3:
		var err error
		f, err := ExtractIOWriter(_env, _args, 0)
		if err != nil {
			return nil, err
		}
		data, err := ExtractSeqable(_env, _args, 1)
		if err != nil {
			return nil, err
		}
		opts, err := ExtractMap(_env, _args, 2)
		if err != nil {
			return nil, err
		}
		_res, err := write(_env, f, data, opts)
		return _res, err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __write_string__P ProcFn = __write_string_
var write_string_ Proc = Proc{Fn: __write_string__P, Name: "write_string_", Package: "std/csv"}

func __write_string_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		data, err := ExtractSeqable(_env, _args, 0)
		if err != nil {
			return nil, err
		}
		_res, err := writeString(_env, data, EmptyArrayMap())
		return MakeString(_res), err

	case _c == 2:
		var err error
		data, err := ExtractSeqable(_env, _args, 0)
		if err != nil {
			return nil, err
		}
		opts, err := ExtractMap(_env, _args, 1)
		if err != nil {
			return nil, err
		}
		_res, err := writeString(_env, data, opts)
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

func Init() {

	InternsOrThunks()
}

var csvNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.csv"))

func init() {
	csvNamespace.Lazy = Init
}
