package core

import (
	"bytes"
	"fmt"
	"io"
)

type (
	Seq interface {
		Seqable
		Object
		First(env *Env) (Object, error)
		Rest(env *Env) (Seq, error)
		IsEmpty(env *Env) (bool, error)
		Cons(obj Object) Seq
	}
	Seqable interface {
		Seq() Seq
	}
	SeqIterator struct {
		seq Seq
	}
	ConsSeq struct {
		InfoHolder
		MetaHolder
		first Object
		rest  Seq
	}
	ArraySeq struct {
		InfoHolder
		MetaHolder
		arr   []Object
		index int
	}
	LazySeq struct {
		InfoHolder
		MetaHolder

		fn  Callable
		seq Seq
	}
	MappingSeq struct {
		InfoHolder
		MetaHolder
		seq Seq
		fn  func(env *Env, obj Object) (Object, error)
	}
)

func SeqsEqual(env *Env, seq1, seq2 Seq) bool {
	iter2 := iter(seq2)
	for iter1 := iter(seq1); iter1.HasNext(env); {
		v, err := iter1.Next(env)
		if err != nil {
			return false
		}

		if !iter2.HasNext(env) {
			v2, err := iter2.Next(env)
			if err != nil {
				return false
			}
			if !v2.Equals(env, v) {
				return false
			}
		}
	}
	return !iter2.HasNext(env)
}

func IsSeqEqual(env *Env, seq Seq, other interface{}) bool {
	if seq == other {
		return true
	}
	switch other := other.(type) {
	case Sequential:
		switch other := other.(type) {
		case Seqable:
			return SeqsEqual(env, seq, other.Seq())
		}
	}
	return false
}

func (seq *MappingSeq) Seq() Seq {
	return seq
}

func (seq *MappingSeq) Equals(env *Env, other interface{}) bool {
	return IsSeqEqual(env, seq, other)
}

func (seq *MappingSeq) ToString(env *Env, escape bool) (string, error) {
	return SeqToString(env, seq, escape)
}

func (seq *MappingSeq) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintSeq(env, seq, w, indent)
}

func (seq *MappingSeq) WithMeta(env *Env, meta Map) (Object, error) {
	res := *seq
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (seq *MappingSeq) GetType() *Type {
	return TYPE.MappingSeq
}

func (seq *MappingSeq) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, seq)
}

func (seq *MappingSeq) First(env *Env) (Object, error) {
	f, err := seq.seq.First(env)
	if err != nil {
		return nil, err
	}

	return seq.fn(env, f)
}

func (seq *MappingSeq) Rest(env *Env) (Seq, error) {
	x, err := seq.seq.Rest(env)
	if err != nil {
		return nil, err
	}

	return &MappingSeq{
		seq: x,
		fn:  seq.fn,
	}, nil
}

func (seq *MappingSeq) IsEmpty(env *Env) (bool, error) {
	return seq.seq.IsEmpty(env)
}

func (seq *MappingSeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: seq}
}

func (seq *MappingSeq) sequential() {}

func (seq *LazySeq) Seq() Seq {
	return seq
}

func (seq *LazySeq) realize(env *Env) error {
	if seq.seq == nil {
		o, err := seq.fn.Call(env, []Object{})
		if err != nil {
			return err
		}
		v, err := AssertSeqable(env, o, "")
		if err != nil {
			return err
		}
		seq.seq = v.Seq()
	}

	return nil
}

func (seq *LazySeq) IsRealized() bool {
	return seq.seq != nil
}

func (seq *LazySeq) Equals(env *Env, other interface{}) bool {
	return IsSeqEqual(env, seq, other)
}

func (seq *LazySeq) ToString(env *Env, escape bool) (string, error) {
	return SeqToString(env, seq, escape)
}

func (seq *LazySeq) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintSeq(env, seq, w, indent)
}

func (seq *LazySeq) WithMeta(env *Env, meta Map) (Object, error) {
	res := *seq
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (seq *LazySeq) GetType() *Type {
	return TYPE.LazySeq
}

func (seq *LazySeq) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, seq)
}

func (seq *LazySeq) First(env *Env) (Object, error) {
	err := seq.realize(env)
	if err != nil {
		return nil, err
	}
	return seq.seq.First(env)
}

