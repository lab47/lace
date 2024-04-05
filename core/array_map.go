package core

import "io"

type (
	ArrayMap struct {
		InfoHolder
		MetaHolder
		arr []Object
	}
	ArrayMapIterator struct {
		m       *ArrayMap
		current int
	}
	ArrayMapSeq struct {
		InfoHolder
		MetaHolder
		m     *ArrayMap
		index int
	}
)

var (
	HASHMAP_THRESHOLD int64 = 16
)

func EmptyArrayMap() *ArrayMap {
	return &ArrayMap{}
}

func ArraySeqFromArrayMap(m *ArrayMap) *ArraySeq {
	return &ArraySeq{arr: m.arr}
}

func SafeMerge(m1, m2 Map) (Map, error) {
	if m1 == nil {
		return m2, nil
	}
	if m2 == nil {
		return m1, nil
	}
	return m1.Merge(m2)
}

func (seq *ArrayMapSeq) sequential() {}

func (seq *ArrayMapSeq) Equals(other interface{}) bool {
	return IsSeqEqual(seq, other)
}

func (seq *ArrayMapSeq) ToString(escape bool) string {
	return SeqToString(seq, escape)
}

func (seq *ArrayMapSeq) Pprint(w io.Writer, indent int) int {
	return pprintSeq(seq, w, indent)
}

