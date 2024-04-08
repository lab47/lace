package core

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"sort"

	"github.com/davecgh/go-spew/spew"
)

type NativeSetup func(env *Env) error

var NativeRegistry = map[string]NativeSetup{}

func AddNativeNamespace(name string, setup NativeSetup) {
	NativeRegistry[name] = setup
}

func PopulateNativeNamespacesToEnv(env *Env) error {
	for _, setup := range NativeRegistry {
		err := setup(env)
		if err != nil {
			return err
		}
	}

	return nil
}

func PopulateNativeNamespaceToEnv(env *Env, name string) (bool, error) {
	setup, ok := NativeRegistry[name]
	if !ok {
		return false, nil
	}

	err := setup(env)
	if err != nil {
		return false, err
	}

	return true, nil
}

type NSBuilder struct {
	env *Env
	ns  *Namespace
	pkg string
}

func NewNSBuilder(env *Env, name string) *NSBuilder {
	ns := env.EnsureNamespace(MakeSymbol(name))

	return &NSBuilder{
		env: env,
		ns:  ns,
		pkg: name,
	}
}

func (b *NSBuilder) Namespace() *Namespace {
	return b.ns
}

type DefnTag struct {
	Name string
	Tag  string
}

type ArityFn struct {
	Args []string
	Fn   any
}

type DefnInfo struct {
	Name  string
	Args  []string
	Rest  bool
	Doc   string
	Added string
	Tag   string
	Tags  []DefnTag
	Fn    any

	Aliases []string

	Fns []ArityFn
}

// from error to Error
func convertErrorOut(env *Env, s reflect.Value) (reflect.Value, error) {
	err, ok := s.Interface().(error)
	if !ok {
		return reflect.Value{}, env.RT.NewError("wrong type, expecting error")
	}

	return reflect.ValueOf(env.RT.NewError(err.Error())), nil
}

// from string to String
func convertStringOut(env *Env, s reflect.Value) (reflect.Value, error) {
	gs, ok := s.Interface().(string)
	if !ok {
		return reflect.Value{}, env.RT.NewError("wrong type, expecting string")
	}

	return reflect.ValueOf(MakeString(gs)), nil
}

// from []byte to String
func convertBytesOut(env *Env, s reflect.Value) (reflect.Value, error) {
	gs, ok := s.Interface().([]byte)
	if !ok {
		return reflect.Value{}, env.RT.NewError("wrong type, expecting string")
	}

	return reflect.ValueOf(MakeString(string(gs))), nil
}

// from *regexp.Regexp to *Regex
func convertRegexpOut(env *Env, s reflect.Value) (reflect.Value, error) {
	gs, ok := s.Interface().(*regexp.Regexp)
	if !ok {
		return reflect.Value{}, env.RT.NewError("wrong type, expecting string")
	}

	return reflect.ValueOf(MakeRegex(gs)), nil
}

// from Object to Object
func convertObjectOut(env *Env, s reflect.Value) (reflect.Value, error) {
	gs, ok := s.Interface().(Object)
	if !ok {
		return reflect.Value{}, env.RT.NewError("wrong type, expecting Object")
	}

	return reflect.ValueOf(gs), nil
}

// from bool to Boolean
func convertBoolOut(env *Env, s reflect.Value) (reflect.Value, error) {
	gb, ok := s.Interface().(bool)
	if !ok {
		return reflect.Value{}, env.RT.NewError("wrong type, expecting string")
	}

	return reflect.ValueOf(MakeBoolean(gb)), nil
}

// from String to string
func convertStringIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(String)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "String")
	}

	return reflect.ValueOf(ls.S), nil
}

// from Symbol to Symbol
func convertSymbolIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Symbol)
	if !ok {
		spew.Dump(o)
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Symbol")
	}

	return reflect.ValueOf(ls), nil
}

// from Symbol to Symbol
func convertKeywordIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Keyword)
	if !ok {
		spew.Dump(o)
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Symbol")
	}

	return reflect.ValueOf(ls), nil
}

// from String to []byte
func convertBytesIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(String)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "String")
	}

	return reflect.ValueOf([]byte(ls.S)), nil
}

// from Int to int
func convertIntIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Int)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Int")
	}

	return reflect.ValueOf(ls.Int()), nil
}

// from Object to Object
func convertObjectIn(env *Env, index int, o Object) (reflect.Value, error) {
	return reflect.ValueOf(o), nil
}

// from Env* to Env*
func convertEnvIn(env *Env, index int, o Object) (reflect.Value, error) {
	return reflect.ValueOf(env), nil
}

// from Seqable to Seqable
func convertSeqableIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Seqable)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Seqable")
	}

	return reflect.ValueOf(ls), nil
}

func convertCallableIn(env *Env, index int, o Object) (reflect.Value, error) {
	ls, ok := o.(Callable)
	if !ok {
		return reflect.Value{}, env.RT.NewArgTypeError(index, o, "Seqable")
	}

	return reflect.ValueOf(ls), nil
}

