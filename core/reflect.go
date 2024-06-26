package core

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	_ "embed"

	"github.com/lab47/lace/pkg/pkgreflect"
)

func (r *Type) Call(env *Env, args []any) (any, error) {
	rv := reflect.New(r.rType)
	return rv.Interface(), nil
}

func CastReflect[T any](env *Env, obj any, v *T) error {
	if x, ok := obj.(T); ok {
		*v = x
		return nil
	}

	return fmt.Errorf("not a reflect value containing a *%T", v)
}

func structSetInPlace(env *Env, val reflect.Value, name string, fval any) error {
	val = reflect.Indirect(val)

	if val.Kind() != reflect.Struct {
		return env.NewError(fmt.Sprintf("value must be a struct, is a %T", val.Interface()))
	}

	field := val.FieldByName(name)
	if !field.IsValid() {
		return env.NewError(fmt.Sprintf("unknown struct field %s", name))
	}

	var (
		frv reflect.Value
		err error
	)

	if field.Type().Kind() == reflect.Func {
		call, ok := fval.(Callable)
		if !ok {
			return env.TypeError(TCContext{Context: "struct value"}, fval, "Callable")
		}
		frv = convReg.makeFuncConvertIn(env, call, field.Type())
	} else {
		cv, _ := convReg.convArg(field.Type())

		frv, err = cv(env, -1, fval)
		if err != nil {
			return err
		}
	}

	if !frv.Type().AssignableTo(field.Type()) {
		return env.NewError(
			fmt.Sprintf("needed type %s, had %T", field.Type(), fval))
	}

	field.Set(frv)

	return nil
}

func structGet(env *Env, obj any, name string) (any, error) {
	val := reflect.ValueOf(obj)

	val = reflect.Indirect(val)

	if val.Kind() != reflect.Struct {
		return nil, env.NewError(fmt.Sprintf("value must be a struct, is a %T", val.Interface()))
	}

	field := val.FieldByName(name)
	if !field.IsValid() {
		return nil, env.NewError(fmt.Sprintf("unknown struct field %s", name))
	}

	rt := field.Type()

	obj, err := convReg.convRet(rt)(field)
	if err != nil {
		if ce, ok := err.(OutConvError); ok {
			return nil, env.NewError(string(ce))
		}
	}

	return obj, nil
}

func structList(env *Env, val reflect.Value) (any, error) {
	val = reflect.Indirect(val)

	if val.Kind() != reflect.Struct {
		return nil, env.NewError(fmt.Sprintf("value must be a struct, is a %T", val.Interface()))
	}

	var ret []any

	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		ret = append(ret, MakeKeyword(t.Field(i).Name))
	}

	return NewListFrom(ret...), nil
}

func structFromMap(env *Env, rt *Type, m Map) (any, error) {
	t := rt.ReflectType()
	if t.Kind() != reflect.Struct {
		return nil, env.NewError("type is not a struct")
	}

	ret := reflect.New(t)

	rv := ret.Elem()

	iter := m.Iter()

	for iter.HasNext() {
		p := iter.Next()

		var key string

		if err := CoerceString(env, p.Key, &key); err != nil {
			return nil, err
		}

		key = convertMethodName(key)

		field := rv.FieldByName(key)
		if !field.IsValid() {
			return nil, env.NewError("Unknown field: %s", key)
		}

		fval := p.Value

		var (
			frv reflect.Value
			err error
		)

		if field.Type().Kind() == reflect.Func {
			call, ok := fval.(Callable)
			if !ok {
				return nil, env.TypeError(TCContext{Context: "struct value"}, fval, "Callable")
			}
			frv = convReg.makeFuncConvertIn(env, call, field.Type())
		} else {
			cv, _ := convReg.convArg(field.Type())

			frv, err = cv(env, -1, fval)
			if err != nil {
				return nil, err
			}
		}

		if !frv.Type().AssignableTo(field.Type()) {
			return nil, env.NewError(
				fmt.Sprintf("needed type %s, had %T", field.Type(), fval))
		}

		field.Set(frv)
	}

	return ret.Interface(), nil
}

