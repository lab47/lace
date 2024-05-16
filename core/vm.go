package core

import (
	"fmt"

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
	switch sv := v.obj.(type) {
	case String:
		return fmt.Sprintf("literal: %q", sv.S())
	case Int:
		return fmt.Sprintf("literal: %q", sv.I())
	case Symbol:
		return fmt.Sprintf("literal: %s", sv.String())
	case Keyword:
		return fmt.Sprintf("literal: %s", sv.String())
	default:
		return fmt.Sprintf("literal: %T{}", v.obj)
	}
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
	S    string
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
	PushNil
	PushInt
	SetLabel // only used during bytecode generation
	SetLine  // only used during bytecode generation
	SetFile  // only used during bytecode generation
)

type handler struct {
	ip, sp int
}

type Compiler struct {
	insns []*Instruction

	fnExpr *FnExpr
	fn     *fnFrame

	filename string

	lines      []int
	macroLines []int

	files      []string
	fileFromIp []int

	nextLabel int32

	curLine int
	curFile string
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
		S:  p.Filename(),
	})

	if p.Filename() != e.curFile {
		e.curFile = p.Filename()
		e.insn(Instruction{
			Op: SetFile,
			S:  e.curFile,
		})
	}

}

func (e *Compiler) setLine(line, ip int) {
	if len(e.lines) == 0 {
		e.lines = append(e.lines, 0, line)
	} else if e.lines[len(e.lines)-1] != line {
		e.lines = append(e.lines, ip, line)
	}
}

func (e *Compiler) setMacroLine(line, ip int) {
	if len(e.macroLines) == 0 {
		e.macroLines = append(e.macroLines, 0, line)
	} else if e.macroLines[len(e.macroLines)-1] != line {
		e.macroLines = append(e.macroLines, ip, line)
	}
}

func (e *Compiler) setFile(file string, ip int) {
	idx := slices.Index(e.files, file)
	if idx == -1 {
		idx = len(e.files)
		e.files = append(e.files, file)
	}

	if len(e.fileFromIp) == 0 {
		e.fileFromIp = append(e.fileFromIp, 0, idx)
	} else if e.fileFromIp[len(e.fileFromIp)-1] != idx {
		e.fileFromIp = append(e.fileFromIp, ip, idx)
	}
}

func printInsn(ip int, insn *Instruction) {
	if insn.S != "" {
		fmt.Printf("% 4d | %s %s %d\n", ip, insn.Op, insn.S, insn.A0)
	} else if insn.Data != nil {
		fmt.Printf("% 4d | %s %d %s\n", ip, insn.Op, insn.A0, insn.Data)
	} else {
		fmt.Printf("% 4d | %s %d\n", ip, insn.Op, insn.A0)
	}

}

func printInsns(stream []*Instruction) {
	for ip, insn := range stream {
		printInsn(ip, insn)
	}
}

type stackEffect struct {
	pop, push int
}

const (
	va  = -2
	vaa = -3
)

var stackEffects = [...]*stackEffect{
	Pop:          {1, 0},
	Dup:          {0, 1},
	Jump:         {0, 0},
	JumpIfTrue:   {1, 0},
	JumpIfFalse:  {1, 0},
	ResolveVar:   {0, 1},
	SetMacro:     {0, 1},
	GetBinding:   {0, 1},
	GetLocal:     {0, 1},
	SetLocal:     {1, 0},
	GetUpval:     {0, 1},
	SetUpval:     {0, 1},
	PushLiteral:  {0, 1},
	MakeVector:   {va, 1},
	MakeSmallMap: {va, 1},
	MakeLargeMap: {va, 1},
	MakeSet:      {va, 1},
	DefValue:     {2, 1},
	DefValue3:    {3, 1},
	Def:          {1, 1},
	Def3:         {2, 1},
	SetMeta:      {2, 1},
	Call:         {vaa, 1},
	Apply:        {2, 1},
	MethodCall:   {vaa, 1},
	Throw:        {1, 0},

	// it doesn't actually push, but the stack is reset here and there needs to be
	// room on the stack for tha value, so we reserve that spot using this effect
	PushHandler:     {0, 1},
	PopHandler:      {0, 0},
	CheckType:       {1, 1},
	Return:          {1, 0},
	CheckArityFixed: {0, 1},
	CheckArityMin:   {0, 1},
	ThrowArity:      {0, 0},
	MakeFn:          {va, 1},
	RefUpval:        {0, 1},
	PushSelfFn:      {0, 1},
	PushNil:         {0, 1},
	PushInt:         {0, 1},
}

