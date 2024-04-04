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
	JOKER
	EDN
	UNKNOWN
)

func ExtractCallable(env *Env, args []Object, index int) Callable {
	return EnsureCallable(env, args, index)
}

func ExtractObject(env *Env, args []Object, index int) Object {
	return args[index]
}

func ExtractString(env *Env, args []Object, index int) string {
	return EnsureString(env, args, index).S
}

func ExtractKeyword(env *Env, args []Object, index int) string {
	return EnsureKeyword(env, args, index).ToString(false)
}

func ExtractStringable(env *Env, args []Object, index int) string {
	return EnsureStringable(args, index).S
}

func ExtractStrings(env *Env, args []Object, index int) []string {
	strs := make([]string, 0)
	for i := index; i < len(args); i++ {
		strs = append(strs, EnsureString(env, args, i).S)
	}
	return strs
}

func ExtractInt(env *Env, args []Object, index int) int {
	return EnsureInt(env, args, index).I
}

func ExtractBoolean(env *Env, args []Object, index int) bool {
	return EnsureBoolean(env, args, index).B
}

func ExtractChar(env *Env, args []Object, index int) rune {
	return EnsureChar(env, args, index).Ch
}

func ExtractTime(env *Env, args []Object, index int) time.Time {
	return EnsureTime(env, args, index).T
}

func ExtractDouble(env *Env, args []Object, index int) float64 {
	return EnsureDouble(env, args, index).D
}

func ExtractNumber(env *Env, args []Object, index int) Number {
	return EnsureNumber(env, args, index)
}

func ExtractRegex(env *Env, args []Object, index int) *regexp.Regexp {
	return EnsureRegex(env, args, index).R
}

func ExtractSeqable(env *Env, args []Object, index int) Seqable {
	return EnsureSeqable(env, args, index)
}

func ExtractMap(env *Env, args []Object, index int) Map {
	return EnsureMap(env, args, index)
}

func ExtractIOReader(env *Env, args []Object, index int) io.Reader {
	return Ensureio_Reader(env, args, index)
}

func ExtractIOWriter(env *Env, args []Object, index int) io.Writer {
	return Ensureio_Writer(env, args, index)
}

var procMeta = func(env *Env, args []Object) Object {
	switch obj := args[0].(type) {
	case Meta:
		meta := obj.GetMeta()
		if meta != nil {
			return meta
		}
	case *Type:
		meta := obj.GetMeta()
		if meta != nil {
			return meta
		}
	}
	return NIL
}

var procWithMeta = func(env *Env, args []Object) Object {
	CheckArity(env, args, 2, 2)
	m := EnsureMeta(env, args, 0)
	if args[1].Equals(NIL) {
		return args[0]
	}
	return m.WithMeta(EnsureMap(env, args, 1))
}

var procIsZero = func(env *Env, args []Object) Object {
	n := EnsureNumber(env, args, 0)
	ops := GetOps(n)
	return Boolean{B: ops.IsZero(n)}
}

var procIsPos = func(env *Env, args []Object) Object {
	n := EnsureNumber(env, args, 0)
	ops := GetOps(n)
	return Boolean{B: ops.Gt(n, Int{I: 0})}
}

var procIsNeg = func(env *Env, args []Object) Object {
	n := EnsureNumber(env, args, 0)
	ops := GetOps(n)
	return Boolean{B: ops.Lt(n, Int{I: 0})}
}

var procAdd = func(env *Env, args []Object) Object {
	x := AssertNumber(env, args[0], "")
	y := AssertNumber(env, args[1], "")
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Add(x, y)
}

var procAddEx = func(env *Env, args []Object) Object {
	x := AssertNumber(env, args[0], "")
	y := AssertNumber(env, args[1], "")
	ops := GetOps(x).Combine(GetOps(y)).Combine(BIGINT_OPS)
	return ops.Add(x, y)
}

var procMultiply = func(env *Env, args []Object) Object {
	x := AssertNumber(env, args[0], "")
	y := AssertNumber(env, args[1], "")
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Multiply(x, y)
}

var procMultiplyEx = func(env *Env, args []Object) Object {
	x := AssertNumber(env, args[0], "")
	y := AssertNumber(env, args[1], "")
	ops := GetOps(x).Combine(GetOps(y)).Combine(BIGINT_OPS)
	return ops.Multiply(x, y)
}

var procSubtract = func(env *Env, args []Object) Object {
	var a, b Object
	if len(args) == 1 {
		a = Int{I: 0}
		b = args[0]
	} else {
		a = args[0]
		b = args[1]
	}
	ops := GetOps(a).Combine(GetOps(b))
	return ops.Subtract(AssertNumber(env, a, ""), AssertNumber(env, b, ""))
}

var procSubtractEx = func(env *Env, args []Object) Object {
	var a, b Object
	if len(args) == 1 {
		a = Int{I: 0}
		b = args[0]
	} else {
		a = args[0]
		b = args[1]
	}
	ops := GetOps(a).Combine(GetOps(b)).Combine(BIGINT_OPS)
	return ops.Subtract(AssertNumber(env, a, ""), AssertNumber(env, b, ""))
}

var procDivide = func(env *Env, args []Object) Object {
	x := EnsureNumber(env, args, 0)
	y := EnsureNumber(env, args, 1)
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Divide(x, y)
}

var procQuot = func(env *Env, args []Object) Object {
	x := EnsureNumber(env, args, 0)
	y := EnsureNumber(env, args, 1)
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Quotient(x, y)
}

var procRem = func(env *Env, args []Object) Object {
	x := EnsureNumber(env, args, 0)
	y := EnsureNumber(env, args, 1)
	ops := GetOps(x).Combine(GetOps(y))
	return ops.Rem(x, y)
}

var procBitNot = func(env *Env, args []Object) Object {
	x := AssertInt(env, args[0], "Bit operation not supported for "+args[0].GetType().ToString(false))
	return Int{I: ^x.I}
}

func AssertInts(env *Env, args []Object) (Int, Int) {
	x := AssertInt(env, args[0], "Bit operation not supported for "+args[0].GetType().ToString(false))
	y := AssertInt(env, args[1], "Bit operation not supported for "+args[1].GetType().ToString(false))
	return x, y
}

var procBitAnd = func(env *Env, args []Object) Object {
	x, y := AssertInts(env, args)
	return Int{I: x.I & y.I}
}

var procBitOr = func(env *Env, args []Object) Object {
	x, y := AssertInts(env, args)
	return Int{I: x.I | y.I}
}

var procBitXor = func(env *Env, args []Object) Object {
	x, y := AssertInts(env, args)
	return Int{I: x.I ^ y.I}
}

var procBitAndNot = func(env *Env, args []Object) Object {
	x, y := AssertInts(env, args)
	return Int{I: x.I &^ y.I}
}

var procBitClear = func(env *Env, args []Object) Object {
	x, y := AssertInts(env, args)
	return Int{I: x.I &^ (1 << uint(y.I))}
}

var procBitSet = func(env *Env, args []Object) Object {
	x, y := AssertInts(env, args)
	return Int{I: x.I | (1 << uint(y.I))}
}

