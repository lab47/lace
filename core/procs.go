package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var coreNamespaces []string

type (
	Phase        int
	Dialect      int
	StringReader interface {
		ReadString(delim byte) (s string, e error)
	}
)

const (
	READ Phase = iota
	PARSE
	EVAL
	PRINT_IF_NOT_NIL
)

const VERSION = "v0.14.2"

const (
	CLJ Dialect = iota
	CLJS
	LACE
	EDN
	UNKNOWN
)

func ExtractCallable(env *Env, args []Object, index int) (Callable, error) {
	return EnsureCallable(env, args, index)
}

func ExtractObject(env *Env, args []Object, index int) (Object, error) {
	return args[index], nil
}

func ExtractString(env *Env, args []Object, index int) (string, error) {
	s, err := EnsureString(env, args, index)
	if err != nil {
		return "", err
	}

	return s.S, nil
}

func ExtractKeyword(env *Env, args []Object, index int) (string, error) {
	k, err := EnsureKeyword(env, args, index)
	if err != nil {
		return "", err
	}
	return k.ToString(env, false)
}

func ExtractStringable(env *Env, args []Object, index int) (string, error) {
	s, err := EnsureStringable(args, index)
	if err != nil {
		return "", err
	}

	return s.S, nil
}

func ExtractStrings(env *Env, args []Object, index int) ([]string, error) {
	strs := make([]string, 0)
	for i := index; i < len(args); i++ {
		s, err := EnsureString(env, args, i)
		if err != nil {
			return nil, err
		}
		strs = append(strs, s.S)
	}
	return strs, nil
}

func ExtractInt(env *Env, args []Object, index int) (int, error) {
	i, err := EnsureInt(env, args, index)
	if err != nil {
		return 0, err
	}
	return i.I, nil
}

func ExtractBoolean(env *Env, args []Object, index int) (bool, error) {
	b, err := EnsureBoolean(env, args, index)
	if err != nil {
		return false, err
	}

	return b.B, nil
}

func ExtractChar(env *Env, args []Object, index int) (rune, error) {
	c, err := EnsureChar(env, args, index)
	if err != nil {
		return 0, err
	}

	return c.Ch, nil
}

func ExtractTime(env *Env, args []Object, index int) (time.Time, error) {
	t, err := EnsureTime(env, args, index)
	if err != nil {
		return time.Time{}, err
	}

	return t.T, nil
}

func ExtractDouble(env *Env, args []Object, index int) (float64, error) {
	d, err := EnsureDouble(env, args, index)
	if err != nil {
		return 0, err
	}

	return d.D, nil
}

func ExtractNumber(env *Env, args []Object, index int) (Number, error) {
	return EnsureNumber(env, args, index)
}

func ExtractRegex(env *Env, args []Object, index int) (*regexp.Regexp, error) {
	r, err := EnsureRegex(env, args, index)
	if err != nil {
		return nil, err
	}
	return r.R, nil
}

func ExtractSeqable(env *Env, args []Object, index int) (Seqable, error) {
	return EnsureSeqable(env, args, index)
}

func ExtractMap(env *Env, args []Object, index int) (Map, error) {
	return EnsureMap(env, args, index)
}

func ExtractIOReader(env *Env, args []Object, index int) (io.Reader, error) {
	return Ensureio_Reader(env, args, index)
}

func ExtractIOWriter(env *Env, args []Object, index int) (io.Writer, error) {
	return Ensureio_Writer(env, args, index)
}

var procMeta = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	switch obj := args[0].(type) {
	case Meta:
		meta := obj.GetMeta()
		if meta != nil {
			return meta, nil
		}
	case *Type:
		meta := obj.GetMeta()
		if meta != nil {
			return meta, nil
		}
	}
	return NIL, nil
}

var procWithMeta = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}
	m, err := EnsureMeta(env, args, 0)
	if err != nil {
		return nil, err
	}
	if args[1].Equals(env, NIL) {
		return args[0], nil
	}
	mm, err := EnsureMap(env, args, 1)
	if err != nil {
		return nil, err
	}

	return m.WithMeta(env, mm)
}

var procIsZero = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	n, err := EnsureNumber(env, args, 0)
	if err != nil {
		return nil, err
	}
	ops := GetOps(n)
	return Boolean{B: ops.IsZero(n)}, nil
}

var procIsPos = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	n, err := EnsureNumber(env, args, 0)
	if err != nil {
		return nil, err
	}
	ops := GetOps(n)
	return Boolean{B: ops.Gt(n, Int{I: 0})}, nil
}

var procIsNeg = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	n, err := EnsureNumber(env, args, 0)
	if err != nil {
		return nil, err
	}
	ops := GetOps(n)
	return Boolean{B: ops.Lt(n, Int{I: 0})}, nil
}

var procAdd = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	x, err := AssertNumber(env, args[0], "")
	if err != nil {
		return nil, err
	}
	y, err := AssertNumber(env, args[1], "")
	if err != nil {
		return nil, err
	}
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Add(x, y)
}

var procAddEx = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	x, err := AssertNumber(env, args[0], "")
	if err != nil {
		return nil, err
	}
	y, err := AssertNumber(env, args[1], "")
	if err != nil {
		return nil, err
	}
	ops := GetOps(x).Combine(GetOps(y)).Combine(BIGINT_OPS)
	return ops.Add(x, y)
}

var procMultiply = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	x, err := AssertNumber(env, args[0], "")
	if err != nil {
		return nil, err
	}
	y, err := AssertNumber(env, args[1], "")
	if err != nil {
		return nil, err
	}
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Multiply(x, y)
}

var procMultiplyEx = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	x, err := AssertNumber(env, args[0], "")
	if err != nil {
		return nil, err
	}
	y, err := AssertNumber(env, args[1], "")
	if err != nil {
		return nil, err
	}
	ops := GetOps(x).Combine(GetOps(y)).Combine(BIGINT_OPS)
	return ops.Multiply(x, y)
}

var procSubtract = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 2); err != nil {
		return nil, err
	}

	var a, b Object
	if len(args) == 1 {
		a = Int{I: 0}
		b = args[0]
	} else {
		a = args[0]
		b = args[1]
	}
	ops := GetOps(a).Combine(GetOps(b))
	av, err := AssertNumber(env, a, "")
	if err != nil {
		return nil, err
	}
	bv, err := AssertNumber(env, b, "")
	if err != nil {
		return nil, err
	}
	return ops.Subtract(av, bv)
}

var procSubtractEx = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 2); err != nil {
		return nil, err
	}

	var a, b Object
	if len(args) == 1 {
		a = Int{I: 0}
		b = args[0]
	} else {
		a = args[0]
		b = args[1]
	}
	ops := GetOps(a).Combine(GetOps(b)).Combine(BIGINT_OPS)
	av, err := AssertNumber(env, a, "")
	if err != nil {
		return nil, err
	}
	bv, err := AssertNumber(env, b, "")
	if err != nil {
		return nil, err
	}
	return ops.Subtract(av, bv)
}

var procDivide = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	x, err := EnsureNumber(env, args, 0)
	if err != nil {
		return nil, err
	}
	y, err := EnsureNumber(env, args, 1)
	if err != nil {
		return nil, err
	}
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Divide(x, y)
}

var procQuot = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	x, err := EnsureNumber(env, args, 0)
	if err != nil {
		return nil, err
	}
	y, err := EnsureNumber(env, args, 1)
	if err != nil {
		return nil, err
	}
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Quotient(x, y)
}

var procRem = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	x, err := EnsureNumber(env, args, 0)
	if err != nil {
		return nil, err
	}
	y, err := EnsureNumber(env, args, 1)
	if err != nil {
		return nil, err
	}
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Rem(x, y)
}

var procBitNot = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	x, err := AssertInt(env, args[0], "Bit operation not supported for "+args[0].GetType().Name())
	if err != nil {
		return nil, err
	}
	return Int{I: ^x.I}, nil
}

func AssertInts(env *Env, args []Object) (Int, Int, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return Int{}, Int{}, err
	}

	x, err := AssertInt(env, args[0], "Bit operation not supported for "+args[0].GetType().Name())
	if err != nil {
		return Int{}, Int{}, err
	}
	y, err := AssertInt(env, args[1], "Bit operation not supported for "+args[1].GetType().Name())
	if err != nil {
		return Int{}, Int{}, err
	}
	return x, y, nil
}

var procBitAnd = func(env *Env, args []Object) (Object, error) {
	x, y, err := AssertInts(env, args)
	if err != nil {
		return nil, err
	}
	return Int{I: x.I & y.I}, nil
}

var procBitOr = func(env *Env, args []Object) (Object, error) {
	x, y, err := AssertInts(env, args)
	if err != nil {
		return nil, err
	}
	return Int{I: x.I | y.I}, nil
}

var procBitXor = func(env *Env, args []Object) (Object, error) {
	x, y, err := AssertInts(env, args)
	if err != nil {
		return nil, err
	}
	return Int{I: x.I ^ y.I}, nil
}

var procBitAndNot = func(env *Env, args []Object) (Object, error) {
	x, y, err := AssertInts(env, args)
	if err != nil {
		return nil, err
	}
	return Int{I: x.I &^ y.I}, nil
}

