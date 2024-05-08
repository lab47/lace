package core

import (
	"fmt"
)

// Create a new lace List from the given arguments
//
//lace:export List
func MakeList(env *Env, args Args) (Object, error) {
	l := NewListFrom(args.Objects...)

	show(env, "list", l)

	return l, nil
}

var cnt = 0

func show(env *Env, prefix string, obj Object) {
	return
	os, _ := obj.ToString(env, true)
	fmt.Printf("%s: %s\n", prefix, os)
}

// Add an element to a Seq value, returning a new Seq
//
//lace:export
func Cons(env *Env, val Object, seq Seqable) (Object, error) {
	s := seq.Seq()

	show(env, "cons<-", val)
	show(env, "cons", s)

	//ss, _ := s.ToString(env, true)
	//os, _ := val.ToString(env, true)

	//fmt.Printf("%d cons %s %s\n", cnt, os, ss)
	cnt++

	return s.Cons(val), nil
}

// Return the first element in a Seq
//
//lace:export
func First(env *Env, s Seqable) (Object, error) {
	q := s.Seq()
	show(env, "first", q)
	return q.First(env)
}

// Return elements other than the first one in a Seq
//
//lace:export
func Next(env *Env, s Seqable) (Object, error) {
	q := s.Seq()
	show(env, "next", q)

	res, err := q.Rest(env)
	if err != nil {
		return nil, err
	}
	empty, err := res.IsEmpty(env)
	if err != nil {
		return nil, err
	}

	if empty {
		return NIL, nil
	}
	return res, nil
}

// Return all elements of a seq except for the first one.
//
//lace:export
func Rest(env *Env, s Seqable) (Object, error) {
	q := s.Seq()
	show(env, "rest", q)

	return q.Rest(env)
}

// Create a new Sequence by combine the value with the collection.
//
//lace:export
func Conj(env *Env, col Object, val Object) (Object, error) {
	show(env, "conj", col)
	show(env, "conj<-", val)

	switch c := col.(type) {
	case Conjable:
		return c.Conj(env, val)
	case Seq:
		return c.Cons(val), nil
	default:
		return nil, env.RT.NewError("conj's first argument must be a collection, got " + c.GetType().Name())
	}
}

// Convert the given value to a Seq
//
//lace:export Seq
func ConvertToSeq(env *Env, s Seqable) (Object, error) {
	sq := s.Seq()
	show(env, "seq", sq)
	empty, err := sq.IsEmpty(env)
	if err != nil {
		return nil, err
	}

	if empty {
		return NIL, nil
	}

	return sq, nil
}

// Concatinate N sequences together
//
//lace:export
func ConcatSimple(env *Env, args Args) (Object, error) {
	var data []Object

	for _, o := range args.Objects {
		if s, ok := o.(Seqable); ok {
			eles, err := ToSlice(env, s.Seq())
			if err != nil {
				return nil, err
			}

			data = append(data, eles...)
		}
	}

	l := NewListFrom(data...)

	show(env, "cs-out", l)

	return l, nil
}

// Compare two values returning a boolean if they are equal or not
//
//lace:export
func Equals(env *Env, a, b Object) (Object, error) {
	return MakeBoolean(a.Equals(env, b)), nil
}

// Add given bindings to the set of current Var bindings, returning
// the original set.
//
//lace:export
func PushBindings(env *Env, assoc Map) (Object, error) {
	cur := env.CurrentVar
	if cur == nil {
		cur = EmptyArrayMap()
	}

	orig := cur

	mi := assoc.Iter()

	var err error

	for mi.HasNext() {
		pair := mi.Next()

		cur, err = cur.Assoc(env, pair.Key, pair.Value)
		if err != nil {
			return nil, err
		}
	}

	return orig, nil
}

// Reset the local var bindings to the given value.
//
//lace:export
func SetBindings(env *Env, assoc Associative) (Object, error) {
	env.CurrentVar = assoc
	return assoc, nil
}
