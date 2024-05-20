package core

import "reflect"

func WrapToProc0_0(fn func()) any {
	return func(env *Env, args []any) (any, error) {
		if len(args) != 0 {
			return nil, ErrorArityMinMax(env, len(args), 0, 0)
		}

		fn()

		return NIL, nil
	}
}

func match[A, B any]() bool {
	return reflect.TypeFor[A]() == reflect.TypeFor[B]()
}

func WrapToProc0_1[O any](fn func() O) any {
	if xfn, ok := any(fn).(func() error); ok {
		return func(env *Env, args []any) (any, error) {
			if len(args) != 0 {
				return nil, ErrorArityMinMax(env, len(args), 0, 0)
			}

			return NIL, WrapError(env, xfn())
		}
	}

	if xfn, ok := any(fn).(func() any); ok {
		return func(env *Env, args []any) (any, error) {
			if len(args) != 0 {
				return nil, ErrorArityMinMax(env, len(args), 0, 0)
			}

			return xfn(), nil
		}
	}

	return fn
}

func WrapToProc0_2[O, O2 any](fn func() (O, O2)) any {
	if xfn, ok := any(fn).(func() (any, error)); ok {
		return func(env *Env, args []any) (any, error) {
			if len(args) != 0 {
				return nil, ErrorArityMinMax(env, len(args), 0, 0)
			}

			obj, err := xfn()
			if err != nil {
				return nil, WrapError(env, err)
			}

			return obj, err
		}
	}

	return fn
}

func WrapToProc1_0[A any](fn func(a A)) any {
	cs := convReg.buildCS(reflect.TypeFor[func(A)]())

	afn := cs.argIn[0]

	return func(env *Env, args []any) (any, error) {
		if len(args) != cs.arity {
			return nil, ErrorArityMinMax(env, len(args), cs.arity, cs.arity)
		}

		arv, err := afn(env, 0, args[0])
		if err != nil {
			return nil, WrapError(env, err)
		}

		fn(arv.Interface().(A))

		return NIL, nil
	}
}

func WrapToProc1_1[A, O any](fn func(a A) O) any {
	if xfn, ok := any(fn).(func(A) error); ok {
		cs := convReg.buildCS(reflect.TypeOf(xfn))
		afn := cs.argIn[0]

		return func(env *Env, args []any) (any, error) {
			if len(args) != 1 {
				return nil, ErrorArityMinMax(env, len(args), 0, 0)
			}

			arv, err := afn(env, 0, args[0])
			if err != nil {
				return nil, WrapError(env, err)
			}

			err = xfn(arv.Interface().(A))
			return NIL, WrapError(env, err)
		}
	}

	if xfn, ok := any(fn).(func(A) any); ok {
		cs := convReg.buildCS(reflect.TypeOf(xfn))
		afn := cs.argIn[0]

		return func(env *Env, args []any) (any, error) {
			if len(args) != 1 {
				return nil, ErrorArityMinMax(env, len(args), 0, 0)
			}

			arv, err := afn(env, 0, args[0])
			if err != nil {
				return nil, WrapError(env, err)
			}

			obj := xfn(arv.Interface().(A))

			return obj, nil
		}
	}

	if xfn, ok := any(fn).(func(*Env) error); ok {
		return func(env *Env, args []any) (any, error) {
			if len(args) != 0 {
				return nil, ErrorArityMinMax(env, len(args), 0, 0)
			}

			err := xfn(env)
			return NIL, WrapError(env, err)
		}
	}

	if xfn, ok := any(fn).(func(*Env) any); ok {
		return func(env *Env, args []any) (any, error) {
			if len(args) != 0 {
				return nil, ErrorArityMinMax(env, len(args), 0, 0)
			}

			obj := xfn(env)
			return obj, nil
		}
	}

	return nil
}

func WrapToProc1_2[A, O, O2 any](fn func(A) (O, O2)) any {
	if xfn, ok := any(fn).(func(A) (any, error)); ok {
		cs := convReg.buildCS(reflect.TypeOf(xfn))
		afn := cs.argIn[0]

		return func(env *Env, args []any) (any, error) {
			if len(args) != 1 {
				return nil, ErrorArityMinMax(env, len(args), 0, 0)
			}

			arv, err := afn(env, 0, args[0])
			if err != nil {
				return nil, WrapError(env, err)
			}

			obj, err := xfn(arv.Interface().(A))
			if err != nil {
				return nil, WrapError(env, err)
			}

			return obj, nil
		}
	}

	if xfn, ok := any(fn).(func(*Env) (any, error)); ok {
		return func(env *Env, args []any) (any, error) {
			if len(args) != 0 {
				return nil, ErrorArityMinMax(env, len(args), 0, 0)
			}

			obj, err := xfn(env)
			if err != nil {
				return nil, WrapError(env, err)
			}

			return obj, nil
		}
	}

	return fn
}

func WrapToProc2_0[E, A any](fn func(e E, a A)) any {
	if xfn, ok := any(fn).(func(*Env, A)); ok {
		cs := convReg.buildCS(reflect.TypeOf(xfn))

		afn := cs.argIn[1]

		return func(env *Env, args []any) (any, error) {
			if len(args) != 1 {
				return nil, ErrorArityMinMax(env, len(args), cs.arity, cs.arity)
			}

			arv, err := afn(env, 0, args[0])
			if err != nil {
				return nil, WrapError(env, err)
			}

			xfn(
				env,
				arv.Interface().(A),
			)

			return NIL, nil
		}
	}

	return fn
}

