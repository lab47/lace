package core

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
)

type ConvRegistry struct {
	converters   map[reflect.Type]ProcFn
	cacheCS      map[reflect.Type]*conversionSet
	cacheWrapper map[reflect.Value]cachedProc
}

type cachedProc struct {
	fn ProcFn
	cs *conversionSet
}

var convReg = &ConvRegistry{
	converters:   make(map[reflect.Type]ProcFn),
	cacheCS:      make(map[reflect.Type]*conversionSet),
	cacheWrapper: make(map[reflect.Value]cachedProc),
}

func (c *ConvRegistry) ConverterForFunc(v reflect.Value) (ProcFn, *conversionSet, error) {
	if f, ok := c.cacheWrapper[v]; ok {
		return f.fn, f.cs, nil
	}

	f, cs, err := c.buildProc(v)
	if err != nil {
		return nil, nil, err
	}

	c.cacheWrapper[v] = cachedProc{fn: f, cs: cs}

	return f, cs, nil
}

type Args struct {
	Objects []Object
}

func (c *ConvRegistry) convArg(at reflect.Type) (inConv, bool) {
	switch at {
	case reflect.TypeFor[*Env]():
		return convertEnvIn, false
	case reflect.TypeFor[Object]():
		return convertObjectIn, true
	case reflect.TypeFor[Seqable]():
		return convertSeqableIn, true
	case reflect.TypeFor[Map]():
		return convertMapIn, true
	case reflect.TypeFor[Callable]():
		return convertCallableIn, true
	case reflect.TypeFor[string]():
		return convertStringIn, true
	case reflect.TypeFor[Symbol]():
		return convertSymbolIn, true
	case reflect.TypeFor[Keyword]():
		return convertKeywordIn, true
	case reflect.TypeFor[[]byte]():
		return convertBytesIn, true
	case reflect.TypeFor[int]():
		return convertIntIn, true
	case reflect.TypeFor[int8]():
		return convertIntIn8, true
	case reflect.TypeFor[int16]():
		return convertIntIn16, true
	case reflect.TypeFor[int32]():
		return convertIntIn32, true
	case reflect.TypeFor[int64]():
		return convertIntIn64, true
	case reflect.TypeFor[uint]():
		return convertUintIn, true
	case reflect.TypeFor[uint8]():
		return convertUintIn8, true
	case reflect.TypeFor[uint16]():
		return convertUintIn16, true
	case reflect.TypeFor[uint32]():
		return convertUintIn32, true
	case reflect.TypeFor[uint64]():
		return convertUintIn64, true
	case reflect.TypeFor[reflect.Type]():
		return convertReflectTypeIn, true
	case reflect.TypeFor[reflect.Value]():
		return convertReflectValueInAny, true
	case reflect.TypeFor[*ReflectValue]():
		return convertReflectValueInDirect, true
	case reflect.TypeFor[any]():
		return convertAnyIn, true
	default:
		return func(e *Env, i int, o Object) (reflect.Value, error) {
			return convertReflectValueIn(e, i, o, at)
		}, true
	}
}

func (c *ConvRegistry) convRet(at reflect.Type) outConv {
	switch at {
	case reflect.TypeFor[Object]():
		return convertObjectOut
	case reflect.TypeFor[string]():
		return convertStringOut
	case reflect.TypeFor[[]byte]():
		return convertBytesOut
	case reflect.TypeFor[bool]():
		return convertBoolOut
	case reflect.TypeFor[int](),
		reflect.TypeFor[int64](),
		reflect.TypeFor[int32](),
		reflect.TypeFor[int16](),
		reflect.TypeFor[int8]():

		return convertFromInt
	case reflect.TypeFor[uint](),
		reflect.TypeFor[uint64](),
		reflect.TypeFor[uint32](),
		reflect.TypeFor[uint16](),
		reflect.TypeFor[uint8]():

		return convertFromUInt
	case reflect.TypeFor[reflect.Type]():
		return convertGoReflectTypeOut
	case reflect.TypeFor[*regexp.Regexp]():
		return convertRegexpOut
	default:
		return convertReflectValueOut
	}
}

var nilObject = reflect.Zero(reflect.TypeFor[Object]())

type inConv func(*Env, int, Object) (reflect.Value, error)
type outConv func(*Env, reflect.Value) (Object, error)

type conversionSet struct {
	ft reflect.Type

	argIn    []inConv
	arityMap []int
	arity    int

	rets     []outConv
	values   int
	errorPos int
}