var procBitFlip = func(env *Env, args []Object) Object {
	x, y := AssertInts(env, args)
	return Int{I: x.I ^ (1 << uint(y.I))}
}

var procBitTest = func(env *Env, args []Object) Object {
	x, y := AssertInts(env, args)
	return Boolean{B: x.I&(1<<uint(y.I)) != 0}
}

var procBitShiftLeft = func(env *Env, args []Object) Object {
	x, y := AssertInts(env, args)
	return Int{I: x.I << uint(y.I)}
}

var procBitShiftRight = func(env *Env, args []Object) Object {
	x, y := AssertInts(env, args)
	return Int{I: x.I >> uint(y.I)}
}

var procUnsignedBitShiftRight = func(env *Env, args []Object) Object {
	x, y := AssertInts(env, args)
	return Int{I: int(uint(x.I) >> uint(y.I))}
}

var procExInfo = func(env *Env, args []Object) Object {
	CheckArity(env, args, 2, 3)
	res := &ExInfo{
		rt: env.RT.clone(),
	}
	res.Add(KEYWORDS.message, EnsureString(env, args, 0))
	res.Add(KEYWORDS.data, EnsureMap(env, args, 1))
	if len(args) == 3 {
		res.Add(KEYWORDS.cause, EnsureError(env, args, 2))
	}
	return res
}

var procExData = func(env *Env, args []Object) Object {
	if ok, res := args[0].(*ExInfo).Get(KEYWORDS.data); ok {
		return res
	}
	return NIL
}

var procExCause = func(env *Env, args []Object) Object {
	if ok, res := args[0].(*ExInfo).Get(KEYWORDS.cause); ok {
		return res
	}
	return NIL
}

var procExMessage = func(env *Env, args []Object) Object {
	return args[0].(Error).Message()
}

var procRegex = func(env *Env, args []Object) Object {
	r, err := regexp.Compile(EnsureString(env, args, 0).S)
	if err != nil {
		panic(env.RT.NewError("Invalid regex: " + err.Error()))
	}
	return &Regex{R: r}
}

func reGroups(s string, indexes []int) Object {
	if indexes == nil {
		return NIL
	} else if len(indexes) == 2 {
		if indexes[0] == -1 {
			return NIL
		} else {
			return String{S: s[indexes[0]:indexes[1]]}
		}
	} else {
		v := EmptyVector()
		for i := 0; i < len(indexes); i += 2 {
			if indexes[i] == -1 {
				v = v.Conjoin(NIL)
			} else {
				v = v.Conjoin(String{S: s[indexes[i]:indexes[i+1]]})
			}
		}
		return v
	}
}

var procReSeq = func(env *Env, args []Object) Object {
	re := EnsureRegex(env, args, 0)
	s := EnsureString(env, args, 1)
	matches := re.R.FindAllStringSubmatchIndex(s.S, -1)
	if matches == nil {
		return NIL
	}
	res := make([]Object, len(matches))
	for i, match := range matches {
		res[i] = reGroups(s.S, match)
	}
	return &ArraySeq{arr: res}
}

var procReFind = func(env *Env, args []Object) Object {
	re := EnsureRegex(env, args, 0)
	s := EnsureString(env, args, 1)
	match := re.R.FindStringSubmatchIndex(s.S)
	return reGroups(s.S, match)
}

var procRand = func(env *Env, args []Object) Object {
	r := rand.Float64()
	return Double{D: r}
}

var procIsSpecialSymbol = func(env *Env, args []Object) Object {
	return Boolean{B: IsSpecialSymbol(args[0])}
}

var procSubs = func(env *Env, args []Object) Object {
	s := EnsureString(env, args, 0).S
	start := EnsureInt(env, args, 1).I
	slen := utf8.RuneCountInString(s)
	end := slen
	if len(args) > 2 {
		end = EnsureInt(env, args, 2).I
	}
	if start < 0 || start > slen {
		panic(env.RT.NewError(fmt.Sprintf("String index out of range: %d", start)))
	}
	if end < 0 || end > slen {
		panic(env.RT.NewError(fmt.Sprintf("String index out of range: %d", end)))
	}
	return String{S: string([]rune(s)[start:end])}
}

var procIntern = func(env *Env, args []Object) Object {
	ns := EnsureNamespace(env, args, 0)
	sym := EnsureSymbol(env, args, 1)
	vr := ns.Intern(sym)
	if len(args) == 3 {
		vr.Value = args[2]
	}
	return vr
}

var procSetMeta = func(env *Env, args []Object) Object {
	vr := EnsureVar(env, args, 0)
	meta := EnsureMap(env, args, 1)
	vr.meta = meta
	return NIL
}

var procAtom = func(env *Env, args []Object) Object {
	res := &Atom{
		value: args[0],
	}
	if len(args) > 1 {
		m := NewHashMap(args[1:]...)
		if ok, v := m.Get(KEYWORDS.meta); ok {
			res.meta = AssertMap(env, v, "")
		}
	}
	return res
}

var procDeref = func(env *Env, args []Object) Object {
	return EnsureDeref(env, args, 0).Deref()
}

var procSwap = func(env *Env, args []Object) Object {
	a := EnsureAtom(env, args, 0)
	f := EnsureCallable(env, args, 1)
	fargs := append([]Object{a.value}, args[2:]...)
	a.value = f.Call(env, fargs)
	return a.value
}

var procSwapVals = func(env *Env, args []Object) Object {
	a := EnsureAtom(env, args, 0)
	f := EnsureCallable(env, args, 1)
	fargs := append([]Object{a.value}, args[2:]...)
	oldValue := a.value
	a.value = f.Call(env, fargs)
	return NewVectorFrom(oldValue, a.value)
}

var procReset = func(env *Env, args []Object) Object {
	a := EnsureAtom(env, args, 0)
	a.value = args[1]
	return a.value
}

var procResetVals = func(env *Env, args []Object) Object {
	a := EnsureAtom(env, args, 0)
	oldValue := a.value
	a.value = args[1]
	return NewVectorFrom(oldValue, a.value)
}

var procAlterMeta = func(env *Env, args []Object) Object {
	r := EnsureRef(env, args, 0)
	f := EnsureFn(env, args, 1)
	return r.AlterMeta(env, f, args[2:])
}

var procResetMeta = func(env *Env, args []Object) Object {
	r := EnsureRef(env, args, 0)
	m := EnsureMap(env, args, 1)
	return r.ResetMeta(m)
}

var procEmpty = func(env *Env, args []Object) Object {
	switch c := args[0].(type) {
	case Collection:
		return c.Empty()
	default:
		return NIL
	}
}

var procIsBound = func(env *Env, args []Object) Object {
	vr := EnsureVar(env, args, 0)
	return Boolean{B: vr.Value != nil}
}

func toNative(obj Object) interface{} {
	switch obj := obj.(type) {
	case Native:
		return obj.Native()
	default:
		return obj.ToString(false)
	}
}

