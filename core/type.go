package core

import (
	"fmt"
	"reflect"

	"github.com/lab47/reflectx"
)

// A value that describes a set of values.
//
//lace:export
type Type struct {
	rType reflect.Type
}

func (t Type) ToString(env *Env, escape bool) (string, error) {
	return t.Name(), nil
}

func (tt Type) Name() string {
	t := tt.rType

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	pkg := t.PkgPath()
	if pkg == "github.com/lab47/lace/core" {
		return t.Name()
	}

	return fmt.Sprintf("%s.%s", pkg, t.Name())
}

type HasGetType interface {
	GetType() Type
}

func TypeName(obj any) string {
	if obj == nil {
		return "nil"
	}

	t := reflect.Indirect(reflect.ValueOf(obj)).Type()
	return Type{
		rType: t,
	}.Name()
}

func GetType(obj any) any {
	if hgt, ok := obj.(HasGetType); ok {
		return hgt.GetType()
	}

	return Type{
		rType: reflect.TypeOf(obj),
	}
}

type HasReflectType interface {
	ReflectType() reflect.Type
}

func (t Type) ReflectType() reflect.Type {
	return t.rType
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

	val := reflect.ValueOf(obj)

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