type inConv func(*Env, int, Object) (reflect.Value, error)
type outConv func(*Env, reflect.Value) (reflect.Value, error)

var procFnType = reflect.TypeFor[ProcFn]()

func (n *NSBuilder) buildProc(fn any) (ProcFn, int, error) {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		panic("procs can only be built from Go funcs")
	}

	t := v.Type()

	argIn := make([]inConv, t.NumIn())

	var passed int
	var envIn bool

	for i := 0; i < t.NumIn(); i++ {
		at := t.In(i)

		switch at {
		case reflect.TypeFor[*Env]():
			envIn = true
			argIn[i] = convertEnvIn
			continue
		case reflect.TypeFor[Object]():
			argIn[i] = convertObjectIn
		case reflect.TypeFor[Seqable]():
			argIn[i] = convertSeqableIn
		case reflect.TypeFor[Callable]():
			argIn[i] = convertCallableIn
		case reflect.TypeFor[string]():
			argIn[i] = convertStringIn
		case reflect.TypeFor[Symbol]():
			argIn[i] = convertSymbolIn
		case reflect.TypeFor[Keyword]():
			argIn[i] = convertKeywordIn
		case reflect.TypeFor[[]byte]():
			argIn[i] = convertBytesIn
		case reflect.TypeFor[int]():
			argIn[i] = convertIntIn
		default:
			return nil, 0, fmt.Errorf("unable to find converter for %s", at)
		}

		passed++
	}

	rets := make([]outConv, t.NumOut())

	var valueReturns int

	errorPos := -1

	for i := 0; i < t.NumOut(); i++ {
		at := t.Out(i)

		switch at {
		case reflect.TypeFor[Object]():
			rets[i] = convertObjectOut
		case reflect.TypeFor[string]():
			rets[i] = convertStringOut
		case reflect.TypeFor[[]byte]():
			rets[i] = convertBytesOut
		case reflect.TypeFor[bool]():
			rets[i] = convertBoolOut
		case reflect.TypeFor[*regexp.Regexp]():
			rets[i] = convertRegexpOut
		case reflect.TypeFor[error]():
			rets[i] = convertErrorOut
			errorPos = i
			continue
		default:
			return nil, 0, fmt.Errorf("unable to find converter for %s", at)
		}

		valueReturns++
	}

	nilErr := reflect.Zero(reflect.TypeFor[error]())
	nilObject := reflect.Zero(reflect.TypeFor[Object]())

	fnVal := reflect.ValueOf(fn)

	var out reflect.Value

	if t.NumIn() == 1 && t.NumOut() == 1 {
		out = reflect.MakeFunc(reflect.TypeFor[ProcFn](), func(args []reflect.Value) (results []reflect.Value) {
			env := args[0].Interface().(*Env)
			objArgs := args[1].Interface().([]Object)

			if len(objArgs) != 1 {
				return []reflect.Value{nilObject, reflect.ValueOf(ErrorArity(env, 1))}
			}

			arg, err := argIn[0](env, 0, objArgs[0])
			if err != nil {
				return []reflect.Value{reflect.Zero(reflect.TypeFor[Object]()), reflect.ValueOf(err)}
			}

			ret := fnVal.Call([]reflect.Value{arg})

			convRet, err := rets[0](env, ret[0])
			if err != nil {
				return []reflect.Value{reflect.Zero(reflect.TypeFor[Object]()), reflect.ValueOf(err)}
			}

			return []reflect.Value{convRet, nilErr}
		})
	} else {
		out = reflect.MakeFunc(reflect.TypeFor[ProcFn](), func(args []reflect.Value) (results []reflect.Value) {
			env := args[0].Interface().(*Env)
			objArgs := args[1].Interface().([]Object)

			if len(objArgs) != passed {
				return []reflect.Value{nilObject, reflect.ValueOf(ErrorArity(env, passed))}
			}

			dest := make([]reflect.Value, t.NumIn())

			var destOffset int
			if envIn {
				dest[0] = reflect.ValueOf(env)
				destOffset = 1
			}

			for i, a := range objArgs {
				sub, err := argIn[destOffset+i](env, i, a)
				if err != nil {
					return []reflect.Value{reflect.Zero(reflect.TypeFor[Object]()), reflect.ValueOf(err)}
				}

				dest[destOffset+i] = sub
			}

			output := make([]reflect.Value, 2)

			ret := fnVal.Call(dest)

			if errorPos >= 0 {
				output[1] = ret[errorPos]
			} else {
				output[1] = nilErr
			}

			switch valueReturns {
			case 0:
				output[0] = reflect.ValueOf(NIL)
			case 1:
				sub, err := rets[0](env, ret[0])
				if err != nil {
					return []reflect.Value{reflect.Zero(reflect.TypeFor[Object]()), reflect.ValueOf(err)}
				}
				output[0] = sub
			default:
				var objects []Object

				for i, rv := range ret {
					if i == errorPos {
						continue
					}

					v, err := rets[i](env, rv)
					if err != nil {
						if o, ok := v.Interface().(Object); ok {
							objects = append(objects, o)
						}
					}
				}

				output[0] = reflect.ValueOf(NewListFrom(objects...))
			}

			return output
		})
	}

	return out.Interface().(ProcFn), passed, nil
}

