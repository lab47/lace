package lang

import (
	"fmt"

	"github.com/lab47/lace/core"
)

// Create a new lace List from the given arguments
func List(env *core.Env, args core.Args) (core.Object, error) {
	l := core.NewListFrom(args.Objects...)

	show(env, "list", l)

	return l, nil
}

var cnt = 0

func show(env *core.Env, prefix string, obj core.Object) {
	return
	os, _ := obj.ToString(env, true)
	fmt.Printf("%s: %s\n", prefix, os)
}

// Add an element to a Seq value, returning a new Seq
func Cons(env *core.Env, val core.Object, seq core.Seqable) (core.Object, error) {
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
func First(env *core.Env, s core.Seqable) (core.Object, error) {
	q := s.Seq()
	show(env, "first", q)
	return q.First(env)
}

// Return elements other than the first one in a Seq
func Next(env *core.Env, s core.Seqable) (core.Object, error) {
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
		return core.NIL, nil
	}
	return res, nil
}

func Rest(env *core.Env, s core.Seqable) (core.Object, error) {
	q := s.Seq()
	show(env, "rest", q)

	return q.Rest(env)
}

func Conj(env *core.Env, col core.Object, val core.Object) (core.Object, error) {
	show(env, "conj", col)
	show(env, "conj<-", val)

	switch c := col.(type) {
	case core.Conjable:
		return c.Conj(env, val)
	case core.Seq:
		return c.Cons(val), nil
	default:
		return nil, env.RT.NewError("conj's first argument must be a collection, got " + c.GetType().Name())
	}
}

func Seq(env *core.Env, s core.Seqable) (core.Object, error) {
	sq := s.Seq()
	show(env, "seq", sq)
	empty, err := sq.IsEmpty(env)
	if err != nil {
		return nil, err
	}

	if empty {
		return core.NIL, nil
	}

	return sq, nil
}

func ConcatSimple(env *core.Env, args core.Args) (core.Object, error) {
	var data []core.Object

	for _, o := range args.Objects {
		if s, ok := o.(core.Seqable); ok {
			eles, err := core.ToSlice(env, s.Seq())
			if err != nil {
				return nil, err
			}

			data = append(data, eles...)
		}
	}

	l := core.NewListFrom(data...)

	show(env, "cs-out", l)

	return l, nil
}

// Compare two values returning a boolean if they are equal or not
func Equals(env *core.Env, a, b core.Object) (core.Object, error) {
	return core.MakeBoolean(a.Equals(env, b)), nil
}
