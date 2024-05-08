package core

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lab47/lace/core/insn"
	"golang.org/x/exp/slices"
)

type MethodSite struct {
	Method string
	Arity  uint
}

type CodeData struct {
	vars        []*Var
	varNames    []Symbol
	defVars     []*Var
	defVarNames []Symbol
	literals    []Object
	codes       []*Code
	methods     []MethodSite

	insns []uint32
}

type BytecodeEncoder struct {
	CodeData
}

func (e *BytecodeEncoder) addVar(vr *Var) uint {
	idx := slices.Index(e.vars, vr)
	if idx != -1 {
		return uint(idx)
	}

	idx = len(e.vars)
	e.vars = append(e.vars, vr)

	sym := AssembleSymbol(vr.ns.Name.Name(), vr.name.String())

	e.varNames = append(e.varNames, sym)

	return uint(idx)
}

func (e *BytecodeEncoder) addDefVar(vr *Var) uint {
	idx := uint(len(e.defVars))
	e.defVars = append(e.defVars, vr)
	e.defVarNames = append(e.defVarNames, vr.name)

	return idx
}

func (e *BytecodeEncoder) addLiteral(obj Object) uint {
	idx := uint(len(e.literals))
	e.literals = append(e.literals, obj)

	return idx
}

func (e *BytecodeEncoder) addCode(c *Code) uint {
	idx := uint(len(e.codes))
	e.codes = append(e.codes, c)

	return idx
}

func (e *BytecodeEncoder) addMethod(mc MethodSite) uint {
	idx := uint(len(e.methods))
	e.methods = append(e.methods, mc)
	return idx
}

func (e *BytecodeEncoder) Encode(i *Instruction) error {
	var a uint

	if i.Data == nil {
		a = uint(i.A0)
	} else {
		switch i.Op {
		case ResolveVar, SetMacro:
			d := i.Data.(*VarData)
			a = e.addVar(d.vr)
		case PushLiteral:
			d := i.Data.(*LiteralData)
			a = e.addLiteral(d.obj)
		case Def, Def3, DefValue, DefValue3:
			d := i.Data.(*VarData)
			a = e.addDefVar(d.vr)
		case CheckType:
			d := i.Data.(*CheckTypeData)
			a = e.addLiteral(d.Type)
		case MakeFn:
			d := i.Data.(*FnData)
			a = e.addCode(d.Code)
		case MethodCall:
			d := i.Data.(*MethodCallData)
			a = e.addMethod(MethodSite{
				Method: d.Method,
				Arity:  uint(i.A0),
			})
		}
	}

	b, err := insn.MakeA(byte(i.Op), a)
	if err != nil {
		return err
	}

	e.insns = append(e.insns, b)

	return nil
}

type EngineFrame struct {
	Ip       int
	Stack    []Object
	SP       int32
	Code     *Fn
	Bindings []Object
	Upvals   []*NamedPair
	Arity    int32
	Args     []Object

	Handlers []handler
}

type Engine struct {
	frames     []*EngineFrame
	allocstack []Object
}

func (e *Engine) printStack(env *Env) {
	var parts []string

	for _, o := range e.allocstack {
		parts = append(parts, fmt.Sprintf("%v", o))
	}

	fmt.Println("[ " + strings.Join(parts, ", ") + " ]")
}

func (f *EngineFrame) stackPush(obj Object) {
	f.SP++
	f.Stack[f.SP] = obj
	//f.Stack = append(f.Stack, obj)
}

func (f *EngineFrame) stackPop() Object {
	if f.SP < 0 {
		panic("stack underflow")
	}

	val := f.Stack[f.SP]
	f.SP--

	return val

	/*
		idx := len(f.Stack) - 1
		if idx == -1 {
			panic("stack underflow")
		}
		val := f.Stack[idx]
		f.Stack = f.Stack[:idx]

		return val
	*/
}

func (f *EngineFrame) stackTop() Object {
	if f.SP < 0 {
		panic("stack underflow")
	}

	return f.Stack[f.SP]

	/*
		idx := len(f.Stack) - 1
		return f.Stack[idx]
	*/
}

func (f *EngineFrame) stackPopN(cnt int) []Object {
	start := int(f.SP) - (cnt - 1)
	if start < 0 {
		panic(fmt.Sprintf("bad top slice request %d (total: %d)", cnt, len(f.Stack)))
	}

	objs := f.Stack[start : f.SP+1]
	f.SP -= int32(cnt)

	return objs

	/*
		idx := len(f.Stack) - cnt
		if idx < 0 {
			panic(fmt.Sprintf("bad top slice request %d (total: %d)", cnt, len(f.Stack)))
		}

		objs := f.Stack[idx:]
		f.Stack = f.Stack[:idx]

		return objs
	*/
}

