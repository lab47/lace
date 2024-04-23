package core

import (
	"fmt"

	"github.com/lab47/lace/core/insn"
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
	idx := uint(len(e.vars))
	e.vars = append(e.vars, vr)

	sym := AssembleSymbol(vr.ns.Name.Name(), vr.name.String())

	e.varNames = append(e.varNames, sym)

	return idx
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

const debugBC = false

func (e *Engine) RunBC(env *Env, fn *Fn) (Object, error) {
	c := fn.code.data

	//fmt.Printf("==== enter frame %d %s:%d =====\n", len(e.frames), fn.code.filename, fn.code.lineForIp(0))
	//defer fmt.Printf("==== exit frame %d =====\n", len(e.frames))
	defer e.popFrame()

	var (
		ip  int
		tmp Object
		err error
	)

loop:
	for {
		//nsn := c.insns[ip]

		op, a := insn.Decode(c.insns[ip])

		if debugBC {
			fmt.Printf("% 2d|% 4d|% 4d| %s %d\n", len(e.frames), ip, fn.code.lineForIp(ip), OpCode(op), a)
			// e.printStack(env)
		}

		switch OpCode(op) {
		case Pop:
			e.stackPop()
		case Return:
			return e.stackPop(), nil
		case Jump:
			ip = int(a)
			continue loop
		case JumpIfTrue:
			if ToBool(e.stackPop()) {
				ip = int(a)
				continue loop
			}
		case JumpIfFalse:
			if !ToBool(e.stackPop()) {
				ip = int(a)
				continue loop
			}
		case GetUpval:
			e.stackPush(fn.upvals[a].Value)
		case RefUpval:
			e.stackPush(fn.upvals[a])
		case SetUpval:
			uv := fn.upvals[a]
			if uv == nil {
				uv = &NamedPair{}
				fn.upvals[a] = uv
			}

			uv.Value = e.stackPop()
		case ResolveVar:
			tmp = c.vars[a].Resolve()
			if debugBC {
				fmt.Printf("             | %s\n", c.vars[a].Name())
			}
			e.stackPush(tmp)
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
			e.stack = append(e.stack, vr)
		case GetLocal:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			e.stackPush(fr.Bindings[a])
		case SetLocal:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			fr.Bindings[a] = e.stackPop()
		case PushLiteral:
			e.stackPush(c.literals[a])
		case PushSelfFn:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}
			e.stackPush(fr.Code)
		case MakeVector:
			vec := NewVectorFrom(e.topSlackSlice(int(a))...)

			e.stackPush(vec)
		case MakeLargeMap:
			data := e.topSlackSlice(int(a))

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

			e.stackPush(res)
		case MakeSmallMap:
			cnt := a
			res := EmptyArrayMap()

			if cnt > 0 {
				data := e.stackPopN(int(cnt))

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

			e.stackPush(res)
		case MakeSet:
			data := e.topSlackSlice(int(a))

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

			e.stackPush(res)
		case Def:
			vr := c.defVars[a]

			basic := e.stackPop()
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

			e.stackPush(vr)
		case Def3:
			vr := c.defVars[a]

			v := e.stackPop()

			if err := Cast(env, e.stackPop(), &vr.meta); err != nil {
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

			e.stackPush(vr)
		case DefValue:
			vr := c.defVars[a]

			val := e.stackPop()

			if err := Cast(env, e.stackPop(), &vr.meta); err != nil {
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

			e.stackPush(vr)
		case DefValue3:
			vr := c.defVars[a]
			v := e.stackPop()

			if err := Cast(env, e.stackPop(), &vr.meta); err != nil {
				return nil, err
			}

			val := e.stackPop()
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

			e.stackPush(vr)
		case SetMeta:
			res := e.stackPop()
			meta := e.stackPop()

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

			e.stackPush(mo)
		case Call:
			obj := e.stackPop()

			switch callable := obj.(type) {
			case Callable:
				args := e.stackPopN(int(a))

				fr, err := e.frameBack(0)
				if err != nil {
					panic(err)
				}

				fr.Ip = ip

				obj, err := e.call(env, callable, args)
				if err != nil {
					newIp, err := e.unwind(err)
					if err != nil {
						return nil, err
					}

					ip = newIp
					continue loop
				}

				e.stackPush(obj)
			default:
				e.printBacktrace(1000)
				s, err := callable.ToString(env, false)
				if err != nil {
					return nil, err
				}

				return nil, env.RT.NewError(s + " is not a Fn")
			}
		case Apply:
			args := e.stackPop()
			obj := e.stackPop()

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

			switch callable := obj.(type) {
			case Callable:
				fr, err := e.frameBack(0)
				if err != nil {
					panic(err)
				}

				fr.Ip = ip

				obj, err := e.call(env, callable, callArgs)
				if err != nil {
					newIp, err := e.unwind(err)
					if err != nil {
						return nil, err
					}

					ip = newIp
					continue loop
				}

				e.stackPush(obj)
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

			obj := e.stackPop()

			res, err := e.methodCall(env, ms.Method, obj, e.stackPopN(int(a)))
			if err != nil {
				newIp, err := e.unwind(err)
				if err != nil {
					return nil, err
				}

				ip = newIp
				continue loop
			}

			e.stackPush(res)
		case Throw:
			newIp, err := e.unwind(&VMError{obj: e.stackPop()})
			if err != nil {
				return nil, err
			}

			ip = newIp
			continue loop

		case PushHandler:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			fr.Handlers = append(fr.Handlers, handler{
				ip: int(a),
				sp: len(e.stack),
			})

		case PopHandler:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			fr.Handlers = fr.Handlers[:len(fr.Handlers)-1]
		case MakeFn:
			code := c.codes[a]

			upvals := make([]*NamedPair, code.totalUpvals)

			for i := code.importUpvals - 1; i >= 0; i-- {
				upvals[i] = e.stackPop().(*NamedPair)
			}

			fn := &Fn{
				code:   code,
				upvals: upvals,
			}

			e.stackPush(fn)
		case CheckArityFixed:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			if fr.Arity == int32(a) {
				copy(fr.Bindings, fr.Args)
				e.stackPush(MakeBoolean(true))
			} else {
				e.stackPush(MakeBoolean(false))
			}
		case CheckArityMin:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			if fr.Arity >= int32(a) {
				copy(fr.Bindings, fr.Args[:a])

				fr.Bindings[a] = &ArraySeq{arr: fr.Args, index: int(a)}

				e.stackPush(MakeBoolean(true))
			} else {
				e.stackPush(MakeBoolean(false))
			}
		case ThrowArity:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			return nil, ErrorArity(env, int(fr.Arity))
		default:
			return nil, fmt.Errorf("unimplemented instruction: %d", op)
		}

		ip++
	}
}
