package core

import (
	"bytes"
	"fmt"
	"io"
)

type (
	// A collection that can store unique values.
	//
	//lace:export
	Set interface {
		Conjable
		Gettable
		Has(key Equ) bool
		Disjoin(env *Env, key any) (Set, error)
		SetIter() SetIter
	}
	SetIter interface {
		HasNext(*Env) bool
		Next(*Env) (any, error)
	}

	// A Set implementation that uses a Map.
	//
	//lace:export
	MapSet struct {
		InfoHolder
		MetaHolder
		m Map
	}
	MapSetIter struct {
		ms *MapSet
		s  *SeqIterator
	}
)

func (v *MapSet) WithMeta(env *Env, meta Map) (any, error) {
	res := *v
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}

	res.meta = m
	return &res, nil
}

func (set *MapSet) Disjoin(env *Env, key any) (Set, error) {
	nm, err := set.m.Without(env, key)
	if err != nil {
		return nil, err
	}
	return &MapSet{m: nm}, nil
}

func (set *MapSet) Add(env *Env, obj any) (bool, error) {
	switch m := set.m.(type) {
	case *ArrayMap:
		return m.Add(env, obj, Boolean(true)), nil
	case *HashMap:
		if m.containsKey(env, obj) {
			return false, nil
		}
		v, err := set.m.Assoc(env, obj, Boolean(true))
		if err != nil {
			return false, err
		}
		set.m = v.(Map)
		return true, nil
	default:
		return false, nil
	}
}

func (set *MapSet) Conj(env *Env, obj any) (Conjable, error) {
	v, err := set.m.Assoc(env, obj, Boolean(true))
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
		s, err := ToString(env, v)
		if err != nil {
			return "", err
		}
		b.WriteString(s)
		if iter.HasNext(env) {
			b.WriteRune(' ')
		}
	}
	b.WriteRune('}')
	return b.String(), nil
}

type EmptySetIterator struct{}

var (
	emptySetIterator = &EmptySetIterator{}
)

func (iter *EmptySetIterator) HasNext(env *Env) bool {
	return false
}

func (iter *EmptySetIterator) Next(env *Env) (any, error) {
	panic(newIteratorError())
}

func (set *MapSet) SetIter() SetIter {
	iter := iter(set.m.Keys())

	return &MapSetIter{
		ms: set,
		s:  iter,
	}
}

func (i *MapSetIter) HasNext(env *Env) bool {
	return i.s.HasNext(env)
}

func (i *MapSetIter) Next(env *Env) (any, error) {
	k, err := i.s.Next(env)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (set *MapSet) Equals(env *Env, other interface{}) bool {
	switch otherSet := other.(type) {
	case *MapSet:
		return Equals(env, set.m, otherSet.m)
	default:
		return false
	}
}

func (set *MapSet) Get(env *Env, key any) (bool, any, error) {
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

func (set *MapSet) Hash(env *Env) (uint32, error) {
	return hashUnordered(env, set.Seq(), 2)
}

func (set *MapSet) Seq() Seq {
	return set.m.Keys()
}

func (set *MapSet) Count() int {
	return set.m.Count()
}

func (set *MapSet) Call(env *Env, args []any) (any, error) {
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
	for {
		empty, err := s.IsEmpty(env)
		if err != nil {
			return nil, err
		}
		if empty {
			break
		}
		v, err := s.First(env)
		if err != nil {
			return nil, err
		}
		_, err = res.Add(env, v)
		if err != nil {
			return nil, err
		}
		s, err = s.Rest(env)
		if err != nil {
			return nil, err
		}
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
			err = writeIndent(w, indent+2)
			if err != nil {
				return 0, err
			}
		}
	}
	fmt.Fprint(w, "}")
	return i + 1, nil
}
