package core

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

type VarData struct {
	vr *Var
}

func (v *VarData) insData() {}

func (v *VarData) String() string {
	return "resolve: " + v.vr.Name()
}

type BindingData struct {
	up int
	in int
}

func (v *BindingData) insData() {}

func (v *BindingData) String() string {
	return fmt.Sprintf("%d:%d", v.up, v.in)
}

type InstructionData interface {
	insData()
}

type LiteralData struct {
	obj Object
}

func (v *LiteralData) insData() {}

func (v *LiteralData) String() string {
	str, _ := v.obj.ToString(nil, false)
	return fmt.Sprintf("literal: %q", str)
}

type DefData struct {
	Meta *ArrayMap
	vr   *Var
}

func (v *DefData) insData() {}

func (v *DefData) String() string {
	return fmt.Sprintf("<def: %s>", v.vr.Name())
}

type MethodCallData struct {
	Method string
}

func (v *MethodCallData) insData() {}

type FnData struct {
	Code *Code
}

func (v *FnData) insData() {}

func (v *FnData) String() string {
	return fmt.Sprintf("<fn %s upvals:%d>", v.Code.Position(), v.Code.importUpvals)
}

type CheckTypeData struct {
	Type *Type
}

func (v *CheckTypeData) insData() {}

//go:generate stringer -type OpCode
type OpCode byte

type Instruction struct {
	Op   OpCode
	A0   int32
	Data InstructionData
}

const (
	Noop OpCode = iota
	Pop
	Dup
	Jump
	JumpIfTrue
	JumpIfFalse
	ResolveVar
	SetMacro
	GetBinding
	GetUpval
	SetUpval
	PushLiteral
	MakeVector
	MakeSmallMap
	MakeLargeMap
	MakeSet
	DefValue
	DefValue3
	Def
	Def3
	SetMeta
	Call
	Apply
	MethodCall
	Throw
	PushHandler
	PopHandler
	CheckType
	SetLocal
	GetLocal
	Return
	CheckArityFixed
	CheckArityMin
	ThrowArity
	MakeFn
	RefUpval
	PushSelfFn
	SetLabel // only used during bytecode generation
	SetLine  // only used during bytecode generation
)

type handler struct {
	ip, sp int
}

type EngineFrame struct {
	Ip       int
	Code     *Fn
	Bindings []Object
	Arity    int32
	Args     []Object

	Handlers []handler
}

type Engine struct {
	frames []EngineFrame
	stack  []Object
}

type Compiler struct {
	insns []*Instruction

	fnExpr *FnExpr
	fn     *fnFrame

	filename string
	lines    []int

	nextLabel int32

	curLine int
}

func (e *Compiler) newLabel() int32 {
	e.nextLabel++
	return e.nextLabel
}

func (e *Compiler) updatePos(p Position) {
	if p.startLine == 0 {
		return
	}

	if e.filename == "" {
		e.filename = p.Filename()
	}

	if e.curLine == p.startLine {
		return
	}

	e.curLine = p.startLine

	e.insn(Instruction{
		Op: SetLine,
		A0: int32(p.startLine),
	})

}

func (e *Compiler) setLine(line, ip int) {
	if len(e.lines) == 0 {
		e.lines = append(e.lines, 0, line)
	} else if e.lines[len(e.lines)-1] != line {
		e.lines = append(e.lines, ip, line)
	}
}

func printInsns(stream []*Instruction) {
	for ip, insn := range stream {
		if insn.Data != nil {
			fmt.Printf("% 4d | %s %d %s\n", ip, insn.Op, insn.A0, insn.Data)
		} else {
			fmt.Printf("% 4d | %s %d\n", ip, insn.Op, insn.A0)
		}
	}
}

const debugCompile = false

