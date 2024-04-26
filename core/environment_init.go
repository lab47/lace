//go:build !fast_init
// +build !fast_init

package core

import "fmt"

func NewEnv() (*Env, error) {
	features := EmptySet()
	res := &Env{
		Namespaces: make(map[string]*Namespace),
		Features:   features,
	}
	_, err := features.Add(res, MakeKeyword("default"))
	if err != nil {
		return nil, err
	}
	_, err = features.Add(res, MakeKeyword("lace"))
	if err != nil {
		return nil, err
	}

	res.CoreNamespace = res.ensureNamespace(criticalSymbols.lace_core)
	res.CoreNamespace.core = true
	res.CoreNamespace.meta = MakeMeta(nil, "Core library of Lace.", "1.0")
	res.NS_VAR, err = res.CoreNamespace.Intern(res, MakeSymbol("ns"))
	if err != nil {
		return nil, err
	}
	res.IN_NS_VAR, err = res.CoreNamespace.Intern(res, MakeSymbol("in-ns"))
	if err != nil {
		return nil, err
	}
	res.ns, err = res.CoreNamespace.Intern(res, MakeSymbol("*ns*"))
	if err != nil {
		return nil, err
	}
	res.stdin, err = res.CoreNamespace.Intern(res, MakeSymbol("*in*"))
	if err != nil {
		return nil, err
	}
	res.stdout, err = res.CoreNamespace.Intern(res, MakeSymbol("*out*"))
	if err != nil {
		return nil, err
	}
	res.stderr, err = res.CoreNamespace.Intern(res, MakeSymbol("*err*"))
	if err != nil {
		return nil, err
	}
	res.file, err = res.CoreNamespace.Intern(res, MakeSymbol("*file*"))
	if err != nil {
		return nil, err
	}
	res.MainFile, err = res.CoreNamespace.Intern(res, MakeSymbol("*main-file*"))
	if err != nil {
		return nil, err
	}
	res.version, err = res.CoreNamespace.InternVar(res, "*lace-version*", versionMap(res),
		MakeMeta(nil, `The version info for Clojure core, as a map containing :major :minor
			:incremental and :qualifier keys. Feature releases may increment
			:minor and/or :major, bugfix releases will increment :incremental.`, "1.0"))
	if err != nil {
		return nil, err
	}
	res.args, err = res.CoreNamespace.Intern(res, MakeSymbol("*command-line-args*"))
	if err != nil {
		return nil, err
	}
	res.classPath, err = res.CoreNamespace.Intern(res, MakeSymbol("*classpath*"))
	if err != nil {
		return nil, err
	}
	res.classPath.Value = NIL
	res.classPath.isPrivate = true
	res.printReadably, err = res.CoreNamespace.Intern(res, MakeSymbol("*print-readably*"))
	if err != nil {
		return nil, err
	}
	res.printReadably.Value = Boolean{B: true}
	_, err = res.CoreNamespace.InternVar(res, "*linter-mode*", Boolean{B: LINTER_MODE},
		MakeMeta(nil, "true if Lace is running in linter mode", "1.0"))
	if err != nil {
		return nil, err
	}
	_, err = res.CoreNamespace.InternVar(res, "*linter-config*", EmptyArrayMap(),
		MakeMeta(nil, "Map of configuration key/value pairs for linter mode", "1.0"))
	if err != nil {
		return nil, err
	}
	res.SetCurrentNamespace(res.EnsureNamespace(MakeSymbol("user")))
	res.RT = &Runtime{
		callstack: &Callstack{frames: make([]Frame, 1, 50)},
	}

	err = initEnv(res)
	if err != nil {
		return nil, err
	}

	// Pull in lace.string because it needs to be defined to pack repl.clj

	res.EnsureNamespace(MakeSymbol("lace.string"))

	err = SetupPkgReflect(res)
	if err != nil {
		return nil, err
	}

	builtinNS := []string{"lace.core"}

	for _, name := range builtinNS {
		if fn, ok := builtinNSSetup[name]; ok {
			err := fn(res)
			if err != nil {
				panic(fmt.Sprintf("error loading %s: %s", name, err))
			}
		}
	}

	return res, nil
}

func (env *Env) ReferCoreToUser() {
	env.FindNamespace(MakeSymbol("user")).ReferAll(env.CoreNamespace)
}
