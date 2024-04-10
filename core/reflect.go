package core

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"

	"github.com/lab47/lace/pkg/pkgreflect"
	_ "github.com/lab47/lace/reflect"
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

type ReflectValue struct {
	InfoHolder
	MetaHolder

	val reflect.Value
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

func (r *ReflectValue) ToString(env *Env, escape bool) (string, error) {
	t := r.val.Type()
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	pkg := t.PkgPath()
	name := t.Name()

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

func convertToInt64(env *Env, index int, obj Object) (int64, error) {
	switch sv := obj.(type) {
	case Int:
		return int64(sv.I), nil
	case *BigInt:
		return sv.b.Int64(), nil
	default:
		return 0, env.RT.NewArgTypeError(index, obj, "int")
	}
}

func convertToUInt64(env *Env, index int, obj Object) (uint64, error) {
	switch sv := obj.(type) {
	case Int:
		return uint64(sv.I), nil
	case *BigInt:
		return sv.b.Uint64(), nil
	default:
		return 0, env.RT.NewArgTypeError(index, obj, "int")
	}
}

func convertToString(env *Env, index int, obj Object) (string, error) {
	switch sv := obj.(type) {
	case String:
		return sv.S, nil
	case Symbol:
		return sv.Name(), nil
	case Keyword:
		return sv.Name(), nil
	default:
		return "", env.RT.NewArgTypeError(index, obj, "string/symbol/keyword")
	}
}

func convertToBytes(env *Env, index int, obj Object) ([]byte, error) {
	switch sv := obj.(type) {
	case String:
		return []byte(sv.S), nil
	default:
		return nil, env.RT.NewArgTypeError(index, obj, "string/symbol/keyword")
	}
}

func convertFromInt(env *Env, rv reflect.Value) (Object, error) {
	i := rv.Int()

	if i > math.MaxInt {
		return MakeBigInt(i), nil
	}

	return MakeInt(int(i)), nil
}

func convertFromUInt(env *Env, rv reflect.Value) (Object, error) {
	i := rv.Uint()

	if i > math.MaxUint {
		return MakeBigInt(int64(i)), nil
	}

	return MakeInt(int(i)), nil
}

func convertFromString(env *Env, rv reflect.Value) (Object, error) {
	return MakeString(rv.String()), nil
}

func convertFromBool(env *Env, rv reflect.Value) (Object, error) {
	return MakeBoolean(rv.Bool()), nil
}

func convertFromAny(env *Env, rv reflect.Value) (Object, error) {
	return &ReflectValue{val: rv}, nil
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

	return nil
}
