package core

import (
	"bytes"
	"fmt"
	"strings"
	"unsafe"
)

type (
	Traceable interface {
		Name() string
		Pos() Position
	}
	EvalError struct {
		msg  string
		pos  Position
		rt   *Runtime
		hash uint32
	}
	Frame struct {
		traceable Traceable
	}
	Callstack struct {
		frames []Frame
	}
	Runtime struct {
		callstack   *Callstack
		currentExpr Expr
		// GIL         sync.Mutex
	}
)

func init() {
	GLOBAL_ENV.RT = &Runtime{
		callstack: &Callstack{frames: make([]Frame, 0, 50)},
	}
}

func (rt *Runtime) clone() *Runtime {
	return &Runtime{
		callstack:   rt.callstack.clone(),
		currentExpr: rt.currentExpr,
	}
}

func StubNewError(msg string) *EvalError {
	res := &EvalError{
		msg: msg,
	}
	return res
}

func (rt *Runtime) NewError(msg string) *EvalError {
	res := &EvalError{
		msg: msg,
		rt:  rt.clone(),
	}
	if rt.currentExpr != nil {
		res.pos = rt.currentExpr.Pos()
	}
	return res
}

func StubNewArgTypeError(index int, obj Object, expectedType string) *EvalError {
	return StubNewError(fmt.Sprintf("Arg[%d] of <<func_name>> must have type %s, got %s", index, expectedType, obj.GetType().ToString(false)))
}

func (rt *Runtime) NewArgTypeError(index int, obj Object, expectedType string) *EvalError {
	name := rt.currentExpr.(Traceable).Name()
	return rt.NewError(fmt.Sprintf("Arg[%d] of %s must have type %s, got %s", index, name, expectedType, obj.GetType().ToString(false)))
}

func StubNewErrorWithPos(msg string, pos Position) *EvalError {
	return &EvalError{
		msg: msg,
		pos: pos,
	}
}

func (rt *Runtime) NewErrorWithPos(msg string, pos Position) *EvalError {
	return &EvalError{
		msg: msg,
		pos: pos,
		rt:  rt.clone(),
	}
}

func (rt *Runtime) stacktrace() string {
	var b bytes.Buffer
	pos := Position{}
	if rt.currentExpr != nil {
		pos = rt.currentExpr.Pos()
	}
	name := "global"
	for _, f := range rt.callstack.frames {
		framePos := f.traceable.Pos()
		b.WriteString(fmt.Sprintf("  %s %s:%d:%d\n", name, framePos.Filename(), framePos.startLine, framePos.startColumn))
		name = f.traceable.Name()
		if strings.HasPrefix(name, "#'") {
			name = name[2:]
		}
	}
	b.WriteString(fmt.Sprintf("  %s %s:%d:%d", name, pos.Filename(), pos.startLine, pos.startColumn))
	return b.String()
}

func (rt *Runtime) pushFrame() {
	// TODO: this is all wrong. We cannot rely on
	// currentExpr for stacktraces. Instead, each Callable
	// should know it's name / position.
	var tr Traceable
	if rt.currentExpr != nil {
		tr = rt.currentExpr.(Traceable)
	} else {
		tr = &CallExpr{}
	}
	rt.callstack.pushFrame(Frame{traceable: tr})
}

func (rt *Runtime) popFrame() {
	rt.callstack.popFrame()
}

