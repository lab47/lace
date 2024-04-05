package core

import (
	"bytes"
	"fmt"
	"io"
)

type (
	Set interface {
		Conjable
		Gettable
		Disjoin(key Object) Set
	}
	MapSet struct {
		InfoHolder
		MetaHolder
		m Map
	}
)

func (v *MapSet) WithMeta(meta Map) (Object, error) {
	res := *v
	m, err := SafeMerge(res.meta, meta)
	if err != nil {
		return nil, err
	}

	res.meta = m
	return &res, nil
}

func (set *MapSet) Disjoin(key Object) Set {
	return &MapSet{m: set.m.Without(key)}
}

func (set *MapSet) Add(obj Object) (bool, error) {
	switch m := set.m.(type) {
	case *ArrayMap:
		return m.Add(obj, Boolean{B: true}), nil
	case *HashMap:
		if m.containsKey(obj) {
			return false, nil
		}
		v, err := set.m.Assoc(obj, Boolean{B: true})
		if err != nil {
			return false, err
		}
		set.m = v.(Map)
		return true, nil
	default:
		return false, nil
	}
}

func (set *MapSet) Conj(obj Object) (Conjable, error) {
	v, err := set.m.Assoc(obj, Boolean{B: true})
	if err != nil {
		return nil, err
	}

	return &MapSet{m: v.(Map)}, nil
}

func EmptySet() *MapSet {
	return &MapSet{m: EmptyArrayMap()}
}

func (set *MapSet) ToString(escape bool) string {
	var b bytes.Buffer
	b.WriteString("#{")
	for iter := iter(set.m.Keys()); iter.HasNext(); {
		b.WriteString(iter.Next().ToString(escape))
		if iter.HasNext() {
			b.WriteRune(' ')
		}
	}
	b.WriteRune('}')
	return b.String()
}

func (set *MapSet) Equals(other interface{}) bool {
	switch otherSet := other.(type) {
	case *MapSet:
		return set.m.Equals(otherSet.m)
	default:
		return false
	}
}

func (set *MapSet) Get(key Object) (bool, Object) {
	if ok, _ := set.m.Get(key); ok {
		return true, key
	}
	return false, nil
}

func (seq *MapSet) GetType() *Type {
	return TYPE.MapSet
}

func (set *MapSet) Hash() uint32 {
	return hashUnordered(set.Seq(), 2)
}

func (set *MapSet) Seq() Seq {
	return set.m.Keys()
}

func (set *MapSet) Count() int {
	return set.m.Count()
}

func (set *MapSet) Call(env *Env, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 1); err != nil {
		return nil, err
	}

	if ok, _ := set.Get(args[0]); ok {
		return args[0], nil
	}
	return NIL, nil
}

var _ Callable = (*MapSet)(nil)

func (set *MapSet) Empty() Collection {
	return EmptySet()
}

func NewSetFromSeq(s Seq) *MapSet {
	res := EmptySet()
	for !s.IsEmpty() {
		res.Add(s.First())
		s = s.Rest()
	}
	return res
}

func (set *MapSet) Pprint(w io.Writer, indent int) int {
	i := indent + 1
	fmt.Fprint(w, "#{")
	for iter := iter(set.m.Keys()); iter.HasNext(); {
		i = pprintObject(iter.Next(), indent+2, w)
		if iter.HasNext() {
			fmt.Fprint(w, "\n")
			writeIndent(w, indent+2)
		}
	}
	fmt.Fprint(w, "}")
	return i + 1
}
