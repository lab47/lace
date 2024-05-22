package core

import (
	"sync"
)

// A value that holds another value and can be changed. Ie, a Var is a variable.
//
//lace:export
type Var struct {
	InfoHolder
	MetaHolder
	ns             *Namespace
	name           Symbol
	staticVal      any
	mu             sync.Mutex
	isMacro        bool
	isPrivate      bool
	isDynamic      bool
	isUsed         bool
	isGloballyUsed bool
	taggedType     any
}

func (v *Var) Name() string {
	return v.ns.Name.String() + "/" + v.name.String()
}

func (v *Var) ToString(env *Env, escape bool) (string, error) {
	return "#'" + v.Name(), nil
}

func (v *Var) String() string {
	return "#'" + v.Name()
}

func (v *Var) Equals(env *Env, other interface{}) bool {
	// TODO: revisit this
	return v == other
}

func (v *Var) WithMeta(env *Env, meta Map) (any, error) {
	res := &Var{
		ns:         v.ns,
		name:       v.name,
		staticVal:  v.staticVal,
		isMacro:    v.isMacro,
		isPrivate:  v.isPrivate,
		isDynamic:  v.isDynamic,
		taggedType: v.taggedType,
	}
	res.meta = v.meta
	res.info = v.info
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return res, nil
}

func (v *Var) ResetMeta(newMeta Map) Map {
	v.meta = newMeta
	return v.meta
}

func (v *Var) AlterMeta(env *Env, fn *Fn, args []any) (Map, error) {
	return AlterMeta(env, &v.MetaHolder, fn, args)
}

func (v *Var) Hash(env *Env) (uint32, error) {
	return HashPtr(v), nil
}

func (v *Var) lock() {
	v.mu.Lock()
}

func (v *Var) unlock() {
	v.mu.Unlock()
}

func (v *Var) Resolve(env *Env) any {
	v.lock()
	isDyn := v.isDynamic
	sval := v.staticVal
	v.unlock()

	if isDyn {
		obj, ok, err := env.FindInCurrentVars(v)
		if err != nil {
			panic(WrapError(env, err))
		}

		if ok {
			return obj
		}
	}
	if sval == nil {
		return NIL
	}
	return sval
}

func (v *Var) Call(env *Env, args []any) (any, error) {
	vl := v.Resolve(env)
	vs, err := ToString(env, v)
	if err != nil {
		return nil, err
	}

	vls, err := ToString(env, v)
	if err != nil {
		return nil, err
	}

	call, err := AssertCallable(env,
		vl,
		"Var "+vs+" resolves to "+vls+", which is not a Fn")
	if err != nil {
		return nil, err
	}

	return call.Call(env, args)
}

var _ Callable = (*Var)(nil)

func (v *Var) Deref(env *Env) (any, error) {
	return v.Resolve(env), nil
}

func (v *Var) SetValue(env *Env, val any) error {
	if v.isDynamic {
		as, err := env.CurrentVar.Assoc(env, v, val)
		if err != nil {
			return nil
		}

		env.CurrentVar = as
		return nil
	}

	v.staticVal = val
	return nil
}

func (v *Var) SetStatic(val any) {
	v.lock()
	defer v.unlock()

	v.staticVal = val
}

func (v *Var) GetStatic() any {
	v.lock()
	defer v.unlock()

	return v.staticVal
}

func (v *Var) Set() bool {
	v.lock()
	defer v.unlock()

	return v.staticVal != nil && v.staticVal != NIL
}
