package core

import (
	"fmt"
	"reflect"
)

func (e *Env) pushTreeEval(expr Expr) {
	e.treeEvalStack = append(e.treeEvalStack, expr)
}

func (e *Env) popTreeEval() {
	e.treeEvalStack = e.treeEvalStack[:len(e.treeEvalStack)-1]
}

func Eval(genv *Env, expr Expr, env *LocalEnv) (Object, error) {
	genv.pushTreeEval(expr)
	defer genv.popTreeEval()

	obj, err := expr.Eval(genv, env)
	return obj, err
}

func (expr *VarRefExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	return expr.vr.Resolve(genv), nil
}

func (expr *SetMacroExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	expr.vr.isMacro = true
	expr.vr.isUsed = false
	if fn, ok := expr.vr.GetStatic().(*Fn); ok {
		fn.isMacro = true
	}
	err := setMacroMeta(genv, expr.vr)
	if err != nil {
		return nil, err
	}
	return expr.vr, nil
}

func (expr *BindingExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	for i := env.frame; i > expr.binding.frame; i-- {
		env = env.parent
	}
	return env.bindings[expr.binding.index], nil
}

func (expr *LiteralExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	return expr.obj, nil
}

func (expr *VectorExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	res := EmptyVector()
	for _, e := range expr.v {
		o, err := Eval(genv, e, env)
		if err != nil {
			return nil, err
		}
		res, _ = res.Conjoin(o)
	}
	return res, nil
}

func (expr *MapExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	if int64(len(expr.keys)) > HASHMAP_THRESHOLD/2 {
		res := EmptyHashMap
		for i := range expr.keys {
			key, err := Eval(genv, expr.keys[i], env)
			if err != nil {
				return nil, err
			}
			if res.containsKey(genv, key) {
				s, err := key.ToString(genv, false)
				if err != nil {
					return nil, err
				}
				return nil, genv.NewError("Duplicate key: " + s)
			}
			v, err := Eval(genv, expr.values[i], env)
			if err != nil {
				return nil, err
			}
			v, err = res.Assoc(genv, key, v)
			if err != nil {
				return nil, err
			}
			if err := Cast(genv, v, &res); err != nil {
				return nil, err
			}
		}
		return res, nil
	}
	res := EmptyArrayMap()
	for i := range expr.keys {
		key, err := Eval(genv, expr.keys[i], env)
		if err != nil {
			return nil, err
		}
		v, err := Eval(genv, expr.values[i], env)
		if err != nil {
			return nil, err
		}
		if !res.Add(genv, key, v) {
			s, err := key.ToString(genv, false)
			if err != nil {
				return nil, err
			}

			return nil, genv.NewError("Duplicate key: " + s)
		}
	}
	return res, nil
}

func (expr *SetExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	res := EmptySet()
	for _, elemExpr := range expr.elements {
		el, err := Eval(genv, elemExpr, env)
		if err != nil {
			return nil, err
		}
		ok, err := res.Add(genv, el)
		if err != nil {
			return nil, err
		}
		if !ok {
			s, err := el.ToString(genv, false)
			if err != nil {
				return nil, err
			}

			return nil, genv.NewError("Duplicate set element: " + s)
		}
	}
	return res, nil
}

func iEval(dst *Object, genv *Env, obj Expr, env *LocalEnv) error {
	x, err := Eval(genv, obj, env)
	if err != nil {
		return err
	}

	*dst = x
	return nil
}

func (expr *DefExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	if expr.value != nil {
		x, err := Eval(genv, expr.value, env)
		if err != nil {
			return nil, err
		}
		expr.vr.SetStatic(x)
	}
	meta := EmptyArrayMap()
	meta.Add(genv, criticalKeywords.line, MakeInt(expr.startLine))
	meta.Add(genv, criticalKeywords.column, MakeInt(expr.startColumn))
	meta.Add(genv, criticalKeywords.file, MakeString(expr.filename))
	meta.Add(genv, criticalKeywords.ns, expr.vr.ns)
	fullName := AssembleSymbol(expr.vr.ns.Name.Name(), expr.vr.name.Name())
	meta.Add(genv, criticalKeywords.name, fullName)
	expr.vr.meta = meta
	if expr.meta != nil {
		v, err := Eval(genv, expr.meta, env)
		if err != nil {
			return nil, err
		}
		var m Map
		if err := Cast(genv, v, &m); err != nil {
			return nil, err
		}
		expr.vr.meta, err = expr.vr.meta.Merge(genv, m)
		if err != nil {
			return nil, err
		}
	}
	// isMacro can be set by set-macro__ during parse stage
	if expr.vr.isMacro {
		v, err := expr.vr.meta.Assoc(genv, criticalKeywords.macro, Boolean(true))
		if err != nil {
			return nil, err
		}
		var m Map
		if err := Cast(genv, v, &m); err != nil {
			return nil, err
		}
		expr.vr.meta = m
	}

	if m, ok := expr.vr.GetStatic().(*Fn); ok {
		if m.meta == nil {
			m.meta = expr.vr.meta
		} else {
			nm, err := m.meta.Assoc(genv, criticalKeywords.name, fullName)
			if err == nil {
				m.meta = nm.(Map)
			}
		}
	}

	return expr.vr, nil
}

func (expr *MetaExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	meta, err := Eval(genv, expr.meta, env)
	if err != nil {
		return nil, err
	}
	res, err := Eval(genv, expr.expr, env)
	if err != nil {
		return nil, err
	}

	var metao Meta
	if err := Cast(genv, res, &metao); err != nil {
		return nil, err
	}

	var m Map
	if err := Cast(genv, meta, &m); err != nil {
		return nil, err
	}

	return metao.WithMeta(genv, m)
}

