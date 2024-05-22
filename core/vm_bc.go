package core

import (
	"errors"
	"fmt"
	"reflect"

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
	literals    []any
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

func (e *BytecodeEncoder) addLiteral(obj any) uint {
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
	Ip        int
	Stack     []any
	SP        int32
	Code      *Fn
	Bindings  []any
	Upvals    []any
	Arity     int32
	Args      []any
	FrameSize int

	Handlers []handler
}

type FrameRope struct {
	curChunk []EngineFrame
	chunks   [][]EngineFrame
	chunkPtr int
	curPtr   int
	total    int
}

func (f *FrameRope) init() {
	f.curChunk = make([]EngineFrame, 50)
	f.chunks = append(f.chunks, f.curChunk)
	f.chunkPtr = 0
	f.curPtr = -1
}

func (f *FrameRope) pushFrame() *EngineFrame {
	chunk := f.curChunk
	f.curPtr++
	pos := f.curPtr

	f.total++

	if len(chunk) <= pos {
		chunk = make([]EngineFrame, 50)
		f.chunks = append(f.chunks, chunk)
		f.curChunk = chunk
		f.chunkPtr++
		f.curPtr = 0
		pos = 0
	}

	return &chunk[pos]
}

func (f *FrameRope) popFrame() {
	f.total--
	if f.curPtr == 0 {
		if f.chunkPtr > 0 {
			f.chunkPtr--
			f.curChunk = f.chunks[f.chunkPtr]
			f.curPtr = len(f.chunks[f.chunkPtr]) - 1
		} else {
			f.curPtr = -1
		}
	} else {
		f.curPtr--
	}
}

func (f *FrameRope) frameTop() *EngineFrame {
	return &f.chunks[f.chunkPtr][f.curPtr]
}

type Engine struct {
	frope      FrameRope
	allocstack []any
	stackTop   int
}

/*
func (e *Engine) printStack(_ *Env) {
	var parts []string

	for _, o := range e.allocstack {
		parts = append(parts, fmt.Sprintf("%v", o))
	}

	fmt.Println("[ " + strings.Join(parts, ", ") + " ]")
}
*/

func (f *EngineFrame) stackPush(obj any) {
	f.SP++
	f.Stack[f.SP] = obj
	//f.Stack = append(f.Stack, obj)
}

func (f *EngineFrame) stackPop() any {
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

func (f *EngineFrame) stackTop() any {
	if f.SP < 0 {
		panic("stack underflow")
	}

	return f.Stack[f.SP]

	/*
		idx := len(f.Stack) - 1
		return f.Stack[idx]
	*/
}

func (f *EngineFrame) stackPopN(cnt int) []any {
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
		str, err := ToString(env, o)
		if err == nil {
			strs = append(strs, str)
		} else {
			strs = append(strs, fmt.Sprint(o))
		}
	}

	fmt.Printf("      [ %s ]\n", strings.Join(strs, ", "))
}
*/

const defStackSize = 8000

var ErrStackOverflow = errors.New("stack overflow, not enough stack room")

func (e *Engine) pushFrame(fn *Fn, args []any) (*EngineFrame, error) {
	var stack []any

	remaining := len(e.allocstack) - e.stackTop

	stackNeeded := fn.code.stackSize + uint32(fn.code.numBindings) + uint32(fn.code.totalUpvals)
	if remaining < int(stackNeeded) {
		return nil, ErrStackOverflow
	}

	frameSpace := e.allocstack[e.stackTop : e.stackTop+int(stackNeeded)]

	// Divide up our slice of the stack into the actual stack, the slice of
	// bindings the function will use, and the slice of upvals that are available
	off := 0
	stack = frameSpace[:fn.code.stackSize]
	off += int(fn.code.stackSize)
	bindings := frameSpace[off : off+fn.code.numBindings]
	off += fn.code.numBindings
	upvals := frameSpace[off : off+fn.code.totalUpvals]

	e.stackTop += int(stackNeeded)

	// NOTE: we don't currently clear the stack because the code can't use
	// stack entries that are beyond the SP. But if you're directly inspecting
	// the stack, you CAN see values from previous invocations.
	// NOTE: we could probably get away with not clearing bindings as
	// the generated code always sets a binding before it's read (obvi),
	// but it's possible that a future compiler could have a bug where it did
	// do this and it would be better that it observe a nil than a random object
	// in that case.
	clear(bindings)
	clear(upvals)
	for i, uv := range fn.importedUpvals {
		upvals[i] = uv
	}

	frame := e.frope.pushFrame()

	frame.Code = fn
	frame.Stack = stack
	frame.SP = -1
	frame.Ip = 0
	frame.Upvals = upvals
	frame.Bindings = bindings
	frame.Handlers = nil
	frame.Args = args
	frame.Arity = int32(len(args))
	frame.FrameSize = int(stackNeeded)

	return frame, nil
}

func (e *Engine) popFrame(env *Env, fr *EngineFrame) {
	if v := recover(); v != nil {
		if _, ok := v.(*EvalError); ok {
			panic(v)
		}

		panic(env.NewError(fmt.Sprint(v)))
	}
	e.stackTop -= fr.FrameSize
	e.frope.popFrame()
}

func (e *Engine) frameTop() (*EngineFrame, error) {
	return e.frope.frameTop(), nil
}

func NewEngine() *Engine {
	eng := &Engine{
		allocstack: make([]any, defStackSize),
	}
	eng.frope.init()

	return eng
}

func EngineRun(env *Env, fn *Fn) (any, error) {
	e := env.Engine
	if e == nil {
		e = NewEngine()
		env.Engine = e
	}

	_, err := e.pushFrame(fn, nil)
	if err != nil {
		return nil, err
	}

	return e.RunBC(env, fn)
}

func (e *Engine) RunWithArgs(env *Env, fn *Fn, args []any) (any, error) {
	_, err := e.pushFrame(fn, args)
	if err != nil {
		return nil, err
	}
	return e.RunBC(env, fn)
}

func (e *Engine) unwind(oerr error) (int, error) {
	obj, ok := oerr.(any)
	if !ok {
		return 0, oerr
	}

	fr, err := e.frameTop()
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

func (e *Engine) frames() []*EngineFrame {
	var ret []*EngineFrame

	for i := 0; i < e.frope.chunkPtr; i++ {
		chunk := e.frope.chunks[i]

		for i := range chunk {
			ret = append(ret, &chunk[i])
		}
	}

	for i := 0; i <= e.frope.curPtr; i++ {
		ret = append(ret, &e.frope.curChunk[i])
	}

	return ret
}

func (e *Engine) printBacktrace(last int) {
	frames := e.frames()

	if last >= len(frames) {
		last = len(frames)
	}

	for _, fr := range frames[len(frames)-last:] {
		fmt.Printf("| %s:%d (ip: %d)\n", fr.Code.code.filename, fr.Code.code.lineForIp(fr.Ip), fr.Ip)
	}
}

func (e *Engine) makeStackTrace() any {
	var vals []any

	frames := e.frames()
	for i := len(frames) - 1; i >= 0; i-- {
		fr := frames[i]

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

func (e *Engine) methodCall(env *Env, methName string, obj any, objArgs []any) (any, error) {
	rv := reflect.ValueOf(obj)

	rt := rv.Type()

	meth := rv.MethodByName(methName)

	if !meth.IsValid() {
		return nil, env.NewError(fmt.Sprintf("unknown method %s on %s", methName, rt))
	}

	procFn, _, err := convReg.ConverterForFunc(meth)
	if err != nil {
		return nil, err
	}

	return procFn(env, objArgs)
}

const debugBC = false

func (e *Engine) RunBC(env *Env, fn *Fn) (any, error) {
	c := fn.code.data

	frame, err := e.frameTop()
	if err != nil {
		return nil, err
	}

	var idx int

	if debugBC {
		idx = e.frope.total
		fmt.Printf("==== enter frame %d %s:%d (upvals: %d) =====\n", idx, fn.code.filename, fn.code.lineForIp(0), len(frame.Upvals))
		defer fmt.Printf("==== exit frame %d =====\n", idx)
	}

	defer e.popFrame(env, frame)

loop:
	for {
		//nsn := c.insns[frame.Ip]

		op, a := insn.Decode(c.insns[frame.Ip])

		if debugBC {
			fmt.Printf("% 3d|% 4d|% 4d| %s %d\n", idx, frame.Ip, fn.code.lineForIp(frame.Ip), OpCode(op), a)
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
			frame.Ip = int(a)
			continue loop
		case JumpIfTrue:
			if ToBool(frame.stackPop()) {
				frame.Ip = int(a)
				continue loop
			}
		case JumpIfFalse:
			if !ToBool(frame.stackPop()) {
				frame.Ip = int(a)
				continue loop
			}
		case GetUpval:
			frame.stackPush(frame.Upvals[a].(*NamedPair).Value)
		case RefUpval:
			frame.stackPush(frame.Upvals[a])
		case SetUpval:
			uv := frame.Upvals[a]
			if uv == nil {
				uv = &NamedPair{}
				frame.Upvals[a] = uv
			}

			uv.(*NamedPair).Value = frame.stackPop()
		case ResolveVar:
			tmp := c.vars[a].Resolve(env)
			if debugBC {
				fmt.Printf("             | %s\n", c.vars[a].Name())
			}
			frame.stackPush(tmp)
		case SetMacro:
			vr := c.vars[a]
			vr.isMacro = true
			vr.isUsed = false
			if fn, ok := vr.GetStatic().(*Fn); ok {
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
					s, err := ToString(env, key)
					if err != nil {
						return nil, err
					}
					return nil, env.NewError("Duplicate key: " + s)
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
						s, err := ToString(env, key)
						if err != nil {
							return nil, err
						}

						return nil, env.NewError("Duplicate key: " + s)
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
					s, err := ToString(env, ele)
					if err != nil {
						return nil, err
					}

					return nil, env.NewError("Duplicate set element: " + s)
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
				v, err := vr.meta.Assoc(env, criticalKeywords.macro, Boolean(true))
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
				v, err := vr.meta.Assoc(env, criticalKeywords.macro, Boolean(true))
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

			meta := frame.stackPop()
			val := frame.stackPop()

			if err := Cast(env, meta, &vr.meta); err != nil {
				return nil, AddContext(env, err, "casting value for Var meta: %s", vr.Name())
			}

			vr.SetStatic(val)

			// isMacro can be set by set-macro__ during parse stage
			if vr.isMacro {
				v, err := vr.meta.Assoc(env, criticalKeywords.macro, Boolean(true))
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
			vr.SetStatic(val)

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
				v, err := vr.meta.Assoc(env, criticalKeywords.macro, Boolean(true))
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

			switch sv := obj.(type) {
			case *Fn:
				_, ferr := e.pushFrame(sv, args)
				if ferr != nil {
					return nil, err
				}

				obj, err = e.RunBC(env, sv)
			case Proc:
				obj, err = sv.Fn(env, args)
			case Callable:
				obj, err = sv.Call(env, args)
			default:
				err = Errorf(env, "value is not callable: %T", sv)
			}

			if err != nil {
				err = env.populateStackTrace(err)
				newIp, err := e.unwind(err)
				if err != nil {
					return nil, err
				}

				frame.Ip = newIp
				continue loop
			}

			frame.stackPush(obj)
		case Apply:
			args := frame.stackPop()
			obj := frame.stackPop()

			seqable, ok := args.(Seqable)
			if !ok {
				newIp, err := e.unwind(env.NewArgTypeError(1, args, "Seqable"))
				if err != nil {
					return nil, err
				}

				frame.Ip = newIp
				continue loop
			}

			sq := seqable.Seq()

			callArgs, err := ToSlice(env, sq)
			if err != nil {
				newIp, err := e.unwind(err)
				if err != nil {
					return nil, err
				}

				frame.Ip = newIp
				continue loop
			}

			switch callable := obj.(type) {
			case Callable:
				if fn, ok := callable.(*Fn); ok {
					_, ferr := e.pushFrame(fn, callArgs)
					if ferr != nil {
						return nil, ferr
					}

					obj, err = e.RunBC(env, fn)
				} else {

					obj, err = callable.Call(env, callArgs)
				}

				if err != nil {
					newIp, err := e.unwind(err)
					if err != nil {
						return nil, err
					}

					frame.Ip = newIp
					continue loop
				}

				frame.stackPush(obj)
			default:
				e.printBacktrace(1000)
				s, err := ToString(env, callable)
				if err != nil {
					return nil, err
				}

				return nil, env.NewError(s + " is not a Fn")
			}
		case MethodCall:
			ms := c.methods[a]

			args := frame.stackPopN(int(ms.Arity))
			obj := frame.stackPop()

			res, err := e.methodCall(env, ms.Method, obj, args)
			if err != nil {
				newIp, err := e.unwind(err)
				if err != nil {
					return nil, err
				}

				frame.Ip = newIp
				continue loop
			}

			frame.stackPush(res)
		case Throw:
			obj := frame.stackPop()

			var err error
			if ee, ok := obj.(error); ok {
				err = ee
			} else {
				ee := NewEvalError(env, "thrown non-error value")
				ee.AddData(env, obj)
				err = ee
			}

			newIp, err := e.unwind(err)
			if err != nil {
				return nil, err
			}

			frame.Ip = newIp
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

					frame.Ip = newIp
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
			return nil, ErrorArity(env, int(frame.Arity))
		case CheckType:
			obj := frame.stackPop()

			if typ, ok := c.literals[a].(Type); ok {
				frame.stackPush(MakeBoolean(IsInstance(env, typ, obj)))
			} else {
				panic("nope")
			}

		default:
			return nil, fmt.Errorf("unimplemented instruction: %s", OpCode(op))
		}

		frame.Ip++
	}
}
