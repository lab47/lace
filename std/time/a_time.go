// This file is generated by generate-std.joke script. Do not edit manually!

package time

import (
	. "github.com/candid82/joker/core"
	"time"
)

var ansi_c_ String
var hour_ *BigInt
var kitchen_ String
var microsecond_ Int
var millisecond_ Int
var minute_ *BigInt
var nanosecond_ Int
var rfc1123_ String
var rfc1123_z_ String
var rfc3339_ String
var rfc3339_nano_ String
var rfc822_ String
var rfc822_z_ String
var rfc850_ String
var ruby_date_ String
var second_ Int
var stamp_ String
var stamp_micro_ String
var stamp_milli_ String
var stamp_nano_ String
var unix_date_ String
var __add__P ProcFn = __add_
var add_ Proc = Proc{Fn: __add__P, Name: "add_", Package: "std/time"}

func __add_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		t := ExtractTime(_args, 0)
		d := ExtractInt(_args, 1)
		_res := t.Add(time.Duration(d))
		return MakeTime(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __add_date__P ProcFn = __add_date_
var add_date_ Proc = Proc{Fn: __add_date__P, Name: "add_date_", Package: "std/time"}

func __add_date_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 4:
		t := ExtractTime(_args, 0)
		years := ExtractInt(_args, 1)
		months := ExtractInt(_args, 2)
		days := ExtractInt(_args, 3)
		_res := t.AddDate(years, months, days)
		return MakeTime(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __format__P ProcFn = __format_
var format_ Proc = Proc{Fn: __format__P, Name: "format_", Package: "std/time"}

func __format_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		t := ExtractTime(_args, 0)
		layout := ExtractString(_args, 1)
		_res := t.Format(layout)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __from_unix__P ProcFn = __from_unix_
var from_unix_ Proc = Proc{Fn: __from_unix__P, Name: "from_unix_", Package: "std/time"}

func __from_unix_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		sec := ExtractInt(_args, 0)
		nsec := ExtractInt(_args, 1)
		_res := time.Unix(int64(sec), int64(nsec))
		return MakeTime(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __hours__P ProcFn = __hours_
var hours_ Proc = Proc{Fn: __hours__P, Name: "hours_", Package: "std/time"}

func __hours_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		d := ExtractInt(_args, 0)
		_res := time.Duration(d).Hours()
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __minutes__P ProcFn = __minutes_
var minutes_ Proc = Proc{Fn: __minutes__P, Name: "minutes_", Package: "std/time"}

func __minutes_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		d := ExtractInt(_args, 0)
		_res := time.Duration(d).Minutes()
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __now__P ProcFn = __now_
var now_ Proc = Proc{Fn: __now__P, Name: "now_", Package: "std/time"}

func __now_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 0:
		_res := time.Now()
		return MakeTime(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __parse__P ProcFn = __parse_
var parse_ Proc = Proc{Fn: __parse__P, Name: "parse_", Package: "std/time"}

func __parse_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		layout := ExtractString(_args, 0)
		value := ExtractString(_args, 1)
		_res, err := time.Parse(layout, value)
		PanicOnErr(err)
		return MakeTime(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __parse_duration__P ProcFn = __parse_duration_
var parse_duration_ Proc = Proc{Fn: __parse_duration__P, Name: "parse_duration_", Package: "std/time"}

func __parse_duration_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		t, err := time.ParseDuration(s)
		PanicOnErr(err)
		_res := int(t)
		return MakeInt(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __round__P ProcFn = __round_
var round_ Proc = Proc{Fn: __round__P, Name: "round_", Package: "std/time"}

func __round_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		d := ExtractInt(_args, 0)
		m := ExtractInt(_args, 1)
		_res := int(time.Duration(d).Round(time.Duration(m)))
		return MakeInt(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __seconds__P ProcFn = __seconds_
var seconds_ Proc = Proc{Fn: __seconds__P, Name: "seconds_", Package: "std/time"}

func __seconds_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		d := ExtractInt(_args, 0)
		_res := time.Duration(d).Seconds()
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __since__P ProcFn = __since_
var since_ Proc = Proc{Fn: __since__P, Name: "since_", Package: "std/time"}

func __since_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		t := ExtractTime(_args, 0)
		_res := int(time.Since(t))
		return MakeInt(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __sleep__P ProcFn = __sleep_
var sleep_ Proc = Proc{Fn: __sleep__P, Name: "sleep_", Package: "std/time"}

func __sleep_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		d := ExtractInt(_args, 0)
		 RT.GIL.Unlock()
		time.Sleep(time.Duration(d))
		RT.GIL.Lock()
		_res := NIL
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __string__P ProcFn = __string_
var string_ Proc = Proc{Fn: __string__P, Name: "string_", Package: "std/time"}

func __string_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		d := ExtractInt(_args, 0)
		_res := time.Duration(d).String()
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __sub__P ProcFn = __sub_
var sub_ Proc = Proc{Fn: __sub__P, Name: "sub_", Package: "std/time"}

func __sub_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		t := ExtractTime(_args, 0)
		u := ExtractTime(_args, 1)
		_res := int(t.Sub(u))
		return MakeInt(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __truncate__P ProcFn = __truncate_
var truncate_ Proc = Proc{Fn: __truncate__P, Name: "truncate_", Package: "std/time"}

func __truncate_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		d := ExtractInt(_args, 0)
		m := ExtractInt(_args, 1)
		_res := int(time.Duration(d).Truncate(time.Duration(m)))
		return MakeInt(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __unix__P ProcFn = __unix_
var unix_ Proc = Proc{Fn: __unix__P, Name: "unix_", Package: "std/time"}

func __unix_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		t := ExtractTime(_args, 0)
		_res := int(t.Unix())
		return MakeInt(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __until__P ProcFn = __until_
var until_ Proc = Proc{Fn: __until__P, Name: "until_", Package: "std/time"}

func __until_(_env *Env, _args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		t := ExtractTime(_args, 0)
		_res := int(time.Until(t))
		return MakeInt(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

func Init() {
	ansi_c_ = MakeString(time.ANSIC)
	hour_ = MakeBigInt(int64(time.Hour))
	kitchen_ = MakeString(time.Kitchen)
	microsecond_ = MakeInt(int(time.Microsecond))
	millisecond_ = MakeInt(int(time.Millisecond))
	minute_ = MakeBigInt(int64(time.Minute))
	nanosecond_ = MakeInt(int(time.Nanosecond))
	rfc1123_ = MakeString(time.RFC1123)
	rfc1123_z_ = MakeString(time.RFC1123Z)
	rfc3339_ = MakeString(time.RFC3339)
	rfc3339_nano_ = MakeString(time.RFC3339Nano)
	rfc822_ = MakeString(time.RFC822)
	rfc822_z_ = MakeString(time.RFC822Z)
	rfc850_ = MakeString(time.RFC850)
	ruby_date_ = MakeString(time.RubyDate)
	second_ = MakeInt(int(time.Second))
	stamp_ = MakeString(time.Stamp)
	stamp_micro_ = MakeString(time.StampMicro)
	stamp_milli_ = MakeString(time.StampMilli)
	stamp_nano_ = MakeString(time.StampNano)
	unix_date_ = MakeString(time.UnixDate)
	InternsOrThunks()
}

var timeNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.time"))

func init() {
	timeNamespace.Lazy = Init
}
