package core

import "io"

type List struct {
	InfoHolder
	MetaHolder
	first Object
	rest  *List
	count int
}

func NewList(first Object, rest *List) *List {
	result := List{
		first: first,
		rest:  rest,
	}
	if rest != nil {
		result.count = rest.count + 1
	}
	return &result
}

func NewListFrom(objs ...Object) *List {
	res := EmptyList
	for i := len(objs) - 1; i >= 0; i-- {
		res = res.conj(objs[i])
	}
	return res
}

func (list *List) WithMeta(env *Env, meta Map) (Object, error) {
	res := *list
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (list *List) conj(obj Object) *List {
	return NewList(obj, list)
}

func (list *List) Conj(env *Env, obj Object) (Conjable, error) {
	return list.conj(obj), nil
}

func (list *List) ToString(env *Env, escape bool) (string, error) {
	return SeqToString(env, list, escape)
}

func (seq *List) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintSeq(env, seq, w, indent)
}

func (list *List) Equals(env *Env, other interface{}) bool {
	return IsSeqEqual(env, list, other)
}

func (list *List) GetType() *Type {
	return TYPE.List
}

func (list *List) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, list)
}

func (list *List) First(env *Env) (Object, error) {
	return list.first, nil
}

func (list *List) Rest() Seq {
	return list.rest
}

func (list *List) IsEmpty() bool {
	return list.count == 0
}

func (list *List) Cons(obj Object) Seq {
	return list.conj(obj)
}

func (list *List) Seq() Seq {
	return list
}

func (list *List) Second() Object {
	return list.rest.first
}

func (list *List) Third() Object {
	return list.rest.rest.first
}

func (list *List) Forth() Object {
	return list.rest.rest.rest.first
}

func (list *List) Count() int {
	return list.count
}

func (list *List) Empty() Collection {
	return EmptyList
}

func (list *List) Peek(env *Env) (Object, error) {
	return list.first, nil
}

func (list *List) Pop(env *Env) (Stack, error) {
	if list.count == 0 {
		return nil, env.RT.NewError("Can't pop empty list")
	}
	return list.rest, nil
}

func (list *List) sequential() {}

var EmptyList = NewList(Nil{}, nil)

func init() {
	EmptyList.rest = EmptyList
}
