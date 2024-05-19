package core

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"sort"
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
	ns, err := env.ProtoNamespace(MakeSymbol(name))
	if err != nil {
		panic(err)
	}

	return &NSBuilder{
		env: env,
		ns:  ns,
		pkg: name,
	}
}

func (b *NSBuilder) ReferCore() {
	b.ns.ReferAll(b.env.CoreNamespace, true)
}

func (b *NSBuilder) NSMeta(doc, added string) {
	m := MakeMeta(NIL, doc, added)
	b.ns.meta = m
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

	ArgTags map[string]string

	Aliases []string

	Fns []ArityFn
}

type DefTypeInfo struct {
	Name  string
	Doc   string
	Added string
	Tag   string

	Type reflect.Type

	Aliases []string
}

type DefVarInfo struct {
	Name  string
	Doc   string
	Added string
	Tag   string

	Value reflect.Value
}

func toAny(env *Env, o Object) (any, error) {
	switch sv := o.(type) {
	case Number:
		return sv.NativeNumber(), nil
	case String:
		return sv.S, nil
	case Symbol:
		return sv.Name(), nil
	case Keyword:
		return sv.Name(), nil
	case Boolean:
		return bool(sv), nil
	case Map:
		m := map[any]any{}
		i := sv.Iter()
		strKey := true
		for i.HasNext() {
			n := i.Next()

			k, err := toAny(env, n.Key)
			if err != nil {
				return nil, err
			}

			_, isStr := k.(string)
			if !isStr {
				strKey = false
			}
			v, err := toAny(env, n.Value)
			if err != nil {
				return nil, err
			}

			m[k] = v
		}

		if strKey {
			ms := map[string]any{}

			for k, v := range m {
				ms[k.(string)] = v
			}

			return ms, nil
		}

		return m, nil
	case Seq:
		var ret []any

		i := iter(sv)

		for i.HasNext(env) {
			o, err := i.Next(env)
			if err != nil {
				return nil, err
			}

			a, err := toAny(env, o)
			if err != nil {
				return nil, err
			}

			ret = append(ret, a)
		}

		return ret, nil
	case Seqable:
		return toAny(env, sv.Seq())
	case *ReflectValue:
		return sv.val, nil
	default:
		return o, nil
	}
}

// from any to any
func convertAnyIn(env *Env, index int, o Object) (reflect.Value, error) {
	if rv, ok := o.(*ReflectValue); ok {
		return rv.val, nil
	}

	a, err := toAny(env, o)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(a), nil
}

func fromAny(env *Env, v any) (Object, error) {
	switch sv := v.(type) {
	case int:
		return MakeInt(sv), nil
	case int8:
		return MakeInt(int(sv)), nil
	case int16:
		return MakeInt(int(sv)), nil
	case int32:
		return MakeInt(int(sv)), nil
	case int64:
		return MakeInt(int(sv)), nil
	case uint:
		return MakeInt(int(sv)), nil
	case uint8:
		return MakeInt(int(sv)), nil
	case uint16:
		return MakeInt(int(sv)), nil
	case uint32:
		return MakeInt(int(sv)), nil
	case uint64:
		return MakeInt(int(sv)), nil
	case string:
		return MakeString(sv), nil
	case bool:
		return MakeBoolean(sv), nil
	case []byte:
		return MakeString(string(sv)), nil
	case float32:
		return MakeDouble(float64(sv)), nil
	case float64:
		return MakeDouble(float64(sv)), nil
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Map:
			m := EmptyArrayMap()

			iter := rv.MapRange()

			for iter.Next() {
				k, err := fromAny(env, iter.Key().Interface())
				if err != nil {
					return nil, err
				}
				v, err := fromAny(env, iter.Value().Interface())
				if err != nil {
					return nil, err
				}

				m.Set(env, k, v)
			}

			return m, nil
		case reflect.Slice, reflect.Array:
			var objs []Object
			for i := 0; i < rv.Len(); i++ {
				o, err := fromAny(env, rv.Index(i).Interface())
				if err != nil {
					return nil, err
				}

				objs = append(objs, o)
			}

			return NewVectorFrom(objs...), nil
		}
		return &ReflectValue{val: reflect.ValueOf(v)}, nil
	}
}

func (n *NSBuilder) buildProc(fn any) (ProcFn, *conversionSet, error) {
	v, ok := fn.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(fn)
	}

	return convReg.ConverterForFunc(v)
}