func (c *Compiler) Export() *Code {
	// We now go through all the upvals (and their instructions) and
	// set the upval indexes. We set this up such that all imported
	// upvals are first, then all self upvals are second.

	if debugCompile {
		fmt.Printf("====== export %s\n", c.fnExpr.Pos())
		printInsns(c.insns)
	}

	var insns []*Instruction

	// TODO here is where peephole optimizations on c.insns should be done.

	// Now go through and fixup all the labels!

	label2ip := map[int]int{}

	var ip int
	for _, i := range c.insns {
		switch i.Op {
		case SetLabel:
			label2ip[int(i.A0)] = ip
		case SetLine:
			c.setLine(int(i.A0), ip)
		default:
			insns = append(insns, i)
			ip++
		}
	}

	insns = append(insns, &Instruction{Op: Return})

	if debugCompile {
		fmt.Println("--------")
		printInsns(insns)
	}

	var be BytecodeEncoder

	for idx := range insns {
		i := insns[idx]

		switch i.Op {
		case Jump, JumpIfFalse, JumpIfTrue, PushHandler:
			i.A0 = int32(label2ip[int(i.A0)])
		}

		err := be.Encode(i)
		if err != nil {
			panic(fmt.Sprintf("unable to encode instruction at %d: %s", idx, err))
		}
	}

	c.lines = append(c.lines, len(c.insns))

	if debugCompile {
		fmt.Println("====== export done")
	}

	var importNames []Symbol

	for _, vb := range c.fn.importUpvals {
		importNames = append(importNames, vb.name)
	}

	return &Code{
		numBindings:    c.fn.totalBindings,
		lines:          c.lines,
		filename:       c.filename,
		data:           be.CodeData,
		importBindings: importNames,
	}
}

func printBindings(fn *fnFrame) {
	i := 0
	for ; fn != nil; fn = fn.parent {
		i++
		fmt.Printf("frame %d:\n", i)

		for vf := fn.top; vf != nil; vf = vf.parent {
			fmt.Printf("  %v\n", vf.bindings)
		}
	}
}

