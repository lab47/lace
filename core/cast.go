package core

func Cast[A any](env *Env, obj Object, a *A) error {
	x, ok := obj.(A)
	if !ok {
		return TypeError[A](env, obj)
	}

	*a = x
	return nil
}
