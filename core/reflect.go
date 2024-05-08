package core

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/lab47/lace/pkg/pkgreflect"
)

type ReflectType struct {
	InfoHolder
	MetaHolder

	typ reflect.Type
}

var _ Object = &ReflectType{}

func (r *ReflectType) Equals(env *Env, other any) bool {
	if ov, ok := other.(*ReflectType); ok {
		return ov.typ == r.typ
	}

	return false
}

func (r *ReflectType) GetType() *Type {
	return TYPE.ReflectType
}

func (r *ReflectType) ToString(env *Env, escape bool) (string, error) {
	return fmt.Sprintf("#reflect.Type[%s]", r.typ), nil
}

func (r *ReflectType) Hash(env *Env) (uint32, error) {
	h := getHash()
	h.Write([]byte(r.typ.Name()))
	return h.Sum32(), nil
}

func (r *ReflectType) WithInfo(i *ObjectInfo) Object {
	d := *r
	d.info = i
	return &d
}

func (r *ReflectType) Call(env *Env, args []Object) (Object, error) {
	rv := reflect.New(r.typ)
	return &ReflectValue{val: rv}, nil
}

type ReflectValue struct {
	InfoHolder
	MetaHolder

	val reflect.Value
}

func MakeReflectValue(v any) *ReflectValue {
	return &ReflectValue{
		val: reflect.ValueOf(v),
	}
}

func CastReflect[T any](env *Env, obj Object, v *T) error {
	var rv *ReflectValue

	err := Cast(env, obj, &rv)
	if err != nil {
		return err
	}

	x, ok := rv.Value().(T)
	if !ok {
		return fmt.Errorf("not a reflect value containing a *%T", v)
	}

	*v = x

	return nil
}

func (r *ReflectValue) Value() any {
	return r.val.Interface()
}

var _ Object = &ReflectValue{}

func (r *ReflectValue) Equals(env *Env, other any) bool {
	if ov, ok := other.(*ReflectValue); ok {
		return ov.val == r.val
	}

	return false
}

func (r *ReflectValue) GetType() *Type {
	return TYPE.ReflectValue
}

func structPut(env *Env, r *ReflectValue, name string, fval Object) error {
	val := r.val

	for val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return env.RT.NewError(fmt.Sprintf("value must be a struct, is a %T", val.Interface()))
	}

	field := val.FieldByName(name)
	if !field.IsValid() {
		return env.RT.NewError(fmt.Sprintf("unknown struct field %s", name))
	}

	var (
		frv reflect.Value
		err error
	)

	if field.Type().Kind() == reflect.Func {
		call, ok := fval.(Callable)
		if !ok {
			return env.RT.TypeError(TCContext{Context: "struct value"}, fval, "Callable")
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
		return env.RT.NewError(
			fmt.Sprintf("needed type %s, had %T", field.Type(), fval))
	}

	field.Set(frv)

	return nil
}

func structGet(env *Env, r *ReflectValue, name string) (Object, error) {
	val := r.val

	for val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, env.RT.NewError(fmt.Sprintf("value must be a struct, is a %T", val.Interface()))
	}

	field := val.FieldByName(name)
	if !field.IsValid() {
		return nil, env.RT.NewError(fmt.Sprintf("unknown struct field %s", name))
	}

	rt := field.Type()

	return convReg.convRet(rt)(env, field)
}

func (r *ReflectValue) ToString(env *Env, escape bool) (string, error) {
	t := r.val.Type()
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	pkg := t.PkgPath()
	name := t.Name()

	if name == "" {
		return fmt.Sprintf("#go.%s[%s]", t.Kind(), r.val.String()), nil
	}

	return fmt.Sprintf("#%s.%s[%s]", pkg, name, r.val), nil
}

func (r *ReflectValue) Hash(env *Env) (uint32, error) {
	h := getHash()
	h.Write([]byte(r.val.String()))
	return h.Sum32(), nil
}

func (r *ReflectValue) WithInfo(i *ObjectInfo) Object {
	d := *r
	d.info = i
	return &d
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

func listMethods(env *Env, obj Object, reg map[reflect.Type]reifiedType) Seq {
	rv, ok := obj.(*ReflectValue)
	if !ok {
		return NIL
	}

	t := rv.val.Type()

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	var objs []Object
	if meths, ok := reg[t]; ok {
		return meths.MethodVec.Seq()
	} else {
		for i := 0; i < t.NumMethod(); i++ {
			objs = append(objs, MakeKeyword(t.Method(i).Name))
		}
	}

	return NewListFrom(objs...)
}

func castObjectToRef(env *Env, typ reflect.Type, obj Object) (Object, error) {
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := AssertNumber(env, obj, "")
		if err != nil {
			return nil, err
		}

		v := reflect.New(typ).Elem()
		v.SetInt(int64(num.Int().I))

		return &ReflectValue{val: v}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := AssertNumber(env, obj, "")
		if err != nil {
			return nil, err
		}

		v := reflect.New(typ).Elem()
		v.SetInt(int64(num.Int().I))

		return &ReflectValue{val: v}, nil
	case reflect.Slice:
		if typ.Elem() == reflect.TypeFor[byte]() {
			str, err := AssertString(env, obj, "")
			if err != nil {
				return nil, err
			}

			v := reflect.MakeSlice(typ, 0, len(str.S))
			v = reflect.AppendSlice(v, reflect.ValueOf([]byte(str.S)))

			return &ReflectValue{val: v}, nil
		}
	case reflect.String:
		if rv, ok := obj.(*ReflectValue); ok {
			if sv, ok := rv.val.Interface().([]byte); ok {
				return MakeString(string(sv)), nil
			}
		}

		if _, ok := obj.(String); ok {
			return obj, nil
		}

	case reflect.Interface:
		if rv, ok := obj.(*ReflectValue); ok {
			if rv.val.Type().AssignableTo(typ) {
				return rv, nil
			}
		}
	}

	return nil, env.RT.NewError("unable to cast to type: " + typ.Name())
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
			f.Name = sv.S
		default:
			return nil, env.RT.NewError("name must be symbol/keyword/string only")
		}

		vt, ok := p.Value.(*ReflectType)
		if !ok {
			return nil, env.RT.NewError("value must be a ReflectType")
		}

		f.Type = vt.typ

		fields = append(fields, f)
	}

	return reflect.StructOf(fields), nil
}

func derefPtr(env *Env, rv reflect.Value) (Object, error) {
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

func SetupPkgReflect(env *Env) error {
	typedMethods := map[reflect.Type]reifiedType{}

	var pkgs []Object

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

			var objs []Object
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
	}

	b := NewNSBuilder(env, "lace.reflect")
	b.Defn(&DefnInfo{
		Name:  "methods",
		Doc:   "Returns the list of methods on the given value.",
		Added: "1.0",
		Tag:   "Seq",
		Fn: func(env *Env, obj Object) Seq {
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
		Name:  "put",
		Doc:   "Set a field by name in the given value.",
		Added: "1.0",
		Fn:    structPut,
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

	return nil
}
