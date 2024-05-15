package core

import (
	"sync"
)

type Var struct {
	InfoHolder
	MetaHolder
	ns             *Namespace
	name           Symbol
	staticVal      Object
	mu             sync.Mutex
	isMacro        bool
	isPrivate      bool
	isDynamic      bool
	isUsed         bool
	isGloballyUsed bool
	taggedType     *Type
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

func (v *Var) WithMeta(env *Env, meta Map) (Object, error) {
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

func (v *Var) AlterMeta(env *Env, fn *Fn, args []Object) (Map, error) {
	return AlterMeta(env, &v.MetaHolder, fn, args)
}

func (v *Var) GetType() *Type {
	return TYPE.Var
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

func (v *Var) Resolve(env *Env) Object {
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

func (v *Var) Call(env *Env, args []Object) (Object, error) {
	vl := v.Resolve(env)
	vs, err := v.ToString(env, false)
	if err != nil {
		return nil, err
	}

	vls, err := v.ToString(env, false)
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

func (v *Var) Deref(env *Env) (Object, error) {
	return v.Resolve(env), nil
}

func (v *Var) SetValue(env *Env, val Object) error {
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

func (v *Var) SetStatic(val Object) {
	v.lock()
	defer v.unlock()

	v.staticVal = val
}

func (v *Var) GetStatic() Object {
	v.lock()
	defer v.unlock()

	return v.staticVal
}

func (v *Var) Set() bool {
	v.lock()
	defer v.unlock()

	return v.staticVal != nil && v.staticVal != NIL
}