func (c *ConvRegistry) wrapFunc(fnVal reflect.Value, cs *conversionSet) reflect.Value {
	return reflect.ValueOf(ProcFn(func(env *Env, objArgs []Object) (Object, error) {
		if len(objArgs) != cs.arity {
			return nil, ErrorArityMinMax(env, len(objArgs), cs.arity, cs.arity)
		}

		var ret []reflect.Value

		// Optimization for a normal (ie, fewer than 10) number of args to avoid heap escape of the
		// input to .Call
		if cs.ft.NumIn() <= 10 {
			var dest [10]reflect.Value

			for destIdx, inIdx := range cs.arityMap {
				var a Object

				if inIdx >= 0 {
					a = objArgs[inIdx]
				}
				dest[destIdx] = nilObject

				sub, err := cs.argIn[destIdx](env, inIdx, a)
				if err != nil {
					return nil, err
				}

				dest[destIdx] = sub
			}

			ret = fnVal.Call(dest[:cs.ft.NumIn()])
		} else {
			dest := make([]reflect.Value, cs.ft.NumIn())

			for destIdx, inIdx := range cs.arityMap {
				var a Object

				if inIdx >= 0 {
					a = objArgs[inIdx]
				}
				dest[destIdx] = nilObject

				sub, err := cs.argIn[destIdx](env, inIdx, a)
				if err != nil {
					return nil, err
				}

				dest[destIdx] = sub
			}

			ret = fnVal.Call(dest[:cs.ft.NumIn()])
		}

		if cs.errorPos >= 0 && !ret[cs.errorPos].IsNil() {
			return nil, ret[cs.errorPos].Interface().(error)
		}

		switch cs.values {
		case 0:
			return NIL, nil
		case 1:
			return cs.rets[0](env, ret[0])
		default:
			var objects []Object

			for i, rv := range ret {
				if i == cs.errorPos {
					continue
				}

				v, err := cs.rets[i](env, rv)
				if err != nil {
					objects = append(objects, v)
				}
			}

			return NewListFrom(objects...), nil
		}
	}))
}

func (c *ConvRegistry) buildCS(t reflect.Type) *conversionSet {
	if cs, ok := c.cacheCS[t]; ok {
		return cs
	}

	var argIn []inConv

	var arityMap []int

	var arity int
	for i := 0; i < t.NumIn(); i++ {
		at := t.In(i)

		conv, inputArg := c.convArg(at)
		if inputArg {
			arityMap = append(arityMap, arity)
			arity++
		} else {
			arityMap = append(arityMap, -1)
		}

		argIn = append(argIn, conv)
	}

	rets := make([]outConv, t.NumOut())

	var valueReturns int

	errorPos := -1

	for i := 0; i < t.NumOut(); i++ {
		at := t.Out(i)

		switch at {
		case reflect.TypeFor[error]():
			errorPos = i
		default:
			rets[i] = c.convRet(at)
			valueReturns++
		}
	}

	cs := &conversionSet{
		ft:       t,
		argIn:    argIn,
		arity:    arity,
		arityMap: arityMap,

		rets:     rets,
		errorPos: errorPos,
		values:   valueReturns,
	}

	c.cacheCS[t] = cs

	return cs
}

func (c *ConvRegistry) buildProc(v reflect.Value) (ProcFn, *conversionSet, error) {
	if v.Kind() != reflect.Func {
		return nil, nil, fmt.Errorf("procs can only be built from Go funcs, is: %s (%s)", v, v.Kind())
	}

	vt := v.Type()

	if vt == rawFunc {
		vc := v.Interface().(func(*Env, Args) (Object, error))

		return func(env *Env, args []Object) (Object, error) {
			return vc(env, Args{Objects: args})
		}, nil, nil
	}

	cs := c.buildCS(v.Type())

	wrapper := c.wrapFunc(v, cs)

	return wrapper.Interface().(ProcFn), cs, nil
}

// from string to String
func convertStringOut(env *Env, s reflect.Value) (Object, error) {
	gs, ok := s.Interface().(string)
	if !ok {
		return nil, env.RT.NewError("wrong type, expecting string")
	}

	return MakeString(gs), nil
}

// from []byte to String
func convertBytesOut(env *Env, s reflect.Value) (Object, error) {
	_, ok := s.Interface().([]byte)
	if !ok {
		return nil, env.RT.NewError("wrong type, expecting string")
	}

	return &ReflectValue{val: s}, nil
}