var procBitClear = func(env *Env, args []Object) (Object, error) {
	x, y, err := AssertInts(env, args)
	if err != nil {
		return nil, err
	}
	return Int{I: x.I &^ (1 << uint(y.I))}, nil
}

var procBitSet = func(env *Env, args []Object) (Object, error) {
	x, y, err := AssertInts(env, args)
	if err != nil {
		return nil, err
	}
	return Int{I: x.I | (1 << uint(y.I))}, nil
}

var procBitFlip = func(env *Env, args []Object) (Object, error) {
	x, y, err := AssertInts(env, args)
	if err != nil {
		return nil, err
	}
	return Int{I: x.I ^ (1 << uint(y.I))}, nil
}

var procBitTest = func(env *Env, args []Object) (Object, error) {
	x, y, err := AssertInts(env, args)
	if err != nil {
		return nil, err
	}
	return Boolean{B: x.I&(1<<uint(y.I)) != 0}, nil
}

var procBitShiftLeft = func(env *Env, args []Object) (Object, error) {
	x, y, err := AssertInts(env, args)
	if err != nil {
		return nil, err
	}
	return Int{I: x.I << uint(y.I)}, nil
}

var procBitShiftRight = func(env *Env, args []Object) (Object, error) {
	x, y, err := AssertInts(env, args)
	if err != nil {
		return nil, err
	}
	return Int{I: x.I >> uint(y.I)}, nil
}

var procUnsignedBitShiftRight = func(env *Env, args []Object) (Object, error) {
	x, y, err := AssertInts(env, args)
	if err != nil {
		return nil, err
	}
	return Int{I: int(uint(x.I) >> uint(y.I))}, nil
}

var procExInfo = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 3); err != nil {
		return nil, err
	}

	res := &ExInfo{
		rt: env.RT.clone(),
	}
	s, err := EnsureString(env, args, 0)
	if err != nil {
		return nil, err
	}
	m, err := EnsureMap(env, args, 1)
	if err != nil {
		return nil, err
	}

	res.AddEqu(criticalKeywords.message, s)
	res.AddEqu(criticalKeywords.data, m)
	if len(args) == 3 {
		e, err := EnsureError(env, args, 2)
		if err != nil {
			return nil, err
		}
		res.AddEqu(criticalKeywords.cause, e)
	}
	return res, nil
}

var procExData = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	var ex *ExInfo
	if err := Cast(env, args[0], &ex); err != nil {
		return nil, err
	}

	if ok, res := ex.GetEqu(criticalKeywords.data); ok {
		return res, nil
	}
	return NIL, nil
}

var procExCause = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	var ex *ExInfo
	if err := Cast(env, args[0], &ex); err != nil {
		return nil, err
	}

	if ok, res := ex.GetEqu(criticalKeywords.cause); ok {
		return res, nil
	}
	return NIL, nil
}

var procExMessage = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	var ex Error
	if err := Cast(env, args[0], &ex); err != nil {
		return nil, err
	}

	return ex.Message(), nil
}

var procRegex = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	s, err := EnsureString(env, args, 0)
	if err != nil {
		return nil, err
	}
	r, err := regexp.Compile(s.S)
	if err != nil {
		return nil, env.RT.NewError("Invalid regex: " + err.Error())
	}
	return &Regex{R: r}, nil
}

func reGroups(s string, indexes []int) (Object, error) {
	if indexes == nil {
		return NIL, nil
	} else if len(indexes) == 2 {
		if indexes[0] == -1 {
			return NIL, nil
		} else {
			return String{S: s[indexes[0]:indexes[1]]}, nil
		}
	} else {
		v := EmptyVector()
		var err error
		for i := 0; i < len(indexes); i += 2 {
			if indexes[i] == -1 {
				v, err = v.Conjoin(NIL)
				if err != nil {
					return nil, err
				}
			} else {
				v, err = v.Conjoin(String{S: s[indexes[i]:indexes[i+1]]})
				if err != nil {
					return nil, err
				}
			}
		}
		return v, nil
	}
}

var procReSeq = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	re, err := EnsureRegex(env, args, 0)
	if err != nil {
		return nil, err
	}
	s, err := EnsureString(env, args, 1)
	if err != nil {
		return nil, err
	}
	matches := re.R.FindAllStringSubmatchIndex(s.S, -1)
	if matches == nil {
		return NIL, nil
	}
	res := make([]Object, len(matches))
	for i, match := range matches {
		grp, err := reGroups(s.S, match)
		if err != nil {
			return nil, err
		}
		res[i] = grp
	}
	return &ArraySeq{arr: res}, nil
}

var procReFind = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	re, err := EnsureRegex(env, args, 0)
	if err != nil {
		return nil, err
	}
	s, err := EnsureString(env, args, 1)
	if err != nil {
		return nil, err
	}
	match := re.R.FindStringSubmatchIndex(s.S)
	return reGroups(s.S, match)
}

var procRand = func(env *Env, args []Object) (Object, error) {
	r := rand.Float64()
	return Double{D: r}, nil
}

var procIsSpecialSymbol = func(env *Env, args []Object) (Object, error) {
	return Boolean{B: IsSpecialSymbol(args[0])}, nil
}

var procSubs = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 3); err != nil {
		return nil, err
	}

	s, err := EnsureString(env, args, 0)
	if err != nil {
		return nil, err
	}
	start, err := EnsureInt(env, args, 1)
	if err != nil {
		return nil, err
	}
	slen := utf8.RuneCountInString(s.S)
	end := slen
	if len(args) > 2 {
		x, err := EnsureInt(env, args, 2)
		if err != nil {
			return nil, err
		}
		end = x.I
	}
	if start.I < 0 || start.I > slen {
		return nil, env.RT.NewError(fmt.Sprintf("String index out of range: %d", start.I))
	}
	if end < 0 || end > slen {
		return nil, env.RT.NewError(fmt.Sprintf("String index out of range: %d", end))
	}
	return String{S: string([]rune(s.S)[start.I:end])}, nil
}

var procIntern = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 3); err != nil {
		return nil, err
	}

	ns, err := EnsureNamespace(env, args, 0)
	if err != nil {
		return nil, err
	}
	sym, err := EnsureSymbol(env, args, 1)
	if err != nil {
		return nil, err
	}
	vr, err := ns.Intern(env, sym)
	if err != nil {
		return nil, err
	}
	if len(args) == 3 {
		vr.Value = args[2]
	}
	return vr, nil
}

var procSetMeta = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	vr, err := EnsureVar(env, args, 0)
	if err != nil {
		return nil, err
	}
	meta, err := EnsureMap(env, args, 1)
	if err != nil {
		return nil, err
	}
	vr.meta = meta
	return NIL, nil
}

var procAtom = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 10000); err != nil {
		return nil, err
	}

	res := &Atom{
		value: args[0],
	}
	if len(args) > 1 {
		m, err := NewHashMap(env, args[1:]...)
		if err != nil {
			return nil, err
		}
		if ok, v := m.GetEqu(criticalKeywords.meta); ok {
			mm, err := AssertMap(env, v, "")
			if err != nil {
				return nil, err
			}
			res.meta = mm
		}
	}
	return res, nil
}

var procDeref = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	ed, err := EnsureDeref(env, args, 0)
	if err != nil {
		return nil, err
	}
	return ed.Deref(env)
}

var procSwap = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 10000); err != nil {
		return nil, err
	}

	a, err := EnsureAtom(env, args, 0)
	if err != nil {
		return nil, err
	}
	f, err := EnsureCallable(env, args, 1)
	if err != nil {
		return nil, err
	}
	fargs := append([]Object{a.value}, args[2:]...)
	v, err := f.Call(env, fargs)
	if err != nil {
		return nil, err
	}

	a.value = v
	return a.value, nil
}

var procSwapVals = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 10000); err != nil {
		return nil, err
	}

	a, err := EnsureAtom(env, args, 0)
	if err != nil {
		return nil, err
	}
	f, err := EnsureCallable(env, args, 1)
	if err != nil {
		return nil, err
	}
	fargs := append([]Object{a.value}, args[2:]...)
	oldValue := a.value
	v, err := f.Call(env, fargs)
	if err != nil {
		return nil, err
	}
	a.value = v
	return NewVectorFrom(oldValue, a.value), nil
}

var procReset = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	a, err := EnsureAtom(env, args, 0)
	if err != nil {
		return nil, err
	}
	a.value = args[1]
	return a.value, nil
}

var procResetVals = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	a, err := EnsureAtom(env, args, 0)
	if err != nil {
		return nil, err
	}
	oldValue := a.value
	a.value = args[1]
	return NewVectorFrom(oldValue, a.value), nil
}

var procAlterMeta = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 10000); err != nil {
		return nil, err
	}

	r, err := EnsureRef(env, args, 0)
	if err != nil {
		return nil, err
	}
	f, err := EnsureFn(env, args, 1)
	if err != nil {
		return nil, err
	}
	return r.AlterMeta(env, f, args[2:])
}

var procResetMeta = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	r, err := EnsureRef(env, args, 0)
	if err != nil {
		return nil, err
	}
	m, err := EnsureMap(env, args, 1)
	if err != nil {
		return nil, err
	}
	return r.ResetMeta(m), nil
}

var procEmpty = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	switch c := args[0].(type) {
	case Collection:
		return c.Empty(), nil
	default:
		return NIL, nil
	}
}

