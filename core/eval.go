package core

import (
	"bytes"
	"fmt"
	"strings"
)

type (
	Traceable interface {
		Name() string
		Pos() Position
	}
	Frame struct {
		traceable Traceable
		current   Position
	}
	Callstack struct {
		frames []Frame
	}
	Runtime struct {
		callstack *Callstack
		// GIL         sync.Mutex
	}
)

func (rt *Runtime) clone() *Runtime {
	return &Runtime{
		callstack: rt.callstack.clone(),
	}
}

func StubNewError(msg string) *EvalError {
	res := &EvalError{
		msg: msg,
	}
	return res
}

func (rt *Runtime) NewErrorFlushed(expr Expr, msg string) *EvalError {
	rt.setTopPosition(expr.Pos())

	res := &EvalError{
		msg: msg,
		rt:  rt.clone(),
	}
	return res
}

func (rt *Runtime) NewError(msg string) *EvalError {
	res := &EvalError{
		msg: msg,
		rt:  rt.clone(),
	}
	return res
}

func StubNewArgTypeError(index int, obj Object, expectedType string) *EvalError {
	return StubNewError(fmt.Sprintf("Arg[%d] of <<func_name>> must have type %s, got %s", index, expectedType, obj.GetType().ToString(false)))
}

func (rt *Runtime) NewArgTypeError(index int, obj Object, expectedType string) *EvalError {
	name := rt.topName()
	return rt.NewError(fmt.Sprintf("Arg[%d] of %s must have type %s, got %s", index, name, expectedType, obj.GetType().ToString(false)))
}

func (rt *Runtime) topName() string {
	if len(rt.callstack.frames) == 0 {
		return ""
	}
	fr := &rt.callstack.frames[len(rt.callstack.frames)-1]

	if fr.traceable == nil {
		return ""
	}

	return fr.traceable.Name()
}

func (rt *Runtime) topPos() Position {
	if len(rt.callstack.frames) == 0 {
		return Position{}
	}
	fr := &rt.callstack.frames[len(rt.callstack.frames)-1]

	if fr.traceable == nil {
		return Position{}
	}

	return fr.current
}

func (rt *Runtime) stacktrace() string {
	var b bytes.Buffer
	for _, f := range rt.callstack.frames {
		framePos := f.current
		name := "global"
		if f.traceable != nil {
			name = f.traceable.Name()
		}
		name = strings.TrimPrefix(name, "#'")

		b.WriteString(fmt.Sprintf("  %s %s:%d:%d\n", name, framePos.Filename(), framePos.startLine, framePos.startColumn))
	}
	return b.String()
}

func (rt *Runtime) setTopPosition(pos Position) {
	fr := &rt.callstack.frames[len(rt.callstack.frames)-1]
	fr.current = pos
}

func (rt *Runtime) popFrame() {
	rt.callstack.popFrame()
}

func TopEval(genv *Env, expr Expr, env *LocalEnv) (Object, error) {
	tr, _ := expr.(Traceable)

	genv.RT.callstack.pushFrame(Frame{
		traceable: tr,
		current:   expr.Pos(),
	})

	return expr.Eval(genv, env)
}

func Eval(genv *Env, expr Expr, env *LocalEnv) (Object, error) {
	obj, err := expr.Eval(genv, env)
	if err != nil {
		if ee, ok := err.(*EvalError); ok {
			if ee.rt == nil {
				ee.rt = genv.RT.clone()
			}
		}
	}

	return obj, err
}

func (s *Callstack) pushFrame(frame Frame) {
	s.frames = append(s.frames, frame)
}

func (s *Callstack) popFrame() {
	s.frames = s.frames[:len(s.frames)-1]
}

func (s *Callstack) clone() *Callstack {
	res := &Callstack{frames: make([]Frame, len(s.frames))}
	copy(res.frames, s.frames)
	return res
}

