package core

func Cast[A any](env *Env, obj Object, a *A) error {
	x, ok := obj.(A)
	if !ok {
		return TypeError[A](env, obj)
	}

	*a = x
	return nil
}

func CoerceString(env *Env, obj Object, a *string) error {
	switch sv := obj.(type) {
	case String:
		*a = sv.S()
	case Keyword:
		*a = sv.RawString()
	case Symbol:
		*a = sv.String()
	default:
		return TypeError[String](env, obj)
	}

	return nil
}