func (c *Compiler) Export(show bool) *Code {
	// We now go through all the upvals (and their instructions) and
	// set the upval indexes. We set this up such that all imported
	// upvals are first, then all self upvals are second.

	if show {
		fmt.Printf("====== export %s\n", c.fnExpr.Pos())
		//printInsns(c.insns)
	}

	var insns []*Instruction

	// TODO here is where peephole optimizations on c.insns should be done.

	// Now go through and fixup all the labels and calculate the max stack needed!

	label2ip := map[int]int{}

	var curStack, maxStack int

	se := func(effect int) {
		curStack += effect
		if curStack > maxStack {
			maxStack = curStack
		}
	}

	var ip int
	for idx, i := range c.insns {
		if show {
			printInsn(idx, i)
		}

		switch i.Op {
		case SetLabel:
			label2ip[int(i.A0)] = ip
		case SetLine:
			if i.S == c.filename {
				c.setLine(int(i.A0), ip)
			} else {
				c.setMacroLine(int(i.A0), ip)
			}
		case SetFile:
			c.setFile(i.S, ip)
		default:
			effect := stackEffects[i.Op]
			if effect == nil {
				panic("missing stack effect for " + i.Op.String())
			}

			pop := effect.pop

			switch effect.pop {
			case va:
				pop = int(i.A0)
			case vaa:
				pop = int(i.A0 + 1)
			}

			se(-pop)
			se(effect.push)

			//fmt.Printf("  SE | -%d + %d => %d (%d)\n", pop, effect.push, curStack, maxStack)

			insns = append(insns, i)
			ip++
		}
	}

	insns = append(insns, &Instruction{Op: Return})

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

	if show {
		fmt.Println("--------")
		printInsns(insns)
	}

	c.lines = append(c.lines, len(insns))
	if len(c.macroLines) > 0 {
		c.macroLines = append(c.macroLines, len(insns))
	}

	if len(c.fileFromIp) > 0 {
		c.fileFromIp = append(c.fileFromIp, len(insns))
	}

	if show {
		fmt.Println("====== export done")
	}

	var importNames []Symbol

	for _, vb := range c.fn.importUpvals {
		importNames = append(importNames, vb.name)
	}

	return &Code{
		fnId:           nextFnId.Add(1),
		numBindings:    c.fn.totalBindings,
		lines:          c.lines,
		macroLines:     c.macroLines,
		files:          c.files,
		fileFromIp:     c.fileFromIp,
		filename:       c.filename,
		data:           be.CodeData,
		importBindings: importNames,
		stackSize:      uint32(maxStack),
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
			// The parsing of the BindingExpr validated that the name
			// does exist, we just can't see it here so we assume it's
			// in the parent.
			vb = c.fn.createUnknownUpval(e.name)
			/*
				pos := e.Pos()
				fmt.Printf("Unable to find binding for %s at %s:%d\n", e.name.Name(), pos.Filename(), pos.startLine)
				printBindings(c.fn)
				return fmt.Errorf("Unable to find binding for %s at %s:%d", e.name.Name(), pos.Filename(), pos.startLine)
			*/
		}

		c.append(vb.read())
	case *LiteralExpr:
		specialized := false

		switch sv := e.obj.(type) {
		case Nil:
			specialized = true
			c.insn(Instruction{
				Op: PushNil,
			})
		case Int:
			if sv.I() < 1024 {
				specialized = true
				c.insn(Instruction{
					Op: PushInt,
					A0: int32(sv.I()),
				})
			}
		}

		if !specialized {
			c.insns = append(c.insns, &Instruction{
				Op:   PushLiteral,
				Data: &LiteralData{obj: e.obj},
			})
		}
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
		meta.Add(env, criticalKeywords.line, MakeInt(e.startLine))
		meta.Add(env, criticalKeywords.column, MakeInt(e.startColumn))
		meta.Add(env, criticalKeywords.file, MakeString(e.filename))
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
		tgt := e.callable

		isApply := false
		// 99% of all calls are to varrefexprs
		if rv, ok := e.callable.(*VarRefExpr); ok {
			switch rv.vr {
			case env.CoreNamespace.Resolve("apply__"):
				isApply = true
			}
		}

		if isApply {
			for _, a := range e.args {
				err := c.Process(env, a)
				if err != nil {
					return err
				}
			}
			c.insn(Instruction{Op: Apply})
		} else {
			err := c.Process(env, tgt)
			if err != nil {
				return err
			}

			for _, a := range e.args {
				err := c.Process(env, a)
				if err != nil {
					return err
				}
			}

			c.insns = append(c.insns, &Instruction{
				Op: Call,
				A0: int32(len(e.args)),
			})
		}
	case *MethodExpr:
		err := c.Process(env, e.obj)
		if err != nil {
			return err
		}

		for _, a := range e.args {
			err := c.Process(env, a)
			if err != nil {
				return err
			}
		}

		c.insn(Instruction{
			Op: MethodCall,
			A0: int32(len(e.args)),
			Data: &MethodCallData{
				Method: e.method,
			},
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
		var finallyWhileUnwinding int32 = -1
		var cacheSetups []int32

		if e.finallyExpr != nil {
			finallyWhileUnwinding = c.newLabel()
			c.insn(Instruction{
				Op: PushHandler,
				A0: finallyWhileUnwinding,
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

			c.insn(Instruction{
				Op: JumpIfFalse,
				A0: nextCase,
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
			}

			lbl := c.newLabel()
			catchDonePos = append(catchDonePos, lbl)

			c.insn(Instruction{
				Op: Jump,
				A0: lbl,
			})

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
			// We output the finaly code twice, once for while unwinding and once for a simple return
			// We could make the rethrow conditionallized but it's cleaner to just emit the code twice.

			c.insn(Instruction{
				Op: PopHandler,
			})

			for _, b := range e.finallyExpr {
				err := c.Process(env, b)
				if err != nil {
					return err
				}
				c.insns = append(c.insns, &Instruction{
					Op: Pop,
				})
			}

			bottom := c.newLabel()

			c.insn(Instruction{
				Op: Jump,
				A0: bottom,
			})

			c.patchToHere(finallyWhileUnwinding)

			for _, b := range e.finallyExpr {
				err := c.Process(env, b)
				if err != nil {
					return err
				}
				c.insns = append(c.insns, &Instruction{
					Op: Pop,
				})
			}

			c.insn(Instruction{
				Op: Throw,
			})

			c.insn(Instruction{
				Op: SetLabel,
				A0: bottom,
			})
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

		var (
			closure *fnClosure
			code    *Code
		)

		if e.closure != nil {
			closure = e.closure
			code = e.compiled
		} else {
			panic("detected fn not compiled during parsing")
		}

		for _, name := range closure.importedVars {
			vb := c.fn.lookup(name)
			if vb == nil {
				vb = c.fn.createUnknownUpval(name)
			}
			c.append(vb.refUpval())
		}

		c.insn(Instruction{
			Op:   MakeFn,
			A0:   int32(len(closure.importedVars)),
			Data: &FnData{Code: code},
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

		for _, a := range e.args {
			err := c.Process(env, a)
			if err != nil {
				return nil
			}
		}

		for i := len(e.args) - 1; i >= 0; i-- {
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

type fnClosure struct {
	importedVars []Symbol
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
		if fn.fnExpr.self.name != "" {
			locals++
		}

		af := fnf.bindingFrame(locals)
		defer fnf.popBindingFrame()

		af.args = arity.args

		for _, a := range arity.args {
			fnf.set(a)
		}

		if fn.fnExpr.self.name != "" {
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

		body := arity.body
		if len(arity.body) == 0 {
			body = []Expr{&LiteralExpr{obj: NIL, Position: arity.Position}}
		}

		err := c.Process(env, &DoExpr{body: body})
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

	fn.code = c.Export(env.DebugBytecode)
	fn.code.totalUpvals = totalUpvals
	fn.importedUpvals = make([]*NamedPair, len(parentVBs))
	fn.code.importUpvals = len(parentVBs)

	cl := &fnClosure{
		importedVars: parentVBs,
	}
	return cl, nil
}