// from *regexp.Regexp to *Regex
func convertRegexpOut(env *Env, s reflect.Value) (Object, error) {
	gs, ok := s.Interface().(*regexp.Regexp)
	if !ok {
		return nil, env.RT.NewError("wrong type, expecting string")
	}

	return MakeRegex(gs), nil
}

// from Object to Object
func convertObjectOut(env *Env, s reflect.Value) (Object, error) {
	gs, ok := s.Interface().(Object)
	if !ok {
		return nil, env.RT.NewError("wrong type, expecting Object")
	}

	return gs, nil
}

// from bool to Boolean
func convertBoolOut(env *Env, s reflect.Value) (Object, error) {
	gb, ok := s.Interface().(bool)
	if !ok {
		return nil, env.RT.NewError("wrong type, expecting string")
	}

	return MakeBoolean(gb), nil
}

// from reflect.Type to ReflectType
func convertGoReflectTypeOut(env *Env, s reflect.Value) (Object, error) {
	gb, ok := s.Interface().(reflect.Type)
	if !ok {
		return nil, env.RT.NewError("wrong type, expecting reflect.Type")
	}

	return &ReflectType{typ: gb}, nil
}

// from String to string
func convertStringIn(env *Env, index int, o Object) (reflect.Value, error) {
	switch sv := o.(type) {
	case String:
		return reflect.ValueOf(sv.S), nil
	case Symbol:
		return reflect.ValueOf(sv.Name()), nil
	case Keyword:
		return reflect.ValueOf(sv.Name()), nil
	default:
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "String")
	}
}

// from Symbol to Symbol
func convertSymbolIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Symbol)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Symbol")
	}

	return reflect.ValueOf(ls), nil
}

// from Symbol to Symbol
func convertKeywordIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Keyword)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Symbol")
	}

	return reflect.ValueOf(ls), nil
}

// from String to []byte
func convertBytesIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(String)
	if !ok {
		if rv, ok := o.(*ReflectValue); ok {
			if rv.val.Type() == reflect.TypeFor[[]byte]() {
				return rv.val, nil
			}
		}
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "String")
	}

	return reflect.ValueOf([]byte(ls.S)), nil
}

// from Int to int
func convertIntIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Int)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Int")
	}

	return reflect.ValueOf(ls.Int().I), nil
}

// from Int to int
func convertIntIn8(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Int)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Int")
	}

	return reflect.ValueOf(int8(ls.Int().I)), nil
}

// from Int to int
func convertIntIn16(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Int)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Int")
	}

	return reflect.ValueOf(int16(ls.Int().I)), nil
}

// from Int to int
func convertIntIn32(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Int)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Int")
	}

	return reflect.ValueOf(int32(ls.Int().I)), nil
}

// from Int to int
func convertIntIn64(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Int)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Int")
	}

	return reflect.ValueOf(int64(ls.Int().I)), nil
}

// from Int to int
func convertUintIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Int)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Int")
	}

	return reflect.ValueOf(uint(ls.Int().I)), nil
}

// from Int to int
func convertUintIn8(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Int)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Int")
	}

	return reflect.ValueOf(uint8(ls.Int().I)), nil
}

// from Int to int
func convertUintIn16(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Int)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Int")
	}

	return reflect.ValueOf(uint16(ls.Int().I)), nil
}

// from Int to int
func convertUintIn32(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Int)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Int")
	}

	return reflect.ValueOf(uint32(ls.Int().I)), nil
}

// from Int to int
func convertUintIn64(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Int)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Int")
	}

	return reflect.ValueOf(uint64(ls.Int().I)), nil
}

// from Object to Object
func convertObjectIn(env *Env, index int, o Object) (reflect.Value, error) {
	return reflect.ValueOf(o), nil
}

// from Env* to Env*
func convertEnvIn(env *Env, index int, o Object) (reflect.Value, error) {
	return reflect.ValueOf(env), nil
}

// from Seqable to Seqable
func convertSeqableIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Seqable)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Seqable")
	}

	return reflect.ValueOf(ls), nil
}

// from Map to Map
func convertMapIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Map)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Map")
	}

	return reflect.ValueOf(ls), nil
}

// from ReflectValue to *type
func convertReflectTypeIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(*ReflectType)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "ReflectType")
	}

	return reflect.ValueOf(ls.typ), nil
}