func (c *Compiler) Process(env *Env, expr Expr) error {
	c.updatePos(expr.Pos())

	switch e := expr.(type) {
	case *VarRefExpr:
		c.insns = append(c.insns, &Instruction{
			Op:   ResolveVar,
			Data: &VarData{vr: e.vr},
		})
	case *SetMacroExpr:
		c.insns = append(c.insns, &Instruction{
			Op:   SetMacro,
			Data: &VarData{vr: e.vr},
		})
	case *BindingExpr:
		vb := c.fn.lookup(e.name)
		if vb == nil {
			pos := e.Pos()
			fmt.Printf("Unable to find binding for %s at %s:%d\n", e.name.Name(), pos.Filename(), pos.startLine)
			printBindings(c.fn)
			return fmt.Errorf("Unable to find binding for %s at %s:%d", e.name.Name(), pos.Filename(), pos.startLine)
		}

		c.append(vb.read())
	case *LiteralExpr:
		c.insns = append(c.insns, &Instruction{
			Op:   PushLiteral,
			Data: &LiteralData{obj: e.obj},
		})
	case *VectorExpr:
		for _, e := range e.v {
			err := c.Process(env, e)
			if err != nil {
				return err
			}
		}

		c.insn(Instruction{
			Op: MakeVector,
			A0: int32(len(e.v)),
		})
	case *MapExpr:
		for i, k := range e.keys {
			err := c.Process(env, k)
			if err != nil {
				return err
			}

			err = c.Process(env, e.values[i])
			if err != nil {
				return err
			}
		}
		if int64(len(e.keys)) > HASHMAP_THRESHOLD/2 {
			c.insn(Instruction{
				Op: MakeLargeMap,
				A0: int32(len(e.keys) * 2),
			})
		} else {
			c.insn(Instruction{
				Op: MakeSmallMap,
				A0: int32(len(e.keys) * 2),
			})
		}
	case *SetExpr:
		for _, e := range e.elements {
			err := c.Process(env, e)
			if err != nil {
				return err
			}
		}

		c.insn(Instruction{
			Op: MakeSet,
			A0: int32(len(e.elements)),
		})

	case *DefExpr:
		var (
			op   OpCode
			data VarData
		)

		data.vr = e.vr

		meta := EmptyArrayMap()
		meta.Add(env, criticalKeywords.line, Int{I: e.startLine})
		meta.Add(env, criticalKeywords.column, Int{I: e.startColumn})
		meta.Add(env, criticalKeywords.file, String{S: *e.filename})
		meta.Add(env, criticalKeywords.ns, e.vr.ns)
		meta.Add(env, criticalKeywords.name, e.vr.name)

		//data.Meta = meta

		if e.value != nil {
			err := c.Process(env, e.value)
			if err != nil {
				return err
			}
		}

		c.insn(Instruction{
			Op:   PushLiteral,
			Data: &LiteralData{obj: meta},
		})

		if e.meta != nil {
			err := c.Process(env, e.meta)
			if err != nil {
				return err
			}
		}

		if e.value != nil {
			if e.meta != nil {
				op = DefValue3
			} else {
				op = DefValue
			}
		} else if e.meta != nil {
			op = Def3
		} else {
			op = Def
		}

		c.insns = append(c.insns, &Instruction{
			Op:   op,
			Data: &data,
		})
	case *MetaExpr:
		err := c.Process(env, e.meta)
		if err != nil {
			return err
		}
		err = c.Process(env, e.expr)
		if err != nil {
			return err
		}

		c.insns = append(c.insns, &Instruction{
			Op: SetMeta,
		})
	case *CallExpr:
		isApply := false
		// 99% of all calls are to varrefexprs
		if rv, ok := e.callable.(*VarRefExpr); ok {
			if rv.vr == env.CoreNamespace.Resolve("apply__") {
				isApply = true
			}
		}

		for _, a := range e.args {
			err := c.Process(env, a)
			if err != nil {
				return err
			}
		}

		if isApply {
			c.insn(Instruction{Op: Apply})
		} else {
			err := c.Process(env, e.callable)
			if err != nil {
				return err
			}

			c.insns = append(c.insns, &Instruction{
				Op: Call,
				A0: int32(len(e.args)),
			})
		}
	case *MethodExpr:
		for _, a := range e.args {
			err := c.Process(env, a)
			if err != nil {
				return err
			}
		}

		err := c.Process(env, e.obj)
		if err != nil {
			return err
		}

		c.insns = append(c.insns, &Instruction{
			Op: MethodCall,
			A0: int32(len(e.args)),
		})
	case *ThrowExpr:
		err := c.Process(env, e.e)
		if err != nil {
			return err
		}

		c.insns = append(c.insns, &Instruction{
			Op: Throw,
		})
	case *TryExpr:
		var finallySetupPos int32 = -1
		var cacheSetups []int32

		if e.finallyExpr != nil {
			finallySetupPos = c.newLabel()
			c.insn(Instruction{
				Op: PushHandler,
				A0: finallySetupPos,
			})
		}

		for range e.catches {
			lbl := c.newLabel()

			cacheSetups = append(cacheSetups, lbl)

			c.insn(Instruction{
				Op: PushHandler,
				A0: lbl,
			})
		}

		var bodyDonePos int32

		for i, b := range e.body {
			if i > 0 {
				c.insns = append(c.insns, &Instruction{
					Op: Pop,
				})
			}

			err := c.Process(env, b)
			if err != nil {
				return err
			}

			bodyDonePos = c.newLabel()
			c.insn(Instruction{
				Op: Jump,
				A0: int32(bodyDonePos),
			})
		}

		var nextCase int32 = -1

		var catchDonePos []int32

		for i, ce := range e.catches {
			c.insn(Instruction{
				Op: SetLabel,
				A0: int32(cacheSetups[i]),
			})

			if nextCase != -1 {
				c.insn(Instruction{
					Op: SetLabel,
					A0: nextCase,
				})
			}

			c.insn(Instruction{
				Op: Dup,
			})

			nextCase = c.newLabel()

			c.insn(Instruction{
				Op: CheckType,
				Data: &CheckTypeData{
					Type: ce.excType,
				},
			})

			c.fn.bindingFrame(1)

			vb := c.fn.set(ce.excSymbol)

			c.append(vb.set())

			for i, b := range ce.body {
				if i > 0 {
					c.insns = append(c.insns, &Instruction{
						Op: Pop,
					})
				}

				err := c.Process(env, b)
				if err != nil {
					return err
				}

				lbl := c.newLabel()
				catchDonePos = append(catchDonePos, lbl)

				c.insn(Instruction{
					Op: Jump,
					A0: lbl,
				})
			}

			c.fn.popBindingFrame()
		}

		if nextCase != -1 {
			c.insn(Instruction{
				Op: SetLabel,
				A0: nextCase,
			})

			c.insn(Instruction{
				Op: Throw,
			})
		}

		c.patchToHere(bodyDonePos)
		for _, p := range catchDonePos {
			c.patchToHere(p)
		}

		if e.finallyExpr != nil {
			c.insn(Instruction{
				Op: PopHandler,
			})

			c.patchToHere(finallySetupPos)

			for _, b := range e.finallyExpr {
				err := c.Process(env, b)
				if err != nil {
					return err
				}
				c.insns = append(c.insns, &Instruction{
					Op: Pop,
				})
			}
		}
	case *DoExpr:
		for i, b := range e.body {
			if i > 0 {
				c.insns = append(c.insns, &Instruction{
					Op: Pop,
				})
			}

			err := c.Process(env, b)
			if err != nil {
				return err
			}
		}
	case *IfExpr:
		err := c.Process(env, e.cond)
		if err != nil {
			return err
		}

		ifPos := c.newLabel()
		c.insn(Instruction{
			Op: JumpIfFalse,
			A0: ifPos,
		})

		err = c.Process(env, e.positive)
		if err != nil {
			return err
		}

		donePos := c.newLabel()
		c.insn(Instruction{
			Op: Jump,
			A0: donePos,
		})

		c.patchToHere(ifPos)

		err = c.Process(env, e.negative)
		if err != nil {
			return err
		}

		c.patchToHere(donePos)
	case *FnExpr:
		fn := &Fn{fnExpr: e}

		closure, err := compileFn(env, fn, c)
		if err != nil {
			return err
		}

		for _, vb := range closure.importedVars {
			c.append(vb.refUpval())
		}

		c.insn(Instruction{
			Op:   MakeFn,
			Data: &FnData{Code: fn.code},
		})
	case *LetExpr:
		c.fn.bindingFrame(len(e.names))
		defer c.fn.popBindingFrame()

		for i, be := range e.values {
			err := c.Process(env, be)
			if err != nil {
				return err
			}

			vb := c.fn.set(e.names[i])

			c.append(vb.set())
		}

		for i, b := range e.body {
			if i > 0 {
				c.insns = append(c.insns, &Instruction{
					Op: Pop,
				})
			}

			err := c.Process(env, b)
			if err != nil {
				return err
			}
		}
	case *LoopExpr:
		vf := c.fn.bindingFrame(len(e.names))
		defer c.fn.popBindingFrame()

		// TODO this doesn't match up with the indexes
		// that GetBinding uses!
		for i, be := range e.values {
			err := c.Process(env, be)
			if err != nil {
				return err
			}

			vb := c.fn.set(e.names[i])

			c.append(vb.set())
		}

		vf.args = e.names
		vf.recurDest = int(c.newLabel())

		c.insn(Instruction{
			Op: SetLabel,
			A0: int32(vf.recurDest),
		})

		for i, b := range e.body {
			if i > 0 {
				c.insns = append(c.insns, &Instruction{
					Op: Pop,
				})
			}

			err := c.Process(env, b)
			if err != nil {
				return err
			}
		}
	case *RecurExpr:
		v := c.fn.top

		for v != nil {
			if v.recurDest >= 0 {
				break
			}

			v = v.parent
		}

		if v == nil {
			return fmt.Errorf("unable to find loop to recur to")
		}

		for i, a := range e.args {
			err := c.Process(env, a)
			if err != nil {
				return nil
			}

			vb := v.bindings[v.args[i].Name()]

			vb.uses = append(vb.uses, c.insn(Instruction{
				Op: SetLocal,
				A0: int32(vb.index),
			}))
		}

		c.insn(Instruction{
			Op: Jump,
			A0: int32(v.recurDest),
		})

	default:
		return fmt.Errorf("unable to handle type: %T", e)
	}

	return nil
}