var procFormat = func(env *Env, args []Object) Object {
	s := EnsureString(env, args, 0)
	objs := args[1:]
	fargs := make([]interface{}, len(objs))
	for i, v := range objs {
		fargs[i] = toNative(v)
	}
	res := fmt.Sprintf(s.S, fargs...)
	return String{S: res}
}

var procList = func(env *Env, args []Object) Object {
	return NewListFrom(args...)
}

var procCons = func(env *Env, args []Object) Object {
	CheckArity(env, args, 2, 2)
	s := EnsureSeqable(env, args, 1).Seq()
	return s.Cons(args[0])
}

var procFirst = func(env *Env, args []Object) Object {
	CheckArity(env, args, 1, 1)
	s := EnsureSeqable(env, args, 0).Seq()
	return s.First()
}

var procNext = func(env *Env, args []Object) Object {
	CheckArity(env, args, 1, 1)
	s := EnsureSeqable(env, args, 0).Seq()
	res := s.Rest()
	if res.IsEmpty() {
		return NIL
	}
	return res
}

var procRest = func(env *Env, args []Object) Object {
	CheckArity(env, args, 1, 1)
	s := EnsureSeqable(env, args, 0).Seq()
	return s.Rest()
}

var procConj = func(env *Env, args []Object) Object {
	switch c := args[0].(type) {
	case Conjable:
		return c.Conj(args[1])
	case Seq:
		return c.Cons(args[1])
	default:
		panic(env.RT.NewError("conj's first argument must be a collection, got " + c.GetType().ToString(false)))
	}
}

var procSeq = func(env *Env, args []Object) Object {
	CheckArity(env, args, 1, 1)
	s := EnsureSeqable(env, args, 0).Seq()
	if s.IsEmpty() {
		return NIL
	}
	return s
}

var procIsInstance = func(env *Env, args []Object) Object {
	CheckArity(env, args, 2, 2)
	t := EnsureType(env, args, 0)
	return Boolean{B: IsInstance(t, args[1])}
}

var procAssoc = func(env *Env, args []Object) Object {
	return EnsureAssociative(env, args, 0).Assoc(args[1], args[2])
}

var procEquals = func(env *Env, args []Object) Object {
	return Boolean{B: args[0].Equals(args[1])}
}

var procCount = func(env *Env, args []Object) Object {
	switch obj := args[0].(type) {
	case Counted:
		return Int{I: obj.Count()}
	default:
		s := AssertSeqable(env, obj, "count not supported on this type: "+obj.GetType().ToString(false))
		return Int{I: SeqCount(s.Seq())}
	}
}

var procSubvec = func(env *Env, args []Object) Object {
	// TODO: implement proper Subvector structure
	v := EnsureVector(env, args, 0)
	start := EnsureInt(env, args, 1).I
	end := EnsureInt(env, args, 2).I
	if start > end {
		panic(env.RT.NewError(fmt.Sprintf("subvec's start index (%d) is greater than end index (%d)", start, end)))
	}
	subv := make([]Object, 0, end-start)
	for i := start; i < end; i++ {
		subv = append(subv, v.at(i))
	}
	return NewVectorFrom(subv...)
}

var procCast = func(env *Env, args []Object) Object {
	t := EnsureType(env, args, 0)
	if t.reflectType.Kind() == reflect.Interface &&
		args[1].GetType().reflectType.Implements(t.reflectType) ||
		args[1].GetType().reflectType == t.reflectType {
		return args[1]
	}
	panic(env.RT.NewError("Cannot cast " + args[1].GetType().ToString(false) + " to " + t.ToString(false)))
}

var procVec = func(env *Env, args []Object) Object {
	return NewVectorFromSeq(EnsureSeqable(env, args, 0).Seq())
}

var procHashMap = func(env *Env, args []Object) Object {
	if len(args)%2 != 0 {
		panic(env.RT.NewError("No value supplied for key " + args[len(args)-1].ToString(false)))
	}
	return NewHashMap(args...)
}

var procHashSet = func(env *Env, args []Object) Object {
	res := EmptySet()
	for i := 0; i < len(args); i++ {
		res.Add(args[i])
	}
	return res
}

var procStr = func(env *Env, args []Object) Object {
	var buffer bytes.Buffer
	for _, obj := range args {
		if !obj.Equals(NIL) {
			t := obj.GetType()
			// TODO: this is a hack. Rethink escape parameter in ToString
			escaped := (t == TYPE.String) || (t == TYPE.Char) || (t == TYPE.Regex)
			buffer.WriteString(obj.ToString(!escaped))
		}
	}
	return String{S: buffer.String()}
}

var procSymbol = func(env *Env, args []Object) Object {
	if len(args) == 1 {
		return MakeSymbol(EnsureString(env, args, 0).S)
	}
	var ns *string = nil
	if !args[0].Equals(NIL) {
		ns = STRINGS.Intern(EnsureString(env, args, 0).S)
	}
	return Symbol{
		ns:   ns,
		name: STRINGS.Intern(EnsureString(env, args, 1).S),
	}
}

var procKeyword = func(env *Env, args []Object) Object {
	if len(args) == 1 {
		switch obj := args[0].(type) {
		case String:
			return MakeKeyword(obj.S)
		case Symbol:
			return Keyword{
				ns:   obj.ns,
				name: obj.name,
				hash: hashSymbol(obj.ns, obj.name) ^ KeywordHashMask,
			}
		default:
			return NIL
		}
	}
	var ns *string = nil
	if !args[0].Equals(NIL) {
		ns = STRINGS.Intern(EnsureString(env, args, 0).S)
	}
	name := STRINGS.Intern(EnsureString(env, args, 1).S)
	return Keyword{
		ns:   ns,
		name: name,
		hash: hashSymbol(ns, name) ^ KeywordHashMask,
	}
}

var procGensym = func(env *Env, args []Object) Object {
	return genSym(EnsureString(env, args, 0).S, "")
}

var procApply = func(env *Env, args []Object) Object {
	// TODO:
	// Stacktrace is broken. Need to somehow know
	// the name of the function passed ...
	f := EnsureCallable(env, args, 0)
	return f.Call(env, ToSlice(EnsureSeqable(env, args, 1).Seq()))
}

var procLazySeq = func(env *Env, args []Object) Object {
	return &LazySeq{
		fn: args[0].(*Fn),
	}
}

var procDelay = func(env *Env, args []Object) Object {
	return &Delay{
		fn: args[0].(*Fn),
	}
}

var procForce = func(env *Env, args []Object) Object {
	switch d := args[0].(type) {
	case *Delay:
		return d.Force(env)
	default:
		return d
	}
}

var procIdentical = func(env *Env, args []Object) Object {
	return Boolean{B: args[0] == args[1]}
}

var procCompare = func(env *Env, args []Object) Object {
	k1, k2 := args[0], args[1]
	if k1.Equals(k2) {
		return Int{I: 0}
	}
	switch k2.(type) {
	case Nil:
		return Int{I: 1}
	}
	switch k1 := k1.(type) {
	case Nil:
		return Int{I: -1}
	case Comparable:
		return Int{I: k1.Compare(k2)}
	}
	panic(env.RT.NewError(fmt.Sprintf("%s (type: %s) is not a Comparable", k1.ToString(true), k1.GetType().ToString(false))))
}