var procIsBound = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	vr, err := EnsureVar(env, args, 0)
	if err != nil {
		return nil, err
	}
	return Boolean{B: vr.Value != nil}, nil
}

func toNative(env *Env, obj Object) (interface{}, error) {
	switch obj := obj.(type) {
	case Native:
		return obj.Native(), nil
	default:
		return obj.ToString(env, false)
	}
}

var procFormat = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 10000); err != nil {
		return nil, err
	}

	s, err := EnsureString(env, args, 0)
	if err != nil {
		return nil, err
	}
	objs := args[1:]
	fargs := make([]interface{}, len(objs))
	for i, v := range objs {
		fargs[i], err = toNative(env, v)
		if err != nil {
			return nil, err
		}
	}
	res := fmt.Sprintf(s.S, fargs...)
	return String{S: res}, nil
}

var procList = func(env *Env, args []Object) (Object, error) {
	return NewListFrom(args...), nil
}

var procCons = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}
	s, err := EnsureSeqable(env, args, 1)
	if err != nil {
		return nil, err
	}
	return s.Seq().Cons(args[0]), nil
}

var procFirst = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}
	s, err := EnsureSeqable(env, args, 0)
	if err != nil {
		return nil, err
	}
	return s.Seq().First(env)
}

var procNext = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}
	s, err := EnsureSeqable(env, args, 0)
	if err != nil {
		return nil, err
	}
	res, err := s.Seq().Rest(env)
	if err != nil {
		return nil, err
	}
	empty, err := res.IsEmpty(env)
	if err != nil {
		return nil, err
	}

	if empty {
		return NIL, nil
	}
	return res, nil
}

var procRest = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}
	s, err := EnsureSeqable(env, args, 0)
	if err != nil {
		return nil, err
	}
	return s.Seq().Rest(env)
}

var procConj = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	switch c := args[0].(type) {
	case Conjable:
		return c.Conj(env, args[1])
	case Seq:
		return c.Cons(args[1]), nil
	default:
		return nil, env.RT.NewError("conj's first argument must be a collection, got " + c.GetType().Name())
	}
}

var procSeq = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}
	s, err := EnsureSeqable(env, args, 0)
	if err != nil {
		return nil, err
	}
	sq := s.Seq()
	empty, err := sq.IsEmpty(env)
	if err != nil {
		return nil, err
	}

	if empty {
		return NIL, nil
	}
	return sq, nil
}

var procIsInstance = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}
	t, err := EnsureType(env, args, 0)
	if err != nil {
		return nil, err
	}
	return Boolean{B: IsInstance(env, t, args[1])}, nil
}

var procAssoc = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 3, 3); err != nil {
		return nil, err
	}

	ea, err := EnsureAssociative(env, args, 0)
	if err != nil {
		return nil, err
	}

	return ea.Assoc(env, args[1], args[2])
}

var procEquals = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	return Boolean{B: args[0].Equals(env, args[1])}, nil
}

var procCount = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	switch obj := args[0].(type) {
	case Counted:
		return Int{I: obj.Count()}, nil
	default:
		s, err := AssertSeqable(env, obj, "count not supported on this type: "+obj.GetType().Name())
		if err != nil {
			return nil, err
		}
		c, err := SeqCount(env, s.Seq())
		if err != nil {
			return nil, err
		}
		return Int{I: c}, nil
	}
}

var procSubvec = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 3, 3); err != nil {
		return nil, err
	}

	// TODO: implement proper Subvector structure
	v, err := EnsureVector(env, args, 0)
	if err != nil {
		return nil, err
	}
	start, err := EnsureInt(env, args, 1)
	if err != nil {
		return nil, err
	}
	end, err := EnsureInt(env, args, 2)
	if err != nil {
		return nil, err
	}
	if start.I > end.I {
		return nil, env.RT.NewError(fmt.Sprintf("subvec's start index (%d) is greater than end index (%d)", start.I, end.I))
	}
	subv := make([]Object, 0, end.I-start.I)
	for i := start.I; i < end.I; i++ {
		subv = append(subv, v.at(i))
	}
	return NewVectorFrom(subv...), nil
}

func mustStr(env *Env, obj Object) string {
	str, err := obj.ToString(env, false)
	if err != nil {
		return fmt.Sprintf("%T(%p)", obj, obj)
	}

	return str
}

var procCast = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	t, err := EnsureType(env, args, 0)
	if err != nil {
		return nil, err
	}
	if t.reflectType.Kind() == reflect.Interface &&
		args[1].GetType().reflectType.Implements(t.reflectType) ||
		args[1].GetType().reflectType == t.reflectType {
		return args[1], nil
	}
	return nil, env.RT.NewError("Cannot cast " + args[1].GetType().Name() + " to " + mustStr(env, t))
}

var procVec = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	sq, err := EnsureSeqable(env, args, 0)
	if err != nil {
		return nil, err
	}
	return NewVectorFromSeq(env, sq.Seq())
}

var procHashMap = func(env *Env, args []Object) (Object, error) {
	if len(args)%2 != 0 {
		return nil, env.RT.NewError("No value supplied for key " + mustStr(env, args[len(args)-1]))
	}
	return NewHashMap(env, args...)
}

var procHashSet = func(env *Env, args []Object) (Object, error) {
	res := EmptySet()
	for i := 0; i < len(args); i++ {
		_, err := res.Add(env, args[i])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

var procStr = func(env *Env, args []Object) (Object, error) {
	var buffer bytes.Buffer
	for _, obj := range args {
		if !obj.Equals(env, NIL) {
			t := obj.GetType()
			// TODO: this is a hack. Rethink escape parameter in ToString
			escaped := (t == TYPE.String) || (t == TYPE.Char) || (t == TYPE.Regex)
			s, err := obj.ToString(env, !escaped)
			if err != nil {
				return nil, err
			}
			buffer.WriteString(s)
		}
	}
	return String{S: buffer.String()}, nil
}

var procSymbol = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 2); err != nil {
		return nil, err
	}

	if len(args) == 1 {
		s, err := EnsureString(env, args, 0)
		if err != nil {
			return nil, err
		}
		return MakeSymbol(s.S), nil
	}

	var ns *string = nil
	if !args[0].Equals(env, NIL) {
		se, err := EnsureString(env, args, 0)
		if err != nil {
			return nil, err
		}
		ns = STRINGS.Intern(se.S)
	}
	name, err := EnsureString(env, args, 1)
	if err != nil {
		return nil, err
	}
	return Symbol{
		ns:   ns,
		name: STRINGS.Intern(name.S),
	}, nil
}

var procKeyword = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 2); err != nil {
		return nil, err
	}

	if len(args) == 1 {
		switch obj := args[0].(type) {
		case String:
			return MakeKeyword(obj.S), nil
		case Symbol:
			return Keyword{
				ns:   obj.ns,
				name: obj.name,
				hash: hashSymbol(obj.ns, obj.name) ^ KeywordHashMask,
			}, nil
		default:
			return NIL, nil
		}
	}

	var ns *string = nil
	if !args[0].Equals(env, NIL) {
		s, err := EnsureString(env, args, 0)
		if err != nil {
			return nil, err
		}
		ns = STRINGS.Intern(s.S)
	}
	sn, err := EnsureString(env, args, 1)
	if err != nil {
		return nil, err
	}
	name := STRINGS.Intern(sn.S)
	return Keyword{
		ns:   ns,
		name: name,
		hash: hashSymbol(ns, name) ^ KeywordHashMask,
	}, nil
}

var procGensym = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	s, err := EnsureString(env, args, 0)
	if err != nil {
		return nil, err
	}
	return genSym(s.S, ""), nil
}

var procApply = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	// TODO:
	// Stacktrace is broken. Need to somehow know
	// the name of the function passed ...
	f, err := EnsureCallable(env, args, 0)
	if err != nil {
		return nil, err
	}
	sq, err := EnsureSeqable(env, args, 1)
	if err != nil {
		return nil, err
	}

	slice, err := ToSlice(env, sq.Seq())
	if err != nil {
		return nil, err
	}

	return f.Call(env, slice)
}

var procLazySeq = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	var fn *Fn
	if err := Cast(env, args[0], &fn); err != nil {
		return nil, err
	}

	return &LazySeq{
		fn: fn,
	}, nil
}

var procDelay = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	var fn *Fn
	if err := Cast(env, args[0], &fn); err != nil {
		return nil, err
	}

	return &Delay{
		fn: fn,
	}, nil
}

var procForce = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	switch d := args[0].(type) {
	case *Delay:
		return d.Force(env)
	default:
		return d, nil
	}
}

var procIdentical = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	return Boolean{B: args[0] == args[1]}, nil
}

var procCompare = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	k1, k2 := args[0], args[1]
	if k1.Equals(env, k2) {
		return Int{I: 0}, nil
	}
	switch k2.(type) {
	case Nil:
		return Int{I: 1}, nil
	}
	switch k1 := k1.(type) {
	case Nil:
		return Int{I: -1}, nil
	case Comparable:
		cmp, err := k1.Compare(env, k2)
		if err != nil {
			return nil, err
		}
		return Int{I: cmp}, nil
	}
	return nil, env.RT.NewError(fmt.Sprintf("%s (type: %s) is not a Comparable", mustStr(env, k1), k1.GetType().Name()))
}

