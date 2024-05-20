package core

import "io"

type (
	ArrayMap struct {
		InfoHolder
		MetaHolder
		arr []any
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

func SafeMerge(env *Env, m1, m2 Map) (Map, error) {
	if m1 == nil {
		return m2, nil
	}
	if m2 == nil {
		return m1, nil
	}
	return m1.Merge(env, m2)
}

func (seq *ArrayMapSeq) sequential() {}

func (seq *ArrayMapSeq) Equals(env *Env, other interface{}) bool {
	return IsSeqEqual(env, seq, other)
}

func (seq *ArrayMapSeq) ToString(env *Env, escape bool) (string, error) {
	return SeqToString(env, seq, escape)
}

func (seq *ArrayMapSeq) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintSeq(env, seq, w, indent)
}

func (seq *ArrayMapSeq) WithMeta(env *Env, meta Map) (any, error) {
	res := *seq
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (seq *ArrayMapSeq) GetType() *Type {
	return TYPE.ArrayMapSeq
}

func (seq *ArrayMapSeq) Hash(env *Env) (uint32, error) {
	return hashOrdered(env, seq)
}

func (seq *ArrayMapSeq) Seq() Seq {
	return seq
}

func (seq *ArrayMapSeq) First(env *Env) (any, error) {
	if seq.index < len(seq.m.arr) {
		return NewVectorFrom(seq.m.arr[seq.index], seq.m.arr[seq.index+1]), nil
	}
	return NIL, nil
}

func (seq *ArrayMapSeq) Rest(env *Env) (Seq, error) {
	if seq.index < len(seq.m.arr) {
		return &ArrayMapSeq{m: seq.m, index: seq.index + 2}, nil
	}
	return EmptyList, nil
}

func (seq *ArrayMapSeq) IsEmpty(env *Env) (bool, error) {
	return seq.index >= len(seq.m.arr), nil
}

func (seq *ArrayMapSeq) Cons(obj any) Seq {
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

func (v *ArrayMap) WithMeta(env *Env, meta Map) (any, error) {
	res := *v
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (m *ArrayMap) indexOf(env *Env, key any) int {
	for i := 0; i < len(m.arr); i += 2 {
		if Equals(env, m.arr[i], key) {
			return i
		}
	}
	return -1
}

func (m *ArrayMap) Get(env *Env, key any) (bool, any, error) {
	i := m.indexOf(env, key)
	if i != -1 {
		return true, m.arr[i+1], nil
	}
	return false, nil, nil
}

type Equ interface {
	any

	Is(o any) bool
	IsHash() uint32
}

func (m *ArrayMap) indexOfEqu(key Equ) int {
	for i := 0; i < len(m.arr); i += 2 {
		if key.Is(m.arr[i]) {
			return i
		}
	}
	return -1
}

func (m *ArrayMap) GetEqu(key Equ) (bool, any) {
	i := m.indexOfEqu(key)
	if i != -1 {
		return true, m.arr[i+1]
	}
	return false, nil
}

func (m *ArrayMap) Set(env *Env, key any, value any) {
	i := m.indexOf(env, key)
	if i != -1 {
		m.arr[i+1] = value
	} else {
		m.arr = append(m.arr, key)
		m.arr = append(m.arr, value)
	}
}

func (m *ArrayMap) Add(env *Env, key any, value any) bool {
	i := m.indexOf(env, key)
	if i != -1 {
		return false
	}
	m.arr = append(m.arr, key)
	m.arr = append(m.arr, value)
	return true
}

func (m *ArrayMap) AddEqu(key Equ, value any) bool {
	i := m.indexOfEqu(key)
	if i != -1 {
		return false
	}
	m.arr = append(m.arr, key)
	m.arr = append(m.arr, value)
	return true
}

func (m *ArrayMap) Plus(env *Env, key any, value any) *ArrayMap {
	i := m.indexOf(env, key)
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
	result := ArrayMap{arr: make([]any, len(m.arr), cap(m.arr))}
	copy(result.arr, m.arr)
	return &result
}

func NewArrayMap(key Equ, value any) (Associative, error) {
	m := EmptyArrayMap()

	i := m.indexOfEqu(key)
	if i != -1 {
		m.arr[i+1] = value
		return m, nil
	}
	m.arr = append(m.arr, key, value)
	return m, nil
}

func (m *ArrayMap) Assoc(env *Env, key any, value any) (Associative, error) {
	i := m.indexOf(env, key)
	if i != -1 {
		res := m.Clone()
		res.arr[i+1] = value
		return res, nil
	}
	if int64(len(m.arr)) >= HASHMAP_THRESHOLD {
		hm, err := NewHashMap(env, m.arr...)
		if err != nil {
			return nil, err
		}

		return hm.Assoc(env, key, value)
	}
	res := m.Clone()
	res.arr = append(res.arr, key)
	res.arr = append(res.arr, value)
	return res, nil
}

func (m *ArrayMap) EntryAt(env *Env, key any) (*Vector, error) {
	i := m.indexOf(env, key)
	if i != -1 {
		return NewVectorFrom(key, m.arr[i+1]), nil
	}
	return nil, nil
}

func (m *ArrayMap) Without(env *Env, key any) (Map, error) {
	result := ArrayMap{arr: make([]any, len(m.arr), cap(m.arr))}
	var i, j int
	for i, j = 0, 0; i < len(m.arr); i += 2 {
		if Equals(env, m.arr[i], key) {
			continue
		}
		result.arr[j] = m.arr[i]
		result.arr[j+1] = m.arr[i+1]
		j += 2
	}
	if i != j {
		result.arr = result.arr[:j]
	}
	return &result, nil
}

func (m *ArrayMap) Merge(env *Env, other Map) (Map, error) {
	if other.Count() == 0 {
		return m, nil
	}
	if m.Count() == 0 {
		return other, nil
	}
	res := m.Clone()
	for iter := other.Iter(); iter.HasNext(); {
		p := iter.Next()
		res.Set(env, p.Key, p.Value)
		if int64(len(res.arr)) > HASHMAP_THRESHOLD {
			hm, err := NewHashMap(env, m.arr...)
			if err != nil {
				return nil, err
			}

			return hm.Merge(env, other)
		}
	}

	return res, nil
}

func (m *ArrayMap) Keys() Seq {
	mlen := len(m.arr) / 2
	res := make([]any, mlen)
	for i := 0; i < mlen; i++ {
		res[i] = m.arr[i*2]
	}
	return &ArraySeq{arr: res}
}

func (m *ArrayMap) Vals() Seq {
	mlen := len(m.arr) / 2
	res := make([]any, mlen)
	for i := 0; i < mlen; i++ {
		res[i] = m.arr[i*2+1]
	}
	return &ArraySeq{arr: res}
}

func (m *ArrayMap) Iter() MapIterator {
	return &ArrayMapIterator{m: m}
}

func (m *ArrayMap) Conj(env *Env, obj any) (Conjable, error) {
	return mapConj(env, m, obj)
}

func (m *ArrayMap) ToString(env *Env, escape bool) (string, error) {
	return mapToString(env, m, escape)
}

func (m *ArrayMap) Equals(env *Env, other interface{}) bool {
	return mapEquals(env, m, other)
}

func (m *ArrayMap) GetType() *Type {
	return TYPE.ArrayMap
}

func (m *ArrayMap) Hash(env *Env) (uint32, error) {
	return hashUnordered(env, m.Seq(), 1)
}

func (m *ArrayMap) Seq() Seq {
	return &ArrayMapSeq{m: m, index: 0}
}

func (m *ArrayMap) Call(env *Env, args []any) (any, error) {
	return callMap(env, m, args)
}

var _ Callable = (*ArrayMap)(nil)

func (m *ArrayMap) Empty() Collection {
	return EmptyArrayMap()
}

func (m *ArrayMap) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	return pprintMap(env, m, w, indent)
}
