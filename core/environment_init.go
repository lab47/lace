//go:build !fast_init
// +build !fast_init

package core

import (
	"context"
	"fmt"
)

func NewEnv() (*Env, error) {
	features := EmptySet()
	res := &Env{
		State: &State{
			Namespaces: make(map[string]*Namespace),
			Features:   features,
		},
		Context: context.Background(),
	}
	_, err := features.Add(res, MakeKeyword("default"))
	if err != nil {
		return nil, err
	}
	_, err = features.Add(res, MakeKeyword("lace"))
	if err != nil {
		return nil, err
	}

	coreNs, err := NewNamespace(res, criticalSymbols.lace_core)
	if err != nil {
		return nil, err
	}

	res.Namespaces[criticalSymbols.lace_core.Name()] = coreNs

	res.CoreNamespace = coreNs
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

	vr, err := res.CoreNamespace.Intern(res, MakeSymbol("*context*"))
	if err != nil {
		return nil, err
	}

	vr.isDynamic = true
	vr.SetStatic(MakeReflectValue(res.Context))

	res.ctx = vr

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
	res.classPath.SetStatic(NIL)
	res.classPath.isPrivate = true
	res.printReadably, err = res.CoreNamespace.Intern(res, MakeSymbol("*print-readably*"))
	if err != nil {
		return nil, err
	}
	res.printReadably.SetStatic(Boolean(true))
	_, err = res.CoreNamespace.InternVar(res, "*linter-mode*", Boolean(LINTER_MODE),
		MakeMeta(nil, "true if Lace is running in linter mode", "1.0"))
	if err != nil {
		return nil, err
	}
	_, err = res.CoreNamespace.InternVar(res, "*linter-config*", EmptyArrayMap(),
		MakeMeta(nil, "Map of configuration key/value pairs for linter mode", "1.0"))
	if err != nil {
		return nil, err
	}

	res.classPath.SetStatic(NewVectorFrom())

	userNs := res.EnsureNamespace(MakeSymbol("user"))

	res.SetCurrentNamespace(userNs)

	err = createCoreFns(res)
	if err != nil {
		return nil, err
	}

	reflectBuilder, err := SetupPkgReflect(res)
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

	// this happens when doing gen_data currently.
	if res.CoreNamespace.Resolve("defmacro") != nil {
		reflectBuilder.Run(reflectCode)
	}

	userNs.ReferAll(res.CoreNamespace, true)

	return res, nil
}