var procInt = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	switch obj := args[0].(type) {
	case Char:
		return Int{I: int(obj.Ch)}, nil
	case Number:
		return obj.Int(), nil
	default:
		return nil, env.RT.NewError(fmt.Sprintf("Cannot cast %s (type: %s) to Int", mustStr(env, obj), obj.GetType().Name()))
	}
}

var procNumber = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	return AssertNumber(env, args[0], fmt.Sprintf("Cannot cast %s (type: %s) to Number", mustStr(env, args[0]), args[0].GetType().Name()))
}

var procDouble = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	n, err := AssertNumber(env, args[0], fmt.Sprintf("Cannot cast %s (type: %s) to Double", mustStr(env, args[0]), args[0].GetType().Name()))
	if err != nil {
		return nil, err
	}
	return n.Double(), nil
}

var procChar = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	switch c := args[0].(type) {
	case Char:
		return c, nil
	case Number:
		i := c.Int().I
		if i < MIN_RUNE || i > MAX_RUNE {
			return nil, env.RT.NewError(fmt.Sprintf("Value out of range for char: %d", i))
		}
		return Char{Ch: rune(i)}, nil
	default:
		return nil, env.RT.NewError(fmt.Sprintf("Cannot cast %s (type: %s) to Char", mustStr(env, c), c.GetType().Name()))
	}
}

var procBoolean = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	return Boolean{B: ToBool(args[0])}, nil
}

var procNumerator = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	bi, err := EnsureRatio(env, args, 0)
	if err != nil {
		return nil, err
	}
	return &BigInt{b: *bi.r.Num()}, nil
}

var procDenominator = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	bi, err := EnsureRatio(env, args, 0)
	if err != nil {
		return nil, err
	}
	return &BigInt{b: *bi.r.Num()}, nil
}

var procBigInt = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	switch n := args[0].(type) {
	case Number:
		return &BigInt{b: *n.BigInt()}, nil
	case String:
		bi := big.Int{}
		if _, ok := bi.SetString(n.S, 10); ok {
			return &BigInt{b: bi}, nil
		}
		return nil, env.RT.NewError("Invalid number format " + n.S)
	default:
		return nil, env.RT.NewError(fmt.Sprintf("Cannot cast %s (type: %s) to BigInt", mustStr(env, n), n.GetType().Name()))
	}
}

var procBigFloat = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	switch n := args[0].(type) {
	case Number:
		return &BigFloat{b: *n.BigFloat()}, nil
	case String:
		b := big.Float{}
		if _, ok := b.SetString(n.S); ok {
			return &BigFloat{b: b}, nil
		}
		return nil, env.RT.NewError("Invalid number format " + n.S)
	default:
		return nil, env.RT.NewError(fmt.Sprintf("Cannot cast %s (type: %s) to BigFloat", mustStr(env, n), n.GetType().Name()))
	}
}

var procNth = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 3); err != nil {
		return nil, err
	}

	ni, err := EnsureNumber(env, args, 1)
	if err != nil {
		return nil, err
	}

	n := ni.Int().I

	switch coll := args[0].(type) {
	case Indexed:
		if len(args) == 3 {
			return coll.TryNth(env, n, args[2])
		}
		return coll.Nth(env, n)
	case Nil:
		return NIL, nil
	case Sequential:
		switch coll := args[0].(type) {
		case Seqable:
			if len(args) == 3 {
				return SeqTryNth(env, coll.Seq(), n, args[2])
			}
			return SeqNth(env, coll.Seq(), n)
		}
	}
	return nil, env.RT.NewError("nth not supported on this type: " + args[0].GetType().Name())
}

var procLt = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	a, err := AssertNumber(env, args[0], "")
	if err != nil {
		return nil, err
	}
	b, err := AssertNumber(env, args[1], "")
	if err != nil {
		return nil, err
	}
	return Boolean{B: GetOps(a).Combine(GetOps(b)).Lt(a, b)}, nil
}

var procLte = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	a, err := AssertNumber(env, args[0], "")
	if err != nil {
		return nil, err
	}
	b, err := AssertNumber(env, args[1], "")
	if err != nil {
		return nil, err
	}
	return Boolean{B: GetOps(a).Combine(GetOps(b)).Lte(a, b)}, nil
}

var procGt = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	a, err := AssertNumber(env, args[0], "")
	if err != nil {
		return nil, err
	}
	b, err := AssertNumber(env, args[1], "")
	if err != nil {
		return nil, err
	}
	return Boolean{B: GetOps(a).Combine(GetOps(b)).Gt(a, b)}, nil
}

var procGte = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	a, err := AssertNumber(env, args[0], "")
	if err != nil {
		return nil, err
	}
	b, err := AssertNumber(env, args[1], "")
	if err != nil {
		return nil, err
	}
	return Boolean{B: GetOps(a).Combine(GetOps(b)).Gte(a, b)}, nil
}

var procEq = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	a, err := AssertNumber(env, args[0], "")
	if err != nil {
		return nil, err
	}
	b, err := AssertNumber(env, args[1], "")
	if err != nil {
		return nil, err
	}
	return MakeBoolean(numbersEq(a, b)), nil
}

var procMax = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	a, err := AssertNumber(env, args[0], "")
	if err != nil {
		return nil, err
	}
	b, err := AssertNumber(env, args[1], "")
	if err != nil {
		return nil, err
	}
	return Max(a, b), nil
}

var procMin = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	a, err := AssertNumber(env, args[0], "")
	if err != nil {
		return nil, err
	}
	b, err := AssertNumber(env, args[1], "")
	if err != nil {
		return nil, err
	}
	return Min(a, b), nil
}

var procIncEx = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	x, err := EnsureNumber(env, args, 0)
	if err != nil {
		return nil, err
	}
	ops := GetOps(x).Combine(BIGINT_OPS)
	return ops.Add(x, Int{I: 1})
}

var procDecEx = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	x, err := EnsureNumber(env, args, 0)
	if err != nil {
		return nil, err
	}
	ops := GetOps(x).Combine(BIGINT_OPS)
	return ops.Subtract(x, Int{I: 1})
}

var procInc = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	x, err := EnsureNumber(env, args, 0)
	if err != nil {
		return nil, err
	}
	ops := GetOps(x).Combine(INT_OPS)
	return ops.Add(x, Int{I: 1})
}

var procDec = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	x, err := EnsureNumber(env, args, 0)
	if err != nil {
		return nil, err
	}
	ops := GetOps(x).Combine(INT_OPS)
	return ops.Subtract(x, Int{I: 1})
}

var procPeek = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	s, err := AssertStack(env, args[0], "")
	if err != nil {
		return nil, err
	}
	return s.Peek(env)
}

var procPop = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	s, err := AssertStack(env, args[0], "")
	if err != nil {
		return nil, err
	}
	obj, err := s.Pop(env)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

var procContains = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	switch c := args[0].(type) {
	case Gettable:
		ok, _, err := c.Get(env, args[1])
		if err != nil {
			return nil, err
		}
		if ok {
			return Boolean{B: true}, nil
		}
		return Boolean{B: false}, nil
	}
	return nil, env.RT.NewError("contains? not supported on type " + args[0].GetType().Name())
}

var procGet = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 3); err != nil {
		return nil, err
	}

	switch c := args[0].(type) {
	case Gettable:
		ok, v, err := c.Get(env, args[1])
		if err != nil {
			return nil, err
		}
		if ok {
			return v, nil
		}
	}
	if len(args) == 3 {
		return args[2], nil
	}
	return NIL, nil
}

var procDissoc = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	m, err := EnsureMap(env, args, 0)
	if err != nil {
		return nil, err
	}
	return m.Without(env, args[1])
}

var procDisj = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	s, err := EnsureSet(env, args, 0)
	if err != nil {
		return nil, err
	}
	return s.Disjoin(env, args[1])
}

var procFind = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	a, err := EnsureAssociative(env, args, 0)
	if err != nil {
		return nil, err
	}
	res, err := a.EntryAt(env, args[1])
	if err != nil {
		return nil, err
	}
	if res == nil {
		return NIL, nil
	}
	return res, nil
}

var procKeys = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	m, err := EnsureMap(env, args, 0)
	if err != nil {
		return nil, err
	}
	return m.Keys(), nil
}

var procVals = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	m, err := EnsureMap(env, args, 0)
	if err != nil {
		return nil, err
	}
	return m.Vals(), nil
}

var procRseq = func(env *Env, args []Object) (Object, error) {
	r, err := EnsureReversible(env, args, 0)
	if err != nil {
		return nil, err
	}
	return r.Rseq(), nil
}

var procName = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	n, err := EnsureNamed(env, args, 0)
	if err != nil {
		return nil, err
	}
	return String{S: n.Name()}, nil
}

var procNamespace = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	n, err := EnsureNamed(env, args, 0)
	if err != nil {
		return nil, err
	}
	ns := n.Namespace()
	if ns == "" {
		return NIL, nil
	}
	return String{S: ns}, nil
}

var procFindVar = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	sym, err := EnsureSymbol(env, args, 0)
	if err != nil {
		return nil, err
	}
	if sym.ns == nil {
		return nil, env.RT.NewError("find-var argument must be namespace-qualified symbol")
	}
	if v, ok := env.Resolve(sym); ok {
		return v, nil
	}
	return NIL, nil
}