func evalSeq(genv *Env, exprs []Expr, env *LocalEnv) ([]Object, error) {
	res := make([]Object, len(exprs))
	for i, expr := range exprs {
		v, err := Eval(genv, expr, env)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func (expr *CallExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	obj, err := Eval(genv, expr.callable, env)
	if err != nil {
		return nil, err
	}

	switch callable := obj.(type) {
	case Callable:
		args, err := evalSeq(genv, expr.args, env)
		if err != nil {
			return nil, err
		}

		return callable.Call(genv, args)
	default:
		s, err := callable.ToString(genv, false)
		if err != nil {
			return nil, err
		}

		return nil, genv.NewError(s + " is not a Fn")
	}
}

func (expr *MethodExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	obj, err := Eval(genv, expr.obj, env)
	if err != nil {
		return nil, err
	}

	var rv reflect.Value
	methName := expr.method

	if orv, ok := obj.(*ReflectValue); ok {
		rv = orv.val
	} else {
		rv = reflect.ValueOf(obj)
	}

	rt := rv.Type()

	if expr.lastType == rt {
		objArgs, err := evalSeq(genv, expr.args, env)
		if err != nil {
			return nil, err
		}

		return expr.lastFn(genv, objArgs)
	}

	meth := rv.MethodByName(methName)

	if !meth.IsValid() {
		return nil, genv.NewError(fmt.Sprintf("unknown method %s on %s", expr.method, rt))
	}

	procFn, _, err := convReg.ConverterForFunc(meth)
	if err != nil {
		return nil, err
	}

	expr.lastType = rt
	expr.lastFn = procFn

	objArgs, err := evalSeq(genv, expr.args, env)
	if err != nil {
		return nil, err
	}

	return procFn(genv, objArgs)
}

func varCallableString(v *Var) string {
	if v.ns.CoreP() {
		return "core/" + v.name.String()
	}
	return v.ns.Name.String() + "/" + v.name.String()
}

func (expr *CallExpr) Name() string {
	switch c := expr.callable.(type) {
	case *VarRefExpr:
		return varCallableString(c.vr)
	case *BindingExpr:
		return c.binding.name.String()
	case *LiteralExpr:
		return "<literal>"
	default:
		return "fn"
	}
}

func (expr *ThrowExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	e, err := Eval(genv, expr.e, env)
	if err != nil {
		return nil, err
	}

	switch sv := e.(type) {
	case Error:
		return nil, sv
	default:
		s, err := e.ToString(genv, false)
		if err != nil {
			return nil, err
		}
		return nil, genv.NewError("Cannot throw " + s)
	}
}

func (expr *TryExpr) Eval(genv *Env, env *LocalEnv) (obj Object, err error) {
	if expr.finallyExpr != nil {
		defer func() {
			_, err = evalBody(genv, expr.finallyExpr, env)
		}()
	}

	obj, err = evalBody(genv, expr.body, env)
	if r, ok := err.(Error); ok {
		for _, catchExpr := range expr.catches {
			if IsInstance(genv, catchExpr.excType, r) {
				obj, err = evalBody(genv, catchExpr.body, env.addFrame([]Object{r}))
				break
			}
		}
	}

	return obj, err
}

func (expr *CatchExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	return nil, genv.NewError("This should never happen!")
}

func evalBody(genv *Env, body []Expr, env *LocalEnv) (Object, error) {
	var res Object = NIL
	var err error
	for _, expr := range body {
		res, err = Eval(genv, expr, env)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func evalLoop(genv *Env, body []Expr, env *LocalEnv) (Object, error) {
	var res Object = NIL
	var err error
loop:
	for _, expr := range body {
		res, err = Eval(genv, expr, env)
		if err != nil {
			return nil, err
		}
	}
	switch res := res.(type) {
	default:
		return res, nil
	case RecurBindings:
		env = env.replaceFrame(res)
		goto loop
	}
}

func (doExpr *DoExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	return evalBody(genv, doExpr.body, env)
}

func (expr *IfExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	v, err := Eval(genv, expr.cond, env)
	if err != nil {
		return nil, err
	}

	if ToBool(v) {
		return Eval(genv, expr.positive, env)
	}
	return Eval(genv, expr.negative, env)
}

func (expr *FnExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	res := &Fn{fnExpr: expr, code: expr.compiled}
	if expr.self != nil {
		env = env.addFrame([]Object{res})
	}
	res.env = env
	return res, nil
}

func (expr *FnArityExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	return nil, genv.NewError("This should never happen!")
}

func (expr *LetExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	env = env.addEmptyFrame(len(expr.names))
	for _, bindingExpr := range expr.values {
		v, err := Eval(genv, bindingExpr, env)
		if err != nil {
			return nil, err
		}
		env.addBinding(v)
	}
	return evalBody(genv, expr.body, env)
}

func (expr *LoopExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	env = env.addEmptyFrame(len(expr.names))
	for _, bindingExpr := range expr.values {
		v, err := Eval(genv, bindingExpr, env)
		if err != nil {
			return nil, err
		}
		env.addBinding(v)
	}
	return evalLoop(genv, expr.body, env)
}

func (expr *RecurExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	v, err := evalSeq(genv, expr.args, env)
	if err != nil {
		return nil, err
	}
	return RecurBindings(v), nil
}

func (expr *MacroCallExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	return expr.macro.Call(genv, expr.args)
}

func (expr *MacroCallExpr) Name() string {
	return expr.name
}

func TryEval(env *Env, expr Expr) (obj Object, err error) {
	return Eval(env, expr, nil)
}

type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("exiting with code: %d", e.Code)
}