func WrapToProc2_1[E, A, O any](fn func(e E, a A) O) any {
	if xfn, ok := any(fn).(func(*Env, A) error); ok {
		cs := convReg.buildCS(reflect.TypeOf(xfn))
		afn := cs.argIn[1]

		return func(env *Env, args []any) (any, error) {
			if len(args) != 1 {
				return nil, ErrorArityMinMax(env, len(args), cs.arity, cs.arity)
			}

			arv, err := afn(env, 0, args[0])
			if err != nil {
				return nil, WrapError(env, err)
			}

			err = xfn(
				env,
				arv.Interface().(A),
			)

			return NIL, WrapError(env, err)
		}
	}

	if xfn, ok := any(fn).(func(*Env, A) any); ok {
		cs := convReg.buildCS(reflect.TypeOf(xfn))
		afn := cs.argIn[1]

		return func(env *Env, args []any) (any, error) {
			if len(args) != 1 {
				return nil, ErrorArityMinMax(env, len(args), cs.arity, cs.arity)
			}

			arv, err := afn(env, 0, args[0])
			if err != nil {
				return nil, WrapError(env, err)
			}

			return xfn(
				env,
				arv.Interface().(A),
			), nil
		}
	}

	return fn
}

func WrapToProc2_2[E, A, O, O2 any](fn func(e E, a A) (O, O2)) any {
	if match[func(E, A) (O, O2), func(*Env, []any) (any, error)]() {
		return fn
	}

	if xfn, ok := any(fn).(func(*Env, A) (any, error)); ok {
		cs := convReg.buildCS(reflect.TypeOf(xfn))
		afn := cs.argIn[1]

		return func(env *Env, args []any) (any, error) {
			if len(args) != 1 {
				return nil, ErrorArityMinMax(env, len(args), cs.arity, cs.arity)
			}

			arv, err := afn(env, 0, args[0])
			if err != nil {
				return nil, err
			}

			return xfn(
				env,
				arv.Interface().(A),
			)
		}
	}

	return fn
}

func WrapToProc3_0[E, A, B any](fn func(e E, a A, b B)) any {
	if xfn, ok := any(fn).(func(*Env, A, B)); ok {
		cs := convReg.buildCS(reflect.TypeOf(xfn))

		afn := cs.argIn[1]
		bfn := cs.argIn[2]

		return func(env *Env, args []any) (any, error) {
			if len(args) != cs.arity {
				return nil, ErrorArityMinMax(env, len(args), cs.arity, cs.arity)
			}

			arv, err := afn(env, 0, args[0])
			if err != nil {
				return nil, WrapError(env, err)
			}

			brv, err := bfn(env, 1, args[1])
			if err != nil {
				return nil, WrapError(env, err)
			}

			xfn(
				env,
				arv.Interface().(A),
				brv.Interface().(B),
			)

			return NIL, nil
		}
	}

	return fn
}

func WrapToProc3_1[E, A, B, O any](fn func(e E, a A, b B) O) any {
	if xfn, ok := any(fn).(func(*Env, A, B) error); ok {
		cs := convReg.buildCS(reflect.TypeOf(xfn))

		afn := cs.argIn[1]
		bfn := cs.argIn[2]

		return func(env *Env, args []any) (any, error) {
			if len(args) != cs.arity {
				return nil, ErrorArityMinMax(env, len(args), cs.arity, cs.arity)
			}

			arv, err := afn(env, 0, args[0])
			if err != nil {
				return nil, WrapError(env, err)
			}

			brv, err := bfn(env, 1, args[1])
			if err != nil {
				return nil, WrapError(env, err)
			}
			err = xfn(
				env,
				arv.Interface().(A),
				brv.Interface().(B),
			)

			return NIL, WrapError(env, err)
		}
	}

	if xfn, ok := any(fn).(func(*Env, A, B) any); ok {
		cs := convReg.buildCS(reflect.TypeOf(xfn))

		afn := cs.argIn[1]
		bfn := cs.argIn[2]

		return func(env *Env, args []any) (any, error) {
			if len(args) != cs.arity {
				return nil, ErrorArityMinMax(env, len(args), cs.arity, cs.arity)
			}

			arv, err := afn(env, 0, args[0])
			if err != nil {
				return nil, WrapError(env, err)
			}

			brv, err := bfn(env, 1, args[1])
			if err != nil {
				return nil, WrapError(env, err)
			}

			obj := xfn(
				env,
				arv.Interface().(A),
				brv.Interface().(B),
			)

			return obj, nil
		}
	}

	return fn
}

func WrapToProc3_2[E, A, B, O, O2 any](fn func(e E, a A, b B) (O, O2)) any {
	if xfn, ok := any(fn).(func(*Env, A, B) (any, error)); ok {
		cs := convReg.buildCS(reflect.TypeOf(xfn))

		afn := cs.argIn[1]
		bfn := cs.argIn[2]

		return func(env *Env, args []any) (any, error) {
			if len(args) != cs.arity {
				return nil, ErrorArityMinMax(env, len(args), cs.arity, cs.arity)
			}

			arv, err := afn(env, 0, args[0])
			if err != nil {
				return nil, WrapError(env, err)
			}

			brv, err := bfn(env, 1, args[1])
			if err != nil {
				return nil, WrapError(env, err)
			}

			o, err := xfn(
				env,
				arv.Interface().(A),
				brv.Interface().(B),
			)

			if err != nil {
				return nil, WrapError(env, err)
			}

			return o, nil
		}
	}

	return fn
}