func (nb *NSBuilder) makeMeta(b *DefnInfo) *ArrayMap {
	var m *ArrayMap

	if len(b.Fns) == 0 {
		var args []Object

		for _, n := range b.Args {
			if t, ok := b.ArgTags[n]; ok {
				args = append(args, MakeTaggedSymbol(n, MakeSymbol(t)))
			} else {
				args = append(args, MakeSymbol(n))
			}
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
		m = m.Plus(nb.env, MakeKeyword("tag"), MakeString(b.Tag))
	}

	return m
}

func (b *NSBuilder) Def(name string, obj Object) {
	m := MakeMeta(
		NewListFrom(),
		"", "x",
	)
	_, err := b.ns.InternVar(b.env, name, obj, m)
	if err != nil {
		panic(err)
	}
}

func (b *NSBuilder) DefType(i *DefTypeInfo) *NSBuilder {
	m := MakeMeta(
		NewListFrom(),
		i.Doc, i.Added,
	)

	obj := &ReflectType{typ: i.Type}

	_, err := b.ns.InternVar(b.env, i.Name, obj, m)
	if err != nil {
		panic(err)
	}
	for _, a := range i.Aliases {
		_, err = b.ns.InternVar(b.env, a, obj, m)
		if err != nil {
			panic(err)
		}
	}

	return b
}

func (b *NSBuilder) DefVar(i *DefVarInfo) *NSBuilder {
	m := MakeMeta(
		NewListFrom(),
		i.Doc, i.Added,
	)

	var obj Object
	switch i.Value.Kind() {
	case reflect.String:
		obj = MakeString(i.Value.String())
	case reflect.Int:
		obj = MakeInt(int(i.Value.Int()))
	case reflect.Bool:
		obj = MakeBoolean(i.Value.Bool())
	default:
		var ok bool
		obj, ok = i.Value.Interface().(Object)
		if !ok {
			obj = &ReflectValue{val: i.Value}
		}
	}

	_, err := b.ns.InternVar(b.env, i.Name, obj, m)
	if err != nil {
		panic(err)
	}

	return b
}

func (n *NSBuilder) Defn(b *DefnInfo) *NSBuilder {
	if b.Fn != nil {
		procFn, _, err := n.buildProc(b.Fn)
		if err != nil {
			panic(fmt.Sprintf("unable to define %s: %s", b.Name, err))
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

		if file == "" {
			_, file, line, _ = runtime.Caller(1)
		}

		p := Proc{
			Fn:      procFn,
			Name:    b.Name,
			Package: n.pkg,
			File:    file,
			Line:    line,
		}

		meta := n.makeMeta(b)

		_, err = n.ns.InternVar(n.env, b.Name, p, meta)
		if err != nil {
			panic(err)
		}
		for _, a := range b.Aliases {
			_, err = n.ns.InternVar(n.env, a, p, meta)
			if err != nil {
				panic(err)
			}
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
		pf, cs, err := n.buildProc(f.Fn)
		if err != nil {
			panic(err)
		}

		for _, pf := range fns {
			if pf.arity == cs.arity {
				panic("supplied two functions with same arity")
			}
		}

		if cs == nil {
			fns = append(fns, fn{-1, pf})
		}

		fns = append(fns, fn{cs.arity, pf})
	}

	sort.Slice(fns, func(i, j int) bool {
		return fns[i].arity < fns[j].arity
	})

	var procFn ProcFn

	if len(fns) == 1 || fns[0].arity == -1 {
		procFn = fns[0].proc
	} else {
		procFn = ProcFn(func(env *Env, args []Object) (Object, error) {
			for _, fn := range fns {
				if fn.arity == len(args) {
					return fn.proc(env, args)
				}
			}

			return nil, ErrorArityMinMax(env, len(args), fns[0].arity, fns[len(fns)-1].arity)
		})
	}

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

	if file == "" {
		_, file, line, _ = runtime.Caller(1)
	}

	p := Proc{
		Fn:      procFn,
		Name:    b.Name,
		Package: n.pkg,
		File:    file,
		Line:    line,
	}

	meta := n.makeMeta(b)

	_, err := n.ns.InternVar(n.env, b.Name, p, meta)
	if err != nil {
		panic(err)
	}
	for _, a := range b.Aliases {
		_, err = n.ns.InternVar(n.env, a, p, meta)
		panic(err)
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
	b.ns.ReferAll(b.env.CoreNamespace, true)
	_, err := ProcessReader(b.env, reader, filename)
	if err != nil {
		return err
	}

	return nil
}
