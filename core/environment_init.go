//go:build !fast_init
// +build !fast_init

package core

func NewEnv() *Env {
	features := EmptySet()
	features.Add(MakeKeyword("default"))
	features.Add(MakeKeyword("lace"))
	res := &Env{
		Namespaces: make(map[*string]*Namespace),
		Features:   features,
	}
	res.CoreNamespace = res.EnsureNamespace(SYMBOLS.lace_core)
	res.CoreNamespace.core = true
	res.CoreNamespace.meta = MakeMeta(nil, "Core library of Joker.", "1.0")
	res.NS_VAR = res.CoreNamespace.Intern(MakeSymbol("ns"))
	res.IN_NS_VAR = res.CoreNamespace.Intern(MakeSymbol("in-ns"))
	res.ns = res.CoreNamespace.Intern(MakeSymbol("*ns*"))
	res.stdin = res.CoreNamespace.Intern(MakeSymbol("*in*"))
	res.stdout = res.CoreNamespace.Intern(MakeSymbol("*out*"))
	res.stderr = res.CoreNamespace.Intern(MakeSymbol("*err*"))
	res.file = res.CoreNamespace.Intern(MakeSymbol("*file*"))
	res.MainFile = res.CoreNamespace.Intern(MakeSymbol("*main-file*"))
	res.version = res.CoreNamespace.InternVar("*lace-version*", versionMap(),
		MakeMeta(nil, `The version info for Clojure core, as a map containing :major :minor
			:incremental and :qualifier keys. Feature releases may increment
			:minor and/or :major, bugfix releases will increment :incremental.`, "1.0"))
	res.args = res.CoreNamespace.Intern(MakeSymbol("*command-line-args*"))
	res.classPath = res.CoreNamespace.Intern(MakeSymbol("*classpath*"))
	res.classPath.Value = NIL
	res.classPath.isPrivate = true
	res.printReadably = res.CoreNamespace.Intern(MakeSymbol("*print-readably*"))
	res.printReadably.Value = Boolean{B: true}
	res.CoreNamespace.InternVar("*linter-mode*", Boolean{B: LINTER_MODE},
		MakeMeta(nil, "true if Joker is running in linter mode", "1.0"))
	res.CoreNamespace.InternVar("*linter-config*", EmptyArrayMap(),
		MakeMeta(nil, "Map of configuration key/value pairs for linter mode", "1.0"))
	res.SetCurrentNamespace(res.EnsureNamespace(MakeSymbol("user")))
	res.RT = &Runtime{
		callstack: &Callstack{frames: make([]Frame, 0, 50)},
	}

	initEnv(res)

	builtinNS := []string{"core", "repl"}

	for _, name := range builtinNS {
		if fn, ok := builtinNSSetup[name]; ok {
			err := fn(res)
			if err != nil {
				panic(err)
			}
		}
	}
	return res
}

func (env *Env) ReferCoreToUser() {
	env.FindNamespace(MakeSymbol("user")).ReferAll(env.CoreNamespace)
}
