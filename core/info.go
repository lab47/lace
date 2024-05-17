package core

type InfoWrapper[T Object] struct {
	InfoHolder

	val T
}

var _ Object = (*InfoWrapper[Boolean])(nil)

func (i *InfoWrapper[T]) Equals(env *Env, other any) bool {
	return i.val.Equals(env, other)
}

func (i *InfoWrapper[T]) GetType() *Type {
	return i.val.GetType()
}

func (i *InfoWrapper[T]) Hash(env *Env) (uint32, error) {
	return i.val.Hash(env)
}

func (i *InfoWrapper[T]) ToString(env *Env, escape bool) (string, error) {
	return i.val.ToString(env, escape)
}

func (i *InfoWrapper[T]) WithInfo(info *ObjectInfo) Object {
	x := *i
	i.info = info
	return &x
}

func (i *InfoWrapper[T]) Unwrap() Object {
	return i.val
}

type Unwrapper interface {
	Unwrap() Object
}

func Unwrap(obj Object) Object {
	if v, ok := obj.(Unwrapper); ok {
		return v.Unwrap()
	}

	return obj
}
