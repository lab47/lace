package core

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	Stdin          io.Reader = os.Stdin
	Stdout         io.Writer = os.Stdout
	Stderr         io.Writer = os.Stderr
	VerbosityLevel           = 0
)

type (
	Env struct {
		Namespaces    map[*string]*Namespace
		CoreNamespace *Namespace
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
		Features      Set

		RT *Runtime
	}
)

func versionMap() Map {
	res := EmptyArrayMap()
	parts := strings.Split(VERSION[1:], ".")
	i, _ := strconv.ParseInt(parts[0], 10, 64)
	res.Add(MakeKeyword("major"), Int{I: int(i)})
	i, _ = strconv.ParseInt(parts[1], 10, 64)
	res.Add(MakeKeyword("minor"), Int{I: int(i)})
	i, _ = strconv.ParseInt(parts[2], 10, 64)
	res.Add(MakeKeyword("incremental"), Int{I: int(i)})
	return res
}

func (env *Env) SetEnvArgs(newArgs []string) {
	args := EmptyVector()
	for _, arg := range newArgs {
		args, _ = args.Conjoin(MakeString(arg))
	}
	if args.Count() > 0 {
		env.args.Value = args.Seq()
	} else {
		env.args.Value = NIL
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
	env.classPath.Value = cpVec
}

/*
This runs after invariant initialization, which includes calling

	NewEnv().
*/
func (env *Env) InitEnv(stdin io.Reader, stdout, stderr io.Writer, args []string) {
	env.stdin.Value = MakeBufferedReader(stdin)
	env.stdout.Value = MakeIOWriter(stdout)
	env.stderr.Value = MakeIOWriter(stderr)
	env.SetEnvArgs(args)
}

func (env *Env) SetStdIO(stdin, stdout, stderr Object) {
	env.stdin.Value = stdin
	env.stdout.Value = stdout
	env.stderr.Value = stderr
}

func (env *Env) StdIO() (stdin, stdout, stderr Object) {
	return env.stdin.Value, env.stdout.Value, env.stderr.Value
}

/*
This runs after invariant initialization, which includes calling

	NewEnv().
*/
func (env *Env) SetMainFilename(filename string) {
	env.MainFile.Value = MakeString(filename)
}

/*
This runs after invariant initialization, which includes calling

	NewEnv().
*/
func (env *Env) SetFilename(obj Object) {
	env.file.Value = obj
}

func (env *Env) IsStdIn(obj Object) bool {
	return env.stdin.Value == obj
}

func (env *Env) CurrentNamespace() *Namespace {
	ns, err := AssertNamespace(env, env.ns.Value, "")
	if err != nil {
		panic(err) // this is extremely rare, we should probably make it not possible
	}

	return ns
}

func (env *Env) SetCurrentNamespace(ns *Namespace) {
	env.ns.Value = ns
}

func (env *Env) EnsureNamespace(sym Symbol) *Namespace {
	if sym.ns != nil {
		panic(env.RT.NewError("Namespace's name cannot be qualified: " + sym.ToString(false)))
	}
	if env.Namespaces[sym.name] == nil {
		env.Namespaces[sym.name] = NewNamespace(sym)
	}
	return env.Namespaces[sym.name]
}

func (env *Env) NamespaceFor(ns *Namespace, s Symbol) *Namespace {
	var res *Namespace
	if s.ns == nil {
		res = ns
	} else {
		res = ns.aliases[s.ns]
		if res == nil {
			res = env.Namespaces[s.ns]
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
	if v, ok := ns.mappings[s.name]; ok {
		return v, true
	}
	if s.Equals(env.IN_NS_VAR.name) {
		return env.IN_NS_VAR, true
	}
	if s.Equals(env.NS_VAR.name) {
		return env.NS_VAR, true
	}
	return nil, false
}

func (env *Env) Resolve(s Symbol) (*Var, bool) {
	ns := env.CurrentNamespace()
	return env.ResolveIn(ns, s)
}

func (env *Env) FindNamespace(s Symbol) *Namespace {
	if s.ns != nil {
		return nil
	}
	ns := env.Namespaces[s.name]
	if ns != nil {
		ns.MaybeLazy(env, "FindNameSpace")
	}
	return ns
}

func (env *Env) RemoveNamespace(s Symbol) *Namespace {
	if s.ns != nil {
		return nil
	}
	if s.Equals(criticalSymbols.lace_core) {
		panic(env.RT.NewError("Cannot remove core namespace"))
	}
	ns := env.Namespaces[s.name]
	delete(env.Namespaces, s.name)
	return ns
}

func (env *Env) ResolveSymbol(s Symbol) (Symbol, error) {
	if strings.ContainsRune(*s.name, '.') {
		return s, nil
	}
	if s.ns == nil && TYPES[s.name] != nil {
		return s, nil
	}
	currentNs := env.CurrentNamespace()

	if s.ns != nil {
		ns := env.NamespaceFor(currentNs, s)
		if ns == nil || ns.Name.name == s.ns {
			if ns != nil {
				ns.isUsed = true
				ns.isGloballyUsed = true
			}
			return s, nil
		}
		ns.isUsed = true
		ns.isGloballyUsed = true
		return Symbol{
			name: s.name,
			ns:   ns.Name.name,
		}, nil
	}
	vr, ok := currentNs.mappings[s.name]
	if !ok {
		return Symbol{
			name: s.name,
			ns:   currentNs.Name.name,
		}, nil
	}
	vr.isUsed = true
	vr.isGloballyUsed = true
	vr.ns.isUsed = true
	vr.ns.isGloballyUsed = true
	return Symbol{
		name: vr.name.name,
		ns:   vr.ns.Name.name,
	}, nil
}