var procSort = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	cmp, err := EnsureComparator(env, args, 0)
	if err != nil {
		return nil, err
	}
	coll, err := EnsureSeqable(env, args, 1)
	if err != nil {
		return nil, err
	}
	sl, err := ToSlice(env, coll.Seq())
	if err != nil {
		return nil, err
	}
	s := SortableSlice{
		env: env,
		s:   sl,
		cmp: cmp,
	}
	sort.Sort(&s)
	if s.err != nil {
		return nil, err
	}
	return &ArraySeq{arr: s.s}, nil
}

var procEval = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	parseContext := &ParseContext{Env: env}
	expr, err := Parse(args[0], parseContext)
	if err != nil {
		return nil, err
	}
	return Eval(env, expr, nil)
}

var procType = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	return args[0].GetType(), nil
}

var procPprint = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	obj := args[0]
	w, err := Assertio_Writer(env, env.stdout.Value, "")
	if err != nil {
		return nil, err
	}
	_, err = pprintObject(env, obj, 0, w)
	if err != nil {
		return nil, err
	}
	fmt.Fprint(w, "\n")
	return NIL, nil
}

func PrintObject(env *Env, obj Object, w io.Writer) {
	printReadably := ToBool(env.printReadably.Value)
	switch obj := obj.(type) {
	case Pprinter:
		obj.Pprint(env, w, 2)
	case Printer:
		obj.Print(w, printReadably)
	default:
		s, err := obj.ToString(env, printReadably)
		if err != nil {
			s = fmt.Sprintf("%T(%p)", obj, obj)
		}
		fmt.Fprint(w, s)
	}
}

var procPr = func(env *Env, args []Object) (Object, error) {
	n := len(args)
	if n > 0 {
		f, err := Assertio_Writer(env, env.stdout.Value, "")
		if err != nil {
			return nil, err
		}
		for _, arg := range args[:n-1] {
			PrintObject(env, arg, f)
			fmt.Fprint(f, " ")
		}
		PrintObject(env, args[n-1], f)
	}
	return NIL, nil
}

var procNewline = func(env *Env, args []Object) (Object, error) {
	f, err := Assertio_Writer(env, env.stdout.Value, "")
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(f)
	return NIL, nil
}

var procFlush = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	switch f := args[0].(type) {
	case *File:
		return NIL, f.Sync()
	}
	return NIL, nil
}

func readFromReader(env *Env, reader io.RuneReader) (Object, error) {
	r := NewReader(reader, "<>")
	obj, err := TryRead(env, r)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

var procRead = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	f, err := Ensureio_RuneReader(env, args, 0)
	if err != nil {
		return nil, err
	}
	return readFromReader(env, f)
}

var procReadString = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}
	s, err := EnsureString(env, args, 0)
	if err != nil {
		return nil, err
	}
	return readFromReader(env, strings.NewReader(s.S))
}

func readLine(r StringReader) (s string, e error) {
	s, e = r.ReadString('\n')
	if e == nil {
		l := len(s)
		if s[l-1] == '\n' {
			l -= 1
			if l > 0 && s[l-1] == '\r' {
				l -= 1
			}
		}
		s = s[0:l]
	} else if s != "" && e == io.EOF {
		e = nil
	}
	return
}

var procReadLine = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 0, 0); err != nil {
		return nil, err
	}
	f, err := AssertStringReader(env, env.stdin.Value, "")
	if err != nil {
		return nil, err
	}
	line, err := readLine(f)
	if err != nil {
		return NIL, nil
	}
	return String{S: line}, nil
}

var procReaderReadLine = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}
	rdr, err := EnsureStringReader(env, args, 0)
	if err != nil {
		return nil, err
	}
	line, err := readLine(rdr)
	if err != nil {
		return NIL, nil
	}
	return String{S: line}, nil
}

var procNanoTime = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 0, 0); err != nil {
		return nil, err
	}

	return &BigInt{b: *big.NewInt(time.Now().UnixNano())}, nil
}

var procMacroexpand1 = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	switch s := args[0].(type) {
	case Seq:
		parseContext := &ParseContext{Env: env}
		return macroexpand1(env, s, parseContext)
	default:
		return s, nil
	}
}

func loadReader(env *Env, reader *Reader) (Object, error) {
	parseContext := &ParseContext{Env: env}
	var lastObj Object = NIL
	for {
		obj, err := TryRead(env, reader)
		if err == io.EOF {
			return lastObj, nil
		}
		if err != nil {
			return nil, err
		}
		expr, err := TryParse(obj, parseContext)
		if err != nil {
			return nil, err
		}
		lastObj, err = TryEval(env, expr)
		if err != nil {
			return nil, err
		}
	}
}

var procLoadString = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	s, err := EnsureString(env, args, 0)
	if err != nil {
		return nil, err
	}
	obj, err := loadReader(env, NewReader(strings.NewReader(s.S), "<string>"))
	if err != nil {
		return nil, err
	}
	return obj, nil
}

var procFindNamespace = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	s, err := EnsureSymbol(env, args, 0)
	if err != nil {
		return nil, err
	}

	ns := env.FindNamespace(s)
	if ns == nil {
		return NIL, nil
	}
	return ns, nil
}

var procCreateNamespace = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	sym, err := EnsureSymbol(env, args, 0)
	if err != nil {
		return nil, err
	}
	res := env.EnsureNamespace(sym)
	// In linter mode the latest create-ns call overrides position info.
	// This is for the cases when (ns ...) is called in .laced/linter.clj file and alike.
	// Also, isUsed needs to be reset in this case.
	if LINTER_MODE {
		if err := Cast(env, res.Name.WithInfo(sym.GetInfo()), &res.Name); err != nil {
			return nil, err
		}
		res.isUsed = false
	}
	return res, nil
}

var procInjectNamespace = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	sym, err := EnsureSymbol(env, args, 0)
	if err != nil {
		return nil, err
	}
	ns := env.EnsureNamespace(sym)
	ns.isUsed = true
	ns.isGloballyUsed = true
	return ns, nil
}

var procRemoveNamespace = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	s, err := EnsureSymbol(env, args, 0)
	if err != nil {
		return nil, err
	}
	ns := env.RemoveNamespace(s)
	if ns == nil {
		return NIL, nil
	}
	return ns, nil
}

var procAllNamespaces = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 0, 0); err != nil {
		return nil, err
	}

	s := make([]Object, 0, len(env.Namespaces))
	for _, ns := range env.Namespaces {
		s = append(s, ns)
	}
	return &ArraySeq{arr: s}, nil
}

var procNamespaceName = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	ns, err := EnsureNamespace(env, args, 0)
	if err != nil {
		return nil, err
	}
	return ns.Name, nil
}

var procNamespaceMap = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	r := &ArrayMap{}

	ns, err := EnsureNamespace(env, args, 0)
	if err != nil {
		return nil, err
	}
	for k, v := range ns.mappings {
		r.Add(env, MakeSymbol(*k), v)
	}
	return r, nil
}

var procNamespaceUnmap = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	ns, err := EnsureNamespace(env, args, 0)
	if err != nil {
		return nil, err
	}
	sym, err := EnsureSymbol(env, args, 1)
	if err != nil {
		return nil, err
	}
	if sym.ns != nil {
		return nil, env.RT.NewError("Can't unintern namespace-qualified symbol")
	}
	delete(ns.mappings, sym.name)
	return NIL, nil
}

var procVarNamespace = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	v, err := EnsureVar(env, args, 0)
	if err != nil {
		return nil, err
	}

	return v.ns, nil
}

var procRefer = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 3, 3); err != nil {
		return nil, err
	}

	ns, err := EnsureNamespace(env, args, 0)
	if err != nil {
		return nil, err
	}
	sym, err := EnsureSymbol(env, args, 1)
	if err != nil {
		return nil, err
	}
	v, err := EnsureVar(env, args, 2)
	if err != nil {
		return nil, err
	}
	return ns.Refer(env, sym, v)
}

var procAlias = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 3, 3); err != nil {
		return nil, err
	}

	ns, err := EnsureNamespace(env, args, 0)
	if err != nil {
		return nil, err
	}
	sym, err := EnsureSymbol(env, args, 1)
	if err != nil {
		return nil, err
	}

	ns2, err := EnsureNamespace(env, args, 2)
	if err != nil {
		return nil, err
	}

	err = ns.AddAlias(env, sym, ns2)
	if err != nil {
		return nil, err
	}
	return NIL, nil
}

var procNamespaceAliases = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	r := &ArrayMap{}
	ns, err := EnsureNamespace(env, args, 0)
	if err != nil {
		return nil, err
	}
	for k, v := range ns.aliases {
		r.Add(env, MakeSymbol(*k), v)
	}
	return r, nil
}

var procNamespaceUnalias = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	ns, err := EnsureNamespace(env, args, 0)
	if err != nil {
		return nil, err
	}
	sym, err := EnsureSymbol(env, args, 1)
	if err != nil {
		return nil, err
	}
	if sym.ns != nil {
		return nil, env.RT.NewError("Alias can't be namespace-qualified")
	}
	delete(ns.aliases, sym.name)
	return NIL, nil
}