func structAsMap(env *Env, val reflect.Value) (any, error) {
	val = reflect.Indirect(val)
	if val.Kind() != reflect.Struct {
		return nil, env.NewError("Value is not a struct, is kind: %s", val.Kind())
	}

	return StructMap{val: val}, nil
}

func structToMap(env *Env, val reflect.Value) (any, error) {
	for val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, env.NewError(fmt.Sprintf("value must be a struct, is a %s (%T)", val.Kind(), val.Interface()))
	}

	var ret Associative = NIL

	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		if !tf.IsExported() {
			continue
		}

		field := val.Field(i)
		rt := field.Type()

		cv, err := convReg.convRet(rt)(field)
		if err != nil {
			if ce, ok := err.(OutConvError); ok {
				return nil, env.NewError(string(ce))
			}
			return nil, err
		}

		ret, err = ret.Assoc(env,
			MakeKeyword(tf.Name),
			cv,
		)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

type reifiedFunc struct {
	pkgreflect.Func
	Name Symbol
}

type reifiedType struct {
	Namespace string
	Methods   map[string]reifiedFunc
	MethodVec *Vector
}

func listMethods(env *Env, obj any, reg map[reflect.Type]reifiedType) Seq {
	t := reflect.TypeOf(obj)

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	var objs []any
	if meths, ok := reg[t]; ok {
		return meths.MethodVec.Seq()
	} else {
		for i := 0; i < t.NumMethod(); i++ {
			objs = append(objs, MakeKeyword(t.Method(i).Name))
		}
	}

	return NewListFrom(objs...)
}

func castObjectToRef(env *Env, typ reflect.Type, obj any) (any, error) {
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := AssertNumber(env, obj, "")
		if err != nil {
			return nil, err
		}

		v := reflect.New(typ).Elem()
		v.SetInt(int64(num.Int().I()))

		return v, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := AssertNumber(env, obj, "")
		if err != nil {
			return nil, err
		}

		v := reflect.New(typ).Elem()
		v.SetInt(int64(num.Int().I()))

		return v, nil
	case reflect.Slice:
		if typ.Elem() == reflect.TypeFor[byte]() {
			str, err := AssertString(env, obj, "")
			if err != nil {
				return nil, err
			}

			v := reflect.MakeSlice(typ, 0, len(str.S()))
			v = reflect.AppendSlice(v, reflect.ValueOf([]byte(str.S())))

			return v, nil
		}
	case reflect.String:
		if slice, ok := obj.([]byte); ok {
			return string(slice), nil
		}

		if s, ok := obj.(String); ok {
			return s.S(), nil
		}

	case reflect.Interface:
		rt := reflect.TypeOf(obj)
		if rt.AssignableTo(typ) {
			return obj, nil
		}
	}

	return nil, env.NewError("unable to cast to type: " + typ.Name())
}

func makePtr(t reflect.Type) reflect.Value {
	return reflect.New(t)
}

func newVal(t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Slice:
		return reflect.MakeSlice(t, 0, 0)
	case reflect.Map:
		return reflect.MakeMap(t)
	case reflect.Chan:
		return reflect.MakeChan(t, 0)
	default:
		return reflect.New(t)
	}
}

func newVal2(t reflect.Type, size int) reflect.Value {
	switch t.Kind() {
	case reflect.Slice:
		return reflect.MakeSlice(t, size, size)
	case reflect.Map:
		return reflect.MakeMap(t)
	case reflect.Chan:
		return reflect.MakeChan(t, 0)
	default:
		return reflect.New(t)
	}
}

func makeMapType(k, v reflect.Type) reflect.Type {
	return reflect.MapOf(k, v)
}

func makeSliceType(e reflect.Type) reflect.Type {
	return reflect.SliceOf(e)
}