func (seq *ArrayMapSeq) WithMeta(meta Map) (Object, error) {
	res := *seq
	m, err := SafeMerge(res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (seq *ArrayMapSeq) GetType() *Type {
	return TYPE.ArrayMapSeq
}

func (seq *ArrayMapSeq) Hash() uint32 {
	return hashOrdered(seq)
}

func (seq *ArrayMapSeq) Seq() Seq {
	return seq
}

func (seq *ArrayMapSeq) First() Object {
	if seq.index < len(seq.m.arr) {
		return NewVectorFrom(seq.m.arr[seq.index], seq.m.arr[seq.index+1])
	}
	return NIL
}

func (seq *ArrayMapSeq) Rest() Seq {
	if seq.index < len(seq.m.arr) {
		return &ArrayMapSeq{m: seq.m, index: seq.index + 2}
	}
	return EmptyList
}

func (seq *ArrayMapSeq) IsEmpty() bool {
	return seq.index >= len(seq.m.arr)
}

func (seq *ArrayMapSeq) Cons(obj Object) Seq {
	return &ConsSeq{first: obj, rest: seq}
}

func (iter *ArrayMapIterator) Next() *Pair {
	res := Pair{
		Key:   iter.m.arr[iter.current],
		Value: iter.m.arr[iter.current+1],
	}
	iter.current += 2
	return &res
}

func (iter *ArrayMapIterator) HasNext() bool {
	return iter.current < len(iter.m.arr)
}

func (v *ArrayMap) WithMeta(meta Map) (Object, error) {
	res := *v
	m, err := SafeMerge(res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (m *ArrayMap) indexOf(key Object) int {
	for i := 0; i < len(m.arr); i += 2 {
		if m.arr[i].Equals(key) {
			return i
		}
	}
	return -1
}

func (m *ArrayMap) Get(key Object) (bool, Object) {
	i := m.indexOf(key)
	if i != -1 {
		return true, m.arr[i+1]
	}
	return false, nil
}

func (m *ArrayMap) Set(key Object, value Object) {
	i := m.indexOf(key)
	if i != -1 {
		m.arr[i+1] = value
	} else {
		m.arr = append(m.arr, key)
		m.arr = append(m.arr, value)
	}
}

func (m *ArrayMap) Add(key Object, value Object) bool {
	i := m.indexOf(key)
	if i != -1 {
		return false
	}
	m.arr = append(m.arr, key)
	m.arr = append(m.arr, value)
	return true
}

func (m *ArrayMap) Plus(key Object, value Object) *ArrayMap {
	i := m.indexOf(key)
	if i != -1 {
		return m
	}
	m.arr = append(m.arr, key)
	m.arr = append(m.arr, value)
	return m
}

func (m *ArrayMap) Count() int {
	return len(m.arr) / 2
}

func (m *ArrayMap) Clone() *ArrayMap {
	result := ArrayMap{arr: make([]Object, len(m.arr), cap(m.arr))}
	copy(result.arr, m.arr)
	return &result
}

func (m *ArrayMap) Assoc(key Object, value Object) (Associative, error) {
	i := m.indexOf(key)
	if i != -1 {
		res := m.Clone()
		res.arr[i+1] = value
		return res, nil
	}
	if int64(len(m.arr)) >= HASHMAP_THRESHOLD {
		hm, err := NewHashMap(m.arr...)
		if err != nil {
			return nil, err
		}

		return hm.Assoc(key, value)
	}
	res := m.Clone()
	res.arr = append(res.arr, key)
	res.arr = append(res.arr, value)
	return res, nil
}

func (m *ArrayMap) EntryAt(key Object) (*Vector, error) {
	i := m.indexOf(key)
	if i != -1 {
		return NewVectorFrom(key, m.arr[i+1]), nil
	}
	return nil, nil
}

func (m *ArrayMap) Without(key Object) Map {
	result := ArrayMap{arr: make([]Object, len(m.arr), cap(m.arr))}
	var i, j int
	for i, j = 0, 0; i < len(m.arr); i += 2 {
		if m.arr[i].Equals(key) {
			continue
		}
		result.arr[j] = m.arr[i]
		result.arr[j+1] = m.arr[i+1]
		j += 2
	}
	if i != j {
		result.arr = result.arr[:j]
	}
	return &result
}

func (m *ArrayMap) Merge(other Map) (Map, error) {
	if other.Count() == 0 {
		return m, nil
	}
	if m.Count() == 0 {
		return other, nil
	}
	res := m.Clone()
	for iter := other.Iter(); iter.HasNext(); {
		p := iter.Next()
		res.Set(p.Key, p.Value)
		if int64(len(res.arr)) > HASHMAP_THRESHOLD {
			hm, err := NewHashMap(m.arr...)
			if err != nil {
				return nil, err
			}

			return hm.Merge(other)
		}
	}

	return res, nil
}

func (m *ArrayMap) Keys() Seq {
	mlen := len(m.arr) / 2
	res := make([]Object, mlen)
	for i := 0; i < mlen; i++ {
		res[i] = m.arr[i*2]
	}
	return &ArraySeq{arr: res}
}

func (m *ArrayMap) Vals() Seq {
	mlen := len(m.arr) / 2
	res := make([]Object, mlen)
	for i := 0; i < mlen; i++ {
		res[i] = m.arr[i*2+1]
	}
	return &ArraySeq{arr: res}
}

func (m *ArrayMap) Iter() MapIterator {
	return &ArrayMapIterator{m: m}
}

func (m *ArrayMap) Conj(obj Object) (Conjable, error) {
	return mapConj(m, obj)
}

func (m *ArrayMap) ToString(escape bool) string {
	return mapToString(m, escape)
}

func (m *ArrayMap) Equals(other interface{}) bool {
	return mapEquals(m, other)
}

func (m *ArrayMap) GetType() *Type {
	return TYPE.ArrayMap
}

func (m *ArrayMap) Hash() uint32 {
	return hashUnordered(m.Seq(), 1)
}

func (m *ArrayMap) Seq() Seq {
	return &ArrayMapSeq{m: m, index: 0}
}

func (m *ArrayMap) Call(env *Env, args []Object) (Object, error) {
	return callMap(env, m, args)
}

var _ Callable = (*ArrayMap)(nil)

func (m *ArrayMap) Empty() Collection {
	return EmptyArrayMap()
}

func (m *ArrayMap) Pprint(w io.Writer, indent int) int {
	return pprintMap(m, w, indent)
}