func (seq *LazySeq) Rest(env *Env) (Seq, error) {
	err := seq.realize(env)
	if err != nil {
		return nil, err
	}
	return seq.seq.Rest(env)
}

func (seq *LazySeq) IsEmpty(env *Env) (bool, error) {
	if err := seq.realize(env); err != nil {
		return false, err
	}
	return seq.seq.IsEmpty(env)
}

func (seq *LazySeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: seq}
}

func (seq *LazySeq) sequential() {}

func NewLazySeq(c Callable) *LazySeq {
	return &LazySeq{fn: c}
}

func (seq *ArraySeq) Seq() Seq {
	return seq
}

func (seq *ArraySeq) Equals(env *Env, other interface{}) bool {
	return IsSeqEqual(env, seq, other)
}

func (seq *ArraySeq) ToString(env *Env, escape bool) (string, error) {
	return SeqToString(env, seq, escape)
}

func (seq *ArraySeq) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintSeq(env, seq, w, indent)
}

func (seq *ArraySeq) WithMeta(env *Env, meta Map) (Object, error) {
	res := *seq
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (seq *ArraySeq) GetType() *Type {
	return TYPE.ArraySeq
}

func (seq *ArraySeq) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, seq)
}

func (seq *ArraySeq) First(env *Env) (Object, error) {
	ok, err := seq.IsEmpty(env)
	if err != nil {
		return nil, err
	}

	if ok {
		return NIL, nil
	}

	return seq.arr[seq.index], nil
}

func (seq *ArraySeq) Rest(env *Env) (Seq, error) {
	if seq.index+1 < len(seq.arr) {
		return &ArraySeq{index: seq.index + 1, arr: seq.arr}, nil
	}
	return EmptyList, nil
}

func (seq *ArraySeq) IsEmpty(env *Env) (bool, error) {
	return seq.index >= len(seq.arr), nil
}

func (seq *ArraySeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: seq}
}

func (seq *ArraySeq) sequential() {}

func SeqToString(env *Env, seq Seq, escape bool) (string, error) {
	var b bytes.Buffer
	b.WriteRune('(')
	for iter := iter(seq); iter.HasNext(env); {
		v, err := iter.Next(env)
		if err != nil {
			return "", err
		}
		s, err := v.ToString(env, escape)
		if err != nil {
			return "", err
		}
		b.WriteString(s)
		if iter.HasNext(env) {
			b.WriteRune(' ')
		}
	}
	b.WriteRune(')')
	return b.String(), nil
}

func (seq *ConsSeq) WithMeta(env *Env, meta Map) (Object, error) {
	res := *seq
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (seq *ConsSeq) Seq() Seq {
	return seq
}

func (seq *ConsSeq) Equals(env *Env, other interface{}) bool {
	return IsSeqEqual(env, seq, other)
}

func (seq *ConsSeq) ToString(env *Env, escape bool) (string, error) {
	return SeqToString(env, seq, escape)
}

func (seq *ConsSeq) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintSeq(env, seq, w, indent)
}

func (seq *ConsSeq) GetType() *Type {
	return TYPE.ConsSeq
}

func (seq *ConsSeq) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, seq)
}

func (seq *ConsSeq) First(env *Env) (Object, error) {
	return seq.first, nil
}

func (seq *ConsSeq) Rest(env *Env) (Seq, error) {
	return seq.rest, nil
}

func (seq *ConsSeq) IsEmpty(env *Env) (bool, error) {
	return false, nil
}

func (seq *ConsSeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: seq}
}

func (seq *ConsSeq) sequential() {}

func NewConsSeq(first Object, rest Seq) *ConsSeq {
	return &ConsSeq{
		first: first,
		rest:  rest,
	}
}

func iter(seq Seq) *SeqIterator {
	return &SeqIterator{seq: seq}
}