// from ReflectValue to *type
func convertReflectValueInDirect(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(*ReflectValue)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "ReflectValue")
	}

	return reflect.ValueOf(ls), nil
}

// from ReflectValue to *type
func convertReflectValueInAny(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(*ReflectValue)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "ReflectValue")
	}

	return reflect.ValueOf(ls.val), nil
}

// from ReflectValue to *type
func convertReflectValueIn(env *Env, index int, o Object, at reflect.Type) (reflect.Value, error) {
	ls, ok := o.(*ReflectValue)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, at.Name())
	}

	if at.Kind() == reflect.Interface {
		if ls.val.Type().AssignableTo(at) {
			return ls.val, nil
		}
	}

	if ls.val.Type() != at {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, at.Name())
	}

	return ls.val, nil
}

// from any to Object|ReflectValue
func convertReflectValueOut(env *Env, s reflect.Value) (Object, error) {
	o := s.Interface()
	if sv, ok := o.(Object); ok {
		return sv, nil
	}

	switch sv := s.Interface().(type) {
	case reflect.Value:
		return &ReflectValue{val: sv}, nil
	case *ReflectValue:
		return sv, nil
	default:
		return &ReflectValue{val: s}, nil
	}
}

func convertCallableIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Callable)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Seqable")
	}

	return reflect.ValueOf(ls), nil
}

func convertFromInt(env *Env, rv reflect.Value) (Object, error) {
	i := rv.Int()

	if i > math.MaxInt {
		return MakeBigInt(i), nil
	}

	return MakeInt(int(i)), nil
}

func convertFromUInt(env *Env, rv reflect.Value) (Object, error) {
	i := rv.Uint()

	if i > math.MaxUint {
		return MakeBigInt(int64(i)), nil
	}

	return MakeInt(int(i)), nil
}

var rawFunc = reflect.TypeFor[func(*Env, Args) (Object, error)]()

func (c *ConvRegistry) makeFuncConvertIn(env *Env, target Callable, ft reflect.Type) reflect.Value {
	var (
		go2lace []outConv
		lace2go []inConv
	)

	for i := 0; i < ft.NumIn(); i++ {
		at := ft.In(i)
		go2lace = append(go2lace, c.convRet(at))
	}

	var (
		errIdx    int = -1
		zeroRet   []reflect.Value
		retIdx    int
		retValues int
	)

	for i := 0; i < ft.NumOut(); i++ {
		at := ft.Out(i)
		zeroRet = append(zeroRet, reflect.Zero(at))

		if at == reflect.TypeFor[error]() {
			errIdx = i
		} else {
			retValues++
			retIdx = i
			fn, _ := c.convArg(at)
			lace2go = append(lace2go, fn)
		}
	}

	trampoline := reflect.MakeFunc(ft, func(args []reflect.Value) (results []reflect.Value) {
		var objs []Object

		var (
			err  error
			sequ Seqable
			ok   bool
			seq  Seq
			ret  Object
			obj  Object
		)

		rets := make([]reflect.Value, ft.NumOut())
		copy(rets, zeroRet)

		for i, f := range go2lace {
			obj, err := f(env, args[i])
			if err != nil {
				goto returnErr
			}

			objs = append(objs, obj)
		}

		ret, err = target.Call(env, objs)
		if err != nil {
			goto returnErr
		}

		if retValues == 0 {
			return rets
		}

		if retValues == 1 {
			v, err := lace2go[0](env, -1, ret)
			if err != nil {
				goto returnErr
			}

			rets[retIdx] = v

			return rets
		}

		sequ, ok = ret.(Seqable)
		if !ok {
			err = fmt.Errorf("function needed a seqable, did not return one")
			goto returnErr
		}

		seq = sequ.Seq()

		for i, oc := range lace2go {
			var v reflect.Value

			if i != errIdx {
				obj, err = seq.First(env)
				if err != nil {
					goto returnErr
				}

				v, err = oc(env, -1, obj)
				if err != nil {
					goto returnErr
				}

				rets[i] = reflect.ValueOf(v)
				seq, err = seq.Rest(env)
				if err != nil {
					goto returnErr
				}
			}
		}

		return rets

	returnErr:
		if errIdx >= 0 {
			rets[errIdx] = reflect.ValueOf(err)
			return rets
		} else {
			panic(err)
		}
	})

	return trampoline
}
