package core

import (
	"bytes"
	"fmt"
	"io"
)

type (
	Map interface {
		Associative
		Seqable
		Counted
		Without(env *Env, key Object) (Map, error)
		Keys() Seq
		Vals() Seq
		Merge(env *Env, m Map) (Map, error)
		Iter() MapIterator
		GetEqu(key Equ) (bool, Object)
	}
	MapIterator interface {
		HasNext() bool
		Next() *Pair
	}
	EmptyMapIterator struct {
	}
	Pair struct {
		Key   Object
		Value Object
	}
)

var (
	emptyMapIterator = &EmptyMapIterator{}
)

func (iter *EmptyMapIterator) HasNext() bool {
	return false
}

func (iter *EmptyMapIterator) Next() *Pair {
	panic(newIteratorError())
}

func mapConj(env *Env, m Map, obj Object) (Conjable, error) {
	switch obj := obj.(type) {
	case *Vector:
		if obj.count != 2 {
			return nil, StubNewError("Vector argument to map's conj must be a vector with two elements")
		}
		return m.Assoc(env, obj.at(0), obj.at(1))
	case Map:
		return m.Merge(env, obj)
	default:
		return nil, StubNewError("Argument to map's conj must be a vector with two elements or a map")
	}
}

func mapEquals(env *Env, m Map, other interface{}) bool {
	if m == other {
		return true
	}
	switch otherMap := other.(type) {
	case Nil:
		return false
	case Map:
		defer env.enableCycleDetection()()

		if env.cycling(m, otherMap) {
			return true
		}

		if m.Count() != otherMap.Count() {
			return false
		}
		for iter := m.Iter(); iter.HasNext(); {
			p := iter.Next()
			success, value, err := otherMap.Get(env, p.Key)
			if err != nil {
				return false
			}
			if !success || !Equals(env, value, p.Value) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func mapToString(env *Env, m Map, escape bool) (string, error) {
	env.enableCycleDetection()()

	if env.cycling(m, NIL) {
		return "{...}", nil
	}

	var b bytes.Buffer
	b.WriteRune('{')
	if m.Count() > 0 {
		for iter := m.Iter(); ; {
			p := iter.Next()
			ks, err := ToString(env, p.Key)
			if err != nil {
				return "", err
			}

			vs, err := ToString(env, p.Value)
			if err != nil {
				return "", err
			}

			b.WriteString(ks)
			b.WriteRune(' ')
			b.WriteString(vs)
			if iter.HasNext() {
				b.WriteString(", ")
			} else {
				break
			}
		}
	}
	b.WriteRune('}')
	return b.String(), nil
}

func callMap(env *Env, m Map, args []Object) (Object, error) {
	if err := CheckArity(env, args, 1, 2); err != nil {
		return nil, err
	}

	ok, v, err := m.Get(env, args[0])
	if err != nil {
		return nil, err
	}

	if ok {
		return v, nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return NIL, nil
}

func pprintMap(env *Env, m Map, w io.Writer, indent int) (int, error) {
	i := indent + 1
	fmt.Fprint(w, "{")
	var err error
	if m.Count() > 0 {
		for iter := m.Iter(); ; {
			p := iter.Next()
			i, err = pprintObject(env, p.Key, indent+1, w)
			if err != nil {
				return 0, err
			}

			fmt.Fprint(w, " ")
			i, err = pprintObject(env, p.Value, i+1, w)
			if err != nil {
				return 0, err
			}
			if iter.HasNext() {
				fmt.Fprint(w, ",\n")
				err = writeIndent(w, indent+1)
				if err != nil {
					return 0, err
				}
			} else {
				break
			}
		}
	}
	fmt.Fprint(w, "}")
	return i + 1, nil
}
