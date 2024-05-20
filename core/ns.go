package core

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
)

type (
	Namespace struct {
		MetaHolder
		Name           Symbol
		Lazy           func(env *Env, ns *Namespace)
		mu             sync.Mutex
		mappings       map[string]*Var
		aliases        map[string]*Namespace
		isUsed         bool
		isGloballyUsed bool
		hash           uint32
		core           bool
	}
)

func (env *Env) AllNamespaces() []string {
	env.mu.Lock()
	defer env.mu.Unlock()

	var names []string

	for k := range env.Namespaces {
		names = append(names, k)
	}

	sort.Strings(names)

	return names
}

func (env *Env) AllNamespaceValues() []any {
	env.mu.Lock()
	defer env.mu.Unlock()

	var vals []any

	for _, v := range env.Namespaces {
		vals = append(vals, v)
	}

	return vals
}

func (ns *Namespace) ToString(env *Env, escape bool) (string, error) {
	return ToString(env, ns.Name)
}

func (ns *Namespace) Qual() string {
	return ns.Name.String()
}

func (ns *Namespace) Print(w io.Writer, printReadably bool) {
	fmt.Fprint(w, "#object[Namespace \""+ns.Name.String()+"\"]")
}

func (ns *Namespace) Equals(env *Env, other interface{}) bool {
	return ns == other
}

func (ns *Namespace) GetInfo() *ObjectInfo {
	return nil
}

func (ns *Namespace) WithInfo(info *ObjectInfo) any {
	return ns
}

func (ns *Namespace) GetType() *Type {
	return TYPE.Namespace
}

func (ns *Namespace) WithMeta(env *Env, meta Map) (any, error) {
	res := &Namespace{
		Name:     ns.Name,
		mappings: make(map[string]*Var),
		aliases:  make(map[string]*Namespace),
	}

	ns.mu.Lock()
	defer ns.mu.Unlock()

	for k, v := range ns.mappings {
		res.mappings[k] = v
	}

	for k, v := range ns.aliases {
		res.aliases[k] = v
	}

	v, err := SafeMerge(env, res.meta, meta)
	if err != nil {
		return nil, err
	}

	ns.meta = v
	return res, nil
}

func (ns *Namespace) ResetMeta(newMeta Map) Map {
	ns.meta = newMeta
	return ns.meta
}

func (ns *Namespace) AlterMeta(env *Env, fn *Fn, args []any) (Map, error) {
	return AlterMeta(env, &ns.MetaHolder, fn, args)
}

func (ns *Namespace) Hash(env *Env) (uint32, error) {
	return ns.hash, nil
}

func (ns *Namespace) MaybeLazy(env *Env, doc string) {
	if ns.Lazy != nil {
		lazyFn := ns.Lazy
		ns.Lazy = nil
		lazyFn(env, ns)
		if VerbosityLevel > 0 {
			fmt.Fprintf(Stderr, "NamespaceFor: Lazily initialized %s for %s\n", ns.Name.Name(), doc)
		}
	}
}

func (ns *Namespace) CoreP() bool {
	return ns.core
}

const nsHashMask uint32 = 0x90569f6f

func NewNamespace(env *Env, sym Symbol) (*Namespace, error) {
	h, err := HashValue(env, sym)
	if err != nil {
		return nil, err
	}

	return &Namespace{
		Name:     sym,
		mappings: make(map[string]*Var),
		aliases:  make(map[string]*Namespace),
		hash:     h ^ nsHashMask,
	}, nil
}

func (ns *Namespace) Refer(env *Env, sym Symbol, vr *Var) (*Var, error) {
	if sym.Namespace() != "" {
		return nil, env.NewError("Can't intern namespace-qualified symbol " + sym.String())
	}
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.mappings[sym.Name()] = vr
	return vr, nil
}

func (ns *Namespace) ReferAll(other *Namespace, safe bool) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	other.mu.Lock()
	defer other.mu.Unlock()

	for name, vr := range other.mappings {
		if !vr.isPrivate {
			if !safe || ns.mappings[name] == nil {
				ns.mappings[name] = vr
			}
		}
	}
}