var procInt = func(env *Env, args []Object) Object {
	switch obj := args[0].(type) {
	case Char:
		return Int{I: int(obj.Ch)}
	case Number:
		return obj.Int()
	default:
		panic(env.RT.NewError(fmt.Sprintf("Cannot cast %s (type: %s) to Int", obj.ToString(true), obj.GetType().ToString(false))))
	}
}

var procNumber = func(env *Env, args []Object) Object {
	return AssertNumber(env, args[0], fmt.Sprintf("Cannot cast %s (type: %s) to Number", args[0].ToString(true), args[0].GetType().ToString(false)))
}

var procDouble = func(env *Env, args []Object) Object {
	n := AssertNumber(env, args[0], fmt.Sprintf("Cannot cast %s (type: %s) to Double", args[0].ToString(true), args[0].GetType().ToString(false)))
	return n.Double()
}

var procChar = func(env *Env, args []Object) Object {
	switch c := args[0].(type) {
	case Char:
		return c
	case Number:
		i := c.Int().I
		if i < MIN_RUNE || i > MAX_RUNE {
			panic(env.RT.NewError(fmt.Sprintf("Value out of range for char: %d", i)))
		}
		return Char{Ch: rune(i)}
	default:
		panic(env.RT.NewError(fmt.Sprintf("Cannot cast %s (type: %s) to Char", c.ToString(true), c.GetType().ToString(false))))
	}
}

var procBoolean = func(env *Env, args []Object) Object {
	return Boolean{B: ToBool(args[0])}
}

var procNumerator = func(env *Env, args []Object) Object {
	bi := EnsureRatio(env, args, 0).r.Num()
	return &BigInt{b: *bi}
}

var procDenominator = func(env *Env, args []Object) Object {
	bi := EnsureRatio(env, args, 0).r.Denom()
	return &BigInt{b: *bi}
}

var procBigInt = func(env *Env, args []Object) Object {
	switch n := args[0].(type) {
	case Number:
		return &BigInt{b: *n.BigInt()}
	case String:
		bi := big.Int{}
		if _, ok := bi.SetString(n.S, 10); ok {
			return &BigInt{b: bi}
		}
		panic(env.RT.NewError("Invalid number format " + n.S))
	default:
		panic(env.RT.NewError(fmt.Sprintf("Cannot cast %s (type: %s) to BigInt", n.ToString(true), n.GetType().ToString(false))))
	}
}

var procBigFloat = func(env *Env, args []Object) Object {
	switch n := args[0].(type) {
	case Number:
		return &BigFloat{b: *n.BigFloat()}
	case String:
		b := big.Float{}
		if _, ok := b.SetString(n.S); ok {
			return &BigFloat{b: b}
		}
		panic(env.RT.NewError("Invalid number format " + n.S))
	default:
		panic(env.RT.NewError(fmt.Sprintf("Cannot cast %s (type: %s) to BigFloat", n.ToString(true), n.GetType().ToString(false))))
	}
}

var procNth = func(env *Env, args []Object) Object {
	n := EnsureNumber(env, args, 1).Int().I
	switch coll := args[0].(type) {
	case Indexed:
		if len(args) == 3 {
			return coll.TryNth(n, args[2])
		}
		return coll.Nth(n)
	case Nil:
		return NIL
	case Sequential:
		switch coll := args[0].(type) {
		case Seqable:
			if len(args) == 3 {
				return SeqTryNth(coll.Seq(), n, args[2])
			}
			return SeqNth(coll.Seq(), n)
		}
	}
	panic(env.RT.NewError("nth not supported on this type: " + args[0].GetType().ToString(false)))
}

var procLt = func(env *Env, args []Object) Object {
	a := AssertNumber(env, args[0], "")
	b := AssertNumber(env, args[1], "")
	return Boolean{B: GetOps(a).Combine(GetOps(b)).Lt(a, b)}
}

var procLte = func(env *Env, args []Object) Object {
	a := AssertNumber(env, args[0], "")
	b := AssertNumber(env, args[1], "")
	return Boolean{B: GetOps(a).Combine(GetOps(b)).Lte(a, b)}
}

var procGt = func(env *Env, args []Object) Object {
	a := AssertNumber(env, args[0], "")
	b := AssertNumber(env, args[1], "")
	return Boolean{B: GetOps(a).Combine(GetOps(b)).Gt(a, b)}
}

var procGte = func(env *Env, args []Object) Object {
	a := AssertNumber(env, args[0], "")
	b := AssertNumber(env, args[1], "")
	return Boolean{B: GetOps(a).Combine(GetOps(b)).Gte(a, b)}
}

var procEq = func(env *Env, args []Object) Object {
	a := AssertNumber(env, args[0], "")
	b := AssertNumber(env, args[1], "")
	return MakeBoolean(numbersEq(a, b))
}

var procMax = func(env *Env, args []Object) Object {
	a := AssertNumber(env, args[0], "")
	b := AssertNumber(env, args[1], "")
	return Max(a, b)
}

var procMin = func(env *Env, args []Object) Object {
	a := AssertNumber(env, args[0], "")
	b := AssertNumber(env, args[1], "")
	return Min(a, b)
}

var procIncEx = func(env *Env, args []Object) Object {
	x := EnsureNumber(env, args, 0)
	ops := GetOps(x).Combine(BIGINT_OPS)
	return ops.Add(x, Int{I: 1})
}

var procDecEx = func(env *Env, args []Object) Object {
	x := EnsureNumber(env, args, 0)
	ops := GetOps(x).Combine(BIGINT_OPS)
	return ops.Subtract(x, Int{I: 1})
}

var procInc = func(env *Env, args []Object) Object {
	x := EnsureNumber(env, args, 0)
	ops := GetOps(x).Combine(INT_OPS)
	return ops.Add(x, Int{I: 1})
}

var procDec = func(env *Env, args []Object) Object {
	x := EnsureNumber(env, args, 0)
	ops := GetOps(x).Combine(INT_OPS)
	return ops.Subtract(x, Int{I: 1})
}

var procPeek = func(env *Env, args []Object) Object {
	s := AssertStack(env, args[0], "")
	return s.Peek()
}

var procPop = func(env *Env, args []Object) Object {
	s := AssertStack(env, args[0], "")
	return s.Pop().(Object)
}

var procContains = func(env *Env, args []Object) Object {
	switch c := args[0].(type) {
	case Gettable:
		ok, _ := c.Get(args[1])
		if ok {
			return Boolean{B: true}
		}
		return Boolean{B: false}
	}
	panic(env.RT.NewError("contains? not supported on type " + args[0].GetType().ToString(false)))
}

var procGet = func(env *Env, args []Object) Object {
	switch c := args[0].(type) {
	case Gettable:
		ok, v := c.Get(args[1])
		if ok {
			return v
		}
	}
	if len(args) == 3 {
		return args[2]
	}
	return NIL
}