func Eval(genv *Env, expr Expr, env *LocalEnv) Object {
	parentExpr := genv.RT.currentExpr
	genv.RT.currentExpr = expr
	defer (func() { genv.RT.currentExpr = parentExpr })()
	return expr.Eval(genv, env)
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

func MakeEvalError(msg string, pos Position, rt *Runtime) *EvalError {
	res := &EvalError{msg, pos, rt, 0}
	res.hash = HashPtr(uintptr(unsafe.Pointer(res)))
	return res
}

func (err *EvalError) ToString(escape bool) string {
	return err.Error()
}

func (err *EvalError) Equals(other interface{}) bool {
	return err == other
}

func (err *EvalError) GetInfo() *ObjectInfo {
	return nil
}

func (err *EvalError) GetType() *Type {
	return TYPE.EvalError
}

func (err *EvalError) Hash() uint32 {
	return err.hash
}

func (err *EvalError) WithInfo(info *ObjectInfo) Object {
	return err
}

func (err *EvalError) Message() Object {
	return MakeString(err.msg)
}

func (err *EvalError) Error() string {
	pos := err.pos
	if err.rt == nil {
		return fmt.Sprintf("%s:%d:%d: Eval error: %s", pos.Filename(), pos.startLine, pos.startColumn, err.msg)
	}

	if len(err.rt.callstack.frames) > 0 && !LINTER_MODE {
		return fmt.Sprintf("%s:%d:%d: Eval error: %s\nStacktrace:\n%s", pos.Filename(), pos.startLine, pos.startColumn, err.msg, err.rt.stacktrace())
	} else {
		if len(err.rt.callstack.frames) > 0 {
			pos = err.rt.callstack.frames[0].traceable.Pos()
		}
		return fmt.Sprintf("%s:%d:%d: Eval error: %s", pos.Filename(), pos.startLine, pos.startColumn, err.msg)
	}
}

func (expr *VarRefExpr) Eval(genv *Env, env *LocalEnv) Object {
	return expr.vr.Resolve()
}

func (expr *SetMacroExpr) Eval(genv *Env, env *LocalEnv) Object {
	expr.vr.isMacro = true
	expr.vr.isUsed = false
	if fn, ok := expr.vr.Value.(*Fn); ok {
		fn.isMacro = true
	}
	setMacroMeta(expr.vr)
	return expr.vr
}

func (expr *BindingExpr) Eval(genv *Env, env *LocalEnv) Object {
	for i := env.frame; i > expr.binding.frame; i-- {
		env = env.parent
	}
	return env.bindings[expr.binding.index]
}

func (expr *LiteralExpr) Eval(genv *Env, env *LocalEnv) Object {
	return expr.obj
}

func (expr *VectorExpr) Eval(genv *Env, env *LocalEnv) Object {
	res := EmptyVector()
	for _, e := range expr.v {
		res = res.Conjoin(Eval(genv, e, env))
	}
	return res
}

func (expr *MapExpr) Eval(genv *Env, env *LocalEnv) Object {
	if int64(len(expr.keys)) > HASHMAP_THRESHOLD/2 {
		res := EmptyHashMap
		for i := range expr.keys {
			key := Eval(genv, expr.keys[i], env)
			if res.containsKey(key) {
				panic(genv.RT.NewError("Duplicate key: " + key.ToString(false)))
			}
			res = res.Assoc(key, Eval(genv, expr.values[i], env)).(*HashMap)
		}
		return res
	}
	res := EmptyArrayMap()
	for i := range expr.keys {
		key := Eval(genv, expr.keys[i], env)
		if !res.Add(key, Eval(genv, expr.values[i], env)) {
			panic(genv.RT.NewError("Duplicate key: " + key.ToString(false)))
		}
	}
	return res
}

func (expr *SetExpr) Eval(genv *Env, env *LocalEnv) Object {
	res := EmptySet()
	for _, elemExpr := range expr.elements {
		el := Eval(genv, elemExpr, env)
		if !res.Add(el) {
			panic(genv.RT.NewError("Duplicate set element: " + el.ToString(false)))
		}
	}
	return res
}

func (expr *DefExpr) Eval(genv *Env, env *LocalEnv) Object {
	if expr.value != nil {
		expr.vr.Value = Eval(genv, expr.value, env)
	}
	meta := EmptyArrayMap()
	meta.Add(KEYWORDS.line, Int{I: expr.startLine})
	meta.Add(KEYWORDS.column, Int{I: expr.startColumn})
	meta.Add(KEYWORDS.file, String{S: *expr.filename})
	meta.Add(KEYWORDS.ns, expr.vr.ns)
	meta.Add(KEYWORDS.name, expr.vr.name)
	expr.vr.meta = meta
	if expr.meta != nil {
		expr.vr.meta = expr.vr.meta.Merge(Eval(genv, expr.meta, env).(Map))
	}
	// isMacro can be set by set-macro__ during parse stage
	if expr.vr.isMacro {
		expr.vr.meta = expr.vr.meta.Assoc(KEYWORDS.macro, Boolean{B: true}).(Map)
	}
	return expr.vr
}

func (expr *MetaExpr) Eval(genv *Env, env *LocalEnv) Object {
	meta := Eval(genv, expr.meta, env)
	res := Eval(genv, expr.expr, env)
	return res.(Meta).WithMeta(meta.(Map))
}

func evalSeq(genv *Env, exprs []Expr, env *LocalEnv) []Object {
	res := make([]Object, len(exprs))
	for i, expr := range exprs {
		res[i] = Eval(genv, expr, env)
	}
	return res
}

func (expr *CallExpr) Eval(genv *Env, env *LocalEnv) Object {
	callable := Eval(genv, expr.callable, env)
	switch callable := callable.(type) {
	case Callable:
		args := evalSeq(genv, expr.args, env)
		return callable.Call(genv, args)
	default:
		panic(genv.RT.NewErrorWithPos(callable.ToString(false)+" is not a Fn", expr.callable.Pos()))
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

func (expr *ThrowExpr) Eval(genv *Env, env *LocalEnv) Object {
	e := Eval(genv, expr.e, env)
	switch e.(type) {
	case Error:
		panic(e)
	default:
		panic(genv.RT.NewError("Cannot throw " + e.ToString(false)))
	}
}

func (expr *TryExpr) Eval(genv *Env, env *LocalEnv) (obj Object) {
	defer func() {
		defer func() {
			if expr.finallyExpr != nil {
				evalBody(genv, expr.finallyExpr, env)
			}
		}()
		if r := recover(); r != nil {
			switch r := r.(type) {
			case Error:
				for _, catchExpr := range expr.catches {
					if IsInstance(catchExpr.excType, r) {
						obj = evalBody(genv, catchExpr.body, env.addFrame([]Object{r}))
						return
					}
				}
				panic(r)
			default:
				panic(r)
			}
		}
	}()
	return evalBody(genv, expr.body, env)
}

func (expr *CatchExpr) Eval(genv *Env, env *LocalEnv) (obj Object) {
	panic(genv.RT.NewError("This should never happen!"))
}

func evalBody(genv *Env, body []Expr, env *LocalEnv) Object {
	var res Object = NIL
	for _, expr := range body {
		res = Eval(genv, expr, env)
	}
	return res
}

func evalLoop(genv *Env, body []Expr, env *LocalEnv) Object {
	var res Object = NIL
loop:
	for _, expr := range body {
		res = Eval(genv, expr, env)
	}
	switch res := res.(type) {
	default:
		return res
	case RecurBindings:
		env = env.replaceFrame(res)
		goto loop
	}
}

func (doExpr *DoExpr) Eval(genv *Env, env *LocalEnv) Object {
	return evalBody(genv, doExpr.body, env)
}

func (expr *IfExpr) Eval(genv *Env, env *LocalEnv) Object {
	if ToBool(Eval(genv, expr.cond, env)) {
		return Eval(genv, expr.positive, env)
	}
	return Eval(genv, expr.negative, env)
}

func (expr *FnExpr) Eval(genv *Env, env *LocalEnv) Object {
	res := &Fn{fnExpr: expr}
	if expr.self.name != nil {
		env = env.addFrame([]Object{res})
	}
	res.env = env
	return res
}

func (expr *FnArityExpr) Eval(genv *Env, env *LocalEnv) Object {
	panic(genv.RT.NewError("This should never happen!"))
}

func (expr *LetExpr) Eval(genv *Env, env *LocalEnv) Object {
	env = env.addEmptyFrame(len(expr.names))
	for _, bindingExpr := range expr.values {
		env.addBinding(Eval(genv, bindingExpr, env))
	}
	return evalBody(genv, expr.body, env)
}

func (expr *LoopExpr) Eval(genv *Env, env *LocalEnv) Object {
	env = env.addEmptyFrame(len(expr.names))
	for _, bindingExpr := range expr.values {
		env.addBinding(Eval(genv, bindingExpr, env))
	}
	return evalLoop(genv, expr.body, env)
}

func (expr *RecurExpr) Eval(genv *Env, env *LocalEnv) Object {
	return RecurBindings(evalSeq(genv, expr.args, env))
}

func (expr *MacroCallExpr) Eval(genv *Env, env *LocalEnv) Object {
	return expr.macro.Call(genv, expr.args)
}

func (expr *MacroCallExpr) Name() string {
	return expr.name
}

func TryEval(env *Env, expr Expr) (obj Object, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case *EvalError:
				err = r.(error)
			case *ExInfo:
				err = r.(error)
			default:
				panic(r)
			}
		}
	}()

	return Eval(env, expr, nil), nil
}

func PanicOnErr(err error) {
	if err != nil {
		panic(StubNewError(err.Error()))
	}
}