/*
func (e *Engine) printStack(env *Env) {
	var strs []string

	for _, o := range e.stack {
		str, err := o.ToString(env, false)
		if err == nil {
			strs = append(strs, str)
		} else {
			strs = append(strs, fmt.Sprint(o))
		}
	}

	fmt.Printf("      [ %s ]\n", strings.Join(strs, ", "))
}
*/

func (e *Engine) pushFrame(fn *Fn) *EngineFrame {
	var stack []Object

	stack = make([]Object, fn.code.stackSize)
	idx := len(e.frames)

	upvals := make([]*NamedPair, fn.code.totalUpvals)
	copy(upvals, fn.importedUpvals)

	e.frames = append(e.frames, &EngineFrame{
		Code:     fn,
		Stack:    stack,
		SP:       -1,
		Upvals:   upvals,
		Bindings: make([]Object, fn.code.numBindings),
	})

	return e.frames[idx]
}

func (e *Engine) popFrame() {
	e.frames = e.frames[:len(e.frames)-1]
}

func (e *Engine) frameBack(cnt int) (*EngineFrame, error) {
	idx := len(e.frames) - cnt - 1

	if idx < 0 {
		return nil, fmt.Errorf("invalid upward frame request %d (have %d)", cnt, len(e.frames))
	}

	return e.frames[idx], nil
}

func EngineCode(env *Env, c *Code) (Object, error) {
	return EngineRun(env, &Fn{code: c})
}

func EngineRun(env *Env, fn *Fn) (Object, error) {
	var e Engine
	env.Engine = &e
	e.allocstack = make([]Object, 0, 100)
	e.pushFrame(fn)
	return e.RunBC(env, fn)
}

func (e *Engine) RunWithArgs(env *Env, fn *Fn, args []Object) (Object, error) {
	fr := e.pushFrame(fn)
	fr.Args = slices.Clone(args)
	fr.Arity = int32(len(args))
	return e.RunBC(env, fn)
}

type VMError struct {
	obj Object
}

func (VMError) Error() string {
	return "a VM error"
}

func (e *Engine) unwind(oerr error) (int, error) {
	obj, ok := oerr.(Object)
	if !ok {
		return 0, oerr
	}

	fr, err := e.frameBack(0)
	if err != nil {
		return 0, err
	}

	if len(fr.Handlers) == 0 {
		return 0, oerr
	}

	eh := fr.Handlers[len(fr.Handlers)-1]

	fr.Handlers = fr.Handlers[:len(fr.Handlers)-1]

	fr.SP = int32(eh.sp)

	fr.stackPush(obj)

	return eh.ip, nil
}

func (e *Engine) printBacktrace(last int) {
	if last >= len(e.frames) {
		last = len(e.frames)
	}

	for _, fr := range e.frames[len(e.frames)-last:] {
		fmt.Printf("| %s:%d (ip: %d)\n", fr.Code.code.filename, fr.Code.code.lineForIp(fr.Ip), fr.Ip)
	}
}

func (e *Engine) makeStackTrace() Object {
	var vals []Object

	for i := len(e.frames) - 1; i >= 0; i-- {
		fr := e.frames[i]

		vals = append(vals, NewVectorFrom(
			fr.Code,
			MakeInt(fr.Ip),
		))
		/*
			var name string
			if fr.Code.meta != nil {
				if ok, val := fr.Code.meta.GetEqu(criticalKeywords.name); ok {
					if sym, ok := val.(Symbol); ok {
						name = sym.String()
					}
				}

				if ok, val := fr.Code.meta.GetEqu(criticalKeywords.ns); ok {
					if ns, ok := val.(Symbol); ok {
						name = ns.Name() + "/" + name
					}
				}
			}
			vals = append(vals, MakeString(
				fmt.Sprintf("% 30s %s:%d (ip: %d)", name, fr.Code.code.filename, fr.Code.code.lineForIp(fr.Ip), fr.Ip),
			))
		*/
	}

	return NewListFrom(vals...)
}

func (e *Engine) assembleBacktrace() []string {
	var ret []string

	for _, fr := range e.frames {
		ret = append(ret,
			fmt.Sprintf("%s:%d (ip: %d)", fr.Code.code.filename, fr.Code.code.lineForIp(fr.Ip), fr.Ip),
		)
	}

	return ret
}