func makeChanType(e reflect.Type) reflect.Type {
	return reflect.ChanOf(reflect.BothDir, e)
}

func makeStructType(env *Env, m Map) (reflect.Type, error) {
	iter := m.Iter()

	var fields []reflect.StructField

	for iter.HasNext() {
		p := iter.Next()

		var f reflect.StructField

		switch sv := p.Key.(type) {
		case Symbol:
			f.Name = sv.Name()
		case Keyword:
			f.Name = sv.Name()
		case String:
			f.Name = sv.S()
		default:
			return nil, env.NewError("name must be symbol/keyword/string only")
		}

		vt, ok := p.Value.(*Type)
		if !ok {
			return nil, env.NewError("value must be a ReflectType")
		}

		f.Type = vt.ReflectType()

		fields = append(fields, f)
	}

	return reflect.StructOf(fields), nil
}

func derefPtr(env *Env, rv reflect.Value) (any, error) {
	if rv.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("derefPtr only takes pointers")
	}

	return fromAny(env, rv.Elem().Interface())
}

func appendVal(env *Env, rv reflect.Value, ov reflect.Value) (reflect.Value, error) {
	if rv.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("append only takes slices")
	}

	if ov.Kind() == reflect.Slice && ov.Elem().Type() != rv.Elem().Type() {
		return reflect.AppendSlice(rv, ov), nil
	}

	if rv.Elem().Type() != ov.Type() {
		return reflect.Value{}, fmt.Errorf("value not of correct type: %s", rv.Elem().Type().String())
	}

	return reflect.Append(rv, ov), nil
}

func copySlice(env *Env, rv reflect.Value, ov reflect.Value) (int, error) {
	if rv.Kind() != reflect.Slice {
		return 0, fmt.Errorf("append only takes slices")
	}

	if ov.Kind() == reflect.Slice && ov.Type().Elem() == rv.Type().Elem() {
		return reflect.Copy(rv, ov), nil
	}

	if rv.Type().Elem() != ov.Type() {
		return 0, fmt.Errorf("value not of correct type: %s, was: %s", rv.Type().Elem().String(), ov.Type())
	}

	rv.Index(0).Set(ov)

	return 1, nil
}

func sliceSlice2(env *Env, rv reflect.Value, start, end int) (reflect.Value, error) {
	if rv.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("append only takes slices")
	}

	return rv.Slice(start, end), nil
}

func sliceSlice3(env *Env, rv reflect.Value, start, end, capa int) (reflect.Value, error) {
	if rv.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("append only takes slices")
	}

	return rv.Slice3(start, end, capa), nil
}

var nsSubs = strings.NewReplacer(
	"github.com", "github",
	"gitlab.com", "gitlab",
	"/", ".",
)

