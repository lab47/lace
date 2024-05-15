package core

import (
	"fmt"
	"os"
	"reflect"
)

type ProcFn func(env *Env, args []Object) (Object, error)
type Proc struct {
	Fn      ProcFn
	Name    string
	Package string // "" for core (this package), else e.g. "std/string"
	File    string
	Line    int
}

var _ Callable = Proc{}

func (p Proc) Compare(env *Env, a, b Object) (int, error) {
	return compare(env, p, a, b)
}

func (p Proc) ToString(env *Env, escape bool) (string, error) {
	pkg := p.Package
	if pkg != "" {
		pkg += "."
	}

	file := p.File
	if file == "" {
		file = "<unknown>"
	}

	return fmt.Sprintf("#Proc[%s%s %s:%d]", pkg, p.Name, file, p.Line), nil
}

func (p Proc) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case Proc:
		return reflect.ValueOf(p.Fn).Pointer() == reflect.ValueOf(other.Fn).Pointer()
	}
	return false
}

func (p Proc) GetInfo() *ObjectInfo {
	return nil
}

func (p Proc) WithInfo(*ObjectInfo) Object {
	return p
}

func (p Proc) GetType() *Type {
	return TYPE.Proc
}

func (p Proc) Hash(env *Env) (uint32, error) {
	return HashPtr(&p.Fn), nil
}

func (p Proc) Call(env *Env, args []Object) (Object, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr,
				"\nPanic from proc: %s at %s:%d\nerror: %s\n\n",
				p.Name, p.File, p.Line, err,
			)

			panic(err)
		}
	}()
	ret, err := p.Fn(env, args)
	if err != nil {
		err = env.populateStackTrace(err)
	}

	return ret, err
}