func (e *Engine) call(env *Env, callable Callable, objArgs []Object) (Object, error) {
	if fn, ok := callable.(*Fn); ok {
		if fn.code == nil {
			if fn.env != nil && len(fn.env.bindings) > 0 {
				panic("can't do this one")
			}

			_, err := compileFn(env, fn, nil)
			if err != nil {
				return nil, err
			}
		}

		//e.printBacktrace(2)
		//printInsns(fn.code.insns)
		fr := e.pushFrame(fn)
		fr.Args = slices.Clone(objArgs)
		fr.Arity = int32(len(objArgs))

		return e.RunBC(env, fn)
	}

	return callable.Call(env, objArgs)
}

func (e *Engine) methodCall(env *Env, methName string, obj Object, objArgs []Object) (Object, error) {
	var rv reflect.Value

	if orv, ok := obj.(*ReflectValue); ok {
		rv = orv.val
	} else {
		rv = reflect.ValueOf(obj)
	}

	rt := rv.Type()

	meth := rv.MethodByName(methName)

	if !meth.IsValid() {
		return nil, env.RT.NewError(fmt.Sprintf("unknown method %s on %s", methName, rt))
	}

	procFn, _, err := convReg.ConverterForFunc(meth)
	if err != nil {
		return nil, err
	}

	return procFn(env, objArgs)
}

const debugBC = false