func (iter *SeqIterator) Next(env *Env) (Object, error) {
	res, err := iter.seq.First(env)
	if err != nil {
		return nil, err
	}

	iter.seq, err = iter.seq.Rest(env)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (iter *SeqIterator) HasNext(env *Env) bool {
	ok, err := iter.seq.IsEmpty(env)
	if err != nil {
		return false
	}

	return !ok
}

func Second(env *Env, seq Seq) (Object, error) {
	x, err := seq.Rest(env)
	if err != nil {
		return nil, err
	}
	return x.First(env)
}

func Third(env *Env, seq Seq) (Object, error) {
	x, err := seq.Rest(env)
	if err != nil {
		return nil, err
	}

	x, err = x.Rest(env)
	if err != nil {
		return nil, err
	}

	return x.First(env)
}

func Fourth(env *Env, seq Seq) (Object, error) {
	x, err := seq.Rest(env)
	if err != nil {
		return nil, err
	}

	x, err = x.Rest(env)
	if err != nil {
		return nil, err
	}

	x, err = x.Rest(env)
	if err != nil {
		return nil, err
	}

	return x.First(env)
}

func ToSlice(env *Env, seq Seq) ([]Object, error) {
	res := make([]Object, 0)
	for {
		ok, err := seq.IsEmpty(env)
		if err != nil {
			return nil, err
		}
		if ok {
			break
		}

		v, err := seq.First(env)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
		seq, err = seq.Rest(env)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func SeqCount(env *Env, seq Seq) (int, error) {
	c := 0
	for {
		ok, err := seq.IsEmpty(env)
		if err != nil {
			return 0, err
		}
		if ok {
			break
		}

		switch obj := seq.(type) {
		case Counted:
			return c + obj.Count(), nil
		}
		c++
		seq, err = seq.Rest(env)
		if err != nil {
			return 0, err
		}
	}

	return c, nil
}

func SeqNth(env *Env, seq Seq, n int) (Object, error) {
	if n < 0 {
		return nil, StubNewError(fmt.Sprintf("Negative index: %d", n))
	}
	i := n
	for {
		ok, err := seq.IsEmpty(env)
		if err != nil {
			return nil, err
		}
		if ok {
			break
		}

		if i == 0 {
			return seq.First(env)
		}
		seq, err = seq.Rest(env)
		if err != nil {
			return nil, err
		}
		i--
	}
	return nil, env.RT.NewError(fmt.Sprintf("Index %d exceeds seq's length %d", n, (n - i)))
}

func SeqTryNth(env *Env, seq Seq, n int, d Object) (Object, error) {
	if n < 0 {
		return d, nil
	}
	i := n
	for {
		ok, err := seq.IsEmpty(env)
		if err != nil {
			return nil, err
		}
		if ok {
			break
		}

		if i == 0 {
			return seq.First(env)
		}
		seq, err = seq.Rest(env)
		if err != nil {
			return nil, err
		}

		i--
	}
	return d, nil
}

func hashUnordered(env *Env, seq Seq, seed uint32) (uint32, error) {
	for {
		ok, err := seq.IsEmpty(env)
		if err != nil {
			return 0, err
		}
		if ok {
			break
		}

		v, err := seq.First(env)
		if err == nil {
			return 0, err
		}
		sv, err := v.Hash(env)
		if err != nil {
			return 0, err
		}
		seed += sv
		seq, err = seq.Rest(env)
		if err != nil {
			return 0, err
		}
	}
	h := getHash()
	h.Write(uint32ToBytes(seed))
	return h.Sum32(), nil
}

func hashOrdered(env *Env, seq Seq) (uint32, error) {
	h := getHash()
	for {
		ok, err := seq.IsEmpty(env)
		if err != nil {
			return 0, err
		}
		if ok {
			break
		}

		v, err := seq.First(env)
		if err == nil {
			return 0, err
		}
		sv, err := v.Hash(env)
		if err != nil {
			return 0, err
		}
		h.Write(uint32ToBytes(sv))
		seq, err = seq.Rest(env)
		if err != nil {
			return 0, err
		}
	}
	return h.Sum32(), nil
}

var cnt int

func pprintSeq(env *Env, seq Seq, w io.Writer, indent int) (int, error) {
	cnt++
	i := indent + 1
	fmt.Fprint(w, "(")
	for iter := iter(seq); iter.HasNext(env); {
		v, err := iter.Next(env)
		if err != nil {
			return 0, err
		}
		i, err = pprintObject(env, v, indent+1, w)
		if err != nil {
			return 0, err
		}
		if iter.HasNext(env) {
			fmt.Fprint(w, "\n")
			writeIndent(w, indent+1)
		}
	}
	fmt.Fprint(w, ")")
	return i + 1, nil
}