func (c *Compiler) nextIP() int {
	return len(c.insns)
}

func (c *Compiler) append(in *Instruction) {
	c.insns = append(c.insns, in)
}

func (c *Compiler) insn(in Instruction) *Instruction {
	idx := len(c.insns)
	c.insns = append(c.insns, &in)

	return c.insns[idx]
}

func (c *Compiler) patchToHere(pos int32) {
	c.insns = append(c.insns, &Instruction{
		Op: SetLabel,
		A0: pos,
	})
}

type Upval struct {
	Obj Object
}

func (e *Engine) stackPush(obj Object) {
	e.stack = append(e.stack, obj)
}

func (e *Engine) stackPop() Object {
	idx := len(e.stack) - 1
	val := e.stack[idx]
	e.stack = e.stack[:idx]

	return val
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
	idx := len(e.frames)
	e.frames = append(e.frames, EngineFrame{
		Code:     fn,
		Bindings: make([]Object, fn.code.numBindings),
	})

	return &e.frames[idx]
}

func (e *Engine) popFrame() {
	e.frames = e.frames[:len(e.frames)-1]
}

func (e *Engine) frameBack(cnt int) (*EngineFrame, error) {
	idx := len(e.frames) - cnt - 1

	if idx < 0 {
		return nil, fmt.Errorf("invalid upward frame request %d (have %d)", cnt, len(e.frames))
	}

	return &e.frames[idx], nil
}