var procDissoc = func(env *Env, args []Object) Object {
	return EnsureMap(env, args, 0).Without(args[1])
}

var procDisj = func(env *Env, args []Object) Object {
	return EnsureSet(env, args, 0).Disjoin(args[1])
}

var procFind = func(env *Env, args []Object) Object {
	res := EnsureAssociative(env, args, 0).EntryAt(args[1])
	if res == nil {
		return NIL
	}
	return res
}

var procKeys = func(env *Env, args []Object) Object {
	return EnsureMap(env, args, 0).Keys()
}

var procVals = func(env *Env, args []Object) Object {
	return EnsureMap(env, args, 0).Vals()
}

var procRseq = func(env *Env, args []Object) Object {
	return EnsureReversible(env, args, 0).Rseq()
}

var procName = func(env *Env, args []Object) Object {
	return String{S: EnsureNamed(env, args, 0).Name()}
}

var procNamespace = func(env *Env, args []Object) Object {
	ns := EnsureNamed(env, args, 0).Namespace()
	if ns == "" {
		return NIL
	}
	return String{S: ns}
}

var procFindVar = func(env *Env, args []Object) Object {
	sym := EnsureSymbol(env, args, 0)
	if sym.ns == nil {
		panic(env.RT.NewError("find-var argument must be namespace-qualified symbol"))
	}
	if v, ok := env.Resolve(sym); ok {
		return v
	}
	return NIL
}

var procSort = func(env *Env, args []Object) Object {
	cmp := EnsureComparator(env, args, 0)
	coll := EnsureSeqable(env, args, 1)
	s := SortableSlice{
		s:   ToSlice(coll.Seq()),
		cmp: cmp,
	}
	sort.Sort(s)
	return &ArraySeq{arr: s.s}
}

var procEval = func(env *Env, args []Object) Object {
	parseContext := &ParseContext{Env: env}
	expr := Parse(args[0], parseContext)
	return Eval(env, expr, nil)
}

var procType = func(env *Env, args []Object) Object {
	return args[0].GetType()
}

var procPprint = func(env *Env, args []Object) Object {
	obj := args[0]
	w := Assertio_Writer(env, env.stdout.Value, "")
	pprintObject(obj, 0, w)
	fmt.Fprint(w, "\n")
	return NIL
}

func PrintObject(env *Env, obj Object, w io.Writer) {
	printReadably := ToBool(env.printReadably.Value)
	switch obj := obj.(type) {
	case Printer:
		obj.Print(w, printReadably)
	default:
		fmt.Fprint(w, obj.ToString(printReadably))
	}
}

var procPr = func(env *Env, args []Object) Object {
	n := len(args)
	if n > 0 {
		f := Assertio_Writer(env, env.stdout.Value, "")
		for _, arg := range args[:n-1] {
			PrintObject(env, arg, f)
			fmt.Fprint(f, " ")
		}
		PrintObject(env, args[n-1], f)
	}
	return NIL
}

var procNewline = func(env *Env, args []Object) Object {
	f := Assertio_Writer(env, env.stdout.Value, "")
	fmt.Fprintln(f)
	return NIL
}

var procFlush = func(env *Env, args []Object) Object {
	switch f := args[0].(type) {
	case *File:
		f.Sync()
	}
	return NIL
}

func readFromReader(env *Env, reader io.RuneReader) Object {
	r := NewReader(reader, "<>")
	obj, err := TryRead(env, r)
	PanicOnErr(err)
	return obj
}

var procRead = func(env *Env, args []Object) Object {
	f := Ensureio_RuneReader(env, args, 0)
	return readFromReader(env, f)
}

