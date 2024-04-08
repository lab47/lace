// This file is generated by generate-std.clj script. Do not edit manually!

package string

import (
	. "github.com/lab47/lace/core"
	"regexp"
	"strings"
	"unicode"
)


var __isblank__P ProcFn = __isblank_
var isblank_ Proc = Proc{Fn: __isblank__P, Name: "isblank_", Package: "std/string"}

func __isblank_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractObject(_env, _args, 0); if err != nil { return nil, err }
		_res, err := isBlank(_env, s)
		return MakeBoolean(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __capitalize__P ProcFn = __capitalize_
var capitalize_ Proc = Proc{Fn: __capitalize__P, Name: "capitalize_", Package: "std/string"}

func __capitalize_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractStringable(_env, _args, 0); if err != nil { return nil, err }
		_res, err := capitalize(s)
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __isends_with__P ProcFn = __isends_with_
var isends_with_ Proc = Proc{Fn: __isends_with__P, Name: "isends_with_", Package: "std/string"}

func __isends_with_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 2:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		substr, err := ExtractStringable(_env, _args, 1); if err != nil { return nil, err }
		_res, err := strings.HasSuffix(s, substr), nil
		return MakeBoolean(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __escape__P ProcFn = __escape_
var escape_ Proc = Proc{Fn: __escape__P, Name: "escape_", Package: "std/string"}

func __escape_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 2:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		cmap, err := ExtractCallable(_env, _args, 1); if err != nil { return nil, err }
		_res, err := escape(_env, s, cmap)
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __isincludes__P ProcFn = __isincludes_
var isincludes_ Proc = Proc{Fn: __isincludes__P, Name: "isincludes_", Package: "std/string"}

func __isincludes_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 2:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		substr, err := ExtractStringable(_env, _args, 1); if err != nil { return nil, err }
		_res, err := strings.Contains(s, substr), nil
		return MakeBoolean(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __index_of__P ProcFn = __index_of_
var index_of_ Proc = Proc{Fn: __index_of__P, Name: "index_of_", Package: "std/string"}

func __index_of_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 2:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		value, err := ExtractObject(_env, _args, 1); if err != nil { return nil, err }
		_res, err := indexOf(s, value, 0)
		return _res, err

	case _c == 3:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		value, err := ExtractObject(_env, _args, 1); if err != nil { return nil, err }
		from, err := ExtractInt(_env, _args, 2); if err != nil { return nil, err }
		_res, err := indexOf(s, value, from)
		return _res, err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __join__P ProcFn = __join_
var join_ Proc = Proc{Fn: __join__P, Name: "join_", Package: "std/string"}

func __join_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		coll, err := ExtractSeqable(_env, _args, 0); if err != nil { return nil, err }
		_res, err := join("", coll)
		return MakeString(_res), err

	case _c == 2:
		var err error
		separator, err := ExtractStringable(_env, _args, 0); if err != nil { return nil, err }
		coll, err := ExtractSeqable(_env, _args, 1); if err != nil { return nil, err }
		_res, err := join(separator, coll)
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __last_index_of__P ProcFn = __last_index_of_
var last_index_of_ Proc = Proc{Fn: __last_index_of__P, Name: "last_index_of_", Package: "std/string"}

func __last_index_of_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 2:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		value, err := ExtractObject(_env, _args, 1); if err != nil { return nil, err }
		_res, err := lastIndexOf(s, value, 0)
		return _res, err

	case _c == 3:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		value, err := ExtractObject(_env, _args, 1); if err != nil { return nil, err }
		from, err := ExtractInt(_env, _args, 2); if err != nil { return nil, err }
		_res, err := lastIndexOf(s, value, from)
		return _res, err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __lower_case__P ProcFn = __lower_case_
var lower_case_ Proc = Proc{Fn: __lower_case__P, Name: "lower_case_", Package: "std/string"}

func __lower_case_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractStringable(_env, _args, 0); if err != nil { return nil, err }
		_res, err := strings.ToLower(s), nil
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __pad_left__P ProcFn = __pad_left_
var pad_left_ Proc = Proc{Fn: __pad_left__P, Name: "pad_left_", Package: "std/string"}

func __pad_left_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 3:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		pad, err := ExtractStringable(_env, _args, 1); if err != nil { return nil, err }
		n, err := ExtractInt(_env, _args, 2); if err != nil { return nil, err }
		_res, err := padLeft(s, pad, n)
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __pad_right__P ProcFn = __pad_right_
var pad_right_ Proc = Proc{Fn: __pad_right__P, Name: "pad_right_", Package: "std/string"}

func __pad_right_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 3:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		pad, err := ExtractStringable(_env, _args, 1); if err != nil { return nil, err }
		n, err := ExtractInt(_env, _args, 2); if err != nil { return nil, err }
		_res, err := padRight(s, pad, n)
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __re_quote__P ProcFn = __re_quote_
var re_quote_ Proc = Proc{Fn: __re_quote__P, Name: "re_quote_", Package: "std/string"}

func __re_quote_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		_res, err := regexp.Compile(regexp.QuoteMeta(s))
		return MakeRegex(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __replace__P ProcFn = __replace_
var replace_ Proc = Proc{Fn: __replace__P, Name: "replace_", Package: "std/string"}

func __replace_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 3:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		match, err := ExtractObject(_env, _args, 1); if err != nil { return nil, err }
		repl, err := ExtractStringable(_env, _args, 2); if err != nil { return nil, err }
		_res, err := replace(s, match, repl)
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __replace_first__P ProcFn = __replace_first_
var replace_first_ Proc = Proc{Fn: __replace_first__P, Name: "replace_first_", Package: "std/string"}

func __replace_first_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 3:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		match, err := ExtractObject(_env, _args, 1); if err != nil { return nil, err }
		repl, err := ExtractStringable(_env, _args, 2); if err != nil { return nil, err }
		_res, err := replaceFirst(s, match, repl)
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __reverse__P ProcFn = __reverse_
var reverse_ Proc = Proc{Fn: __reverse__P, Name: "reverse_", Package: "std/string"}

func __reverse_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		_res, err := reverse(s)
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __split__P ProcFn = __split_
var split_ Proc = Proc{Fn: __split__P, Name: "split_", Package: "std/string"}

func __split_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 2:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		sep, err := ExtractObject(_env, _args, 1); if err != nil { return nil, err }
		_res, err := splitOnStringOrRegex(s, sep, 0)
		return _res, err

	case _c == 3:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		sep, err := ExtractObject(_env, _args, 1); if err != nil { return nil, err }
		n, err := ExtractInt(_env, _args, 2); if err != nil { return nil, err }
		_res, err := splitOnStringOrRegex(s, sep, n)
		return _res, err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __split_lines__P ProcFn = __split_lines_
var split_lines_ Proc = Proc{Fn: __split_lines__P, Name: "split_lines_", Package: "std/string"}

func __split_lines_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		_res, err := split(s, newLine, 0)
		return _res, err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __isstarts_with__P ProcFn = __isstarts_with_
var isstarts_with_ Proc = Proc{Fn: __isstarts_with__P, Name: "isstarts_with_", Package: "std/string"}

func __isstarts_with_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 2:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		substr, err := ExtractStringable(_env, _args, 1); if err != nil { return nil, err }
		_res, err := strings.HasPrefix(s, substr), nil
		return MakeBoolean(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __trim__P ProcFn = __trim_
var trim_ Proc = Proc{Fn: __trim__P, Name: "trim_", Package: "std/string"}

func __trim_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		_res, err := strings.TrimSpace(s), nil
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __trim_left__P ProcFn = __trim_left_
var trim_left_ Proc = Proc{Fn: __trim_left__P, Name: "trim_left_", Package: "std/string"}

func __trim_left_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		_res, err := strings.TrimLeftFunc(s, unicode.IsSpace), nil
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __trim_newline__P ProcFn = __trim_newline_
var trim_newline_ Proc = Proc{Fn: __trim_newline__P, Name: "trim_newline_", Package: "std/string"}

func __trim_newline_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		_res, err := strings.TrimRight(s, "\n\r"), nil
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __trim_right__P ProcFn = __trim_right_
var trim_right_ Proc = Proc{Fn: __trim_right__P, Name: "trim_right_", Package: "std/string"}

func __trim_right_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		_res, err := strings.TrimRightFunc(s, unicode.IsSpace), nil
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __triml__P ProcFn = __triml_
var triml_ Proc = Proc{Fn: __triml__P, Name: "triml_", Package: "std/string"}

func __triml_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		_res, err := strings.TrimLeftFunc(s, unicode.IsSpace), nil
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __trimr__P ProcFn = __trimr_
var trimr_ Proc = Proc{Fn: __trimr__P, Name: "trimr_", Package: "std/string"}

func __trimr_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractString(_env, _args, 0); if err != nil { return nil, err }
		_res, err := strings.TrimRightFunc(s, unicode.IsSpace), nil
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

var __upper_case__P ProcFn = __upper_case_
var upper_case_ Proc = Proc{Fn: __upper_case__P, Name: "upper_case_", Package: "std/string"}

func __upper_case_(_env *Env, _args []Object) (Object, error) {
	_c := len(_args)
	switch {
	case _c == 1:
		var err error
		s, err := ExtractStringable(_env, _args, 0); if err != nil { return nil, err }
		_res, err := strings.ToUpper(s), nil
		return MakeString(_res), err

	default:
		return nil, ErrorArity(_env, _c)
	}
}

func Init(env *Env, ns *Namespace) {

	InternsOrThunks(env, ns)
}

func init() {
	AddNativeNamespace("lace.string", func(env *Env) error {
		ns := env.EnsureNamespace(MakeSymbol("lace.string"))
		Init(env, ns)
		return nil
	})
}
