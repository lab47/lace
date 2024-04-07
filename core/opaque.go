package core

import (
	"fmt"
	"unsafe"
)

type (
	Opaque[T any] struct {
		Value T
	}
)

func (f *Opaque[T]) ToString(escape bool) string {
	var t T
	return fmt.Sprintf("#object[Opaque[%T]]", t)
}

func (f *Opaque[T]) Equals(other interface{}) bool {
	return f == other
}

func (f *Opaque[T]) GetInfo() *ObjectInfo {
	return nil
}

func (f *Opaque[T]) GetType() *Type {
	return TYPE.Opaque
}

func (f *Opaque[T]) Hash() uint32 {
	return HashPtr(uintptr(unsafe.Pointer(f)))
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
		return env.RT.NewError("expected typed opaque, was not")
	}

	*dest = box.Value
	return nil
}