func (n *NSBuilder) makeMeta(b *DefnInfo) *ArrayMap {
	var m *ArrayMap

	if len(b.Fns) == 0 {
		var args []Object

		for _, n := range b.Args {
			args = append(args, MakeSymbol(n))
		}

		m = MakeMeta(
			NewListFrom(NewVectorFrom(args...)),
			b.Doc, b.Added,
		)
	} else {
		var vecs []Object

		for _, fn := range b.Fns {
			var args []Object

			for _, n := range fn.Args {
				args = append(args, MakeSymbol(n))
			}

			vecs = append(vecs, NewVectorFrom(args...))
		}

		m = MakeMeta(
			NewListFrom(vecs...),
			b.Doc, b.Added,
		)
	}

	if b.Tag != "" {
		m = m.Plus(MakeKeyword("tag"), MakeString(b.Tag))
	}

	return m
}

func (b *NSBuilder) Def(name string, obj Object) {
	m := MakeMeta(
		NewListFrom(),
		"", "x",
	)
	b.ns.InternVar(name, obj, m)
}

func (n *NSBuilder) Defn(b *DefnInfo) *NSBuilder {
	if b.Fn != nil {
		procFn, _, err := n.buildProc(b.Fn)
		if err != nil {
			panic(err)
		}

		var (
			file string
			line int
		)

		rv := reflect.ValueOf(b.Fn)
		if rv.Kind() == reflect.Func {
			rf := runtime.FuncForPC(rv.Pointer())
			file, line = rf.FileLine(rv.Pointer())
		}

		p := Proc{
			Fn:      procFn,
			Name:    b.Name,
			Package: n.pkg,
			File:    file,
			Line:    line,
		}

		meta := n.makeMeta(b)

		n.ns.InternVar(b.Name, p, meta)
		for _, a := range b.Aliases {
			n.ns.InternVar(a, p, meta)
		}

		return n
	}

	if b.Fns == nil {
		panic("didn't provide fn or fns: " + b.Name)
	}

	type fn struct {
		arity int
		proc  ProcFn
	}

	var fns []fn

	for _, f := range b.Fns {
		pf, arity, err := n.buildProc(f.Fn)
		if err != nil {
			panic(err)
		}

		for _, pf := range fns {
			if pf.arity == arity {
				panic("supplied two functions with same arity")
			}
		}

		fns = append(fns, fn{arity, pf})
	}

	sort.Slice(fns, func(i, j int) bool {
		return fns[i].arity < fns[j].arity
	})

	dispatch := reflect.MakeFunc(reflect.TypeFor[ProcFn](), func(args []reflect.Value) (results []reflect.Value) {
		env := args[0].Interface().(*Env)
		objArgs := args[1].Interface().([]Object)

		for _, fn := range fns {
			if fn.arity == len(objArgs) {
				ret, err := fn.proc(env, objArgs)
				return []reflect.Value{reflect.ValueOf(ret), reflect.ValueOf(err)}
			}
		}

		return []reflect.Value{
			reflect.Zero(reflect.TypeFor[Object]()),
			reflect.ValueOf(ErrorArityMinMax(env, len(objArgs), fns[0].arity, fns[len(fns)-1].arity)),
		}
	})

	procFn := dispatch.Interface().(ProcFn)

	// TODO we can do better than reporting the first function only as the location.
	var (
		file string
		line int
	)

	rv := reflect.ValueOf(b.Fns[0])
	if rv.Kind() == reflect.Func {
		rf := runtime.FuncForPC(rv.Pointer())
		file, line = rf.FileLine(rv.Pointer())
	}

	p := Proc{
		Fn:      procFn,
		Name:    b.Name,
		Package: n.pkg,
		File:    file,
		Line:    line,
	}

	meta := n.makeMeta(b)

	n.ns.InternVar(b.Name, p, meta)
	for _, a := range b.Aliases {
		n.ns.InternVar(a, p, meta)
	}

	return n
}

func (b *NSBuilder) Run(code []byte) error {
	filename := fmt.Sprintf("<%s>", b.pkg)
	reader := NewReader(bytes.NewReader(code), filename)
	cur := b.env.CurrentNamespace()
	defer func() {
		b.env.SetCurrentNamespace(cur)
	}()
	b.env.SetCurrentNamespace(b.ns)
	b.ns.ReferAll(b.env.CoreNamespace)
	err := ProcessReader(b.env, reader, filename, EVAL)
	if err != nil {
		return err
	}

	return nil
}