var procReadString = func(env *Env, args []Object) Object {
	CheckArity(env, args, 1, 1)
	return readFromReader(env, strings.NewReader(EnsureString(env, args, 0).S))
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

var procReadLine = func(env *Env, args []Object) Object {
	CheckArity(env, args, 0, 0)
	f := AssertStringReader(env, env.stdin.Value, "")
	line, err := readLine(f)
	if err != nil {
		return NIL
	}
	return String{S: line}
}

var procReaderReadLine = func(env *Env, args []Object) Object {
	CheckArity(env, args, 1, 1)
	rdr := EnsureStringReader(env, args, 0)
	line, err := readLine(rdr)
	if err != nil {
		return NIL
	}
	return String{S: line}
}

var procNanoTime = func(env *Env, args []Object) Object {
	return &BigInt{b: *big.NewInt(time.Now().UnixNano())}
}

var procMacroexpand1 = func(env *Env, args []Object) Object {
	switch s := args[0].(type) {
	case Seq:
		parseContext := &ParseContext{Env: env}
		return macroexpand1(env, s, parseContext)
	default:
		return s
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

var procLoadString = func(env *Env, args []Object) Object {
	s := EnsureString(env, args, 0)
	obj, err := loadReader(env, NewReader(strings.NewReader(s.S), "<string>"))
	if err != nil {
		panic(err)
	}
	return obj
}

var procFindNamespace = func(env *Env, args []Object) Object {
	ns := env.FindNamespace(EnsureSymbol(env, args, 0))
	if ns == nil {
		return NIL
	}
	return ns
}

var procCreateNamespace = func(env *Env, args []Object) Object {
	sym := EnsureSymbol(env, args, 0)
	res := env.EnsureNamespace(sym)
	// In linter mode the latest create-ns call overrides position info.
	// This is for the cases when (ns ...) is called in .jokerd/linter.clj file and alike.
	// Also, isUsed needs to be reset in this case.
	if LINTER_MODE {
		res.Name = res.Name.WithInfo(sym.GetInfo()).(Symbol)
		res.isUsed = false
	}
	return res
}

var procInjectNamespace = func(env *Env, args []Object) Object {
	sym := EnsureSymbol(env, args, 0)
	ns := env.EnsureNamespace(sym)
	ns.isUsed = true
	ns.isGloballyUsed = true
	return ns
}

var procRemoveNamespace = func(env *Env, args []Object) Object {
	ns := env.RemoveNamespace(EnsureSymbol(env, args, 0))
	if ns == nil {
		return NIL
	}
	return ns
}

var procAllNamespaces = func(env *Env, args []Object) Object {
	s := make([]Object, 0, len(env.Namespaces))
	for _, ns := range env.Namespaces {
		s = append(s, ns)
	}
	return &ArraySeq{arr: s}
}

var procNamespaceName = func(env *Env, args []Object) Object {
	return EnsureNamespace(env, args, 0).Name
}

var procNamespaceMap = func(env *Env, args []Object) Object {
	r := &ArrayMap{}
	for k, v := range EnsureNamespace(env, args, 0).mappings {
		r.Add(MakeSymbol(*k), v)
	}
	return r
}

var procNamespaceUnmap = func(env *Env, args []Object) Object {
	ns := EnsureNamespace(env, args, 0)
	sym := EnsureSymbol(env, args, 1)
	if sym.ns != nil {
		panic(env.RT.NewError("Can't unintern namespace-qualified symbol"))
	}
	delete(ns.mappings, sym.name)
	return NIL
}

var procVarNamespace = func(env *Env, args []Object) Object {
	v := EnsureVar(env, args, 0)
	return v.ns
}

var procRefer = func(env *Env, args []Object) Object {
	ns := EnsureNamespace(env, args, 0)
	sym := EnsureSymbol(env, args, 1)
	v := EnsureVar(env, args, 2)
	return ns.Refer(sym, v)
}

var procAlias = func(env *Env, args []Object) Object {
	EnsureNamespace(env, args, 0).AddAlias(EnsureSymbol(env, args, 1), EnsureNamespace(env, args, 2))
	return NIL
}

var procNamespaceAliases = func(env *Env, args []Object) Object {
	r := &ArrayMap{}
	for k, v := range EnsureNamespace(env, args, 0).aliases {
		r.Add(MakeSymbol(*k), v)
	}
	return r
}

var procNamespaceUnalias = func(env *Env, args []Object) Object {
	ns := EnsureNamespace(env, args, 0)
	sym := EnsureSymbol(env, args, 1)
	if sym.ns != nil {
		panic(env.RT.NewError("Alias can't be namespace-qualified"))
	}
	delete(ns.aliases, sym.name)
	return NIL
}

var procVarGet = func(env *Env, args []Object) Object {
	return EnsureVar(env, args, 0).Resolve()
}

var procVarSet = func(env *Env, args []Object) Object {
	EnsureVar(env, args, 0).Value = args[1]
	return args[1]
}

var procNsResolve = func(env *Env, args []Object) Object {
	ns := EnsureNamespace(env, args, 0)
	sym := EnsureSymbol(env, args, 1)
	if sym.ns == nil && TYPES[sym.name] != nil {
		return TYPES[sym.name]
	}
	if vr, ok := env.ResolveIn(ns, sym); ok {
		return vr
	}
	return NIL
}

var procArrayMap = func(env *Env, args []Object) Object {
	if len(args)%2 == 1 {
		panic(env.RT.NewError("No value supplied for key " + args[len(args)-1].ToString(false)))
	}
	res := EmptyArrayMap()
	for i := 0; i < len(args); i += 2 {
		res.Set(args[i], args[i+1])
	}
	return res
}

const bufferHashMask uint32 = 0x5ed19e84

var procBuffer = func(env *Env, args []Object) Object {
	if len(args) > 0 {
		s := EnsureString(env, args, 0)
		return MakeBuffer(bytes.NewBufferString(s.S))
	}
	return MakeBuffer(&bytes.Buffer{})
}

var procBufferedReader = func(env *Env, args []Object) Object {
	switch rdr := args[0].(type) {
	case io.Reader:
		return MakeBufferedReader(rdr)
	default:
		panic(env.RT.NewArgTypeError(0, args[0], "IOReader"))
	}
}

var procSlurp = func(env *Env, args []Object) Object {
	b, err := os.ReadFile(EnsureString(env, args, 0).S)
	PanicOnErr(err)
	return String{S: string(b)}
}

var procSpit = func(env *Env, args []Object) Object {
	filename := EnsureString(env, args, 0)
	content := EnsureString(env, args, 1)
	opts := EnsureMap(env, args, 2)
	appendFile := false
	if ok, append := opts.Get(MakeKeyword("append")); ok {
		appendFile = ToBool(append)
	}
	flags := os.O_CREATE | os.O_WRONLY
	if appendFile {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	f, err := os.OpenFile(filename.S, flags, 0644)
	PanicOnErr(err)
	defer f.Close()
	_, err = f.WriteString(content.S)
	PanicOnErr(err)
	return NIL
}

var procShuffle = func(env *Env, args []Object) Object {
	s := ToSlice(EnsureSeqable(env, args, 0).Seq())
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
	return NewVectorFrom(s...)
}

var procIsRealized = func(env *Env, args []Object) Object {
	return Boolean{B: EnsurePending(env, args, 0).IsRealized()}
}

var procDeriveInfo = func(env *Env, args []Object) Object {
	dest := args[0]
	src := args[1]
	return dest.WithInfo(src.GetInfo())
}

var procJokerVersion = func(env *Env, args []Object) Object {
	return String{S: VERSION[1:]}
}

var procHash = func(env *Env, args []Object) Object {
	return Int{I: int(args[0].Hash())}
}

func loadFile(env *Env, filename string) Object {
	var reader *Reader
	f, err := os.Open(filename)
	PanicOnErr(err)
	reader = NewReader(bufio.NewReader(f), filename)
	ProcessReaderFromEval(env, reader, filename)
	return NIL
}

var procLoadFile = func(env *Env, args []Object) Object {
	filename := EnsureString(env, args, 0)
	return loadFile(env, filename.S)
}

var procLoadLibFromPath = func(env *Env, args []Object) Object {
	libname := EnsureSymbol(env, args, 0).Name()
	pathname := EnsureString(env, args, 1).S
	cp := env.classPath.Value
	cpvec := AssertVector(env, cp, "*classpath* must be a Vector, not a "+cp.GetType().ToString(false))
	count := cpvec.Count()
	var f *os.File
	var err error
	var canonicalErr error
	var filename string
	for i := 0; i < count; i++ {
		elem := cpvec.at(i)
		cpelem := AssertString(env, elem, "*classpath* must contain only Strings, not a "+elem.GetType().ToString(false)+" (at element "+strconv.Itoa(i)+")")
		s := cpelem.S
		if s == "" {
			filename = pathname
		} else {
			filename = filepath.Join(s, filepath.Join(strings.Split(libname, ".")...)) + ".joke" // could cache inner join....
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
	PanicOnErr(canonicalErr)
	PanicOnErr(err)
	reader := NewReader(bufio.NewReader(f), filename)
	ProcessReaderFromEval(env, reader, filename)
	return NIL
}

var procReduceKv = func(env *Env, args []Object) Object {
	f := EnsureCallable(env, args, 0)
	init := args[1]
	coll := EnsureKVReduce(env, args, 2)
	return coll.kvreduce(f, init)
}

var procIndexOf = func(env *Env, args []Object) Object {
	s := EnsureString(env, args, 0)
	ch := EnsureChar(env, args, 1)
	for i, r := range s.S {
		if r == ch.Ch {
			return Int{I: i}
		}
	}
	return Int{I: -1}
}

func libExternalPath(env *Env, sym Symbol) (path string, ok bool) {
	nsSourcesVar, _ := env.Resolve(MakeSymbol("joker.core/*ns-sources*"))
	nsSources := ToSlice(nsSourcesVar.Value.(*Vector).Seq())

	var sourceKey string
	var sourceMap Map
	for _, source := range nsSources {
		sourceKey = source.(*Vector).Nth(0).ToString(false)
		match, _ := regexp.MatchString(sourceKey, sym.Name())
		if match {
			sourceMap = source.(*Vector).Nth(1).(Map)
			break
		}
	}
	if sourceMap != nil {
		ok, url := sourceMap.Get(MakeKeyword("url"))
		if !ok {
			panic(env.RT.NewError("Key :url not found in ns-sources for: " + sourceKey))
		} else {
			return externalSourceToPath(env, sym.Name(), url.ToString(false)), true
		}
	}
	return
}

var procLibPath = func(env *Env, args []Object) Object {
	sym := EnsureSymbol(env, args, 0)
	var path string

	path, ok := libExternalPath(env, sym)

	if !ok {
		var file string
		if env.file.Value == nil {
			var err error
			file, err = filepath.Abs("user")
			PanicOnErr(err)
		} else {
			file = AssertString(env, env.file.Value, "").S
			if linkDest, err := os.Readlink(file); err == nil {
				file = linkDest
			}
		}
		ns := env.CurrentNamespace().Name

		parts := strings.Split(ns.Name(), ".")
		for _ = range parts {
			file, _ = filepath.Split(file)
			file = file[:len(file)-1]
		}
		path = filepath.Join(append([]string{file}, strings.Split(sym.Name(), ".")...)...) + ".joke"
	}
	return String{S: path}
}

var procInternFakeVar = func(env *Env, args []Object) Object {
	nsSym := EnsureSymbol(env, args, 0)
	sym := EnsureSymbol(env, args, 1)
	isMacro := ToBool(args[2])
	res := InternFakeSymbol(env, env.FindNamespace(nsSym), sym)
	res.isMacro = isMacro
	return res
}

var procParse = func(env *Env, args []Object) Object {
	lm, _ := env.Resolve(MakeSymbol("joker.core/*linter-mode*"))
	lm.Value = Boolean{B: true}
	LINTER_MODE = true
	defer func() {
		LINTER_MODE = false
		lm.Value = Boolean{B: false}
	}()
	parseContext := &ParseContext{Env: env}
	res := Parse(args[0], parseContext)
	return res.Dump(false)
}

var procTypes = func(env *Env, args []Object) Object {
	CheckArity(env, args, 0, 0)
	res := EmptyArrayMap()
	for k, v := range TYPES {
		res.Add(String{S: *k}, v)
	}
	return res
}

var procCreateChan = func(env *Env, args []Object) Object {
	CheckArity(env, args, 1, 1)
	n := EnsureInt(env, args, 0)
	ch := make(chan FutureResult, n.I)
	return MakeChannel(ch)
}

var procCloseChan = func(env *Env, args []Object) Object {
	CheckArity(env, args, 1, 1)
	EnsureChannel(env, args, 0).Close()
	return NIL
}

var procSend = func(env *Env, args []Object) (obj Object) {
	CheckArity(env, args, 2, 2)
	ch := EnsureChannel(env, args, 0)
	v := args[1]
	if v.Equals(NIL) {
		panic(env.RT.NewError("Can't put nil on channel"))
	}
	if ch.isClosed {
		return MakeBoolean(false)
	}
	obj = MakeBoolean(true)
	defer func() {
		if r := recover(); r != nil {
			//env.RT.GIL.Lock()
			obj = MakeBoolean(false)
		}
	}()
	//env.RT.GIL.Unlock()
	ch.ch <- MakeFutureResult(v, nil)
	//env.RT.GIL.Lock()
	return
}

var procReceive = func(env *Env, args []Object) Object {
	CheckArity(env, args, 1, 1)
	ch := EnsureChannel(env, args, 0)
	//env.RT.GIL.Unlock()
	res, ok := <-ch.ch
	//env.RT.GIL.Lock()
	if !ok {
		return NIL
	}
	if res.err != nil {
		panic(res.err)
	}
	return res.value
}

var procGo = func(env *Env, args []Object) Object {
	CheckArity(env, args, 1, 1)
	f := EnsureCallable(env, args, 0)
	ch := MakeChannel(make(chan FutureResult, 1))
	go func() {

		defer func() {
			if r := recover(); r != nil {
				switch r := r.(type) {
				case Error:
					ch.ch <- MakeFutureResult(NIL, r)
					ch.Close()
				default:
					//env.RT.GIL.Unlock()
					panic(r)
				}
			}
			//env.RT.GIL.Unlock()
		}()

		//env.RT.GIL.Lock()
		res := f.Call(env, []Object{})
		ch.ch <- MakeFutureResult(res, nil)
		ch.Close()
	}()
	return ch
}

var procVerbosityLevel = func(env *Env, args []Object) Object {
	CheckArity(env, args, 0, 0)
	return MakeInt(VerbosityLevel)
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
		PanicOnErr(err)
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

var procIncProblemCount = func(env *Env, args []Object) Object {
	PROBLEM_COUNT++
	return NIL
}

func ProcessReader(env *Env, reader *Reader, filename string, phase Phase) error {
	parseContext := &ParseContext{Env: env}
	if filename != "" {
		currentFilename := parseContext.Env.file.Value
		defer func() {
			parseContext.Env.SetFilename(currentFilename)
		}()
		s, err := filepath.Abs(filename)
		PanicOnErr(err)
		parseContext.Env.SetFilename(MakeString(s))
	}
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
		obj, err = TryEval(env, expr)
		if err != nil {
			fmt.Fprintln(Stderr, err)
			return err
		}
		if phase == EVAL {
			continue
		}
		if _, ok := obj.(Nil); !ok {
			fmt.Fprintln(Stdout, obj.ToString(true))
		}
	}
}

func ProcessReaderFromEval(env *Env, reader *Reader, filename string) {
	parseContext := &ParseContext{Env: env}
	if filename != "" {
		currentFilename := parseContext.Env.file.Value
		defer func() {
			parseContext.Env.SetFilename(currentFilename)
		}()
		s, err := filepath.Abs(filename)
		PanicOnErr(err)
		parseContext.Env.SetFilename(MakeString(s))
	}
	for {
		obj, err := TryRead(env, reader)
		if err == io.EOF {
			return
		}
		PanicOnErr(err)
		expr, err := TryParse(obj, parseContext)
		PanicOnErr(err)
		obj, err = TryEval(env, expr)
		PanicOnErr(err)
	}
}

func processData(data []byte) {
	processInEnv(GLOBAL_ENV, data)
}

func processInEnv(env *Env, data []byte) error {
	ns := env.CurrentNamespace()
	env.SetCurrentNamespace(env.CoreNamespace)
	defer func() { env.SetCurrentNamespace(ns) }()
	header, p := UnpackHeader(data, env)
	for len(p) > 0 {
		var expr Expr
		expr, p = UnpackExpr(env, p, header)
		_, err := TryEval(env, expr)
		PanicOnErr(err)
	}
	if VerbosityLevel > 0 {
		fmt.Fprintf(Stderr, "processData: Evaluated code for %s\n", env.CurrentNamespace().ToString(false))
	}

	return nil
}

func setCoreNamespaces(env *Env) {
	ns := env.CoreNamespace
	ns.MaybeLazy("joker.core")

	vr := ns.Resolve("*core-namespaces*")
	set := vr.Value.(*MapSet)
	for _, ns := range coreNamespaces {
		set = set.Conj(MakeSymbol(ns)).(*MapSet)
	}
	vr.Value = set

	// Add 'joker.core to *loaded-libs*, now that it's loaded.
	vr = ns.Resolve("*loaded-libs*")
	set = vr.Value.(*MapSet).Conj(ns.Name).(*MapSet)
	vr.Value = set
}

var procIsNamespaceInitialized = func(env *Env, args []Object) Object {
	sym := EnsureSymbol(env, args, 0)
	if sym.ns != nil {
		panic(env.RT.NewError("Can't ask for namespace info on namespace-qualified symbol"))
	}
	// First look for registered (e.g. std) libs
	ns, found := env.Namespaces[sym.name]
	return MakeBoolean(found && ns.Lazy == nil)
}

func findConfigFile(filename string, workingDir string, findDir bool) string {
	var err error
	configName := ".joker"
	if findDir {
		configName = ".jokerd"
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

func knownMacrosToMap(km Object) (Map, error) {
	s := km.(Seqable).Seq()
	res := EmptyArrayMap()
	for !s.IsEmpty() {
		obj := s.First()
		switch obj := obj.(type) {
		case Symbol:
			res.Add(obj, NIL)
		case *Vector:
			if obj.Count() != 2 {
				return nil, errors.New(":known-macros item must be a symbol or a vector with two elements")
			}
			res.Add(obj.at(0), obj.at(1))
		default:
			return nil, errors.New(":known-macros item must be a symbol or a vector, got " + obj.GetType().ToString(false))
		}
		s = s.Rest()
	}
	return res, nil
}

func ReadConfig(env *Env, filename string, workingDir string) {
	LINTER_CONFIG = env.CoreNamespace.Intern(MakeSymbol("*linter-config*"))
	LINTER_CONFIG.Value = EmptyArrayMap()
	configFileName := findConfigFile(filename, workingDir, false)
	if configFileName == "" {
		return
	}
	f, err := os.Open(configFileName)
	if err != nil {
		printConfigError(configFileName, err.Error())
		return
	}
	r := NewReader(bufio.NewReader(f), configFileName)
	config, err := TryRead(env, r)
	if err != nil {
		printConfigError(configFileName, err.Error())
		return
	}
	configMap, ok := config.(Map)
	if !ok {
		printConfigError(configFileName, "config root object must be a map, got "+config.GetType().ToString(false))
		return
	}
	ok, ignoredUnusedNamespaces := configMap.Get(MakeKeyword("ignored-unused-namespaces"))
	if ok {
		seq, ok1 := ignoredUnusedNamespaces.(Seqable)
		if ok1 {
			WARNINGS.ignoredUnusedNamespaces = NewSetFromSeq(seq.Seq())
		} else {
			printConfigError(configFileName, ":ignored-unused-namespaces value must be a vector, got "+ignoredUnusedNamespaces.GetType().ToString(false))
			return
		}
	}
	ok, ignoredFileRegexes := configMap.Get(MakeKeyword("ignored-file-regexes"))
	if ok {
		seq, ok1 := ignoredFileRegexes.(Seqable)
		if ok1 {
			s := seq.Seq()
			for !s.IsEmpty() {
				regex, ok2 := s.First().(*Regex)
				if !ok2 {
					printConfigError(configFileName, ":ignored-file-regexes elements must be regexes, got "+s.First().GetType().ToString(false))
					return
				}
				WARNINGS.IgnoredFileRegexes = append(WARNINGS.IgnoredFileRegexes, regex.R)
				s = s.Rest()
			}
		} else {
			printConfigError(configFileName, ":ignored-file-regexes value must be a vector, got "+ignoredFileRegexes.GetType().ToString(false))
			return
		}
	}
	ok, entryPoints := configMap.Get(MakeKeyword("entry-points"))
	if ok {
		seq, ok1 := entryPoints.(Seqable)
		if ok1 {
			WARNINGS.entryPoints = NewSetFromSeq(seq.Seq())
		} else {
			printConfigError(configFileName, ":entry-points value must be a vector, got "+entryPoints.GetType().ToString(false))
			return
		}
	}
	ok, knownNamespaces := configMap.Get(MakeKeyword("known-namespaces"))
	if ok {
		if _, ok1 := knownNamespaces.(Seqable); !ok1 {
			printConfigError(configFileName, ":known-namespaces value must be a vector, got "+knownNamespaces.GetType().ToString(false))
			return
		}
	}
	ok, knownTags := configMap.Get(MakeKeyword("known-tags"))
	if ok {
		if _, ok1 := knownTags.(Seqable); !ok1 {
			printConfigError(configFileName, ":known-tags value must be a vector, got "+knownTags.GetType().ToString(false))
			return
		}
	}
	ok, knownMacros := configMap.Get(KEYWORDS.knownMacros)
	if ok {
		_, ok1 := knownMacros.(Seqable)
		if !ok1 {
			printConfigError(configFileName, ":known-macros value must be a vector, got "+knownMacros.GetType().ToString(false))
			return
		}
		m, err := knownMacrosToMap(knownMacros)
		if err != nil {
			printConfigError(configFileName, err.Error())
			return
		}
		configMap = configMap.Assoc(KEYWORDS.knownMacros, m).(Map)
	}
	ok, rules := configMap.Get(KEYWORDS.rules)
	if ok {
		m, ok := rules.(Map)
		if !ok {
			printConfigError(configFileName, ":rules value must be a map, got "+rules.GetType().ToString(false))
			return
		}
		if ok, v := m.Get(KEYWORDS.ifWithoutElse); ok {
			WARNINGS.ifWithoutElse = ToBool(v)
		}
		if ok, v := m.Get(KEYWORDS.unusedFnParameters); ok {
			WARNINGS.unusedFnParameters = ToBool(v)
		}
		if ok, v := m.Get(KEYWORDS.fnWithEmptyBody); ok {
			WARNINGS.fnWithEmptyBody = ToBool(v)
		}
	}
	LINTER_CONFIG.Value = configMap
}

func removeJokerNamespaces(env *Env) {
	for k, ns := range env.Namespaces {
		if ns != env.CoreNamespace && strings.HasPrefix(*k, "joker.") {
			delete(env.Namespaces, k)
		}
	}
}

func markJokerNamespacesAsUsed(env *Env) {
	for k, ns := range env.Namespaces {
		if ns != env.CoreNamespace && strings.HasPrefix(*k, "joker.") {
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

func ProcessLinterFile(env *Env, configDir string, filename string) {
	linterFileName := filepath.Join(configDir, filename)
	if _, err := os.Stat(linterFileName); err == nil {
		if reader, err := NewReaderFromFile(linterFileName); err == nil {
			ProcessReader(env, reader, linterFileName, EVAL)
		}
	}
}

func ProcessLinterFiles(env *Env, dialect Dialect, filename string, workingDir string) {
	if dialect == EDN || dialect == JOKER {
		return
	}
	configDir := findConfigFile(filename, workingDir, true)
	if configDir == "" {
		return
	}
	ProcessLinterFile(env, configDir, "linter.cljc")
	switch dialect {
	case CLJS:
		ProcessLinterFile(env, configDir, "linter.cljs")
	case CLJ:
		ProcessLinterFile(env, configDir, "linter.clj")
	}
}