func (ns *Namespace) Intern(env *Env, sym Symbol) (*Var, error) {
	if sym.Namespace() != "" {
		return nil, StubNewError("Can't intern namespace-qualified symbol " + sym.String())
	}

	ns.mu.Lock()
	defer ns.mu.Unlock()

	ClearMeta(sym)
	existingVar, ok := ns.mappings[sym.Name()]
	if !ok {
		newVar := &Var{
			ns:   ns,
			name: sym,
		}
		ns.mappings[sym.Name()] = newVar
		return newVar, nil
	}
	if existingVar.ns != ns {
		if Equals(env, existingVar.ns.Name, criticalSymbols.lace_core) {
			newVar := &Var{
				ns:   ns,
				name: sym,
			}
			ns.mappings[sym.Name()] = newVar
			if !strings.HasPrefix(ns.Name.Name(), "lace.") {
				printParseWarning(GetInfo(sym).Pos(), fmt.Sprintf("WARNING: %s already refers to: %s in namespace %s, being replaced by: %s\n",
					sym.String(), existingVar.String(), ns.Name.String(), newVar.String()))
			}
			return newVar, nil
		}
		return nil, env.NewError(fmt.Sprintf("WARNING: %s already refers to: %s in namespace %s",
			sym.String(), existingVar.String(), ns.Qual()))
	}
	return existingVar, nil
}

func (ns *Namespace) InternVar(env *Env, name string, val any, meta *ArrayMap) (*Var, error) {
	vr, err := ns.Intern(env, MakeSymbol(name))
	if err != nil {
		return nil, err
	}
	vr.SetStatic(val)
	if meta == nil {
		meta = &ArrayMap{}
	}
	meta.Add(env, criticalKeywords.ns, ns)
	meta.Add(env, criticalKeywords.name, vr.name)
	vr.meta = meta
	return vr, nil
}

func (ns *Namespace) AddAlias(env *Env, alias Symbol, namespace *Namespace) error {
	if alias.Namespace() != "" {
		return env.NewError("Alias can't be namespace-qualified")
	}
	existing := ns.aliases[alias.Name()]
	if existing != nil && existing != namespace {
		msg := "Alias " + alias.String() + " already exists in namespace " + ns.Name.String() + ", aliasing " + existing.Name.String()
		if LINTER_MODE {
			printParseError(GetPosition(alias), msg)
			return nil
		}
		return env.NewError(msg)
	}
	ns.aliases[alias.Name()] = namespace
	return nil
}

func (ns *Namespace) Resolve(name string) *Var {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	return ns.mappings[name]
}

func (ns *Namespace) Mappings() map[string]*Var {
	return ns.mappings
}

func (ns *Namespace) AliasNames() []string {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	var ret []string
	for k := range ns.aliases {
		ret = append(ret, k)
	}

	return ret
}

func (ns *Namespace) LookupVar(name string) (*Var, bool) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	v, ok := ns.mappings[name]
	return v, ok
}

func (ns *Namespace) DeleteVar(name string) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	delete(ns.mappings, name)
}

func (ns *Namespace) VarNames() []string {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	out := make([]string, 0, len(ns.mappings))

	for k := range ns.mappings {
		out = append(out, k)
	}

	sort.Strings(out)

	return out
}

func (ns *Namespace) MappingsAsMap(env *Env) Map {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	r := &ArrayMap{}

	for k, v := range ns.mappings {
		r.Add(env, MakeSymbol(k), v)
	}

	return r
}

func WarnOnUnusedVars(env *Env) {
	var names []string
	positions := make(map[string]Position)

	for _, ns := range env.Namespaces {
		if ns == env.CoreNamespace {
			continue
		}
		ns.mu.Lock()
		for _, vr := range ns.mappings {
			if vr.ns == ns && !vr.isUsed && vr.isPrivate {
				pos := vr.GetInfo()
				if pos != nil {
					names = append(names, vr.name.Name())
					positions[vr.name.Name()] = pos.Position
				}
			}
		}
		ns.mu.Lock()
	}

	sort.Strings(names)
	for _, name := range names {
		printParseWarning(positions[name], "unused var "+name)
	}
}

func WarnOnGloballyUnusedVars(env *Env) {
	var names []string
	positions := make(map[string]Position)

	for _, ns := range env.Namespaces {
		if ns == env.CoreNamespace {
			continue
		}
		ns.mu.Lock()
		for _, vr := range ns.mappings {
			if vr.ns == ns && !vr.isGloballyUsed && !vr.isPrivate && !isRecordConstructor(vr.name) && !isEntryPointVar(vr) {
				pos := vr.GetInfo()
				if pos != nil {
					varName := vr.Name()
					names = append(names, varName)
					positions[varName] = pos.Position
				}
			}
		}
		ns.mu.Unlock()
	}

	sort.Strings(names)
	for _, name := range names {
		printParseWarning(positions[name], "globally unused var "+name)
	}
}

func ResetUsage(env *Env) {
	for _, ns := range env.Namespaces {
		if ns == env.CoreNamespace {
			continue
		}
		ns.isUsed = true
		ns.mu.Lock()
		for _, vr := range ns.mappings {
			vr.isUsed = true
		}
		ns.mu.Unlock()
	}
}
