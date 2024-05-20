package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//go:generate go run .././pkg/pkgreflect/cmd/pkgreflect -lace-name lace.lang -honor-directive -in-core -specialized github.com/lab47/lace/core binding.go

// Create a new lace List from the given arguments
//
//lace:export List
func MakeList(env *Env, args []Object) (Object, error) {
	l := NewListFrom(args...)

	show(env, "list", l)

	return l, nil
}

var cnt = 0

func show(env *Env, prefix string, obj Object) {
	return
	os, _ := ToString(env, obj)
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
		return nil, env.NewError("conj's first argument must be a collection, got " + TypeName(c))
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
func ConcatSimple(env *Env, args []Object) (Object, error) {
	var data []Object

	for _, o := range args {
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
//lace:export Equals
func EqualsValues(env *Env, a, b Object) (Object, error) {
	return MakeBoolean(Equals(env, a, b)), nil
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

// Attempt to load a given lib from a given path.
//
//lace:export
func LoadLibFromPath(env *Env, libnamev Symbol, pathnamev String) (Object, error) {
	// Sometimes we load namespaces without telling the clojure code,
	// so see if it's already loaded and if so, use it.

	if env.FindNamespace(libnamev) != nil {
		return NIL, nil
	}

	libname := libnamev.Name()
	pathname := pathnamev.S()

	cp := env.classPath.GetStatic()
	cpvec, err := AssertVector(env, cp, "*classpath* must be a Vector, not a "+TypeName(cp))
	if err != nil {
		return nil, err
	}

	count := cpvec.Count()
	var f *os.File
	var canonicalErr error
	var filename string
	for i := 0; i < count; i++ {
		elem := cpvec.at(i)
		cpelem, err := AssertString(env, elem, "*classpath* must contain only Strings, not a "+TypeName(elem)+" (at element "+strconv.Itoa(i)+")")
		if err != nil {
			return nil, err
		}
		s := cpelem.S()
		if s == "" {
			filename = pathname
		} else {
			filename = filepath.Join(s, filepath.Join(strings.Split(libname, ".")...)) + ".clj" // could cache inner join....
		}

		f, err = os.Open(filename)
		if err == nil {
			canonicalErr = nil
			break
		}
		if s == "" {
			canonicalErr = err
		}
	}
	if canonicalErr != nil {
		return nil, canonicalErr
	}
	if f == nil || filename == "" {
		return nil, SError(env, "LoadError", "unable to find path for library", "library", libname)
	}
	reader := NewReader(bufio.NewReader(f), filename)
	err = ProcessReaderFromEval(env, reader, filename)
	if err != nil {
		return nil, SError(env, "LoadError", "error loading file", "path", filename, "error", err.Error())
	}
	return NIL, nil
}

// StartGoRoutine runs the given callable in a new goroutine, returning a channel
// that can be used to retrieve the return value.
//
//lace:export
func StartGoRoutine(parent *Env, callable Callable) (Object, error) {
	ch := MakeChannel(make(chan FutureResult, 1))
	env := parent.Child()
	go func() {
		var cerr Error
		res, err := callable.Call(env, []Object{})
		if err != nil {
			cerr = env.NewError(err.Error())
		}

		if cerr != nil {
			DisplayError(env, cerr)
		}

		ch.ch <- MakeFutureResult(res, cerr)
		ch.Close()
	}()
	return ch, nil
}
