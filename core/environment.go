package core

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var (
	Stdin          io.Reader = os.Stdin
	Stdout         io.Writer = os.Stdout
	Stderr         io.Writer = os.Stderr
	VerbosityLevel           = 0
)

type (
	State struct {
		mu            sync.Mutex
		Namespaces    map[string]*Namespace
		CoreNamespace *Namespace
		LangNamespace *Namespace
		stdout        *Var
		stdin         *Var
		stderr        *Var
		printReadably *Var
		file          *Var
		MainFile      *Var
		args          *Var
		classPath     *Var
		ns            *Var
		NS_VAR        *Var
		IN_NS_VAR     *Var
		version       *Var
		ctx           *Var
		Features      Set
	}

	Env struct {
		*State

		Parent *Env

		CurrentVar     Associative
		Context        context.Context
		cycleDetection map[[2]any]struct{}

		Engine *Engine

		DebugBytecode bool

		treeEvalStack []Expr
	}
)

func (env *Env) Child() *Env {
	ret := &Env{
		State:      env.State,
		Parent:     env,
		Engine:     NewEngine(),
		CurrentVar: NIL,
	}

	err := ret.SetContext(context.Background())
	if err != nil {
		panic(err)
	}

	return ret
}

func (env *Env) enableCycleDetection() func() {
	if env.cycleDetection != nil {
		return func() {}
	}

	env.cycleDetection = make(map[[2]any]struct{})
	return func() {
		env.cycleDetection = nil
	}
}

func (env *Env) cycling(a, b any) bool {
	if env.cycleDetection == nil {
		return false
	}

	key := [2]any{a, b}

	_, ok := env.cycleDetection[key]
	if ok {
		return true
	}

	env.cycleDetection[key] = struct{}{}

	return false
}

func (env *Env) FindInCurrentVars(vr *Var) (any, bool, error) {
	if env.CurrentVar == nil {
		return nil, false, nil
	}

	found, val, err := env.CurrentVar.Get(env, vr)
	if err != nil {
		return nil, false, err
	}

	return val, found, nil
}

func versionMap(env *Env) Map {
	res := EmptyArrayMap()
	parts := strings.Split(VERSION[1:], ".")
	i, _ := strconv.ParseInt(parts[0], 10, 64)
	res.Add(env, MakeKeyword("major"), MakeInt(int(i)))
	i, _ = strconv.ParseInt(parts[1], 10, 64)
	res.Add(env, MakeKeyword("minor"), MakeInt(int(i)))
	i, _ = strconv.ParseInt(parts[2], 10, 64)
	res.Add(env, MakeKeyword("incremental"), MakeInt(int(i)))
	return res
}

func (env *Env) SetEnvArgs(newArgs []string) {
	args := EmptyVector()
	for _, arg := range newArgs {
		args, _ = args.Conjoin(MakeString(arg))
	}
	if args.Count() > 0 {
		env.args.SetStatic(args.Seq())
	} else {
		env.args.SetStatic(NIL)
	}
}

/*
This runs after invariant initialization, which includes calling

	NewEnv().
*/
func (env *Env) SetClassPath(cp string) {
	cpArray := filepath.SplitList(cp)
	cpVec := EmptyVector()
	for _, cpelem := range cpArray {
		cpVec, _ = cpVec.Conjoin(MakeString(cpelem))
	}
	if cpVec.Count() == 0 {
		cpVec, _ = cpVec.Conjoin(MakeString(""))
	}
	env.classPath.SetStatic(cpVec)
}

/*
This runs after invariant initialization, which includes calling

	NewEnv().
*/
func (env *Env) InitEnv(stdin io.Reader, stdout, stderr io.Writer, args []string) {
	env.stdin.SetStatic(MakeBufferedReader(stdin))
	env.stdout.SetStatic(MakeIOWriter(stdout))
	env.stderr.SetStatic(MakeIOWriter(stderr))
	env.SetEnvArgs(args)
}

func (env *Env) SetStdIO(stdin, stdout, stderr any) {
	env.stdin.SetStatic(stdin)
	env.stdout.SetStatic(stdout)
	env.stderr.SetStatic(stderr)
}

func (env *Env) StdIO() (stdin, stdout, stderr any) {
	return env.stdin.GetStatic(), env.stdout.GetStatic(), env.stderr.GetStatic()
}

/*
This runs after invariant initialization, which includes calling

	NewEnv().
*/
func (env *Env) SetMainFilename(filename string) {
	env.MainFile.SetStatic(MakeString(filename))
}

/*
This runs after invariant initialization, which includes calling

	NewEnv().
*/
func (env *Env) SetFilename(obj any) {
	env.file.SetStatic(obj)
}

func (env *Env) IsStdIn(obj any) bool {
	return env.stdin.GetStatic() == obj
}

func (env *Env) CurrentNamespace() *Namespace {
	ns, err := AssertNamespace(env, env.ns.GetStatic(), "")
	if err != nil {
		panic(err) // this is extremely rare, we should probably make it not possible
	}

	return ns
}

func (env *Env) SetCurrentNamespace(ns *Namespace) {
	env.ns.SetStatic(ns)
}

func (env *Env) ProtoNamespace(sym Symbol) (*Namespace, error) {
	if sym.Namespace() != "" {
		return nil, env.NewError("Namespace's name cannot be qualified: " + sym.String())
	}

	env.mu.Lock()
	ns := env.Namespaces[sym.Name()]
	env.mu.Unlock()

	var err error

	if ns != nil {
		return ns, nil
	}

	ns, err = NewNamespace(env, sym)
	if err != nil {
		return nil, err
	}

	env.mu.Lock()
	env.Namespaces[sym.Name()] = ns
	env.mu.Unlock()

	return ns, nil
}

