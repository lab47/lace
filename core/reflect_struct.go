package core

import (
	"fmt"
	"io"
	"reflect"
)

type StructMap struct {
	val reflect.Value
}

var _ Map = StructMap{}

func (r StructMap) Clone() StructMap {
	val := reflect.Indirect(r.val)
	t := val.Type()

	nv := reflect.New(t)

	nv.Set(val)

	return StructMap{
		val: nv,
	}
}

func (r StructMap) Assoc(env *Env, key any, value any) (Associative, error) {
	val := r.val

	var name string

	if err := CoerceString(env, key, &name); err != nil {
		return nil, err
	}

	// We only can access public fields, so convert the key to uppercase for
	// clojure-esque ergonomics.
	name = convertMethodName(name)

	field := val.FieldByName(name)
	if !field.IsValid() {
		return nil, env.NewError(fmt.Sprintf("unknown struct field %s", name))
	}

	// Ok, the field is there, switch to a new value now to set

	nv := r.Clone()
	val = nv.val
	val = reflect.Indirect(val)
	field = val.FieldByName(name)

	var (
		frv reflect.Value
		err error
	)

	if field.Type().Kind() == reflect.Func {
		call, ok := value.(Callable)
		if !ok {
			return nil, env.TypeError(TCContext{Context: "struct value"}, value, "Callable")
		}
		frv = convReg.makeFuncConvertIn(env, call, field.Type())
	} else {
		cv, _ := convReg.convArg(field.Type())

		frv, err = cv(env, -1, value)
		if err != nil {
			return nil, err
		}
	}

	if !frv.Type().AssignableTo(field.Type()) {
		return nil, env.NewError(
			fmt.Sprintf("needed type %s, had %T", field.Type(), value))
	}

	field.Set(frv)

	return nv, nil
}

func (r StructMap) Set(env *Env, key, value any) error {
	val := r.val

	var name string

	if err := CoerceString(env, key, &name); err != nil {
		return err
	}

	// We only can access public fields, so convert the key to uppercase for
	// clojure-esque ergonomics.
	name = convertMethodName(name)

	field := val.FieldByName(name)
	if !field.IsValid() {
		return env.NewError(fmt.Sprintf("unknown struct field %s", name))
	}

	var (
		frv reflect.Value
		err error
	)

	if field.Type().Kind() == reflect.Func {
		call, ok := value.(Callable)
		if !ok {
			return env.TypeError(TCContext{Context: "struct value"}, value, "Callable")
		}
		frv = convReg.makeFuncConvertIn(env, call, field.Type())
	} else {
		cv, _ := convReg.convArg(field.Type())

		frv, err = cv(env, -1, value)
		if err != nil {
			return err
		}
	}

	if !frv.Type().AssignableTo(field.Type()) {
		return env.NewError(
			fmt.Sprintf("needed type %s, had %T", field.Type(), value))
	}

	field.Set(frv)

	return nil
}

func (r StructMap) Conj(env *Env, obj any) (Conjable, error) {
	return mapConj(env, r, obj)
}

func (r StructMap) Count() int {
	return r.val.NumField()
}

func (r StructMap) EntryAt(env *Env, key any) (*Vector, error) {
	val := r.val

	var name string

	if err := CoerceString(env, key, &name); err != nil {
		return nil, err
	}

	// We only can access public fields, so convert the key to uppercase for
	// clojure-esque ergonomics.
	name = convertMethodName(name)

	field := val.FieldByName(name)
	if !field.IsValid() {
		return nil, env.NewError(fmt.Sprintf("unknown struct field %s", name))
	}

	rt := field.Type()

	rval, err := convReg.convRet(rt)(field)
	if err != nil {
		if ce, ok := err.(OutConvError); ok {
			err = env.NewError(string(ce))
		}
		return nil, err
	}

	return NewVectorFrom(key, rval), nil
}

func (r StructMap) Get(env *Env, key any) (bool, any, error) {
	val := r.val

	var name string

	if err := CoerceString(env, key, &name); err != nil {
		return false, nil, err
	}

	// We only can access public fields, so convert the key to uppercase for
	// clojure-esque ergonomics.
	name = convertMethodName(name)

	field := val.FieldByName(name)
	if !field.IsValid() {
		return false, nil, nil
	}

	rt := field.Type()

	rval, err := convReg.convRet(rt)(field)
	if err != nil {
		if ce, ok := err.(OutConvError); ok {
			err = env.NewError(string(ce))
		}
		return false, nil, err
	}

	return true, rval, nil
}

