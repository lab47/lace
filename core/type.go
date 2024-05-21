package core

import (
	"fmt"
	"reflect"

	"github.com/lab47/reflectx"
)

type HasGetType interface {
	GetType() *Type
}

func TypeName(obj any) string {
	if hgt, ok := obj.(HasGetType); ok {
		return hgt.GetType().Name()
	}

	t := reflect.Indirect(reflect.ValueOf(obj)).Type()

	return t.PkgPath() + "." + t.Name()
}

func GetType(obj any) any {
	if hgt, ok := obj.(HasGetType); ok {
		return hgt.GetType()
	}

	return &ReflectType{
		typ: reflect.TypeOf(obj),
	}
}

type HasReflectType interface {
	ReflectType() reflect.Type
}

func (t *Type) ReflectType() reflect.Type {
	return t.reflectType
}

type HasHash interface {
	Hash(env *Env) (uint32, error)
}

func HashValue(env *Env, obj any) (uint32, error) {
	if hh, ok := obj.(HasHash); ok {
		return hh.Hash(env)
	}

	return uint32(reflectx.Hash(obj)), nil
}

type HasEquals interface {
	Equals(env *Env, other any) bool
}

func Equals(env *Env, a, b any) bool {
	if a == b {
		return true
	}

	if he, ok := a.(HasEquals); ok {
		return he.Equals(env, b)
	}

	if he, ok := b.(HasEquals); ok {
		return he.Equals(env, a)
	}

	return a == b
}

type HasToString interface {
	ToString(env *Env, escape bool) (string, error)
}

func ToString(env *Env, obj any) (string, error) {
	if hs, ok := obj.(HasToString); ok {
		return hs.ToString(env, false)
	}

	var val reflect.Value

	if rv, ok := obj.(*ReflectValue); ok {
		val = rv.val
	} else {
		val = reflect.ValueOf(obj)
	}

	t := val.Type()
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	pkg := t.PkgPath()
	name := t.Name()

	if name == "" {
		return fmt.Sprintf("#go.%s[%s]", t.Kind(), val.String()), nil
	}

	for val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Pointer:
		return fmt.Sprintf("#%s.%s[%p]", pkg, name, val.Interface()), nil
	default:
		return fmt.Sprintf("#%s.%s[%s]", pkg, name, val.Interface()), nil
	}
}

func SimpleToString(obj any) string {
	switch s := obj.(type) {
	case String:
		return s.S()
	default:
		return fmt.Sprint(obj)
	}
}

func AsMap(obj any) (Map, bool) {
	if m, ok := obj.(Map); ok {
		return m, true
	}

	rv := reflect.ValueOf(obj)
	if reflect.Indirect(rv).Kind() == reflect.Struct {
		return &ReflectValue{val: rv}, true
	}

	return nil, false
}

func AsGettable(obj any) (Gettable, bool) {
	if m, ok := obj.(Gettable); ok {
		return m, true
	}

	rv := reflect.ValueOf(obj)
	if reflect.Indirect(rv).Kind() == reflect.Struct {
		return &ReflectValue{val: rv}, true
	}

	return nil, false
}

func AsAssociative(obj any) (Associative, bool) {
	if m, ok := obj.(Associative); ok {
		return m, true
	}

	rv := reflect.ValueOf(obj)
	if reflect.Indirect(rv).Kind() == reflect.Struct {
		return &ReflectValue{val: rv}, true
	}

	return nil, false
}