var procVarGet = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	v, err := EnsureVar(env, args, 0)
	if err != nil {
		return nil, err
	}
	return v.Resolve(), nil
}

var procVarSet = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	v, err := EnsureVar(env, args, 0)
	if err != nil {
		return nil, err
	}
	v.Value = args[1]
	return args[1], nil
}

var procNsResolve = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	ns, err := EnsureNamespace(env, args, 0)
	if err != nil {
		return nil, err
	}
	sym, err := EnsureSymbol(env, args, 1)
	if err != nil {
		return nil, err
	}
	if sym.ns == nil && TYPES[sym.name] != nil {
		return TYPES[sym.name], nil
	}
	if vr, ok := env.ResolveIn(ns, sym); ok {
		return vr, nil
	}
	return NIL, nil
}

var procArrayMap = func(env *Env, args []Object) (Object, error) {
	if len(args)%2 == 1 {
		return nil, env.RT.NewError("No value supplied for key " + mustStr(env, args[len(args)-1]))
	}
	res := EmptyArrayMap()
	for i := 0; i < len(args); i += 2 {
		res.Set(env, args[i], args[i+1])
	}
	return res, nil
}

var procBuffer = func(env *Env, args []Object) (Object, error) {
	if len(args) > 0 {
		s, err := EnsureString(env, args, 0)
		if err != nil {
			return nil, err
		}
		return MakeBuffer(bytes.NewBufferString(s.S)), nil
	}
	return MakeBuffer(&bytes.Buffer{}), nil
}

var procBufferedReader = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	switch rdr := args[0].(type) {
	case io.Reader:
		return MakeBufferedReader(rdr), nil
	default:
		return nil, env.RT.NewArgTypeError(0, args[0], "IOReader")
	}
}

var procSlurp = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	s, err := EnsureString(env, args, 0)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(s.S)
	if err != nil {
		return nil, err
	}
	return String{S: string(b)}, nil
}

var procSpit = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 3, 3); err != nil {
		return nil, err
	}

	filename, err := EnsureString(env, args, 0)
	if err != nil {
		return nil, err
	}
	content, err := EnsureString(env, args, 1)
	if err != nil {
		return nil, err
	}
	opts, err := EnsureMap(env, args, 2)
	if err != nil {
		return nil, err
	}
	appendFile := false
	ok, appe, err := opts.Get(env, MakeKeyword("append"))
	if err != nil {
		return nil, err
	}

	if ok {
		appendFile = ToBool(appe)
	}
	flags := os.O_CREATE | os.O_WRONLY
	if appendFile {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	f, err := os.OpenFile(filename.S, flags, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.WriteString(content.S)
	if err != nil {
		return nil, err
	}
	return NIL, nil
}

var procShuffle = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	seq, err := EnsureSeqable(env, args, 0)
	if err != nil {
		return nil, err
	}
	s, err := ToSlice(env, seq.Seq())
	if err != nil {
		return nil, err
	}
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
	return NewVectorFrom(s...), nil
}

var procIsRealized = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	p, err := EnsurePending(env, args, 0)
	if err != nil {
		return nil, err
	}
	return Boolean{B: p.IsRealized()}, nil
}

var procDeriveInfo = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	dest := args[0]
	src := args[1]
	return dest.WithInfo(src.GetInfo()), nil
}

var procLaceVersion = func(env *Env, args []Object) (Object, error) {
	return String{S: VERSION[1:]}, nil
}

var procHash = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	h, err := args[0].Hash(env)
	if err != nil {
		return nil, err
	}
	return Int{I: int(h)}, nil
}

func loadFile(env *Env, filename string) (Object, error) {
	var reader *Reader
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader = NewReader(bufio.NewReader(f), filename)
	err = ProcessReaderFromEval(env, reader, filename)
	if err != nil {
		return nil, err
	}
	return NIL, nil
}

var procLoadFile = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	filename, err := EnsureString(env, args, 0)
	if err != nil {
		return nil, err
	}
	return loadFile(env, filename.S)
}

var procLoadLibFromPath = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	libnamev, err := EnsureSymbol(env, args, 0)
	if err != nil {
		return nil, err
	}
	libname := libnamev.Name()
	pathnamev, err := EnsureString(env, args, 1)
	if err != nil {
		return nil, err
	}
	pathname := pathnamev.S

	cp := env.classPath.Value
	cpvec, err := AssertVector(env, cp, "*classpath* must be a Vector, not a "+cp.GetType().Name())
	if err != nil {
		return nil, err
	}

	count := cpvec.Count()
	var f *os.File
	var canonicalErr error
	var filename string
	for i := 0; i < count; i++ {
		elem := cpvec.at(i)
		cpelem, err := AssertString(env, elem, "*classpath* must contain only Strings, not a "+elem.GetType().Name()+" (at element "+strconv.Itoa(i)+")")
		if err != nil {
			return nil, err
		}
		s := cpelem.S
		if s == "" {
			filename = pathname
		} else {
			filename = filepath.Join(s, filepath.Join(strings.Split(libname, ".")...)) + ".clj" // could cache inner join....
		}
		f, err = os.Open(filename)
		if err == nil {
			canonicalErr = nil
			break
		}
		if s == "" {
			canonicalErr = err
		}
	}
	if canonicalErr != nil {
		return nil, canonicalErr
	}
	if err != nil {
		return nil, err
	}
	reader := NewReader(bufio.NewReader(f), filename)
	err = ProcessReaderFromEval(env, reader, filename)
	if err != nil {
		return nil, err
	}
	return NIL, nil
}

var procReduceKv = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 3, 3); err != nil {
		return nil, err
	}

	f, err := EnsureCallable(env, args, 0)
	if err != nil {
		return nil, err
	}
	init := args[1]
	coll, err := EnsureKVReduce(env, args, 2)
	if err != nil {
		return nil, err
	}
	return coll.kvreduce(env, f, init)
}

var procIndexOf = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}

	s, err := EnsureString(env, args, 0)
	if err != nil {
		return nil, err
	}
	ch, err := EnsureChar(env, args, 1)
	if err != nil {
		return nil, err
	}
	for i, r := range s.S {
		if r == ch.Ch {
			return Int{I: i}, nil
		}
	}
	return Int{I: -1}, nil
}

func libExternalPath(env *Env, sym Symbol) (path string, ok bool, err error) {
	nsSourcesVar, _ := env.Resolve(MakeSymbol("lace.core/*ns-sources*"))
	var vec *Vector
	if err := Cast(env, nsSourcesVar.Value, &vec); err != nil {
		return "", false, err
	}

	nsSources, err := ToSlice(env, vec.Seq())
	if err != nil {
		return "", false, err
	}

	var sourceKey string
	var sourceMap Map
	for _, source := range nsSources {
		var svec *Vector
		if err := Cast(env, source, &svec); err != nil {
			return "", false, err
		}

		n, err := svec.Nth(env, 0)
		if err != nil {
			return "", false, err
		}

		sourceKey, err = n.ToString(env, false)
		if err != nil {
			return "", false, err
		}
		match, _ := regexp.MatchString(sourceKey, sym.Name())
		if match {
			n, err := svec.Nth(env, 0)
			if err != nil {
				return "", false, err
			}
			if err := Cast(env, n, &sourceMap); err != nil {
				return "", false, err
			}
			break
		}
	}
	if sourceMap != nil {
		ok, url, err := sourceMap.Get(env, MakeKeyword("url"))
		if err != nil {
			return "", false, err
		}
		if !ok {
			return "", false, env.RT.NewError("Key :url not found in ns-sources for: " + sourceKey)
		} else {
			s, err := url.ToString(env, false)
			if err != nil {
				return "", false, err
			}

			path, err := externalSourceToPath(env, sym.Name(), s)
			if err != nil {
				return "", false, err
			}
			return path, true, nil
		}
	}
	return
}

var procLibPath = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	sym, err := EnsureSymbol(env, args, 0)
	if err != nil {
		return nil, err
	}
	var path string

	path, ok, err := libExternalPath(env, sym)
	if err != nil {
		return nil, err
	}

	if !ok {
		var file string
		if env.file.Value == nil {
			var err error
			file, err = filepath.Abs("user")
			if err != nil {
				return nil, err
			}
		} else {
			filev, err := AssertString(env, env.file.Value, "")
			if err != nil {
				return nil, err
			}
			file = filev.S
			if linkDest, err := os.Readlink(file); err == nil {
				file = linkDest
			}
		}
		ns := env.CurrentNamespace().Name

		parts := strings.Split(ns.Name(), ".")
		for range parts {
			file, _ = filepath.Split(file)
			file = file[:len(file)-1]
		}
		path = filepath.Join(append([]string{file}, strings.Split(sym.Name(), ".")...)...) + ".clj"
	}
	return String{S: path}, nil
}

var procInternFakeVar = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 3, 3); err != nil {
		return nil, err
	}

	nsSym, err := EnsureSymbol(env, args, 0)
	if err != nil {
		return nil, err
	}
	sym, err := EnsureSymbol(env, args, 1)
	if err != nil {
		return nil, err
	}
	isMacro := ToBool(args[2])
	res, err := InternFakeSymbol(env, env.FindNamespace(nsSym), sym)
	if err != nil {
		return nil, err
	}
	res.isMacro = isMacro
	return res, nil
}

