package core

import (
	"fmt"
)

type (
	Opaque[T any] struct {
		Value T
	}
)

func (f *Opaque[T]) ToString(env *Env, escape bool) (string, error) {
	var t T
	return fmt.Sprintf("#object[Opaque[%T]]", t), nil
}

func MakeOpaque[T any](f T) *Opaque[T] {
	return &Opaque[T]{f}
}

func ExtractOpaque[T any](env *Env, obj any, dest *T) error {
	box, ok := obj.(*Opaque[T])
	if !ok {
		return env.NewError("expected typed opaque, was not")
	}

	*dest = box.Value
	return nil
}
