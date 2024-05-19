package core

import (
	"fmt"
	"io"
	"reflect"
)

func (r *ReflectValue) Clone() *ReflectValue {
	val := reflect.Indirect(r.val)
	t := val.Type()

	nv := reflect.New(t)

	nv.Set(val)

	return &ReflectValue{
		val: nv,
	}
}

func (r *ReflectValue) isStruct() bool {
	return reflect.Indirect(r.val).Kind() == reflect.Struct
}

func (r *ReflectValue) checkStruct(env *Env) error {
	if !r.isStruct() {
		return env.TypeError(TCContext{Context: "using map functions"}, r, "Map")
	}

	return nil
}

func (r *ReflectValue) Assoc(env *Env, key Object, value Object) (Associative, error) {
	if err := r.checkStruct(env); err != nil {
		return nil, err
	}

	val := reflect.Indirect(r.val)

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

func (r *ReflectValue) Set(env *Env, key, value Object) error {
	if err := r.checkStruct(env); err != nil {
		return err
	}

	val := reflect.Indirect(r.val)

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

func (r *ReflectValue) Conj(env *Env, obj Object) (Conjable, error) {
	if err := r.checkStruct(env); err != nil {
		return nil, err
	}

	return mapConj(env, r, obj)
}

func (r *ReflectValue) Count() int {
	if !r.isStruct() {
		return 0
	}
	val := reflect.Indirect(r.val)
	return val.NumField()
}

func (r *ReflectValue) EntryAt(env *Env, key Object) (*Vector, error) {
	if err := r.checkStruct(env); err != nil {
		return nil, err
	}

	val := reflect.Indirect(r.val)

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

func (r *ReflectValue) Get(env *Env, key Object) (bool, Object, error) {
	if err := r.checkStruct(env); err != nil {
		return false, nil, err
	}

	val := reflect.Indirect(r.val)

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

func (r *ReflectValue) GetEqu(key Equ) (bool, Object) {
	if !r.isStruct() {
		return false, nil
	}

	val := reflect.Indirect(r.val)

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

func (r *ReflectValue) Iter() MapIterator {
	if !r.isStruct() {
		return emptyMapIterator
	}

	v := reflect.Indirect(r.val)

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
	val := MakeReflectValue(fv.Interface())

	return &Pair{Key: name, Value: val}
}

func (r *ReflectValue) Keys() Seq {
	if !r.isStruct() {
		return NIL
	}

	t := reflect.Indirect(r.val).Type()

	var ret []Object

	for i := 0; i < t.NumField(); i++ {
		ret = append(ret, MakeKeyword(t.Field(i).Name))
	}

	return &ArraySeq{arr: ret}
}

func (r *ReflectValue) Vals() Seq {
	if !r.isStruct() {
		return NIL
	}

	v := reflect.Indirect(r.val)

	var ret []Object

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

func (r *ReflectValue) Merge(env *Env, other Map) (Map, error) {
	if err := r.checkStruct(env); err != nil {
		return nil, err
	}

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

func (r *ReflectValue) Without(env *Env, key Object) (Map, error) {
	if err := r.checkStruct(env); err != nil {
		return nil, err
	}

	if r.Count() == 0 {
		return r, nil
	}
	res := r.Clone()
	for iter := r.Iter(); iter.HasNext(); {
		p := iter.Next()

		if p.Key.Equals(env, key) {
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
	m     *ReflectValue
	index int
}

func (r *ReflectValue) Seq() Seq {
	if !r.isStruct() {
		return NIL
	}

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

func (seq *reflectStructSeq) WithMeta(env *Env, meta Map) (Object, error) {
	res := *seq
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (seq *reflectStructSeq) GetType() *Type {
	return TYPE.ArrayMapSeq
}

func (seq *reflectStructSeq) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, seq)
}

func (seq *reflectStructSeq) Seq() Seq {
	return seq
}

func (seq *reflectStructSeq) First(env *Env) (Object, error) {
	t := seq.m.val.Type()

	if seq.index >= t.NumField() {
		return NIL, nil
	}

	tf := seq.m.val.Type().Field(seq.index)
	fv := seq.m.val.Field(seq.index)

	name := MakeKeyword(tf.Name)
	val := MakeReflectValue(fv.Interface())

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

func (seq *reflectStructSeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: seq}
}