func (e *Engine) RunBC(env *Env, fn *Fn) (Object, error) {
	c := fn.code.data

	frame, err := e.frameBack(0)
	if err != nil {
		return nil, err
	}

	if debugBC {
		fmt.Printf("==== enter frame %d %s:%d (upvals: %d) =====\n", len(e.frames), fn.code.filename, fn.code.lineForIp(0), len(frame.Upvals))
		defer fmt.Printf("==== exit frame %d =====\n", len(e.frames))
	}

	defer e.popFrame()

	var (
		ip  int
		tmp Object
	)

	defer func() {
		frame.Ip = ip
	}()

loop:
	for {
		//nsn := c.insns[ip]

		op, a := insn.Decode(c.insns[ip])

		if debugBC {
			fmt.Printf("% 2d|% 4d|% 4d| %s %d\n", len(e.frames), ip, fn.code.lineForIp(ip), OpCode(op), a)
			//e.printStack(env)
		}

		switch OpCode(op) {
		case Pop:
			frame.stackPop()
		case Dup:
			frame.stackPush(frame.stackTop())
		case Return:
			return frame.stackPop(), nil
		case Jump:
			ip = int(a)
			continue loop
		case JumpIfTrue:
			if ToBool(frame.stackPop()) {
				ip = int(a)
				continue loop
			}
		case JumpIfFalse:
			if !ToBool(frame.stackPop()) {
				ip = int(a)
				continue loop
			}
		case GetUpval:
			frame.stackPush(frame.Upvals[a].Value)
		case RefUpval:
			frame.stackPush(frame.Upvals[a])
		case SetUpval:
			uv := frame.Upvals[a]
			if uv == nil {
				uv = &NamedPair{}
				frame.Upvals[a] = uv
			}

			uv.Value = frame.stackPop()
		case ResolveVar:
			tmp = c.vars[a].Resolve()
			if debugBC {
				fmt.Printf("             | %s\n", c.vars[a].Name())
			}
			frame.stackPush(tmp)
		case SetMacro:
			vr := c.vars[a]
			vr.isMacro = true
			vr.isUsed = false
			if fn, ok := vr.Value.(*Fn); ok {
				fn.isMacro = true
			}
			err := setMacroMeta(env, vr)
			if err != nil {
				return nil, err
			}
			frame.stackPush(vr)
		case GetLocal:
			frame.stackPush(frame.Bindings[a])
		case SetLocal:
			frame.Bindings[a] = frame.stackPop()
		case PushLiteral:
			frame.stackPush(c.literals[a])
		case PushSelfFn:
			frame.stackPush(frame.Code)
		case PushInt:
			frame.stackPush(MakeInt(int(a)))
		case PushNil:
			frame.stackPush(NIL)
		case MakeVector:
			vec := NewVectorFrom(frame.stackPopN(int(a))...)

			frame.stackPush(vec)
		case MakeLargeMap:
			data := frame.stackPopN(int(a))

			res := EmptyHashMap
			for i := 0; i < len(data); i += 2 {
				key := data[i]
				val := data[i+1]

				if res.containsKey(env, key) {
					s, err := key.ToString(env, false)
					if err != nil {
						return nil, err
					}
					return nil, env.RT.NewError("Duplicate key: " + s)
				}

				up, err := res.Assoc(env, key, val)
				if err != nil {
					return nil, err
				}

				if err := Cast(env, up, &res); err != nil {
					return nil, err
				}
			}

			frame.stackPush(res)
		case MakeSmallMap:
			cnt := a
			res := EmptyArrayMap()

			if cnt > 0 {
				data := frame.stackPopN(int(cnt))

				for i := 0; i < len(data); i += 2 {
					key := data[i]
					val := data[i+1]

					if !res.Add(env, key, val) {
						s, err := key.ToString(env, false)
						if err != nil {
							return nil, err
						}

						return nil, env.RT.NewError("Duplicate key: " + s)
					}
				}
			}

			frame.stackPush(res)
		case MakeSet:
			data := frame.stackPopN(int(a))

			res := EmptySet()

			for i := 0; i < len(data); i++ {
				ele := data[i]

				ok, err := res.Add(env, ele)
				if err != nil {
					return nil, err
				}

				if !ok {
					s, err := ele.ToString(env, false)
					if err != nil {
						return nil, err
					}

					return nil, env.RT.NewError("Duplicate set element: " + s)
				}
			}

			frame.stackPush(res)
		case Def:
			vr := c.defVars[a]

			basic := frame.stackPop()
			if err := Cast(env, basic, &vr.meta); err != nil {
				return nil, err
			}

			// isMacro can be set by set-macro__ during parse stage
			if vr.isMacro {
				v, err := vr.meta.Assoc(env, criticalKeywords.macro, Boolean{B: true})
				if err != nil {
					return nil, err
				}
				var m Map
				if err := Cast(env, v, &m); err != nil {
					return nil, err
				}
				vr.meta = m
			}

			frame.stackPush(vr)
		case Def3:
			vr := c.defVars[a]

			v := frame.stackPop()

			if err := Cast(env, frame.stackPop(), &vr.meta); err != nil {
				return nil, err
			}

			var m Map
			if err := Cast(env, v, &m); err != nil {
				return nil, err
			}
			vr.meta, err = vr.meta.Merge(env, m)
			if err != nil {
				return nil, err
			}

			// isMacro can be set by set-macro__ during parse stage
			if vr.isMacro {
				v, err := vr.meta.Assoc(env, criticalKeywords.macro, Boolean{B: true})
				if err != nil {
					return nil, err
				}
				var m Map
				if err := Cast(env, v, &m); err != nil {
					return nil, err
				}
				vr.meta = m
			}

			frame.stackPush(vr)
		case DefValue:
			vr := c.defVars[a]

			val := frame.stackPop()

			if err := Cast(env, frame.stackPop(), &vr.meta); err != nil {
				return nil, err
			}

			vr.Value = val

			// isMacro can be set by set-macro__ during parse stage
			if vr.isMacro {
				v, err := vr.meta.Assoc(env, criticalKeywords.macro, Boolean{B: true})
				if err != nil {
					return nil, err
				}
				var m Map
				if err := Cast(env, v, &m); err != nil {
					return nil, err
				}
				vr.meta = m
			}

			if m, ok := val.(*Fn); ok {
				m.meta = vr.meta
			}

			frame.stackPush(vr)
		case DefValue3:
			vr := c.defVars[a]
			v := frame.stackPop()

			if err := Cast(env, frame.stackPop(), &vr.meta); err != nil {
				return nil, err
			}

			val := frame.stackPop()
			vr.Value = val

			var m Map
			if err := Cast(env, v, &m); err != nil {
				return nil, err
			}

			vr.meta, err = vr.meta.Merge(env, m)
			if err != nil {
				return nil, err
			}

			// isMacro can be set by set-macro__ during parse stage
			if vr.isMacro {
				v, err := vr.meta.Assoc(env, criticalKeywords.macro, Boolean{B: true})
				if err != nil {
					return nil, err
				}
				var m Map
				if err := Cast(env, v, &m); err != nil {
					return nil, err
				}
				vr.meta = m
			}

			if m, ok := val.(*Fn); ok {
				m.meta = vr.meta
			}

			frame.stackPush(vr)
		case SetMeta:
			res := frame.stackPop()
			meta := frame.stackPop()

			var metao Meta
			if err := Cast(env, res, &metao); err != nil {
				return nil, err
			}

			var m Map
			if err := Cast(env, meta, &m); err != nil {
				return nil, err
			}

			mo, err := metao.WithMeta(env, m)
			if err != nil {
				return nil, err
			}

			frame.stackPush(mo)
		case Call:
			args := frame.stackPopN(int(a))
			obj := frame.stackPop()

			frame.Ip = ip

			switch callable := obj.(type) {
			case Callable:
				if fn, ok := callable.(*Fn); ok {
					fr := e.pushFrame(fn)
					fr.Args = slices.Clone(args)
					fr.Arity = int32(len(args))

					obj, err = e.RunBC(env, fn)
				} else {
					obj, err = callable.Call(env, args)
				}

				if err != nil {
					newIp, err := e.unwind(err)
					if err != nil {
						return nil, err
					}

					ip = newIp
					continue loop
				}

				frame.stackPush(obj)
			default:
				e.printBacktrace(1000)
				s, err := callable.ToString(env, false)
				if err != nil {
					return nil, err
				}

				return nil, env.RT.NewError(s + " is not a Fn")
			}
		case Apply:
			frame.Ip = ip

			args := frame.stackPop()
			obj := frame.stackPop()

			seqable, ok := args.(Seqable)
			if !ok {
				newIp, err := e.unwind(env.RT.NewArgTypeError(1, args, "Seqable"))
				if err != nil {
					return nil, err
				}

				ip = newIp
				continue loop
			}

			sq := seqable.Seq()

			callArgs, err := ToSlice(env, sq)
			if err != nil {
				newIp, err := e.unwind(env.RT.NewArgTypeError(1, args, "Seqable"))
				if err != nil {
					return nil, err
				}

				ip = newIp
				continue loop
			}

			frame.Ip = ip

			switch callable := obj.(type) {
			case Callable:
				if fn, ok := callable.(*Fn); ok {
					fr := e.pushFrame(fn)
					fr.Args = slices.Clone(callArgs)
					fr.Arity = int32(len(callArgs))

					obj, err = e.RunBC(env, fn)
				} else {

					obj, err = callable.Call(env, callArgs)
				}

				if err != nil {
					newIp, err := e.unwind(err)
					if err != nil {
						return nil, err
					}

					ip = newIp
					continue loop
				}

				frame.stackPush(obj)
			default:
				e.printBacktrace(1000)
				s, err := callable.ToString(env, false)
				if err != nil {
					return nil, err
				}

				return nil, env.RT.NewError(s + " is not a Fn")
			}
		case MethodCall:
			ms := c.methods[a]

			args := frame.stackPopN(int(a))
			obj := frame.stackPop()

			frame.Ip = ip

			res, err := e.methodCall(env, ms.Method, obj, args)
			if err != nil {
				newIp, err := e.unwind(err)
				if err != nil {
					return nil, err
				}

				ip = newIp
				continue loop
			}

			frame.stackPush(res)
		case Throw:
			frame.Ip = ip

			newIp, err := e.unwind(&VMError{obj: frame.stackPop()})
			if err != nil {
				return nil, err
			}

			ip = newIp
			continue loop

		case PushHandler:
			frame.Handlers = append(frame.Handlers, handler{
				ip: int(a),
				sp: int(frame.SP),
			})

		case PopHandler:
			frame.Handlers = frame.Handlers[:len(frame.Handlers)-1]
		case MakeFn:
			code := c.codes[a]

			upvals := make([]*NamedPair, code.importUpvals)

			imports := frame.stackPopN(code.importUpvals)

			for i, o := range imports {
				np, ok := o.(*NamedPair)
				if !ok {
					newIp, err := e.unwind(fmt.Errorf("value should have been NamedPair for upval"))
					if err != nil {
						return nil, err
					}

					ip = newIp
					continue loop
				}

				upvals[i] = np
			}

			fn := &Fn{
				code:           code,
				importedUpvals: upvals,
			}

			frame.stackPush(fn)
		case CheckArityFixed:
			if frame.Arity == int32(a) {
				copy(frame.Bindings, frame.Args)
				frame.stackPush(MakeBoolean(true))
			} else {
				frame.stackPush(MakeBoolean(false))
			}
		case CheckArityMin:
			if frame.Arity >= int32(a) {
				copy(frame.Bindings, frame.Args[:a])

				if frame.Arity == int32(a) {
					frame.Bindings[a] = NIL
				} else {
					frame.Bindings[a] = &ArraySeq{arr: frame.Args, index: int(a)}
				}

				frame.stackPush(MakeBoolean(true))
			} else {
				frame.stackPush(MakeBoolean(false))
			}
		case ThrowArity:
			frame.Ip = ip
			return nil, ErrorArity(env, int(frame.Arity))
		default:
			return nil, fmt.Errorf("unimplemented instruction: %d", op)
		}

		ip++
	}
}
