package core

import (
	"fmt"
	"reflect"
	"strings"

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
	return "literal: " + str
}

type VectorData struct {
	cnt int
}

func (v *VectorData) insData() {}

type MapData struct {
	cnt int
}

func (v *MapData) insData() {}

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
	Code   *Code
	Upvals int
}

func (v *FnData) insData() {}

func (v *FnData) String() string {
	return fmt.Sprintf("<fn %s upvals:%d>", v.Code.Position(), v.Upvals)
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

type varBind struct {
	index int

	name Symbol
	home *fnFrame

	upval bool
	upidx int

	uses []*Instruction

	ref *varBind
}

type varFrame struct {
	parent   *varFrame
	args     []Symbol
	bindings map[string]*varBind
	top      int
	expr     Expr
}

type fnFrame struct {
	parent *fnFrame
	top    *varFrame

	importUpvals []*varBind
	upvals       []*varBind
	letUpvals    []*varBind

	selfUpvals int
	refUpvals  int
}

type Compiler struct {
	arity         int
	varadic       bool
	insns         []*Instruction
	curBindings   int
	totalBindings int

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

func (c *Compiler) Init(expr Expr, args []Symbol) {
	var vf varFrame
	vf.bindings = make(map[string]*varBind)
	vf.args = args
	vf.expr = expr
	vf.top = int(c.newLabel())

	c.insn(Instruction{
		Op: SetLabel,
		A0: int32(vf.top),
	})

	var fn fnFrame
	fn.parent = c.fn
	fn.top = &vf

	c.fn = &fn

	if len(args) > 0 {
		c.arity = len(args)
		start := c.addBindings(len(args))

		for i, a := range args {
			vf.bindings[a.Name()] = &varBind{index: start + i, name: a, home: c.fn}
		}
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

func (c *Compiler) Export() *Code {
	// We now go through all the upvals (and their instructions) and
	// set the upval indexes. We set this up such that all imported
	// upvals are first, then all self upvals are second.

	fmt.Printf("====== export %s\n", c.fnExpr.Pos())
	printInsns(c.insns)

	/*
		upstream := map[*varBind]int{}

		var importNames []Symbol

		for _, b := range c.fn.importUpvals {
			idx, ok := upstream[b.ref]
			if !ok {
				idx = len(upstream)
				upstream[b.ref] = idx
				importNames = append(importNames, b.name)
			}

			for _, i := range b.uses {
				*i = Instruction{
					Op: GetUpval,
					A0: int32(idx),
				}
			}
		}

		for i, b := range c.fn.upvals {
			idx := len(upstream) + i
			upstream[b] = idx
			fmt.Printf("! upval remapping local:%d => upval:%d\n", b.index, idx)
			for _, j := range b.uses {
				*j = Instruction{
					Op: GetUpval,
					A0: int32(idx),
				}
			}
		}

		// if any arguments are now upvals, we need to insert instructions
		// to set them from the locals now.

		var insns []Instruction

		for _, vb := range c.fn.top.bindings {
			if idx, ok := upstream[vb]; ok {
				insns = append(insns,
					Instruction{
						Op: GetLocal,
						A0: int32(vb.index),
					},
					Instruction{
						Op: SetUpval,
						A0: int32(idx),
					})
			}
		}
	*/

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

	for idx := range insns {
		i := insns[idx]

		switch i.Op {
		case Jump, JumpIfFalse, JumpIfTrue, PushHandler:
			i.A0 = int32(label2ip[int(i.A0)])
		}
	}

	insns = append(insns, &Instruction{Op: Return})

	c.lines = append(c.lines, len(c.insns))

	fmt.Println("--------")
	printInsns(insns)
	fmt.Println("====== export done")

	return &Code{
		arity:       c.arity,
		insns:       insns,
		numBindings: c.totalBindings,
		lines:       c.lines,
		filename:    c.filename,

		//importBindings: importNames,
	}
}

var fnExprCompiles = map[*FnExpr]*Fn{}

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

proc:
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
		fn := c.fn

		// see if it's a binding in the current frame
		for vf := fn.top; vf != nil; vf = vf.parent {
			if b, ok := vf.bindings[e.name.Name()]; ok {
				if b.upval {
					b.uses = append(b.uses, c.insn(Instruction{
						Op: GetUpval,
					}))
				} else {
					b.uses = append(b.uses, c.insn(Instruction{
						Op: GetLocal,
						A0: int32(b.index),
					}))
				}
				break proc
			}
		}

		// Ok, it's not local, let's see if we can find it in a parent
		fn = c.fn.parent

		for fn != nil {
			for vf := fn.top; vf != nil; vf = vf.parent {
				if b, ok := vf.bindings[e.name.Name()]; ok {
					if !b.upval {
						b.upval = true
						fn.upvals = append(fn.upvals, b)
					}

					vb := &varBind{
						ref:   b,
						upval: true,
						name:  e.name,
						home:  c.fn,
					}

					c.fn.importUpvals = append(c.fn.importUpvals, vb)

					vb.uses = append(vb.uses, c.insn(Instruction{
						Op: GetUpval,
					}))

					// Add a binding between the current fn and the origin fn.
					// This is so we know to pass the upval along the function chain
					// even if the in-between ones don't use it.

					for parent := c.fn.parent; parent != fn; parent = parent.parent {
						vb := &varBind{
							ref:   b,
							upval: true,
							name:  e.name,
							home:  parent,
						}

						parent.importUpvals = append(parent.importUpvals, vb)
					}

					/*
						// ok! found it in a parent, let's tag the parents binding as being an upval
						if !b.upval {
							b.upval = true
							b.upidx = len(fn.upvals)
							fn.upvals = append(fn.upvals, b)

							// insert a binding the current frame that reference the binding
							// in the parent.
							vb := &varBind{
								ref:   b,
								upval: true,
								upidx: len(c.fn.upvals),
							}
							c.fn.upvals = append(c.fn.upvals, vb)

							fn.top.bindings[e.name.Name()] = vb

							//fn.selfUpvals++

							//c.fn.refUpvals++

							// Backpatch the previous users to now get this variable from
							// an upval
							for _, i := range b.uses {
								switch c.insns[i].Op {
								case GetBinding:
									c.insns[i] = Instruction{
										Op:    GetUpval,
										Count: int32(b.upidx),
									}
								case SetLocal:
									c.insns[i] = Instruction{
										Op:    SetUpval,
										Count: int32(b.upidx),
									}
								default:
									return fmt.Errorf("Unknown op in varbind uses: %s", c.insns[i])
								}
							}
						}

						c.insn(Instruction{
							Op:    GetUpval,
							Count: int32(b.upidx),
						})
					*/

					break proc
				}
			}

			fn = fn.parent
		}

		pos := e.Pos()
		fmt.Printf("Unable to find binding for %s at %s:%d\n", e.name.Name(), pos.Filename(), pos.startLine)
		printBindings(c.fn)
		return fmt.Errorf("Unable to find binding for %s at %s:%d", e.name.Name(), pos.Filename(), pos.startLine)
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

		c.insns = append(c.insns, &Instruction{
			Op:   MakeVector,
			Data: &VectorData{cnt: len(e.v)},
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
			c.insns = append(c.insns, &Instruction{
				Op:   MakeLargeMap,
				Data: &MapData{cnt: len(e.keys) * 2},
			})
		} else {
			c.insns = append(c.insns, &Instruction{
				Op:   MakeSmallMap,
				Data: &MapData{cnt: len(e.keys) * 2},
			})
		}
	case *SetExpr:
		for _, e := range e.elements {
			err := c.Process(env, e)
			if err != nil {
				return err
			}
		}

		c.insns = append(c.insns, &Instruction{
			Op:   MakeVector,
			Data: &VectorData{cnt: len(e.elements)},
		})

	case *DefExpr:
		var (
			op   OpCode
			data DefData
		)

		data.vr = e.vr

		meta := EmptyArrayMap()
		meta.Add(env, criticalKeywords.line, Int{I: e.startLine})
		meta.Add(env, criticalKeywords.column, Int{I: e.startColumn})
		meta.Add(env, criticalKeywords.file, String{S: *e.filename})
		meta.Add(env, criticalKeywords.ns, e.vr.ns)
		meta.Add(env, criticalKeywords.name, e.vr.name)

		data.Meta = meta

		if e.value != nil {
			err := c.Process(env, e.value)
			if err != nil {
				return err
			}
		}

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

			var vf varFrame
			vf.top = -1
			vf.args = []Symbol{ce.excSymbol}
			vf.bindings = map[string]*varBind{}
			vf.parent = c.fn.top
			vf.expr = e
			c.fn.top = &vf

			start := c.addBindings(1)

			vb := &varBind{index: start, name: ce.excSymbol, home: c.fn}
			vf.bindings[ce.excSymbol.Name()] = vb

			vb.uses = append(vb.uses, c.insn(Instruction{
				Op: SetLocal,
				A0: int32(start),
			}))

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

			c.fn.top = c.fn.top.parent
			c.removeBindings(1)
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

		vbs := make([]*varBind, len(closure.importedVars))

		for vb, pos := range closure.importedVars {
			vbs[pos] = vb
		}

		for _, vb := range vbs {
			vb.uses = append(vb.uses, c.insn(Instruction{
				Op: RefUpval,
			}))
		}

		c.insn(Instruction{
			Op:   MakeFn,
			Data: &FnData{Code: fn.code, Upvals: len(vbs)},
		})
	case *LetExpr:
		start := c.addBindings(len(e.names))

		var vf varFrame
		vf.top = -1
		vf.args = e.names
		vf.bindings = map[string]*varBind{}
		vf.parent = c.fn.top
		vf.expr = e
		c.fn.top = &vf

		for i, be := range e.values {
			err := c.Process(env, be)
			if err != nil {
				return err
			}

			vb := &varBind{index: start + i, name: e.names[i], home: c.fn}

			vf.bindings[e.names[i].Name()] = vb

			vb.uses = append(vb.uses, c.insn(Instruction{
				Op: SetLocal,
				A0: int32(start + i),
			}))
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

		for _, vb := range vf.bindings {
			if vb.upval {
				c.fn.letUpvals = append(c.fn.letUpvals, vb)
			}
		}

		// See if we need to promote any of these to upvals

		// TODO could clear the locals here, for gc and safety reasons

		c.removeBindings(len(e.names))
		c.fn.top = c.fn.top.parent

	case *LoopExpr:
		start := c.addBindings(len(e.names))

		var vf varFrame
		vf.bindings = map[string]*varBind{}
		vf.args = e.names
		vf.parent = c.fn.top
		vf.expr = e
		c.fn.top = &vf

		// TODO this doesn't match up with the indexes
		// that GetBinding uses!
		for i, be := range e.values {
			err := c.Process(env, be)
			if err != nil {
				return err
			}

			vb := &varBind{index: start + i, name: e.names[i], home: c.fn}

			vf.bindings[e.names[i].Name()] = vb

			vb.uses = append(vb.uses, c.insn(Instruction{
				Op: SetLocal,
				A0: int32(start + i),
			}))
		}

		vf.top = int(c.newLabel())

		c.insn(Instruction{
			Op: SetLabel,
			A0: int32(vf.top),
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

		// TODO could clear the locals here, for gc and safety reasons

		c.removeBindings(len(e.names))
		c.fn.top = c.fn.top.parent

	case *RecurExpr:
		v := c.fn.top

		for v != nil {
			if v.top >= 0 {
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
			A0: int32(v.top),
		})

	default:
		return fmt.Errorf("unable to handle type: %T", e)
	}

	return nil
}

func (c *Compiler) addBindings(cnt int) int {
	cur := c.curBindings

	c.curBindings += cnt
	if c.curBindings > c.totalBindings {
		c.totalBindings = c.curBindings
	}

	return cur
}

func (c *Compiler) removeBindings(cnt int) {
	c.curBindings -= cnt
}

func (c *Compiler) currentIP() int {
	return len(c.insns) - 1
}

func (c *Compiler) nextIP() int {
	return len(c.insns)
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

type Code struct {
	arity       int
	varadic     bool
	insns       []*Instruction
	numBindings int

	totalUpvals int

	filename string
	lines    []int

	importBindings []Symbol
}

func (c *Code) Position() string {
	return fmt.Sprintf("%s:%d", c.filename, c.lineForIp(0))
}

func (c *Code) lineForIp(ip int) int {
	for i := 0; i < len(c.lines); i += 2 {
		if c.lines[i] <= ip && c.lines[i+2] > ip {
			return c.lines[i+1]
		}
	}

	return -1
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
	return e.Run(env, fn)
}

func (e *Engine) RunWithArgs(env *Env, fn *Fn, args []Object) (Object, error) {
	fr := e.pushFrame(fn)
	fr.Args = slices.Clone(args)
	fr.Arity = int32(len(args))
	return e.Run(env, fn)
}

func (e *Engine) Run(env *Env, fn *Fn) (Object, error) {
	c := fn.code

	//fmt.Printf("==== enter frame %d %s:%d =====\n", len(e.frames), c.filename, c.lineForIp(0))
	//defer fmt.Printf("==== exit frame %d =====\n", len(e.frames))
	defer e.popFrame()

	var (
		ip  int
		tmp Object
		err error
	)

loop:
	for {
		insn := c.insns[ip]

		//fmt.Printf("% 2d |% 4d |% 3d| %s %d %s\n", len(e.frames), ip, c.lineForIp(ip), insn.Op, insn.A0, insn.Data)
		//e.printStack(env)

		switch insn.Op {
		case Pop:
			e.stackPop()
		case Return:
			return e.stackPop(), nil
		case Jump:
			ip = int(insn.A0)
			continue loop
		case JumpIfTrue:
			if ToBool(e.stackPop()) {
				ip = int(insn.A0)
				continue loop
			}
		case JumpIfFalse:
			if !ToBool(e.stackPop()) {
				ip = int(insn.A0)
				continue loop
			}
		case GetUpval:
			e.stackPush(fn.upvals[insn.A0].Value)
		case RefUpval:
			e.stackPush(fn.upvals[insn.A0])
		case SetUpval:
			uv := fn.upvals[insn.A0]
			if uv == nil {
				uv = &NamedPair{}
				fn.upvals[insn.A0] = uv
			}

			uv.Value = e.stackPop()
		case ResolveVar:
			tmp = insn.Data.(*VarData).vr.Resolve()
			e.stackPush(tmp)
		case SetMacro:
			vr := insn.Data.(*VarData).vr
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
		case GetBinding:
			// THIS IS WRONG. ITS NOT THE CALLER, IT'S THE CLOSURE
			b := insn.Data.(*BindingData)
			fr, err := e.frameBack(b.up)
			if err != nil {
				return nil, err
			}

			e.stackPush(fr.Bindings[b.in])
		case GetLocal:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			e.stackPush(fr.Bindings[insn.A0])
		case SetLocal:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			fr.Bindings[insn.A0] = e.stackPop()
		case PushLiteral:
			d := insn.Data.(*LiteralData)
			e.stackPush(d.obj)
		case MakeVector:
			d := insn.Data.(*VectorData)

			vec := NewVectorFrom(e.topSlackSlice(d.cnt)...)

			e.stackPush(vec)
		case MakeLargeMap:
			d := insn.Data.(*MapData)

			data := e.topSlackSlice(d.cnt)

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
			d := insn.Data.(*MapData)

			res := EmptyArrayMap()

			if d.cnt > 0 {
				data := e.stackPopN(d.cnt)

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
			d := insn.Data.(*VectorData)

			data := e.topSlackSlice(d.cnt)

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
			d := insn.Data.(*DefData)

			// isMacro can be set by set-macro__ during parse stage
			if d.vr.isMacro {
				v, err := d.vr.meta.Assoc(env, criticalKeywords.macro, Boolean{B: true})
				if err != nil {
					return nil, err
				}
				var m Map
				if err := Cast(env, v, &m); err != nil {
					return nil, err
				}
				d.vr.meta = m
			}

			e.stackPush(d.vr)
		case Def3:
			d := insn.Data.(*DefData)

			d.vr.meta = d.Meta

			v := e.stackPop()

			var m Map
			if err := Cast(env, v, &m); err != nil {
				return nil, err
			}
			d.vr.meta, err = d.vr.meta.Merge(env, m)
			if err != nil {
				return nil, err
			}

			// isMacro can be set by set-macro__ during parse stage
			if d.vr.isMacro {
				v, err := d.vr.meta.Assoc(env, criticalKeywords.macro, Boolean{B: true})
				if err != nil {
					return nil, err
				}
				var m Map
				if err := Cast(env, v, &m); err != nil {
					return nil, err
				}
				d.vr.meta = m
			}

			e.stackPush(d.vr)
		case DefValue:
			d := insn.Data.(*DefData)
			val := e.stackPop()

			d.vr.Value = val

			// isMacro can be set by set-macro__ during parse stage
			if d.vr.isMacro {
				v, err := d.vr.meta.Assoc(env, criticalKeywords.macro, Boolean{B: true})
				if err != nil {
					return nil, err
				}
				var m Map
				if err := Cast(env, v, &m); err != nil {
					return nil, err
				}
				d.vr.meta = m
			}

			e.stackPush(d.vr)
		case DefValue3:
			d := insn.Data.(*DefData)
			v := e.stackPop()
			val := e.stackPop()
			d.vr.Value = val

			d.vr.meta = d.Meta

			var m Map
			if err := Cast(env, v, &m); err != nil {
				return nil, err
			}

			d.vr.meta, err = d.vr.meta.Merge(env, m)
			if err != nil {
				return nil, err
			}

			// isMacro can be set by set-macro__ during parse stage
			if d.vr.isMacro {
				v, err := d.vr.meta.Assoc(env, criticalKeywords.macro, Boolean{B: true})
				if err != nil {
					return nil, err
				}
				var m Map
				if err := Cast(env, v, &m); err != nil {
					return nil, err
				}
				d.vr.meta = m
			}

			e.stackPush(d.vr)
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
				args := e.stackPopN(int(insn.A0))

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
				s, err := callable.ToString(env, false)
				if err != nil {
					return nil, err
				}

				return nil, env.RT.NewError(s + " is not a Fn")
			}
		case MethodCall:
			d := insn.Data.(*MethodCallData)

			obj := e.stackPop()

			res, err := e.methodCall(env, d.Method, obj, e.stackPopN(int(insn.A0)))
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
				ip: int(insn.A0),
				sp: len(e.stack),
			})

		case PopHandler:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			fr.Handlers = fr.Handlers[:len(fr.Handlers)-1]
		case MakeFn:
			d := insn.Data.(*FnData)

			upvals := make([]*NamedPair, d.Code.totalUpvals)

			for i := 0; i < d.Upvals; i++ {
				upvals[i] = e.stackPop().(*NamedPair)
			}

			fn := &Fn{
				code:   d.Code,
				upvals: upvals,
			}

			e.stackPush(fn)
		case CheckArityFixed:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			if fr.Arity == insn.A0 {
				for i, a := range fr.Args {
					fr.Bindings[i] = a
				}

				e.stackPush(MakeBoolean(true))
			} else {
				e.stackPush(MakeBoolean(false))
			}
		case CheckArityMin:
			fr, err := e.frameBack(0)
			if err != nil {
				return nil, err
			}

			if fr.Arity >= insn.A0 {
				for i, a := range fr.Args[:insn.A0] {
					fr.Bindings[i] = a
				}

				fr.Bindings[insn.A0] = &ArraySeq{arr: fr.Args, index: int(insn.A0)}

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
			return nil, fmt.Errorf("unimplemented instruction: %d", insn.Op)
		}

		ip++
	}
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

		return e.Run(env, fn)
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
	importedVars map[*varBind]int
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

func compileFn(env *Env, fn *Fn, parent *Compiler) (*fnClosure, error) {
	var c Compiler
	c.fnExpr = fn.fnExpr

	var parentFn *fnFrame
	if parent != nil {
		parentFn = parent.fn
	}

	var nextArity int32 = -1

	var maxArgUpvals int

	importBinds := map[*varBind]int{}

	for _, arity := range fn.fnExpr.arities {
		c.updatePos(arity.Position)

		if nextArity != -1 {
			c.insn(Instruction{
				Op: SetLabel,
				A0: nextArity,
			})
		}

		var fnf fnFrame
		fnf.parent = parentFn
		c.fn = &fnf

		fnTop := c.newLabel()

		var af varFrame
		af.bindings = make(map[string]*varBind)
		af.args = arity.args
		af.top = int(fnTop)
		af.expr = &arity

		af.parent = c.fn.top
		c.fn.top = &af

		binds := len(af.args)
		if fn.fnExpr.self.name != nil {
			binds++
		}

		start := c.addBindings(binds)

		for i, a := range af.args {
			af.bindings[a.Name()] = &varBind{index: start + i, name: a, home: c.fn}
		}

		if fn.fnExpr.self.name != nil {
			selfIndex := start + len(af.args)

			c.insn(Instruction{
				Op:   PushLiteral,
				Data: &LiteralData{obj: fn},
			})

			vb := &varBind{index: selfIndex, name: fn.fnExpr.self, home: c.fn}
			af.bindings[fn.fnExpr.self.Name()] = vb

			vb.uses = append(vb.uses, c.insn(Instruction{
				Op: SetLocal,
				A0: int32(selfIndex),
			}))
		}

		c.insn(Instruction{
			Op: CheckArityFixed,
			A0: int32(len(af.args)),
		})

		nextArity = c.newLabel()

		c.insn(Instruction{
			Op: JumpIfFalse,
			A0: nextArity,
		})

		bodyStart := c.nextIP()

		c.insn(Instruction{
			Op: SetLabel,
			A0: fnTop,
		})

		err := c.Process(env, &DoExpr{body: arity.body})
		if err != nil {
			return nil, err
		}

		for _, vb := range fnf.importUpvals {
			idx, ok := importBinds[vb.ref]
			if !ok {
				idx = len(importBinds)
				importBinds[vb.ref] = idx
			}

			//fmt.Printf("* %s arity %d imports an upval: %s (%d)\n", arity.Pos(), len(af.args), vb.name.Name(), idx)
		}

		cnt := len(fnf.importUpvals)

		var prelude []*Instruction

		for _, a := range af.args {
			vfb := af.bindings[a.Name()]
			if vfb.upval {
				//fmt.Printf("* %s arity %d provides an upval: %s (%d)\n", arity.Pos(), len(af.args), vfb.name.Name(), vfb.index)

				prelude = append(prelude,
					&Instruction{
						Op: GetLocal,
						A0: int32(vfb.index),
					},
					&Instruction{
						Op: SetUpval,
						A0: int32(cnt),
					},
				)

				for _, i := range vfb.uses {
					if i.Op == GetLocal {
						i.Op = GetUpval
					}

					i.A0 = int32(cnt)
				}

				cnt++
			}
		}

		for _, vb := range fnf.letUpvals {
			for _, i := range vb.uses {
				switch i.Op {
				case SetLocal:
					i.Op = SetUpval
				case GetLocal:
					i.Op = GetUpval
				}

				i.A0 = int32(cnt)
			}

			cnt++
		}

		c.insns = slices.Insert(c.insns, bodyStart, prelude...)

		if cnt > maxArgUpvals {
			maxArgUpvals = cnt
		}

		c.insn(Instruction{Op: Return})

		c.fn.top = af.parent
		c.removeBindings(len(af.args))
	}

	if v := fn.fnExpr.variadic; v != nil {
		c.updatePos(v.Position)

		if nextArity != -1 {
			c.insn(Instruction{
				Op: SetLabel,
				A0: nextArity,
			})
		}

		var fnf fnFrame
		fnf.parent = parentFn
		c.fn = &fnf

		fnTop := c.newLabel()

		var af varFrame
		af.bindings = make(map[string]*varBind)
		af.args = v.args
		af.top = int(fnTop)
		af.expr = v

		af.parent = c.fn.top
		c.fn.top = &af

		binds := len(af.args)
		if fn.fnExpr.self.name != nil {
			binds++
		}

		start := c.addBindings(binds)

		for i, a := range af.args {
			af.bindings[a.Name()] = &varBind{index: start + i, name: a, home: c.fn}
		}

		/*
			if fn.fnExpr.self.name != nil {
				start := c.addBindings(1)

				c.insn(Instruction{
					Op:   PushLiteral,
					Data: &LiteralData{obj: fn},
				})

				vb := &varBind{index: start, name: fn.fnExpr.self}
				c.fn.top.bindings[fn.fnExpr.self.Name()] = vb

				vb.uses = append(vb.uses, c.insn(Instruction{
					Op: SetLocal,
					A0: int32(start),
				}))
			}
		*/
		if fn.fnExpr.self.name != nil {
			selfIndex := start + len(af.args)

			c.insn(Instruction{
				Op:   PushLiteral,
				Data: &LiteralData{obj: fn},
			})

			vb := &varBind{index: selfIndex, name: fn.fnExpr.self, home: c.fn}
			af.bindings[fn.fnExpr.self.Name()] = vb

			vb.uses = append(vb.uses, c.insn(Instruction{
				Op: SetLocal,
				A0: int32(selfIndex),
			}))
		}

		c.insn(Instruction{
			Op: CheckArityMin,
			A0: int32(len(af.args) - 1),
		})

		nextArity = c.newLabel()

		c.insn(Instruction{
			Op: JumpIfFalse,
			A0: nextArity,
		})

		bodyStart := c.nextIP()

		c.insn(Instruction{
			Op: SetLabel,
			A0: fnTop,
		})

		err := c.Process(env, &DoExpr{body: v.body})
		if err != nil {
			return nil, err
		}

		for _, vb := range fnf.importUpvals {
			idx, ok := importBinds[vb.ref]
			if !ok {
				idx = len(importBinds)
				importBinds[vb.ref] = idx
			}
			//fmt.Printf("* %s arity %d imports an upval: %s (%d)\n", v.Pos(), len(af.args), vb.name.Name(), idx)
		}

		cnt := len(fnf.importUpvals)

		var prelude []*Instruction

		for _, a := range af.args {
			vfb := af.bindings[a.Name()]
			if vfb.upval {
				//fmt.Printf("* %s arity %d fn needs upval: %s (%d)\n", v.Pos(), len(af.args), vfb.name.Name(), vfb.index)

				prelude = append(prelude,
					&Instruction{
						Op: GetLocal,
						A0: int32(vfb.index),
					},
					&Instruction{
						Op: SetUpval,
						A0: int32(cnt),
					},
				)

				for _, i := range vfb.uses {
					if i.Op == GetLocal {
						i.Op = GetUpval
					}

					i.A0 = int32(cnt)
				}

				cnt++
			}
		}

		c.insns = slices.Insert(c.insns, bodyStart, prelude...)

		if cnt > maxArgUpvals {
			maxArgUpvals = cnt
		}

		c.insn(Instruction{Op: Return})

		c.fn.top = c.fn.top.parent
		c.removeBindings(len(af.args))
	}

	c.insn(Instruction{
		Op: SetLabel,
		A0: nextArity,
	})

	c.insn(Instruction{
		Op: ThrowArity,
	})

	//for vb, pos := range importBinds {
	//fmt.Printf("*- binding provided by parent: %s = upval:%d\n", vb.name.Name(), pos)
	//}

	fn.code = c.Export()
	fn.code.totalUpvals = maxArgUpvals
	fn.upvals = make([]*NamedPair, maxArgUpvals)

	cl := &fnClosure{
		importedVars: importBinds,
	}
	return cl, nil
}