var procParse = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	lm, _ := env.Resolve(MakeSymbol("lace.core/*linter-mode*"))
	lm.Value = Boolean{B: true}
	LINTER_MODE = true
	defer func() {
		LINTER_MODE = false
		lm.Value = Boolean{B: false}
	}()
	parseContext := &ParseContext{Env: env}
	res, err := Parse(args[0], parseContext)
	if err != nil {
		return nil, err
	}
	return res.Dump(false), nil
}

var procTypes = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 0, 0); err != nil {
		return nil, err
	}
	res := EmptyArrayMap()
	for k, v := range TYPES {
		res.Add(env, String{S: *k}, v)
	}
	return res, nil
}

var procCreateChan = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}
	n, err := EnsureInt(env, args, 0)
	if err != nil {
		return nil, err
	}
	ch := make(chan FutureResult, n.I)
	return MakeChannel(ch), nil
}

var procCloseChan = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}
	c, err := EnsureChannel(env, args, 0)
	if err != nil {
		return nil, err
	}

	c.Close()
	return NIL, nil
}

var procSend = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 2, 2); err != nil {
		return nil, err
	}
	ch, err := EnsureChannel(env, args, 0)
	if err != nil {
		return nil, err
	}
	v := args[1]
	if v.Equals(env, NIL) {
		return nil, env.RT.NewError("Can't put nil on channel")
	}
	if ch.isClosed {
		return MakeBoolean(false), nil
	}
	obj := MakeBoolean(true)
	//env.RT.GIL.Unlock()
	ch.ch <- MakeFutureResult(v, nil)
	//env.RT.GIL.Lock()
	return obj, nil
}

var procReceive = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}
	ch, err := EnsureChannel(env, args, 0)
	if err != nil {
		return nil, err
	}
	//env.RT.GIL.Unlock()
	res, ok := <-ch.ch
	//env.RT.GIL.Lock()
	if !ok {
		return NIL, nil
	}
	if res.err != nil {
		return nil, err
	}
	return res.value, nil
}

var procGo = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}
	f, err := EnsureCallable(env, args, 0)
	if err != nil {
		return nil, err
	}

	ch := MakeChannel(make(chan FutureResult, 1))
	go func() {
		var cerr Error
		res, err := f.Call(env, []Object{})
		if err != nil {
			cerr = env.RT.NewError(err.Error())
		}

		ch.ch <- MakeFutureResult(res, cerr)
		ch.Close()
	}()
	return ch, nil
}

var procVerbosityLevel = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 0, 0); err != nil {
		return nil, err
	}
	return MakeInt(VerbosityLevel), nil
}

var procExit = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 0, 1); err != nil {
		return nil, err
	}

	if len(args) == 1 {
		i, err := EnsureInt(env, args, 0)
		if err != nil {
			return nil, err
		}

		return nil, &ExitError{Code: i.I}
	}

	return nil, &ExitError{}
}

func PackReader(env *Env, reader *Reader, filename string) ([]byte, error) {
	var p []byte
	packEnv := NewPackEnv(env)
	parseContext := &ParseContext{Env: env}
	if filename != "" {
		currentFilename := parseContext.Env.file.Value
		defer func() {
			parseContext.Env.SetFilename(currentFilename)
		}()
		s, err := filepath.Abs(filename)
		if err != nil {
			return nil, err
		}
		parseContext.Env.SetFilename(MakeString(s))
	}

	for {
		obj, err := TryRead(env, reader)
		if err == io.EOF {
			var hp []byte
			hp = packEnv.Pack(hp)
			return append(hp, p...), nil
		}
		if err != nil {
			fmt.Fprintln(Stderr, err)
			return nil, err
		}
		expr, err := TryParse(obj, parseContext)
		if err != nil {
			fmt.Fprintln(Stderr, err)
			return nil, err
		}
		p = expr.Pack(p, packEnv)
		_, err = TryEval(env, expr)
		if err != nil {
			fmt.Fprintln(Stderr, err)
			return nil, err
		}
	}
}

var procIncProblemCount = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 0, 0); err != nil {
		return nil, err
	}

	PROBLEM_COUNT++
	return NIL, nil
}

func ProcessReader(env *Env, reader *Reader, filename string, phase Phase) error {
	parseContext := &ParseContext{Env: env}
	if filename != "" {
		currentFilename := parseContext.Env.file.Value
		defer func() {
			parseContext.Env.SetFilename(currentFilename)
		}()
		s, err := filepath.Abs(filename)
		if err != nil {
			return err
		}
		parseContext.Env.SetFilename(MakeString(s))
	}

	var exprs []Expr

	for {
		obj, err := TryRead(env, reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(Stderr, err)
			return err
		}
		if phase == READ {
			continue
		}
		expr, err := TryParse(obj, parseContext)
		if err != nil {
			fmt.Fprintln(Stderr, err)
			return err
		}
		if phase == PARSE {
			continue
		}

		exprs = append(exprs, expr)
	}

	/*
		for _, expr := range exprs {
			env.RT.setTopPosition(expr.Pos())

			obj, err := TryEval(env, expr)
			if err != nil {
				fmt.Fprintln(Stderr, err)
				return err
			}
			if phase == EVAL {
				continue
			}
			if _, ok := obj.(Nil); !ok {
				s, err := obj.ToString(env, true)
				if err != nil {
					return err
				}
				fmt.Fprintln(Stdout, s)
			}
		}
	*/

	fn := &Fn{
		fnExpr: &FnExpr{
			arities: []FnArityExpr{
				{
					body: exprs,
				},
			},
		},
	}

	_, err := compileFn(env, fn, nil)
	if err != nil {
		fmt.Printf("error compiling: %s\n", err)
		return err
	}

	//spew.Dump(code)

	_, err = EngineRun(env, fn)
	if err != nil {
		fmt.Printf("error running: %s\n", err)
		return err
	}

	return nil
	/*
		for {
			obj, err := TryRead(env, reader)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				fmt.Fprintln(Stderr, err)
				return err
			}
			if phase == READ {
				continue
			}
			expr, err := TryParse(obj, parseContext)
			if err != nil {
				fmt.Fprintln(Stderr, err)
				return err
			}
			if phase == PARSE {
				continue
			}

			env.RT.setTopPosition(expr.Pos())

			obj, err = TryEval(env, expr)
			if err != nil {
				fmt.Fprintln(Stderr, err)
				return err
			}
			if phase == EVAL {
				continue
			}
			if _, ok := obj.(Nil); !ok {
				s, err := obj.ToString(env, true)
				if err != nil {
					return err
				}
				fmt.Fprintln(Stdout, s)
			}
		}
	*/
}

func ProcessReaderFromEval(env *Env, reader *Reader, filename string) error {
	parseContext := &ParseContext{Env: env}
	if filename != "" {
		currentFilename := parseContext.Env.file.Value
		defer func() {
			parseContext.Env.SetFilename(currentFilename)
		}()
		s, err := filepath.Abs(filename)
		if err != nil {
			return err
		}
		parseContext.Env.SetFilename(MakeString(s))
	}
	for {
		obj, err := TryRead(env, reader)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		expr, err := TryParse(obj, parseContext)
		if err != nil {
			return err
		}
		_, err = TryEval(env, expr)
		if err != nil {
			return err
		}
	}
}

func processInEnvInNS(env *Env, ns *Namespace, data []byte) error {
	cur := env.CurrentNamespace()
	env.SetCurrentNamespace(ns)
	defer func() { env.SetCurrentNamespace(cur) }()

	header, p, err := UnpackHeader(data, env)
	if err != nil {
		return err
	}

	var exprs []Expr

	for len(p) > 0 {
		var expr Expr
		expr, p, err = UnpackExpr(env, p, header)
		if err != nil {
			return err
		}

		exprs = append(exprs, expr)
	}

	/*
		fn := &Fn{
			fnExpr: &FnExpr{
				arities: []FnArityExpr{
					{
						body: exprs,
					},
				},
			},
		}

		_, err = compileFn(env, fn, nil)
		if err != nil {
			fmt.Printf("error compiling: %s\n", err)
			return err
		}

		//spew.Dump(code)

		_, err = EngineRun(env, fn)
		if err != nil {
			fmt.Printf("error running: %s\n", err)
			return err
		}

	*/
	for _, expr := range exprs {
		_, err := TryEval(env, expr)
		if err != nil {
			return err
		}
	}
	if VerbosityLevel > 0 {
		fmt.Fprintf(Stderr, "processData: Evaluated code for %s\n", env.CurrentNamespace().Qual())
	}

	return nil
}

func setCoreNamespaces(env *Env) error {
	ns := env.CoreNamespace
	ns.MaybeLazy(env, "lace.core")

	var set *MapSet

	// Add 'lace.core to *loaded-libs*, now that it's loaded.
	vr := ns.Resolve("*loaded-libs*")
	if err := Cast(env, vr.Value, &set); err != nil {
		return err
	}
	v, err := set.Conj(env, ns.Name)
	if err != nil {
		return err
	}
	if err := Cast(env, v, &set); err != nil {
		return err
	}
	vr.Value = set
	return nil
}