func (e *Engine) topSlackSlice(cnt int) []Object {
	idx := len(e.stack) - cnt
	if idx < 0 {
		panic(fmt.Sprintf("bad top slice request %d (total: %d)", cnt, len(e.stack)))
	}

	return e.stack[idx:]
}

func (e *Engine) stackPopN(cnt int) []Object {
	idx := len(e.stack) - cnt
	if idx < 0 {
		panic(fmt.Sprintf("bad top slice request %d (total: %d)", cnt, len(e.stack)))
	}

	objs := e.stack[idx:]
	e.stack = e.stack[:idx]

	return objs
}

func EngineCode(env *Env, c *Code) (Object, error) {
	return EngineRun(env, &Fn{code: c})
}

func EngineRun(env *Env, fn *Fn) (Object, error) {
	var e Engine
	env.Engine = &e
	e.stack = make([]Object, 0, 100)
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

	e.stack = e.stack[:eh.sp]

	e.stackPush(obj)

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

/*
func (e *Engine) call(env *Env, callable Callable, objArgs []Object) (Object, error) {
	min := math.MaxInt32
	max := -1

	if fn, ok := callable.(*Fn); ok {
		if fn.arities == nil && fn.varadic == nil {
			_, err := compileFn(env, fn, nil)
			if err != nil {
				return nil, err
			}
		}

		for _, c := range fn.arities {
			if len(objArgs) == c.arity {
				fr := e.pushFrame(fn)
				for i, a := range objArgs {
					fr.Bindings[i] = a
				}

				return e.Run(env, fn)
			}

			if c.arity < min {
				min = c.arity
			}

			if c.arity > max {
				max = c.arity
			}
		}

		// - 1 because the last argument is where the rest goes
		if fn.varadic != nil && len(objArgs) >= (fn.varadic.arity-1) {
			c := fn.varadic
			var restArgs Object = NIL
			if len(objArgs) > c.arity-1 {
				restArgs = &ArraySeq{arr: objArgs, index: c.arity}
			}

			restPos := c.arity - 1 // the last arg

			fr := e.pushFrame(fn)
			for i, a := range objArgs[:c.arity] {
				fr.Bindings[i] = a
			}

			fr.Bindings[restPos] = restArgs

			return e.Run(env, fn)
		}

		fmt.Printf("arity error calling fn at %s:%d\n", fn.fnExpr.Filename(), fn.fnExpr.startLine)

		return nil, ErrorArityMinMax(env, 0, 0, 0)
	}

	return callable.Call(env, objArgs)
}
*/

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

type fnClosure struct {
	importedVars []*varBind
}

func Compile(env *Env, exprs []Expr) (*Fn, error) {
	fn := &Fn{
		fnExpr: &FnExpr{
			arities: []FnArityExpr{
				{
					body: exprs,
				},
			},
		},
	}

	_, err := compileFn(env, fn, nil)
	return fn, err
}

func CompileScript(env *Env, exprs []Expr) (*Fn, error) {
	fn := &Fn{
		fnExpr: &FnExpr{
			arities: []FnArityExpr{
				{
					body: exprs,
				},
			},
		},
	}

	_, err := compileFn(env, fn, nil)
	return fn, err
}

func compileFn(env *Env, fn *Fn, parent *Compiler) (*fnClosure, error) {
	var c Compiler
	c.fnExpr = fn.fnExpr

	var parentFn *fnFrame
	if parent != nil {
		parentFn = parent.fn
	}

	var fnf *fnFrame
	if parentFn != nil {
		fnf = parentFn.childFrame(nil)
	} else {
		fnf = newFrame(nil)
	}

	c.fn = fnf

	var nextArity int32 = -1

	genFn := func(arity *FnArityExpr, checkOp OpCode) error {
		c.updatePos(arity.Position)

		if nextArity != -1 {
			c.insn(Instruction{
				Op: SetLabel,
				A0: nextArity,
			})
		}

		locals := len(arity.args)
		if fn.fnExpr.self.name != nil {
			locals++
		}

		af := fnf.bindingFrame(locals)
		defer fnf.popBindingFrame()

		af.args = arity.args

		for _, a := range arity.args {
			fnf.set(a)
		}

		if fn.fnExpr.self.name != nil {
			selfVar := fnf.set(fn.fnExpr.self)

			c.insn(Instruction{
				Op: PushSelfFn,
			})

			c.append(selfVar.set())
		}

		if checkOp == CheckArityMin {
			c.insn(Instruction{
				Op: checkOp,
				A0: int32(len(arity.args) - 1), // because the Nth local is where rest goes
			})
		} else {
			c.insn(Instruction{
				Op: checkOp,
				A0: int32(len(arity.args)),
			})
		}

		nextArity = c.newLabel()

		c.insn(Instruction{
			Op: JumpIfFalse,
			A0: nextArity,
		})

		bodyStart := c.nextIP()

		recurDest := c.newLabel()
		af.recurDest = int(recurDest)

		c.insn(Instruction{
			Op: SetLabel,
			A0: recurDest,
		})

		err := c.Process(env, &DoExpr{body: arity.body})
		if err != nil {
			return err
		}

		var prelude []*Instruction

		for _, a := range arity.args {
			vfb := fnf.lookup(a)
			if vfb.upval {
				//fmt.Printf("* %s arity %d provides an upval: %s (%d)\n", arity.Pos(), len(af.args), vfb.name.Name(), vfb.index)

				prelude = append(prelude,
					&Instruction{
						Op: GetLocal,
						A0: int32(vfb.index),
					},
					vfb.u(&Instruction{
						Op: SetUpval,
						A0: int32(vfb.upidx),
					}),
				)
			}
		}

		c.insns = slices.Insert(c.insns, bodyStart, prelude...)

		c.insn(Instruction{Op: Return})

		return nil
	}

	for _, arity := range fn.fnExpr.arities {
		err := genFn(&arity, CheckArityFixed)
		if err != nil {
			return nil, err
		}
	}

	if v := fn.fnExpr.variadic; v != nil {
		err := genFn(v, CheckArityMin)
		if err != nil {
			return nil, err
		}
	}

	c.insn(Instruction{
		Op: SetLabel,
		A0: nextArity,
	})

	c.insn(Instruction{
		Op: ThrowArity,
	})

	totalUpvals := fnf.assignUpvals()

	parentVBs := fnf.closeFrame()

	fn.code = c.Export()
	fn.code.totalUpvals = totalUpvals
	fn.upvals = make([]*NamedPair, fn.code.totalUpvals)
	fn.code.importUpvals = len(parentVBs)

	cl := &fnClosure{
		importedVars: parentVBs,
	}
	return cl, nil
}