func SetupPkgReflect(env *Env) ([]*Namespace, error) {
	var ret []*Namespace

	typedMethods := map[reflect.Type]reifiedType{}

	var pkgs []any

	for name, pkg := range pkgreflect.Registry() {
		nsName := nsSubs.Replace(name)
		m := EmptyArrayMap()
		m.AddEqu(MakeKeyword("source"), MakeString("golang"))
		pkgs = append(pkgs, MakeSymbolWithMeta(nsName, m))

		b := NewNSBuilder(env, nsName)

		b.NSMeta(pkg.Doc, "1.0")

		for name, typ := range pkg.Types {
			b.DefType(&DefTypeInfo{
				Name:  name,
				Doc:   typ.Doc,
				Added: "1.0",
				Type:  typ.Value,
			})

			var keys []string
			for n := range typ.Methods {
				keys = append(keys, n)
			}

			sort.Strings(keys)

			var objs []any
			methods := map[string]reifiedFunc{}
			for _, n := range keys {
				m := typ.Methods[n]
				var name Symbol
				if m.Tag != "" {
					tag := MakeSymbol(m.Tag)
					name = MakeTaggedSymbol(nsName+"/"+n, tag)
				} else {
					name = MakeSymbol(nsName + "/" + n)
				}

				methods[n] = reifiedFunc{
					Func: m,
					Name: name,
				}

				objs = append(objs, name)
			}

			typedMethods[typ.Value] = reifiedType{
				Namespace: nsName,
				Methods:   methods,
				MethodVec: NewVectorFrom(objs...),
			}
		}

		for name, val := range pkg.Functions {
			var args []string
			tags := map[string]string{}

			for _, a := range val.Args {
				args = append(args, a.Name)
				if a.Tag != "" {
					tags[a.Name] = a.Tag
				}
			}

			b.Defn(&DefnInfo{
				Name:    name,
				Doc:     val.Doc,
				Added:   "1.0",
				Tag:     val.Tag,
				Args:    args,
				ArgTags: tags,
				Fn:      val.Value,
			})
		}

		for name, val := range pkg.Consts {
			b.DefVar(&DefVarInfo{
				Name:  name,
				Doc:   val.Doc,
				Added: "1.0",
				Value: val.Value,
			})
		}

		for name, val := range pkg.Variables {
			b.DefVar(&DefVarInfo{
				Name:  name,
				Doc:   val.Doc,
				Added: "1.0",
				Value: val.Value,
			})
		}

		ret = append(ret, b.ns)
	}

	b := NewNSBuilder(env, "lace.reflect")
	b.Defn(&DefnInfo{
		Name:  "methods",
		Doc:   "Returns the list of methods on the given value.",
		Added: "1.0",
		Tag:   "Seq",
		Fn: func(env *Env, obj any) Seq {
			return listMethods(env, obj, typedMethods)
		},
	})

	b.DefVar(&DefVarInfo{
		Name:  "*golang-packages*",
		Doc:   "Returns the list of golang packages that are loaded.",
		Added: "1.0",
		Tag:   "Seq",

		Value: reflect.ValueOf(NewVectorFrom(pkgs...)),
	})

	b.Defn(&DefnInfo{
		Name:  "cast",
		Doc:   "Cast a given value to a Go type.",
		Added: "1.0",
		Fn:    castObjectToRef,
	})

	b.Defn(&DefnInfo{
		Name:  "deref",
		Doc:   "Read the value of a pointer value.",
		Added: "1.0",
		Fn:    derefPtr,
	})

	b.Defn(&DefnInfo{
		Name:  "ptr",
		Doc:   "Create a pointer value of the given type.",
		Added: "1.0",
		Fn:    makePtr,
	})

	b.Defn(&DefnInfo{
		Name:  "get",
		Doc:   "Retrieve a field by name from the given value.",
		Added: "1.0",
		Fn:    structGet,
	})

	b.Defn(&DefnInfo{
		Name:  "set!",
		Doc:   "Set a field by name in the given value.",
		Added: "1.0",
		Fn:    structSetInPlace,
	})

	b.Defn(&DefnInfo{
		Name:  "fields",
		Doc:   "List the fields in a struct value.",
		Added: "1.0",
		Fn:    structList,
	})

	b.Defn(&DefnInfo{
		Name:  "as-map",
		Doc:   "Allow the given struct to be accessed as a map.",
		Added: "1.0",
		Fn:    structAsMap,
	})

	b.Defn(&DefnInfo{
		Name:  "to-map",
		Doc:   "Convert a struct value to a map.",
		Added: "1.0",
		Fn:    structToMap,
	})

	b.Defn(&DefnInfo{
		Name:  "from-map",
		Doc:   "Create a struct value, populating it from the values in the map.",
		Added: "1.0",
		Fn:    structFromMap,
	})

	b.Defn(&DefnInfo{
		Name:  "append",
		Doc:   "Append data to a slice value.",
		Added: "1.0",
		Fn:    appendVal,
	})

	b.Defn(&DefnInfo{
		Name:  "copy",
		Doc:   "Copy data into a slice.",
		Added: "1.0",
		Fn:    copySlice,
	})

	b.Defn(&DefnInfo{
		Name:  "slice",
		Doc:   "Create a subslice of a given slice.",
		Added: "1.0",
		Fns: []ArityFn{
			{
				Args: []string{"start", "end"},
				Fn:   sliceSlice2,
			},
			{
				Args: []string{"start", "end", "capa"},
				Fn:   sliceSlice3,
			},
		},
	})

	b.Defn(&DefnInfo{
		Name:  "new",
		Doc:   "Make a new value of the given type.",
		Added: "1.0",
		Fns: []ArityFn{
			{
				Args: []string{"type"},
				Fn:   newVal,
			},
			{
				Args: []string{"type", "size"},
				Fn:   newVal2,
			},
		},
	})

	b.Defn(&DefnInfo{
		Name:  "struct-type",
		Doc:   "Create a new Go struct type.",
		Added: "1.0",
		Fn:    makeStructType,
	})

	b.Defn(&DefnInfo{
		Name:  "map-type",
		Doc:   "Create a new Go map type.",
		Added: "1.0",
		Fn:    makeMapType,
	})

	b.Defn(&DefnInfo{
		Name:  "chan-type",
		Doc:   "Create a new Go chan type.",
		Added: "1.0",
		Fn:    makeChanType,
	})

	b.Defn(&DefnInfo{
		Name:  "chan-type",
		Doc:   "Create a new Go chan type.",
		Added: "1.0",
		Fn:    makeSliceType,
	})

	// The go builtin types

	b.DefType(&DefTypeInfo{
		Name:  "any",
		Doc:   "The Go any type",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[any](),
	})

	b.DefType(&DefTypeInfo{
		Name:  "int",
		Doc:   "The Go int",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[int](),
	})

	b.DefType(&DefTypeInfo{
		Name:  "uint",
		Doc:   "The Go uint",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[uint](),
	})

	b.DefType(&DefTypeInfo{
		Name:  "bool",
		Doc:   "The Go bool type",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[bool](),
	})

	b.DefType(&DefTypeInfo{
		Name:  "string",
		Doc:   "The Go string type",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[string](),
	})

	b.DefType(&DefTypeInfo{
		Name:  "int8",
		Doc:   "The Go int8",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[int8](),
	})

	b.DefType(&DefTypeInfo{
		Name:  "int16",
		Doc:   "The Go int16",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[int16](),
	})

	b.DefType(&DefTypeInfo{
		Name:    "int32",
		Doc:     "The Go int32",
		Added:   "1.0",
		Tag:     "ReflectType",
		Type:    reflect.TypeFor[int32](),
		Aliases: []string{"rune"},
	})

	b.DefType(&DefTypeInfo{
		Name:  "int64",
		Doc:   "The Go int64",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[int64](),
	})

	b.DefType(&DefTypeInfo{
		Name:    "uint8",
		Doc:     "The Go uint8",
		Added:   "1.0",
		Tag:     "ReflectType",
		Type:    reflect.TypeFor[uint8](),
		Aliases: []string{"byte"},
	})

	b.DefType(&DefTypeInfo{
		Name:  "uint16",
		Doc:   "The Go uint16",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[uint16](),
	})

	b.DefType(&DefTypeInfo{
		Name:  "uint32",
		Doc:   "The Go uint32",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[uint32](),
	})

	b.DefType(&DefTypeInfo{
		Name:  "uint64",
		Doc:   "The Go uint64",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[uint64](),
	})

	b.DefType(&DefTypeInfo{
		Name:  "float32",
		Doc:   "The Go float32",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[float32](),
	})

	b.DefType(&DefTypeInfo{
		Name:  "float64",
		Doc:   "The Go float64",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[float64](),
	})

	b.DefType(&DefTypeInfo{
		Name:  "bytes",
		Doc:   "A byte slice",
		Added: "1.0",
		Tag:   "ReflectType",
		Type:  reflect.TypeFor[[]byte](),
	})

	ret = append(ret, b.ns)

	return ret, nil
}
