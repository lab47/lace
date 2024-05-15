package core

import (
	"fmt"
)

type (
	Opaque[T any] struct {
		Value T
	}
)

var _ Object = &Opaque[bool]{}

func (f *Opaque[T]) ToString(env *Env, escape bool) (string, error) {
	var t T
	return fmt.Sprintf("#object[Opaque[%T]]", t), nil
}

func (f *Opaque[T]) Equals(env *Env, other interface{}) bool {
	return f == other
}

func (f *Opaque[T]) GetInfo() *ObjectInfo {
	return nil
}

func (f *Opaque[T]) GetType() *Type {
	return TYPE.Opaque
}

func (f *Opaque[T]) Hash(env *Env) (uint32, error) {
	return HashPtr(f), nil
}

func (f *Opaque[T]) WithInfo(info *ObjectInfo) Object {
	return f
}

func MakeOpaque[T any](f T) *Opaque[T] {
	return &Opaque[T]{f}
}

func ExtractOpaque[T any](env *Env, obj Object, dest *T) error {
	box, ok := obj.(*Opaque[T])
	if !ok {
		return env.NewError("expected typed opaque, was not")
	}

	*dest = box.Value
	return nil
}
