package core

import "fmt"

// A value that contains code and can be called to run that code.
//
//lace:export
type Fn struct {
	InfoHolder
	MetaHolder
	isMacro bool
	fnExpr  *FnExpr
	env     *LocalEnv

	code           *Code
	importedUpvals []*NamedPair
}

func (fn *Fn) ToString(env *Env, escape bool) (string, error) {
	if fn.code == nil {
		return "#Fn[]", nil
	}

	return fmt.Sprintf("#Fn[%s:%d]", fn.code.filename, fn.code.lineForIp(0)), nil
}

func (fn *Fn) String() string {
	if fn.code != nil {
		return fmt.Sprintf("<fn bc @ %s:%d>", fn.code.filename, fn.code.lineForIp(0))
	} else {
		pos := fn.fnExpr.Pos()
		return fmt.Sprintf("<fn tree @ %s:%d>", pos.filename, pos.startLine)
	}
}

func (fn *Fn) Equals(env *Env, other interface{}) bool {
	switch other := other.(type) {
	case *Fn:
		return fn == other
	default:
		return false
	}
}

func (fn *Fn) WithMeta(env *Env, meta Map) (any, error) {
	res := *fn
	m, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}
	res.meta = m
	return &res, nil
}

func (fn *Fn) Hash(env *Env) (uint32, error) {
	return HashPtr(fn), nil
}

func (fn *Fn) Call(env *Env, args []any) (any, error) {
	obj, err := env.Engine.RunWithArgs(env, fn, args)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (fn *Fn) Compare(env *Env, a, b any) (int, error) {
	return compare(env, fn, a, b)
}

var _ Callable = (*Fn)(nil)