func (r StructMap) GetEqu(key Equ) (bool, any) {
	val := r.val

	var name string

	if ok := TryCoerceString(key, &name); !ok {
		return false, nil
	}

	// We only can access public fields, so convert the key to uppercase for
	// clojure-esque ergonomics.
	name = convertMethodName(name)

	field := val.FieldByName(name)
	if !field.IsValid() {
		return false, nil
	}

	rt := field.Type()

	rval, err := convReg.convRet(rt)(field)
	if err != nil {
		// TODO change GetEqu to return the error
		return false, nil
	}

	return true, rval
}

type structIterator struct {
	rv  reflect.Value
	idx int
}

func (r StructMap) Iter() MapIterator {
	v := r.val

	si := &structIterator{
		rv:  v,
		idx: -1,
	}

	return si
}

func (r *structIterator) HasNext() bool {
	return r.rv.NumField() > r.idx
}

func (r *structIterator) Next() *Pair {
	r.idx++
	tf := r.rv.Type().Field(r.idx)
	fv := r.rv.Field(r.idx)

	name := MakeKeyword(tf.Name)
	val := fv.Interface()

	return &Pair{Key: name, Value: val}
}

func (r StructMap) Keys() Seq {
	t := r.val.Type()

	var ret []any

	for i := 0; i < t.NumField(); i++ {
		ret = append(ret, MakeKeyword(t.Field(i).Name))
	}

	return &ArraySeq{arr: ret}
}

func (r StructMap) Vals() Seq {
	v := r.val

	var ret []any

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		rval, err := convReg.convRet(f.Type())(f)
		if err != nil {
			continue
		}

		ret = append(ret, rval)
	}

	return &ArraySeq{arr: ret}
}

func (r StructMap) Merge(env *Env, other Map) (Map, error) {
	if other.Count() == 0 {
		return r, nil
	}
	if r.Count() == 0 {
		return other, nil
	}
	res := r.Clone()
	for iter := other.Iter(); iter.HasNext(); {
		p := iter.Next()
		err := res.Set(env, p.Key, p.Value)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (r StructMap) Without(env *Env, key any) (Map, error) {
	if r.Count() == 0 {
		return r, nil
	}
	res := r.Clone()
	for iter := r.Iter(); iter.HasNext(); {
		p := iter.Next()

		if Equals(env, p.Key, key) {
			continue
		}

		err := res.Set(env, p.Key, p.Value)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

type reflectStructSeq struct {
	InfoHolder
	MetaHolder
	m     StructMap
	index int
}

func (r StructMap) Seq() Seq {
	return &reflectStructSeq{
		m: r,
	}
}

func (seq *reflectStructSeq) sequential() {}

func (seq *reflectStructSeq) Equals(env *Env, other interface{}) bool {
	return IsSeqEqual(env, seq, other)
}

func (seq *reflectStructSeq) ToString(env *Env, escape bool) (string, error) {
	return SeqToString(env, seq, escape)
}

func (seq *reflectStructSeq) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintSeq(env, seq, w, indent)
}

func (seq *reflectStructSeq) WithMeta(env *Env, meta Map) (any, error) {
	res := *seq
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (seq *reflectStructSeq) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, seq)
}

func (seq *reflectStructSeq) Seq() Seq {
	return seq
}

func (seq *reflectStructSeq) First(env *Env) (any, error) {
	t := seq.m.val.Type()

	if seq.index >= t.NumField() {
		return NIL, nil
	}

	tf := seq.m.val.Type().Field(seq.index)
	fv := seq.m.val.Field(seq.index)

	name := MakeKeyword(tf.Name)
	val := fv.Interface()

	return NewVectorFrom(name, val), nil
}

func (seq *reflectStructSeq) Rest(env *Env) (Seq, error) {
	t := seq.m.val.Type()

	if seq.index >= t.NumField() {
		return EmptyList, nil
	}

	return &reflectStructSeq{m: seq.m, index: seq.index + 1}, nil
}

func (seq *reflectStructSeq) IsEmpty(env *Env) (bool, error) {
	t := seq.m.val.Type()
	return seq.index >= t.NumField(), nil
}

func (seq *reflectStructSeq) Cons(obj any) Seq {
	return &ConsSeq{first: obj, rest: seq}
}
