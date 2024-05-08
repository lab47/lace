package core

import "fmt"

type EvalError struct {
	msg  string
	rt   *Runtime
	hash uint32
}

var _ Object = &EvalError{}

func (err *EvalError) ToString(env *Env, escape bool) (string, error) {
	return err.Error(), nil
}

func (err *EvalError) Equals(env *Env, other interface{}) bool {
	return err == other
}

func (err *EvalError) GetInfo() *ObjectInfo {
	return nil
}

func (err *EvalError) GetType() *Type {
	return TYPE.EvalError
}

func (err *EvalError) Hash(env *Env) (uint32, error) {
	return err.hash, nil
}

func (err *EvalError) WithInfo(info *ObjectInfo) Object {
	return err
}

func (err *EvalError) Message() Object {
	return MakeString(err.msg)
}

func (err *EvalError) Error() string {
	pos := err.rt.topPos()
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

func Errorf(env *Env, str string, args ...any) error {
	return env.populateStackTrace(fmt.Errorf(str, args...))
}
