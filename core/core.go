package core

import "fmt"

func CallVar(env *Env, varName string, args ...Object) (Object, error) {
	sym := MakeSymbol(varName)

	nsName := sym.ns

	var ns *Namespace

	if nsName == "" {
		ns = env.CoreNamespace
	} else {
		ns = env.FindNamespace(MakeSymbol(nsName))
	}

	if ns == nil {
		return nil, fmt.Errorf("unknown namespace: %s", nsName)
	}

	vr, err := ns.Intern(env, MakeSymbol(sym.name))
	if err != nil {
		return nil, err
	}

	callable, ok := vr.Value.(Callable)
	if !ok {
		return nil, fmt.Errorf("var %s is not callable", varName)
	}

	obj, err := callable.Call(env, args)
	if err != nil {
		err = env.populateStackTrace(err)
	}

	return obj, err
}

func Load(env *Env, libname string) (Object, error) {
	return CallVar(env, "lace.core/load", MakeSymbol(libname))
}