var procIsNamespaceInitialized = func(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	sym, err := EnsureSymbol(env, args, 0)
	if err != nil {
		return nil, err
	}

	if sym.ns != nil {
		return nil, env.RT.NewError("Can't ask for namespace info on namespace-qualified symbol")
	}
	// First look for registered (e.g. std) libs
	ns, found := env.Namespaces[sym.name]
	return MakeBoolean(found && ns.Lazy == nil), nil
}

func findConfigFile(filename string, workingDir string, findDir bool) string {
	var err error
	configName := ".lace"
	if findDir {
		configName = ".laced"
	}
	if filename != "" {
		filename, err = filepath.Abs(filename)
		if err != nil {
			fmt.Fprintln(Stderr, "Error reading config file "+filename+": ", err)
			return ""
		}
	}

	if workingDir != "" {
		workingDir, err := filepath.Abs(workingDir)
		if err != nil {
			fmt.Fprintln(Stderr, "Error resolving working directory"+workingDir+": ", err)
			return ""
		}
		filename = filepath.Join(workingDir, configName)
	}
	for {
		oldFilename := filename
		filename = filepath.Dir(filename)
		if filename == oldFilename {
			home, ok := os.LookupEnv("HOME")
			if !ok {
				home, ok = os.LookupEnv("USERPROFILE")
				if !ok {
					return ""
				}
			}
			p := filepath.Join(home, configName)
			if info, err := os.Stat(p); err == nil {
				if !findDir || info.IsDir() {
					return p
				}
			}
			return ""
		}
		p := filepath.Join(filename, configName)
		if info, err := os.Stat(p); err == nil {
			if !findDir || info.IsDir() {
				return p
			}
		}
	}
}

func printConfigError(filename, msg string) {
	fmt.Fprintln(Stderr, "Error reading config file "+filename+": ", msg)
}

func knownMacrosToMap(env *Env, km Object) (Map, error) {
	var seqo Seqable
	if err := Cast(env, km, &seqo); err != nil {
		return nil, err
	}

	s := seqo.Seq()
	res := EmptyArrayMap()
	for {
		empty, err := s.IsEmpty(env)
		if err != nil {
			return nil, err
		}

		if empty {
			break
		}

		obj, err := s.First(env)
		if err != nil {
			return nil, err
		}
		switch obj := obj.(type) {
		case Symbol:
			res.Add(env, obj, NIL)
		case *Vector:
			if obj.Count() != 2 {
				return nil, errors.New(":known-macros item must be a symbol or a vector with two elements")
			}
			res.Add(env, obj.at(0), obj.at(1))
		default:
			return nil, errors.New(":known-macros item must be a symbol or a vector, got " + obj.GetType().Name())
		}
		s, err = s.Rest(env)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func ReadConfig(env *Env, filename string, workingDir string) error {
	var err error
	LINTER_CONFIG, err = env.CoreNamespace.Intern(env, MakeSymbol("*linter-config*"))
	if err != nil {
		return err
	}
	LINTER_CONFIG.Value = EmptyArrayMap()
	configFileName := findConfigFile(filename, workingDir, false)
	if configFileName == "" {
		return nil
	}
	f, err := os.Open(configFileName)
	if err != nil {
		printConfigError(configFileName, err.Error())
		return err
	}
	r := NewReader(bufio.NewReader(f), configFileName)
	config, err := TryRead(env, r)
	if err != nil {
		printConfigError(configFileName, err.Error())
		return err
	}
	configMap, ok := config.(Map)
	if !ok {
		printConfigError(configFileName, "config root object must be a map, got "+config.GetType().Name())
		return nil
	}
	ok, ignoredUnusedNamespaces := configMap.GetEqu(MakeKeyword("ignored-unused-namespaces"))
	if ok {
		seq, ok1 := ignoredUnusedNamespaces.(Seqable)
		if ok1 {
			WARNINGS.ignoredUnusedNamespaces, err = NewSetFromSeq(env, seq.Seq())
			if err != nil {
				return err
			}
		} else {
			printConfigError(configFileName, ":ignored-unused-namespaces value must be a vector, got "+ignoredUnusedNamespaces.GetType().Name())
			return nil
		}
	}
	ok, ignoredFileRegexes := configMap.GetEqu(MakeKeyword("ignored-file-regexes"))
	if ok {
		seq, ok1 := ignoredFileRegexes.(Seqable)
		if ok1 {
			s := seq.Seq()
			for {
				empty, err := s.IsEmpty(env)
				if err != nil {
					return err
				}

				if empty {
					break
				}

				f, err := s.First(env)
				if err != nil {
					return err
				}
				regex, ok2 := f.(*Regex)
				if !ok2 {
					f, err := s.First(env)
					if err != nil {
						return err
					}

					printConfigError(configFileName, ":ignored-file-regexes elements must be regexes, got "+f.GetType().Name())
					return nil
				}
				WARNINGS.IgnoredFileRegexes = append(WARNINGS.IgnoredFileRegexes, regex.R)
				s, err = s.Rest(env)
				if err != nil {
					return err
				}
			}
		} else {
			printConfigError(configFileName, ":ignored-file-regexes value must be a vector, got "+ignoredFileRegexes.GetType().Name())
			return nil
		}
	}
	ok, entryPoints := configMap.GetEqu(MakeKeyword("entry-points"))
	if ok {
		seq, ok1 := entryPoints.(Seqable)
		if ok1 {
			WARNINGS.entryPoints, err = NewSetFromSeq(env, seq.Seq())
			if err != nil {
				return err
			}
		} else {
			printConfigError(configFileName, ":entry-points value must be a vector, got "+entryPoints.GetType().Name())
			return nil
		}
	}
	ok, knownNamespaces := configMap.GetEqu(MakeKeyword("known-namespaces"))
	if ok {
		if _, ok1 := knownNamespaces.(Seqable); !ok1 {
			printConfigError(configFileName, ":known-namespaces value must be a vector, got "+knownNamespaces.GetType().Name())
			return nil
		}
	}
	ok, knownTags := configMap.GetEqu(MakeKeyword("known-tags"))
	if ok {
		if _, ok1 := knownTags.(Seqable); !ok1 {
			printConfigError(configFileName, ":known-tags value must be a vector, got "+knownTags.GetType().Name())
			return nil
		}
	}
	ok, knownMacros := configMap.GetEqu(criticalKeywords.knownMacros)
	if ok {
		_, ok1 := knownMacros.(Seqable)
		if !ok1 {
			printConfigError(configFileName, ":known-macros value must be a vector, got "+knownMacros.GetType().Name())
			return nil
		}
		m, err := knownMacrosToMap(env, knownMacros)
		if err != nil {
			printConfigError(configFileName, err.Error())
			return nil
		}
		v, err := configMap.Assoc(env, criticalKeywords.knownMacros, m)
		if err != nil {
			return err
		}
		if err := Cast(env, v, &configMap); err != nil {
			return err
		}
	}
	ok, rules := configMap.GetEqu(criticalKeywords.rules)
	if ok {
		m, ok := rules.(Map)
		if !ok {
			printConfigError(configFileName, ":rules value must be a map, got "+rules.GetType().Name())
			return nil
		}
		if ok, v := m.GetEqu(criticalKeywords.ifWithoutElse); ok {
			WARNINGS.ifWithoutElse = ToBool(v)
		}
		if ok, v := m.GetEqu(criticalKeywords.unusedFnParameters); ok {
			WARNINGS.unusedFnParameters = ToBool(v)
		}
		if ok, v := m.GetEqu(criticalKeywords.fnWithEmptyBody); ok {
			WARNINGS.fnWithEmptyBody = ToBool(v)
		}
	}
	LINTER_CONFIG.Value = configMap
	return nil
}

func removeLaceNamespaces(env *Env) {
	for k, ns := range env.Namespaces {
		if ns != env.CoreNamespace && strings.HasPrefix(*k, "lace.") {
			delete(env.Namespaces, k)
		}
	}
}

func markLaceNamespacesAsUsed(env *Env) {
	for k, ns := range env.Namespaces {
		if ns != env.CoreNamespace && strings.HasPrefix(*k, "lace.") {
			ns.isUsed = true
			ns.isGloballyUsed = true
		}
	}
}

func NewReaderFromFile(filename string) (*Reader, error) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(Stderr, "Error: ", err)
		return nil, err
	}
	return NewReader(bufio.NewReader(f), filename), nil
}

func ProcessLinterFile(env *Env, configDir string, filename string) error {
	linterFileName := filepath.Join(configDir, filename)
	if _, err := os.Stat(linterFileName); err == nil {
		if reader, err := NewReaderFromFile(linterFileName); err == nil {
			err := ProcessReader(env, reader, linterFileName, EVAL)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ProcessLinterFiles(env *Env, dialect Dialect, filename string, workingDir string) error {
	if dialect == EDN || dialect == LACE {
		return nil
	}
	configDir := findConfigFile(filename, workingDir, true)
	if configDir == "" {
		return nil
	}
	if err := ProcessLinterFile(env, configDir, "linter.cljc"); err != nil {
		return err
	}
	switch dialect {
	case CLJS:
		if err := ProcessLinterFile(env, configDir, "linter.cljs"); err != nil {
			return err
		}
	case CLJ:
		if err := ProcessLinterFile(env, configDir, "linter.clj"); err != nil {
			return err
		}
	}

	return nil
}
