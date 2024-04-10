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
		Rest() Seq
		IsEmpty() bool
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

		env *Env
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
	for iter1 := iter(seq1); iter1.HasNext(); {
		v, err := iter1.Next(env)
		if err != nil {
			return false
		}

		if !iter2.HasNext() {
			v2, err := iter2.Next(env)
			if err != nil {
				return false
			}
			if !v2.Equals(env, v) {
				return false
			}
		}
	}
	return !iter2.HasNext()
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

func (seq *MappingSeq) Rest() Seq {
	return &MappingSeq{
		seq: seq.seq.Rest(),
		fn:  seq.fn,
	}
}

func (seq *MappingSeq) IsEmpty() bool {
	return seq.seq.IsEmpty()
}

func (seq *MappingSeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: seq}
}

func (seq *MappingSeq) sequential() {}

func (seq *LazySeq) Seq() Seq {
	return seq
}

func (seq *LazySeq) realize() error {
	if seq.seq == nil {
		o, err := seq.fn.Call(seq.env, []Object{})
		if err != nil {
			return err
		}
		v, err := AssertSeqable(seq.env, o, "")
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
	seq.realize()
	return seq.seq.First(env)
}

func (seq *LazySeq) Rest() Seq {
	seq.realize()
	return seq.seq.Rest()
}

func (seq *LazySeq) IsEmpty() bool {
	if err := seq.realize(); err != nil {
		panic(err)
	}
	return seq.seq.IsEmpty()
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
	if seq.IsEmpty() {
		return NIL, nil
	}
	return seq.arr[seq.index], nil
}

func (seq *ArraySeq) Rest() Seq {
	if seq.index+1 < len(seq.arr) {
		return &ArraySeq{index: seq.index + 1, arr: seq.arr}
	}
	return EmptyList
}

func (seq *ArraySeq) IsEmpty() bool {
	return seq.index >= len(seq.arr)
}

func (seq *ArraySeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: seq}
}

func (seq *ArraySeq) sequential() {}

func SeqToString(env *Env, seq Seq, escape bool) (string, error) {
	var b bytes.Buffer
	b.WriteRune('(')
	for iter := iter(seq); iter.HasNext(); {
		v, err := iter.Next(env)
		if err != nil {
			return "", err
		}
		s, err := v.ToString(env, escape)
		if err != nil {
			return "", err
		}
		b.WriteString(s)
		if iter.HasNext() {
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

func (seq *ConsSeq) Rest() Seq {
	return seq.rest
}

func (seq *ConsSeq) IsEmpty() bool {
	return false
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

	iter.seq = iter.seq.Rest()
	return res, nil
}

func (iter *SeqIterator) HasNext() bool {
	return !iter.seq.IsEmpty()
}

func Second(env *Env, seq Seq) (Object, error) {
	return seq.Rest().First(env)
}

func Third(env *Env, seq Seq) (Object, error) {
	return seq.Rest().Rest().First(env)
}

func Fourth(env *Env, seq Seq) (Object, error) {
	return seq.Rest().Rest().Rest().First(env)
}

func ToSlice(env *Env, seq Seq) ([]Object, error) {
	res := make([]Object, 0)
	for !seq.IsEmpty() {
		v, err := seq.First(env)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
		seq = seq.Rest()
	}
	return res, nil
}

func SeqCount(seq Seq) int {
	c := 0
	for !seq.IsEmpty() {
		switch obj := seq.(type) {
		case Counted:
			return c + obj.Count()
		}
		c++
		seq = seq.Rest()
	}
	return c
}

func SeqNth(env *Env, seq Seq, n int) (Object, error) {
	if n < 0 {
		return nil, StubNewError(fmt.Sprintf("Negative index: %d", n))
	}
	i := n
	for !seq.IsEmpty() {
		if i == 0 {
			return seq.First(env)
		}
		seq = seq.Rest()
		i--
	}
	return nil, StubNewError(fmt.Sprintf("Index %d exceeds seq's length %d", n, (n - i)))
}

func SeqTryNth(env *Env, seq Seq, n int, d Object) (Object, error) {
	if n < 0 {
		return d, nil
	}
	i := n
	for !seq.IsEmpty() {
		if i == 0 {
			return seq.First(env)
		}
		seq = seq.Rest()
		i--
	}
	return d, nil
}

func hashUnordered(env *Env, seq Seq, seed uint32) (uint32, error) {
	for !seq.IsEmpty() {
		v, err := seq.First(env)
		if err == nil {
			return 0, err
		}
		sv, err := v.Hash(env)
		if err != nil {
			return 0, err
		}
		seed += sv
		seq = seq.Rest()
	}
	h := getHash()
	h.Write(uint32ToBytes(seed))
	return h.Sum32(), nil
}

func hashOrdered(env *Env, seq Seq) (uint32, error) {
	h := getHash()
	for !seq.IsEmpty() {
		v, err := seq.First(env)
		if err == nil {
			return 0, err
		}
		sv, err := v.Hash(env)
		if err != nil {
			return 0, err
		}
		h.Write(uint32ToBytes(sv))
		seq = seq.Rest()
	}
	return h.Sum32(), nil
}

func pprintSeq(env *Env, seq Seq, w io.Writer, indent int) (int, error) {
	i := indent + 1
	fmt.Fprint(w, "(")
	for iter := iter(seq); iter.HasNext(); {
		v, err := iter.Next(env)
		if err != nil {
			return 0, err
		}
		i, err = pprintObject(env, v, indent+1, w)
		if err != nil {
			return 0, err
		}
		if iter.HasNext() {
			fmt.Fprint(w, "\n")
			writeIndent(w, indent+1)
		}
	}
	fmt.Fprint(w, ")")
	return i + 1, nil
}