func (env *Env) SetupNamespace(ns *Namespace) error {
	sym := ns.Name

	if _, ok := ns.LookupVar("*loaded-libs*"); !ok {
		ns.ReferAll(env.CoreNamespace, true)
	}

	if setup, ok := builtinNSSetup[sym.Name()]; ok {
		err := setup(env)
		if err != nil {
			return errors.Wrapf(err, "loading builtin ns: %s", sym.Name())
		}
	} else {
		_, err := PopulateNativeNamespaceToEnv(env, sym.Name())
		if err != nil {
			return errors.Wrapf(err, "loading builtin ns: %s", sym.Name())
		}
	}

	return nil
}

func (env *Env) EnsureNamespace(sym Symbol) *Namespace {
	ns, err := env.InitNamespace(sym)
	if err != nil {
		panic(WrapError(env, err))
	}

	return ns
}

func (env *Env) InitNamespace(sym Symbol) (*Namespace, error) {
	if sym.Namespace() != "" {
		return nil, env.NewError("Namespace's name cannot be qualified: " + sym.String())
	}

	env.mu.Lock()
	ns := env.Namespaces[sym.Name()]
	env.mu.Unlock()

	var err error

	if ns == nil {
		ns, err = NewNamespace(env, sym)
		if err != nil {
			return nil, err
		}

		env.mu.Lock()
		env.Namespaces[sym.Name()] = ns
		env.mu.Unlock()

		err := env.SetupNamespace(ns)
		if err != nil {
			return nil, err
		}
	}

	return ns, nil
}

func (env *Env) NamespaceFor(ns *Namespace, s Symbol) *Namespace {
	var res *Namespace
	if s.Namespace() == "" {
		res = ns
	} else {
		ns.mu.Lock()
		res = ns.aliases[s.Namespace()]
		ns.mu.Unlock()

		if res == nil {
			env.mu.Lock()
			res = env.Namespaces[s.Namespace()]
			env.mu.Unlock()
		}
	}
	if res != nil {
		res.MaybeLazy(env, "NamespaceFor")
	}
	return res
}

func (env *Env) ResolveIn(n *Namespace, s Symbol) (*Var, bool) {
	ns := env.NamespaceFor(n, s)
	if ns == nil {
		return nil, false
	}
	if v, ok := ns.LookupVar(s.Name()); ok {
		return v, true
	}
	if s.Is(env.IN_NS_VAR.name) {
		return env.IN_NS_VAR, true
	}
	if s.Is(env.NS_VAR.name) {
		return env.NS_VAR, true
	}
	return nil, false
}

func (env *Env) Resolve(s Symbol) (*Var, bool) {
	ns := env.CurrentNamespace()
	return env.ResolveIn(ns, s)
}

func (env *Env) MakeVar(s Symbol) (*Var, error) {
	ns := env.CurrentNamespace()
	return ns.Intern(env, s)
}

func (env *Env) LookupNamespace(s Symbol) (*Namespace, bool) {
	if s.Namespace() != "" {
		return nil, false
	}

	env.mu.Lock()
	ns := env.Namespaces[s.Name()]
	env.mu.Unlock()

	if ns != nil {
		ns.MaybeLazy(env, "FindNameSpace")
	} else {
		if _, ok := builtinNSSetup[s.Name()]; ok {
			// don't call setup! just create the namespace because EnsureNamespace will call
			// the setup.
			ns = env.EnsureNamespace(s)
		} else {
			_, err := PopulateNativeNamespaceToEnv(env, s.Name())
			if err != nil {
				panic(WrapError(env, err))
			}
		}
	}

	return ns, ns == nil
}

func (env *Env) FindNamespace(s Symbol) *Namespace {
	ns, _ := env.LookupNamespace(s)
	return ns
}

func (env *Env) RemoveNamespace(s Symbol) *Namespace {
	if s.Namespace() != "" {
		return nil
	}

	env.mu.Lock()
	defer env.mu.Unlock()

	if s.Is(criticalSymbols.lace_core) {
		panic(env.NewError("Cannot remove core namespace"))
	}
	ns := env.Namespaces[s.Name()]
	delete(env.Namespaces, s.Name())
	return ns
}

func (env *Env) ResolveSymbol(s Symbol) (Symbol, error) {
	if strings.ContainsRune(s.Name(), '.') {
		return s, nil
	}
	if s.Namespace() == "" && env.ResolveType(s) != nil {
		return s, nil
	}
	currentNs := env.CurrentNamespace()

	if s.Namespace() != "" {
		ns := env.NamespaceFor(currentNs, s)
		if ns == nil || ns.Name.Name() == s.Namespace() {
			if ns != nil {
				ns.isUsed = true
				ns.isGloballyUsed = true
			}
			return s, nil
		}
		ns.isUsed = true
		ns.isGloballyUsed = true
		return AssembleSymbol(ns.Name.Name(), s.Name()), nil
	}
	vr, ok := currentNs.LookupVar(s.Name())
	if !ok {
		return AssembleSymbol(currentNs.Name.Name(), s.Name()), nil
	}
	vr.isUsed = true
	vr.isGloballyUsed = true
	vr.ns.isUsed = true
	vr.ns.isGloballyUsed = true
	return AssembleSymbol(vr.ns.Name.Name(), vr.name.Name()), nil
}

func (env *Env) Eval(str string) (any, error) {
	reader := NewReader(strings.NewReader(str), "<expr>")
	return ProcessReader(env, reader, "")
}

func (e *Env) SetContext(ctx context.Context) error {
	e.Context = ctx
	if v := e.ctx; v != nil {
		return v.SetValue(e, ctx)
	}
	return nil
}
