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
		Has(key Equ) bool
		Disjoin(env *Env, key Object) (Set, error)
	}
	MapSet struct {
		InfoHolder
		MetaHolder
		m Map
	}
)

func (v *MapSet) WithMeta(env *Env, meta Map) (Object, error) {
	res := *v
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}

	res.meta = m
	return &res, nil
}

func (set *MapSet) Disjoin(env *Env, key Object) (Set, error) {
	nm, err := set.m.Without(env, key)
	if err != nil {
		return nil, err
	}
	return &MapSet{m: nm}, nil
}

func (set *MapSet) Add(env *Env, obj Object) (bool, error) {
	switch m := set.m.(type) {
	case *ArrayMap:
		return m.Add(env, obj, Boolean{B: true}), nil
	case *HashMap:
		if m.containsKey(env, obj) {
			return false, nil
		}
		v, err := set.m.Assoc(env, obj, Boolean{B: true})
		if err != nil {
			return false, err
		}
		set.m = v.(Map)
		return true, nil
	default:
		return false, nil
	}
}

func (set *MapSet) Conj(env *Env, obj Object) (Conjable, error) {
	v, err := set.m.Assoc(env, obj, Boolean{B: true})
	if err != nil {
		return nil, err
	}

	return &MapSet{m: v.(Map)}, nil
}

func EmptySet() *MapSet {
	return &MapSet{m: EmptyArrayMap()}
}

func (set *MapSet) ToString(env *Env, escape bool) (string, error) {
	var b bytes.Buffer
	b.WriteString("#{")
	for iter := iter(set.m.Keys()); iter.HasNext(env); {
		v, err := iter.Next(env)
		if err != nil {
			return "", err
		}
		s, err := v.ToString(env, escape)
		b.WriteString(s)
		if iter.HasNext(env) {
			b.WriteRune(' ')
		}
	}
	b.WriteRune('}')
	return b.String(), nil
}

func (set *MapSet) Equals(env *Env, other interface{}) bool {
	switch otherSet := other.(type) {
	case *MapSet:
		return set.m.Equals(env, otherSet.m)
	default:
		return false
	}
}

func (set *MapSet) Get(env *Env, key Object) (bool, Object, error) {
	ok, _, err := set.m.Get(env, key)
	if err != nil {
		return false, nil, err
	}

	if ok {
		return true, key, nil
	}

	return false, nil, nil
}

func (set *MapSet) Has(key Equ) bool {
	ok, _ := set.m.GetEqu(key)
	return ok
}

func (seq *MapSet) GetType() *Type {
	return TYPE.MapSet
}

func (set *MapSet) Hash(env *Env) (uint32, error) {
	return hashUnordered(env, set.Seq(), 2)
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

	ok, _, err := set.Get(env, args[0])
	if err != nil {
		return nil, err
	}

	if ok {
		return args[0], nil
	}

	return NIL, nil
}

var _ Callable = (*MapSet)(nil)

func (set *MapSet) Empty() Collection {
	return EmptySet()
}

func NewSetFromSeq(env *Env, s Seq) (*MapSet, error) {
	res := EmptySet()
	for !s.IsEmpty(env) {
		v, err := s.First(env)
		if err != nil {
			return nil, err
		}
		res.Add(env, v)
		s = s.Rest(env)
	}
	return res, nil
}

func (set *MapSet) Pprint(env *Env, w io.Writer, indent int) (int, error) {
	i := indent + 1
	fmt.Fprint(w, "#{")
	for iter := iter(set.m.Keys()); iter.HasNext(env); {
		v, err := iter.Next(env)
		if err != nil {
			return 0, err
		}
		i, err = pprintObject(env, v, indent+2, w)
		if err != nil {
			return 0, err
		}
		if iter.HasNext(env) {
			fmt.Fprint(w, "\n")
			writeIndent(w, indent+2)
		}
	}
	fmt.Fprint(w, "}")
	return i + 1, nil
}
