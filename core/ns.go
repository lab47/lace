package core

import (
	"fmt"
	"io"
	"strings"
)

type (
	Namespace struct {
		MetaHolder
		Name           Symbol
		Lazy           func(env *Env, ns *Namespace)
		mappings       map[*string]*Var
		aliases        map[*string]*Namespace
		isUsed         bool
		isGloballyUsed bool
		hash           uint32
		core           bool
	}
)

func (ns *Namespace) ToString(escape bool) string {
	return ns.Name.ToString(escape)
}

func (ns *Namespace) Print(w io.Writer, printReadably bool) {
	fmt.Fprint(w, "#object[Namespace \""+ns.Name.ToString(true)+"\"]")
}

func (ns *Namespace) Equals(other interface{}) bool {
	return ns == other
}

func (ns *Namespace) GetInfo() *ObjectInfo {
	return nil
}

func (ns *Namespace) WithInfo(info *ObjectInfo) Object {
	return ns
}

func (ns *Namespace) GetType() *Type {
	return TYPE.Namespace
}

func (ns *Namespace) WithMeta(meta Map) (Object, error) {
	res := *ns
	v, err := SafeMerge(res.meta, meta)
	if err != nil {
		return nil, err
	}

	ns.meta = v
	return &res, nil
}

func (ns *Namespace) ResetMeta(newMeta Map) Map {
	ns.meta = newMeta
	return ns.meta
}

func (ns *Namespace) AlterMeta(env *Env, fn *Fn, args []Object) (Map, error) {
	return AlterMeta(env, &ns.MetaHolder, fn, args)
}

func (ns *Namespace) Hash() uint32 {
	return ns.hash
}

func (ns *Namespace) MaybeLazy(env *Env, doc string) {
	if ns.Lazy != nil {
		lazyFn := ns.Lazy
		ns.Lazy = nil
		lazyFn(env, ns)
		if VerbosityLevel > 0 {
			fmt.Fprintf(Stderr, "NamespaceFor: Lazily initialized %s for %s\n", *ns.Name.name, doc)
		}
	}
}

func (ns *Namespace) CoreP() bool {
	return ns.core
}

const nsHashMask uint32 = 0x90569f6f

func NewNamespace(sym Symbol) *Namespace {
	return &Namespace{
		Name:     sym,
		mappings: make(map[*string]*Var),
		aliases:  make(map[*string]*Namespace),
		hash:     sym.Hash() ^ nsHashMask,
	}
}

func (ns *Namespace) Refer(env *Env, sym Symbol, vr *Var) (*Var, error) {
	if sym.ns != nil {
		return nil, env.RT.NewError("Can't intern namespace-qualified symbol " + sym.ToString(false))
	}
	ns.mappings[sym.name] = vr
	return vr, nil
}

func (ns *Namespace) ReferAll(other *Namespace) {
	for name, vr := range other.mappings {
		if !vr.isPrivate {
			ns.mappings[name] = vr
		}
	}
}

func (ns *Namespace) Intern(env *Env, sym Symbol) (*Var, error) {
	if sym.ns != nil {
		return nil, StubNewError("Can't intern namespace-qualified symbol " + sym.ToString(false))
	}
	sym.meta = nil
	existingVar, ok := ns.mappings[sym.name]
	if !ok {
		newVar := &Var{
			ns:   ns,
			name: sym,
		}
		ns.mappings[sym.name] = newVar
		return newVar, nil
	}
	if existingVar.ns != ns {
		if existingVar.ns.Name.Equals(criticalSymbols.lace_core) {
			newVar := &Var{
				ns:   ns,
				name: sym,
			}
			ns.mappings[sym.name] = newVar
			if !strings.HasPrefix(ns.Name.Name(), "lace.") {
				printParseWarning(sym.GetInfo().Pos(), fmt.Sprintf("WARNING: %s already refers to: %s in namespace %s, being replaced by: %s\n",
					sym.ToString(false), existingVar.ToString(false), ns.Name.ToString(false), newVar.ToString(false)))
			}
			return newVar, nil
		}
		return nil, StubNewError(fmt.Sprintf("WARNING: %s already refers to: %s in namespace %s",
			sym.ToString(false), existingVar.ToString(false), ns.ToString(false)))
	}
	if LINTER_MODE && existingVar.expr != nil && !existingVar.ns.Name.Equals(criticalSymbols.lace_core) {
		printParseWarning(sym.GetInfo().Pos(), "Duplicate def of "+existingVar.ToString(false))
	}
	return existingVar, nil
}

func (ns *Namespace) InternVar(env *Env, name string, val Object, meta *ArrayMap) (*Var, error) {
	vr, err := ns.Intern(env, MakeSymbol(name))
	if err != nil {
		return nil, err
	}
	vr.Value = val
	meta.Add(criticalKeywords.ns, ns)
	meta.Add(criticalKeywords.name, vr.name)
	vr.meta = meta
	return vr, nil
}

func (ns *Namespace) AddAlias(env *Env, alias Symbol, namespace *Namespace) error {
	if alias.ns != nil {
		return env.RT.NewError("Alias can't be namespace-qualified")
	}
	existing := ns.aliases[alias.name]
	if existing != nil && existing != namespace {
		msg := "Alias " + alias.ToString(false) + " already exists in namespace " + ns.Name.ToString(false) + ", aliasing " + existing.Name.ToString(false)
		if LINTER_MODE {
			printParseError(GetPosition(alias), msg)
			return nil
		}
		return env.RT.NewError(msg)
	}
	ns.aliases[alias.name] = namespace
	return nil
}

func (ns *Namespace) Resolve(name string) *Var {
	return ns.mappings[STRINGS.Intern(name)]
}

func (ns *Namespace) Mappings() map[*string]*Var {
	return ns.mappings
}

func (ns *Namespace) Aliases() map[*string]*Var {
	return ns.mappings
}
