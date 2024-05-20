package core

import "slices"

type varBind struct {
	index int

	name Symbol
	home *fnFrame

	upval bool
	upidx int

	uses []*Instruction
}

type varFrame struct {
	parent    *varFrame
	args      []Symbol
	bindings  map[string]*varBind
	recurDest int

	arity int
	start int

	fn *fnFrame
}

type fnFrame struct {
	parent *fnFrame
	top    *varFrame
	root   *varFrame

	importUpvals []*varBind
	letUpvals    []*varBind

	arity int

	totalBindings int
	curBindings   int
}

func (f *fnFrame) addBindings(cnt int) int {
	cur := f.curBindings

	f.curBindings += cnt
	if f.curBindings > f.totalBindings {
		f.totalBindings = f.curBindings
	}

	return cur
}

func (f *fnFrame) removeBindings(cnt int) {
	f.curBindings -= cnt
}

func newFrame(args []Symbol) *fnFrame {
	fr := &fnFrame{
		root: &varFrame{
			bindings:  make(map[string]*varBind),
			args:      args,
			recurDest: -1,
		},
	}

	fr.top = fr.root
	fr.top.fn = fr

	if len(args) > 0 {
		fr.arity = len(args)
		start := fr.addBindings(len(args))

		vf := fr.top

		for i, a := range args {
			vf.bindings[a.Name()] = &varBind{index: start + i, name: a, home: fr}
		}
	}

	return fr
}

func (f *fnFrame) childFrame(args []Symbol) *fnFrame {
	chfr := newFrame(args)
	chfr.parent = f

	return chfr
}

func (f *fnFrame) bindingFrame(size int) *varFrame {
	start := f.addBindings(size)

	bf := &varFrame{
		parent:    f.top,
		bindings:  make(map[string]*varBind),
		arity:     size,
		start:     start,
		fn:        f,
		recurDest: -1,
	}

	f.top = bf

	return bf
}

func (f *fnFrame) popBindingFrame() {
	f.top.processUpvals()

	f.removeBindings(f.top.arity)

	f.top = f.top.parent
}

func (b *fnFrame) set(sym Symbol) *varBind {
	return b.top.set(sym)
}

func (b *varFrame) set(sym Symbol) *varBind {
	vb := &varBind{
		index: b.start,
		name:  sym,
	}

	b.start++

	b.bindings[sym.Name()] = vb

	return vb
}

func (f *fnFrame) lookup(sym Symbol) *varBind {
	for vf := f.top; vf != nil; vf = vf.parent {
		if vb, ok := vf.bindings[sym.Name()]; ok {
			return vb
		}
	}

	// check parents and perform upval marking

	for fn := f.parent; fn != nil; fn = fn.parent {
		for vf := fn.top; vf != nil; vf = vf.parent {
			if vb, ok := vf.bindings[sym.Name()]; ok {
				if !vb.upval {
					vb.upval = true
					vb.upidx = -1
				}

				chvb := &varBind{
					upval: true,
					name:  sym,
					upidx: -1,
					home:  f,
				}

				f.importUpvals = append(f.importUpvals, chvb)

				// We install the binding in the root varFrame for the function
				// because the binding is common for all binding frames in the function.
				f.root.bindings[sym.Name()] = chvb
				return chvb
			}
		}
	}

	return nil
}

func (f *fnFrame) createUnknownUpval(name Symbol) *varBind {
	chvb := &varBind{
		upval: true,
		name:  name,
		upidx: -1,
		home:  f,
	}

	f.importUpvals = append(f.importUpvals, chvb)

	// We install the binding in the root varFrame for the function
	// because the binding is common for all binding frames in the function.
	f.root.bindings[name.Name()] = chvb
	return chvb
}

func (f *varFrame) processUpvals() {
	for _, vb := range f.bindings {
		if vb.upval {
			f.fn.letUpvals = append(f.fn.letUpvals, vb)
		}
	}
}

func (f *fnFrame) updateUsers(vb *varBind) {
	for _, i := range vb.uses {
		switch i.Op {
		case GetLocal:
			i.Op = GetUpval
			i.A0 = int32(vb.upidx)
		case GetUpval:
			i.A0 = int32(vb.upidx)
		case SetUpval:
			i.A0 = int32(vb.upidx)
		case SetLocal:
			i.Op = SetUpval
			i.A0 = int32(vb.upidx)
		case RefUpval:
			i.A0 = int32(vb.upidx)
		}
	}
}

func (f *fnFrame) assignUpvals() int {
	upidx := 0

	for _, vb := range f.importUpvals {
		vb.upidx = upidx
		f.updateUsers(vb)
		upidx++
	}

	for _, a := range f.top.args {
		vb := f.root.bindings[a.Name()]
		if vb.upval {
			vb.upidx = upidx
			f.updateUsers(vb)
			upidx++
		}
	}

	for _, vb := range f.letUpvals {
		vb.upidx = upidx
		f.updateUsers(vb)
		upidx++
	}

	return upidx
}

func (f *fnFrame) importedNames() []Symbol {
	var ret []Symbol

	for _, vb := range f.importUpvals {
		ret = append(ret, vb.name)
	}

	return ret
}

func (f *fnFrame) closeFrame() []Symbol {
	return slices.Clone(f.importedNames())
}

func (v *varBind) u(i *Instruction) *Instruction {
	v.uses = append(v.uses, i)
	return i
}

func (v *varBind) read() *Instruction {
	if v.upval {
		return v.u(&Instruction{
			Op: GetUpval,
			A0: -1,
		})
	}

	return v.u(&Instruction{
		Op: GetLocal,
		A0: int32(v.index),
	})
}

func (v *varBind) set() *Instruction {
	if v.upval {
		return v.u(&Instruction{
			Op: SetUpval,
			A0: -1,
		})
	}

	return v.u(&Instruction{
		Op: SetLocal,
		A0: int32(v.index),
	})
}

func (v *varBind) refUpval() *Instruction {
	v.upval = true
	return v.u(&Instruction{
		Op: RefUpval,
		A0: -1,
	})
}