func (s *Callstack) String() string {
	var b bytes.Buffer
	for _, f := range s.frames {
		pos := f.traceable.Pos()
		b.WriteString(fmt.Sprintf("%s %s:%d:%d\n", f.traceable.Name(), pos.Filename(), pos.startLine, pos.startColumn))
	}
	if b.Len() > 0 {
		b.Truncate(b.Len() - 1)
	}
	return b.String()
}

func (expr *VarRefExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	return expr.vr.Resolve(), nil
}

func (expr *SetMacroExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	expr.vr.isMacro = true
	expr.vr.isUsed = false
	if fn, ok := expr.vr.Value.(*Fn); ok {
		fn.isMacro = true
	}
	setMacroMeta(expr.vr)
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
			if res.containsKey(key) {
				return nil, genv.RT.NewErrorFlushed(expr, "Duplicate key: "+key.ToString(false))
			}
			v, err := Eval(genv, expr.values[i], env)
			if err != nil {
				return nil, err
			}
			v, err = res.Assoc(key, v)
			if err != nil {
				return nil, err
			}
			res = v.(*HashMap)
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
		if !res.Add(key, v) {
			return nil, genv.RT.NewErrorFlushed(expr, "Duplicate key: "+key.ToString(false))
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
		ok, err := res.Add(el)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, genv.RT.NewErrorFlushed(expr, "Duplicate set element: "+el.ToString(false))
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
		err := iEval(&expr.vr.Value, genv, expr.value, env)
		if err != nil {
			return nil, err
		}
	}
	meta := EmptyArrayMap()
	meta.Add(criticalKeywords.line, Int{I: expr.startLine})
	meta.Add(criticalKeywords.column, Int{I: expr.startColumn})
	meta.Add(criticalKeywords.file, String{S: *expr.filename})
	meta.Add(criticalKeywords.ns, expr.vr.ns)
	meta.Add(criticalKeywords.name, expr.vr.name)
	expr.vr.meta = meta
	if expr.meta != nil {
		v, err := Eval(genv, expr.meta, env)
		if err != nil {
			return nil, err
		}
		expr.vr.meta, err = expr.vr.meta.Merge(v.(Map))
		if err != nil {
			return nil, err
		}
	}
	// isMacro can be set by set-macro__ during parse stage
	if expr.vr.isMacro {
		v, err := expr.vr.meta.Assoc(criticalKeywords.macro, Boolean{B: true})
		if err != nil {
			return nil, err
		}
		expr.vr.meta = v.(Map)
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
	return res.(Meta).WithMeta(meta.(Map))
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
	genv.RT.setTopPosition(expr.Position)

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

		genv.RT.callstack.pushFrame(Frame{
			traceable: expr,
			current:   expr.callable.Pos(),
		})

		defer genv.RT.popFrame()

		return callable.Call(genv, args)
	default:
		return nil, genv.RT.NewErrorFlushed(expr, callable.ToString(false)+" is not a Fn")
	}
}

func varCallableString(v *Var) string {
	if v.ns.CoreP() {
		return "core/" + v.name.ToString(false)
	}
	return v.ns.Name.ToString(false) + "/" + v.name.ToString(false)
}

func (expr *CallExpr) Name() string {
	switch c := expr.callable.(type) {
	case *VarRefExpr:
		return varCallableString(c.vr)
	case *BindingExpr:
		return c.binding.name.ToString(false)
	case *LiteralExpr:
		return c.obj.ToString(false)
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
		return nil, genv.RT.NewErrorFlushed(expr, "Cannot throw "+e.ToString(false))
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
			if IsInstance(catchExpr.excType, r) {
				obj, err = evalBody(genv, catchExpr.body, env.addFrame([]Object{r}))
				break
			}
		}
	}

	return obj, err
}

func (expr *CatchExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	return nil, genv.RT.NewErrorFlushed(expr, "This should never happen!")
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
	res := &Fn{fnExpr: expr}
	if expr.self.name != nil {
		env = env.addFrame([]Object{res})
	}
	res.env = env
	return res, nil
}

func (expr *FnArityExpr) Eval(genv *Env, env *LocalEnv) (Object, error) {
	return nil, genv.RT.NewErrorFlushed(expr, "This should never happen!")
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
